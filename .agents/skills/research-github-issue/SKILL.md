---
name: research-github-issue
description: |
  Phase 1 of the Research → Plan → Implement workflow.
  Use when a GitHub Issue number is provided. Maps the codebase, fetches design context,
  captures visual baselines, and outputs a structured research document.
  Keywords: research, github issue, phase 1, codebase discovery, baseline
---

# GitHub Issue Research Protocol

You are in **Phase 1 (Research)**. Your primary role is to act as a **documentarian**. Your goal is to map the loosely defined GitHub Issue to the current state of the codebase. **Do NOT write or modify any application code, and do NOT suggest improvements, evaluate architecture, or prescribe solutions.** Map the "what is" cleanly to inform the planning phase.

## Inputs

- **Issue number** (e.g. `#42`)
- **Repository**: `DoGab/VitalStack`

## Steps

### 1. Fetch & Analyze the Issue

Use the GitHub MCP tools to read the issue from `DoGab/VitalStack`:
- Read the body of the issue and understand the unstructured feature description.
- Extract the core intent to formulate an **Objective**.
- Based on your knowledge of the stack, infer the **Scope** (Frontend, Backend, Database, Design).
- Derive a specific **Definition of Done (DoD)** based on what the user wants to achieve.

### 2. Codebase Discovery

Act as an investigative documentarian using file search and grep tools to locate relevant files:

- **If scope includes "Backend (Go API)":**
  - Find relevant controllers, services, router registrations, and models in `apps/api-go/`.
  - Trace the data flow from route → controller → service. Document exactly how they are structured right now.
  - Read `apps/api-go/architecture.md` for context on layered architecture patterns.

- **If scope includes "Frontend (SvelteKit)":**
  - Find relevant pages (`+page.svelte`), components (`.svelte`), and stores.
  - Document their current implementation and state.
  - Read `apps/web/architecture.md` for VitalStack design language context.
  - Check `apps/web/src/lib/config/nutrition-config.ts` if macros are involved.
  - List installed shadcn-svelte components: `ls apps/web/src/lib/components/ui/`.

### 3. Fetch Design Context (If scope includes Frontend or Design)

- Read `apps/web/architecture.md` for VitalStack design tokens
- Check if existing custom components (`CircularProgress`, `MacroBars`, `SectionHeader`, `StatCard`) are relevant
- Read the `shadcn-svelte-catalog` skill to identify any uninstalled components that could help natively.
- If Stitch MCP is available, query for relevant design system context

### 4. Visual Baseline (If scope includes Frontend)

- Check if the dev server is running on port 5173
- If not, start it with `make dev`
- Navigate to the relevant page in the browser
- Take a screenshot and save it to `thoughts/assets/baseline_issue_<number>.png`

### 5. Output Generation

Create a new file at `thoughts/research/research_issue_<number>.md` with this **exact structure**:

```markdown
# Research: Issue #<number> — <title>

## 🎯 Objective
<Synthesized objective based on the issue description>

## 📋 Scope
<List: Frontend / Backend / Database / Design>

## ✅ Derived Definition of Done
- [ ] <Derived testable acceptance criteria based on issue description>
- [ ] <More criteria>

## 🗺️ Discovered Codebase State

### Backend (if applicable)
- `apps/api-go/...:line` — <Detailed, neutral description of how this file functions currently>

### Frontend (if applicable)
- `apps/web/src/...:line` — <Detailed, neutral description of how this component functions currently>

## 🎨 Design Context (if applicable)
- Relevant VitalStack tokens: <colors, fonts, themes>
- Recommended components: <shadcn-svelte + custom components>
- Components to install: <any uninstalled shadcn-svelte components>

## 📸 Visual Baseline (if applicable)
![Baseline UI](../assets/baseline_issue_<number>.png)
<1-sentence neutral visual description of the current state>

## 🔗 Architecture References
- Backend: `apps/api-go/architecture.md`
- Frontend: `apps/web/architecture.md`
```

### 6. Handoff

End the file with:
> **Research complete. Ready for Planning phase.**

Then tell the user: "Research is saved at `thoughts/research/research_issue_<number>.md`. Ask me to plan this issue."

## Context Size Notes

- This skill runs as its own conversation — do NOT carry over into planning
- Keep the research file concise but precise — include specific file:line references.
- The research file is the ONLY artifact passed to the next phase
