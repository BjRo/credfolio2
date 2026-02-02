---
# credfolio2-ibxx
title: Add resilience to MinIO storage operations
status: todo
type: task
priority: high
created_at: 2026-02-02T08:35:46Z
updated_at: 2026-02-02T08:39:28Z
parent: credfolio2-abtx
---

The MinIO storage layer lacks retry logic and timeouts, causing random failures under network issues and potential hangs on slow operations.

## Problem

In `src/backend/internal/infrastructure/storage/minio.go`, storage operations have no resilience:

1. **No retry logic** - A single network hiccup causes immediate failure
2. **No operation timeouts** - Large file uploads/downloads can hang indefinitely
3. **No graceful shutdown** - MinIO client not explicitly closed (minor)

## Impact

- **Job failures**: Resume/reference letter processing fails if storage has momentary issues
- **Stuck workers**: Large file operations can hang indefinitely, blocking job workers
- **Poor UX**: Users see random upload failures that would succeed on retry

## Current Code

```go
// minio.go - no retry, no timeout
func (s *MinIOStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
    _, err := s.client.PutObject(ctx, s.bucket, key, reader, size, minio.PutObjectOptions{
        ContentType: contentType,
    })
    return err  // ❌ Fails immediately on any error
}

func (s *MinIOStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
    obj, err := s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
    // ❌ No timeout - can hang forever
    return obj, err
}
```

## Solution

Wrap storage operations with failsafe retry pattern (similar to LLM resilience):

```go
import "github.com/failsafe-go/failsafe-go"
import "github.com/failsafe-go/failsafe-go/retrypolicy"

func NewMinIOStorage(cfg config.MinIOConfig) (*MinIOStorage, error) {
    // ... existing setup ...
    
    retryPolicy := retrypolicy.Builder[any]().
        WithMaxAttempts(3).
        WithBackoff(100*time.Millisecond, 5*time.Second).
        Build()
    
    return &MinIOStorage{
        client:      client,
        retryPolicy: retryPolicy,
        // ...
    }, nil
}

func (s *MinIOStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    return failsafe.Run(func() error {
        _, err := s.client.PutObject(ctx, s.bucket, key, reader, size, minio.PutObjectOptions{
            ContentType: contentType,
        })
        return err
    }, s.retryPolicy)
}
```

## Checklist

- [ ] Add failsafe-go dependency (already used in LLM layer)
- [ ] Create retry policy in MinIOStorage constructor
- [ ] Wrap Upload() with retry and 30s timeout
- [ ] Wrap Download() with retry and 30s timeout
- [ ] Wrap Delete() with retry and 10s timeout
- [ ] Wrap Exists() with retry and 5s timeout
- [ ] Add explicit Close() method and call from main.go shutdown
- [ ] Add tests for retry behavior

## Files to Modify

- `src/backend/internal/infrastructure/storage/minio.go`
- `src/backend/cmd/server/main.go` (shutdown cleanup)
- `src/backend/internal/infrastructure/storage/minio_test.go` (new tests)

## Definition of Done
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review