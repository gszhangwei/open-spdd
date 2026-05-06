---
name: /spdd-generate
id: spdd-generate
category: Development
description: Generate code from a structured SPDD prompt file following the REASONS Canvas methodology
---

Generate implementation code from a structured SPDD (Structured Prompt-Driven Development) prompt file, strictly following the Operations sequence and coding norms defined in the prompt.

**Input**: The argument after `/spdd-generate` is the path to the structured prompt file (e.g., `@spdd/prompt/SPDD-XXX-202602271430-[Feat]-api-create-agent-endpoint.md`).

**Steps**

1. **If no input provided, ask for the prompt file**

   Use the **AskUserQuestion tool** to ask:

   > "Please provide the path to the structured prompt file (e.g., `@spdd/prompt/xxx.md`)."

   **IMPORTANT**: Do NOT proceed without a valid prompt file path.

2. **Read and parse the structured prompt file**

   Read the prompt file and extract the REASONS Canvas sections:

   | Section              | Purpose                        | Usage                             |
   | -------------------- | ------------------------------ | --------------------------------- |
   | **R** - Requirements | Overall goal and DoD           | Understand the business context   |
   | **E** - Entities     | Domain model and relationships | Reference for class design        |
   | **A** - Approach     | Implementation strategy        | Guide architectural decisions     |
   | **S** - Structure    | Components and dependencies    | Verify layering and relationships |
   | **O** - Operations   | Concrete implementation tasks  | **Execute in defined order**      |
   | **N** - Norms        | Engineering standards          | Apply to all generated code       |
   | **S** - Safeguards   | Non-negotiable constraints     | Enforce strictly                  |

   **IMPORTANT**: Read the ENTIRE file carefully. Each section provides critical guidance.

3. **Analyze project context**

   Before generating code:
   - Identify the project's technology stack (e.g., Spring Boot, Java version)
   - Locate existing similar patterns in the codebase for reference
   - Identify the correct package structure and directory layout
   - Check for existing base classes, utilities, or configurations to reuse

   **IMPORTANT**: Generated code MUST align with existing project conventions.

4. **Validate the Operations sequence**

   Review the **Operations** section to verify:

   a. **Dependency order is correct**:
   - Classes with no dependencies come first (enums, constants)
   - Classes depend only on previously defined classes
   - No circular dependencies exist

   b. **Task decomposition is complete**:
   - Each operation is atomic and testable
   - No logical gaps between operations
   - All components mentioned in Structure are covered

   c. **Consistency with Structure section**:
   - Inheritance relationships match
   - Dependencies match
   - Layered architecture is respected

   **If issues are found**: Report to user and suggest prompt modifications before proceeding.

   **IMPORTANT**: Do NOT re-plan the sequence. The Operations order is the designed execution order from the Abstraction phase.

5. **Generate code following Operations sequence**

   For each operation in the **Operations** section (in order):

   a. **Read the operation specification**:
   - Responsibility: What the component does
   - Attributes/Methods: Exact fields and signatures
   - Annotations: Required annotations
   - Validation rules: Bean validation or custom logic
   - Business logic: Step-by-step implementation details

   b. **Apply Norms**:
   - Annotation standards (e.g., @RestController, @Service)
   - Dependency injection style (constructor injection)
   - Exception handling patterns
   - Logging conventions
   - Response format standards

   c. **Enforce Safeguards**:
   - Field validation constraints
   - HTTP status code requirements
   - **Exact error messages** (do not modify)
   - Security constraints
   - Data integrity rules

   d. **Generate the code**:
   - Use correct package path based on project structure
   - Include all required imports
   - Implement exact method signatures as specified
   - Follow the exact validation messages from Safeguards

   **IMPORTANT**:
   - Do NOT deviate from the specifications in Operations
   - Do NOT add features or methods not specified
   - Do NOT change error messages from Safeguards
   - DO reference existing project patterns for consistency

6. **Batch validation after generation**

   After ALL code is generated, perform unified validation:

   a. **Compilation check**:
   - Run linter to check for syntax errors
   - Verify all imports are correct
   - Fix any type mismatches

   b. **Acceptance Criteria verification**:
   - Cross-check with the **Acceptance Criteria Traceability** table (if present)
   - Ensure each AC is addressed by the implementation
   - Verify error codes, HTTP status codes, and messages match exactly

   c. **Structure verification**:
   - Verify layered architecture is respected
   - Confirm dependency injection is correct
   - Check interface/implementation relationships

   d. **Fix any issues found**:
   - Fix compilation errors
   - Correct import statements
   - Ensure code follows project formatting standards

7. **Report generation summary**

   Provide a summary to the user:
   - List of created files with their responsibilities
   - Any deviations or assumptions made
   - Validation results (pass/fail for each check)
   - Story metadata mutations performed in Step 8 (assignee + refine-loop counter)

