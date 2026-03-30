# Research: Issue #1 — Update food logging

## Objective
Redesign the food logging and meal details experience to give users more control over their entries. Key additions include allowing users to explicitly select a timestamp for logs, adjust ingredient amounts before saving, and deselect incorrectly detected ingredients via a toggle. The ingredients list needs a cleaner, collapsible UI pattern maximizing screen real estate. The same design language must be applied to viewing previously logged meals.

## Scope
- Design (Stitch)
- Frontend (SvelteKit)

## Definition of Done
- [ ] UI is redesigned in a new Stitch MCP project and approved.
- [ ] "Add Entry" screen includes a Date/Time picker (defaults to `now()`).
- [ ] Ingredients in the list have editable amounts (e.g., changing 100g to 150g).
- [ ] Ingredients have a selection toggle (default: all checked); unchecked items are excluded from the final API payload.
- [ ] Ingredients list is wrapped in a collapsible shadcn Accordion/container (default: open/expanded).
- [ ] Ingredients list header has the title aligned left, and a disabled `+` icon button aligned right (for future use).
- [ ] Ingredients are rendered as full-width list items rather than separate floating containers. 
- [ ] Each ingredient list item still displays its individual macro breakdown (Calories, Protein, Carbs, Fat) just as it currently does.
- [ ] The existing `MacroBars` component remains untouched and continues to display the totals.
- [ ] The exact same component/design logic is used for viewing a historical logged meal.
- [ ] Visual verification passes via the `/verify-frontend` workflow.

## Discovered Files

### Frontend
- `apps/web/src/lib/components/navigation/AddEntryModal.svelte` — The component handling the logging flow ("Add Entry").
- `apps/web/src/lib/components/dashboard/MealDetailsModal.svelte` — The component used for viewing a historically logged meal.
- `apps/web/src/lib/components/dashboard/TodaysMeals.svelte` — The parent view displaying historical meals that opens `MealDetailsModal`.
- `apps/web/src/lib/components/ui/macro-bars.svelte` — Existing macro distribution component that must be kept exactly as it is.

## Design Context
- Relevant VitalStack tokens: Organic Premium (Light theme) is the baseline. Primary: Deep Arboretum (#1B3022), Secondary: Burnished Gold (#C5A059), Base: Cream Paper (#F9F7F2). Typography uses `Inter` for body UI text and `JetBrains Mono` for macro numbers. Using shadcn-svelte aesthetics (`rounded-md`, clean borders, `shadow-sm`).
- Recommended components: Use Svelte 5 runes (`$state`) for tracking the ingredient selection map and amounts locally. Core logic goes into a highly reusable `IngredientEditor.svelte`.
- Components to install: `checkbox` and `accordion` (Note: `accordion` allows collapsible groups, while `collapsible` is already installed).

## Visual Baseline
![Baseline UI](../assets/baseline_issue_1.png)
The current interface shows ingredients in isolated cards; this change will streamline them into full-width list item rows within a collapsible section.

## Architecture References
- Frontend: `apps/web/architecture.md`

> **Research complete. Ready for Planning phase.**
