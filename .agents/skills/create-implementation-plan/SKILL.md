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

- Read the specified file from `thoughts/research/`
- Extract: Objective, Scope, DoD, Discovered Files, Design Context

### 2. Read Project Rules

Since this phase doesn't edit application files directly, the project rules won't auto-load. You **must** manually read the relevant rules to inform the plan:

- **If scope includes Backend:** Read `.agent/rules/backend.md`
- **If scope includes Frontend:** Read `.agent/rules/frontend.md`

These rules define testing requirements, architectural constraints, component usage, and verification steps that your plan must account for.

### 3. Draft the Plan

Break the required changes into a **numbered checklist** of atomic steps. Each step should:
- Modify at most 1-2 files
- Be independently verifiable
- Include the exact file path(s) being changed
- Have a clear success criterion
- Comply with the rules read in Step 2

Ensure the plan includes steps for all mandatory activities from the rules (e.g., tests, linting, OpenAPI sync, visual verification, architecture docs updates).

### 4. Output

Create a new file at `thoughts/plan/plan_<name>.md` with this **exact structure**:

```markdown
# Plan: <Feature/Issue Name>

> Source: `thoughts/research/research_<name>.md`

## 🎯 Objective
<Copy from research>

## ✅ Definition of Done
<Copy the DoD checkboxes from research>

## 📝 Implementation Steps

### Step 1: <Short description>
**Files:** `path/to/file`
**Action:** <What to do>
**Success:** <How to verify this step worked>
- [ ] Done

### Step 2: <Short description>
**Files:** `path/to/file`
**Action:** <What to do>
**Success:** <How to verify this step worked>
- [ ] Done

... (continue for all steps)

### Step N: Verify & Clean Up
**Action:** Run linting, tests, and visual verification
**Success:** All tests pass, no lint errors, screenshots look correct
- [ ] Done
```

### 5. Pause for Review

After saving the file, explicitly ask the user:

> "The implementation plan is saved at `thoughts/plan/plan_<name>.md`. Please review it. Should I proceed to implementation?"

**Do NOT continue until the user approves.**

## Context Size Notes

- Read ONLY the research file and the relevant rules — do not re-discover the codebase
- The plan file is the ONLY artifact passed to the Implementation phase
- Keep steps atomic and concise — the implementation agent reads this file fresh
