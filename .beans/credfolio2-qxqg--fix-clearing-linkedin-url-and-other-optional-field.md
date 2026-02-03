---
# credfolio2-qxqg
title: Fix clearing LinkedIn URL (and other optional fields) in author edit
status: completed
type: bug
created_at: 2026-02-03T13:50:44Z
updated_at: 2026-02-03T13:50:44Z
---

## Problem

When clearing the LinkedIn URL field in the AuthorEditModal and saving, the LinkedIn icon still appears. The URL is not actually cleared in the database.

## Root Cause

The frontend sends `null` when clearing optional fields like `linkedInUrl`, `title`, and `company`:
```javascript
linkedInUrl: formData.linkedInUrl.trim() || null,
```

In gqlgen (Go GraphQL), explicit `null` becomes a `nil` pointer, which is indistinguishable from "field not provided". The backend resolver only updates when the pointer is non-nil:
```go
if input.LinkedInURL != nil {
    author.LinkedInURL = input.LinkedInURL
}
```

The `imageId` field uses the correct pattern - it sends an empty string `""` to signal "clear this field", and the backend checks for that.

## Solution

1. **Frontend**: Send empty string instead of null for optional string fields
2. **Backend**: Treat empty string as "clear this field" (consistent with imageId pattern)

## Checklist

- [x] Update AuthorEditModal to send empty string instead of null for optional fields
- [x] Update backend UpdateAuthor resolver to clear fields when empty string is received
- [x] Add test for clearing LinkedIn URL
- [x] `pnpm lint` passes
- [x] `pnpm test` passes
- [x] Visual verification with agent-browser
