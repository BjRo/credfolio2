---
# credfolio2-jvgq
title: Add theme switching with site header
status: in-progress
type: feature
created_at: 2026-01-27T16:49:17Z
updated_at: 2026-01-27T16:49:17Z
---

Add light/dark/system theme switching to the app using next-themes, with a new site-wide header navbar.

## Context
- shadcn/ui is already set up with dark mode CSS variables in globals.css
- No header/navbar component exists yet
- Using Next.js 16, React 19, Tailwind CSS 4

## Checklist
- [x] Install next-themes dependency
- [x] Create ThemeProvider wrapper component
- [x] Create ThemeToggle component (Sun/Moon icons, cycles light→dark→system)
- [x] Create SiteHeader component (app name + theme toggle)
- [x] Wire ThemeProvider and SiteHeader into root layout
- [x] Write tests for ThemeToggle
- [x] Write tests for SiteHeader

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed