---
name: /spdd-block
id: spdd-block
category: Development
description: Mark an SPDD story as BLOCKED, append a block_records entry with reason and blocked_at timestamp, and remember the prior status for /spdd-unblock to restore
---

Record an external blocker against a story without losing its current workflow state. The previous `status` is preserved inside the block record so `/spdd-unblock` can restore it cleanly.

**Input**: `/spdd-block @requirements/<story-file>.md -m "<reason>"`

Both arguments are required. Examples:

```
/spdd-block @requirements/[User-story-7]token-billing.md -m "Waiting on payments API quota approval"
/spdd-block @requirements/[User-story-3]invoice-export.md -m "Blocked by upstream story User-story-2"
```

**Steps**

1. **Validate input**

   a. **If no `@` story file is provided**, use the **AskUserQuestion tool** (open-ended) to ask:
   > "Which story should be marked as BLOCKED? Provide a path under `requirements/`."

   b. **If no `-m "<reason>"` is provided**, use the **AskUserQuestion tool** to ask:
   > "What is the blocker reason? (Be concrete — this is logged in the audit trail.)"

   **IMPORTANT**: Do NOT proceed without both inputs. Reasons MUST NOT be empty strings or placeholders.

2. **Read the story file and parse its frontmatter**

   - Read the entire file with the Read tool.
   - Confirm a YAML frontmatter block exists at the top with the canonical schema defined by `/spdd-story` (must contain `status:` and `block_records:`).
   - If the file has no frontmatter, abort and instruct the user to backfill it (this command does NOT create frontmatter — that is `/spdd-story`'s responsibility).

3. **Capture the previous status and the reporter identity**

   - Snapshot the current `status` value into a local variable `prev_status`.
   - If `prev_status` is already `BLOCKED`, do NOT push a new block record. Instead, surface the most recent open block_record (one whose `unblocked_at` is null) and exit with a message telling the user to call `/spdd-unblock` first.
   - Run `git config user.name` to capture `reporter` (fall back to `$USER`, then `"unknown"`).
   - Capture the current local timestamp `now = YYYY-MM-DD HH:MM`.

4. **Mutate the frontmatter** in-place via Read → StrReplace (never reorder keys, never touch the body):

   a. Set `status: BLOCKED`.

   b. Append a new entry to `block_records` (create the array if currently `[]`):

   ```yaml
   block_records:
     - reason: "{reason from -m}"
       blocked_at: "{now}"
       unblocked_at: null
       previous_status: "{prev_status}"
       reporter: "{reporter}"
   ```

   - YAML strings containing `:`, `#`, or leading `-` MUST be double-quoted.
   - Preserve the indentation style already used in the file (2-space indent inside arrays).
   - Do NOT modify any earlier `block_records` entries — they are immutable history.

5. **Report the change** to the user:

   ```
   Story marked as BLOCKED.
      File: requirements/<story-file>.md
      Status: <prev_status> → BLOCKED
      Reason: <reason>
      Logged at: <now> by <reporter>

   When the blocker is resolved, run:
      /spdd-unblock @requirements/<story-file>.md
   ```

**Output**

The story file's YAML frontmatter is mutated to set `status: BLOCKED` and to append an open `block_records` entry. No body changes, no other field changes.

**Guardrails**

- Do NOT proceed without both a story file and a non-empty reason
- Do NOT create frontmatter — abort and refer the user back to `/spdd-story` if the block is missing
- Do NOT touch any field other than `status` and `block_records`
- Do NOT modify previously closed `block_records` entries — they are immutable audit history
- Do NOT push a second open block record on top of an existing open one — fail fast and require `/spdd-unblock` first
- Do NOT prompt for the reporter — derive from `git config user.name`
- Always preserve the existing key order and indentation style of the frontmatter
