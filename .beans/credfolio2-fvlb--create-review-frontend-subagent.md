---
# credfolio2-fvlb
title: Create review-frontend subagent
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T15:57:44Z
updated_at: 2026-02-07T17:08:42Z
parent: credfolio2-ynmd
blocking:
    - credfolio2-xflj
---

Convert the `review-frontend` skill into a proper `.claude/agents/` subagent, and preload the frontend-specific review skills.

## Why

Same as review-backend: review output is verbose and should be isolated from the main context. Additionally, the `web-design-guidelines` and `vercel-react-best-practices` skills only matter for frontend review — they should be preloaded into this subagent's context rather than being discoverable by the main conversation.

## What

- Create `.claude/agents/review-frontend.md` with appropriate frontmatter
- Move the system prompt from the skill's SKILL.md into the agent's markdown body
- Preload skills via the `skills` frontmatter field:
  - `web-design-guidelines`
  - `vercel-react-best-practices`
- Set `model: inherit`
- Restrict tools to read-only + Bash (for `gh` PR comment posting)
- Remove the old skill (or keep as thin wrapper, TBD)
- Consider whether `web-design-guidelines` and `vercel-react-best-practices` should have `user-invocable: false` and `disable-model-invocation: true` since they're now only consumed by this subagent

## Definition of Done
- [ ] Subagent file created at `.claude/agents/review-frontend.md`
- [ ] `web-design-guidelines` and `vercel-react-best-practices` skills preloaded
- [ ] Old skill removed or converted to delegate to the subagent
- [ ] Subagent can be invoked and posts review comments to a PR
- [ ] Review output stays out of the main conversation context
- [ ] Skills `web-design-guidelines` and `vercel-react-best-practices` reviewed for invocation settings
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review

## Implementation Plan

### Approach

Follow the exact pattern established by `credfolio2-hpu8` (review-backend subagent conversion). Create a new agent file at `.claude/agents/review-frontend.md` containing the full review prompt currently in `.claude/skills/review-frontend/SKILL.md`. The key difference from review-backend is that this agent also preloads two skills (`web-design-guidelines` and `vercel-react-best-practices`) via the `skills` frontmatter field, following the precedent set by the QA agent which preloads `agent-browser`. The old skill directory will be fully deleted. The dev-workflow skill will be updated so the `@review-frontend` invocation uses the named agent directly instead of a general-purpose subagent reading the skill file.

Additionally, `web-design-guidelines` and `vercel-react-best-practices` will have `disable-model-invocation: true` added to their frontmatter, following the same precedent as `agent-browser` (which has this flag because it is preloaded by the QA agent rather than invoked directly).

### Files to Create/Modify

- `.claude/agents/review-frontend.md` — **Create**: New agent file with frontmatter + full review prompt body
- `.claude/skills/review-frontend/SKILL.md` — **Delete**: Old skill file, replaced by the agent
- `.claude/skills/web-design-guidelines/SKILL.md` — **Modify**: Add `disable-model-invocation: true` to frontmatter
- `.claude/skills/vercel-react-best-practices/SKILL.md` — **Modify**: Add `disable-model-invocation: true` to frontmatter
- `.claude/skills/dev-workflow/SKILL.md` — **Modify**: Update the `@review-frontend` Task tool invocation to reference the named agent

### Steps

