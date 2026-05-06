---
name: /spdd-story
id: spdd-story
category: Development
description: Decompose high-level feature requirements into INVEST-compliant, business-focused stories with clear scope boundaries and testable acceptance criteria
---

Decompose high-level feature requirements into structured, business-focused development stories. Each story is independently deliverable, scoped to 1-5 days of work, and includes concrete acceptance criteria written in business language. Stories are split according to the INVEST principle and serve as input for the SPDD workflow (`/spdd-analysis` → `/spdd-reasons-canvas` → `/spdd-generate`). This command focuses on the **"What"** and **"For whom"** — technical analysis and implementation details are left to downstream SPDD phases.

**Input**: The argument after `/spdd-story` is a business requirement description or file reference.

Input can be provided in two ways:

1. **Text description**: Direct text describing the feature or requirement
2. **File/folder reference**: Using `@` to reference files or folders containing requirements

**Examples**:

```
# File reference
/spdd-story @requirements/team-management-feature.md

# Text description
/spdd-story Implement team management API with CRUD operations, member management, and role-based access control

# Combined
/spdd-story @requirements/team-management.md additionally needs invitation flow and audit logging
```

**Steps**

1. **Validate and consolidate business input**

   a. **If no input provided**, use the **AskUserQuestion tool** (open-ended, no preset options) to ask:
   > "Please provide the business requirement document or description for story generation (you can use text, @file references, or both)."

   **IMPORTANT**: Do NOT proceed without business input.

   b. **If input contains `@` file/folder references**:
    - Read ALL referenced files completely using the Read tool
    - For folder references, read all relevant files within the folder
    - Consolidate all file contents into a unified business context

   c. **Combine all context sources**:
    - Merge text descriptions with file contents
    - Preserve the complete information from all sources — do NOT summarize or truncate

   d. **Context Integrity Check**:
    - Verify all `@` references were successfully read
    - If any file cannot be read, report the error and ask user to provide alternative
    - Confirm the consolidated context contains sufficient information to proceed

2. **Scan existing stories for context**

   Do NOT explore the codebase — that is `/spdd-analysis`'s responsibility. Only gather enough context to avoid duplication and determine numbering.

   a. **List existing stories**:
    - List files in `requirements/` directory (if exists)
    - Read existing stories to understand numbering conventions and scope boundaries

   b. **Identify coverage gaps**:
    - Note which features and operations are already covered by existing stories
    - Ensure new stories do not duplicate or conflict with existing ones

   c. **Determine next available story number**:
    - Extract the highest `[User-story-N]` number from existing filenames
    - New stories start from `N+1`

2a. **Capture machine-derivable metadata for YAML frontmatter** (Docs-as-Code automation)

   These values are used in Step 8 to inject the canonical YAML frontmatter on every generated story file. Do NOT prompt the user for any of them — they MUST be derived automatically.

   a. **Capture BA identity** by running:
   ```
   git config user.name
   git config user.email
   ```
   - If `user.name` is empty, fall back to the value of `$USER` / `$USERNAME` from the environment.
   - If both are empty, set `assignees.ba` to `"unknown"` and emit a warning in the final summary instructing the user to run `git config --global user.name "Your Name"`.

   b. **Capture creation timestamp** in `YYYY-MM-DD HH:MM` form using the local system clock (`date +"%Y-%m-%d %H:%M"`). Use this single timestamp consistently across all stories generated in this run so they share a creation moment.

   c. **AI-derive `title`** from the consolidated business input collected in Step 1. The title MUST be:
    - A single line, ≤ 80 characters.
    - In the same natural language as the user input.
    - Free of trailing punctuation and Markdown.
   The user will see and may correct this in the final summary.

   **IMPORTANT**: All values captured here are inputs to Step 8's frontmatter injection — do NOT discuss them with the user mid-flow.

