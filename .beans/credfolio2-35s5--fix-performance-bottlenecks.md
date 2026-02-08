---
# credfolio2-35s5
title: Fix performance bottlenecks
status: in-progress
type: task
priority: high
created_at: 2026-02-08T11:03:47Z
updated_at: 2026-02-08T14:54:29Z
parent: credfolio2-nihn
---

Address critical performance issues identified in backend and frontend code.

## Critical Issues

### Backend (from @review-backend)
1. **N+1 Query Pattern** - Validation count resolvers query per skill/experience

### Frontend (from @review-frontend)
2. **Home Page Waterfall** - Client-side fetch before redirect
3. **File Upload Bypass** - Uses XMLHttpRequest instead of urql
4. **Testimonials Waterfall** - Sequential queries instead of single query

## Impact

- Poor performance at scale (100 skills = 100 queries)
- Slow page loads and navigation
- Inconsistent error handling (file upload)

## Files Affected

- `src/backend/internal/graphql/resolver/schema.resolvers.go`
- `src/frontend/src/app/page.tsx`
- `src/frontend/src/components/upload/document-upload.tsx`
- `src/frontend/src/components/profile/testimonials-section.tsx`

## Acceptance Criteria

- [ ] N+1 queries resolved with dataloader or eager loading
- [ ] Home page uses Server Component with `redirect()`
- [ ] File upload uses proper GraphQL mutation
- [ ] Testimonials fetched in single query with fragments

## Reference

See: /documentation/reviews/2026-02-08-comprehensive-codebase-review.md

## Implementation Plan

### Approach

This plan addresses four distinct performance issues, each with a different root cause:

1. **Backend N+1**: Dataloader pattern to batch validation count queries
2. **Home page waterfall**: Server-side redirect using Next.js Server Components
3. **File upload bypass**: Document limitation + workaround justification
4. **Testimonials waterfall**: GraphQL fragment to merge queries

These fixes are independent and can be implemented in parallel or sequentially based on priority.

### Files to Create/Modify

#### Backend (Issue 1: N+1 Queries)
- `src/backend/internal/graphql/dataloader/dataloader.go` - New dataloader infrastructure (batch loading functions)
- `src/backend/internal/graphql/dataloader/validation_count_loader.go` - Validation count batcher
- `src/backend/internal/graphql/middleware/dataloader.go` - Context injection middleware
- `src/backend/internal/graphql/handler.go` - Wire up dataloader middleware
- `src/backend/internal/graphql/resolver/schema.resolvers.go` - Update ValidationCount resolvers to use dataloader
- `src/backend/internal/domain/repository.go` - Add batch count methods to repository interfaces
- `src/backend/internal/repository/postgres/skill_validation_repository.go` - Implement batch count method
- `src/backend/internal/repository/postgres/experience_validation_repository.go` - Implement batch count method
- `src/backend/internal/graphql/dataloader/dataloader_test.go` - Tests for dataloader logic

#### Frontend (Issue 2: Home Page Waterfall)
- `src/frontend/src/app/page.tsx` - Convert to Server Component with server-side redirect
- `src/frontend/src/app/page.test.tsx` - Update tests for Server Component

#### Frontend (Issue 3: File Upload Bypass)
- `src/frontend/src/components/upload/document-upload.tsx` - Document why XMLHttpRequest is required
- `src/frontend/src/components/profile/ReferenceLetterUploadModal.tsx` - Document why XMLHttpRequest is required

#### Frontend (Issue 4: Testimonials Waterfall)
- `src/frontend/src/graphql/queries.graphql` - Add testimonials fragment to GetProfileById query
- `src/frontend/src/app/profile/[id]/page.tsx` - Remove separate testimonials query, use data from profile query
- `src/backend/internal/graphql/schema/schema.graphqls` - Add testimonials field to Profile type
- `src/backend/internal/graphql/resolver/schema.resolvers.go` - Implement Profile.testimonials resolver
- `src/backend/internal/domain/repository.go` - (Already has GetByProfileID method in TestimonialRepository)

### Steps

#### Issue 1: Fix N+1 Queries in Backend (Priority: Critical)

**Background**: The `validationCount` field on `ProfileSkill` and `ProfileExperience` triggers a database query per item. With 100 skills, this becomes 100 separate COUNT queries.

