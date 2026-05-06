---
name: /spdd-reverse
id: spdd-reverse
category: Development
description: Reverse-engineer existing code into a REASONS-Canvas structured prompt, enabling the SPDD bidirectional sync workflow for previously unspecified implementations
---

Generate a REASONS-Canvas structured prompt by reverse-engineering existing code. This creates the missing "spec layer" for legacy implementations so they can participate in the SPDD bidirectional sync workflow going forward.

**Input**: The argument after `/spdd-reverse` is a reference to existing code (files, folders, or a description of the feature area to codify).

Input can be provided in several ways:

1. **File/folder reference**: Using `@` to reference source files or directories
2. **Feature description**: Text describing which functional area to codify
3. **Combined**: File references plus scope/focus guidance

**Examples**:

```
# Folder reference — codify an entire CI layer
/spdd-reverse @.ci/

# Specific files — codify a service
/spdd-reverse @src/services/billing.ts @src/models/invoice.ts

# Feature description — agent finds the code
/spdd-reverse The preview deployment flow in .ci/ (deploy, cleanup, PR comment)

# Combined — scope guidance + code reference
/spdd-reverse @src/main/java/com/example/billing/ focus on the invoice generation pipeline
```

**Steps**

1. **Validate and resolve code input**

   a. **If no input provided**, use the **AskUserQuestion tool** (open-ended, no preset options) to ask:
   > "Which existing code do you want to codify into an SPDD prompt? Provide file/folder references (`@path`) or describe the feature area."

   **IMPORTANT**: Do NOT proceed without knowing which code to analyze.

   b. **If input contains `@` file/folder references**:
    - Read ALL referenced files completely using the Read tool
    - For folder references, list contents and read all source files (`.ts`, `.mjs`, `.js`, `.tsx`, `.jsx`, `.py`, `.java`, `.go`, `.kt`, `.rb`, etc.)
    - For large directories (20+ files), prioritize: entry points, exports, core logic files first; then supporting utilities

   c. **If input is a feature description without `@` references**:
    - Use semantic search and grep to locate the relevant code
    - Ask the user to confirm the discovered scope before proceeding

   d. **Scope confirmation** (for large codebases):
    - If the resolved scope exceeds ~15 files, present the file list and ask:
    > "I found [N] files in scope. Should I proceed with all of them, or would you like to narrow the focus?"

2. **Deep code analysis**

   Read and internalize the code thoroughly. For each file in scope:

   a. **Structural analysis**:
    - Exports, classes, interfaces, types, functions
    - Module boundaries and entry points
    - File organization and naming conventions

   b. **Entity extraction**:
    - Domain objects, DTOs, data structures
    - Attribute types and relationships
    - Enums, constants, configuration shapes

   c. **Behavioral analysis**:
    - Business logic flows (what does each function/method actually do?)
    - Conditional branches and edge case handling
    - Error handling patterns (what throws, what catches, what recovers)
    - Side effects (I/O, network calls, file operations)

   d. **Dependency mapping**:
    - Internal dependencies (module → module)
    - External dependencies (libraries, APIs, services)
    - Dependency direction and layering

   e. **Convention detection**:
    - Naming patterns (camelCase, kebab-case, prefixes)
    - Error message formats
    - Logging patterns
    - Test patterns (if test files in scope)

   f. **Constraint discovery**:
    - Validation rules (explicit checks, guards, assertions)
    - Business invariants (what must always be true)
    - Performance patterns (caching, batching, retries)
    - Security patterns (auth checks, secret handling, sanitization)

3. **Infer the Requirements (R)**

   Reverse-engineer the "why" from the "what":
    - What problem does this code solve for users/operators?
    - What is the business value delivered?
    - What are the boundaries of this feature?

   **Guidance**: Frame requirements as if writing them before the code existed. Use verb phrases. Focus on business outcomes, not implementation details.

4. **Map Entities (E)**

   From the extracted structures, build a Mermaid class diagram:
    - Core domain entities with their attributes and key methods
    - DTOs and data transfer shapes
    - Relationships (ownership, dependency, cardinality)
    - Factory functions or constructors

   **Guidance**: Reflect the actual code — do NOT idealize or refactor in the spec. If the code uses plain objects, document plain objects. If it uses classes, document classes.