3. **Abstract task analysis and INVEST evaluation**

   Before creating any stories, analyze the feature at an abstract level.

   a. **Identify the abstract task**:

   ```
   ### Abstract Task: "[Feature Name]"

   **Analysis Dimensions**:
   - **Core Responsibility**: [What is the primary purpose of this feature]
   - **Primary Operations**: [List the main operations: create, query, update, delete, list, etc.]
   - **Key Constraints**: [Data uniqueness, permissions, associations, business rules]
   - **Technical Complexity**: [Low/Medium/High — standard CRUD, complex business logic, etc.]
   - **Business Complexity**: [Low/Medium/High — multiple query scenarios, cascading effects, etc.]
   ```

   b. **INVEST compliance check**:

   ```
   ### INVEST Evaluation:
   - ✅/❌ **Independent**: [Can this be developed, tested, and deployed independently?]
   - ✅/❌ **Negotiable**: [Can design details be discussed with the team?]
   - ✅/❌ **Valuable**: [Does this provide clear business value?]
   - ✅/❌ **Estimable**: [Can the team accurately estimate the effort?]
   - ✅/❌ **Small**: [Can it be completed in 1-2 sprints?]
   - ✅/❌ **Testable**: [Are there clear acceptance criteria?]

   **Conclusion**: [Needs splitting / Ready as-is]
   ```

   c. **If splitting is needed**, determine the split strategy:

   **Split dimensions** (choose the most appropriate):
    - **By operation type**: CREATE-READ / UPDATE-DELETE / LIST-SEARCH
    - **By complexity**: Basic features / Advanced features / Admin features
    - **By user role**: Regular user features / Admin features
    - **By technical dependency**: Core features / Extension features

   **Split rules**:
    - Each story contains at most 2-3 core functional points
    - Functional points within a story must have logical relevance
    - Single story workload: 1-5 days
    - Each story must deliver independent business value

4. **Determine story module numbering**

   a. **Scan existing stories** in `requirements/` to determine the next available module number
   b. **Assign module number**: `{MODULE_NUMBER}` (e.g., `003` for the third module)
   c. **Assign sequential story numbers** within the module: `001`, `002`, `003`, etc.

5. **Generate stories following the standard structure**

   For each story determined in the split analysis, generate the complete story using this structure:

   ---

   ### Story Title Format

   ```
   ## [STORY-{MODULE}-{SEQ}] {Operation Description} API Development
   ```

   ### Required Sections

   #### Section 1: Background

   ```markdown
   ### Background
   [Describe the business background, use cases, and role within the overall system]

   Key points:
   - Business value and user needs
   - Relationship with other features
   - Why this capability is needed now
   ```

   #### Section 2: Business Value

   ```markdown
   ### Business Value
   - Provide {specific capability} for {role}
   - Support {specific need} in {business scenario}
   - Enable {key function} of {system goal}
   ```

   #### Section 3: Dependencies and Assumptions

   ```markdown
   ### Dependencies and Assumptions
   - **Prerequisites**: [Features or stories that must be completed first, if any]
   - **Data assumptions**: [What data or entities are expected to already exist]
   - **Integration points**: [External systems, APIs, or services this story interacts with]
   - **Business constraints**: [Regulatory, contractual, or organizational constraints]
   ```

   #### Section 4: Scope In / Scope Out

   ```markdown
   ### Scope In
   - [Feature included in this story, bullet points]

   ### Scope Out
   - [Feature NOT included in this story, bullet points]
   ```

   #### Section 5: Acceptance Criteria (ACs)

   Generate business-focused ACs using Given-When-Then format. The number of ACs depends on the business scenarios — do NOT use a fixed number. Cover these categories as appropriate:

   **a. Happy path ACs** — Core business scenarios with concrete examples:
   ```markdown
   #### AC{N}: {Business Scenario Description}
   **Given** {business precondition with concrete values/examples}
   **When** {user action in business language}
   **Then** {expected business outcome with specific numbers}
   ```

   **b. Validation and business rule ACs** — Input validation and business constraints, expressed as user-facing behavior:
   ```markdown
   #### AC{N}: {Validation Scenario}
   **Given** {invalid or edge-case input condition}
   **When** {user attempts the action}
   **Then** {system rejects with clear user-facing message}
   ```

   **c. Error condition ACs** — Expected failure scenarios from the user's perspective:
   ```markdown
   #### AC{N}: {Error Scenario}
   **Given** {condition that causes failure, e.g., resource does not exist}
   **When** {user attempts the action}
   **Then** {system responds with appropriate error and HTTP status}
   ```

   **d. Non-functional expectations** (optional, only if relevant to the business):
   ```markdown
   #### Non-Functional Expectations
   - [Business-level quality expectation, e.g., "Billing calculation must complete fast enough for real-time API response"]
   - [Business-level quality expectation, e.g., "System must handle concurrent usage submissions from the same customer correctly"]
   ```

   **AC writing guidelines**:
    - Use **business language**, not implementation language (say "system rejects the request" not "return HTTP 400")
    - Include **concrete numbers and examples** (say "100,000 monthly quota, 80,000 used" not "some quota, some usage")
    - Exception: HTTP status codes ARE acceptable in ACs since they are part of the API contract visible to consumers
    - Each AC must be **independently testable** by a QA engineer without reading source code
    - Do NOT prescribe HOW to implement (no "use parameterized queries", "apply cache strategy", etc.)
    - Do NOT specify internal technical details (no JSON response format, no DB schema, no specific error codes)

