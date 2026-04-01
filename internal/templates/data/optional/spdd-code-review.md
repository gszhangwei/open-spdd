---
name: /spdd-code-review
id: spdd-code-review
category: Development
description: Review AI-generated code against REASONS-Canvas structured prompts, detecting intent drift, safeguard violations, and scope boundary issues to reduce human reviewer cognitive load
---

Review implementation code against its corresponding REASONS-Canvas structured prompt, producing a structured review report that highlights alignment gaps, intent drift, safeguard violations, and areas requiring human attention — reducing the cognitive load of reviewing AI-generated code.

**Input**: The argument after `/spdd-code-review` includes the structured prompt file reference and the code scope to review.

Input can be provided in several ways:

1. **Prompt file + code files**: Reference the REASONS-Canvas prompt and specific implementation files
2. **Prompt file + folder**: Reference the prompt and an entire directory of generated code
3. **Prompt file + git diff**: Reference the prompt and use git diff to scope code changes
4. **Prompt file only**: Reference the prompt and let the review infer code scope from Operations section

**Examples**:

```
# Review specific files against prompt
/spdd-code-review @spdd/prompt/GGQPA-169-202603061530-[Feat]-api-user-registration.md @src/main/java/com/example/user/

# Review with explicit file list
/spdd-code-review @spdd/prompt/GGQPA-XXX-202603131530-[Feat]-token-usage-billing.md @src/controllers/BillingController.java @src/services/BillingService.java

# Review recent changes against prompt
/spdd-code-review @spdd/prompt/GGQPA-169-202603061530-[Feat]-api-user-registration.md --git-diff HEAD~3

# Review with prompt only (auto-infer code scope from Operations)
/spdd-code-review @spdd/prompt/GGQPA-169-202603061530-[Feat]-api-user-registration.md
```

**Steps**

1. **Validate and consolidate input**

   a. **If no prompt file provided**, use the **AskUserQuestion tool** (open-ended, no preset options) to ask:

   > "Please provide the path to the REASONS-Canvas prompt file to review against (e.g., `@spdd/prompt/xxx.md`)."

   **IMPORTANT**: Do NOT proceed without a valid REASONS-Canvas prompt file.

   b. **Determine code scope**:
   - **If code files/folders provided**: Read ALL referenced files completely
   - **If `--git-diff` specified**: Extract changed files from the git diff range
   - **If no code scope provided**: Infer code scope from the Operations section of the prompt — locate the files described in each operation by searching the codebase for matching class names, package paths, and file locations

   c. **Context Integrity Check**:
   - Verify the prompt file is a valid REASONS-Canvas document (contains R, E, A, S, O, N, S sections)
   - Verify all code references were successfully read
   - If any file cannot be read, report the error and ask user for alternatives

2. **Parse the REASONS-Canvas prompt**

   Extract and internalize each section for review comparison:

   | Section              | Review Focus                                                  |
   | -------------------- | ------------------------------------------------------------- |
   | **R** - Requirements | Does code solve the stated problem — no more, no less?        |
   | **E** - Entities     | Do code entities match defined model? Any entity drift?       |
   | **A** - Approach     | Does implementation follow chosen strategies?                 |
   | **S** - Structure    | Do inheritance, dependencies, and layers match?               |
   | **O** - Operations   | Does each operation's implementation match its specification? |
   | **N** - Norms        | Are coding standards consistently applied?                    |
   | **S** - Safeguards   | Are all constraints respected? Any violations?                |

3. **Read and analyze the implementation code**

   For each file in the code scope:

   a. **Extract code structure**:
   - Class hierarchy, interfaces, inheritance
   - Method signatures, return types, parameters
   - Field definitions, types, annotations
   - Dependency injection patterns
   - Exception handling patterns
   - Business logic flow

   b. **Build code inventory**:
   - List all components (classes, interfaces, enums) created
   - List all methods and their signatures
   - List all dependencies and injections
   - List all annotations used
   - List all error messages and exception types

