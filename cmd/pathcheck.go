package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// programName is the bare command name we expect to be reachable via $PATH.
const programName = "openspdd"

// pathHintMarkerName is the filename written under the user's config directory
// once the PATH hint has been displayed, so we never nag the user twice.
const pathHintMarkerName = ".path-hint-shown"

// maybePrintPathHint prints a one-time message telling the user how to add the
// `openspdd` binary's directory to $PATH, but ONLY when running `openspdd` by
// bare name would not actually find this (or any) openspdd binary on $PATH.
//
// Why exec.LookPath and not "is binDir in PATH":
//
//   - For Homebrew installs the real binary lives in
//     /opt/homebrew/Cellar/openspdd/<ver>/bin (NOT on PATH), but a symlink
//     under /opt/homebrew/bin (which IS on PATH) makes `openspdd` invokable.
//     Comparing the resolved binary directory to $PATH would falsely flag
//     these users.
//   - LookPath answers the question the user actually cares about:
//     "Will typing `openspdd` work in my next shell?"
//
// Output goes to stderr so it never pollutes pipelines that read stdout.
func maybePrintPathHint() {
	if _, err := exec.LookPath(programName); err == nil {
		return
	}

	binDir := guessInstallDir()

	markerPath := pathHintMarkerPath()
	if markerPath != "" {
		if _, err := os.Stat(markerPath); err == nil {
			return
		}
	}

	printPathInstructions(binDir)

	if markerPath != "" {
		_ = os.MkdirAll(filepath.Dir(markerPath), 0o755)
		_ = os.WriteFile(markerPath, []byte("shown"), 0o644)
	}
}

// guessInstallDir returns the directory we believe the running binary lives in.
// It is best-effort and used only for display purposes in the PATH hint.
func guessInstallDir() string {
	exePath, err := os.Executable()
	if err != nil {
		return "<your Go bin directory>"
	}
	if resolved, err := filepath.EvalSymlinks(exePath); err == nil {
		exePath = resolved
	}
	return filepath.Dir(exePath)
}

// pathHintMarkerPath returns the file path used to remember that we already
// showed the PATH hint. Returns an empty string if the user config dir is
// unavailable, in which case the hint will simply repeat on the next run.
func pathHintMarkerPath() string {
	cfgDir, err := os.UserConfigDir()
	if err != nil || cfgDir == "" {
		return ""
	}
	return filepath.Join(cfgDir, "openspdd", pathHintMarkerName)
}

// printPathInstructions prints shell-specific instructions for adding binDir
// to the user's PATH. Uses stderr so stdout consumers (scripts, pipes) are
// unaffected.
func printPathInstructions(binDir string) {
	w := os.Stderr

	fmt.Fprintln(w)
	fmt.Fprintln(w, "──────────────────────────────────────────────────────────────")
	fmt.Fprintln(w, " openspdd is installed but its directory is not on your PATH.")
	fmt.Fprintln(w, "──────────────────────────────────────────────────────────────")
	fmt.Fprintf(w, " Binary location: %s\n", binDir)
	fmt.Fprintln(w)

	switch runtime.GOOS {
	case "windows":
		fmt.Fprintln(w, " To make `openspdd` available in any terminal, add it to PATH:")
		fmt.Fprintln(w)
		fmt.Fprintf(w, "   PowerShell (current user, persistent):\n"+
			"     [Environment]::SetEnvironmentVariable(\"Path\",\n"+
			"       [Environment]::GetEnvironmentVariable(\"Path\",\"User\") + \";%s\",\n"+
			"       \"User\")\n", binDir)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "   Then open a new terminal window.")
	default:
		shell := filepath.Base(os.Getenv("SHELL"))
		rcFile := "~/.profile"
		switch shell {
		case "zsh":
			rcFile = "~/.zshrc"
		case "bash":
			rcFile = "~/.bashrc"
		case "fish":
			rcFile = "~/.config/fish/config.fish"
		}

		if shell == "fish" {
			fmt.Fprintf(w, " Append this line to %s:\n\n", rcFile)
			fmt.Fprintf(w, "   set -gx PATH %s $PATH\n", binDir)
		} else {
			fmt.Fprintf(w, " Append this line to %s:\n\n", rcFile)
			fmt.Fprintf(w, "   export PATH=\"%s:$PATH\"\n", binDir)
		}
		fmt.Fprintln(w)
		fmt.Fprintf(w, " Then reload your shell:\n\n   source %s\n", rcFile)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, " (This message will not be shown again.)")
	fmt.Fprintln(w, "──────────────────────────────────────────────────────────────")
	fmt.Fprintln(w)
}
