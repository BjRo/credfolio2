# Devcontainer Configuration

This directory contains the development container configuration for Credfolio2.

## Git Configuration

**Important:** The git remote URL must use HTTPS (not SSH) for VSCode's git credential helper to work correctly inside the devcontainer.

### Verify your remote URL

```bash
git remote -v
```

You should see:
```
origin  https://github.com/BjRo/credfolio2.git (fetch)
origin  https://github.com/BjRo/credfolio2.git (push)
```

### If using SSH URL

If your remote is using SSH (`git@github.com:BjRo/credfolio2.git`), switch to HTTPS:

```bash
git remote set-url origin https://github.com/BjRo/credfolio2.git
```

VSCode's credential helper will handle authentication automatically with HTTPS remotes.

## Claude Code Permissions (YOLO Mode)

Claude Code runs in **bypass permissions mode** inside the devcontainer only. This allows Claude to work freely without permission prompts, which is safe because:

1. The devcontainer is isolated from the host machine
2. All work happens in a sandboxed Docker environment
3. The workspace settings outside the container remain restrictive

### Configuration

Bypass permissions is enabled through multiple layers:

1. **VSCode extension settings** ([devcontainer.json](devcontainer.json#L52-L53)):
   ```json
   "claude-code.permissions.defaultMode": "bypassPermissions",
   "claude-code.dangerouslySkipPermissions": true
   ```

2. **Claude config volume** (`~/.claude/` inside container):
   - Mounted from Docker volume: `claude-code-config-${devcontainerId}`
   - Contains `settings.json` with `defaultMode: "bypassPermissions"`
   - Persists across container rebuilds

3. **Workspace settings** (`.claude/settings.json` in project root):
   - Remain restrictive (empty allow list)
   - Only affect Claude Code when running outside the container
   - Ensure permission checks when working on the host

### Why This Works

The devcontainer sets `CLAUDE_CONFIG_DIR=/home/node/.claude`, which points to the Docker volume mount. This means:

- Inside the container: Uses volume settings with bypass enabled
- Outside the container: Uses workspace settings with restrictions
- Complete isolation between environments

### Rebuilding the Container

The bypass permissions settings persist across rebuilds because:
- VSCode settings are in `devcontainer.json` (rebuilt each time)
- Claude config is in a Docker volume (persists across rebuilds)

No manual reconfiguration needed after rebuilding.