1. **Create the agent file at `.claude/agents/review-frontend.md`**

   Use the following frontmatter (matching the pattern from `review-backend.md` and `qa.md`):

   ```yaml
   ---
   name: review-frontend
   description: Staff-level React/Next.js Frontend code reviewer. Reviews frontend code for best practices, accessibility, performance, and maintainability. Posts findings as PR review comments.
   tools: Read, Bash, Glob, Grep
   model: inherit
   skills:
     - web-design-guidelines
     - vercel-react-best-practices
   ---
   ```

   **Frontmatter details:**
   - `name: review-frontend` — matches the existing skill name for consistency
   - `description` — reuse the existing skill description, trimmed of the "Use after creating a PR..." trigger phrase since agents have their own invocation model
   - `tools: Read, Bash, Glob, Grep` — Read-only file access (Read, Glob, Grep) plus Bash for `gh` CLI commands to post PR review comments. Matches review-backend agent's tool set exactly. No Write or Edit tools — the reviewer should never modify code.
   - `model: inherit` — uses whatever model the caller is using
   - `skills` — preloads `web-design-guidelines` and `vercel-react-best-practices` so the agent has access to the full React/Next.js performance rules and web design guidelines without the main conversation needing to load them

   **Note on `WebFetch`:** The `web-design-guidelines` skill references `WebFetch` to fetch external guidelines from GitHub. The agent's tool set does not include `WebFetch`. However, this is acceptable because: (a) the skill content itself describes the guidelines source and approach, (b) the agent can use `Bash` with `curl` as a fallback if needed, and (c) in the devcontainer environment, external network access is restricted anyway (per CLAUDE.md). The review-frontend agent should rely on the preloaded skill content rather than fetching external URLs at review time.

   **Body content:** Copy the entire body of `.claude/skills/review-frontend/SKILL.md` (everything after the YAML frontmatter closing `---`, lines 6-181) into the agent file body. The content starts with `# Frontend Code Review — Staff-Level React/Next.js Engineer` and includes all sections: Review Process, Comment Guidelines, What NOT to Review, and Project-Specific Context.

   No modifications to the prompt content are needed — it is already written in the imperative "you are..." style that works for both skills and agents.

2. **Delete the old skill directory**

   Remove the entire `.claude/skills/review-frontend/` directory:
   - Delete `.claude/skills/review-frontend/SKILL.md`
   - Remove the `.claude/skills/review-frontend/` directory

   **Rationale for full deletion (not keeping as thin wrapper):**
   - The review-frontend agent is now the canonical entry point, just like review-backend.
   - Unlike `agent-browser` (which is kept as a skill because the QA agent preloads it), `review-frontend` is not preloaded by any other agent — it IS the agent.
   - The downstream bean `credfolio2-xflj` will update all references to point at the new agent.

3. **Add `disable-model-invocation: true` to `web-design-guidelines`**

   Edit `.claude/skills/web-design-guidelines/SKILL.md` frontmatter to add:
   ```yaml
   ---
   name: web-design-guidelines
   description: Review UI code for Web Interface Guidelines compliance. Use when asked to "review my UI", "check accessibility", "audit design", "review UX", or "check my site against best practices".
   disable-model-invocation: true
   metadata:
     author: vercel
     version: "1.0.0"
     argument-hint: <file-or-pattern>
   ---
   ```

   **Rationale:** This skill is now consumed exclusively by the review-frontend subagent via the `skills` preload mechanism. There is no reason for the main conversation's model to auto-invoke it — that would waste context and duplicate what the subagent already does. Adding `disable-model-invocation: true` follows the exact precedent set by `agent-browser` (which has this flag because it is preloaded by the QA agent). Users can still manually invoke it with `/skill web-design-guidelines` if they want to.

4. **Add `disable-model-invocation: true` to `vercel-react-best-practices`**

   Edit `.claude/skills/vercel-react-best-practices/SKILL.md` frontmatter to add:
   ```yaml
   ---
   name: vercel-react-best-practices
   description: React and Next.js performance optimization guidelines from Vercel Engineering. This skill should be used when writing, reviewing, or refactoring React/Next.js code to ensure optimal performance patterns. Triggers on tasks involving React components, Next.js pages, data fetching, bundle optimization, or performance improvements.
   disable-model-invocation: true
   license: MIT
   metadata:
     author: vercel
     version: "1.0.0"
   ---
   ```

   **Rationale:** Same as step 3. This skill's description is very broad ("triggers on tasks involving React components...") which means the model would frequently auto-invoke it during normal frontend work, adding substantial context (45 rules across 8 categories) when it is not needed. Now that it is preloaded by the review-frontend agent, it should only appear in that context. Users can still manually invoke it with `/skill vercel-react-best-practices`.

