# SPDD Analysis: Add SPDD Explore Mode and Standardize Template Metadata

## Original Business Requirement

New Explore Mode Workflow:

Create spdd-explore.md to establish a safe, read-only thinking phase ("Explore Mode") for AI coding assistants.
Enforce an absolute read-only policy, prohibiting code generation or file mutation in this mode.
Embed the "Distorted Mirror Technique" to surface user assumptions and validate initial scopes.
Design micro action segments (M1-M4) and clear transition criteria before initiating /spdd-analysis.

Configuration & Paths Alignment:

Update Antigravity configuration directory path mapping to .agents/workflows in internal/detector/types.go.

Template Metadata Standardization:

Standardize core Markdown templates (spdd-analysis, spdd-generate, spdd-prompt-update, spdd-reasons-canvas, spdd-sync) by stripping the leading slash (/) from their frontmatter name fields.

Follow-up refinement (same initiative):

refactor: move spdd-explore template to optional and remove '/' prefix from optional template names

- Move `internal/templates/data/core/spdd-explore.md` → `internal/templates/data/optional/spdd-explore.md`
- Strip leading `/` from optional template frontmatter `name` fields (`aupro-context`, `spdd-api-test`, `spdd-code-review`, `spdd-reverse`, `spdd-story`)

Associated design intent already captured in PR prompt draft:

1. **New Explore Mode Workflow**:
   - Create `spdd-explore.md` to establish a safe, read-only thinking phase ("Explore Mode") for AI coding assistants.
   - Enforce an absolute read-only policy, prohibiting code generation or file mutation in this mode.
   - Embed the "Distorted Mirror Technique" to surface user assumptions and validate initial scopes.
   - Design micro action segments (M1-M4) and clear transition criteria before initiating `/spdd-analysis`.

2. **Configuration & Paths Alignment**:
   - Update `Antigravity` configuration directory path mapping to `.agents/workflows` in `internal/detector/types.go`.

3. **Template Metadata Standardization**:
   - Standardize core Markdown templates (`spdd-analysis`, `spdd-generate`, `spdd-prompt-update`, `spdd-reasons-canvas`, `spdd-sync`) by stripping the leading slash (`/`) from their frontmatter `name` fields.

## Domain Concept Identification

### Existing Concepts (from codebase)
- `AIToolType` + `GetConfigDir()`: closed enum of supported assistants; Antigravity currently maps to `.antigravity/commands`, while Codex already owns `.agents/skills` under the shared `.agents/` parent namespace.
- `EmbeddedTemplateManager` + hierarchical template corpus: `data/core/` (default install), `data/optional/` (opt-in via `list --optional` / explicit generate), `data/tools/<tool>/` (tool-specific assets).
- Core SPDD workflow commands: `spdd-analysis`, `spdd-reasons-canvas`, `spdd-generate`, `spdd-prompt-update`, `spdd-sync` — currently shipped with frontmatter `name: /spdd-*`.
- Optional pre/post pipeline commands: `spdd-story`, `spdd-api-test`, `spdd-code-review`, `spdd-reverse`, `aupro-context` — also currently `name: /...`.
- `TemplateMeta.Name` / `GetByName`: lookup accepts either `Name` or `ID` (case-insensitive); many tests and Codex/OpenCode flows still call `GetByName("/spdd-analysis")`.
- `OpenCodeTemplateAdapter`: generation-time strip of the `name` frontmatter key for OpenCode only (implements GGQPA-006 strategy-local adaptation).
- `CodexSkillStrategy`: strips leading `/` from `Name` when emitting skill `name:` for Codex SKILL.md.
- Flat generation strategies for Cursor / Claude Code / Antigravity / OpenCode: write `<configDir>/<id>.md` using filename-as-command identity for most tools.
- Existing SPDD analysis GGQPA-006: previously rejected globally rewriting core templates to remove slash-prefixed `name`, favoring OpenCode-specific output adaptation instead.

### New Concepts Required
- `Explore Mode (spdd-explore)`: optional, read-mostly pre-analysis stance that matures ambiguous requirements before `/spdd-analysis`; includes Distorted Mirror Technique, micro-segments M1–M4, and transition checkpoint gates.
- `Explore Decision Log artifact`: draft markdown records under `spdd/explore/{timestamp}-[Explore]-{description}.md` as the only intended write surface during exploration.
- `Slash-free template name convention`: frontmatter `name` values without leading `/` across core and optional templates, while slash invocation remains a UI/runtime concern of each assistant.
- `Antigravity workflows path`: config target `.agents/workflows` instead of `.antigravity/commands`, aligning Antigravity install layout with Antigravity’s expected workflows directory.

