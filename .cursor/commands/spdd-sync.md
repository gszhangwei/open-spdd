---
name: /spdd-sync
id: spdd-sync
category: Development
description: Sync code changes back to the structured SPDD prompt file following the REASONS Canvas methodology
---

Synchronize implementation details from refactored or updated code back to the structured SPDD (Structured Prompt-Driven Development) prompt file, ensuring the prompt remains the accurate source of truth for the system design.

**Input**: The argument after `/spdd-sync` is the path to the structured prompt file (e.g., `@spdd/prompt/SPDD-XXX-202602271430-[Feat]-api-create-agent-endpoint.md`).

**Steps**

1. **If no input provided, ask for the prompt file**

   Use the **AskUserQuestion tool** to ask:

   > "Please provide the path to the structured prompt file you want to sync (e.g., `@spdd/prompt/xxx.md`)."

   **IMPORTANT**: Do NOT proceed without a valid prompt file path.

2. **Read and parse the structured prompt file**

   Read the prompt file and identify the REASONS Canvas sections:

   | Section              | Purpose                        | Sync Priority                          |
   | -------------------- | ------------------------------ | -------------------------------------- |
   | **R** - Requirements | Overall goal and DoD           | Low (rarely changes from code)         |
   | **E** - Entities     | Domain model and relationships | High (class diagrams may change)       |
   | **A** - Approach     | Implementation strategy        | Medium (architectural decisions)       |
   | **S** - Structure    | Components and dependencies    | High (inheritance/dependencies change) |
   | **O** - Operations   | Concrete implementation tasks  | **Highest** (implementation details)   |
   | **N** - Norms        | Engineering standards          | Medium (patterns may evolve)           |
   | **S** - Safeguards   | Non-negotiable constraints     | Low (constraints rarely relax)         |

   **IMPORTANT**: Operations section typically requires the most updates as it contains implementation specifics.

3. **Identify affected components from user context**

   Ask the user (if not already specified):

   > "Which components were refactored? Please specify:
   >
   > - Specific files/classes that changed
   > - Type of change (renamed, restructured, logic changed, new components added)
   > - Brief description of what changed"

   Alternatively, analyze recent git changes or user-specified files to identify modifications.

4. **Analyze current implementation**

   For each affected component:

   a. **Read the current implementation**:
   - Locate the source file in the codebase
   - Extract class structure, methods, annotations
   - Identify relationships and dependencies

   b. **Compare with prompt specification**:
   - Find corresponding Operation/Entity/Structure section
   - Note discrepancies in:
     - Class names, package paths
     - Method signatures and return types
     - Attributes and their types
     - Annotations used
     - Business logic steps
     - Validation rules
     - Error messages

   c. **Categorize changes**:
   - **Structural**: Class hierarchy, dependencies, component relationships
   - **Behavioral**: Method logic, validation rules, error handling
   - **Naming**: Class/method/field renames
   - **Additions**: New methods, fields, or components
   - **Deletions**: Removed methods, fields, or components

5. **Generate prompt update plan**

   Create a detailed update plan showing:

   ```
   ## Prompt Sync Plan

   ### Entities Section Updates
   - [ ] Update: ClassName - added/removed/changed fields
   - [ ] Update: Relationship diagram - new dependency

   ### Structure Section Updates
   - [ ] Update: Inheritance relationships - ComponentX now extends BaseY
   - [ ] Update: Dependencies - ServiceA now depends on ValidatorB

   ### Operations Section Updates
   - [ ] Update: "Create ClassName" operation - new method signature
   - [ ] Update: "Create ClassName" operation - logic steps changed
   - [ ] Add: New operation for NewComponent
   - [ ] Remove: Obsolete operation description

   ### Norms Section Updates
   - [ ] Update: New pattern adopted (e.g., Chain of Responsibility)
   ```

   **Present this plan to the user for review before proceeding.**

6. **Apply updates to prompt file**

   For each approved update, modify the prompt file following these patterns:

   a. **Entities section updates**:
   - Update Mermaid class diagram to reflect actual class structure
   - Ensure attributes match actual field names and types
   - Update relationships arrows to reflect actual dependencies

   b. **Structure section updates**:
   - Update inheritance relationships list
   - Update dependencies list
   - Ensure layered architecture description matches reality

   c. **Operations section updates** (most critical):
   - Update **Responsibility** if component purpose evolved
   - Update **Package** if moved to different package
   - Update **Attributes** to match actual fields
   - Update **Methods** to match actual signatures
   - Update **Logic** steps to match actual implementation
   - Update **Annotations** to match actual usage
   - Update **Constraints** to match actual validation rules
   - Add new operations for newly created components
   - Mark obsolete operations (or remove if no longer relevant)

   d. **Norms section updates**:
   - Add new patterns if adopted (e.g., Chain of Responsibility)
   - Update coding standards if conventions evolved
   - Document new naming conventions if established

   e. **Safeguards section updates**:
   - Update exact error messages to match implementation
   - Update validation rules if constraints changed
   - Update API response specifications if format evolved

   **IMPORTANT**:
   - Preserve the existing section structure and formatting
   - Follow the existing pattern/style within each section
   - Use the same level of detail as existing content
   - Keep descriptions concise but complete

