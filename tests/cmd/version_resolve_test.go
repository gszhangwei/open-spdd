package cmd_test

import (
	"runtime/debug"
	"testing"

	"github.com/gszhangwei/open-spdd/cmd"
)

// stubBuildInfo returns a readBuildInfo seam implementation that always reports
// the supplied module version, so tests can drive the build-info fallback path
// without invoking the real runtime/debug.ReadBuildInfo.
func stubBuildInfo(version string) func() (*debug.BuildInfo, bool) {
	return func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{Main: debug.Module{Version: version}}, true
	}
}

func TestResolveVersion_LdflagSetWins(t *testing.T) {
	restore := cmd.SetReadBuildInfoForTest(stubBuildInfo("v9.9.9"))
	t.Cleanup(restore)

	if got, want := cmd.ResolveVersion("v1.2.3"), "v1.2.3"; got != want {
		t.Fatalf("ResolveVersion(%q) = %q, want %q (ldflag must win over build info)", "v1.2.3", got, want)
	}
}

func TestResolveVersion_FallsBackToBuildInfo_WhenInjectedIsDev(t *testing.T) {
	restore := cmd.SetReadBuildInfoForTest(stubBuildInfo("v0.4.12"))
	t.Cleanup(restore)

	if got, want := cmd.ResolveVersion("dev"), "v0.4.12"; got != want {
		t.Fatalf("ResolveVersion(\"dev\") = %q, want %q (build info must take over from \"dev\")", got, want)
	}
}

func TestResolveVersion_FallsBackToBuildInfo_WhenInjectedIsEmpty(t *testing.T) {
	restore := cmd.SetReadBuildInfoForTest(stubBuildInfo("v0.4.12"))
	t.Cleanup(restore)

	if got, want := cmd.ResolveVersion(""), "v0.4.12"; got != want {
		t.Fatalf("ResolveVersion(\"\") = %q, want %q (empty injected must behave like \"dev\")", got, want)
	}
}

func TestResolveVersion_BuildInfoDevelTreatedAsMissing(t *testing.T) {
	restore := cmd.SetReadBuildInfoForTest(stubBuildInfo("(devel)"))
	t.Cleanup(restore)

	if got, want := cmd.ResolveVersion("dev"), "dev"; got != want {
		t.Fatalf("ResolveVersion(\"dev\") = %q, want %q (\"(devel)\" must be treated as no-version)", got, want)
	}
}

func TestResolveVersion_BuildInfoEmptyVersionTreatedAsMissing(t *testing.T) {
	restore := cmd.SetReadBuildInfoForTest(stubBuildInfo(""))
	t.Cleanup(restore)

	if got, want := cmd.ResolveVersion("dev"), "dev"; got != want {
		t.Fatalf("ResolveVersion(\"dev\") = %q, want %q (empty bi.Main.Version must be treated as no-version)", got, want)
	}
}

func TestResolveVersion_BuildInfoUnavailable(t *testing.T) {
	restore := cmd.SetReadBuildInfoForTest(func() (*debug.BuildInfo, bool) {
		return nil, false
	})
	t.Cleanup(restore)

	if got, want := cmd.ResolveVersion("dev"), "dev"; got != want {
		t.Fatalf("ResolveVersion(\"dev\") = %q, want %q ((nil,false) must fall through to injected)", got, want)
	}
}

// TestResolveVersion_PseudoVersionFlowsThrough verifies that go install's
// pseudo-version format (e.g., generated for `@main`-style installs) is
// surfaced unchanged, since the resolver does no normalization beyond the
// emptiness/"(devel)" guards.
func TestResolveVersion_PseudoVersionFlowsThrough(t *testing.T) {
	const pseudo = "v0.4.13-0.20260429154200-abcdef123456"
	restore := cmd.SetReadBuildInfoForTest(stubBuildInfo(pseudo))
	t.Cleanup(restore)

	if got := cmd.ResolveVersion("dev"); got != pseudo {
		t.Fatalf("ResolveVersion(\"dev\") = %q, want %q (pseudo-version must flow through unchanged)", got, pseudo)
	}
}
