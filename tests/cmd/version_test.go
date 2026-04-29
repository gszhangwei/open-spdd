package cmd_test

import (
	"testing"

	"github.com/gszhangwei/open-spdd/cmd"
)

func TestSetVersion_AssignsRootVersion(t *testing.T) {
	cmd.SetVersion("v1.2.3")

	root := cmd.RootCommand()
	if got, want := root.Version, "v1.2.3"; got != want {
		t.Fatalf("rootCmd.Version = %q, want %q", got, want)
	}

	flag := root.Flags().Lookup("version")
	if flag == nil {
		t.Fatal("expected a 'version' flag to be registered after SetVersion")
	}
	if got, want := flag.Shorthand, "v"; got != want {
		t.Fatalf("version flag Shorthand = %q, want %q", got, want)
	}
}

func TestSetVersion_EmptyDefaultsToDev(t *testing.T) {
	cmd.SetVersion("")

	root := cmd.RootCommand()
	if got, want := root.Version, "dev"; got != want {
		t.Fatalf("empty SetVersion should default Version to %q, got %q", want, got)
	}
}

func TestSetVersion_TrimsWhitespace(t *testing.T) {
	cmd.SetVersion("  v9.9.9  ")

	root := cmd.RootCommand()
	if got, want := root.Version, "v9.9.9"; got != want {
		t.Fatalf("rootCmd.Version = %q, want %q (whitespace must be trimmed)", got, want)
	}
}

func TestSetVersion_Idempotent(t *testing.T) {
	cmd.SetVersion("v0.0.1")
	cmd.SetVersion("v0.0.1")

	root := cmd.RootCommand()
	if got, want := root.Version, "v0.0.1"; got != want {
		t.Fatalf("rootCmd.Version = %q, want %q", got, want)
	}
	flag := root.Flags().Lookup("version")
	if flag == nil {
		t.Fatal("expected version flag to remain registered")
	}
	if flag.Shorthand != "v" {
		t.Fatalf("expected shorthand 'v' to remain after second SetVersion, got %q", flag.Shorthand)
	}
}
