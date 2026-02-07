# Implement Subagent for Bean Execution

**Date**: 2026-02-07
**Bean**: credfolio2-g3f0

## Context

The agentic dev workflow already has dedicated subagents for planning (refine), visual verification (qa), and code review (review-backend, review-frontend). However, the most context-intensive phase — actual implementation (reading files, writing code, running tests, iterating) — was happening directly in the main conversation. This consumed the context window with low-level details and made it harder to maintain strategic oversight.

## Decision

Introduced an `@implement` subagent that executes a bean's implementation plan in an isolated context. The agent:

- Reads the bean and its implementation plan (produced by the `@refine` agent)
- Sets up the feature branch if needed
- Follows TDD strictly (preloads the TDD skill)
- Commits frequently with bean checklist updates
- Runs lint and tests
- Reports a summary when done

The agent explicitly does NOT: create PRs, launch QA/review agents, mark beans complete, or push branches. These orchestration steps remain in the main conversation.

A corresponding `/implement <bean-id>` skill provides a convenient invocation shorthand.

## Reasoning

- **Context management**: Implementation is the most token-heavy phase. Isolating it preserves the main conversation for high-level steering and decision-making.
- **Pipeline completeness**: Fills the gap in refine → **implement** → qa → review.
- **Composability**: The main conversation can run the implement agent in the background and continue working, or chain it with other agents.
- **Prerequisite quality**: The refine agent already produces detailed plans with specific file paths, steps, and testing strategies — exactly the input the implement agent needs.

Alternatives considered:
- Running everything in the main conversation (status quo): works but wastes context on low-level details.
- Full pipeline agent (implement + QA + PR in one agent): too much scope, reduces the user's ability to steer between phases.

## Consequences

- Beans should be refined (via `@refine`) before being passed to `@implement`. Poorly specified beans will produce poor results.
- The main conversation takes on an orchestration role: launch refine → review plan → launch implement → launch QA → create PR → launch reviews.
- The `AskUserQuestion` tool is available to the implement agent for surfacing ambiguities, but the interaction is less fluid than in the main conversation.
- The dev-workflow skill now documents this option in step 4 (Develop Using TDD).
