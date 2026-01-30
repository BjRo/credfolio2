---
# credfolio2-o8pk
title: Session-based authentication for GraphQL mutations
status: todo
type: feature
priority: normal
created_at: 2026-01-30T10:13:41Z
updated_at: 2026-01-30T10:13:41Z
parent: credfolio2-wxn8
---

Currently, user-scoped GraphQL mutations and queries accept a `userId` parameter directly in the request. This is a security concern because:

1. **No authentication enforcement**: Any client can specify any user ID and potentially access or modify another user's data
2. **Authorization bypass risk**: Without proper checks, malicious actors could enumerate user IDs and access sensitive data

## Proposed Solutions

### Option A: Session-based user identification (Recommended)
- Get the current user ID from the authenticated session/JWT token context
- Remove `userId` parameter from mutations/queries that operate on "current user" data
- The resolver extracts user identity from the request context

### Option B: Policy-based authorization system
- Keep `userId` parameter but implement authorization policies
- Policies determine whether the requesting user can access/modify the target user's data
- Useful for admin users or delegated access scenarios

## Affected Mutations/Queries

- `uploadFile(userId, file)` → should use session user
- `uploadResume(userId, file)` → should use session user  
- `createExperience(userId, input)` → should use session user
- `createEducation(userId, input)` → should use session user
- `createSkill(userId, input)` → should use session user
- `applyReferenceLetterValidations(userId, input)` → should use session user
- `files(userId)` → should use session user or require authorization
- `referenceLetters(userId)` → should use session user or require authorization
- `resumes(userId)` → should use session user or require authorization
- `profile(userId)` → may need public access for profile viewing

## Checklist

- [ ] Design authentication/session handling approach
- [ ] Implement session context extraction in GraphQL middleware
- [ ] Update mutations to use session user instead of parameter
- [ ] Update queries to use session user or implement authorization
- [ ] Add authorization checks for cross-user access where needed
- [ ] Update frontend to not pass userId (use session)
- [ ] Add tests for authorization scenarios

## Definition of Done

- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All checklist items above are completed
- [ ] Branch pushed and PR created for human review