---
trigger: glob
glob: apps/web/**
---

# Frontend Rules (VitalStack Web)

## Custom Components — Always Use When Applicable

Before creating new UI patterns, check if one of these existing components already handles the use case:

| Component | Import | When to Use |
|-----------|--------|-------------|
| `CircularProgress` | `$lib/components/ui/circular-progress.svelte` | Any circular progress indicator (calorie ring, macro completion, etc.). Supports `percent`, `addedPercent`, and child content. |
| `MacroBars` | `$lib/components/ui/macro-bars.svelte` | Displaying macro nutrient progress bars with current/goal values. Automatically uses `NUTRITION_CONFIG` colors and icons. Supports `added` values for incoming food items. |
| `SectionHeader` | `$lib/components/ui/section-header.svelte` | Page section headings with optional subtitle, action link, or action snippet. |
| `StatCard` | `$lib/components/ui/stat-card.svelte` | Compact stat display cards (e.g., "120 kcal", "32g protein"). Supports `colorValue` for tinted backgrounds. |

## Component Extraction Rule

When a shadcn-svelte component (e.g., `Card`, `Dialog`) is used **multiple times** in a very similar pattern across pages, extract it into a dedicated custom Svelte component in `$lib/components/ui/` with configurable props. Follow the pattern of `StatCard`, `MacroBars`, etc. This keeps pages DRY and ensures consistent usage.

## Nutrition Color System

**Always** use `NUTRITION_CONFIG` from `$lib/config/nutrition-config.ts` for macro-related colors, icons, and labels. Never hardcode nutrition colors.

| Macro | Color Class | Bar Color |
|-------|-------------|-----------|
| Calories | `text-[#D65A31]` (Terracotta) | `#D65A31` |
| Protein | `text-red-500` | `oklch(63.7% 0.237 25.331)` |
| Carbs | `text-amber-500` | `oklch(76.9% 0.188 70.08)` |
| Fat | `text-blue-500` | `oklch(62.3% 0.214 259.815)` |
| Fiber | `text-green-500` | `oklch(72.3% 0.219 149.579)` |

## shadcn-svelte Component Usage

1. **Use Unmodified**: Use shadcn-svelte components as-is without modification. No additional spacing should be added unless exceptionally needed.
2. **Discover Installed**: Run `ls apps/web/src/lib/components/ui/` to see what's already available.
3. **Discover Available**: Read the `shadcn-svelte-catalog` skill for the full catalog of 55+ components available for installation.
4. **Install New**: From the `apps/web` directory: `pnpm dlx shadcn-svelte@latest add <component>`
5. **Namespace Import**: Import with namespace pattern: `import * as Card from '$lib/components/ui/card'`

## Design Language

Follow the **VitalStack Design Language** defined in `apps/web/architecture.md`:

- **Organic Premium** (Light): Primary #1B3022, Secondary #C5A059, Accent #7B506F, Base #F9F7F2
- **Dark Organic** (Dark): Primary #C5A059, Secondary #4E8056, Accent #D4913D, Base #1B3022
- **Classic Font Theme**: Playfair Display (display) + Inter (body) + JetBrains Mono (mono)
- **Modern Font Theme**: Outfit (display) + Inter (body) + JetBrains Mono (mono)
- Theme via CSS variables with `@theme` (without `inline` for multi-theme support)
- Use `data-theme="organic"` / `data-theme="darkorganic"` for theme switching

## Svelte 5 Patterns

- Use Runes exclusively: `$state`, `$derived`, `$props`, `$bindable`
- Use `onclick` not `on:click`; callback props not `createEventDispatcher`
- Use `{#snippet}` / `{@render}` not `<slot>`
- Keep components modular and leverage Svelte transitions

## Mobile-First UX

- **48px minimum** tap targets for thumb-friendly interactions
- Bottom navigation on mobile, sidebar on desktop
- `overscroll-behavior: none` for native app feel
- Tactile feedback with `active:scale-95` on interactive elements
- `user-select: none` on buttons for PWA standards

## Architecture Documentation

Update `apps/web/architecture.md` whenever meaningful architecture changes are made (new components, new patterns, routing changes, state management changes, design system updates).

## Post-Implementation Verification

After implementing frontend changes:
1. Start the full stack with `make dev` (if not already running)
2. Take a **desktop screenshot** (1280px viewport) at `https://localhost:5173`
3. Take a **mobile screenshot** (375px viewport) at `https://localhost:5173`
4. Verify both for layout, colors, fonts, and component rendering
5. Flag any visual issues or confirm success
