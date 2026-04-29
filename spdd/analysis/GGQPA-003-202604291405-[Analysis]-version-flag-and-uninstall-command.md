# SPDD Analysis: Version Flag (`-v`) and `uninstall` Command for openspdd CLI

## Original Business Requirement

Add two features to the `openspdd` CLI:

1. A `-v` flag that prints the current version of `openspdd`.
2. An `uninstall` subcommand that lets users easily uninstall `openspdd`. Support both Homebrew installs and `go install` installs.

## Domain Concept Identification

### Existing Concepts (from codebase)

- **CLI Root Command (`rootCmd`)** — the Cobra command tree rooted at `openspdd` (defined in `cmd/root.go`). Owns persistent flags (currently only `--tool` / `-t`). All subcommands attach to it via `rootCmd.AddCommand(...)` in their respective `init()` functions. New version flag and new `uninstall` subcommand attach here.
- **CLI Subcommand** — the convention used by `init`, `list`, `generate` for each top-level verb. Each lives in its own file under `cmd/`, declares its `cobra.Command`, registers flags, and attaches in `init()`. The new `uninstall` verb will follow the same convention.
- **Binary Path Resolution** — already implemented in `cmd/pathcheck.go` via `os.Executable()` + `filepath.EvalSymlinks(...)` (see `guessInstallDir`). This same primitive is the foundation for _Install Method Detection_ (a new concept).
- **PATH Reachability Check** — `exec.LookPath("openspdd")` in `cmd/pathcheck.go`. Useful to confirm whether the user has a _second_ binary still resolvable after uninstall (e.g., user uninstalls the brew copy but a `go install` copy still shadows on PATH).
- **User Config Directory Marker** — `pathHintMarkerPath()` in `cmd/pathcheck.go` writes `<UserConfigDir>/openspdd/.path-hint-shown`. This is openspdd-owned residual state that an uninstall should clean up.
- **UI Renderer (`uiRenderer`)** — `internal/ui` exposes `RenderSuccess/Warning/Error` and a `Confirm(prompt)` interactive prompt powered by `charmbracelet/huh`. The uninstall flow will reuse `Confirm` for destructive-action confirmation and the renderer methods for output.
- **GoReleaser Build Pipeline (`.goreleaser.yaml`)** — already injects `-s -w -X main.version={{.Version}}` at build time, **but the `main.version` symbol does not currently exist** in `cmd/openspdd/main.go`. The ldflag is silently dropped today. The Homebrew formula's smoke test (`system "#{bin}/openspdd", "--version"`) would also fail today — this requirement implicitly fixes a latent bug.
- **Homebrew Tap (`gszhangwei/homebrew-tools`, formula name `openspdd`)** — established distribution channel. Documented user-facing install command is `brew install gszhangwei/tools/openspdd`. The matching uninstall command must use the same fully-qualified formula name.
- **Generated Template Files** — `.cursor/commands/`, `.claude/commands/`, etc., produced inside _user projects_ by `openspdd generate`. **These are user-owned project artifacts** and explicitly out of scope for uninstall (see Risk & Gap Analysis).

### New Concepts Required

- **Version Metadata** — a build-time-injectable string identifying the running binary (e.g., `v1.2.3`, `dev`). Surfaced via the `-v` / `--version` flag and consumable by the `uninstall` flow if it wants to display "you are about to uninstall openspdd vX.Y.Z". Wired through the existing goreleaser ldflag to `main.version`.
- **Install Method** — a classification of how the running binary got onto the user's machine. Required values cover: Homebrew, `go install`, manual/unknown. Derived from the resolved binary path (Homebrew installs live under `Cellar/openspdd/.../bin`; `go install` installs live under `$GOBIN` or `$GOPATH/bin`). Drives which uninstall _Plan_ applies.
- **Uninstall Plan** — an install-method-specific recipe describing what will be removed (binary, symlinks, residual marker file) and which external command (if any) must be invoked (`brew uninstall ...` for Homebrew; direct file removal for `go install`). Each plan is presented to the user before execution so the action is transparent.
- **Uninstall Executor** — the actor that carries out a plan: shells out to `brew` for Homebrew, removes files for `go install`, and (on best-effort basis) cleans up the openspdd-owned config marker. Owns the safety boundary: it never touches user-project files and never deletes paths it cannot confidently classify.