8. **Update the source story's YAML frontmatter** (Docs-as-Code automation)

   `/spdd-generate` is the natural pulse of the AI refine loop — every invocation against the same story counts as one iteration. Update the originating story file before reporting the summary.

   a. **Locate the source story file**:
    - Resolve from the prompt file's JIRA / story id (e.g., the prompt path `spdd/prompt/SPDD-003-...` ↔ story under `requirements/` whose frontmatter `id` shares the same trailing `STORY-{MODULE}-{SEQ}` or that the prompt was originally generated from).
    - If the link is ambiguous, list the candidates and ask the user to confirm via the **AskUserQuestion tool**.
    - If no story file can be matched (e.g., prompt was hand-written), skip steps 8b–8d and emit a warning so the team knows automation lacks an anchor.

   b. **Capture Dev identity** by running `git config user.name` (fall back to `$USER`, then `"unknown"`).

   c. **Mutate the frontmatter** in-place (Read → StrReplace, never rewrite the body or reorder keys):
    - `assignees.dev`: set ONLY if it is currently `null`. Preserve the original developer across re-runs (multi-dev pairing should be expressed in `inline_defects`/`related_bugs` or via a comment, not by stomping this field).
    - `quality_metrics.ai_refine_loops`: read the current integer and write back `current + 1`. Treat `null` or missing as `0` before incrementing.
    - **Do NOT** touch `status`, `timestamps.developed_at`, or `timestamps.tested_at` — those transitions belong to Git hooks / CI (PR open → `IN_TEST` + `developed_at`; PR merge → `DONE` + `delivered_at`).

   d. **Surface the mutations** in the run summary (Step 7), e.g.:

   ```
   Story metadata updated (requirements/[User-story-7]token-billing.md):
      - assignees.dev: null → wendi.zhang
      - quality_metrics.ai_refine_loops: 2 → 3
   ```

**Review & Iteration Loop**

If issues are discovered after generation (during testing or code review), follow the SPDD principle:

> **"When reality diverges, fix the prompt first — then update the code."**

1. **Identify the issue**: What behavior is incorrect or missing?

2. **Trace to prompt section**: Which part of the prompt caused this?
   - Wrong requirement interpretation → Update **Requirements**
   - Missing entity/relationship → Update **Entities**
   - Flawed strategy → Update **Approach**
   - Incorrect component design → Update **Structure**
   - Wrong implementation detail → Update **Operations**
   - Missing standard → Update **Norms**
   - Missing constraint → Update **Safeguards**

3. **Update the prompt first**: Modify the relevant section in the prompt file

4. **Regenerate affected code**: Only regenerate the components affected by the prompt change

5. **Commit together**: Commit the updated prompt and code together to maintain traceability

**Example iteration**:

```
Issue: "AgentService interface shouldn't contain business logic"

1. Trace: Operations section defines AgentService with logic
2. Update prompt: Move logic to AgentServiceImpl in Operations
3. Regenerate: Only regenerate AgentService and AgentServiceImpl
4. Commit: Commit prompt change + code change together
```

**Output**

- All generated source files following the project structure
- Summary of created files and their responsibilities
- Validation results
- Any issues requiring prompt modification

**Guardrails**

- Do NOT generate code without first reading the entire prompt file
- Do NOT re-plan the Operations sequence — execute in the defined order
- Do NOT skip any operation defined in the Operations section
- Do NOT change method signatures, field names, or error messages from the specification
- Do NOT add extra features, endpoints, or fields not specified
- Do NOT patch code directly when issues are found — update prompt first
- Always use the exact error messages from Safeguards
- Always follow Norms for coding style and patterns
- Always verify against Acceptance Criteria after generation
- Always check for and fix linter errors after batch generation
- Always commit prompt and code changes together
- Always update the source story's YAML frontmatter per Step 8 (set `assignees.dev` if null, increment `quality_metrics.ai_refine_loops`)
- NEVER overwrite an existing non-null `assignees.dev`
- NEVER mutate `status`, `timestamps.developed_at`, `timestamps.tested_at`, or `timestamps.delivered_at` from this command — those transitions are owned by Git hooks / CI
- NEVER reset or skip incrementing `quality_metrics.ai_refine_loops` even when generation fails partway — every full invocation is one loop

**Integration with /spdd-context**

This command is the second phase of the SPDD workflow:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           SPDD Workflow                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Phase 1: /spdd-context                                                │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Requirement → Alignment → Abstraction → Structured Prompt      │    │
│  │                                                                 │    │
│  │ Output: spdd/prompt/SPDD-XXX-*.md (REASONS Canvas)           │    │
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
│  Phase 3: Review & Iteration                                            │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Issue Found → Update Prompt → Regenerate → Commit Together     │    │
│  │                                                                 │    │
│  │ Principle: "Fix prompt first, then update code"                │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

The structured prompt serves as the **contract** between design and implementation, and must stay in sync with the code throughout the lifecycle.
