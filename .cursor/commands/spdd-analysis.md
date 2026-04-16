---
name: /spdd-analysis
id: spdd-analysis
category: Development
description: Analyze business requirements against codebase context at a strategic level, producing enriched context (business + domain concepts + strategic direction + risks) for REASONS Canvas generation
---

Analyze a business requirement document against the current codebase, producing a **strategic-level** enriched context that combines business information, domain concept identification, high-level approach decisions, and risk analysis — serving as high-quality input for `/spdd-reasons-canvas`. This phase focuses on the "What" and "Why", leaving the "How" to the REASONS Canvas phase.

**Input**: The argument after `/spdd-analysis` is a business requirement description or file reference.

Input can be provided in two ways:

1. **Text description**: Direct text describing the requirement
2. **File/folder reference**: Using `@` to reference files or folders containing requirements

**Examples**:

```
# File reference
/spdd-analysis @requirements/token-usage-billing-story.md

# Text description
/spdd-analysis Implement monthly billing summary report for customers with usage breakdown

# Combined
/spdd-analysis @requirements/billing-report.md additionally needs CSV export support
```

**Steps**

1. **Validate and consolidate business input**

   a. **If no input provided**, use the **AskUserQuestion tool** (open-ended, no preset options) to ask:
   > "Please provide the business requirement document or description (you can use text, @file references, or both)."

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