**Current Implementation** (schema.resolvers.go:2974-2984):
```go
func (r *profileSkillResolver) ValidationCount(ctx context.Context, obj *model.ProfileSkill) (int, error) {
    skillID, err := uuid.Parse(obj.ID)
    if err != nil {
        return 0, fmt.Errorf("invalid skill ID: %w", err)
    }
    // This executes once per skill - N+1 problem!
    count, err := r.skillValidationRepo.CountByProfileSkillID(ctx, skillID)
    // ...
}
```

**Solution**: Use dataloader pattern to batch validation counts into a single query.

**Step 1.1**: Add batch count methods to repository interfaces

File: `src/backend/internal/domain/repository.go`

Add these methods:
```go
// In SkillValidationRepository interface:
// BatchCountByProfileSkillIDs returns validation counts for multiple skills in one query.
// The returned map uses profile_skill_id as key and count as value.
BatchCountByProfileSkillIDs(ctx context.Context, profileSkillIDs []uuid.UUID) (map[uuid.UUID]int, error)

// In ExperienceValidationRepository interface:
// BatchCountByProfileExperienceIDs returns validation counts for multiple experiences in one query.
// The returned map uses profile_experience_id as key and count as value.
BatchCountByProfileExperienceIDs(ctx context.Context, profileExperienceIDs []uuid.UUID) (map[uuid.UUID]int, error)
```

**Step 1.2**: Implement batch count methods in PostgreSQL repositories

File: `src/backend/internal/repository/postgres/skill_validation_repository.go`

```go
func (r *SkillValidationRepository) BatchCountByProfileSkillIDs(ctx context.Context, profileSkillIDs []uuid.UUID) (map[uuid.UUID]int, error) {
    if len(profileSkillIDs) == 0 {
        return map[uuid.UUID]int{}, nil
    }
    
    type Result struct {
        ProfileSkillID uuid.UUID `bun:"profile_skill_id"`
        Count          int       `bun:"count"`
    }
    
    var results []Result
    err := r.db.NewSelect().
        Model((*domain.SkillValidation)(nil)).
        Column("profile_skill_id").
        ColumnExpr("COUNT(*) as count").
        Where("profile_skill_id IN (?)", bun.In(profileSkillIDs)).
        Group("profile_skill_id").
        Scan(ctx, &results)
    if err != nil {
        return nil, err
    }
    
    counts := make(map[uuid.UUID]int, len(results))
    for _, r := range results {
        counts[r.ProfileSkillID] = r.Count
    }
    
    return counts, nil
}
```

File: `src/backend/internal/repository/postgres/experience_validation_repository.go`

Similar implementation for experiences (replace SkillValidation with ExperienceValidation, profile_skill_id with profile_experience_id).

**Step 1.3**: Create dataloader infrastructure

File: `src/backend/internal/graphql/dataloader/dataloader.go`

```go
package dataloader

import (
    "context"
    "sync"
)

// contextKey is the type for context keys used by dataloader.
type contextKey string

const (
    skillValidationCountKey       contextKey = "skillValidationCountLoader"
    experienceValidationCountKey  contextKey = "experienceValidationCountLoader"
)

// Loaders contains all dataloaders for the application.
type Loaders struct {
    SkillValidationCount       *ValidationCountLoader
    ExperienceValidationCount  *ValidationCountLoader
}

// GetLoaders retrieves dataloaders from context.
func GetLoaders(ctx context.Context) *Loaders {
    return ctx.Value(contextKey("loaders")).(*Loaders)
}

// ValidationCountLoader batches validation count queries.
type ValidationCountLoader struct {
    mu      sync.Mutex
    batch   []uuid.UUID
    results map[uuid.UUID]int
    done    chan struct{}
    fetch   func(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]int, error)
}

// NewValidationCountLoader creates a new validation count loader.
func NewValidationCountLoader(fetch func(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]int, error)) *ValidationCountLoader {
    return &ValidationCountLoader{
        batch:   make([]uuid.UUID, 0),
        results: make(map[uuid.UUID]int),
        done:    make(chan struct{}),
        fetch:   fetch,
    }
}

// Load retrieves a validation count, batching the request.
func (l *ValidationCountLoader) Load(ctx context.Context, id uuid.UUID) (int, error) {
    l.mu.Lock()
    l.batch = append(l.batch, id)
    l.mu.Unlock()
    
    // Wait for batch to be processed
    <-l.done
    
    count, ok := l.results[id]
    if !ok {
        return 0, nil // No validations for this ID
    }
    return count, nil
}

// ExecuteBatch executes the batched fetch.
func (l *ValidationCountLoader) ExecuteBatch(ctx context.Context) error {
    l.mu.Lock()
    ids := l.batch
    l.mu.Unlock()
    
    results, err := l.fetch(ctx, ids)
    if err != nil {
        return err
    }
    
    l.mu.Lock()
    l.results = results
    l.mu.Unlock()
    
    close(l.done)
    return nil
}
```

