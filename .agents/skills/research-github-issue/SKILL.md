---
name: research-github-issue
description: |
  Researches the codebase for a given GitHub Issue.
  Use when a GitHub Issue number is provided. Maps the codebase, fetches design context,
  captures visual baselines, and outputs a structured research document.
  Keywords: research, github issue, codebase discovery, baseline
---

# Research Codebase (GitHub)

You are tasked with conducting comprehensive research across the codebase to answer user questions by spawning parallel sub-agents and synthesizing their findings.

## ⛔ OUTPUT GUARDRAILS — READ FIRST

> **YOUR OUTPUT FILE MUST be written to `thoughts/research/` — NEVER to `thoughts/plan/`.**
> - Output path: `thoughts/research/YYYY-MM-DD-XXXX-description.md`
> - `thoughts/plan/` is for implementation plans. This skill is ONLY for research.
> - If you find yourself writing to `thoughts/plan/`, you are doing the WRONG thing. STOP and correct.
> - The research file MUST be created using `write_to_file` before you present findings to the user.
> - If the file was NOT created, you have NOT completed the skill. Do NOT summarize findings without writing the file first.

## ⛔ SKILL ISOLATION — DO NOT CROSS

> **This skill is RESEARCH ONLY. It is NOT a pipeline. It does NOT hand off to anything.**
> - **DO NOT** read, invoke, or reference the `create-implementation-plan` skill.
> - **DO NOT** read, invoke, or reference the `execute-implementation-plan` skill.
> - **DO NOT** create any implementation plans, task breakdowns, or phase plans.
> - **DO NOT** write code, propose changes, or suggest architecture.
> - **DO NOT** continue into planning after research is complete.
> - After writing the research file, your job is **DONE**. Full stop.
> - The user will manually invoke the planning skill in a **separate session** if and when they choose to.

## CRITICAL: YOUR ONLY JOB IS TO DOCUMENT AND EXPLAIN THE CODEBASE AS IT EXISTS TODAY

- DO NOT suggest improvements or changes unless the user explicitly asks for them
- DO NOT perform root cause analysis unless the user explicitly asks for them
- DO NOT propose future enhancements unless the user explicitly asks for them
- DO NOT critique the implementation or identify problems
- DO NOT recommend refactoring, optimization, or architectural changes
- ONLY describe what exists, where it exists, how it works, and how components interact
- You are creating a technical map/documentation of the existing system


## Inputs

- **Issue number** (e.g. `#42`)
- **Repository**: `DoGab/VitalStack`

## Initial Setup: Fetch & Analyze the Issue

Use the GitHub MCP tools to read the issue from `DoGab/VitalStack`:
- Read the body of the issue and understand the unstructured feature description.

## Steps to follow after getting the GitHub issue:

1. **Read any directly mentioned files first:**

- If the user mentions specific files (tickets, docs, JSON), read them FULLY first
- **CRITICAL**: Read these files yourself in the main context before spawning any sub-tasks
- This ensures you have full context before decomposing the research

2. **Analyze and decompose the research question:**

- Break down the user's query into composable research areas
- Take time to ultrathink about the underlying patterns, connections, and architectural implications the user might be seeking
- Identify specific components, patterns, or concepts to investigate
- Create a research plan using the `task.md` artifact to track all subtasks
- Consider which directories, files, or architectural patterns are relevant
- Extract the core intent from the issue description to formulate an **Objective**.
- Derive a specific **Definition of Done (DoD)** based on what the user wants to achieve.

3. **Spawn parallel sub-agent tasks for comprehensive research:**

- Create multiple Task agents to research different aspects concurrently
- We now have specialized agents that know how to do specific research tasks:

**For codebase research:**

- Use the **codebase-locator** workflow to find WHERE files and components live
- Use the **codebase-analyzer** workflow to understand HOW specific code works (without critiquing it)
- Use the **codebase-pattern-finder** workflow to find examples of existing patterns (without evaluating them)

**IMPORTANT**: All workflows are documentarians, not critics. They will describe what exists without suggesting improvements or identifying issues.

**For thoughts directory:**

- Use the **thoughts-locator** workflow to discover what documents exist about the topic
- Use the **thoughts-analyzer** workflow to extract key insights from specific documents (only the most relevant ones)

**For web research (only if user explicitly asks):**

- Use the **web-search-researcher** workflow for external documentation and resources
- IF you use web-research agents, instruct them to return LINKS with their findings, and please INCLUDE those links in your final report

The key is to use these workflows intelligently:

- Start with locator workflows to find what exists
- Then use analyzer workflows on the most promising findings to document how they work
- Run multiple workflows in parallel when they're searching for different things
- Each workflow knows its job - just tell it what you're looking for
- Don't write detailed prompts about HOW to search - the workflows already know
- Remind workflows they are documenting, not evaluating or improving

4. **Fetch Design Context (If scope includes Frontend or Design)**

- Read `apps/web/architecture.md` for VitalStack design tokens
- Check if existing custom components (`CircularProgress`, `MacroBars`, `SectionHeader`, `StatCard`) are relevant
- Read the `shadcn-svelte-catalog` skill to identify any uninstalled components that could help natively.
- If Stitch MCP is available, query for relevant design system context
- Check if the dev server is running on https://localhost:5173
- If not, start it with `make dev`
- Navigate to the relevant page in the browser
- Take a screenshot and save it to `thoughts/assets/baseline_issue_<number>.png`


5. **Wait for all sub-agents to complete and synthesize findings:**

- IMPORTANT: Wait for ALL sub-agent tasks to complete before proceeding
- Compile all sub-agent results (both codebase and thoughts findings)
- Prioritize live codebase findings as primary source of truth
- Use thoughts/ findings as supplementary historical context
- Connect findings across different components
- Include specific file paths and line numbers for reference
- Verify all thoughts/ paths are correct (e.g., thoughts/allison/ not thoughts/shared/ for personal files)
- Highlight patterns, connections, and architectural decisions
- Answer the user's specific questions with concrete evidence

6. **Output Generation — MANDATORY (do NOT skip this step)**

> **⛔ PRE-WRITE ASSERTION: Before writing, verify:**
> - [ ] You are writing to `thoughts/research/` (NOT `thoughts/plan/`)
> - [ ] The filename follows the format `YYYY-MM-DD-XXXX-description.md`
> - [ ] You are using `write_to_file` tool (NOT just summarizing in chat)
> - [ ] The content includes all sections from the template below

- **You MUST use the `write_to_file` tool** to create the file at `thoughts/research/YYYY-MM-DD-XXXX-description.md`
- This is NOT optional. The skill is INCOMPLETE without this file.
  - Format: `YYYY-MM-DD-XXXX-description.md` where:
    - YYYY-MM-DD is today's date
    - XXXX is the GitHub Issue number
    - description is a brief kebab-case description of the research topic
  - Example: `2025-04-04-42-parent-child-tracking.md`
  ```markdown
  ---
  date: [Current date and time with timezone in ISO format]
  git_commit: [Current commit hash]
  branch: [Current branch name]
  topic: "[User's Question/Topic]"
  tags: [research, codebase, relevant-component-names]
  status: complete
  last_updated: [Current date in YYYY-MM-DD format]
  ---

  # Research: [User's Question/Topic]

  **Date**: [Current date and time with timezone]
  **Git Commit**: [Current commit hash]
  **Branch**: [Current branch name]
  **GitHub Issue**: #<number> — <title>

  ## Research Question
  [Original user query]

  ## Summary
  [High-level documentation of what was found, answering the user's question by describing what exists]

  ## 🎯 Objective
  [Synthesized objective based on the issue description]

  ## 📋 Scope
  [List: Frontend / Backend / Database / Design]

  ## ✅ Derived Definition of Done
  - [ ] <Derived testable acceptance criteria based on issue description>
  - [ ] <More criteria>

  ## Detailed Findings

  ### [Component/Area 1]
  - Description of what exists ([file.ext:line](link))
  - How it connects to other components
  - Current implementation details (without evaluation)

  ### [Component/Area 2]
  ...

  ## Code References
  - `path/to/file.py:123` - Description of what's there
  - `another/file.ts:45-67` - Description of the code block

  ## 🔗 Architecture Documentation
  [Current patterns, conventions, and design implementations found in the codebase]
  - Backend: `apps/api-go/architecture.md`
  - Frontend: `apps/web/architecture.md`

  ## 🎨 Design Context (if applicable)
  - Relevant VitalStack tokens: <colors, fonts, themes>
  - Recommended components: <shadcn-svelte + custom components>
  - Components to install: <any uninstalled shadcn-svelte components>

  ## 📸 Visual Baseline (if applicable)
  - ![Baseline UI](../assets/baseline_issue_<number>.png)
  - <1-sentence neutral visual description of the current state>

  ## Historical Context (from thoughts/)
  [Relevant insights from thoughts/ directory with references]
  - `thoughts/shared/something.md` - Historical decision about X
  - `thoughts/local/notes.md` - Past exploration of Y
  Note: Paths exclude "searchable/" even if found there

  ## Related Research
  [Links to other research documents in thoughts/shared/research/]

  ## Open Questions
  [Any areas that need further investigation]
  ```

7. **Present findings:**

- Present a concise summary of findings to the user
- Include key file references for easy navigation
- Ask if they have follow-up questions or need clarification

8. **Handle follow-up questions:**

- If the user has follow-up questions, append to the same research document
- Update the frontmatter fields last_updated and last_updated_by to reflect the update
- Add last_updated_note: "Added follow-up research for [brief description]" to frontmatter
- Add a new section: ## Follow-up Research [timestamp]
- Spawn new sub-agents as needed for additional investigation
- Continue updating the document

9. **Post-Write Verification**

- After calling `write_to_file`, verify the file exists by reading it back or listing the directory.
- If the file was NOT created, retry immediately. Do NOT proceed to the handoff step without a confirmed file.

10. **End of Skill — STOP HERE**

End the research file with:
> **Research complete.**

Then tell the user:
> "Research is saved at `thoughts/research/YYYY-MM-DD-XXXX-description.md`.
> When you're ready, start a **new session** and invoke the planning skill manually."

**After delivering this message, you are DONE. Do not continue working. Do not invoke other skills. Do not start planning. Yield control to the user.**

> **⛔ FINAL CHECK**: If you are about to respond to the user but have NOT yet written the research file to `thoughts/research/`, STOP. Go back and write it NOW. The file is the primary deliverable of this skill, not the chat summary.

## Important notes:
- Always use parallel Task agents to maximize efficiency and minimize context usage
- Always run fresh codebase research - never rely solely on existing research documents
- The thoughts/ directory provides historical context to supplement live findings
- Focus on finding concrete file paths and line numbers for developer reference
- Research documents should be self-contained with all necessary context
- Each sub-agent prompt should be specific and focused on read-only documentation operations
- Document cross-component connections and how systems interact
- Include temporal context (when the research was conducted)
- Link to GitHub when possible for permanent references
- Keep the main agent focused on synthesis, not deep file reading
- Have sub-agents document examples and usage patterns as they exist
- Explore all of thoughts/ directory, not just research subdirectory
- **CRITICAL**: You and all sub-agents are documentarians, not evaluators
- **REMEMBER**: Document what IS, not what SHOULD BE
- **NO RECOMMENDATIONS**: Only describe the current state of the codebase
- **File reading**: Always read mentioned files FULLY (no limit/offset) before spawning sub-tasks
- **Critical ordering**: Follow the numbered steps exactly
  - ALWAYS read mentioned files first before spawning sub-tasks (step 1)
  - ALWAYS wait for all sub-agents to complete before synthesizing (step 4)
  - ALWAYS gather metadata before writing the document (step 5 before step 6)
  - NEVER write the research document with placeholder values
- **Path handling**: The thoughts/searchable/ directory contains hard links for searching
  - Always document paths by removing ONLY "searchable/" - preserve all other subdirectories
  - Examples of correct transformations:
    - `thoughts/searchable/allison/old_stuff/notes.md` → `thoughts/allison/old_stuff/notes.md`
    - `thoughts/searchable/shared/prs/123.md` → `thoughts/shared/prs/123.md`
    - `thoughts/searchable/global/shared/templates.md` → `thoughts/global/shared/templates.md`
  - NEVER change allison/ to shared/ or vice versa - preserve the exact directory structure
  - This ensures paths are correct for editing and navigation
- **Frontmatter consistency**:
  - Always include frontmatter at the beginning of research documents
  - Keep frontmatter fields consistent across all research documents
  - Update frontmatter when adding follow-up research
  - Use snake_case for multi-word field names (e.g., `last_updated`, `git_commit`)
  - Tags should be relevant to the research topic and components studied

## Context Size Notes

- This skill runs as its own conversation — do NOT carry over into planning
- The research file is the ONLY output of this skill
- **NEVER** read or invoke other skills (especially `create-implementation-plan` or `execute-implementation-plan`)
- After the research file is written and confirmed, the skill is COMPLETE — yield to the user
