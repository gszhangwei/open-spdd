package cmd_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gszhangwei/open-spdd/cmd"
	"github.com/gszhangwei/open-spdd/internal/templates"
)

// stubRenderer captures every Render*/Confirm call so tests can assert on
// the user-facing output without coupling to terminal control codes.
type stubRenderer struct {
	successes []string
	warnings  []string
	errors    []string
	confirms  []string
	confirm   bool
}

func (s *stubRenderer) RenderSuccess(msg string)                      { s.successes = append(s.successes, msg) }
func (s *stubRenderer) RenderError(msg string)                        { s.errors = append(s.errors, msg) }
func (s *stubRenderer) RenderWarning(msg string)                      { s.warnings = append(s.warnings, msg) }
func (s *stubRenderer) RenderTable(headers []string, rows [][]string) {}
func (s *stubRenderer) SelectTemplate(_ []templates.TemplateMeta) (templates.TemplateMeta, error) {
	return templates.TemplateMeta{}, nil
}
func (s *stubRenderer) Confirm(prompt string) bool {
	s.confirms = append(s.confirms, prompt)
	return s.confirm
}

// noopProgram returns a path to a system binary that exits 0 with no output,
// so we can exercise the run-command path without running brew. Skip the
// test on Windows where /bin/true does not exist.
func noopProgram(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("skipping run-command success test on Windows")
	}
	return "/usr/bin/true"
}

func TestExecutor_RunCommand_Success(t *testing.T) {
	prog := noopProgram(t)

	restoreLook := cmd.SwapLookPathForTest(func(name string) (string, error) {
		if name == "brew" {
			return prog, nil
		}
		return "", errors.New("not used")
	})
	defer restoreLook()

	restoreExec := cmd.SwapExecCommandForTest(func(name string, args ...string) *exec.Cmd {
		return exec.Command(prog)
	})
	defer restoreExec()

	plan := cmd.UninstallPlan{
		Method: cmd.MethodHomebrew,
		Steps: []cmd.PlanStep{{
			Kind:        cmd.StepRunCommand,
			Description: "Run Homebrew",
			Program:     "brew",
			Args:        []string{"uninstall", "gszhangwei/tools/openspdd"},
			Command:     "brew uninstall gszhangwei/tools/openspdd",
		}},
	}

	stub := &stubRenderer{}
	executor := cmd.NewUninstallExecutorForTest(stub)
	if err := executor.Execute(plan); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(stub.successes) == 0 {
		t.Fatal("expected at least one RenderSuccess call")
	}
	if len(stub.errors) != 0 {
		t.Fatalf("expected no errors, got %v", stub.errors)
	}
}

func TestExecutor_RunCommand_LookPathMissing(t *testing.T) {
	restoreLook := cmd.SwapLookPathForTest(func(name string) (string, error) {
		if name == "brew" {
			return "", errors.New("executable file not found in $PATH")
		}
		// post-execute leftover check probes "openspdd"
		return "", errors.New("not found")
	})
	defer restoreLook()

	plan := cmd.UninstallPlan{
		Method: cmd.MethodHomebrew,
		Steps: []cmd.PlanStep{{
			Kind:        cmd.StepRunCommand,
			Description: "Run Homebrew",
			Program:     "brew",
			Args:        []string{"uninstall", "gszhangwei/tools/openspdd"},
			Command:     "brew uninstall gszhangwei/tools/openspdd",
		}},
	}

	stub := &stubRenderer{}
	executor := cmd.NewUninstallExecutorForTest(stub)
	err := executor.Execute(plan)
	if err == nil {
		t.Fatal("expected error when brew is not on PATH")
	}
	if len(stub.errors) == 0 {
		t.Fatal("expected RenderError to be called when brew is missing")
	}
	if !strings.Contains(stub.errors[0], "brew") {
		t.Fatalf("error message should mention 'brew', got %q", stub.errors[0])
	}
}

func TestExecutor_RemoveFile_Success(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "openspdd-fake-binary")
	if err := os.WriteFile(target, []byte("fake"), 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// LookPath stub for the post-execute leftover check (must always succeed
	// to error and skip the advisory).
	restoreLook := cmd.SwapLookPathForTest(func(name string) (string, error) {
		return "", errors.New("not found")
	})
	defer restoreLook()

	plan := cmd.UninstallPlan{
		Method: cmd.MethodGoInstall,
		Steps: []cmd.PlanStep{{
			Kind:        cmd.StepRemoveFile,
			Description: "Remove the openspdd binary",
			Paths:       []string{target},
			Optional:    false,
		}},
	}

	stub := &stubRenderer{}
	executor := cmd.NewUninstallExecutorForTest(stub)
	if err := executor.Execute(plan); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if _, err := os.Stat(target); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected target to be removed; stat err=%v", err)
	}
	if len(stub.successes) == 0 || !strings.Contains(stub.successes[0], "Removed") {
		t.Fatalf("expected RenderSuccess for removal, got %v", stub.successes)
	}
}

func TestExecutor_RemoveFile_OptionalMissingIsSilent(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "does-not-exist")

	restoreLook := cmd.SwapLookPathForTest(func(name string) (string, error) {
		return "", errors.New("not found")
	})
	defer restoreLook()

	plan := cmd.UninstallPlan{
		Method: cmd.MethodHomebrew,
		Steps: []cmd.PlanStep{{
			Kind:        cmd.StepRemoveFile,
			Description: "Remove openspdd's first-run marker file",
			Paths:       []string{missing},
			Optional:    true,
		}},
	}

	stub := &stubRenderer{}
	executor := cmd.NewUninstallExecutorForTest(stub)
	if err := executor.Execute(plan); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(stub.warnings) != 0 {
		t.Fatalf("optional missing path should not emit warnings, got %v", stub.warnings)
	}
	if len(stub.errors) != 0 {
		t.Fatalf("optional missing path should not emit errors, got %v", stub.errors)
	}
}

func TestExecutor_PostExecute_LeftoverWarning(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "openspdd-fake")
	if err := os.WriteFile(target, []byte("fake"), 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	restoreLook := cmd.SwapLookPathForTest(func(name string) (string, error) {
		if name == "openspdd" {
			return "/some/other/openspdd", nil
		}
		return "", errors.New("not found")
	})
	defer restoreLook()

	plan := cmd.UninstallPlan{
		Method: cmd.MethodGoInstall,
		Steps: []cmd.PlanStep{{
			Kind:        cmd.StepRemoveFile,
			Description: "Remove the openspdd binary",
			Paths:       []string{target},
		}},
	}

	stub := &stubRenderer{}
	executor := cmd.NewUninstallExecutorForTest(stub)
	if err := executor.Execute(plan); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	foundLeftover := false
	for _, w := range stub.warnings {
		if strings.Contains(w, "still on PATH") {
			foundLeftover = true
			break
		}
	}
	if !foundLeftover {
		t.Fatalf("expected leftover-on-PATH warning, got warnings=%v", stub.warnings)
	}
}
