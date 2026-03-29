---
name: research-github-issue
description: |
  Phase 1 of the Research → Plan → Implement workflow.
  Use when a GitHub Issue number is provided. Maps the codebase, fetches design context,
  captures visual baselines, and outputs a structured research document.
  Keywords: research, github issue, phase 1, codebase discovery, baseline
---

# GitHub Issue Research Protocol

You are in **Phase 1 (Research)**. Your goal is to document the current state of the codebase regarding a specific GitHub Issue. **Do NOT write or modify any application code.**

## Inputs

- **Issue number** (e.g. `#42`)
- **Repository**: `DoGab/VitalStack`

## Steps

### 1. Fetch the Issue

Use the GitHub MCP tools to read the issue from `DoGab/VitalStack`:
- Title, body, labels
- Extract: **Objective**, **Scope**, **Definition of Done (DoD)**, **UI Description**, **Technical Notes**
- Note the `scope` field — it determines which steps below are relevant

### 2. Codebase Discovery

Use file search and grep tools to locate the relevant files:

- **If scope includes "Backend (Go API)":**
  - Read `apps/api-go/architecture.md` for layered architecture context
  - Find relevant controllers, services, models in `apps/api-go/`
  - Trace the data flow from route → controller → service

- **If scope includes "Frontend (SvelteKit)":**
  - Read `apps/web/architecture.md` for VitalStack design language
  - Find relevant pages, components, stores in `apps/web/src/`
  - Check `apps/web/src/lib/config/nutrition-config.ts` if macros are involved
  - List installed shadcn-svelte components: `ls apps/web/src/lib/components/ui/`

### 3. Fetch Design Context (If scope includes Frontend or Design)

- Read `apps/web/architecture.md` for VitalStack design tokens
- Check if existing custom components (`CircularProgress`, `MacroBars`, `SectionHeader`, `StatCard`) are relevant
- Read the `shadcn-svelte-catalog` skill to identify any uninstalled components that could help
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
<Summary of what the issue asks for>

## 📋 Scope
<List: Frontend / Backend / Database / Design>

## ✅ Definition of Done
<Copy the DoD checkboxes from the issue>

## 🗺️ Discovered Files

### Backend (if applicable)
- `apps/api-go/...` — <role of this file>

### Frontend (if applicable)
- `apps/web/src/...` — <role of this file>

## 🎨 Design Context (if applicable)
- Relevant VitalStack tokens: <colors, fonts, themes>
- Recommended components: <shadcn-svelte + custom components>
- Components to install: <any uninstalled shadcn-svelte components>

## 📸 Visual Baseline (if applicable)
![Baseline UI](../assets/baseline_issue_<number>.png)
<1-sentence visual analysis of the current state>

## 🔗 Architecture References
- Backend: `apps/api-go/architecture.md`
- Frontend: `apps/web/architecture.md`
```

### 6. Handoff

End the file with:
> **Research complete. Ready for Planning phase.**

Then tell the user: "Research is saved at `thoughts/research/research_issue_<number>.md`. Start a new conversation and ask me to plan this issue."

## Context Size Notes

- This skill runs as its own conversation — do NOT carry over into planning
- Keep the research file concise — summarize, don't paste entire file contents
- The research file is the ONLY artifact passed to the next phase
