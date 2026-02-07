---
description: "How to start, stop, and troubleshoot dev servers"
---

# Starting Dev Servers

Before running `pnpm dev`, ensure no stale processes are occupying ports:

```bash
# Kill dev server processes (both the orchestrator and spawned processes)
pkill -f "turbo run dev" 2>/dev/null
pkill -f "go run cmd/server" 2>/dev/null
pkill -f "next dev" 2>/dev/null
fuser -k 8080/tcp 3000/tcp 2>/dev/null
sleep 2

# Verify ports are free (re-run the above if this shows processes)
lsof -i :8080 -i :3000 2>/dev/null || echo "Ports are free"

# Clear Turbopack cache if Next.js hangs (connects but never responds)
rm -rf src/frontend/.next

# Then start
pnpm dev
```

**Why this approach:**
- **`pkill -f`** kills by command pattern, catching both parent and child processes
- **Three-pronged kill**: turbo (orchestrator), go run (backend), and next dev (frontend)
- **`fuser -k`** as fallback catches anything else holding the ports
- **Verification step** with `lsof` confirms ports are actually free before starting

**Common issues:**
- **Port already in use**: Turborepo spawns a process tree; killing just the port holder can leave orphans. Use the full pkill sequence above.
- **Frontend connects but never responds**: Turbopack's cache (`.next/`) can become corrupted. Fix: `rm -rf src/frontend/.next`
- **Backend failure kills frontend**: Turborepo tears down all tasks if one fails, but zombie processes may remain. Always run the full cleanup sequence before retrying.