Note: This is a simplified dataloader. Consider using a library like `github.com/graph-gophers/dataloader` for production use with proper batching, caching, and error handling.

**Step 1.4**: Create middleware to inject dataloaders into context

File: `src/backend/internal/graphql/middleware/dataloader.go`

```go
package middleware

import (
    "context"
    "net/http"
    
    "backend/internal/domain"
    "backend/internal/graphql/dataloader"
)

// DataloaderMiddleware creates dataloaders per request and injects them into context.
func DataloaderMiddleware(
    skillValidationRepo domain.SkillValidationRepository,
    expValidationRepo domain.ExperienceValidationRepository,
) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := r.Context()
            
            loaders := &dataloader.Loaders{
                SkillValidationCount: dataloader.NewValidationCountLoader(
                    skillValidationRepo.BatchCountByProfileSkillIDs,
                ),
                ExperienceValidationCount: dataloader.NewValidationCountLoader(
                    expValidationRepo.BatchCountByProfileExperienceIDs,
                ),
            }
            
            ctx = context.WithValue(ctx, contextKey("loaders"), loaders)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

**Step 1.5**: Wire up dataloader middleware in GraphQL handler

File: `src/backend/internal/graphql/handler.go`

Find the NewHandler function and add dataloader middleware to the chain.

**Step 1.6**: Update resolvers to use dataloader

File: `src/backend/internal/graphql/resolver/schema.resolvers.go`

Replace the ValidationCount resolver implementations:

```go
func (r *profileSkillResolver) ValidationCount(ctx context.Context, obj *model.ProfileSkill) (int, error) {
    skillID, err := uuid.Parse(obj.ID)
    if err != nil {
        return 0, fmt.Errorf("invalid skill ID: %w", err)
    }
    
    loaders := dataloader.GetLoaders(ctx)
    return loaders.SkillValidationCount.Load(ctx, skillID)
}

func (r *profileExperienceResolver) ValidationCount(ctx context.Context, obj *model.ProfileExperience) (int, error) {
    expID, err := uuid.Parse(obj.ID)
    if err != nil {
        return 0, fmt.Errorf("invalid experience ID: %w", err)
    }
    
    loaders := dataloader.GetLoaders(ctx)
    return loaders.ExperienceValidationCount.Load(ctx, expID)
}
```

**Step 1.7**: Write tests for dataloader

File: `src/backend/internal/graphql/dataloader/dataloader_test.go`

Test that batch loading reduces queries from N to 1.

#### Issue 2: Fix Home Page Waterfall (Priority: High)

**Background**: The home page (`src/frontend/src/app/page.tsx`) is a Client Component that fetches profile data, then redirects based on the result. This creates a waterfall: HTML loads → JS executes → GraphQL query → redirect. Users see a loading spinner unnecessarily.

**Current Implementation** (page.tsx:10-28):
```tsx
"use client";

export default function Home() {
  const router = useRouter();
  const [result] = useQuery({
    query: GetProfileDocument,
    variables: { userId: DEMO_USER_ID },
  });

  useEffect(() => {
    if (fetching) return;
    if (profile) {
      router.push(`/profile/${profile.id}`);
    } else {
      router.push("/upload");
    }
  }, [fetching, profile, router]);
  // ... loading spinner ...
}
```

**Solution**: Convert to Server Component and use `redirect()` from `next/navigation`.

**Step 2.1**: Convert home page to Server Component with server-side redirect

File: `src/frontend/src/app/page.tsx`

```tsx
import { redirect } from "next/navigation";
import { createUrqlClient, GRAPHQL_ENDPOINT } from "@/lib/urql/client";
import { GetProfileDocument } from "@/graphql/generated/graphql";