### Conceptual Relationships

- _CLI Root Command_ → owns → _Version Metadata_ (surfaced via `-v` flag) and _Uninstall_ subcommand.
- _Uninstall_ subcommand → uses → _Install Method_ (detection) → selects → _Uninstall Plan_ → executed by → _Uninstall Executor_.
- _Install Method_ detection → reuses → _Binary Path Resolution_ primitive (already in `pathcheck.go`).
- _Uninstall Executor_ → coordinates with → _UI Renderer_ for confirmation and output.
- _Uninstall Plan_ → may include → cleanup of _User Config Directory Marker_.
- _Version Metadata_ → injected by → _GoReleaser Build Pipeline_ (existing wiring, not currently terminating in a real symbol).

### Key Business Rules

- **Version source-of-truth is the build** — version is set at link time via ldflags. A development build (plain `go build` / `go run` without ldflags) MUST still produce a sensible value (e.g., `dev` or `(devel)`) — never a crash, never an empty string.
- **Uninstall is destructive and MUST be confirmed** — before any file removal or `brew uninstall` invocation, the user must explicitly confirm. A non-interactive override is needed for scripted use.
- **Tool MUST detect, not assume, the install method** — auto-detection from the running binary's resolved path. If detection is ambiguous, default to _show plan, ask user to confirm method or abort_; never guess destructively.
- **Tool MUST NOT delete files outside of openspdd's own footprint** — the binary, its symlink (via `brew uninstall`), and openspdd's marker file. Generated templates inside user projects (`.cursor/commands/spdd-*.md`, etc.) are user property and stay untouched.
- **Tool MUST surface what it is about to do before doing it** — the uninstall plan (paths, commands) is printed before confirmation so the user sees the exact effects.
- **Cross-platform behavior MUST be consistent in intent, but adapted in mechanics** — on platforms where self-deletion is unreliable (Windows), the tool falls back to printing the exact manual command for the user to run.

## Strategic Approach

### Solution Direction

- **Feature 1 (`-v` version flag)**: introduce a single `version` package-level variable in `cmd/openspdd/main.go` (matches the existing `-X main.version=...` ldflag wiring) with a sensible default of `dev`. Pass the value into the `cmd` package and configure Cobra's built-in version support on `rootCmd`. Reuse Cobra's auto-generated `--version` flag and override its definition to add the `-v` short alias. Customize the version output template so that `openspdd -v` prints a clean line such as `openspdd v1.2.3`.
- **Feature 2 (`uninstall` subcommand)**: add `cmd/uninstall.go` following the same single-file-per-subcommand convention as `init.go`/`list.go`/`generate.go`. The command's runtime flow is: (a) detect _Install Method_ from the resolved running binary path, (b) build an _Uninstall Plan_ describing exact actions, (c) display the plan and ask `Confirm(...)` via the existing `uiRenderer`, (d) hand off to the _Uninstall Executor_ which either shells out to `brew uninstall <tap-qualified-name>` or removes the binary file directly, then cleans the openspdd config marker. Provide non-interactive flags (`--yes`, `--dry-run`) for scriptability.
- **Reuse over reinvention**: the binary-path/symlink resolution logic in `cmd/pathcheck.go` is the basis of install-method detection — extract a shared helper rather than duplicating `os.Executable() + EvalSymlinks` logic.
- **Channel-honest messaging**: the README documents `brew install gszhangwei/tools/openspdd`; the uninstall plan must mirror that exact form (`brew uninstall gszhangwei/tools/openspdd`) so users see commands they recognize.
- **Latent-bug resolution as a side effect**: the Homebrew formula's smoke test `system "#{bin}/openspdd", "--version"` is currently broken because no version flag exists; landing Feature 1 makes the formula test correct without a separate change.

### Key Design Decisions

