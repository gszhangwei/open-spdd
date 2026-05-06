---
name: /spdd-unblock
id: spdd-unblock
category: Development
description: Close the most recent open block_records entry on an SPDD story, restore the previously recorded status, and stamp unblocked_at
---

Close an open blocker on a story and restore the workflow state that was in effect before `/spdd-block` was called. This is the symmetric counterpart to `/spdd-block`.

**Input**: `/spdd-unblock @requirements/<story-file>.md [-m "<resolution note>"]`

The `-m` note is optional but recommended — it is appended to the closed block record as `resolution`. Examples:

```
/spdd-unblock @requirements/[User-story-7]token-billing.md
/spdd-unblock @requirements/[User-story-3]invoice-export.md -m "Quota approved by ops; proceeding"
```

**Steps**

1. **Validate input**

   a. **If no `@` story file is provided**, use the **AskUserQuestion tool** to ask:
   > "Which story should be unblocked? Provide a path under `requirements/`."

   **IMPORTANT**: Do NOT proceed without a story file. The `-m` note is optional.

2. **Read the story file and parse its frontmatter**

   - Read the entire file.
   - Confirm a YAML frontmatter block exists with the canonical schema.
   - If `status` is not `BLOCKED`, abort with a message: "Story is not currently BLOCKED (status = X). Nothing to unblock."

3. **Locate the open block record**

   - Find the LAST entry in `block_records` whose `unblocked_at` is `null`. This is the active blocker.
   - If no such entry exists (e.g., schema drift, manual edits), abort with a clear error and ask the user to fix the file by hand. Do NOT invent a record.

4. **Capture closing metadata**

   - Run `git config user.name` to capture `closer` (fall back to `$USER`, then `"unknown"`).
   - Capture `now = YYYY-MM-DD HH:MM` from the local clock.
   - Read `previous_status` from the open block record — this is the status to restore.

5. **Mutate the frontmatter** in-place via Read → StrReplace (never reorder keys, never touch the body):

   a. Set `status: <previous_status>` (the value snapshotted by `/spdd-block`).

   b. Update the open block record by setting:
    - `unblocked_at: "{now}"`
    - `closer: "{closer}"`
    - `resolution: "{-m note}"` — only add this key if `-m` was provided; do NOT add it as an empty string.

   c. Do NOT modify any other entries in `block_records` and do NOT change any other top-level field.

6. **Report the change** to the user:

   ```
   Story unblocked.
      File: requirements/<story-file>.md
      Status: BLOCKED → <previous_status>
      Block duration: <blocked_at> → <now>
      Closed by: <closer>
      Resolution: <-m note OR "(none)">
   ```

**Output**

The story file's YAML frontmatter is mutated to restore the prior `status` and to close the active `block_records` entry with `unblocked_at`, `closer`, and (optionally) `resolution`. Block history is preserved.

**Guardrails**

- Do NOT proceed without a story file
- Do NOT touch any field other than `status` and the single open `block_records` entry
- Do NOT modify any closed `block_records` entries — they are immutable audit history
- Do NOT invent or fabricate a block record if none is open — abort and ask the user
- Do NOT add `resolution: ""` when no `-m` note was provided — omit the key instead
- Always restore exactly the `previous_status` recorded by `/spdd-block` — never guess
- Always preserve the existing key order and indentation style of the frontmatter
