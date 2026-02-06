---
# credfolio2-k3cl
title: Refactor profile page to route by profile ID instead of resume ID
status: in-progress
type: task
priority: normal
created_at: 2026-02-06T09:16:34Z
updated_at: 2026-02-06T09:30:43Z
parent: 3ram
---

## Summary
Refactor /profile/[id] from resume ID routing to profile ID routing. Add profile(id: ID!) query, simplify ProfileHeader, update all redirects.

## Checklist
- [x] Add profile(id: ID!) and rename profile(userId) to profileByUserId(userId: ID!) in GraphQL schema
- [x] Implement resolvers with shared loadProfileData helper
- [x] Add GetProfileById frontend query and regenerate codegen
- [x] Rewrite profile page to use profile ID directly
- [x] Simplify ProfileHeader to accept profile prop instead of data/profileOverrides split
- [x] Update ExtractionReview to pass profile ID from import result
- [x] Update home page to query profileByUserId and redirect by profile ID
- [x] Update ResumeUpload redirect to fetch profile by userId then redirect by profile ID
- [x] Update reference letter preview page to use profileId instead of resumeId
- [x] Update ProfileHeader tests for new profile prop interface
- [x] Update home page tests for new profileByUserId query pattern
- [x] Remove unused ProfileData type export

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed