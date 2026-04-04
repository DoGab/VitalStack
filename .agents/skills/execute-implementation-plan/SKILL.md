---
name: execute-implementation-plan
description: |
  Phase 3 of the Research → Plan → Implement workflow.
  Reads an approved plan from thoughts/plan/ and executes it step-by-step,
  writing code, running tests, and verifying at each step.
  Keywords: implement, execute, phase 3, code, build
---

# Implement Plan

You are in **Phase 3 (Implementation)**. You are tasked with implementing an approved technical plan from `thoughts/plan/`. These plans contain phases with specific changes and success criteria.

## Inputs

- An approved plan file from `thoughts/plan/` (the user will specify which one)

## Steps

### 1. Ingest Plan

When given a plan path:

- Read the plan completely and check for any existing checkmarks (- [x])
- Read the original ticket and all files mentioned in the plan
- Read files fully - never use limit/offset parameters, you need complete context
- Think deeply about how the pieces fit together
- Create a todo list to track your progress
- Start implementing if you understand what needs to be done

If no plan path provided, ask for one.

### 2. Execute Phase-by-Phase

**Do not rush.** You must execute the plan **Phase by Phase**. Do not attempt to complete the entire plan in a single prompt response. Plans are carefully designed, but reality can be messy. Your job is to:

For each phase in the plan:

1. **Announce Intent**: Tell the user which phase you are starting.
2. **Execute**: Make the changes specified in the phase (write code, create files). Do NOT make changes outside the scope of the current phase.
3. **Automated Verification**: Run the automated checks defined in the phase (e.g., `make lint`, `go test`, `npm run check`).
  - Fix any issues before proceeding

4. **Update the Plan**: Once automated verification passes, edit the `thoughts/plan/plan_<name>.md` file to check off `[x]` the Automated Verification tasks for that phase. Also update the tasks in the tasks list `tasks.md`.
5. **Pause for Manual Verification**: If the phase includes "Manual Verification" criteria (like testing a UI flow in the browser), you **MUST STOP**.
   - After completing all automated verification for a phase, pause and inform the human that the phase is ready for manual testing. Use this format:
    ```
    Phase [N] Complete - Ready for Manual Verification

    Automated verification passed:
    - [List automated checks that passed]

    Please perform the manual verification steps listed in the plan:
    - [List manual verification items from the plan]

    Let me know when manual testing is complete so I can proceed to Phase [N+1].
    ```
6. **Complete Phase**: Upon human approval, check off `[x]` the Manual Verification tasks in the plan document.

### 3. Handling Deviations

If a phase cannot be implemented as planned (e.g., unforeseen constraints or a complex bug):
- **STOP and think deeply about why the plan can't be followed**
- Present the issue clearly:
  ```
  Issue in Phase [N]:
  Expected: [what the plan says]
  Found: [actual situation]
  Why this matters: [explanation]

  How should I proceed?
  ```
- **Do NOT proceed further** until the user approves the adjusted plan.

### 4. If You Get Stuck
When something isn't working as expected:

- First, make sure you've read and understood all the relevant code
- Consider if the codebase has evolved since the plan was written
- Present the mismatch clearly and ask for guidance
- Use sub-tasks sparingly - mainly for targeted debugging or exploring unfamiliar territory.

### 5. Resuming Work
If the plan has existing checkmarks:

- Trust that completed work is done
- Pick up from the first unchecked item
- Verify previous work only if something seems off
- Remember: You're implementing a solution, not just checking boxes. Keep the end goal in mind and maintain forward momentum.

### 6. Final Verification

After all phases are complete:

1. **Backend (if applicable):** Run `golangci-lint` and all tests.
2. **Frontend (if applicable):** Run the `/verify-frontend` workflow to capture final screenshots.
3. **OpenAPI:** If any controllers changed, confirm `make openapi` was run.
4. **Architecture docs:** If architectural changes were made, confirm the relevant `architecture.md` was updated.
5. **DoD Check:** Go through the Definition of Done checkboxes and verify each one.

### 7. Report

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