6. **Quality check each generated story**

   After generating each story, verify against this checklist:

   **Structure and Completeness**:
    - [ ] Contains all required sections (Background, Business Value, Dependencies and Assumptions, Scope In/Out, ACs)
    - [ ] Each AC uses Given-When-Then format with concrete values and examples
    - [ ] ACs are written in business language — no implementation details leaked in
    - [ ] ACs cover happy path, validation/business rules, and error conditions

   **Business Clarity**:
    - [ ] Business value is clear and stated for a specific audience/role
    - [ ] Scope In and Scope Out clearly delineate boundaries — no ambiguous overlap
    - [ ] No duplication with other existing stories
    - [ ] A QA engineer could write test cases from the ACs without reading source code

   **Sizing and Independence**:
    - [ ] Story contains at most 2-3 core functional points
    - [ ] Story can be developed and delivered independently
    - [ ] Estimated workload is 1-5 days

7. **Final INVEST re-validation**

   After all stories are generated, perform a final check:

   **For each story**:
    - [ ] **Independent**: Does not depend on other stories' implementation details
    - [ ] **Complete**: Contains a complete functional loop and acceptance criteria
    - [ ] **Valuable**: Delivers independent business value when implemented alone
    - [ ] **Estimable**: Development team can accurately estimate development time
    - [ ] **Right-sized**: 1-5 day workload, at most 3 core functional points
    - [ ] **Testable**: Has clear test scenarios and acceptance conditions

   **Anti-patterns to avoid**:
    - Do NOT split by technical layer (frontend/backend/database for the same feature)
    - Do NOT over-fragment a single API into multiple stories
    - Do NOT break business logic completeness
    - Do NOT create complex inter-story dependencies