- **Where the `version` variable lives** — in `cmd/openspdd/main.go` (matching existing goreleaser ldflag `-X main.version`) versus in the `cmd` package (would require updating `.goreleaser.yaml` to `-X github.com/.../cmd.version`).
  - _Trade-offs_: Keeping it in `main` honors the existing pipeline contract and avoids touching CI; it requires `main` to push the value into the `cmd` package (small surface, e.g., `cmd.SetVersion(version)` or pass into `cmd.Execute(version)`). Moving it into `cmd` would be more idiomatic for a "fat cmd, thin main" layout but requires a coordinated goreleaser change and a release-time risk window.
  - _Recommendation_: keep the variable in `main`, expose a one-line setter or accept it as an `Execute` argument. Lower-risk, no CI change, fixes the latent ldflag-no-target issue immediately.

- **`-v` short flag vs. `verbose` convention** — the natural shortcut for `--version` is `-v`, but `-v` is also the customary shortcut for `--verbose` in many CLIs.
  - _Trade-offs_: Claiming `-v` for version now blocks `--verbose -v` later; reserving `-v` for a future verbose flag means version users must type `--version`. The codebase has no verbose flag today and no plan documented in the requirement.
  - _Recommendation_: claim `-v` for `--version` per the explicit user request. If a verbose mode is ever added, use `-V` (capital) or a different shortcut. Document the rationale in code comment so it's not silently re-shadowed later.

- **How the version flag is wired in Cobra** — Cobra adds `--version` automatically when `rootCmd.Version` is non-empty, but the auto-flag is long-only.
  - _Trade-offs_: (A) Set `rootCmd.Version` and _additionally_ override the version flag definition to also bind `-v`; minimal code, uses cobra's built-in version handling. (B) Disable Cobra's auto-flag entirely and implement a fully manual `-v`/`--version` BoolVarP — more code, more flexibility, higher chance of behavioral drift.
  - _Recommendation_: Option A. Set `rootCmd.Version` and replace the auto-generated `version` flag with one that has `Shorthand: "v"`. Customize via `SetVersionTemplate` for output formatting.

- **Uninstall — print-only, execute-with-confirmation, or hybrid** — three UX models.
  - _Trade-offs_: Print-only is safest but doesn't satisfy "easy uninstall". Always-execute (with confirmation) is convenient but risks scripted misuse. Hybrid (default = confirm + execute; opt-in `--dry-run` to print only; opt-out `--yes` to skip confirmation) covers all three audiences.
  - _Recommendation_: Hybrid. Mirrors familiar tooling (e.g., `brew`, `apt`). Aligns with the requirement's "simple uninstall" intent without giving up the confirmation safety net.

- **Self-deletion behavior** — the running binary is the one being uninstalled.
  - _Trade-offs_: On macOS/Linux, deleting an executing binary is well-defined (the inode persists for the running process), so `os.Remove(binaryPath)` works for the `go install` case. On Windows, it fails with sharing violations. For Homebrew, `brew uninstall` is invoked as a subprocess; the brew tool handles symlink + Cellar removal regardless of the running openspdd process state.
  - _Recommendation_: Treat self-deletion as a normal expected case on Unix. On Windows, detect the platform, skip the remove, and print the exact manual `del` / `Remove-Item` command. Print a parting message before exit on all platforms so the user sees confirmation even if the binary vanishes.

- **What constitutes "uninstall scope"** — binary only? plus marker? plus user templates? plus Homebrew tap?
  - _Trade-offs_: Removing more is "cleaner" but invades user property. Removing less risks leaving "something behind".
  - _Recommendation_: in scope = (a) the binary (via `brew uninstall` or file removal), (b) the openspdd config marker file under `<UserConfigDir>/openspdd/`. Out of scope = (c) generated templates in user projects (user content), (d) the Homebrew tap (other tools may use it). Document this scope in the printed plan.

- **What if the `uninstall` command can't classify the install method** — manual binary move, system package, container, etc.
  - _Trade-offs_: Refusing to act protects the user but feels broken. Guessing risks data loss.
  - _Recommendation_: when the install method is _Unknown_, refuse to execute, print the resolved binary path, and tell the user to remove it manually. Exit non-zero so scripts catch the unhandled case.

### Alternatives Considered

