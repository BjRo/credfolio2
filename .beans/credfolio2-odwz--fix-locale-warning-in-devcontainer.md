---
# credfolio2-odwz
title: Fix locale warning in devcontainer
status: in-progress
type: task
priority: normal
created_at: 2026-02-08T10:06:18Z
updated_at: 2026-02-08T10:09:19Z
parent: credfolio2-ynmd
---

Fix the persistent 'setlocale: LC_ALL: cannot change locale (en_US.UTF-8)' warning that appears on bash commands in the devcontainer.

## Problem
Every bash command shows:
```
bash: warning: setlocale: LC_ALL: cannot change locale (en_US.UTF-8)
```

This happens because the locale isn't properly generated or configured in the devcontainer.

## Solution Approach
1. Generate the en_US.UTF-8 locale in the container
2. Set it as the default locale
3. Update devcontainer configuration to ensure it persists across rebuilds

## Checklist
- [x] Verify current locale settings with `locale` and `locale -a`
- [x] Update `.devcontainer/devcontainer.json` or Dockerfile to generate locales
- [ ] Test that the warning no longer appears after rebuild
- [ ] Verify existing functionality still works

## Implementation Plan

### Approach
The locale warning occurs because Debian Bookworm Slim doesn't include the `locales` package or generated locale files by default. The environment variables `LANG` and `LC_ALL` are set to `en_US.UTF-8`, but the actual locale data files don't exist.

The fix requires:
1. Installing the `locales` package in the Dockerfile
2. Generating the `en_US.UTF-8` locale during container build
3. Setting the locale environment variable in the Dockerfile

This approach ensures the fix persists across container rebuilds, as it's baked into the image.

### Files to Create/Modify
- `.devcontainer/Dockerfile` — Add locale generation steps after the initial apt-get install section (around line 9-53)

### Steps

1. **Add locales package installation**
   - In `.devcontainer/Dockerfile`, add `locales` to the existing apt-get install command on line 9
   - Position: Add it alphabetically in the list (after `less` or near other system packages)
   - This package provides `/usr/sbin/locale-gen` and locale definition files

2. **Generate en_US.UTF-8 locale**
   - After the first `apt-get install` block (after line 53, before line 55), add new RUN command:
   ```dockerfile
   # Generate en_US.UTF-8 locale
   RUN sed -i '/en_US.UTF-8/s/^# //g' /etc/locale.gen && \
       locale-gen en_US.UTF-8
   ```
   - This uncomments the `en_US.UTF-8` line in `/etc/locale.gen` and generates the locale
   - Alternative approach (more explicit): `RUN echo "en_US.UTF-8 UTF-8" > /etc/locale.gen && locale-gen`

3. **Set default locale environment variable**
   - After the locale-gen command, add:
   ```dockerfile
   ENV LANG=en_US.UTF-8
   ```
   - This ensures `LANG` is set at the image level (currently it's only set in containerEnv in devcontainer.json)

4. **Rebuild the devcontainer**
   - Command: Rebuild the devcontainer via VS Code command palette or by restarting the devcontainer
   - The Dockerfile changes will be applied during build

5. **Verify the fix**
   - Run `locale` — should show no warnings and all locales set to `en_US.UTF-8`
   - Run `locale -a` — should include `en_US.utf8` in the list
   - Run any bash command (e.g., `ls`, `echo test`) — should produce no locale warnings
   - Test existing functionality: `pnpm dev`, `pnpm test`, etc.

### Testing Strategy
- **Manual verification**: After rebuild, run `locale` and `locale -a` to confirm locale is available
- **Integration test**: Run `pnpm dev` to ensure backend and frontend start without locale warnings
- **Regression test**: Run `pnpm lint` and `pnpm test` to ensure no existing functionality is broken

### Implementation Notes
- The `locales` package is small (approximately 4MB) and is standard for Debian systems
- Locale generation happens at build time, not runtime, so there's no performance impact
- The fix follows Debian best practices for locale configuration
- No changes needed to devcontainer.json — the `LANG` environment variable there can remain, but the Dockerfile ENV takes precedence and ensures consistency

## Definition of Done
- [x] Tests written (TDD: write tests before implementation) — N/A: Infrastructure config
- [x] `pnpm lint` passes with no errors — N/A: Dockerfile change
- [x] `pnpm test` passes with no failures — N/A: Infrastructure config
- [x] Visual verification via `@qa` subagent (via Task tool, for UI changes) — N/A: No UI changes
- [x] ADR written via `/decision` skill (if new dependencies, patterns, or architectural changes were introduced) — N/A: Minor config fix
- [ ] All other checklist items above are completed — Pending devcontainer rebuild
- [x] Branch pushed to remote
- [x] PR created for human review
- [ ] Automated code review passed via `@review-backend`, `@review-frontend`, and/or `@review-ai` (for LLM changes) subagents (via Task tool) — N/A: No code changes
