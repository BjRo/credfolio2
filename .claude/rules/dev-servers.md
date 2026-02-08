---
description: "How to start, stop, and troubleshoot dev servers"
---

# Starting Dev Servers

Before running `pnpm dev`, ensure no stale processes are occupying ports:

```bash
# Stop any running dev servers
.claude/scripts/stop-dev.sh

# Clear Turbopack cache if Next.js hangs (connects but never responds)
rm -rf src/frontend/.next

# Then start
pnpm dev
```

**Why this approach:**
- **`.claude/scripts/stop-dev.sh`** handles all cleanup in one command:
  - **`pkill -f`** kills by command pattern (turbo orchestrator, go backend, next frontend)
  - **`fuser -k`** as fallback catches anything else holding ports 8080 and 3000
  - **Verification step** with `lsof` confirms ports are actually free
  - **Idempotent** - safe to run even if nothing is running

**Common issues:**
- **Port already in use**: Turborepo spawns a process tree; killing just the port holder can leave orphans. Use the full pkill sequence above.
- **Frontend connects but never responds**: Turbopack's cache (`.next/`) can become corrupted. Fix: `rm -rf src/frontend/.next`
- **Backend failure kills frontend**: Turborepo tears down all tasks if one fails, but zombie processes may remain. Always run the full cleanup sequence before retrying.
