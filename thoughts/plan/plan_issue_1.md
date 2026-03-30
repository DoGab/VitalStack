# Plan: Update food logging

> Source: `thoughts/research/research_issue_1.md`

## 🎯 Objective
Redesign the food logging and meal details experience to give users more control over their entries. Key additions include allowing users to explicitly select a timestamp for logs, adjust ingredient amounts before saving, and deselect incorrectly detected ingredients via a toggle. The ingredients list needs a cleaner, collapsible UI pattern maximizing screen real estate. The same design language must be applied to viewing previously logged meals.

## ✅ Definition of Done
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

## 📝 Implementation Steps

### Step 1: Stitch UI Design
**Files:** Stitch MCP Workspace
**Action:** Create a new Stitch project and design the Add Entry / Meal Details screens to visualize the datetime picker and the new list-item based ingredients accordion with toggles and editable amounts.
**Success:** The user approves the UI pattern generated in the Stitch project.
- [ ] Done

### Step 2: Install shadcn-svelte Components
**Files:** `apps/web/`
**Action:** Run `pnpm dlx shadcn-svelte@latest add checkbox accordion` to install necessary components.
**Success:** Components are successfully generated in `apps/web/src/lib/components/ui`.
- [ ] Done

### Step 3: Create IngredientEditor Component
**Files:** `apps/web/src/lib/components/dashboard/IngredientEditor.svelte`
**Action:** Create a reusable component using Svelte 5 runes (`$state` for toggle and quantity) rendering full-width list items. Each item includes a checkbox, editable number input, and macro breakdown string.
**Success:** Component conforms to Svelte 5 syntax and handles local state updates without compiler errors.
- [ ] Done

### Step 4: Update AddEntryModal Component
**Files:** `apps/web/src/lib/components/navigation/AddEntryModal.svelte`
**Action:** Add a Date/Time picker (defaulting to current time). Wrap the `IngredientEditor` component inside the new shadcn `Accordion`. Bind the state from `IngredientEditor` to filter out unchecked ingredients before submitting the API payload. Ensure header layout matches formatting (title left, disabled plus icon right).
**Success:** Unchecked ingredients are safely ignored in logs, updated amounts correctly scale macros, and the component renders properly.
- [ ] Done

### Step 5: Update MealDetailsModal Component
**Files:** `apps/web/src/lib/components/dashboard/MealDetailsModal.svelte`
**Action:** Swap the existing floating card list of ingredients with the new `IngredientEditor` component wrapped in an `Accordion`, mirroring the `AddEntryModal` design exactly. Maintain `MacroBars` usage.
**Success:** Viewing historical meals utilizes visually identical, space-optimized list items inside a collapsible pane.
- [ ] Done

### Step 6: Verify & Clean Up
**Action:** Run the `/verify-frontend` workflow to execute linting, type-checking, and take visual snapshots (desktop and mobile) of the app. Ensure `apps/web/architecture.md` is updated if component usage requires documentation.
**Success:** All tests pass, no lint errors, and screenshots properly reflect the new requirements.
- [ ] Done