2. **Concept-driven codebase exploration**

   Do NOT exhaustively read the entire codebase — this does not scale. Instead, use a **concept-driven** approach: first build a lightweight project fingerprint, then extract search concepts from the business requirement, and finally explore only the relevant parts of the codebase in depth.

   a. **Project fingerprint (lightweight bootstrap — always do first)**:
    - Read the **primary** build/dependency file (ONE file: e.g., `build.gradle`, `pom.xml`, `package.json`, `requirements.txt`, `go.mod`) to detect the tech stack, framework, and key dependencies
    - List the top-level directory structure (directory names only, not file contents) to understand the project layout and layering conventions
    - Read the **main** configuration file (e.g., `application.yml`, `.env`, `next.config.js`) to understand infrastructure choices (DB, caching, messaging, etc.)
    - This step should be fast and touch only 2–3 files

   b. **Extract search concepts from business input**:
   Before touching any domain code, analyze the business requirement from Step 1 to extract:
    - **Domain nouns**: entity/concept names that likely map to code (e.g., "customer", "bill", "pricing plan", "subscription", "quota")
    - **Action verbs**: operations that likely map to endpoints or services (e.g., "submit usage", "calculate bill")
    - **API surfaces**: explicit paths, event names, or queue names mentioned (e.g., `POST /api/usage`)
    - **Technical hints**: mentioned technologies, patterns, or domain-specific terms (e.g., "monthly reset", "overage rate", "token")

   These extracted concepts become the **search scope** for all subsequent exploration.

   c. **Targeted schema exploration (scoped by concepts)**:
    - Search migration/schema files for tables whose names match the extracted domain nouns — do NOT read all migrations
    - Read ONLY the matched migrations or schema definitions
    - Follow foreign key relationships **one hop outward** from matched tables to capture boundary context (e.g., if `bills` references `customers`, read the `customers` table definition too)
    - If using an ORM, search for entity/model classes matching the same concept names

   d. **Targeted code exploration (scoped by concepts)**:
    - Search file names and class names for matches against the extracted concepts (e.g., `*Customer*`, `*Bill*`, `*Usage*`, `*Quota*`)
    - Read matched files to understand existing business logic, validation, and error handling
    - Follow direct code dependencies **one hop** (e.g., if `BillService` injects `SubscriptionRepository`, read that too — but don't keep chaining)
    - From the matched files, observe architecture conventions in use (naming, layering, error format, validation style, test patterns)
    - If no matching domain code exists (greenfield area), note this explicitly and rely on framework conventions inferred from dependencies

   e. **Relevant SPDD context (scoped by concepts)**:
    - List files in `spdd/prompt/` and `spdd/analysis/` (if the directories exist)
    - Read ONLY those files whose filenames suggest relevance to the extracted concepts
    - If none are relevant or the directories are absent, skip this step

   f. **Controlled expansion (one additional hop only)**:
    - If during steps 2c–2e you discover a concept that is clearly essential to the requirement but was NOT in the initial extraction (e.g., an unexpected foreign key, a shared utility), add it to the concept list and do **one more** targeted search for it
    - Do NOT recursively expand beyond this single additional hop — stop and note the boundary

   **IMPORTANT**: Be targeted — explore deeply within the relevant scope, not broadly across the entire codebase. Read actual file contents for the scoped concepts; do not guess. If the scope turns out to be very large (e.g., the requirement touches 10+ existing modules), explicitly list all identified concepts and prioritize the core ones, noting the peripheral ones as "boundary context" to be verified during REASONS Canvas.

3. **Domain Concept Identification**

   Identify the business concepts involved at a **conceptual level** — do NOT drill into specific attributes, data types, method signatures, or DTOs. The goal is to understand the domain landscape, not to design the implementation.

   a. **Concept inventory**:
    - What core business concepts does this requirement involve?
    - Which already exist in the codebase (as DB tables, existing classes, or type definitions)?
    - Which are new and need to be introduced?

   b. **Conceptual relationships**:
    - How do these concepts relate to each other at a business level?
    - What are the ownership and lifecycle boundaries?

   c. **Key business rules**:
    - What invariants must be maintained?
    - What business rules are explicit in the requirement?
    - What business rules are **implicit** and need to be surfaced?

   Output this section as:

   ```
   ### Domain Concept Identification

   #### Existing Concepts (from codebase)
   - [ConceptName]: [business purpose] — [relationship to other concepts]

   #### New Concepts Required
   - [ConceptName]: [business purpose] — [how it relates to existing concepts]

   #### Key Business Rules
   - [Rule]: [which concepts it governs]
   ```

4. **Strategic Approach & Trade-offs**

   Determine the **high-level solution direction** — do NOT specify implementation details like specific queries, annotations, JSON shapes, method signatures, or step-by-step logic. Those belong in the REASONS Canvas phase.

   a. **Solution direction**:
    - What is the overall approach to solving this requirement?
    - Which existing architectural patterns and conventions should be leveraged?
    - What is the general data flow direction (e.g., "REST endpoint → service-layer calculation → persist result")?

   b. **Key design decisions**:
    - What strategic choices need to be made?
    - What are the trade-offs for each choice?
    - What is the recommended direction and why?

   c. **Alternatives considered** (if applicable):
    - What other approaches were considered?
    - Why were they rejected?

   Output this section as:

   ```
   ### Strategic Approach

   #### Solution Direction
   - [High-level description of approach]

   #### Key Design Decisions
   - [Decision]: [trade-offs] → [recommendation and rationale]

   #### Alternatives Considered
   - [Alternative]: [why rejected]
   ```

5. **Risk & Gap Analysis**

   Surface everything that could cause problems or needs clarification **before** detailed design begins in the REASONS Canvas phase.

   a. **Requirement ambiguities**:
    - What is unclear, underspecified, or open to interpretation in the requirement?
    - What implicit assumptions has the requirement made?

   b. **Edge cases**:
    - What scenarios are not explicitly addressed by the requirement or ACs?
    - What boundary conditions need clarification?

   c. **Technical risks**:
    - What technical challenges or constraints could impact the implementation?
    - Are there concurrency, performance, or data integrity concerns?

   d. **Acceptance Criteria coverage**:
    - Are all ACs addressable with the proposed approach?
    - Are there gaps between the ACs and the full scope of the requirement?

   Output this section as:

   ```
   ### Risk & Gap Analysis

   #### Requirement Ambiguities
   - [Ambiguity]: [what needs clarification]

   #### Edge Cases
   - [Scenario]: [why it matters]

   #### Technical Risks
   - [Risk]: [potential impact and mitigation direction]

   #### Acceptance Criteria Coverage
   | AC# | Description | Addressable? | Gaps/Notes |
   |-----|-------------|--------------|------------|
   | [n] | [AC text]   | Yes/Partial  | [any gaps]  |
   ```

6. **Assemble the enriched context document**

   Combine all analysis results into a single, structured document:

   ```markdown
   # SPDD Analysis: [Derived Title]

   ## Original Business Requirement
   [Complete original requirement text — unmodified]

   ## Domain Concept Identification
   [Output from Step 3]

   ## Strategic Approach
   [Output from Step 4]

   ## Risk & Gap Analysis
   [Output from Step 5]
   ```

   **NOTE**: The codebase exploration from Step 2 is a **working process** — its findings are internalized and reflected through the Domain Concept Identification (which references existing vs. new concepts), Strategic Approach (which references existing patterns and conventions), and Risk & Gap Analysis (which surfaces technical constraints). Do NOT output a separate "Codebase Context Summary" section.

   **IMPORTANT**:
    - The original business requirement MUST be included verbatim — do NOT paraphrase
    - Every section must contain concrete, specific content — no placeholders
    - All analysis must be grounded in actual codebase exploration, not assumptions
    - Stay at a **conceptual/strategic** level — do NOT include implementation details (specific queries, JSON shapes, method signatures, annotations, component inventories). Those belong in the REASONS Canvas phase.

7. **Save the enriched context document**

   a. **Derive file name**: `{JIRA}-{TIMESTAMP}-[Analysis]-{description}.md`
    - **JIRA**: Extract from business context if mentioned, otherwise use `GGQPA-XXX`
    - **TIMESTAMP**: `YYYYMMDDHHmm` (current time)
    - **description**: Derive from business context — kebab-case, < 10 words

   Examples:
    - `GGQPA-XXX-202603131530-[Analysis]-token-usage-billing.md`
    - `GGQPA-169-202603131530-[Analysis]-monthly-report-export.md`

   b. **Create directory and write file**:
    - Ensure directory `spdd/analysis/` exists under the project root (create if not)
    - Write the complete enriched context document to `spdd/analysis/<file-name>.md`

   c. **Show summary to user**:

   ```
   ✅ Analysis complete. Enriched context saved to `spdd/analysis/<file-name>.md`

   📋 Analysis summary:
   - Project type: [backend/frontend/fullstack]
   - Existing concepts identified: [count]
   - New concepts required: [count]
   - Key design decisions: [count]
   - Acceptance Criteria coverage: [count]/[total]
   - Open questions/risks: [count]

   🔗 Next step: Use this as input for REASONS Canvas generation:
      /spdd-reasons-canvas @spdd/analysis/<file-name>.md
   ```

8. **Offer to proceed with REASONS Canvas generation**

   > "The enriched context is ready. Would you like me to proceed with `/spdd-reasons-canvas` using this analysis as input?"

   If the user confirms, invoke the `/spdd-reasons-canvas` workflow with the saved analysis file as input.

**Output**

An enriched context document saved to `spdd/analysis/<file-name>.md` that transforms raw business requirements into a **strategic-level** analysis containing:
- Original business requirements (preserved verbatim)
- Domain concept identification (existing and new concepts, conceptual relationships, business rules — grounded in codebase exploration)
- Strategic approach (solution direction, key design decisions, trade-offs, alternatives considered)
- Risk & gap analysis (ambiguities, edge cases, technical risks, AC coverage assessment)

**Guardrails**

- Do NOT proceed without business requirement input
- Do NOT skip codebase exploration — analysis MUST be grounded in actual codebase state
- Do NOT exhaustively read the entire codebase — use concept-driven scoping from the business requirement to target only relevant areas
- Do NOT summarize or truncate the original business requirement — preserve it verbatim
- Do NOT make assumptions about codebase structure without reading actual files
- Do NOT hardcode stack-specific terminology — always detect the project type first and adapt language to the discovered stack
- Do NOT generate code — this command produces analysis only
- Do NOT include implementation-level details (specific queries, JSON shapes, method signatures, annotations, component-layer inventories, step-by-step logic) — those belong in `/spdd-reasons-canvas`
- Do NOT leave placeholders or TODO items — generate complete, specific content
- Do NOT modify any existing files in the codebase
- Always read ALL `@` referenced files completely
- Always create `spdd/analysis/` directory if it does not exist
- File name MUST follow the naming convention defined above
- Use `GGQPA-XXX` if JIRA ticket number cannot be extracted from context
- Acceptance Criteria coverage MUST assess every AC from the requirement
- Risk & Gap Analysis MUST surface any ambiguities — do NOT silently assume

**Context Integrity Guardrails**:

- **MUST read ALL `@` referenced files completely** — do NOT skip or partially read any referenced file
- **MUST read folder contents** when `@` references a folder — scan and read all relevant files
- **Do NOT summarize or truncate** referenced file contents — preserve full information
- **Verify all references resolved** — if any `@` reference fails to read, report error immediately
- **Combine all sources** — merge text descriptions with file contents into unified context
- **Preserve original intent** — do not interpret or modify the meaning of provided context

**Integration with SPDD Workflow**

This command is the **pre-processing phase** of the SPDD workflow, bridging raw business requirements to implementation-ready structured prompts:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           SPDD Workflow                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Phase 0: /spdd-analysis                                                │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Business Requirement                                            │    │
│  │   + Concept-driven Codebase Exploration (targeted, not full)    │    │
│  │   + Domain Concept Identification (conceptual, not detailed)    │    │
│  │   + Strategic Approach & Trade-offs (direction, not specifics)  │    │
│  │   + Risk & Gap Analysis (ambiguities, edge cases, risks)        │    │
│  │   = Enriched Context (Business + Strategic + Risks)             │    │
│  │                                                                 │    │
│  │ Output: spdd/analysis/GGQPA-XXX-*-[Analysis]-*.md              │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  Phase 1: /spdd-reasons-canvas                                         │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Enriched Context → REASONS Canvas Structured Prompt             │    │
│  │                                                                 │    │
│  │ Output: spdd/prompt/GGQPA-XXX-*.md (REASONS Canvas)           │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  Phase 2: /spdd-generate                                               │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Structured Prompt → Validate → Generate → Verify → Code        │    │
│  │                                                                 │    │
│  │ Output: Implementation code following Operations sequence       │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  Phase 3: /spdd-sync                                                   │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Code Changes → Analyze → Update Prompt → Consistency           │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**Why This Phase Matters**

Raw business requirements describe **what** to build but lack the technical context needed to produce high-quality REASONS Canvas prompts. `/spdd-analysis` bridges this gap by:

1. **Grounding in reality**: Analysis is based on actual codebase state, explored through concept-driven scoping rather than exhaustive reading
2. **Surfacing hidden complexity**: Business rules, edge cases, and ambiguities that are implicit in the requirement become explicit
3. **Setting strategic direction**: Key design decisions and trade-offs are resolved before detailed design begins
4. **Reducing hallucination risk**: By feeding `/spdd-reasons-canvas` enriched context with real codebase data, the generated REASONS Canvas is more accurate and implementable
5. **Identifying risks early**: Open questions and ambiguities are surfaced before design, not during implementation

**Separation of Concerns with REASONS Canvas**:

| Concern | `/spdd-analysis` (this phase) | `/spdd-reasons-canvas` (next phase) |
|---------|-------------------------------|--------------------------------------|
| Thinking level | Strategic — "What" & "Why" | Tactical — "How" |
| Domain | Conceptual identification | Detailed entity modeling (E) |
| Solution | Direction & trade-offs | Concrete design & architecture (A, S) |
| Implementation | Out of scope | Specific operations & tasks (O) |
| Standards | Out of scope | Coding norms & safeguards (N, S) |
| Risks | Identify & surface | Resolve via constraints |
