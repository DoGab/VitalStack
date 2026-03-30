# Plan: Editable Meal Log

> Source: `thoughts/research/research_editable-meal-log.md`

## 🎯 Objective
Improve the `MealDetailsModal` to allow editing the timestamp and ingredient amounts of a logged meal, as well as adding new ingredients. The UI should switch to an edit mode, where the "Edit" button becomes a "Save" button. Clicking "Save" persists the changes to the backend via a new API endpoint.

## ✅ Definition of Done

- [ ] Added backend endpoint `PUT /api/nutrition/log/{id}` to update standard meal fields and totally replace ingredients.
- [ ] Upgraded `MealDetailsModal.svelte` with an edit mode state.
- [ ] Users can edit the logged meal timestamp.
- [ ] Users can edit amounts for specific ingredients inside the modal.
- [ ] Users can add more ingredients to an existing log.
- [ ] The "Edit" button transforms into a "Save" button and persists data.
- [ ] Tests pass for all new code.

## 📝 Implementation Steps

### Step 1: Backend Repository
**Files:** `apps/api-go/internal/repository/food_log.go`, `apps/api-go/internal/repository/food_log_test.go`
**Action:** Implement `UpdateFoodLogWithIngredients` wrapping updating the `logged_foods` record and fully replacing linked components in `logged_food_components` inside a database transaction to ensure atomicity. Create test coverage for the repository method.
**Success:** The method exists, successfully uses a transaction, and passes unit tests.
- [ ] Done

### Step 2: Backend Service Layer
**Files:** `apps/api-go/pkg/service/nutrition_types.go`, `apps/api-go/pkg/service/nutrition_service.go`, `apps/api-go/pkg/service/nutrition_service_test.go`
**Action:** Define `UpdateFoodLogInput` service DTO. Implement `UpdateLoggedFood` function in `NutritionService` to orchestrate validation and call the repository layer. Write unit tests for the logic.
**Success:** The service layer method successfully processes the DTO and invokes the repository, with tests passing.
- [ ] Done

### Step 3: Backend Controller Layer
**Files:** `apps/api-go/internal/controller/nutrition_types.go`, `apps/api-go/internal/controller/nutrition_controller.go`, `apps/api-go/internal/controller/nutrition_controller_test.go`
**Action:** 
1. Define typed Huma structs: `UpdateLogInput` (URI), `UpdateLogInputBody` (Body), and `UpdateLogOutput`. Add `doc`, `json`, and `example` tags. Add `ToServiceDTO` converter. 
2. Write `UpdateLogHandler` and register `PUT /api/nutrition/log/{id}` setup via Huma. 
3. Write controller tests verifying HTTP boundaries.
**Success:** The API endpoint handles requests and responds gracefully, and tests pass.
- [ ] Done

### Step 4: OpenAPI Sync & Backend Verification
**Files:** `apps/api-go/api/openapi.yaml` (auto), `apps/web/src/lib/api/client/index.ts` (auto)
**Action:** Run `golangci-lint run` in the backend folder to ensure code quality. Next, run `make openapi` at the repository root to regenerate the OpenAPI spec and the TypeScript client so the frontend can access the new endpoint.
**Success:** No lint errors. The TypeScript client updates to include the `PUT` operation logic.
- [ ] Done

### Step 5: Frontend Interface Update
**Files:** `apps/web/src/lib/components/dashboard/MealDetailsModal.svelte`
**Action:** Use Svelte 5 `$state` runes to introduce `isEditing` to the modal. Expose a date/time picker (e.g. `input type="datetime-local"`) for the timestamp. Use the existing editable ingredient pattern for the list to modify amounts and toggle states. Bind a "Save" action to invoke the new API via the TypeScript client, refetching/revalidating page data upon success.
**Success:** The user can toggle into an edit mode, alter inputs, click "Save", and successfully persist the mutation to the backend.
- [ ] Done

### Step 6: Verify & Clean Up
**Action:** Run frontend linting/formatting (`pnpm check`, `pnpm format`). Start `make dev`, then visually verify the modal on desktop & mobile resolutions (either manually or using the backend dev tools). Provide screenshots confirming visual adherence and test end-to-end functionality.
**Success:** All checks pass, layout scales effectively across device footprints, and logging edits reflect on the history correctly.
- [ ] Done