### Key Business Rules
- Explore mode is opt-in (optional template), not part of default `generate --all` core install.
- Explore mode must not implement product code; it may only produce explore decision-log drafts under `spdd/explore/`.
- Transition to `/spdd-analysis` is suggested only when problem scope, chosen direction, mitigated uncertainties, codebase mapping, distorted-mirror validation, and cleared implicit assumptions are all satisfied.
- Template identity for generation remains `id` / filename (`spdd-explore.md` → command `spdd-explore`); frontmatter `name` is metadata for listing/lookup, not the sole command identity.
- Antigravity generation must land under `.agents/workflows` after the path change.
- Changing source frontmatter `name` values must not break `GetByName` for users/tests that still pass slash-prefixed names (ID fallback already supports `spdd-*` without slash; slash-prefixed Name lookups need an explicit compatibility decision).
- OpenCode’s existing strip-`name` adapter and Codex’s slash-trim logic must remain correct after source metadata standardization.

## Strategic Approach

### Solution Direction
- Treat this as three coordinated, low-code product changes inside the existing openspdd CLI corpus: (1) add an optional explore-mode command template that slots before analysis in the SPDD mental model, (2) realign Antigravity’s install path via `AIToolType.GetConfigDir()`, and (3) normalize frontmatter `name` fields to slash-free values across core and optional templates, updating tests/docs that assert the old slash-prefixed convention.
- Prefer extending the existing template hierarchy and detector mapping rather than introducing new packages, generation strategies, or workflow engines.

### Key Design Decisions
- `Explore placement = optional, not core`: keeps default installs lean and matches “stance/pre-analysis” nature; users who want discovery opt in via optional listing/generation → recommendation: keep optional (as in the follow-up commit).
- `Explore writes decision logs only`: absolute read-only for product code, with an explicit exception for `spdd/explore/*.md` drafts → recommendation: document the exception as the sole mutation surface to avoid policy contradiction.
- `Antigravity path → `.agents/workflows``: aligns with Antigravity workflow conventions, but shares the `.agents/` parent with Codex’s `.agents/skills` → recommendation: proceed with path change, and explicitly verify uninstall/init/detection do not assume exclusive ownership of `.agents/`.
- `Slash removal in source templates`: simplifies cross-tool metadata and matches Codex’s already-trimmed skill names; conflicts with GGQPA-006’s earlier preference to avoid global source edits → recommendation: accept the standardization only if lookup compatibility is preserved (ID match and/or slash-tolerant `GetByName`) and OpenCode/Codex adapters are revalidated.
- `No new Go runtime behavior for Explore Mode`: explore is prompt/policy content consumed by AI assistants after generation; openspdd’s job is to ship and install the markdown template correctly.

### Alternatives Considered
- `Ship Explore Mode as a core template`: rejected for default-install bloat and because explore is situational rather than mandatory pipeline step.
- `Keep Antigravity at `.antigravity/commands``: rejected by the requirement’s path-alignment goal; would leave Antigravity installs misaligned with expected workflows directory.
- `Only adapt slash at generation time (status quo from GGQPA-006)`: rejected for this initiative’s standardization goal, but remains a safety net for OpenCode; source slash removal and OpenCode strip-`name` can coexist.
- `Hard-enforce explore transition gates in CLI code`: rejected as out of scope — gates belong in the command prompt behavior, not in openspdd’s Go runtime.

## Risk & Gap Analysis

### Requirement Ambiguities
- `Read-only vs Decision Log contradiction`: “absolute read-only / never create project files” coexists with creating `spdd/explore/...` drafts — needs a single authoritative exception statement.
- `Antigravity signature vs config path divergence`: detection still keys on `.antigravity` signatures while config writes move to `.agents/workflows`; unclear whether detection/signature set should also expand to `.agents/workflows`.
- `Slash-prefixed lookup compatibility`: requirement strips `name: /spdd-*` to `name: spdd-*` but does not say whether `openspdd generate /spdd-analysis` (or tests using `GetByName("/spdd-analysis")`) must keep working.
- `Relationship to GGQPA-006`: earlier analysis rejected global slash removal; this initiative reverses that product decision without stating whether GGQPA-006’s OpenCode adapter remains necessary (it still is, because OpenCode should not emit `name` at all).
- `Explore artifact lifecycle`: no AC for when explore drafts are required, optional, overwritten, or promoted into `spdd/analysis/`.
- `Distorted Mirror ethics/UX`: technique deliberately injects unconfirmed assumptions without announcing the test — acceptable for power users, but may confuse newcomers; no guidance on frequency/safety boundaries beyond “every 3–4 interactions”.

