---
name: research-adhoc-feature
description: |
  Phase 1 of the Research → Plan → Implement workflow for ad-hoc features.
  Use when the user requests a feature without a GitHub Issue.
  Maps the codebase, fetches design context, captures visual baselines.
  Keywords: research, ad-hoc, feature, phase 1, codebase discovery
---

# Ad-Hoc Feature Research Protocol

You are in **Phase 1 (Research)**. Your primary role is to act as a **documentarian**. Your goal is to map the loosely defined ad-hoc feature request to the current state of the codebase. **Do NOT write or modify any application code, and do NOT suggest improvements, evaluate architecture, or prescribe solutions.** Map the "what is" cleanly to inform the planning phase.

## Inputs

- **User's feature request** (from the conversation prompt)
- Derive a short, descriptive `<feature_name>` for file naming (e.g., `water-tracking`, `meal-history-export`)

## Steps

### 1. Analyze Request

- Understand the user's unstructured prompt and identify the core intent to formulate an **Objective**.
- Based on your knowledge of the stack, infer the **Scope** (Frontend, Backend, Database, Design).
- Derive a specific **Definition of Done (DoD)** based on what the user wants to achieve.

### 2. Codebase Discovery

Act as an investigative documentarian using file search and grep tools to find where this feature would live and how existing patterns are structured:

- **Backend (Go API):**
  - Read `apps/api-go/architecture.md` for layered architecture patterns.
  - Find similar existing endpoints/services to match the pattern. Document exactly how they are structured right now.
  - Identify where new routes, controllers, services would be added.

- **Frontend (SvelteKit):**
  - Read `apps/web/architecture.md` for VitalStack design language.
  - Find similar existing pages/components and document their current implementation.
  - Check `apps/web/src/lib/config/nutrition-config.ts` if macros are involved.
  - List installed shadcn-svelte components: `ls apps/web/src/lib/components/ui/`.

### 3. Fetch Design Context (If UI is involved)

- Read `apps/web/architecture.md` for VitalStack design tokens (Organic Premium / Dark Organic)
- Check if existing custom components (`CircularProgress`, `MacroBars`, `SectionHeader`, `StatCard`) are relevant
- Read the `shadcn-svelte-catalog` skill to identify any uninstalled components that could help natively.
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
<Synthesized objective based on the ad-hoc request>

## 📋 Scope
<List: Frontend / Backend / Database / Design>

## ✅ Derived Definition of Done
- [ ] <Derived testable acceptance criteria based on user's request>
- [ ] <More criteria>

## 🗺️ Discovered Codebase State

### Existing Patterns
- `path/to/similar/file:line` — <Detailed, neutral description of how a similar pattern functions currently>

### Integration Points
- `path/to/target/file:line` — <Description of strictly where the new feature will integrate>

## 🎨 Design Context (if applicable)
- Relevant VitalStack tokens: <colors, fonts, themes>
- Recommended components: <shadcn-svelte + custom components>
- Components to install: <any uninstalled shadcn-svelte components>

## 📸 Visual Baseline (if applicable)
![Baseline UI](../assets/baseline_<feature_name>.png)
<1-sentence neutral visual description explaining where the new feature will integrate>

## 🔗 Architecture References
- Backend: `apps/api-go/architecture.md`
- Frontend: `apps/web/architecture.md`
```

### 6. Handoff

End the file with:
> **Research complete. Ready for Planning phase.**

Then tell the user: "Research is saved at `thoughts/research/research_<feature_name>.md`. Ask me to plan this feature."

## Context Size Notes

- This skill runs as its own conversation — do NOT carry over into planning
- Keep the research file concise but precise — include specific file:line references.
- The research file is the ONLY artifact passed to the next phase
