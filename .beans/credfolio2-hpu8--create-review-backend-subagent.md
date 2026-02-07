---
# credfolio2-hpu8
title: Create review-backend subagent
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T15:57:37Z
updated_at: 2026-02-07T16:59:19Z
parent: credfolio2-ynmd
blocking:
    - credfolio2-xflj
---

Convert the `review-backend` skill into a proper `.claude/agents/` subagent.

## Why

The review-backend skill currently runs in the main conversation context. Review output is verbose (PR comments, code analysis) and pollutes the main context window. Running as a subagent isolates this output — only a summary returns to the caller.

## What

- Create `.claude/agents/review-backend.md` with appropriate frontmatter
- Move the system prompt from the skill's SKILL.md into the agent's markdown body
- Set `model: inherit` (reviews need the best model available)
- Restrict tools to read-only + Bash (for `gh` PR comment posting)
- Remove the old skill (or keep as a thin wrapper that delegates, TBD)

## Definition of Done
- [x] Subagent file created at `.claude/agents/review-backend.md`
- [x] Old skill removed or converted to delegate to the subagent
- [x] Subagent can be invoked and posts review comments to a PR
- [x] Review output stays out of the main conversation context
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Branch pushed and PR created for human review

## Implementation Plan

### Approach

Create a new agent file at `.claude/agents/review-backend.md` following the exact pattern established by the QA subagent (`qa.md`) and refine agent (`refine.md`). The agent file will contain the full review prompt currently in `.claude/skills/review-backend/SKILL.md`. The old skill directory will be deleted entirely since nothing needs to preload it (unlike `agent-browser` which is preloaded by the QA agent). The dev-workflow skill will be updated to point callers at the new agent instead of the old skill path.

This is a straightforward file move + reformat, not a rewrite. The review prompt content is already well-structured and battle-tested.

### Files to Create/Modify

- `.claude/agents/review-backend.md` — **Create**: New agent file with frontmatter + full review prompt
- `.claude/skills/review-backend/SKILL.md` — **Delete**: Old skill file, no longer needed
- `.claude/skills/dev-workflow/SKILL.md` — **Modify**: Update the Task tool invocation instructions for @review-backend to reference the agent instead of the skill path

### Steps

1. **Create the agent file at `.claude/agents/review-backend.md`**

   Use the following frontmatter format (matching the pattern from `qa.md` and `refine.md`):

   ```yaml
   ---
   name: review-backend
   description: Staff-level Go/Backend code reviewer. Reviews backend code (Go, GraphQL API) for maintainability, design, performance, and security. Posts findings as PR review comments.
   tools: Read, Bash, Glob, Grep
   model: inherit
   ---
   ```

   **Frontmatter details:**
   - `name: review-backend` — matches the existing skill name for consistency
   - `description` — reuse the existing skill description (trimmed of the "Use after creating a PR..." suffix since agents have their own invocation model)
   - `tools: Read, Bash, Glob, Grep` — Read-only file access (Read, Glob, Grep) plus Bash for `gh` CLI commands to post PR review comments. This matches the QA agent's tool set exactly. No Write or Edit tools are needed since the reviewer should never modify code.
   - `model: inherit` — uses whatever model the caller is using, ensuring the reviewer gets the best available model
   - No `skills` field needed — the review-backend agent is self-contained and does not depend on any other skills

   **Body content:** Copy the entire body of `.claude/skills/review-backend/SKILL.md` (everything after the YAML frontmatter closing `---`, lines 6-173) into the agent file body. The content starts with `# Backend Code Review — Staff-Level Go Engineer` and includes all sections: Review Process, Comment Guidelines, What NOT to Review, and Project-Specific Context.

   No modifications to the prompt content are needed — it is already written in the imperative "you are..." style that works for both skills and agents.

2. **Delete the old skill directory**

   Remove the entire `.claude/skills/review-backend/` directory:
   - Delete `.claude/skills/review-backend/SKILL.md`
   - Remove the `.claude/skills/review-backend/` directory

   **Rationale for full deletion (not keeping as thin wrapper):**
   - The `agent-browser` skill was kept because the QA agent preloads it via the `skills` frontmatter field. The review-backend agent has no skills to preload — it is self-contained.
   - Keeping a thin wrapper skill that delegates to the agent would add indirection with no benefit. Claude already knows how to invoke agents directly via the Task tool.
   - The downstream bean (credfolio2-xflj) will update all references to point at the new agent.

3. **Update the dev-workflow skill's invocation instructions**

   In `.claude/skills/dev-workflow/SKILL.md`, update the section titled "### 8. Run Automated Code Reviews" (around lines 163-196). Specifically:

   **Change the @review-backend Task tool invocation** (lines 180-185) from:
   ```
   For **@review-backend**:
   Task tool call:
     subagent_type: "general-purpose"
     description: "Backend code review"
     prompt: "You are the @review-backend agent. Read the skill definition at .claude/skills/review-backend/SKILL.md and follow its instructions to review the current PR. Post your findings as PR comments using the gh CLI."
   ```

   To:
   ```
   For **@review-backend**:
   Task tool call:
     subagent: "review-backend"
     prompt: "Review the current PR. Post your findings as PR comments using the gh CLI."
   ```

   **Note:** The exact Task tool invocation syntax for named agents may differ from the `subagent_type: "general-purpose"` pattern used today. The key change is that instead of telling a generic subagent to read the skill file, we reference the named agent directly. The agent's own system prompt (from `review-backend.md`) provides all the review instructions.

   Also update line 196 which says `Never invoke /skill review-backend` — change to reflect that the skill no longer exists and the agent is the canonical path.

4. **Verify the agent can be invoked**

   After creating the agent file:
   - Confirm the file is properly formatted by reading it back
   - The agent should be invocable via the Task tool with `subagent: "review-backend"`
   - Verification that it actually posts PR comments requires a real PR, which will be available when this bean's own PR is created

5. **Run lint and tests**

   ```bash
   pnpm lint
   pnpm test
   ```

   These commands verify that no source code was broken. Since this change only touches `.claude/` configuration files, lint and test should pass trivially.

6. **Create branch, commit, and push**

   Branch name: `feat/credfolio2-hpu8-review-backend-subagent`

   Commits:
   - First commit: bean status change to in-progress
   - Second commit: add agent file, delete old skill, update dev-workflow — single atomic commit with message like `feat: convert review-backend skill to subagent`

### Testing Strategy

- **Structural verification**: Read the created agent file and confirm it has correct frontmatter fields (`name`, `description`, `tools`, `model`) and the full review prompt body
- **Skill deletion verification**: Confirm `.claude/skills/review-backend/` directory no longer exists
- **Dev-workflow verification**: Read the updated dev-workflow skill and confirm it references the agent, not the old skill path
- **Lint/test**: `pnpm lint` and `pnpm test` must pass
- **Functional verification**: Once the PR is created for this bean, the review-backend agent can be invoked against its own PR as a smoke test. This is optional but ideal.

### Notes

- This bean does NOT need to update CLAUDE.md references or existing bean checklists. That is the responsibility of the downstream bean `credfolio2-xflj` ("Update bean checklists and templates to reference new subagents"), which is blocked by this bean.
- The sibling bean `credfolio2-fvlb` ("Create review-frontend subagent") should follow the same pattern. Consistency between the two agents is important.
- No `disable-model-invocation: true` needs to be added to any existing skill (unlike the QA subagent work which added it to `agent-browser`), because the review-backend skill is being fully deleted rather than kept as a preloaded dependency.
