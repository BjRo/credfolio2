# Decisions Context for Claude

This directory contains Architecture Decision Records (ADRs). See [README.md](README.md) for the template and guidelines on when to create decisions.

## Decision Index

| File | Title | Date | Summary |
|------|-------|------|---------|
| [20260119163958-add-decision-documentation-workflow.md](20260119163958-add-decision-documentation-workflow.md) | Add Decision Documentation Workflow | 2026-01-19 | Established ADR system with timestamped files, standard template, and `/decision` skill |
| [20260120044531-go-clean-architecture-structure.md](20260120044531-go-clean-architecture-structure.md) | Go Backend Clean Architecture Structure | 2026-01-20 | Adopted Clean Architecture with domain/repository/service/handler layers in `internal/` |
| [20260120044532-chi-router-adoption.md](20260120044532-chi-router-adoption.md) | Chi Router Adoption | 2026-01-20 | Selected chi v5 as HTTP router for stdlib compatibility and middleware support |
| [20260120134659-devcontainer-network-isolation-and-permission-model.md](20260120134659-devcontainer-network-isolation-and-permission-model.md) | Devcontainer Network Isolation and Permission Model | 2026-01-20 | Implemented iptables firewall in devcontainer with domain allowlist for AI agent sandboxing |
| [20260121140000-bun-orm-for-database-access.md](20260121140000-bun-orm-for-database-access.md) | Bun ORM for Database Access | 2026-01-21 | Adopted Bun as SQL-first ORM built on pgx for PostgreSQL access |
| [20260122001044-dnsmasq-dynamic-firewall-whitelisting.md](20260122001044-dnsmasq-dynamic-firewall-whitelisting.md) | dnsmasq for Dynamic Firewall Whitelisting | 2026-01-22 | Added dnsmasq for dynamic DNS-based firewall rules to handle Go module proxy IP changes |
| [20260123000000-shadcn-ui-component-library.md](20260123000000-shadcn-ui-component-library.md) | shadcn/ui Component Library | 2026-01-23 | Adopted shadcn/ui as the component library for the frontend |
| [20260207200000-implement-subagent-for-bean-execution.md](20260207200000-implement-subagent-for-bean-execution.md) | Implement Subagent for Bean Execution | 2026-02-07 | Introduced @implement subagent to execute bean implementation plans in isolated context, completing the refine → implement → qa → review pipeline |
| [20260208004529-add-review-ai-subagent-for-llm-code-review.md](20260208004529-add-review-ai-subagent-for-llm-code-review.md) | Add @review-ai Subagent for LLM Code Review | 2026-02-08 | Introduced @review-ai specialized reviewer for LLM code focusing on prompt injection, schema design, model selection, and evaluation strategies |

## Maintenance

When creating new decisions with `/decision`, remember to add an entry to this index table.
