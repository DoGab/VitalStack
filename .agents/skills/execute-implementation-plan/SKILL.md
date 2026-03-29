---
name: execute-implementation-plan
description: |
  Phase 3 of the Research → Plan → Implement workflow.
  Reads an approved plan from thoughts/plan/ and executes it step-by-step,
  writing code, running tests, and verifying at each step.
  Keywords: implement, execute, phase 3, code, build
---

# Implementation Protocol

You are in **Phase 3 (Implementation)**. Your goal is to safely and systematically execute an approved plan.

## Inputs

- An approved plan file from `thoughts/plan/` (the user will specify which one)

## Steps

### 1. Ingest Plan

- Read the approved plan from `thoughts/plan/`
- Note the total number of steps and the Definition of Done
- Confirm with the user that the plan is approved before starting

### 2. Execute Step-by-Step

**Do not rush.** Complete Step 1, verify it works, and only then move to Step 2.

For each step in the plan:

1. **Read** the step's description, files, and success criteria
2. **Execute** the changes (write code, create files, install components)
3. **Verify** the success criteria is met
4. **Update** the plan file — check off the step's `- [ ] Done` → `- [x] Done`

### 3. Follow the Project Rules

When you edit files in `apps/web/` or `apps/api-go/`, the project rules in `.agent/rules/frontend.md` and `.agent/rules/backend.md` will automatically activate via their glob triggers. **Follow all rules that are loaded** — they define testing, architecture, component usage, linting, and documentation requirements.

If rules don't auto-load (e.g., you're working on other files), read the relevant rule file manually.

### 4. Final Verification

After all steps are complete:

1. **Backend (if applicable):** Run `golangci-lint` and all tests
2. **Frontend (if applicable):** Run the `/verify-frontend` workflow:
   - Start full stack with `make dev`
   - Take desktop screenshot (1280px) and mobile screenshot (375px)
   - Compare against the baseline screenshot from `thoughts/assets/`
   - Verify VitalStack design language is correct
3. **OpenAPI:** If any controllers changed, confirm `make openapi` was run
4. **Architecture docs:** If architectural changes were made, confirm the relevant `architecture.md` was updated
5. **DoD Check:** Go through the Definition of Done checkboxes and verify each one

### 5. Report

When all steps are complete, report to the user:

> "Implementation complete. All [N] steps executed successfully."
> - ✅ Tests: [pass/fail]
> - ✅ Lint: [pass/fail]
> - ✅ Visual: [verified with screenshots]
> - ✅ DoD: [all items checked]

## Context Size Notes

- Read ONLY the plan file — do not re-read the research file
- If you need more context about a file, read it directly (don't re-do research)
- Update the plan file as you go — this is your persistent state
- If the context grows too large, tell the user to start a new conversation and resume from the plan file (unchecked steps will show remaining work)
