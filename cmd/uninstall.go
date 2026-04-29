package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gszhangwei/open-spdd/internal/ui"
)

// homebrewFormulaName is the canonical fully-qualified Homebrew formula path
// for openspdd. It is intentionally a hard-coded string literal — never a
// configurable value — so the uninstall command can never be redirected to a
// different formula by environment or flag input.
const homebrewFormulaName = "gszhangwei/tools/openspdd"

// InstallMethod classifies how the running openspdd binary was installed.
type InstallMethod string

const (
	MethodHomebrew  InstallMethod = "homebrew"
	MethodGoInstall InstallMethod = "go-install"
	MethodUnknown   InstallMethod = "unknown"
)

// InstallContext is the runtime introspection result used to build an
// uninstall plan. It is pure data; building one performs read-only I/O
// (os.Executable, EvalSymlinks, os.UserConfigDir, optional `go env` shell-out).
type InstallContext struct {
	Method       InstallMethod
	BinaryPath   string // raw os.Executable() result
	ResolvedPath string // after EvalSymlinks (falls back to BinaryPath on error)
	FormulaName  string // populated only for Homebrew, fixed to homebrewFormulaName
	MarkerPath   string // <UserConfigDir>/openspdd/.path-hint-shown (may be empty)
	Reason       string // human-readable classification rationale
}

// StepKind enumerates the kinds of actions a plan step can describe.
type StepKind string

const (
	StepRunCommand StepKind = "run-command"
	StepRemoveFile StepKind = "remove-file"
	StepAdvisory   StepKind = "advisory"
)

// PlanStep is a single action in an uninstall plan. Each step is data only;
// the executor (uninstallExecutor) is the only thing that turns a step into
// real I/O.
type PlanStep struct {
	Kind        StepKind
	Description string
	Command     string   // display form, only for StepRunCommand
	Program     string   // argv[0], only for StepRunCommand
	Args        []string // argv[1:], only for StepRunCommand
	Paths       []string // only for StepRemoveFile
	Optional    bool     // when true, failure does not abort the plan
}

// UninstallPlan groups the steps that will uninstall openspdd, plus a short
// human-readable summary suitable for printing before confirmation.
type UninstallPlan struct {
	Method  InstallMethod
	Steps   []PlanStep
	Summary string
}

// Describe renders the plan as a human-readable, multi-line string. Used by
// the "show before confirm" UX and the --dry-run UX.
func (p UninstallPlan) Describe() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Uninstall plan (%s):\n", p.Method)
	if p.Summary != "" {
		fmt.Fprintf(&b, "  %s\n", p.Summary)
	}
	for i, s := range p.Steps {
		fmt.Fprintf(&b, "  [%d] %s\n", i+1, s.Description)
		switch s.Kind {
		case StepRunCommand:
			fmt.Fprintf(&b, "      $ %s\n", s.Command)
		case StepRemoveFile:
			for _, p := range s.Paths {
				fmt.Fprintf(&b, "      path: %s\n", p)
			}
		case StepAdvisory:
			// description already rendered above
		}
	}
	return b.String()
}

// detectInstallContext probes the running binary's location and produces an
// InstallContext. It performs only read-only I/O and never panics. Callers
// can rely on it being safe to invoke at any point.
func detectInstallContext() InstallContext {
	raw, resolved, err := resolveExecutablePath()
	if err != nil {
		return InstallContext{
			Method: MethodUnknown,
			Reason: "could not determine executable path: " + err.Error(),
		}
	}

	method, reason := classifyByPath(resolved)

	ctx := InstallContext{
		Method:       method,
		BinaryPath:   raw,
		ResolvedPath: resolved,
		MarkerPath:   pathHintMarkerPath(),
		Reason:       reason,
	}
	if method == MethodHomebrew {
		ctx.FormulaName = homebrewFormulaName
	}
	return ctx
}

// classifyByPath inspects a resolved binary path and returns the matching
// install method and a short rationale string. Pure: no I/O, easy to test.
func classifyByPath(resolved string) (InstallMethod, string) {
	if resolved == "" {
		return MethodUnknown, "empty resolved path"
	}
	slashed := filepath.ToSlash(resolved)

	if strings.Contains(slashed, "/Cellar/openspdd/") {
		return MethodHomebrew, "binary resolved under Homebrew Cellar"
	}

	binDir := filepath.ToSlash(filepath.Dir(resolved))
	for _, candidate := range findGoBin() {
		if filepath.ToSlash(filepath.Clean(candidate)) == binDir {
			return MethodGoInstall, "binary located in " + candidate
		}
	}

	return MethodUnknown, "binary path matched neither Homebrew Cellar nor a Go bin directory"
}

