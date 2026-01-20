---
# credfolio2-poa3
title: Authentication System
status: draft
type: epic
created_at: 2026-01-20T11:25:23Z
updated_at: 2026-01-20T11:25:23Z
parent: credfolio2-fq2i
---

Implement user registration and login functionality.

## Goals
- Simple email/password authentication in Go
- Secure password hashing and session management
- Standard auth flows (register, login, logout, reset)

## Checklist
- [ ] Create users table with auth fields
- [ ] Implement password hashing (bcrypt/argon2)
- [ ] Create registration endpoint
- [ ] Create login endpoint with session/JWT
- [ ] Create logout endpoint
- [ ] Implement password reset flow
- [ ] Add session middleware
- [ ] Create auth UI pages (login, register, forgot password)