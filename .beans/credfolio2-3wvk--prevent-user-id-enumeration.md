---
# credfolio2-3wvk
title: Prevent user ID enumeration
status: scrapped
type: task
priority: normal
created_at: 2026-01-30T15:59:57Z
updated_at: 2026-02-06T23:14:26Z
parent: credfolio2-abtx
---

Ensure user IDs cannot be enumerated by attackers through API responses, timing differences, or predictable patterns.

## Rationale

User ID enumeration allows attackers to:
- Discover valid user accounts for targeted attacks
- Build lists of users for credential stuffing
- Identify high-value targets (e.g., admin accounts)

## Attack Vectors to Address

1. **Response differences**: Different error messages for "user not found" vs "wrong password"
2. **Timing attacks**: Different response times when user exists vs doesn't exist
3. **Sequential IDs**: Predictable user IDs that can be guessed (1, 2, 3...)
4. **API enumeration**: Endpoints that reveal user existence (e.g., /api/users/{id})

## Checklist

- [ ] Audit all auth endpoints for response message differences
- [ ] Ensure login/password reset return identical responses regardless of user existence
- [ ] Add constant-time comparison for sensitive lookups where applicable
- [ ] Verify UUIDs are used for user IDs (not sequential integers)
- [ ] Audit public API endpoints that accept user IDs
- [ ] Add rate limiting to prevent brute-force enumeration
- [ ] Write tests verifying consistent responses for valid/invalid users

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review