// findGoBin returns the candidate directories where `go install` may have
// placed binaries on this system. It tries `go env GOBIN`, then
// `<go env GOPATH>/bin`, then `<HomeDir>/go/bin` as a heuristic fallback.
//
// Tolerates `go` not being on PATH: returns whatever heuristic candidates
// can still be derived (e.g., `~/go/bin`). Returns paths cleaned and
// de-duplicated.
func findGoBin() []string {
	var out []string
	add := func(p string) {
		p = strings.TrimSpace(p)
		if p == "" {
			return
		}
		out = append(out, filepath.Clean(p))
	}

	if _, err := exec.LookPath("go"); err == nil {
		if v, err := exec.Command("go", "env", "GOBIN").Output(); err == nil {
			add(strings.TrimSpace(string(v)))
		}
		if v, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
			gopath := strings.TrimSpace(string(v))
			if gopath != "" {
				add(filepath.Join(gopath, "bin"))
			}
		}
	}

	// Heuristic fallback: the conventional location even when `go` is not
	// installed (e.g., user uninstalled Go but the binary remains).
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		add(filepath.Join(home, "go", "bin"))
	}

	// De-duplicate while preserving order.
	seen := make(map[string]struct{}, len(out))
	deduped := out[:0]
	for _, p := range out {
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		deduped = append(deduped, p)
	}
	return deduped
}

// buildPlan turns an InstallContext into a deterministic UninstallPlan using
// the current runtime OS. Tests should prefer buildPlanForOS for full
// cross-platform coverage.
func buildPlan(ctx InstallContext) UninstallPlan {
	return buildPlanForOS(ctx, runtime.GOOS)
}

// buildPlanForOS is the OS-parameterized form of buildPlan, used directly by
// tests so Windows-specific behavior can be exercised on any host.
func buildPlanForOS(ctx InstallContext, goos string) UninstallPlan {
	plan := UninstallPlan{Method: ctx.Method}

	appendMarkerStep := func() {
		if ctx.MarkerPath == "" {
			return
		}
		plan.Steps = append(plan.Steps, PlanStep{
			Kind:        StepRemoveFile,
			Description: "Remove openspdd's first-run marker file",
			Paths:       []string{ctx.MarkerPath},
			Optional:    true,
		})
	}

	switch ctx.Method {
	case MethodHomebrew:
		formula := ctx.FormulaName
		if formula == "" {
			formula = homebrewFormulaName
		}
		plan.Steps = append(plan.Steps, PlanStep{
			Kind:        StepRunCommand,
			Description: "Run Homebrew to uninstall the formula",
			Program:     "brew",
			Args:        []string{"uninstall", formula},
			Command:     "brew uninstall " + formula,
		})
		appendMarkerStep()
		plan.Summary = fmt.Sprintf(
			"openspdd was installed via Homebrew. Will run %q and clean up state.",
			"brew uninstall "+formula,
		)

	case MethodGoInstall:
		if goos == "windows" {
			plan.Steps = append(plan.Steps, PlanStep{
				Kind: StepAdvisory,
				Description: "On Windows, the running .exe cannot self-delete. " +
					"After this command exits, run: del \"" + ctx.ResolvedPath + "\"",
			})
		} else {
			plan.Steps = append(plan.Steps, PlanStep{
				Kind:        StepRemoveFile,
				Description: "Remove the openspdd binary at " + ctx.ResolvedPath,
				Paths:       []string{ctx.ResolvedPath},
				Optional:    false,
			})
		}
		appendMarkerStep()
		plan.Summary = "openspdd was installed via go install. Will remove the binary and clean up state."

	default: // MethodUnknown
		desc := "Could not classify how openspdd was installed"
		if ctx.ResolvedPath != "" {
			desc += " (resolved path: " + ctx.ResolvedPath + ")"
		}
		desc += ". Refusing to auto-uninstall. Remove the binary manually."
		plan.Steps = append(plan.Steps, PlanStep{
			Kind:        StepAdvisory,
			Description: desc,
		})
		plan.Summary = "openspdd installation method could not be detected — manual removal required."
	}

	return plan
}

// uninstallExecutor performs the side effects described by an UninstallPlan.
// It is the only thing in this package that calls os.Remove or exec.Command
// for the uninstall flow.
type uninstallExecutor struct {
	ui ui.UIRenderer
}

// execCommand is the test seam for subprocess execution. Production callers
// always go through exec.Command; tests substitute a fake implementation that
// returns a *exec.Cmd pointing at a no-op program.
var execCommand = exec.Command

// lookPath is a test seam for PATH lookup. Tests can override it to simulate
// a missing binary without manipulating the real PATH environment variable.
var lookPath = exec.LookPath

// Execute runs each step in the plan in order. Non-optional failures abort
// the plan and return the underlying error. Optional failures are surfaced as
// warnings and execution continues.
func (e *uninstallExecutor) Execute(plan UninstallPlan) error {
	for _, step := range plan.Steps {
		switch step.Kind {
		case StepRunCommand:
			if err := e.runCommandStep(step); err != nil {
				return err
			}
		case StepRemoveFile:
			if err := e.runRemoveFileStep(step); err != nil {
				return err
			}
		case StepAdvisory:
			e.ui.RenderWarning(step.Description)
		}
	}

	if leftover, err := lookPath("openspdd"); err == nil && leftover != "" {
		e.ui.RenderWarning("Note: another openspdd binary is still on PATH at " + leftover)
	}

	return nil
}

func (e *uninstallExecutor) runCommandStep(step PlanStep) error {
	path, err := lookPath(step.Program)
	if err != nil {
		msg := step.Program + " is not on PATH; cannot complete uninstall automatically"
		e.ui.RenderError(msg)
		return errors.New(msg)
	}

	cmd := execCommand(path, step.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		e.ui.RenderError("Step failed: " + step.Description + ": " + err.Error())
		return err
	}
	e.ui.RenderSuccess("Done: " + step.Description)
	return nil
}

func (e *uninstallExecutor) runRemoveFileStep(step PlanStep) error {
	for _, p := range step.Paths {
		if err := os.Remove(p); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				if !step.Optional {
					e.ui.RenderWarning("Already absent: " + p)
				}
				continue
			}
			if step.Optional {
				e.ui.RenderWarning("Could not remove " + p + ": " + err.Error())
				continue
			}
			e.ui.RenderError("Could not remove " + p + ": " + err.Error())
			return err
		}
		e.ui.RenderSuccess("Removed: " + p)
	}
	return nil
}

// ---- Test-only exports ------------------------------------------------------
//
// The helpers below expose unexported symbols to tests in the external
// `cmd_test` package (under tests/cmd/). They are deliberately thin and
// labelled "ForTest" so callers cannot mistake them for production API.

// ClassifyByPathForTest exposes classifyByPath to external tests.
func ClassifyByPathForTest(resolved string) (InstallMethod, string) {
	return classifyByPath(resolved)
}

// BuildPlanForOSForTest exposes buildPlanForOS to external tests.
func BuildPlanForOSForTest(ctx InstallContext, goos string) UninstallPlan {
	return buildPlanForOS(ctx, goos)
}

// SwapExecCommandForTest replaces the subprocess factory used by the
// uninstall executor and returns a function that restores the previous
// implementation. Intended for tests only.
func SwapExecCommandForTest(f func(name string, arg ...string) *exec.Cmd) (restore func()) {
	prev := execCommand
	execCommand = f
	return func() { execCommand = prev }
}

// SwapLookPathForTest replaces the PATH lookup used by the uninstall
// executor and returns a function that restores the previous implementation.
// Intended for tests only.
func SwapLookPathForTest(f func(name string) (string, error)) (restore func()) {
	prev := lookPath
	lookPath = f
	return func() { lookPath = prev }
}

// UninstallExecutor is the test-visible alias of uninstallExecutor. Tests
// construct one via NewUninstallExecutorForTest and call Execute on it.
type UninstallExecutor = uninstallExecutor

// NewUninstallExecutorForTest constructs an UninstallExecutor with the given
// renderer. Intended for tests only.
func NewUninstallExecutorForTest(r ui.UIRenderer) *UninstallExecutor {
	return &UninstallExecutor{ui: r}
}

// ---- Cobra subcommand wiring ------------------------------------------------

var (
	uninstallYesFlag    bool
	uninstallDryRunFlag bool
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall openspdd from this machine",
	Long: `Uninstall openspdd by detecting how it was installed (Homebrew or go install)
and either invoking the matching package manager or removing the binary directly.

By default, the planned actions are displayed and confirmation is required.
Use --dry-run to preview without executing, or --yes to skip the confirmation
prompt. If both --dry-run and --yes are set, --dry-run wins (no execution).

Scope:
  In scope:  the openspdd binary, plus openspdd's own first-run marker file
             under <UserConfigDir>/openspdd/.
  Out of scope: generated SPDD command templates inside your projects
             (.cursor/commands/spdd-*.md and friends), and the Homebrew tap
             gszhangwei/tools (which may host other tools).`,
	Args: cobra.NoArgs,
	Run:  runUninstall,
}

func init() {
	uninstallCmd.Flags().BoolVarP(&uninstallYesFlag, "yes", "y", false, "Skip confirmation prompt")
	uninstallCmd.Flags().BoolVar(&uninstallDryRunFlag, "dry-run", false, "Print the plan without executing it")
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(cmd *cobra.Command, args []string) {
	ctx := detectInstallContext()
	plan := buildPlan(ctx)

	fmt.Println(plan.Describe())

	if ctx.Method == MethodUnknown {
		uiRenderer.RenderError("Cannot auto-uninstall: install method could not be detected")
		os.Exit(1)
	}

	if uninstallDryRunFlag {
		uiRenderer.RenderSuccess("Dry run complete — no changes made.")
		return
	}

	if !uninstallYesFlag {
		if !uiRenderer.Confirm("Proceed with uninstall?") {
			uiRenderer.RenderWarning("Uninstall cancelled")
			return
		}
	}

	executor := &uninstallExecutor{ui: uiRenderer}
	if err := executor.Execute(plan); err != nil {
		uiRenderer.RenderError("Uninstall failed: " + err.Error())
		os.Exit(1)
	}

	uiRenderer.RenderSuccess("openspdd has been uninstalled.")
}
