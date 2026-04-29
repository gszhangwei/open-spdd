package cmd_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/gszhangwei/open-spdd/cmd"
)

func TestClassifyByPath_Homebrew(t *testing.T) {
	cases := []string{
		"/opt/homebrew/Cellar/openspdd/0.1.0/bin/openspdd",
		"/usr/local/Cellar/openspdd/0.2.3/bin/openspdd",
		"/home/linuxbrew/.linuxbrew/Cellar/openspdd/1.0.0/bin/openspdd",
	}
	for _, p := range cases {
		t.Run(p, func(t *testing.T) {
			method, reason := cmd.ClassifyByPathForTest(p)
			if method != cmd.MethodHomebrew {
				t.Fatalf("expected MethodHomebrew for %q, got %q (%s)", p, method, reason)
			}
			if reason == "" {
				t.Fatal("expected non-empty reason")
			}
		})
	}
}

func TestClassifyByPath_Unknown(t *testing.T) {
	cases := []string{
		"/usr/local/bin/openspdd",
		"/tmp/openspdd",
		"C:/Program Files/openspdd.exe",
		"/some/random/place/openspdd",
		"",
	}
	for _, p := range cases {
		t.Run(p, func(t *testing.T) {
			method, _ := cmd.ClassifyByPathForTest(p)
			if method != cmd.MethodUnknown {
				t.Fatalf("expected MethodUnknown for %q, got %q", p, method)
			}
		})
	}
}

func TestBuildPlan_Homebrew_HasBrewStepThenMarker(t *testing.T) {
	ctx := cmd.InstallContext{
		Method:       cmd.MethodHomebrew,
		BinaryPath:   "/opt/homebrew/bin/openspdd",
		ResolvedPath: "/opt/homebrew/Cellar/openspdd/1.2.3/bin/openspdd",
		FormulaName:  "gszhangwei/tools/openspdd",
		MarkerPath:   "/Users/test/Library/Application Support/openspdd/.path-hint-shown",
	}
	plan := cmd.BuildPlanForOSForTest(ctx, "darwin")

	if plan.Method != cmd.MethodHomebrew {
		t.Fatalf("expected plan.Method=Homebrew, got %q", plan.Method)
	}
	if got := len(plan.Steps); got != 2 {
		t.Fatalf("expected 2 steps (brew + marker), got %d", got)
	}

	step0 := plan.Steps[0]
	if step0.Kind != cmd.StepRunCommand {
		t.Fatalf("step[0].Kind = %q, want StepRunCommand", step0.Kind)
	}
	if step0.Program != "brew" {
		t.Fatalf("step[0].Program = %q, want \"brew\"", step0.Program)
	}
	if got, want := step0.Args, []string{"uninstall", "gszhangwei/tools/openspdd"}; !equalStringSlice(got, want) {
		t.Fatalf("step[0].Args = %v, want %v", got, want)
	}

	step1 := plan.Steps[1]
	if step1.Kind != cmd.StepRemoveFile {
		t.Fatalf("step[1].Kind = %q, want StepRemoveFile", step1.Kind)
	}
	if !step1.Optional {
		t.Fatal("marker cleanup step must be Optional=true")
	}
}

func TestBuildPlan_Homebrew_OmitsMarkerWhenEmpty(t *testing.T) {
	ctx := cmd.InstallContext{
		Method:       cmd.MethodHomebrew,
		ResolvedPath: "/opt/homebrew/Cellar/openspdd/1.0.0/bin/openspdd",
		FormulaName:  "gszhangwei/tools/openspdd",
	}
	plan := cmd.BuildPlanForOSForTest(ctx, "linux")

	if got := len(plan.Steps); got != 1 {
		t.Fatalf("expected 1 step (no marker step), got %d", got)
	}
	if plan.Steps[0].Kind != cmd.StepRunCommand {
		t.Fatalf("only step should be the brew run-command, got %q", plan.Steps[0].Kind)
	}
}

func TestBuildPlan_Homebrew_FillsMissingFormulaName(t *testing.T) {
	ctx := cmd.InstallContext{
		Method:       cmd.MethodHomebrew,
		ResolvedPath: "/opt/homebrew/Cellar/openspdd/1.0.0/bin/openspdd",
	}
	plan := cmd.BuildPlanForOSForTest(ctx, "darwin")

	if got, want := plan.Steps[0].Args, []string{"uninstall", "gszhangwei/tools/openspdd"}; !equalStringSlice(got, want) {
		t.Fatalf("step[0].Args = %v, want %v (default formula)", got, want)
	}
}

func TestBuildPlan_GoInstall_Unix_RemovesBinary(t *testing.T) {
	ctx := cmd.InstallContext{
		Method:       cmd.MethodGoInstall,
		BinaryPath:   "/Users/test/go/bin/openspdd",
		ResolvedPath: "/Users/test/go/bin/openspdd",
		MarkerPath:   "/Users/test/Library/Application Support/openspdd/.path-hint-shown",
	}
	plan := cmd.BuildPlanForOSForTest(ctx, "linux")

	if got := len(plan.Steps); got != 2 {
		t.Fatalf("expected 2 steps (binary + marker), got %d", got)
	}
	step0 := plan.Steps[0]
	if step0.Kind != cmd.StepRemoveFile {
		t.Fatalf("step[0].Kind = %q, want StepRemoveFile", step0.Kind)
	}
	if step0.Optional {
		t.Fatal("binary removal must NOT be Optional")
	}
	if got, want := step0.Paths, []string{ctx.ResolvedPath}; !equalStringSlice(got, want) {
		t.Fatalf("step[0].Paths = %v, want %v", got, want)
	}
}

func TestBuildPlan_GoInstall_Windows_PrintsAdvisory(t *testing.T) {
	ctx := cmd.InstallContext{
		Method:       cmd.MethodGoInstall,
		BinaryPath:   `C:\Users\test\go\bin\openspdd.exe`,
		ResolvedPath: `C:\Users\test\go\bin\openspdd.exe`,
	}
	plan := cmd.BuildPlanForOSForTest(ctx, "windows")

	if got := len(plan.Steps); got != 1 {
		t.Fatalf("expected 1 advisory step on Windows, got %d", got)
	}
	step0 := plan.Steps[0]
	if step0.Kind != cmd.StepAdvisory {
		t.Fatalf("Windows go-install step must be Advisory, got %q", step0.Kind)
	}
	if !strings.Contains(step0.Description, "del") {
		t.Fatalf("advisory should mention `del` command, got %q", step0.Description)
	}
	if !strings.Contains(step0.Description, ctx.ResolvedPath) {
		t.Fatalf("advisory should mention the resolved path %q, got %q", ctx.ResolvedPath, step0.Description)
	}
}

func TestBuildPlan_Unknown_OnlyAdvisoryNoDestructive(t *testing.T) {
	ctx := cmd.InstallContext{
		Method:       cmd.MethodUnknown,
		ResolvedPath: "/opt/custom/openspdd",
		MarkerPath:   "/Users/test/Library/Application Support/openspdd/.path-hint-shown",
	}
	plan := cmd.BuildPlanForOSForTest(ctx, runtime.GOOS)

	if got := len(plan.Steps); got != 1 {
		t.Fatalf("Unknown method must produce exactly 1 step, got %d", got)
	}
	if plan.Steps[0].Kind != cmd.StepAdvisory {
		t.Fatalf("Unknown plan step must be Advisory, got %q", plan.Steps[0].Kind)
	}
	for _, s := range plan.Steps {
		if s.Kind == cmd.StepRemoveFile || s.Kind == cmd.StepRunCommand {
			t.Fatalf("Unknown plan must NOT contain destructive steps, found %q", s.Kind)
		}
	}
}

func TestUninstallPlan_Describe_ContainsKeyDetails(t *testing.T) {
	ctx := cmd.InstallContext{
		Method:       cmd.MethodHomebrew,
		ResolvedPath: "/opt/homebrew/Cellar/openspdd/1.2.3/bin/openspdd",
		FormulaName:  "gszhangwei/tools/openspdd",
	}
	plan := cmd.BuildPlanForOSForTest(ctx, "darwin")

	out := plan.Describe()

	if !strings.Contains(out, "homebrew") {
		t.Fatalf("Describe() should mention method 'homebrew', got:\n%s", out)
	}
	if !strings.Contains(out, "brew uninstall gszhangwei/tools/openspdd") {
		t.Fatalf("Describe() should include the brew command, got:\n%s", out)
	}
}

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
