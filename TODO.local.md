# Credfolio - Project Breakdown

This document provides an overview of the project structure. **Actual work items are tracked in beans** (see `.beans/` directory).

## Technical Stack Summary

| Layer | Technology |
|-------|------------|
| Frontend | Next.js 16, React 19, TypeScript, Tailwind CSS 4, shadcn/ui, URQL |
| Backend | Go 1.24, gqlgen (GraphQL), Bun (ORM), golang-migrate |
| Database | PostgreSQL |
| Job Queue | River (Postgres-based) |
| File Storage | MinIO (local), S3-compatible (production) |
| LLM Integration | Anthropic Claude via Go gateway service |
| Auth | Simple email/password (Go-based) |
| Deployment | Docker Compose (local), container-ready for K8s |

## Development Process

- **TDD**: Write tests first (Red → Green → Refactor)
- **Iterative**: Build vertical slices, connect early, iterate
- **Branching**: Feature branches from up-to-date main, PRs for review
- **Tracking**: Use `beans` CLI to manage work items

## Epic Structure (Iterative Slices)

### Milestone 1: Foundation & First Vertical Slice

**Goal**: Upload a reference letter → extract text via LLM → display raw results

#### Epic 1.1: Infrastructure Foundation
- Docker Compose setup (Postgres, MinIO)
- Database migrations setup with golang-migrate
- Basic project structure for backend (Clean Architecture)
- Environment configuration via env vars

#### Epic 1.2: File Upload Pipeline
- File upload endpoint (accept PDF/DOCX/TXT)
- Store files in MinIO/local storage
- Job queue integration (River)
- Basic upload UI in Next.js

#### Epic 1.3: LLM Gateway Service
- LLM abstraction layer in Go
- Anthropic Claude integration
- Circuit breaker pattern
- Retry with exponential backoff
- Document text extraction via Claude vision

#### Epic 1.4: Reference Letter Processing
- Extract structured data from reference letters
- Define data schema (company, role, dates, skills, testimonials)
- Store extracted data in PostgreSQL
- Display raw extraction results in UI

### Milestone 2: Profile Display

**Goal**: Display extracted data as a beautiful, interactive profile

#### Epic 2.1: GraphQL API
- Set up gqlgen
- Define schema for profiles, positions, skills
- Implement resolvers
- Connect URQL on frontend

#### Epic 2.2: Profile UI Components
- Profile header (name, photo, summary)
- Position cards with extracted data
- Skills visualization
- Testimonial highlights
- Responsive design

#### Epic 2.3: Interactive Elements
- Animations and micro-interactions
- Expandable position details
- Skill tag filtering
- "Moments of delight" polish

### Milestone 3: Profile Editing

**Goal**: Fine-tune extracted data via LLM instructions or manual edits

#### Epic 3.1: Manual Editing
- Edit profile fields directly
- Edit position details
- Add/remove skills
- Reorder positions

#### Epic 3.2: LLM-Assisted Refinement
- "Refine with AI" prompts
- Apply instructions to regenerate sections
- Preview before applying changes
- Undo/history support

### Milestone 4: Multi-User & Authentication

**Goal**: Support multiple users with accounts

#### Epic 4.1: Authentication System
- User registration (email/password)
- Login/logout
- Password reset flow
- Session management

#### Epic 4.2: User Isolation
- User-scoped data
- Profile ownership
- Public profile URLs (`/u/username`)

### Milestone 5: Profile Cloning & Customization

**Goal**: Create job-specific profile variants

#### Epic 5.1: Profile Cloning
- Clone existing profile
- Name/describe the variant
- Link to job posting (optional)

#### Epic 5.2: Job-Tailored Customization
- Highlight relevant skills/experience
- Hide irrelevant positions
- AI-assisted tailoring based on job description
- Compare original vs tailored view

### Milestone 6: Production Readiness

**Goal**: Prepare for deployment

#### Epic 6.1: Observability
- Structured logging
- Metrics (Prometheus)
- Health checks
- Error tracking

#### Epic 6.2: Security Hardening
- Input validation
- Rate limiting
- CORS configuration
- File upload security

#### Epic 6.3: Deployment
- Production Docker images
- Kubernetes manifests (optional)
- CI/CD pipeline
- Environment configuration

---

## Commands Reference

```bash
# View all beans
beans query '{ beans { id title status type priority } }'

# Find work to do
beans query '{ beans(filter: { excludeStatus: ["completed", "scrapped", "draft"], isBlocked: false }) { id title status type } }'

# Start working on a bean
beans update <bean-id> --status in-progress

# Complete a bean
beans update <bean-id> --status completed

# Create new work
beans create "Title" -t task -d "Description..." -s todo
```

## First Steps

1. Review and approve this breakdown
2. Create the milestone and epic beans
3. Start with Epic 1.1 (Infrastructure Foundation)