8. **Assemble and save the story document(s) with YAML frontmatter**

   a. **Inject canonical YAML frontmatter** at the very top of every generated story file. This is the **single source of truth** for downstream automation (`/spdd-analysis`, `/spdd-generate`, `/spdd-prompt-update`, `/spdd-block`, CI hooks). The schema is fixed — do NOT add, rename, reorder, or omit fields. Use `null` for unknown values; do NOT use empty strings unless the schema explicitly says so.

   ```yaml
   ---
   id: STORY-{MODULE}-{SEQ}              # from Step 4 — primary key for downstream tools
   title: "{ai-derived title from Step 2a.c}"
   status: BACKLOG                        # one of BACKLOG | IN_PROGRESS | IN_TEST | BLOCKED | DONE
   iteration: null                        # filled in by humans during sprint planning
   tags: []                               # may be left empty; humans curate later
   assignees:
     ba: "{git user.name from Step 2a.a}"
     qa: null                             # populated by /spdd-analysis on first run
     dev: null                            # populated by /spdd-generate on first run
   timestamps:
     created_at: "{YYYY-MM-DD HH:MM from Step 2a.b}"
     analyzed_at: null                    # populated by /spdd-analysis
     developed_at: null                   # populated by Git pre-push / PR open hook
     tested_at: null                      # populated by CI on QA approval
     delivered_at: null                   # populated by CI on merge to main
   quality_metrics:
     ai_refine_loops: 0                   # incremented by /spdd-generate and /spdd-prompt-update
     qa_rejections: 0                     # incremented by CI on PR "changes requested"
     test_coverage: null                  # written by CI after unit tests
   block_records: []                      # appended/closed by /spdd-block & /spdd-unblock
   inline_defects: []                     # human-edited (use IDE snippet)
   related_bugs: []                       # human-edited (use IDE snippet)
   ---
   ```

   **Per-file rules**:
    - When the run produces multiple story files, each file gets its own frontmatter with its own `id`. The `created_at` is shared (same run).
    - Preserve key order exactly as shown above so diffs stay readable across automation passes.
    - YAML strings containing `:`, `#`, or leading `-` MUST be double-quoted.

   b. **Construct the complete document body** (placed immediately after the closing `---` of the frontmatter):

   ```markdown
   # Story Decomposition: [Feature Name]

   ## INVEST Analysis

   ### Abstract Task
   [Output from Step 3a]

   ### INVEST Evaluation
   [Output from Step 3b]

   ### Split Strategy
   [Output from Step 3c — if splitting was needed]

   ---

   [All generated stories in sequence]
   ```

   c. **Derive file name**: `[User-story-{N}]{kebab-case-title}.md`
    - **N**: Next available story number based on existing files in `requirements/`
    - **title**: Descriptive kebab-case title derived from the feature

   If the feature splits into multiple stories, generate ONE file per story:
    - `[User-story-{N}]{story-1-title}.md`
    - `[User-story-{N+1}]{story-2-title}.md`
    - etc.

   Alternatively, if the stories are closely related and part of a single feature decomposition, generate a single consolidated file. In that case the consolidated file still carries a single frontmatter block with a single `id` representing the umbrella story; per-section IDs live inside the body as `## [STORY-…]` headers.

   d. **Create directory and write file(s)**:
    - Ensure directory `requirements/` exists under the project root (create if not)
    - Write the complete story document(s) (frontmatter + body) to `requirements/<file-name>.md`

   e. **Show summary to user** (and surface any metadata fallbacks from Step 2a):

   ```
   ✅ Story generation complete. Stories saved to `requirements/`

   📋 Generation summary:
   - Feature: [feature name]
   - Stories generated: [count]
   - Total ACs: [count]
   - Estimated total effort: [X-Y days]
   - BA assignee (auto from git): [git user.name OR ⚠️ unknown — run `git config --global user.name "Your Name"`]
   - Created at: [YYYY-MM-DD HH:MM]
   - Status: BACKLOG (use /spdd-analysis to advance)

   📝 Stories:
   1. [STORY-XXX-001] [title] — [estimated effort]
   2. [STORY-XXX-002] [title] — [estimated effort]
   ...

   🔗 Next step: Process each story through the SPDD workflow:
      /spdd-analysis @requirements/<file-name>.md
   ```

9. **Offer to proceed with SPDD analysis**

   > "Stories are ready. Would you like me to proceed with `/spdd-analysis` for any of the generated stories?"

   If the user selects a story, invoke the `/spdd-analysis` workflow with that story file as input.

**Output**

Structured, INVEST-compliant story document(s) saved to `requirements/`, containing:
- INVEST analysis with abstract task identification and split strategy
- Complete stories with Background, Business Value, Dependencies and Assumptions, Scope In/Out, and business-focused Acceptance Criteria
- Quality-checked against business clarity, sizing, and independence checklists

**Guardrails**

- Do NOT proceed without business requirement input
- Do NOT explore the codebase — codebase analysis is `/spdd-analysis`'s responsibility
- Do NOT summarize or truncate the original business requirement — preserve it verbatim when referenced
- Do NOT generate code — this command produces stories only
- Do NOT leave placeholders or TODO items — generate complete, specific content
- Do NOT modify any existing files in the codebase
- Do NOT split stories by technical layer (frontend/backend/database)
- Do NOT create stories with complex inter-story dependencies
- Do NOT over-fragment a single API endpoint into multiple stories
- Do NOT specify HOW to implement — only WHAT the expected behavior is
- Do NOT prescribe technical solutions in ACs (no caching strategies, indexing approaches, query patterns, error JSON formats, P95/P99 metrics)
- Do NOT include security implementation details in ACs (no "use parameterized queries", "sanitize HTML", etc.) — those are `/spdd-reasons-canvas` Safeguards
- ACs MUST use business language that a QA engineer or Product Owner can understand
- ACs MUST include concrete numbers, examples, and expected outcomes
- Each story MUST pass INVEST compliance checks
- Each story's workload MUST be between 1-5 days
- Each story MUST contain at most 3 core functional points
- Always read ALL `@` referenced files completely
- Always create `requirements/` directory if it does not exist
- File name MUST follow the naming convention: `[User-story-{N}]{kebab-case-title}.md`
- Every generated story file MUST start with the canonical YAML frontmatter defined in Step 8a — no field may be added, renamed, reordered, or omitted
- `assignees.ba` MUST come from `git config user.name`; do NOT prompt the user for it
- `created_at` MUST come from the local system clock at run time; do NOT make it up
- `status` MUST be initialized to `BACKLOG`; never set a later status from this command
- All counters in `quality_metrics` MUST be initialized to `0` (or `null` for `test_coverage`); do NOT pre-populate