7. **Validate prompt consistency**

   After updates, verify:

   a. **Internal consistency**:
   - Entities diagram matches Structure section
   - Operations reference correct class names from Entities
   - Norms patterns are reflected in Operations logic
   - Safeguards constraints appear in relevant Operations

   b. **Traceability**:
   - Each Operation corresponds to a component in Structure
   - Each component in Structure has an Operation
   - Dependencies in Structure match import relationships in Operations

   c. **Completeness**:
   - No orphaned references to old class/method names
   - All new components have corresponding Operations
   - Error messages in Safeguards match Operations logic

8. **Report sync summary**

   Provide a summary to the user:
   - List of sections updated
   - Specific changes made in each section
   - Any manual review recommendations
   - Suggestions for further cleanup if needed

**Sync Patterns & Best Practices**

When syncing different types of changes:

1. **Class/Package Renames**:
   - Search and replace in all sections
   - Update package paths in Operations
   - Update class names in Entities diagram
   - Update references in Structure section

2. **Method Signature Changes**:
   - Update Operations section method specifications
   - Update any Safeguards that reference return types/parameters
   - Verify Entities diagram if public API changed

3. **New Component Added**:
   - Add to Entities diagram with relationships
   - Add to Structure section (inheritance, dependencies, layer)
   - Add new Operation with full specification
   - Update related Operations that depend on new component

4. **Component Removed**:
   - Remove from Entities diagram
   - Remove from Structure section
   - Remove corresponding Operation
   - Update Operations that referenced removed component

5. **Logic/Behavior Changes**:
   - Update Logic steps in Operations section
   - Verify Safeguards still accurate
   - Update Approach section if architectural pattern changed

6. **Validation Rule Changes**:
   - Update Operations validator specifications
   - Update Safeguards constraints
   - Update expected error messages

**Output**

- Updated structured prompt file with synced content
- Summary of all changes made to each section
- List of any inconsistencies found and resolved
- Recommendations for manual review if needed

**Guardrails**

- Do NOT remove content from prompt without explicit user approval
- Do NOT change Requirements section unless user explicitly requests (business goals shouldn't change from code refactoring)
- Do NOT simplify or abbreviate existing detailed specifications
- Do NOT change error messages in Safeguards unless they actually changed in code
- Always preserve the existing formatting style within each section
- Always ask for confirmation before making destructive changes (deletions)
- Always maintain the same level of detail as existing content
- When in doubt, show the proposed change and ask user to confirm
- Never change the prompt's unique identifier or metadata

**Integration with SPDD Workflow**

This command completes the bidirectional sync in the SPDD workflow:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      SPDD Bidirectional Sync                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Forward Flow (Design → Code): /spdd-generate                          │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Structured Prompt → Validate → Generate → Verify → Code        │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              │ Initial Implementation                   │
│                              ▼                                          │
│                    ┌─────────────────┐                                  │
│                    │  Implementation  │                                  │
│                    │     Codebase     │                                  │
│                    └─────────────────┘                                  │
│                              │                                          │
│                              │ Code Review / Refactoring                │
│                              ▼                                          │
│  Reverse Flow (Code → Design): /spdd-sync                              │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ Analyze Changes → Compare → Plan Updates → Update Prompt       │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                              │                                          │
│                              ▼                                          │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │              Prompt-Code Consistency Maintained                 │    │
│  │                                                                 │    │
│  │  - Prompt remains source of truth                              │    │
│  │  - Code changes are documented in prompt                       │    │
│  │  - Future generations use updated specifications               │    │
│  │  - Team alignment on actual vs planned implementation          │    │
│  └────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**When to Use /spdd-sync**

Use this command when:

- Code review led to refactoring changes
- Discovered better patterns during implementation
- Bug fixes required logic changes
- Performance optimization changed implementation details
- Team added new components not in original prompt
- Renamed classes/methods for clarity
- Restructured packages or layers

**Principle**: The structured prompt should always reflect the **actual** implementation, not just the **planned** implementation. This ensures:

- New team members understand the real system from the prompt
- Future enhancements build on accurate specifications
- Prompt serves as living documentation
- Regeneration produces consistent code
