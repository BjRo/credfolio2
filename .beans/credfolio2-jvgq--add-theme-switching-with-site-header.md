---
# credfolio2-jvgq
title: Add theme switching with site header
status: completed
type: feature
priority: normal
created_at: 2026-01-27T16:49:17Z
updated_at: 2026-01-27T17:48:08Z
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
- [x] Update profile page to use theme-aware colors
- [x] Update ProfileHeader to use theme-aware colors
- [x] Update WorkExperienceSection to use theme-aware colors
- [x] Update EducationSection to use theme-aware colors
- [x] Update SkillsSection to use theme-aware colors
- [x] Update CertificationsSection to use theme-aware colors
- [x] Update ProfileSkeleton to use theme-aware colors

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed