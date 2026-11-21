---
slug: 20261121-0035-docs-beta3-refresh
lens: Integration
---

# Feature: Documentation Refresh for 1.0.0-beta.3

## 1. Context & Value
The public documentation (README, CONTRIBUTING and the `docs/` pages in English, Portuguese and Spanish) still describes an older release: it lists a decommissioned agent, wrong file locations, a legacy invocation syntax and omits most of the lifecycle introduced since then (spec sizes, module deltas, independent QA, test-first tasks, auto-update). New adopters follow wrong instructions and experienced users get lost exactly in the non-happy paths (a refused archive, a change after implementation). This feature makes every page match the current release and adds a human-oriented workflow guide with tips and edge cases.

## 2. Out of Scope (Anti-Goals)
- Do not change any product behavior, CLI command, flag, skill blueprint or generated agent guide (including the managed block of `AGENTS.md`).
- Do not document the local kit override (a project-level kit manifest or project-level skills folder); it stays an internal capability until its own proposal is specified.
- Do not add audience labels ("you" vs "agent") to the CLI reference; it remains a plain reference grouped by area.
- Do not create additional documentation pages beyond the workflow guide (no separate concepts or updates page).
- Do not add a link checker, script, CI job or any other versioned verification tooling.
- Do not generate the CLI reference automatically from the command help.
- Do not change the Constitution (`.specforce/docs/`) or module living specs.

## 3. Acceptance Criteria (BDD)

### [US-1] Accurate facts across all pages
**User Story:** AS A developer adopting Specforce, I WANT TO read documentation that matches the release I installed, SO THAT following it never leads me to a missing tool, a wrong folder or an unsupported requirement.

**Scenarios:**
1. **[Happy Path]** GIVEN the current release supports exactly Claude Code, Qwen, Kimi Code, OpenCode, KiloCode, Codex, Antigravity and Cursor WHEN a reader checks any supported-agent list in any language THEN that exact set is listed and no decommissioned agent appears.
2. **[Happy Path]** GIVEN the Constitution lives in the project's `.specforce/docs/` folder with six documents (principles, architecture, ui-ux, security, engineering, governance) plus `modules/` WHEN a reader looks up where the Constitution is THEN every page names that folder and the six documents.
3. **[Happy Path]** GIVEN the build from source requires the toolchain version declared by the project WHEN a reader checks prerequisites in the README and the getting-started page THEN both state the same minimum version.
4. **[Edge Case]** GIVEN a page previously contained duplicated, truncated or misplaced content (a duplicated license sentence, repeated section numbers, a metadata header in the middle of a page, example text in another language) WHEN the refreshed page is read THEN none of these defects remain.
5. **[Edge Case]** GIVEN the root agent guide is linked into a tool folder only for Claude Code WHEN a reader reads the automated configuration section THEN it states that only Claude Code receives that link and that Cursor and the other tools read the root guide or their skills folder directly.

**UI/UX Specifics:**
- **View/Component:** Existing pages edited in place; same headings structure where still valid.
- **Feedback Logic:** Not applicable to static pages; correctness is verified by the scenarios in the roadmap.

**Technical Constraints (NFR):**
- **[Performance]:** Zero occurrences of the decommissioned agent name and of the legacy colon invocation syntax across README, CONTRIBUTING and `docs/` files in all three languages.
- **[Integrity]:** Every fact stated about tools, folders and versions is traceable to the current kit manifest, CLI help or project manifest.

### [US-2] Complete CLI reference
**User Story:** AS A developer, I WANT TO find every Specforce terminal command and flag in one reference page, SO THAT I understand what exists, including the commands my agent runs on my behalf.

