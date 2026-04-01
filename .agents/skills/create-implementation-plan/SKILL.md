---
name: create-implementation-plan
description: |
  Phase 2 of the Research → Plan → Implement workflow.
  Reads a research file from thoughts/research/ and generates a step-by-step
  implementation plan with atomic, verifiable steps.
  Keywords: plan, implementation, phase 2, checklist, steps
---

# Implementation Planning Protocol

You are in **Phase 2 (Planning)**. Your goal is to write a highly specific, atomic execution plan based on prior research. **Do NOT write application code yet.**

## Inputs

- A research file from `thoughts/research/` (the user will specify which one)

## Steps

### 1. Ingest Research

- Read the specified file from `thoughts/research/` COMPLETELY.
- Extract: Objective, Scope, DoD, Discovered Codebase State, Design Context.

### 2. Read Project Rules

Since this phase doesn't edit application files directly, the project rules won't auto-load. You **must** manually read the relevant rules to inform the plan:

- **If scope includes Backend:** Read `.agent/rules/backend.md`
- **If scope includes Frontend:** Read `.agent/rules/frontend.md`

These rules define testing requirements, architectural constraints, component usage, and verification steps that your plan must account for.

### 3. Analyze and Verify Understanding

- Cross-reference the research findings with the actual codebase reality.
- If you encounter ambiguities, missing details, or open questions, **STOP AND ASK THE USER**.
- **CRITICAL**: Do NOT write the plan with unresolved questions. Every technical decision must be made before finalizing the plan. The implementation plan must be complete and actionable.

### 4. Draft the Plan

Break the required changes into **Phases** (rather than long monolithic steps). Each phase should:
- Group related file changes (e.g., Schema -> API -> UI).
- Define explicitly what we're NOT doing (to prevent scope creep).
- Separate the success criteria into **Automated Verification** (e.g., `make test`, `npm run lint`) and **Manual Verification** (e.g., UI checks on the browser).
- Include a note for the execution agent to pause for manual verification when needed.

### 5. Output

Create a new file at `thoughts/plan/plan_<name>.md` with this **exact structure**:

```markdown
# Plan: <Feature/Issue Name>

> Source: `thoughts/research/research_<name>.md`

## 🎯 Objective
<Copy from research>

## 🚫 What We're NOT Doing
<Explicitly list out-of-scope items to prevent scope creep>

## ✅ Definition of Done
<Copy the DoD checkboxes from research>

## 📝 Implementation Phases

### Phase 1: <Descriptive Name>

**Overview:** <What this phase accomplishes>

**Changes Required:**
- `path/to/file` — <What to do>

**Success Criteria:**

#### Automated Verification:
- [ ] <e.g., Unit tests pass: `go test ./...`>
- [ ] <e.g., Linting passes: `make lint`>

#### Manual Verification:
- [ ] <e.g., UI behaves correctly when clicking X>

**Implementation Note**: After completing this phase and all automated verification passes, pause here for manual confirmation from the human that the manual testing was successful before proceeding to the next phase.

---

### Phase 2: <Descriptive Name>
...

### Phase N: Final Clean Up
**Overview:** Final overall project verification (OpenAPI sync, architectural doc updates, etc.)
```

### 6. Pause for Review

After saving the file, explicitly ask the user:

> "The implementation plan is saved at `thoughts/plan/plan_<name>.md`. Please review it. Should I proceed to implementation?"

**Do NOT continue until the user approves.**

## Context Size Notes

- Read ONLY the research file and the relevant rules — do not re-discover the codebase
- The plan file is the ONLY artifact passed to the Implementation phase
- Keep steps atomic and concise — the implementation agent reads this file fresh
