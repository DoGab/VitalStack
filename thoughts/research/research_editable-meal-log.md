# Research: Editable Meal Log

## 🎯 Objective
Improve the `MealDetailsModal` to allow editing the timestamp and ingredient amounts of a logged meal, as well as adding new ingredients. The UI should switch to an edit mode, where the "Edit" button becomes a "Save" button. Clicking "Save" persists the changes to the backend via a new API endpoint.

## 🃋. Scope
- Frontend (SvelteKit)
- Backend (Go API)
- Database (Supabase)
- Design

## ✅ Definition of Done

- [ ] Added backend endpoint `PUT /api/nutrition/log/{id}` to update standard meal fields and totally replace ingredients.
- [ ] Upgraded `MealDetailsModal.svelte` with an edit mode state.
- [ ] Users can edit the logged meal timestamp.
- [ ] Users can edit amounts for specific ingredients inside the modal.
- [ ] Users can add more ingredients to an existing log.
- [ ] The "Edit" button transforms into a "Save" button and persists data.
- [ ] Tests pass for all new code.

## 🙺 Target Files

### Existing Files to Modify
- `apps/api-go/internal/controller/nutrition_controller.go` — Add `UpdateLogHandler` and register `PUT /api/nutrition/log/{id}` in Huma.
- `apps/api-go/internal/controller/nutrition_types.go` — Define `UpdateLogInput`, `UpdateLogOutput`, `UpdateLogInputBody`.
- `apps/api-go/pkg/service/nutrition_service.go` — Add `UpdateLoggedFood` function orchestrating the data mutation.
- `apps/api-go/pkg/service/nutrition_types.go` — Define internal `UpdateFoodLogInput` service type.
- `apps/api-go/internal/repository/food_log.go` — Add `UpdateFoodLogWithIngredients` database interaction.
- `apps/web/src/lib/components/dashboard/MealDetailsModal.svelte` — Implement edit state, bind interactive inputs for timestamp and ingredient amounts, and add "Save" HTTP call.

### New Files to Create
- Database migration to create the `update_food_log_atomic` RPC snippet if needed (or equivalent sequence logic in the repository).

## 🎨 Design Context (if applicable)
- Relevant VitalStack tokens: Base-300 for inputs, Primary (#1B3022 / #C5A059) for "Save" actions.
- Recommended components: `Button`, `Input`, `Label` from `shadcn-svelte`. 
- For the datetime, native `input type="time"` or "datetime-local"` will be used to keep consistency and optimal mobile PWA experience.

## 👸️ Visual Baseline (if applicable)
![Baseline UI](../assets/baseline_editable-meal-log.png)
The historical meal view currently displays a read-only list of ingredients and a static timestamp. The new feature will allow swapping this read-only view into an editable form layout seamlessly.

	## 🔥 Architecture References
- Backend: `apps/api-go/architecture.md`
- Frontend: `apps/web/architecture.md`

> **Research complete. Ready for Planning phase.**