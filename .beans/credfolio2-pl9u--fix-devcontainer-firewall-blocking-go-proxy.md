---
# credfolio2-pl9u
title: Fix devcontainer firewall blocking Go proxy
status: completed
type: bug
priority: normal
created_at: 2026-01-20T13:35:34Z
updated_at: 2026-01-20T13:44:45Z
---

The pre-commit hook blocks commits because backend tests fail. The Go chi router dependency cannot be downloaded because proxy.golang.org is blocked by the devcontainer firewall.

## Root Cause
The firewall in .devcontainer/init-firewall.sh does not include proxy.golang.org in the allowed domains list.

## Checklist
- [x] Add proxy.golang.org to allowed domains in init-firewall.sh
- [x] Rebuild devcontainer to apply changes
- [x] Verify tests pass with pnpm test