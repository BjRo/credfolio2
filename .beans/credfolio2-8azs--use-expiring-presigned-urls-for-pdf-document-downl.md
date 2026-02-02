---
# credfolio2-8azs
title: Use expiring presigned URLs for PDF document downloads
status: completed
type: task
priority: normal
created_at: 2026-02-02T13:17:30Z
updated_at: 2026-02-02T13:29:55Z
---

The File.url resolver currently uses GetPublicURL which returns permanent proxy URLs when storageProxyURL is configured. The GraphQL schema explicitly states the URL should expire.

## Changes
- Modify `fileResolver.URL` in schema.resolvers.go to use `GetPresignedURL` instead of `GetPublicURL`
- This ensures File downloads (including PDFs) get time-limited signed URLs
- Profile photos continue using GetPublicURL (permanent proxy URLs) - they're resolved separately

## Definition of Done
- [x] Update fileResolver.URL to use GetPresignedURL
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review