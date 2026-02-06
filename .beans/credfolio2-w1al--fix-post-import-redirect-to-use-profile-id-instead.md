---
# credfolio2-w1al
title: Fix post-import redirect to use profile ID instead of user ID
status: todo
type: bug
priority: normal
created_at: 2026-02-05T23:16:29Z
updated_at: 2026-02-06T07:29:15Z
parent: credfolio2-3ram
---

## Summary

After importing extracted data, the upload flow redirects to `/profile/${DEMO_USER_ID}` (the hardcoded user UUID). The profile page route is `/profile/[id]` where `[id]` is the **profile ID** (not user ID). The import mutation already returns `profile.id` but it's ignored.

## Current Behavior

`UploadFlow.tsx:70-72`:
```tsx
const handleImportComplete = (_profileId: string) => {
    router.push(\`/profile/\${DEMO_USER_ID}\`);  // Uses hardcoded user ID
```

The `profileId` parameter (returned from `importDocumentResults` mutation) is prefixed with underscore and unused.

## Expected Behavior

```tsx
const handleImportComplete = (profileId: string) => {
    router.push(\`/profile/\${profileId}\`);  // Uses actual profile ID
```

## Checklist

- [ ] Use `profileId` parameter in `handleImportComplete` instead of `DEMO_USER_ID`
- [ ] Update test assertion to verify correct profile ID in redirect URL

### Definition of Done
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review