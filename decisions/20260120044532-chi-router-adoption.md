# Chi Router Adoption

**Date**: 2026-01-20
**Bean**: credfolio2-85vv

## Context

The backend needed an HTTP router. The standard library's `http.ServeMux` is functional but lacks features like method-based routing, middleware chaining, and URL parameters that are essential for building a REST API.

## Decision

Adopted [chi](https://github.com/go-chi/chi) (v5) as the HTTP router.

```go
r := chi.NewRouter()
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
r.Get("/health", healthHandler.ServeHTTP)
```

## Reasoning

Chi was chosen because:

- **100% compatible with `net/http`**: Handlers are standard `http.HandlerFunc`, no lock-in
- **Lightweight**: No reflection, no external dependencies beyond stdlib
- **Method-based routing**: `r.Get()`, `r.Post()`, etc. instead of manual method checks
- **Built-in middleware**: Logger, Recoverer, RequestID, Timeout, etc.
- **URL parameters**: `r.Get("/users/{id}", ...)` with `chi.URLParam(r, "id")`
- **Route grouping**: `r.Route("/api", func(r chi.Router) { ... })`
- **Widely adopted**: Battle-tested in production, active maintenance

Alternatives considered:
- `http.ServeMux` (stdlib): Lacks method routing and middleware
- Gin: More opinionated, uses custom context, heavier
- Echo: Similar to Gin, custom context
- Gorilla Mux: Good but chi is more idiomatic and lighter

## Consequences

- All routes are defined using chi's method-based routing (`r.Get`, `r.Post`, etc.)
- Middleware is applied via `r.Use()` or route-specific `r.With()`
- URL parameters accessed via `chi.URLParam(r, "name")`
- Handlers remain standard `http.HandlerFunc` - easy to test and portable
- The `middleware` package provides Logger, Recoverer, and RequestID by default