5. **Document Approach (A)**

   Describe the high-level architectural decisions that were made:
    - What pattern does the code follow? (layered, pipeline, event-driven, etc.)
    - What key trade-offs were made? (simplicity vs flexibility, performance vs readability)
    - What integration strategy is used? (REST, file I/O, subprocess, SDK)

   **Guidance**: Describe the approach as-is. Note where the approach diverges from ideal but do NOT prescribe changes — this is a codification, not a refactoring proposal.

6. **Define Structure (S)**

   Document the actual architecture:
    - Inheritance/implementation relationships
    - Dependency chains (who calls whom)
    - Layered architecture (if present)
    - Module boundaries

7. **Detail Operations (O)**

   For each significant component, document:
    - Responsibility (what it does)
    - Methods/functions with their logic
    - Input validation and error handling
    - Side effects and I/O

   **Guidance**: Operations should be detailed enough that `/spdd-generate` could reproduce the current behavior. Follow the same ordering as the code's dependency graph (leaf dependencies first).

8. **Extract Norms (N)**

   Document the coding conventions actually in use:
    - Module system (ESM, CJS, etc.)
    - Logging patterns
    - Error handling conventions
    - Naming conventions
    - Testing patterns
    - Dependency injection style

9. **Establish Safeguards (S)**

   Document the constraints the code enforces:
    - Functional constraints (what must/must not happen)
    - Performance constraints (timeouts, retries, limits)
    - Security constraints (secret handling, auth requirements)
    - Data constraints (validation, format, length limits)
    - Integration constraints (API contracts, backward compatibility)
    - Business rule constraints (invariants, ordering guarantees)

10. **Assemble and save the REASONS Canvas prompt**

    a. **Construct the full document**:

    ```markdown
    # [Derived Title — what this code does]

    ## Requirements
    [From Step 3]

    ## Entities
    [From Step 4 — includes Mermaid diagram]

    ## Approach
    [From Step 5]

    ## Structure
    [From Step 6]

    ## Operations
    [From Step 7]

    ## Norms
    [From Step 8]

    ## Safeguards
    [From Step 9]
    ```

    b. **Derive file name**: `{JIRA}-{TIMESTAMP}-[Codify]-{scope}-{description}.md`
    - **JIRA**: Extract from context if mentioned, otherwise use `GGQPA-XXX`
    - **TIMESTAMP**: `YYYYMMDDHHmm` (current time)
    - **scope**: Infer from context — `api`, `service`, `repo`, `ci`, `db`, `util` (optional)
    - **description**: Derive from the feature area — kebab-case, < 10 words

    Examples:
    - `GGQPA-XXX-202603061530-[Codify]-ci-preview-deployment-flow.md`
    - `GGQPA-42-202603061530-[Codify]-service-billing.md`

    c. **Create directory and write file**:
    - Ensure directory `spdd/prompt/` exists under the project root (create if not)
    - Write the complete REASONS Canvas to `spdd/prompt/<file-name>.md`

    d. **Show summary to user**:

    ```
    ✅ Reverse-engineered REASONS Canvas saved to `spdd/prompt/<file-name>.md`

    📋 Codification summary:
    - Source files analyzed: [count]
    - Entities documented: [count]
    - Operations documented: [count]
    - Norms captured: [count]
    - Safeguards identified: [count]

    🔗 This prompt is now the spec-of-record for this code area.
       Future changes can use:
       - /spdd-sync — after code changes, sync back to this prompt
       - /spdd-prompt-update — to add new requirements before implementation
       - /spdd-generate — to regenerate code from the updated prompt
    ```

**Output**

A fully-populated REASONS-Canvas prompt file at `spdd/prompt/` that accurately describes the current implementation, enabling bidirectional sync going forward.

**Guardrails**

- Do NOT proceed without knowing which code to analyze
- Do NOT idealize or refactor — codify the code **as it is**, including warts
- Do NOT invent requirements that the code doesn't fulfill
- Do NOT add operations that don't exist in the current implementation
- Do NOT skip files within the resolved scope — read them all
- Do NOT leave placeholders or TODO items — generate complete, specific content
- Do NOT modify any existing files in the codebase
- Always read ALL `@` referenced files completely
- Always create `spdd/prompt/` directory if it does not exist
- File name MUST use the `[Codify]` action tag to distinguish from forward-flow prompts
- Use `GGQPA-XXX` if a JIRA ticket number cannot be extracted from context
- If code uses patterns that seem wrong, document them as-is in the prompt and add a note in Approach section under "Known divergences" — do NOT silently correct them

