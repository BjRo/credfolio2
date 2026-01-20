---
# credfolio2-dmle
title: Simplify database selection via CREDFOLIO_ENV
status: completed
type: task
priority: normal
created_at: 2026-01-20T15:02:48Z
updated_at: 2026-01-20T15:11:16Z
---

Replace separate dev/test migration commands with a single set that uses CREDFOLIO_ENV (defaults to 'dev'). Database name becomes credfolio_${CREDFOLIO_ENV}.

## Changes

- [x] Update Makefile to use `CREDFOLIO_ENV` instead of separate `-test` commands
- [x] Rename databases to `credfolio_dev` and `credfolio_test`
- [x] Update init-db.sh to create both databases
- [x] Update docker-compose.yml to use postgres default DB
- [x] Update config.go to use CREDFOLIO_ENV
- [x] Update config tests
- [x] Update .env.example
- [x] Update CLAUDE.md

## Usage

```bash
# Default (dev environment)
make migrate-up                    # Runs on credfolio_dev

# Test environment
CREDFOLIO_ENV=test make migrate-up # Runs on credfolio_test

# Check status of all environments
make migrate-status
```
