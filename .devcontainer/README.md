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

### How It Works

Claude Code loads settings from two scopes, merging them at runtime:

| Scope | File | Purpose |
|-------|------|---------|
| **User** | `~/.claude/settings.json` | Global permissions — bypass mode, tool allow list |
| **Project** | `/workspace/.claude/settings.json` | Hooks (`PreToolUse`, `SessionStart`, `PreCompact`), `.env` deny rules |

The **user-level** settings are baked into the Docker image from [`claude-user-settings.json`](claude-user-settings.json) via the Dockerfile:

```dockerfile
COPY --chown=node:node .devcontainer/claude-user-settings.json /home/node/.claude/settings.json
```

The **project-level** settings come from the repo checkout at `/workspace/.claude/settings.json` and are available automatically when the workspace is mounted.

### Editing Permissions

To change which tools are allowed without prompts, edit [`claude-user-settings.json`](claude-user-settings.json) and rebuild the devcontainer. The current allow list:

- `Bash(*)`, `Read(*)`, `Write(*)`, `Edit(*)` — file and shell operations
- `Grep(*)`, `Glob(*)` — search operations
- `Task(*)` — subagent operations
- `WebFetch(*)`, `WebSearch(*)` — web access

### VSCode Extension Settings

Bypass mode is also configured at the extension level in [`devcontainer.json`](devcontainer.json):

```json
"claude-code.permissions.defaultMode": "bypassPermissions",
"claude-code.dangerouslySkipPermissions": true
```

### Rebuilding the Container

No manual reconfiguration needed — user settings are baked into the image and project settings come from the repo.
