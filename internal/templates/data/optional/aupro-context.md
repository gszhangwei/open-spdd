---
name: /aupro-context
id: aupro-context
category: Development
description: Generate structured development context from a requirement using Aupro, save to file, then implement
---

Generate structured development context for a requirement using Aupro MCP tools, save the prompt to a file, then implement upon confirmation.

**Input**: The argument after `/aupro-context` is the requirement description. Can be a brief sentence or a detailed multi-line requirement.

**Steps**

1. **If no input provided, ask what they want to build**

   Use the **AskUserQuestion tool** (open-ended, no preset options) to ask:
   > "What do you want to implement? Describe the requirement in detail."

   **IMPORTANT**: Do NOT proceed without a clear requirement description.

2. **Call Aupro context_generation tool**

   Call the `user-Aupro-context_generation` MCP tool with:
   - `projectAbsolutePath`: the workspace root path
   - `user_message`: the user's requirement (translated to English if needed)

   **IMPORTANT**:
   - The requirement passed to `user_message` MUST be in English. If the user provides it in another language, translate it to English first.
   - **DO NOT summarize or condense the requirement.** Pass the FULL original content.
   - If a file is referenced (e.g., `@requirements/xxx.md`), read the complete file content and pass it verbatim to `user_message`.
   - Preserve ALL details including: Acceptance Criteria, error messages, Given/When/Then scenarios, JSON examples, etc.
   - **Exception**: Only summarize if the content exceeds 2000 characters. In this case, prioritize preserving Acceptance Criteria and exact error messages.

3. **Follow the Aupro tool's returned instructions**

   The tool returns a structured prompt containing:
   - Relevant project memories
   - Framework stages with required output formats
   - Construction guidelines and quality standards

   **You MUST strictly follow the returned instructions to construct the implementation prompt.** Do NOT invent your own structure — use the framework stages and formatting defined in the tool's response.

4. **Save the structured prompt to file**

   a. **Derive file name**: `{JIRA}-{TIMESTAMP}-[{ACTION}]-{scope}-{description}.md`
      
      - **JIRA**: `GGQPA-{number}` or `GGQPA-XXX`
      - **TIMESTAMP**: `YYYYMMDDHHmm`
      - **ACTION**: `[Feat]`, `[Fix]`, `[Refactor]`, `[Test]`, `[Docs]`
      - **scope**: `api`, `service`, `repo`, `bq`, `db`, `util` (optional)
      - **description**: kebab-case, < 10 words
      
      Examples: `GGQPA-169-202511271630-[Feat]-api-assignment-mismatch.md`

   b. **Create directory and write file**:
      - Ensure directory `aupro/prompt/` exists under the project root (create if not)
      - Write the complete structured implementation prompt to `aupro/prompt/<file-name>.md`

   c. **Show the saved file path** to the user:
      > "Structured prompt saved to `aupro/prompt/<file-name>.md`"

**Output**

The structured implementation prompt (saved to `aupro/prompt/<file-name>.md`), then implementation upon confirmation.

**Guardrails**
- Do NOT skip calling the Aupro context_generation tool — it provides project-specific memories and framework
- Do NOT implement code before user confirms the structured prompt
- Do NOT modify the framework structure returned by Aupro — follow it as-is
- The `user_message` parameter MUST be in English
- Always use the actual workspace root for `projectAbsolutePath`
- File name MUST follow SPDD naming convention defined above
- Use `GGQPA-XXX` if JIRA ticket number is unknown
- Always create `aupro/prompt/` directory if it does not exist
- **NEVER summarize, condense, or paraphrase the requirement content** — pass the complete original text to the Aupro tool (only exception: content exceeds 2000 characters)
