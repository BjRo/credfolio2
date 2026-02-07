---
paths:
  - ".devcontainer/**"
---

# Devcontainer Notes

- Based on `debian:bookworm-slim` image
- **Tool management via mise**: Go 1.24.1 and Node 20 installed using mise (defined in mise.toml)
- Mise provides version consistency between local dev and devcontainer
- Other tools installed manually: golangci-lint, migrate, beans, delta, agent-browser
- Includes: git, gh, fzf, zsh, postgresql-client, and development tools
- Rebuild required when Dockerfile changes