**Context Integrity Guardrails**:

- **MUST read ALL `@` referenced files completely** — do NOT skip or partially read any referenced file
- **MUST read folder contents** when `@` references a folder — scan and read all relevant files
- **Do NOT summarize or truncate** referenced file contents — preserve full information
- **Verify all references resolved** — if any `@` reference fails to read, report error immediately
- **Combine all sources** — merge text descriptions with file contents into unified context
- **Preserve original intent** — do not interpret or modify the meaning of provided context

**Integration with SPDD Workflow**

This command is the **story decomposition phase** of the SPDD workflow, transforming high-level feature requirements into structured, implementable stories that feed into the SPDD pipeline:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           SPDD Workflow                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Pre-Phase: /spdd-story    ← YOU ARE HERE                              │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ High-Level Feature Requirement                                  │    │
│  │   + INVEST Analysis (abstract task, compliance, split strategy) │    │
│  │   + Story Decomposition (2-3 functional points per story)       │    │
│  │   + Business-Focused ACs (Given-When-Then with concrete data)   │    │
│  │   = Structured, Business-Ready Stories                           │    │
│  │                                                                  │    │
│  │ Output: requirements/[User-story-N]*.md                         │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                    ┌─────────┴─────────┐                                │
│                    │  For each story:   │                                │
│                    └─────────┬─────────┘                                │
│                              ▼                                          │
│  Phase 0: /spdd-analysis                                                │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Story → Enriched Context (domain concepts + strategy + risks)  │    │
│  │                                                                 │    │
│  │ Output: spdd/analysis/SPDD-XXX-*-[Analysis]-*.md              │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  Phase 1: /spdd-reasons-canvas                                         │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Enriched Context → REASONS Canvas Structured Prompt             │    │
│  │                                                                 │    │
│  │ Output: spdd/prompt/SPDD-XXX-*.md (REASONS Canvas)           │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  Phase 2: /spdd-generate                                               │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Structured Prompt → Validate → Generate → Verify → Code        │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  Phase 3: /spdd-api-test                                               │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Generated Code + ACs → API Test Script                          │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  Phase 4: /spdd-sync                                                   │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Code Changes → Analyze → Update Prompt → Consistency           │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**Why This Phase Matters**

High-level feature requirements often describe complex capabilities that cannot be implemented in a single sprint. `/spdd-story` bridges this gap by:

1. **Right-sizing work**: INVEST analysis ensures each story is independently deliverable within 1-5 days
2. **Business clarity**: ACs use concrete examples and business language that POs, QA, and developers all understand
3. **Reducing rework**: Stories with clear scope boundaries and acceptance criteria minimize ambiguity during implementation
4. **Enabling parallelism**: Independent stories can be developed concurrently by different team members
5. **Feeding the SPDD pipeline**: Each generated story is immediately ready for `/spdd-analysis`, creating a smooth end-to-end workflow

**Separation of Concerns across the SPDD Pipeline**:

| Concern | `/spdd-story` | `/spdd-analysis` | `/spdd-reasons-canvas` |
|---------|---------------|-------------------|------------------------|
| Audience | PO, Scrum Master, QA | Dev Team, Architect | Dev Team |
| Language | Business language | Strategic / conceptual | Technical / implementation |
| Focus | "What to build & for whom" | "What & Why in codebase context" | "How to build it" |
| Scope | Feature decomposition, story boundaries | Domain concepts, design direction, risks | Entity models, architecture, operations |
| ACs | Business scenarios with concrete data | AC coverage & gap analysis | Safeguards with specific metrics |
| Technical depth | None — out of scope | Strategic (trade-offs, approach) | Full (code structure, error formats, P99) |
| Codebase access | None (only reads `requirements/`) | Deep concept-driven exploration | Full exploration for implementation |
| Output | `requirements/[User-story-N]*.md` | `spdd/analysis/*-[Analysis]-*.md` | `spdd/prompt/*.md` (REASONS Canvas) |
