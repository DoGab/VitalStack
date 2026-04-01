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

- Read the approved plan from `thoughts/plan/`.
- Note the total number of phases and the Definition of Done.
- Confirm with the user that the plan is approved before starting.

### 2. Execute Phase-by-Phase

**Do not rush.** You must execute the plan **Phase by Phase**. Do not attempt to complete the entire plan in a single prompt response.

For each phase in the plan:

1. **Announce Intent**: Tell the user which phase you are starting.
2. **Execute**: Make the changes specified in the phase (write code, create files). Do NOT make changes outside the scope of the current phase.
3. **Automated Verification**: Run the automated checks defined in the phase (e.g., `make lint`, `go test`, `npm run check`).
4. **Update the Plan**: Once automated verification passes, edit the `thoughts/plan/plan_<name>.md` file to check off `[x]` the Automated Verification tasks for that phase.
5. **Pause for Manual Verification**: If the phase includes "Manual Verification" criteria (like testing a UI flow in the browser), you **MUST STOP**.
   - Ask the user to perform the manual verification steps.
   - Example: *"I have completed Phase 1 and automated tests pass. Please verify the UI changes in your browser. Let me know when you approve so I can mark the manual verification as complete and proceed to Phase 2."*
6. **Complete Phase**: Upon human approval, check off `[x]` the Manual Verification tasks in the plan document.

### 3. Handling Deviations

If a phase cannot be implemented as planned (e.g., unforeseen constraints or a complex bug):
- **Stop executing.**
- Explain the blocker to the user.
- Propose an adjustment to the plan.
- **Do NOT proceed further** until the user approves the adjusted plan.

### 4. Final Verification

After all phases are complete:

1. **Backend (if applicable):** Run `golangci-lint` and all tests.
2. **Frontend (if applicable):** Run the `/verify-frontend` workflow to capture final screenshots.
3. **OpenAPI:** If any controllers changed, confirm `make openapi` was run.
4. **Architecture docs:** If architectural changes were made, confirm the relevant `architecture.md` was updated.
5. **DoD Check:** Go through the Definition of Done checkboxes and verify each one.

### 5. Report

When all phases are complete, report to the user:

> "Implementation complete. All phases executed successfully."
> - ✅ Tests: [pass/fail]
> - ✅ Lint: [pass/fail]
> - ✅ Visual: [verified with screenshots / human]
> - ✅ DoD: [all items checked]

## Context Size Notes

- Read ONLY the plan file — do not re-read the research file
- If you need more context about a file, read it directly (don't re-do research)
- Update the plan file as you go — this is your persistent state
- If the context grows too large, tell the user to start a new conversation and resume from the plan file (unchecked steps will show remaining work)