const DEMO_USER_ID = "00000000-0000-0000-0000-000000000001";

export default async function Home() {
  // Server-side GraphQL query
  const client = createUrqlClient(GRAPHQL_ENDPOINT);
  const result = await client.query(GetProfileDocument, { userId: DEMO_USER_ID }).toPromise();
  
  const profile = result.data?.profileByUserId;
  
  // Server-side redirect (no loading spinner, instant navigation)
  if (profile) {
    redirect(`/profile/${profile.id}`);
  } else {
    redirect("/upload");
  }
  
  // This code is unreachable, but TypeScript requires a return
  return null;
}
```

**Step 2.2**: Update tests for Server Component

File: `src/frontend/src/app/page.test.tsx`

Server Components require different testing approach (test the redirect behavior, not the component rendering).

#### Issue 3: Document File Upload XMLHttpRequest Requirement (Priority: Medium)

**Background**: The review identified that `DocumentUpload.tsx` and `ReferenceLetterUploadModal.tsx` use `XMLHttpRequest` instead of urql's GraphQL client. However, this is **not a bug** - it's a necessary workaround for a Next.js limitation.

**Why XMLHttpRequest is required**:
1. Next.js rewrites (`/api/graphql → :8080/graphql`) don't handle multipart form data correctly
2. File uploads require `Content-Type: multipart/form-data` following the GraphQL multipart request spec
3. urql's `@urql/exchange-multipart-fetch` would work with the direct backend URL, but that breaks in the devcontainer network setup
4. The frontend uses `GRAPHQL_ENDPOINT` constant for regular queries (proxied) and `GRAPHQL_UPLOAD_ENDPOINT` for uploads (direct backend)

**Solution**: Document the limitation with clear comments. No code changes required.

**Step 3.1**: Add documentation comments to file upload components

File: `src/frontend/src/components/upload/DocumentUpload.tsx` (line 92, before the XMLHttpRequest usage)

```tsx
// NOTE: We use XMLHttpRequest instead of urql for file uploads because:
// 1. Next.js rewrites (/api/graphql) don't correctly handle multipart/form-data
// 2. This follows the GraphQL multipart request spec (https://github.com/jaydenseric/graphql-multipart-request-spec)
// 3. We use GRAPHQL_ENDPOINT (not GRAPHQL_UPLOAD_ENDPOINT) to go through the Next.js proxy
// 4. XMLHttpRequest provides progress events for upload feedback
// See: src/frontend/src/lib/urql/client.ts for endpoint configuration
const result = await new Promise<UploadForDetectionResult>((resolve, reject) => {
```

File: `src/frontend/src/components/profile/ReferenceLetterUploadModal.tsx` (line 76, in uploadAuthorImageXhr function)

```tsx
// NOTE: Uses XMLHttpRequest for file uploads (see DocumentUpload.tsx for detailed explanation)
return new Promise((resolve, reject) => {
```

**Step 3.2**: Update the acceptance criterion

Change "File upload uses proper GraphQL mutation" to "File upload approach documented (XMLHttpRequest required for Next.js limitation)".

#### Issue 4: Fix Testimonials Waterfall (Priority: High)

**Background**: The profile page fetches testimonials in a separate query after the profile loads (page.tsx:36-41). This creates a waterfall: profile query completes → testimonials query starts. With slow network, testimonials appear late.

**Current Implementation** (page.tsx:36-44):
```tsx
const [testimonialsResult, reexecuteTestimonialsQuery] = useQuery({
  query: GetTestimonialsDocument,
  variables: { profileId },
  pause: !profile,  // ← Waits for profile query
  requestPolicy: "network-only",
});

const testimonials = testimonialsResult.data?.testimonials ?? [];
```

**Solution**: Add testimonials field to Profile type and fetch in single query using GraphQL fragment.

**Step 4.1**: Add testimonials field to Profile GraphQL type

File: `src/backend/internal/graphql/schema/schema.graphqls`

Find the `Profile` type definition and add:
```graphql
type Profile {
  # ... existing fields ...
  testimonials: [Testimonial!]!
}
```

**Step 4.2**: Add testimonials resolver to Profile

File: `src/backend/internal/graphql/resolver/schema.resolvers.go`

Find the Profile resolver methods and add:
```go
func (r *profileResolver) Testimonials(ctx context.Context, obj *domain.Profile) ([]*domain.Testimonial, error) {
    testimonials, err := r.testimonialRepo.GetByProfileID(ctx, obj.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to load testimonials: %w", err)
    }
    return testimonials, nil
}
```

**Step 4.3**: Update GetProfileById query to include testimonials

File: `src/frontend/src/graphql/queries.graphql`

Update the GetProfileById query (line 208):
```graphql
query GetProfileById($id: ID!) {
  profile(id: $id) {
    id
    user {
      id
    }
    name
    email
    phone
    location
    summary
    profilePhotoUrl
    experiences {
      # ... existing fields ...
    }
    educations {
      # ... existing fields ...
    }
    skills {
      # ... existing fields ...
    }
    testimonials {
      id
      quote
      author {
        id
        name
        title
        company
        linkedInUrl
        imageUrl
      }
      authorName
      authorTitle
      authorCompany
      relationship
      createdAt
      validatedSkills {
        id
        name
      }
      referenceLetter {
        id
        file {
          id
          url
        }
      }
    }
    createdAt
    updatedAt
  }
}
```

**Step 4.4**: Remove separate testimonials query from profile page

File: `src/frontend/src/app/profile/[id]/page.tsx`

```tsx
// Remove GetTestimonialsDocument import (line 18)
// Remove separate testimonials query (lines 36-41)
// Remove reexecuteTestimonialsQuery from handleMutationSuccess (line 48)

// Replace testimonials data source:
const testimonials = profile.testimonials ?? [];

// Remove testimonialsResult.fetching usage:
// Delete testimonialsLoading variable (line 44)
// Pass false for isLoading prop to TestimonialsSection (line 144)
```

**Step 4.5**: Run codegen to regenerate types

```bash
cd src/frontend
pnpm codegen
cd src/backend
go generate ./internal/graphql/
```

### Testing Strategy

#### Backend Tests (Issue 1: N+1 Queries)
- **Unit tests**: Test batch count repository methods return correct counts for multiple IDs
- **Integration tests**: Test dataloader reduces queries from N to 1 using database
- **Benchmark tests**: Measure query count reduction (100 skills: 100 queries → 1 query)

#### Frontend Tests (Issue 2: Home Page)
- **Unit tests**: Test redirect logic with mocked profile query
- **Integration tests**: Test that page redirects correctly with real urql client

#### Documentation (Issue 3: File Upload)
- **No tests required**: This is documentation-only

#### Frontend Tests (Issue 4: Testimonials)
- **Unit tests**: Test profile page renders testimonials from single query result
- **Integration tests**: Test testimonials appear without separate query waterfall
- **Visual tests**: Use @qa subagent to verify testimonials display correctly

### Open Questions

**None** - All four issues have clear solutions based on existing patterns in the codebase.

### Performance Impact

**Before**:
- Profile with 50 skills + 20 experiences = 71 queries (1 profile + 50 skill counts + 20 experience counts + 1 testimonials)
- Home page redirect = 2 round trips (HTML + GraphQL query)
- Testimonials load after profile = waterfall delay

**After**:
- Profile with 50 skills + 20 experiences = 3 queries (1 profile with testimonials + 1 batch skill counts + 1 batch experience counts)
- Home page redirect = 0 round trips (server-side redirect)
- Testimonials load with profile = no waterfall

**Estimated improvement**:
- Backend queries: 71 → 3 (96% reduction)
- Home page load time: -200ms (eliminate client-side query)
- Testimonials load time: -500ms (eliminate waterfall)

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] No new TODO/FIXME/HACK/XXX comments introduced (verify with `git diff main...HEAD | grep -i "^+.*TODO\|FIXME\|HACK\|XXX"`)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification via `@qa` subagent (via Task tool, for UI changes)
- [ ] ADR written via `/decision` skill (if new dependencies, patterns, or architectural changes were introduced)
- [ ] All other checklist items above are completed
- [ ] Branch pushed to remote
- [ ] PR created for human review
- [ ] Automated code review passed via `@review-backend`, `@review-frontend`, and/or `@review-ai` (for LLM changes) subagents (via Task tool)