- **Auto-update / re-install command instead of uninstall** — out of scope; the requirement asks specifically for uninstall.
- **Embed full Homebrew formula management (tap removal, dependency cleanup) in the CLI** — rejected because Homebrew provides idiomatic commands the user already knows, and the tap may host other tools.
- **Print-only uninstall (no execution)** — rejected because the requirement explicitly wants the uninstall to be "simple"; making the user copy-paste contradicts that goal. Retained as a `--dry-run` mode.
- **Use a third-party self-update / self-uninstall library** — rejected. Adds dependency surface; requirements only need two narrow flows that fit cleanly in ~150 LOC.
- **Reading the installed Homebrew formula's metadata to determine the formula name** — rejected as over-engineering. The tap and formula name are stable, owned by us, and already hard-coded in `.goreleaser.yaml` and the README.

## Risk & Gap Analysis

### Requirement Ambiguities

- **Exact format of `-v` output** — should it print just the version string (`v1.2.3`), `openspdd v1.2.3`, or include commit/date/build info? _Needs clarification before REASONS Canvas finalizes the output template._ Recommendation: `openspdd <version>` on a single line, matching common CLI conventions and the existing Homebrew formula smoke test expectation.
- **Scope of "easy uninstall"** — does this include cleaning up generated templates inside user projects, or only the openspdd binary and its own config? Recommendation in this analysis: binary + openspdd's own marker file only. _Confirm with user before implementation._
- **Should the long flag remain `--version` or become something else** — the requirement only mentions `-v`. Convention strongly favors `--version` as the long form; this analysis assumes both.
- **Should `uninstall` accept any flags by default** — e.g., `--purge` to remove the marker, `--keep-config` to not? Analysis recommends `--dry-run` and `--yes` as the minimum useful set; no `--purge` (marker is always removed because it has no user data, only one byte of state).
- **Behavior on multiple installations** — if both a Homebrew copy and a `go install` copy exist on the system, the running binary is one of them; should `uninstall` warn about the other? Recommendation: after uninstall, run `exec.LookPath("openspdd")` again; if another copy is now reachable on PATH, print a hint identifying it.

### Edge Cases

- **Development builds (`go run ./cmd/openspdd`, plain `go build`)** — `version` variable is empty because no ldflag was applied. The `-v` output must still be sensible (`openspdd dev` or `openspdd (devel)`).
- **Homebrew install where `brew` command is not on PATH** — extremely rare but possible (e.g., user removed brew). The uninstall executor must detect missing `brew` and fall back to printing the manual `rm` paths under the resolved Cellar location.
- **`go install` with custom `GOBIN`** — the binary may not be at `$GOPATH/bin/openspdd`. Detection should resolve the _running binary's_ path rather than infer from `go env`.
- **Symlink farms** — Homebrew uses `/opt/homebrew/bin/openspdd` (Apple Silicon) or `/usr/local/bin/openspdd` (Intel) as a symlink into Cellar. `EvalSymlinks` is required to classify correctly. Both prefixes must be recognized.
- **User installed via the project's `scripts/install.sh`** — that script wraps `go install`, so the binary lives at `$GOBIN/openspdd`. Detection should classify this as _go install_.
- **User installed by manually `go build`-ing and `cp`-ing into `/usr/local/bin`** — neither Homebrew nor `go install` pattern matches. Classified as _Unknown_; uninstall refuses and prints the resolved path with manual instructions.
- **Windows `go install` users** — binary at `%USERPROFILE%\go\bin\openspdd.exe`. Self-removal of a running EXE typically fails. Plan must adapt.
- **Permissions** — Homebrew installs in `/usr/local/...` on Intel macOS may need `sudo` in unusual setups. Surface a meaningful error if `brew uninstall` exits non-zero.
- **First-run path-hint marker not yet written** — uninstall's "remove marker" step must tolerate "file does not exist" silently.
- **User runs `openspdd --version` (long form) vs `openspdd -v`** — both must produce identical output.
- **`-v` collides with a future `--verbose` flag** — locked in for version per user request; documented in Strategic Approach.
- **Cobra prints `Error: unknown flag: -v` if the flag isn't registered before parsing** — must be wired in `init()` or before `Execute`.
- **`brew uninstall` interactive prompts** — Homebrew is non-interactive by default for `uninstall`; should still capture stderr to give the user feedback if it fails.
- **Tap-qualified vs. short formula name** — `brew uninstall openspdd` works only if the formula is unambiguous; `brew uninstall gszhangwei/tools/openspdd` is always safe. Use the qualified form in the executed command for robustness.

