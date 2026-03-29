---
name: research-adhoc-feature
description: |
  Phase 1 of the Research → Plan → Implement workflow for ad-hoc features.
  Use when the user requests a feature without a GitHub Issue.
  Maps the codebase, fetches design context, captures visual baselines.
  Keywords: research, ad-hoc, feature, phase 1, codebase discovery
---

# Ad-Hoc Feature Research Protocol

You are in **Phase 1 (Research)**. Your goal is to document how a requested ad-hoc feature integrates into the current codebase. **Do NOT write or modify any application code.**

## Inputs

- **User's feature request** (from the conversation prompt)
- Derive a short, descriptive `<feature_name>` for file naming (e.g., `water-tracking`, `meal-history-export`)

## Steps

### 1. Analyze Request

- Understand the user's prompt and identify the feature's scope
- Determine if it affects: Frontend, Backend, Database, or Design
- Formulate a clear objective statement

### 2. Codebase Discovery

Use file search and grep tools to find where this feature should live:

- **Backend (Go API):**
  - Read `apps/api-go/architecture.md` for layered architecture patterns
  - Find similar existing endpoints/services to match the pattern
  - Identify where new routes, controllers, services should be added

- **Frontend (SvelteKit):**
  - Read `apps/web/architecture.md` for VitalStack design language
  - Find similar existing pages/components to match the pattern
  - Check `apps/web/src/lib/config/nutrition-config.ts` if macros are involved
  - List installed shadcn-svelte components: `ls apps/web/src/lib/components/ui/`

### 3. Fetch Design Context (If UI is involved)

- Read `apps/web/architecture.md` for VitalStack design tokens (Organic Premium / Dark Organic)
- Check if existing custom components (`CircularProgress`, `MacroBars`, `SectionHeader`, `StatCard`) are relevant
- Read the `shadcn-svelte-catalog` skill to identify any uninstalled components that could help
- If Stitch MCP is available, query for relevant design system context

### 4. Visual Baseline (If UI is involved)

- Check if the dev server is running on port 5173
- If not, start it with `make dev`
- Navigate to the page where the new feature will be injected
- Take a screenshot and save it to `thoughts/assets/baseline_<feature_name>.png`

### 5. Output Generation

Create a new file at `thoughts/research/research_<feature_name>.md` with this **exact structure**:

```markdown
# Research: <Feature Name>

## 🎯 Objective
<Summary of what the feature should do>

## 📋 Scope
<List: Frontend / Backend / Database / Design>

## ✅ Definition of Done
- [ ] <Derived acceptance criteria>
- [ ] <Based on user's request>
- [ ] Tests pass for all new code

## 🗺️ Target Files

### Existing Files to Modify
- `path/to/file` — <what changes are needed>

### New Files to Create
- `path/to/new/file` — <purpose>

## 🎨 Design Context (if applicable)
- Relevant VitalStack tokens: <colors, fonts, themes>
- Recommended components: <shadcn-svelte + custom components>
- Components to install: <any uninstalled shadcn-svelte components>

## 📸 Visual Baseline (if applicable)
![Baseline UI](../assets/baseline_<feature_name>.png)
<1-sentence visual analysis explaining where the new feature will integrate>

## 🔗 Architecture References
- Backend: `apps/api-go/architecture.md`
- Frontend: `apps/web/architecture.md`
```

### 6. Handoff

End the file with:
> **Research complete. Ready for Planning phase.**

Then tell the user: "Research is saved at `thoughts/research/research_<feature_name>.md`. Start a new conversation and ask me to plan this feature."

## Context Size Notes

- This skill runs as its own conversation — do NOT carry over into planning
- Keep the research file concise — summarize, don't paste entire file contents
- The research file is the ONLY artifact passed to the next phase