### Edge Cases
- `Multi-tool repos with both Codex and Antigravity`: both use `.agents/...` children; init/uninstall must not delete sibling tool content.
- `Users with previously generated Antigravity files under `.antigravity/commands``: path change leaves stale commands unless migration/regeneration guidance exists.
- `Optional explore not installed`: users invoking `/spdd-explore` in an assistant that never received the optional template will fail; docs must clarify opt-in install.
- `Explore Decision Log vs analysis naming`: similar timestamped markdown patterns under `spdd/explore` vs `spdd/analysis` may confuse consumers of SPDD artifacts.
- `Tests asserting `name: /spdd-*` or `GetByName("/...")`**: will fail or become misleading after frontmatter standardization unless updated or made slash-tolerant.
- `Cursor/Claude installed copies in user projects`: regenerating templates updates target dirs; local `.cursor/commands` copies in this repo currently still carry slash-prefixed names and are out-of-band relative to embedded sources.

### Technical Risks
- `Regression in template lookup UX`: if `Name` loses the leading `/` and callers still pass `/spdd-*`, only ID match saves them when the argument equals `spdd-*` without slash; `/spdd-analysis` would no longer match `Name` and does not equal `ID` → mitigation: teach `GetByName` to trim a single leading `/` before comparison, or update all call sites/tests.
- `Antigravity uninstall/init blast radius under `.agents/``: shared parent with Codex skills increases accidental cleanup risk → mitigation: scope file operations strictly to `.agents/workflows`, never the entire `.agents` tree.
- `OpenCode/Codex adapter drift`: source `name` without slash changes assumptions in tests that assert presence/absence of specific frontmatter strings → mitigation: re-run OpenCode flat-strategy and Codex skill-strategy suites as required gates.
- `Detector still points at `.antigravity` while outputs go to `.agents/workflows``: projects without `.antigravity` marker but with workflows dir may not auto-detect Antigravity → mitigation: clarify whether signature files should include `.agents/workflows`.
- `Policy enforcement is soft`: explore read-only rules live in prompt text only; assistants can still violate them → mitigation: accept soft enforcement (consistent with other SPDD commands) and keep optional placement.

### Acceptance Criteria Coverage
| AC# | Description | Addressable? | Gaps/Notes |
|-----|-------------|--------------|------------|
| 1 | Add `spdd-explore` command template establishing Explore Mode (read-mostly thinking partner before analysis) | Yes | Content is prompt-level; no Go runtime needed beyond embedding/listing. |
| 2 | Explore Mode prohibits implementable code generation and product-file mutation | Partial | Soft prompt policy only; Decision Log write exception must be clarified. |
| 3 | Distorted Mirror Technique is embedded with rules and example | Yes | Prompt content; no automated verification possible in openspdd CLI. |
| 4 | Micro action segments M1–M4 and transition checkpoint to `/spdd-analysis` are defined | Yes | Prompt content; transition gates are conversational, not machine-checked. |
| 5 | `spdd-explore` is available as an optional (not default/core) template | Yes | Follow-up already places it under `data/optional/`; `ListAvailable`/`generate --all` will not include it by default. |
| 6 | Antigravity `GetConfigDir()` returns `.agents/workflows` | Yes | One-line detector mapping + test expectation update. |
| 7 | Core templates’ frontmatter `name` values have no leading `/` | Yes | Edit five core markdown sources; update dependent assertions. |
| 8 | Optional templates’ frontmatter `name` values have no leading `/` | Yes | Includes explore + existing optional set. |
| 9 | Existing tool generation for non-Antigravity tools remains functionally correct | Partial | Requires regression tests; slash-removal can break slash-prefixed `GetByName` callers. |
| 10 | OpenCode continues to avoid emitting conflicting `name` identity in generated commands | Yes | Existing `OpenCodeTemplateAdapter` still strips `name`; should remain. |
| 11 | Codex skill naming remains slash-free and stable | Yes | `TrimPrefix(meta.Name, "/")` stays valid whether source has slash or not. |
| 12 | Users with stale `.antigravity/commands` artifacts are guided to regenerate under `.agents/workflows` | Partial | Not stated in requirement; recommended docs/migration note. |