4. **Section-by-section alignment analysis**

   Compare each REASONS section against the implementation code:

   ***

   ### R - Requirements Alignment

   **Check**: Does the code solve the stated problem without expanding or narrowing scope?
   - Verify the implementation addresses the core requirement described in Requirements
   - Check if the code introduces functionality beyond what Requirements specifies (scope expansion)
   - Check if the code omits functionality that Requirements implies (scope contraction)
   - Flag any feature that cannot be traced back to a stated requirement

   **Output Format**:

   ```
   ### Requirements Alignment

   **Status**: [✅ Aligned | ⚠️ Partial Drift | ❌ Significant Drift]

   **Alignment**: [What matches]
   **Scope Expansion**: [What code does that the requirement doesn't ask for — if any]
   **Scope Contraction**: [What the requirement asks for but code doesn't implement — if any]
   ```

   ***

   ### E - Entities Alignment

   **Check**: Do code entities match the defined entity model?
   - Compare class definitions against the Mermaid entity diagram
   - Verify attribute names, types, and relationships match
   - Check if new entities were introduced that aren't in the diagram
   - Check if existing entity structures were modified beyond what the prompt specifies
   - **Conservative constraint check**: Verify the code respects "Prohibit Unnecessary Refactoring" — existing simple data types are not wrapped in complex entity classes without justification

   **Output Format**:

   ```
   ### Entities Alignment

   **Status**: [✅ Aligned | ⚠️ Partial Drift | ❌ Significant Drift]

   **Matched Entities**: [List entities that match the spec]
   **Entity Drift**:
   - [Entity]: [What differs — added attributes, changed types, new relationships]
   **Unauthorized Entities**: [Entities in code but not in prompt — if any]
   **Conservative Constraint Violations**: [Cases where existing structures were unnecessarily refactored — if any]
   ```

   ***

   ### A - Approach Alignment

   **Check**: Does the implementation follow the chosen solution strategies?
   - Verify architectural patterns match (e.g., if prompt says "synchronous", code shouldn't use event-driven)
   - Verify framework and technology choices match
   - Verify error handling strategy matches
   - Check for unauthorized architectural decisions (e.g., adding a caching layer not specified in Approach)

   **Output Format**:

   ```
   ### Approach Alignment

   **Status**: [✅ Aligned | ⚠️ Partial Drift | ❌ Significant Drift]

   **Followed Strategies**: [What matches]
   **Approach Drift**:
   - [Strategy area]: [Expected approach] → [Actual approach]
   **Unauthorized Decisions**: [Architectural choices made by code but not in prompt — if any]
   ```

   ***

   ### S - Structure Alignment

   **Check**: Do inheritance, dependencies, and layered architecture match the specification?
   - Verify inheritance relationships (interfaces, abstract classes, implementations)
   - Verify dependency injection matches
   - Verify layered architecture is respected (Controller → Service → Repository → DAO)
   - Check for circular dependencies or layer violations not permitted by Structure
   - Check for components introduced without corresponding Structure definition

   **Output Format**:

   ```
   ### Structure Alignment

   **Status**: [✅ Aligned | ⚠️ Partial Drift | ❌ Significant Drift]

   **Matched Structure**: [Inheritance and dependencies that match]
   **Structure Drift**:
   - [Component]: [Expected relationship] → [Actual relationship]
   **Layer Violations**: [Components that violate the defined layered architecture — if any]
   ```

   ***

   ### O - Operations Alignment

   **Check**: Does each operation's implementation match its specification?

   This is the most detailed check — compare each operation in the prompt against its code implementation:
   - Verify method signatures match (name, parameters, return type)
   - Verify business logic steps match the described logic sequence
   - Verify validation rules are implemented as specified
   - Verify error handling matches the described approach
   - Check for logic steps in code that have no corresponding prompt description
   - Check for prompt-described logic steps that are missing in code

   **Output Format**:

   ```
   ### Operations Alignment

   **Status**: [✅ Aligned | ⚠️ Partial Drift | ❌ Significant Drift]

   **Operation**: [Operation Name from prompt]
   - **Signature**: [✅ Match | ❌ Mismatch: expected X, got Y]
   - **Logic Steps**: [✅ N/N matched | ⚠️ N/M matched — details below]
   - **Validation**: [✅ Match | ❌ Mismatch — details]
   - **Error Handling**: [✅ Match | ❌ Mismatch — details]
   - **Extra Logic**: [Logic in code not described in prompt — if any]
   - **Missing Logic**: [Logic in prompt not found in code — if any]

   [Repeat for each operation]
   ```

   ***

   ### N - Norms Alignment

   **Check**: Are coding standards consistently applied across all generated code?
   - Verify annotation standards (correct annotations on correct component types)
   - Verify dependency injection style matches
   - Verify exception handling follows the defined pattern
   - Verify logging conventions are followed
   - Verify naming conventions are consistent
   - Check for norm violations that may indicate AI deviation from project standards

   **Output Format**:

   ```
   ### Norms Alignment

   **Status**: [✅ Aligned | ⚠️ Partial Drift | ❌ Significant Drift]

   **Followed Norms**: [Standards that are consistently applied]
   **Norm Violations**:
   - [Norm]: [Expected pattern] → [Actual pattern] — [File:Line]
   ```

   ***

   ### S - Safeguards Alignment

   **Check**: Are all constraints respected? Are there any violations of "what NOT to do"?

   This is a critical check — Safeguards define the negative space (boundaries that must not be crossed):
   - Check each Safeguard constraint against the implementation
   - Verify exact error messages match (if specified)
   - Verify performance constraints are respected
   - Verify security constraints are followed
   - Verify API constraints (status codes, response format) match
   - Verify data constraints (validation rules, format requirements) are enforced
   - **Explicit violation flagging**: Any Safeguard violation is a high-priority finding

   **Output Format**:

   ```
   ### Safeguards Alignment

   **Status**: [✅ All Respected | ⚠️ Partial Violations | ❌ Critical Violations]

   **Respected Safeguards**: [Constraints that are properly enforced]
   **Violations**:
   - 🔴 [Safeguard]: [What should NOT happen] → [What code actually does] — [File:Line]
   ```

5. **Intent drift detection**

   Beyond section-by-section alignment, perform a holistic intent drift analysis across three categories:

   a. **Positive Drift (Additions)** — Code does things the prompt doesn't specify:
   - Added methods, endpoints, or fields not in Operations
   - Added validation rules not in Safeguards
   - Added error handling not in Norms
   - Added dependencies not in Structure
   - Added design patterns not in Approach (e.g., AI added caching, event publishing, retry logic)

   b. **Negative Drift (Omissions)** — Prompt specifies things code doesn't implement:
   - Missing operations or methods
   - Missing validation rules
   - Missing error handling
   - Missing constraints enforcement

   c. **Direction Drift (Divergence)** — Code does what the prompt says, but differently:
   - Different architectural pattern than specified
   - Different error handling strategy
   - Different data flow direction
   - Different entity relationships
   - Subtle algorithmic differences in business logic

   **Output Format**:

   ```
   ## Intent Drift Analysis

   ### Positive Drift (Unauthorized Additions)
   | Finding | Severity | Location | Description |
   |---------|----------|----------|-------------|
   | [ID]    | [🔴/🟡] | [File:Line] | [What was added without prompt authorization] |

   ### Negative Drift (Missing Implementations)
   | Finding | Severity | Location | Description |
   |---------|----------|----------|-------------|
   | [ID]    | [🔴/🟡] | [Prompt Section] | [What prompt specifies but code doesn't implement] |

   ### Direction Drift (Divergent Approaches)
   | Finding | Severity | Location | Description |
   |---------|----------|----------|-------------|
   | [ID]    | [🔴/🟡] | [File:Line] | [Expected: X, Actual: Y] |
   ```

6. **Implicit decision detection**

   Identify places where AI made design decisions that were not explicitly delegated in the prompt:
   - Data structure choices not specified in Entities
   - Algorithm choices not specified in Operations
   - Error message wording not specified in Safeguards
   - Configuration defaults not specified anywhere
   - Concurrency or threading decisions not in Approach

   These are not necessarily wrong — but they represent areas where the AI exercised judgment that the human reviewer should be aware of and may want to validate.

   **Output Format**:

   ```
   ## Implicit Decisions (AI Judgment Points)

   The following decisions were made by the AI without explicit guidance in the prompt.
   These require human validation:

   | Decision | Category | Location | AI's Choice | Risk |
   |----------|----------|----------|-------------|------|
   | [What was decided] | [Data/Algorithm/Config/...] | [File:Line] | [What AI chose] | [Low/Medium/High] |
   ```

7. **Scope boundary check**

   Verify the code only modifies files and components within the defined scope:
   - Check if any existing files outside the prompt's scope were modified
   - Check if the code introduces dependencies on components not mentioned in Structure
   - Check if the code modifies shared utilities, configurations, or base classes beyond the prompt's scope
   - Flag any "ripple effect" changes that may have unintended consequences

   **Output Format**:

   ```
   ## Scope Boundary Check

   **Status**: [✅ Within Scope | ⚠️ Minor Boundary Crossing | ❌ Significant Out-of-Scope Changes]

   **In-Scope Components**: [count] components within defined scope
   **Boundary Crossings**:
   - [File/Component]: [Why it's outside scope] — [Risk assessment]
   ```

8. **Assemble the review report**

   Construct a comprehensive review report with cognitive load reduction as the primary design goal:

   ```markdown
   # SPDD Code Review: [Derived Title]

   ## Review Context

   - **Prompt**: [prompt file path]
   - **Code Scope**: [files/folders reviewed]
   - **Review Date**: [timestamp]

   ## Review Summary (Start Here)

   | Dimension      | Status     | Findings       | Priority            |
   | -------------- | ---------- | -------------- | ------------------- |
   | Requirements   | [✅/⚠️/❌] | [count] issues | [High/Medium/Low/-] |
   | Entities       | [✅/⚠️/❌] | [count] issues | [High/Medium/Low/-] |
   | Approach       | [✅/⚠️/❌] | [count] issues | [High/Medium/Low/-] |
   | Structure      | [✅/⚠️/❌] | [count] issues | [High/Medium/Low/-] |
   | Operations     | [✅/⚠️/❌] | [count] issues | [High/Medium/Low/-] |
   | Norms          | [✅/⚠️/❌] | [count] issues | [High/Medium/Low/-] |
   | Safeguards     | [✅/⚠️/❌] | [count] issues | [High/Medium/Low/-] |
   | Intent Drift   | [✅/⚠️/❌] | [count] issues | [High/Medium/Low/-] |
   | Scope Boundary | [✅/⚠️/❌] | [count] issues | [High/Medium/Low/-] |

   **Overall Assessment**: [✅ Ready to Merge | ⚠️ Needs Attention | ❌ Needs Rework]

   ## 🔴 Must Review (Critical)

   [Top-priority findings that must be addressed before merge]

   - [Finding]: [Brief description] — [Location]

   ## 🟡 Should Review (Important)

   [Findings that should be reviewed but may be acceptable]

   - [Finding]: [Brief description] — [Location]

   ## 🟢 Informational (Low Risk)

   [Findings for awareness only — no action needed unless concerned]

   - [Finding]: [Brief description] — [Location]

   ## Detailed Analysis

   [All section-by-section alignment results from Step 4]

   ## Intent Drift Analysis

   [Results from Step 5]

   ## Implicit Decisions

   [Results from Step 6]

   ## Scope Boundary Check

   [Results from Step 7]

   ## Recommended Actions

   1. [Action]: [Which prompt section to update if needed]
   2. [Action]: [Which code to fix if needed]
   ```

   **Cognitive Load Reduction Design Principles**:
   - **Summary table first**: The reviewer can assess overall health in 5 seconds
   - **Traffic light system**: ✅/⚠️/❌ provides instant visual triage
   - **Priority-ranked findings**: 🔴 → 🟡 → 🟢 tells the reviewer exactly where to focus
   - **Overall assessment**: A single verdict reduces decision fatigue
   - **Detailed analysis is opt-in**: Only read sections where issues are flagged
   - **Actionable recommendations**: Each finding includes what to do (fix code or update prompt)

9. **Save the review report**

   a. **Derive file name**: `{JIRA}-{TIMESTAMP}-[Review]-{description}.md`
   - **JIRA**: Extract from prompt file name or use `GGQPA-XXX`
   - **TIMESTAMP**: `YYYYMMDDHHmm` (current time)
   - **description**: Derive from prompt context — kebab-case, < 10 words

   Examples:
   - `GGQPA-169-202604011530-[Review]-api-user-registration.md`
   - `GGQPA-XXX-202604011530-[Review]-token-usage-billing.md`

   b. **Create directory and write file**:
   - Ensure directory `spdd/review/` exists under the project root (create if not)
   - Write the complete review report to `spdd/review/<file-name>.md`

10. **Show review summary to user**

    ```
    ✅ Code review complete. Report saved to `spdd/review/<file-name>.md`

    📋 Review Summary:
    - REASONS Alignment: [N/7] sections fully aligned
    - Intent Drift: [count] findings ([additions], [omissions], [divergences])
    - Safeguard Violations: [count]
    - Implicit Decisions: [count] requiring human validation
    - Scope Boundary: [✅ Within Scope | ⚠️ N boundary crossings]

    🔴 Critical Issues: [count]
    🟡 Important Issues: [count]
    🟢 Informational: [count]

    Overall: [✅ Ready to Merge | ⚠️ Needs Attention | ❌ Needs Rework]
    ```

11. **Offer follow-up actions**

    Based on review findings, suggest appropriate next steps:

    > Based on the review findings:
    >
    > - [If prompt needs update]: "Would you like me to update the prompt using `/spdd-prompt-update`?"
    > - [If code needs fixing]: "Would you like me to fix the identified code issues?"
    > - [If all aligned]: "The code is well-aligned with the prompt. Ready to merge."
    > - [If sync needed]: "Would you like me to sync the accepted changes back to the prompt using `/spdd-sync`?"

**Output**

A structured code review report saved to `spdd/review/<file-name>.md` containing:

- REASONS section-by-section alignment analysis
- Intent drift detection (additions, omissions, divergences)
- Safeguard violation check
- Implicit decision surfacing
- Scope boundary verification
- Priority-ranked findings with traffic light system (🔴/🟡/🟢)
- Overall merge readiness assessment
- Actionable recommendations

**Guardrails**

- Do NOT proceed without a valid REASONS-Canvas prompt file
- Do NOT skip any of the 7 REASONS sections — analyze ALL of them
- Do NOT modify any code or prompt files during review — this command is read-only analysis
- Do NOT generate false positives by flagging stylistic differences that don't affect functionality or violate Norms
- Do NOT ignore Safeguard violations — every defined constraint must be checked
- Do NOT assume code is correct just because it compiles and runs
- Do NOT produce review findings without specifying the location (file and approximate code area)
- Do NOT provide only positive findings — explicitly confirm what is aligned AND what is not
- Always read the ENTIRE prompt file — partial reading leads to incomplete review
- Always read ALL referenced code files completely — do not skim or skip
- Always provide actionable recommendations for each finding
- Always include an overall assessment verdict
- Always save the review report to `spdd/review/` directory
- File name MUST follow the naming convention defined above
- Use `GGQPA-XXX` if JIRA ticket number cannot be extracted from prompt file name

**Context Integrity Guardrails**:

- **MUST read ALL `@` referenced files completely** — do NOT skip or partially read any referenced file
- **MUST read folder contents** when `@` references a folder — scan and read all relevant code files
- **Verify all references resolved** — if any `@` reference fails to read, report error immediately
- **MUST read the complete REASONS-Canvas prompt** — do NOT skip any section
- **Preserve objectivity** — report what the code does vs what the prompt says, not what you think is better

**Review Severity Definitions**:

| Severity      | Symbol | Criteria                                                                                     | Examples                                                                                           |
| ------------- | ------ | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| Critical      | 🔴     | Safeguard violation, core logic mismatch, security risk, significant intent drift            | Prompt says "no caching" but code adds Redis; business logic calculates differently than specified |
| Important     | 🟡     | Partial drift, unauthorized additions, minor deviations, implicit decisions with medium risk | Extra validation not in spec; different error message; additional method parameters                |
| Informational | 🟢     | Minor stylistic differences, benign additions, low-risk implicit decisions                   | Variable naming style; comment additions; import ordering                                          |

**Integration with SPDD Workflow**

This command fills the **quality assurance gap** between code generation and prompt sync:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           SPDD Workflow                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Phase 0: /spdd-analysis                                                │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Business Requirement → Enriched Context                         │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  Phase 1: /spdd-reasons-canvas                                         │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Enriched Context → REASONS Canvas Structured Prompt             │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  Phase 2: /spdd-generate                                               │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Structured Prompt → Implementation Code                         │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  Phase 3: /spdd-code-review  ← YOU ARE HERE                           │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Structured Prompt + Code → Alignment Analysis                    │    │
│  │                                                                 │    │
│  │ Core Review Dimensions:                                          │    │
│  │ 1. REASONS Section-by-Section Alignment (7 dimensions)          │    │
│  │ 2. Intent Drift Detection (additions, omissions, divergences)   │    │
│  │ 3. Safeguard Violation Check                                     │    │
│  │ 4. Implicit Decision Surfacing                                   │    │
│  │ 5. Scope Boundary Verification                                   │    │
│  │                                                                 │    │
│  │ Output: spdd/review/GGQPA-XXX-*-[Review]-*.md                  │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                       ┌──────┴──────┐                                   │
│                       ▼             ▼                                    │
│              [Issues Found]   [All Aligned]                             │
│                       │             │                                    │
│                       ▼             ▼                                    │
│  Fix: /spdd-prompt-update    Phase 4: /spdd-sync                       │
│       or code fix            ┌────────────────────────────────────┐    │
│       │                      │ Code Changes → Update Prompt        │    │
│       │                      └────────────────────────────────────┘    │
│       │                                                                  │
│       └──► Re-run /spdd-code-review                                     │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**Why This Phase Matters**

AI-generated code often passes superficial review — it compiles, follows naming conventions, and looks reasonable. But design-level drift is subtle and expensive to catch later:

1. **Prompt is the contract**: Code review against the prompt catches misalignment that conventional code review misses — not "is this code well-written?" but "does this code do what we decided it should do?"
2. **Negative space enforcement**: Safeguards define what NOT to do, and AI frequently crosses these boundaries in reasonable-looking ways. Systematic checking catches violations that human scanning often misses.
3. **Cognitive load reduction**: Instead of a human cross-referencing a 300-line prompt against 2000 lines of code, the review report tells them exactly where to look and what to verify.
4. **Intent drift is the #1 risk**: The most dangerous AI coding failures aren't bugs — they're subtle design direction changes that look correct but diverge from the team's actual intent.
5. **Implicit decision visibility**: AI makes dozens of micro-decisions during code generation. Surfacing these gives the reviewer awareness of where AI exercised judgment without explicit guidance.

**Relationship to Design Philosophy**

This command operationalizes three key insights from the SPDD design philosophy:

| Insight                                             | How This Command Addresses It                                           |
| --------------------------------------------------- | ----------------------------------------------------------------------- |
| "AI isn't not smart enough — it has too many ideas" | Positive Drift detection catches unauthorized additions                 |
| "Negative space is easily overlooked"               | Safeguard Violation Check systematically verifies every constraint      |
| "Design intent is easily lost over time"            | Review report becomes a traceable record of prompt-code alignment state |

**Common Review Scenarios**

1. **Post-Generation Review**
   - Trigger: Right after `/spdd-generate`
   - Focus: Full REASONS alignment, especially Operations and Safeguards
   - Goal: Catch drift before code enters the codebase

2. **Pull Request Review Assistance**
   - Trigger: During code review of AI-generated PR
   - Focus: Intent drift, scope boundary, implicit decisions
   - Goal: Help human reviewer focus on what matters

3. **Iterative Development Review**
   - Trigger: After manual code modifications
   - Focus: Detect where manual changes diverge from prompt
   - Goal: Decide whether to update prompt (`/spdd-prompt-update`) or revert code

4. **Pre-Sync Validation**
   - Trigger: Before running `/spdd-sync`
   - Focus: Understand the full scope of drift before syncing
   - Goal: Ensure intentional changes are synced while accidental drift is fixed
