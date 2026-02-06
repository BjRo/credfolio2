---
# credfolio2-8w6a
title: Change review skills to submit proper GitHub PR reviews
status: in-progress
type: task
created_at: 2026-02-06T22:24:30Z
updated_at: 2026-02-06T22:24:30Z
---

Update /review-frontend and /review-backend skills to submit a proper (non-blocking) GitHub Pull Request Review instead of posting individual standalone comments.

## Current Behavior
- Inline comments posted individually via `gh api repos/{owner}/{repo}/pulls/{pr_number}/comments` (standalone, not grouped)
- Summary posted via `gh pr review --comment` (separate from inline comments)

## Desired Behavior
- All inline comments collected during review
- Submitted as a single GitHub Pull Request Review via the Reviews API
- Uses `event: COMMENT` (non-blocking)
- Inline comments + summary body grouped in one review

## Checklist
- [x] Update review-backend SKILL.md section 4 (Post Review Comments)
- [x] Update review-frontend SKILL.md section 4 (Post Review Comments)
- [x] Ensure both skills instruct the agent to collect findings first, then submit as a single review
- [x] Verify instructions use the correct GitHub API endpoint and JSON format

## Definition of Done
- [x] Both skill files updated
- [x] Instructions are clear and correct for an AI agent to follow
- [x] Branch pushed and PR created for human review