5. **Update the dev-workflow skill's invocation instructions**

   In `.claude/skills/dev-workflow/SKILL.md`, update section "### 8. Run Automated Code Reviews" (around lines 163-196). Specifically:

   **Change the @review-frontend Task tool invocation** (lines 188-194) from:
   ```
   For **@review-frontend**:
   ```
   Task tool call:
     subagent_type: "general-purpose"
     description: "Frontend code review"
     prompt: "You are the @review-frontend agent. Read the skill definition at .claude/skills/review-frontend/SKILL.md and follow its instructions to review the current PR. Post your findings as PR comments using the gh CLI."
   ```
   ```

   To:
   ```
   For **@review-frontend**:
   ```
   Task tool call:
     subagent: "review-frontend"
     description: "Frontend code review"
     prompt: "Review the current PR. Post your findings as PR comments using the gh CLI."
   ```
   ```

   **Also update the IMPORTANT note on line 196** which currently says:
   ```
   **IMPORTANT**: Always launch these as subagents via the Task tool. Never invoke review skills directly in the main conversation — that defeats the purpose of keeping the context clean. `@review-backend` is a named agent (`.claude/agents/review-backend.md`); `@review-frontend` is invoked via a general-purpose subagent reading its skill definition.
   ```

   Change to:
   ```
   **IMPORTANT**: Always launch these as subagents via the Task tool. Never invoke review skills directly in the main conversation — that defeats the purpose of keeping the context clean. Both `@review-backend` and `@review-frontend` are named agents (`.claude/agents/review-backend.md` and `.claude/agents/review-frontend.md`).
   ```

6. **Verify the agent can be invoked**

   After creating the agent file:
   - Confirm the file is properly formatted by reading it back
   - Confirm the `skills` field correctly references both preloaded skills
   - Confirm the old skill directory is gone
   - The agent should be invocable via the Task tool with `subagent: "review-frontend"`
   - Verification that it actually posts PR comments requires a real PR, which will be available when this bean's own PR is created

7. **Run lint and tests**

   ```bash
   pnpm lint
   pnpm test
   ```

   These commands verify that no source code was broken. Since this change only touches `.claude/` and `.beans/` configuration files, lint and test should pass trivially.

8. **Create branch, commit, and push**

   Branch name: `feat/credfolio2-fvlb-review-frontend-subagent`

   Commits:
   - First commit: bean status change to in-progress
   - Second commit: add agent file, delete old skill, update skill frontmatter, update dev-workflow — single atomic commit with message like `feat: convert review-frontend skill to subagent`

### Testing Strategy

- **Structural verification**: Read the created agent file and confirm it has correct frontmatter fields (`name`, `description`, `tools`, `model`, `skills`) and the full review prompt body
- **Skills preload verification**: Confirm `skills` field lists both `web-design-guidelines` and `vercel-react-best-practices`
- **Skill deletion verification**: Confirm `.claude/skills/review-frontend/` directory no longer exists
- **Skill invocation settings**: Confirm both `web-design-guidelines` and `vercel-react-best-practices` now have `disable-model-invocation: true` in their frontmatter
- **Dev-workflow verification**: Read the updated dev-workflow skill and confirm it references the named `review-frontend` agent, not the old skill path or general-purpose subagent pattern
- **Lint/test**: `pnpm lint` and `pnpm test` must pass
- **Functional verification**: Once the PR is created for this bean, the review-frontend agent can be invoked against its own PR as a smoke test. The review will find no frontend changes (since only `.claude/` files are modified), but it validates that the agent starts, reads the PR, and posts a comment.

### Blockers

None. The bean has no `blockedBy` relationships. The predecessor bean `credfolio2-hpu8` (review-backend subagent) is already completed, establishing the pattern to follow.

### Notes

- This bean does NOT need to update CLAUDE.md references or existing bean checklists. That is the responsibility of the downstream bean `credfolio2-xflj` ("Update bean checklists and templates to reference new subagents"), which is blocked by this bean.
- Consistency with `review-backend.md` is important. The frontmatter structure, tool set, and overall approach should mirror what was done for the backend agent. The only structural difference is the `skills` field (review-backend has no preloaded skills; review-frontend preloads two).
- The `vercel-react-best-practices` skill has substantial supporting files (45+ rule files in `/workspace/.claude/skills/vercel-react-best-practices/rules/`). These files stay in place since the skill itself is not being deleted, only having its invocation settings changed. The skill (and its supporting files) will be preloaded into the review-frontend agent's context.
- The `web-design-guidelines` skill mentions fetching from an external URL. Since the agent runs with `Bash` access, `curl` can serve as a fallback for `WebFetch`. In practice, the devcontainer may have limited external access, so the skill's presence as preloaded context (describing what to check) is the primary value, not the external fetch.