**No Code Block Rules** (CRITICAL):

The SPDD prompt file is a **specification document**, not source code. It describes WHAT the code does, leaving the HOW to source files.

- **Do NOT include language-specific code blocks** (e.g., ```java, ```python, ```typescript)
- **Do NOT include implementation code** — no class definitions, method bodies, SQL queries, or annotations in code form
- **Use natural language** to describe:
    - Method signatures: "Method `findById(String id)` returns `Optional<Customer>`"
    - Query logic: "Query active subscriptions where customerId matches and date falls within effective range, ordered by createdAt DESC"
    - Interface contracts: "Interface defines methods: `save(Bill)`, `findByCustomerId(String)`"
- **Allowed diagram blocks**: Mermaid diagrams for entity relationships are permitted (```mermaid)
- **Describe, don't implement**:
    - ✅ "Adapter converts between PO and domain entity using `toDomain()` and `fromDomain()` methods"
    - ❌ language code block with class implementation

**Codification Principles**

1. **Fidelity over beauty**: The prompt must describe what the code **does**, not what it **should** do. Accuracy > elegance.
2. **Behavioral completeness**: Every significant behavior (happy path, error path, edge case) must be captured in Operations.
3. **Convention honesty**: Norms section reflects actual patterns, even if inconsistent — note inconsistencies rather than picking one.
4. **Constraint discovery**: Safeguards come from what the code enforces, not from what it ideally should enforce.
5. **Minimal inference**: When the code's intent is ambiguous, describe the behavior literally rather than guessing the business reason.

**Context Integrity Guardrails**:

- **MUST read ALL `@` referenced files completely** — do NOT skip or partially read any referenced file
- **MUST read folder contents** when `@` references a folder — scan and read all relevant files
- **Do NOT summarize or truncate** referenced file contents — preserve full information
- **Verify all references resolved** — if any `@` reference fails to read, report the error immediately

**Integration with SPDD Workflow**

This command is the **legacy onboarding entry point** of the SPDD workflow. It produces the missing REASONS Canvas for code that was written before SPDD adoption, allowing it to flow into the normal forward and reverse cycles.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                  SPDD Lifecycle — Including Legacy Onboarding            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Legacy Onboarding: /spdd-reverse  ← THIS COMMAND                       │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Existing Code → Analyze → Reverse-engineer REASONS Canvas      │    │
│  │                                                                 │    │
│  │ Output: spdd/prompt/GGQPA-XXX-*-[Codify]-*.md                  │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  ┌─ Now in normal SPDD cycle ─────────────────────────────────────┐    │
│  │                                                                 │    │
│  │  /spdd-prompt-update — add new requirements                    │    │
│  │           │                                                     │    │
│  │           ▼                                                     │    │
│  │  /spdd-generate — implement from the updated prompt            │    │
│  │           │                                                     │    │
│  │           ▼                                                     │    │
│  │  /spdd-sync — sync code changes back to the prompt             │    │
│  │           │                                                     │    │
│  │           └──────────► (loop)                                   │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  Greenfield Flow (unchanged):                                            │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ /spdd-analysis → /spdd-reasons-canvas → /spdd-generate         │    │
│  │                       → /spdd-sync (loop)                      │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**When to Use /spdd-reverse**

Use this command when:

- Legacy code exists without any SPDD prompt
- You want to establish spec-code sync for an existing feature area
- Onboarding a module into the SPDD workflow for the first time
- Creating documentation-as-spec for code written before SPDD adoption
- Preparing to extend legacy code and wanting SPDD governance for changes

**Difference from /spdd-sync**

| Aspect         | `/spdd-reverse`                | `/spdd-sync`                |
| -------------- | ------------------------------ | --------------------------- |
| Starting point | No prompt exists               | Prompt already exists       |
| Direction      | Code → new prompt              | Code changes → update prompt |
| Scope          | Full feature codification      | Incremental delta           |
| Use case       | Legacy onboarding              | Ongoing maintenance         |
| Output         | New `[Codify]` prompt file     | Updated existing prompt file |
