---
# credfolio2-duun
title: Fix Makefile help output showing wrong target names
status: todo
type: bug
created_at: 2026-02-07T16:30:46Z
updated_at: 2026-02-07T16:30:46Z
parent: credfolio2-ynmd
---

The `make help` output in `src/backend/Makefile` shows "make Makefile" for every command instead of the actual target names.

## Current behavior

```
make Makefile           Show this help
make Makefile           Create a new migration (usage: make migration name=create_users)
make Makefile           Run all pending migrations
make Makefile           Rollback the last migration
...
```

## Expected behavior

```
make help               Show this help
make migration          Create a new migration (usage: make migration name=create_users)
make migrate-up         Run all pending migrations
make migrate-down       Rollback the last migration
...
```

## Impact

- Confusing for developers reading the output
- Claude parses this and tries to run `make Makefile` instead of the actual targets, causing errors and wasted turns
- Falls under devcontainer reliability — basic tooling should work correctly

## Fix

Inspect the `help` target in `src/backend/Makefile`. The issue is likely in the sed/awk/grep pattern that extracts target names from the Makefile comments. Common causes:
- `$@` or `$$@` variable expansion issue
- Incorrect regex in the help parser
- Missing or malformed `##` comment annotations on targets

## Definition of Done
- [ ] `make help` (from `src/backend/`) shows correct target names
- [ ] All documented targets work as shown in help output
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review