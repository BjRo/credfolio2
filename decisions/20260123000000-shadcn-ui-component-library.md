# shadcn/ui as the UI Component Library

## Status

Accepted

## Context

The project needs a consistent approach to building UI components for the React/Next.js frontend. Options considered:

1. **Build components from scratch** - Maximum flexibility but time-consuming
2. **Traditional component libraries** (e.g., Material UI, Chakra UI) - Quick start but locked into their styling paradigms
3. **shadcn/ui** - Copy-paste component collection, owns the code, built on Radix UI primitives

## Decision

Use **shadcn/ui** as the UI component library for the frontend.

### Key characteristics of shadcn/ui:

- **Not a dependency** - Components are copied into your codebase, not installed as npm packages
- **Built on Radix UI** - Accessible primitives with proper ARIA attributes and keyboard navigation
- **Tailwind CSS native** - Uses Tailwind utility classes, integrates with our Tailwind CSS 4 setup
- **Customizable** - Since you own the code, components can be modified without forking
- **TypeScript-first** - Full type safety out of the box
- **New York style variant** - Using the more refined "new-york" style over "default"

### Configuration

- **components.json** - Defines component paths and styling preferences
- **CSS variables** - Theme colors defined in `globals.css` using OKLCH color space
- **Dark mode** - Configured via `.dark` class variant
- **Path aliases** - Components at `@/components/ui`, utilities at `@/lib/utils`

## Consequences

### Positive

- Full ownership of component code - can customize without constraints
- Consistent accessibility through Radix UI primitives
- Seamless Tailwind CSS 4 integration
- Active community and regular updates to copy from
- No runtime CSS-in-JS overhead

### Negative

- Must manually add components as needed (not all available by default)
- Updates require manual re-copying (though usually stable)
- Network restrictions in devcontainer prevent automatic CLI fetching of components

### Adding Components

Due to network restrictions, components must be manually added. Reference the [shadcn/ui documentation](https://ui.shadcn.com/) and copy component code into `/src/frontend/src/components/ui/`.

## Related

- Bean: credfolio2-iev6
- shadcn/ui docs: https://ui.shadcn.com/
- Radix UI: https://www.radix-ui.com/