**Scenarios:**
1. **[Happy Path]** GIVEN the CLI exposes the command groups init, console, constitution, spec, implementation and archive WHEN a reader opens the CLI reference THEN every subcommand of every group is documented with its arguments and every flag it accepts.
2. **[Happy Path]** GIVEN a command accepts machine-readable output WHEN its entry is read THEN the entry states that the flag exists and that it is intended for automation.
3. **[Edge Case]** GIVEN the archive command has an override for pending tasks WHEN a reader reads its entry THEN it states that the override never bypasses the delta check nor a failed QA verdict.
4. **[Edge Case]** GIVEN a command is used only by the agent (for example recording a QA report from standard input) WHEN a reader reads its entry THEN the entry describes its purpose without instructing the reader to run it as part of the normal flow, and carries no audience label.

**UI/UX Specifics:**
- **View/Component:** Single CLI reference page grouped by area (Global, Constitution, Specification, Delta, Implementation & QA, Archive, Global flags & environment).

**Technical Constraints (NFR):**
- **[Performance]:** 100% of subcommands and flags shown by the CLI help of the current release appear on the page; 0 documented commands or flags absent from it.
- **[Integrity]:** The development-only global flag is documented as intended for kit contributors.

### [US-3] Invoking skills in each supported tool
**User Story:** AS A developer using a specific AI tool, I WANT TO know exactly how to call each Specforce skill in my tool, SO THAT I do not type a syntax my tool does not understand.

**Scenarios:**
1. **[Happy Path]** GIVEN the kit installs the skills discovery, constitution, spec, implement and archive (plus the consultative grill used internally by planning) WHEN a reader opens the supported-tools page THEN a table shows, for each supported tool, where skills are installed and how to invoke them explicitly.
2. **[Happy Path]** GIVEN any other page refers to a skill WHEN the reader reads it THEN it uses the neutral skill name and a plain-language request example, linking to the invocation table.
3. **[Edge Case]** GIVEN a tool has no explicit invocation syntax confirmed by its official documentation WHEN its table row is read THEN the row tells the reader to ask for the skill by name in plain language instead of showing an unverified syntax.
4. **[Edge Case]** GIVEN Codex also writes prompt files outside the project (in the user's Codex prompts folder, overridable by environment variables) WHEN a reader reads the Codex row or note THEN that global write and its location are stated explicitly.

**UI/UX Specifics:**
- **View/Component:** Markdown table with columns Tool, Install folder, Explicit invocation, Notes.

**Technical Constraints (NFR):**
- **[Performance]:** One row per supported tool, eight rows total.
- **[Integrity]:** Every explicit invocation shown is confirmed against the tool's official documentation at implementation time; the source consulted is cited in the design inventory, not in the page.

### [US-4] Lifecycle concepts explained in the artifacts page
**User Story:** AS A developer planning a change, I WANT TO understand spec sizes, required artifacts, module deltas, the QA report and test-first tasks, SO THAT I know what the agent will produce and why a step is blocked.

**Scenarios:**
1. **[Happy Path]** GIVEN the four sizes small, medium, large and complex WHEN a reader reads the artifacts page THEN a table states which artifacts each size requires and the default size.
2. **[Happy Path]** GIVEN feature and bug specification types WHEN the reader reads the page THEN it explains the bug variants of requirements and design and that a bug roadmap starts with a failing regression test.
3. **[Happy Path]** GIVEN medium and larger specs must declare module impact WHEN the reader reads the delta section THEN it explains the delta document (added, modified, removed), the explicit "no module impact" declaration, sealing, each delta check state and that archive requires the delta applied (or none/exempt).
4. **[Happy Path]** GIVEN the independent QA step persists a verdict WHEN the reader reads the QA section THEN it explains the verdicts (passed, failed, not verifiable), that a failed verdict blocks archive and that a missing report only warns.
5. **[Edge Case]** GIVEN a task is marked as test-first (RED) WHEN the reader reads the tasks section THEN it states that finishing it skips only the per-task hooks, while phase and all-tasks hooks still run.
6. **[Edge Case]** GIVEN the reader configures instructions per artifact WHEN they read the configuration page THEN it lists the accepted instruction keys and explains that a specific key (for example the bug variant of requirements) is combined with its generic key.

**UI/UX Specifics:**
- **View/Component:** Artifacts page extended with sections Sizes, Spec metadata, Bug specs, Module deltas, QA report, Test-first tasks; configuration page extended with the instruction keys list.

**Technical Constraints (NFR):**
- **[Performance]:** Each concept section is at most 25 lines so the page stays scannable.
- **[Integrity]:** The task format description matches the current roadmap rules (target plus single executable acceptance check, optional constraints, no verification-only tasks, black-box scenarios section).

### [US-5] Updates and installation channels
**User Story:** AS A developer, I WANT TO know how Specforce updates itself and what to do when it cannot, SO THAT a "manual update required" notice does not surprise me.

**Scenarios:**
1. **[Happy Path]** GIVEN an interactive terminal session WHEN the reader reads the updates section of the getting-started page THEN it explains the background check (at most every 6 hours), that updates follow the installation channel (npm tag of the running channel, or the latest release for native binaries) and that the update applies on the next run.
2. **[Edge Case]** GIVEN the install folder is not writable or an update fails WHEN the reader sees a manual update notice THEN the section shows the exact manual command for npm and for native installs.
3. **[Edge Case]** GIVEN the reader wants to understand an updater decision WHEN they read the section THEN it explains the debug environment variable and the per-user state location.
4. **[Edge Case]** GIVEN commands run by an agent WHEN they run THEN the section states that agents never see update notices and their output is never polluted.

**UI/UX Specifics:**
- **View/Component:** New "Updates" section placed next to Installation and Troubleshooting in getting-started.

**Technical Constraints (NFR):**
- **[Performance]:** Section at most 30 lines.
- **[Integrity]:** Statements match the auto-update module rules of the current release.

### [US-6] Human workflow guide with tips
**User Story:** AS A developer driving an AI agent with Specforce, I WANT TO follow ready workflows and tips for common situations, SO THAT I know what I type in my terminal, what I ask my agent and what happens behind the scenes, without confusing the three.

**Scenarios:**
1. **[Happy Path]** GIVEN the workflow guide WHEN the reader opens it THEN it contains these 12 workflows: prototype in discovery before implementing; fix a bug with a failing test first; make a small change; scope grew (resize); change after implementing but before archiving; change after archiving (follow-up spec); QA failed; parallel specs across worktrees; adopt Specforce in an existing project; resume an interrupted implementation; switch AI tools mid-project; enforce TDD through configuration.
2. **[Happy Path]** GIVEN any workflow WHEN the reader reads it THEN it has: when to use it, numbered human steps where each step is clearly marked as either "in your terminal" or "in your agent chat", an optional collapsed "behind the scenes" block, and an "if it goes wrong" note.
3. **[Edge Case]** GIVEN an action the agent performs on its own WHEN it is mentioned in a workflow THEN it appears only in prose or inside the collapsed "behind the scenes" block, never as a step the reader is told to run.
4. **[Edge Case]** GIVEN the prototype workflow WHEN the reader follows it THEN it explains that discovery keeps product code read-only, that disposable evidence lives in the scratch folder (or a probe-marked file that is cleaned up), and how to ask for formalization into a proposal that planning reuses.
5. **[Edge Case]** GIVEN the switch-tools workflow WHEN the reader passes tool keys to init THEN it warns, before the command runs, that every installed tool left out has its whole folder deleted without confirmation.

**UI/UX Specifics:**
- **View/Component:** New `workflows.md` page in each language, linked from README, getting-started and the docs index list.
- **Feedback Logic:** Terminal steps and chat steps use visually distinct labels consistently across all workflows.

**Technical Constraints (NFR):**
- **[Performance]:** Each workflow at most 40 lines including its collapsed block.
- **[Integrity]:** Every terminal command shown to the reader exists in the current CLI and is one a human is expected to run (init, console, spec list/status, implementation status, spec delta check, or install/update commands).

### [US-7] Edge cases: symptom, cause, action
**User Story:** AS A developer whose flow got stuck, I WANT TO look up the symptom I see, SO THAT I understand why and what to do next.

**Scenarios:**
1. **[Happy Path]** GIVEN the "Common situations" section of the workflow guide WHEN the reader scans it THEN each entry follows Symptom → Why → What to do.
2. **[Happy Path]** GIVEN the section WHEN it is read THEN it covers: archive refused because tasks are pending; archive refused because the delta is missing or pending; delta base conflict; archive refused because QA failed; archive warning because no QA report exists; a failing hook blocking a task; the agent wanting to code without a spec; a worker returning a blocked question; resizing in the middle of an implementation; wanting to edit an already approved spec.
3. **[Edge Case]** GIVEN an entry's resolution requires an agent action WHEN the reader reads "What to do" THEN it tells the reader what to ask the agent in chat, not which internal command to run.

**UI/UX Specifics:**
- **View/Component:** Final section of the workflow guide; one sub-heading per symptom.

**Technical Constraints (NFR):**
- **[Performance]:** 10 entries, each at most 10 lines.

### [US-8] README corrected and reorganized
**User Story:** AS A first-time visitor, I WANT TO understand the lifecycle in the right order and reach the detailed docs, SO THAT I start correctly.

**Scenarios:**
1. **[Happy Path]** GIVEN the README in any language WHEN it is read THEN the lifecycle appears in the order Constitution → Discovery → Spec → Implement → Archive with sequential numbering and neutral skill names linking to the invocation table.
2. **[Happy Path]** GIVEN the documentation list in the README WHEN it is read THEN it includes the workflow guide.
3. **[Edge Case]** GIVEN philosophy and benefit sections WHEN refreshed THEN they are kept but trimmed of claims about decommissioned tools or tool-specific rule files.

**Technical Constraints (NFR):**
- **[Performance]:** README length does not grow by more than 10% over the current English README.

### [US-9] Language parity
**User Story:** AS A Portuguese or Spanish speaking developer, I WANT TO read the same up-to-date content in my language, SO THAT I am not stuck on outdated translations.

**Scenarios:**
1. **[Happy Path]** GIVEN English is the source of truth WHEN the English pages are final THEN Portuguese and Spanish versions exist with identical file names, the same sections and the language selector at the top.
2. **[Happy Path]** GIVEN CONTRIBUTING has Portuguese and Spanish files whose content is currently English WHEN refreshed THEN they are fully translated and carry the same corrected facts.
3. **[Edge Case]** GIVEN a translated page links to another page WHEN the link is followed THEN it resolves to the page in the same language folder.

**Technical Constraints (NFR):**
- **[Performance]:** The three `docs/` language folders contain the same set of 7 file names.
- **[Integrity]:** Commands, flags, file names and skill names are never translated.

## 4. Business Invariants
- English is the source of truth; translations never contain facts absent from the English page.
- File names are identical across the `docs/en`, `docs/pt` and `docs/es` folders.
- Every relative link in README, CONTRIBUTING and `docs/` resolves to an existing file.
- No page lists a decommissioned agent or uses the legacy colon invocation syntax.
- No page documents a command or flag that the current CLI does not expose.
- In the workflow guide, only commands a human is expected to run appear as steps for the reader.
- The local kit override is not documented anywhere.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Workflow guide readable section by section: 12 workflows at most 40 lines each and 10 edge-case entries at most 10 lines each.
- **[Reliability]:** Every factual statement is checked against the current CLI help, kit manifest, project manifest or module living specs before it is written.
- **[Security]:** Examples never include real credentials or tokens; secrets are referenced as environment files only.
- **[Maintainability]:** Pages keep the existing language-selector header format and heading hierarchy so future translations stay aligned.