### Technical Risks

- **Latent broken ldflag wiring (existing bug surfaced by this work)** — `.goreleaser.yaml` injects into `main.version` but no such symbol exists, and the Homebrew formula's smoke test calls a non-existent `--version` flag. _Mitigation_: Feature 1 implements the missing symbol and flag; the very next release will exercise both paths, so the formula test should be re-validated as part of release verification.
- **Self-deletion races on Windows** — `os.Remove` on the running EXE returns `Access is denied`. _Mitigation_: detect `runtime.GOOS == "windows"`, skip removal, print exact manual command, exit 0 with a clear message so the user knows the next step.
- **`brew` not found on PATH** — uninstall path fails. _Mitigation_: `exec.LookPath("brew")` before invocation; if missing, fall back to manual instructions including the resolved Cellar path.
- **User has multiple copies (brew + go install)** — uninstalling one leaves the other; user may be confused. _Mitigation_: post-uninstall `LookPath("openspdd")` check + advisory message identifying the still-resolvable copy.
- **`brew uninstall` prompts or fails non-interactively** — surfaces as a non-zero exit. _Mitigation_: capture and print stdout/stderr from the `brew` invocation; propagate exit code.
- **Cleaning the user-config marker on a system where `os.UserConfigDir()` returns a different path than where the marker was originally written** — possible across user account migrations. _Mitigation_: tolerate "not found", do not error out the whole uninstall.
- **Risk of breaking existing Cobra flag parsing** — adding `-v` short alias could collide with a subcommand-level flag. _Mitigation_: install only on root persistent flags, scan existing subcommand flags (none currently use `-v`).
- **Goreleaser produces `dev` builds with non-semver versions** — version display must not assume semver shape (no parsing).
- **Cross-platform path classification heuristics drift** — Homebrew prefix may change in future macOS/Linux setups. _Mitigation_: classification is best-effort; _Unknown_ fallback always exists.
- **Test coverage gap** — the project has `tests/` directory but the new logic touches process-level concerns (running binary path, subprocess `brew` invocation) that require either fakes or skipping; if integration tests are expected, a `--dry-run`-driven test mode is the path of least resistance.

### Acceptance Criteria Coverage

The original requirement is informally specified — there are no numbered ACs. The two stated capabilities are decomposed below as implicit ACs.

| AC# | Description                                                                                              | Addressable? | Gaps / Notes                                                                                                                                                            |
| --- | -------------------------------------------------------------------------------------------------------- | ------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | Running `openspdd -v` prints the current version of the binary.                                          | Yes          | Output format not specified; analysis recommends `openspdd <version>`. Long alias `--version` assumed.                                                                  |
| 2   | Version value reflects the actual build (release builds show release tag, dev builds show a sane value). | Yes          | Requires fixing the existing latent gap where `.goreleaser.yaml` injects into a non-existent symbol.                                                                    |
| 3   | `openspdd uninstall` works for Homebrew installations.                                                   | Yes          | Will shell out to `brew uninstall gszhangwei/tools/openspdd`. Requires `brew` on PATH; fallback path designed for missing brew.                                         |
| 4   | `openspdd uninstall` works for `go install` installations.                                               | Yes          | Removes the running binary file. Windows self-removal restriction handled by manual-instruction fallback.                                                               |
| 5   | Uninstall is "simple" for the user.                                                                      | Partial      | "Simple" is subjective. Hybrid execute-with-confirmation model (with `--yes` to skip) is the analysis recommendation. Confirm with user that this matches their intent. |
| 6   | Tool does not destroy user data outside its own footprint.                                               | Yes          | Scope explicitly limited to binary + openspdd's own marker. Generated templates in user projects untouched.                                                             |
| 7   | Behavior is predictable when install method cannot be detected.                                          | Yes          | _Unknown_ method ⇒ refuse to execute, print resolved binary path and manual instructions, exit non-zero.                                                                |
