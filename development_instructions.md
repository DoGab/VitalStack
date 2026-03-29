# VitalStack Development Instructions

This document describes how the AI agent is configured for this project, including the rules, skills, workflows, and the Research → Plan → Implement framework for structured feature development.

---

## Table of Contents

- [Agent Rules](#agent-rules)
- [Agent Skills](#agent-skills)
- [Agent Workflows](#agent-workflows)
- [Research → Plan → Implement Framework](#research--plan--implement-framework)
- [GitHub Issues](#github-issues)
- [Quick Reference](#quick-reference)

---

## Agent Rules

Rules are located in `.agent/rules/` and automatically load based on triggers. They define coding standards the agent must always follow.

| File | Trigger | Description |
|------|---------|-------------|
| `identity.md` | Always on | Shared Lead Architect identity, tech stack proficiency, and cross-stack synchronization directives. |
| `frontend.md` | `apps/web/**` | Frontend rules: custom component catalog, nutrition color system, shadcn-svelte usage, VitalStack design language, Svelte 5 runes, mobile-first UX, and post-implementation verification. |
| `backend.md` | `apps/api-go/**` | Backend rules: mandatory testing, strict layered architecture, DTOs, Huma type safety, documentation comments, linting, and OpenAPI spec sync. |

---

## Agent Skills

Skills are located in `.agents/skills/` and provide specialized instructions the agent follows for specific tasks.

### R/P/I Workflow Skills

| Skill | Description |
|-------|-------------|
| `research-github-issue` | Phase 1: Reads a GitHub Issue, discovers relevant codebase files, fetches design context, captures visual baselines, and outputs a structured research document to `thoughts/research/`. |
| `research-adhoc-feature` | Phase 1: Same as above but for ad-hoc feature requests when no GitHub Issue exists. Derives scope and acceptance criteria from the user's prompt. |
| `create-implementation-plan` | Phase 2: Reads a research file and generates a step-by-step implementation plan with atomic, verifiable steps. Reads the project rules to ensure the plan accounts for all requirements. |
| `execute-implementation-plan` | Phase 3: Reads an approved plan and executes it step-by-step, writing code, running tests, verifying at each step, and checking off completed steps in the plan file. |

### Frontend Skills

| Skill | Description |
|-------|-------------|
| `shadcn-svelte-catalog` | Complete catalog of all 55+ shadcn-svelte components available for installation, organized by category with descriptions and documentation links. |
| `shadcn-svelte-components` | Patterns and best practices for using shadcn-svelte components: Bits UI primitives, trigger snippets, child props, and accessibility. |
| `svelte5-best-practices` | Svelte 5 runes, snippets, SvelteKit data loading, form actions, and TypeScript patterns. |
| `svelte-code-writer` | CLI tools for Svelte 5 documentation lookup and code analysis via the `@sveltejs/mcp` integration. |

### Design Skills

| Skill | Description |
|-------|-------------|
| `enhance-prompt` | Transforms vague UI ideas into polished, specific prompts optimized for Stitch MCP screen generation. |
| `stitch-design` | Unified entry point for Stitch design work: prompt enhancement, design system synthesis, and screen generation/editing via Stitch MCP. |
| `stitch-loop` | Teaches the agent to iteratively build UIs using Stitch with an autonomous baton-passing loop pattern. |

### Other Skills

| Skill | Description |
|-------|-------------|
| `supabase-postgres-best-practices` | Postgres performance optimization: indexes, query patterns, schema design, and RLS policies for Supabase. |
| `find-skills` | Meta-skill for discovering and installing new agent skills from the community. |

---

## Agent Workflows

Workflows are located in `.agent/workflows/` and define repeatable multi-step procedures.

| Workflow | Trigger | Description |
|----------|---------|-------------|
| `verify-frontend` | `/verify-frontend` | Starts the full stack with `make dev`, takes desktop (1280px) and mobile (375px) screenshots, and validates the VitalStack design language. |

---

## Research → Plan → Implement Framework

This framework structures feature development into 3 isolated phases to manage context size and ensure quality.

### Why 3 Phase?

Each phase runs as a **separate conversation**. The `thoughts/` directory acts as persistent file-based memory between phases. This prevents context explosion — each conversation only carries the focused output from the previous phase.

```
thoughts/
├── research/      Phase 1 outputs (codebase analysis, design context)
├── plan/          Phase 2 outputs (step-by-step implementation plans)
└── assets/        Screenshots (baselines, verifications)
```

> `thoughts/` is gitignored — these are ephemeral agent artifacts.

### Phase 1: Research

**Conversation 1** — Understand the problem, don't write code.

```
User: "Research issue #42"
      — or —
User: "Research a new water tracking feature"
```

The agent will:
1. Read the GitHub Issue (or analyze the user's request)
2. Discover relevant files in `apps/api-go/` and `apps/web/`
3. Fetch design context from architecture docs and Stitch MCP (if UI)
4. Take a visual baseline screenshot (if UI)
5. Output: `thoughts/research/research_issue_42.md`

**End this conversation.** Start a new one for Phase 2.

### Phase 2: Plan

**Conversation 2** — Create a step-by-step plan, don't write code.

```
User: "Plan issue #42"
      — or —
User: "Create a plan from thoughts/research/research_water_tracking.md"
```

The agent will:
1. Read the research file
2. Read the relevant project rules (`.agent/rules/backend.md`, `.agent/rules/frontend.md`)
3. Break the work into atomic, verifiable steps
4. Output: `thoughts/plan/plan_issue_42.md`
5. Ask for user approval before proceeding

**Review the plan. Approve or request changes.** Then start a new conversation for Phase 3.

### Phase 3: Implement

**Conversation 3** — Execute the approved plan step-by-step.

```
User: "Implement the plan from thoughts/plan/plan_issue_42.md"
```

The agent will:
1. Read the approved plan
2. Execute each step, writing code and verifying as it goes
3. Check off completed steps in the plan file
4. Run final verification (tests, lint, screenshots, DoD)
5. Report results

If context grows too large, tell the agent to stop. Start a new conversation and resume — the plan file tracks which steps are done.

---

## GitHub Issues

### Issue Template

A structured YAML form template is configured at `.github/ISSUE_TEMPLATE/feature.yml`. It creates issues with these sections:

| Section | Required | Purpose |
|---------|----------|---------|
| **Objective** | ✅ | What the feature should do |
| **Scope** | ✅ | Which stack layers are affected (Frontend, Backend, Database, Design) |
| **Definition of Done** | ✅ | Checkbox acceptance criteria |
| **UI Description** | ❌ | Visual design details, component references |
| **Technical Notes** | ❌ | Implementation hints, API contracts |

Issues are auto-labeled with `enhancement` and `agent-ready`.

### Example Issue

```markdown
### Objective
Users should be able to track their daily water intake. The app should
show how many glasses of water they've consumed vs their daily goal,
and allow adding water in 250ml increments.

### Scope
- Frontend (SvelteKit)
- Backend (Go API)

### Definition of Done
- [ ] API endpoint POST /api/water-intake to record water consumption
- [ ] API endpoint GET /api/water-intake/today to fetch today's intake
- [ ] Frontend shows a water tracking card on the dashboard with
      CircularProgress showing percentage of daily goal
- [ ] Users can tap a "+" button to add 250ml
- [ ] Tests pass for all new backend code
- [ ] OpenAPI spec is regenerated
- [ ] Visual verification via /verify-frontend

### UI Description
Add a water tracking card below the macro section on the dashboard.
Use CircularProgress for the daily goal percentage (blue color).
Use a StatCard to show the total ml consumed.
Follow the Organic Premium theme.

### Technical Notes
Follow the existing NutritionService pattern for the new WaterService.
Store data in a new `water_intake` table in Supabase.
Daily goal defaults to 2000ml (configurable later).
```

---

## Quick Reference

### Key Commands

| Command | Purpose |
|---------|---------|
| `make dev` | Start full stack (Supabase + Go API + SvelteKit) |
| `make openapi` | Regenerate OpenAPI spec + TypeScript client |
| `make lint-api` | Run Go linter |
| `make lint-web` | Run frontend linter |

### Key File Paths

| Path | Purpose |
|------|---------|
| `apps/api-go/architecture.md` | Backend architecture documentation |
| `apps/web/architecture.md` | Frontend architecture documentation |
| `apps/web/src/lib/config/nutrition-config.ts` | Macro color/icon/label configuration |
| `.agent/rules/` | Agent rules (auto-loaded) |
| `.agents/skills/` | Agent skills (invoked on demand) |
| `.agent/workflows/` | Agent workflows (triggered with `/command`) |
| `thoughts/` | R/P/I framework artifacts (gitignored) |
| `.github/ISSUE_TEMPLATE/feature.yml` | GitHub Issue template |
