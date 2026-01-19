---
name: tdd
description: Develop features using Test-Driven Development (Red-Green-Refactor). Use when implementing new features, fixing bugs, or when the user mentions TDD, tests first, or test-driven.
---

# Test-Driven Development (TDD)

Follow the classic Red-Green-Refactor cycle for all feature development.

## The TDD Cycle

### 1. RED: Write a Failing Test First

Before writing any production code:

1. **Understand the requirement** - What behavior are we adding?
2. **Write the smallest test** that demonstrates the missing behavior
3. **Run the test** - Confirm it fails for the right reason
4. **Never skip this step** - The failing test proves the test is valid

```
# Frontend (Vitest)
pnpm --filter frontend test:watch

# Backend (Go)
cd src/backend && go test ./...
```

### 2. GREEN: Make the Test Pass

Write the **minimum code** necessary to make the test pass:

1. **Don't over-engineer** - Solve only what the test requires
2. **It's okay to be ugly** - Clean code comes in refactor phase
3. **Run the test** - Confirm it passes
4. **Run all tests** - Ensure no regressions

### 3. REFACTOR: Improve the Code

Now that tests are green, improve the implementation:

1. **Remove duplication** - DRY up code if needed
2. **Improve naming** - Make intent clear
3. **Simplify** - Remove unnecessary complexity
4. **Keep tests green** - Run tests after each change

## Key Principles

### Test Behavior, Not Implementation

```typescript
// BAD: Testing implementation details
expect(component.state.isLoading).toBe(true)

// GOOD: Testing observable behavior
expect(screen.getByRole('progressbar')).toBeInTheDocument()
```

### One Assertion Per Test (When Practical)

Each test should verify one logical concept. Multiple assertions are fine if they verify the same behavior.

### Descriptive Test Names

```typescript
// BAD
it('works', () => { ... })

// GOOD
it('displays error message when login fails with invalid credentials', () => { ... })
```

### Test Structure: Arrange-Act-Assert

```typescript
it('calculates total with discount applied', () => {
  // Arrange
  const cart = createCart([{ price: 100, qty: 2 }])
  const discount = { percent: 10 }

  // Act
  const total = cart.calculateTotal(discount)

  // Assert
  expect(total).toBe(180)
})
```

## Project-Specific Testing

### Frontend (TypeScript/React)

- **Framework**: Vitest + Testing Library
- **Location**: Co-locate tests as `*.test.tsx` or `*.test.ts`
- **Run**: `pnpm --filter frontend test` or `test:watch` for TDD

```typescript
import { render, screen } from '@testing-library/react'
import { MyComponent } from './MyComponent'

describe('MyComponent', () => {
  it('renders greeting with provided name', () => {
    render(<MyComponent name="Alice" />)
    expect(screen.getByText('Hello, Alice!')).toBeInTheDocument()
  })
})
```

### Backend (Go)

- **Framework**: Go standard testing
- **Location**: `*_test.go` files alongside source
- **Run**: `pnpm --filter backend test` or `go test ./...`

```go
func TestGreet(t *testing.T) {
    got := Greet("Alice")
    want := "Hello, Alice!"

    if got != want {
        t.Errorf("Greet() = %q, want %q", got, want)
    }
}
```

## TDD Workflow Checklist

When implementing a feature:

- [ ] Understand the requirement clearly
- [ ] Write a failing test that describes the behavior
- [ ] Verify the test fails (RED)
- [ ] Write minimum code to pass the test
- [ ] Verify the test passes (GREEN)
- [ ] Refactor while keeping tests green
- [ ] Repeat for next behavior

## Common Mistakes to Avoid

1. **Writing tests after code** - Defeats the purpose of TDD
2. **Writing too much test at once** - Take small steps
3. **Skipping the red phase** - Always see the test fail first
4. **Not running all tests** - Catch regressions early
5. **Testing private methods** - Test public behavior instead
6. **Mocking too much** - Prefer real objects when practical
