# Memorial (Active Memory)

> **FOR AI AGENTS: RULES OF ENGAGEMENT**
> This is your cross-session memory file. You MUST obey these rules:
> 1. **READ THIS FIRST** before starting any task, debugging, or writing code in this repository.
> 2. **TRIAGE:** Read `Critical Now` first, then `Last Actions`. Only read lessons if they match your current scope.
> 3. **STRICT LIMITS (DO NOT BLOAT):**
>    - Max 3 `Critical Now` bullets.
>    - Max 5 `Last Actions` entries.
>    - Max 7 `Active Lessons`.
>    - Max 5 `Pending Decisions`.
> 4. **GARBAGE COLLECTION:** If a list is full, you MUST delete the oldest/least relevant item before adding a new one.
> 5. **DISTILLATION:** If a lesson or decision becomes a permanent rule, you MUST move it to the official Constitution files (`engineering.md`, `architecture.md`, etc.) and completely delete it from here.

## Critical Now
- [2026-05-05] [Agent/Protocol] SDD Trigger Optimization implemented. Agents now have explicit mandates for triggering `spf.discovery` (vague intent/bug investigation) and `spf.spec` (planning). Direct edits to codebase without an approved spec are now strictly prohibited in `AGENTS.md`.
- [2026-05-05] [Project/Init] Conditional tool folder creation implemented. `.gemini`, `.claude`, and `.agent` directories are now only created if selected or already existing.
- [2026-05-04] [Agent/Kit] Agent/Skill reorganization (v1.x) implemented. Dynamic Registry with multi-source discovery (Embedded + Local) and variable injection (`{{context}}`) are now active.

## Last Actions
- **Date:** 2026-05-06
- **Scope:** Shields.io & Navigation Integration
- **Completed:** Integrated 8 technical metric badges (CI, Release, NPM, Go, Go Report, Issues, PRs, License) and established a visual Language Switcher system across 3 READMEs and 18 technical documentation files. Implemented localized relative navigation headers in all docs subdirectories.
- **Next:** Monitor for link breakage if documentation files are renamed or moved across language directories.
- **Relevant Files:** README.md, README.pt.md, README.es.md, docs/**/*

- **Date:** 2026-05-06
- **Scope:** Specialized Bugfix Templates & YAML Metadata
- **Completed:** Introduced `spec.Metadata` (`spec.yaml`) to track spec types (feature/bug) and lens. Updated `Registry` to support category-aware artifacts with `{type}-{name}.yaml` naming convention. Enhanced `spec init` with `--type` flag and `spec status` with type-aware filtering. Added specialized `bug-requirements` and `bug-design` templates.
- **Relevant Files:** src/internal/spec/metadata.go, src/internal/spec/registry.go, src/internal/cli/spec.go, src/internal/agent/artifacts/spec/bug-*.yaml

## Active Lessons & Anti-Patterns
- **First Seen:** 2026-05-06
- **Last Seen:** 2026-05-06
- **Scope:** UI-UX / Documentation
- **Symptom:** Users get "lost" when navigating technical documentation in non-English languages due to a lack of immediate "language-back" or "language-cross" links.
- **Avoid:** Relying on a single global language switcher at the project root (README).
- **Do Instead:** Prepend a standardized relative navigation header to EVERY technical documentation file to ensure immediate context switching at the point of consumption.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** ui-ux.md, governance.md

- **First Seen:** 2026-05-06
- **Last Seen:** 2026-05-06
- **Scope:** SDD Protocol / Verification

- **First Seen:** 2026-04-30
- **Last Seen:** 2026-04-30
- **Scope:** Proxy / Environment
- **Symptom:** Users encounter "command not found" even after successful global install if npm bin is missing from PATH.
- **Avoid:** Silently failing or giving generic errors when native binaries are missing.
- **Do Instead:** Run an automated diagnostic check that identifies the exact missing PATH entry and provides the specific OS-level fix command.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** engineering.md

## Pending Decisions (Need Distillation)
- **Date:** 2026-05-04
- **Scope:** Parser / SDD Protocol
- **Decision:** Transition the CLI parser to support "Natural LLM Task Formats" (standard markdown checklists `- [ ]`) instead of forcing strict `#### T1.1` headers, to reduce LLM formatting drift and align with OpenSpec.
- **Why:** Imposing rigid formatting constraints wastes tokens, breaks automation when LLMs drift to natural checklists, and creates a high "translation tax".
- **Validate By:** `src/internal/spec/tasks_test.go`
- **Distill To:** engineering.md

- **Date:** 2026-04-28
- **Scope:** Architecture / Kit
- **Decision:** Support `MappingConfigs` as a slice in `kit.yaml` to allow multiple destinations for the same category.
- **Why:** Enables "Commands as Skills" without logic duplication.
- **Validate By:** `src/internal/agent/dual_install_test.go`
- **Distill To:** architecture.md

- **Date:** 2026-05-06
- **Scope:** SDD Protocol / Requirements
- **Decision:** Shift from global UI/UX and NFR blocks to localized context within each functional requirement ([REQ-x]), utilizing a flexible "Attribute: Value" tagging system.
- **Why:** Reduces implementation "hallucination" by providing localized technical and visual constraints exactly where the behavior is defined. Prevents "N/A" bloat in templates while ensuring "Performance" and "Safety" are explicitly addressed.
- **Validate By:** `src/internal/agent/artifacts/spec/requirements.yaml` and resulting `requirements.md` artifacts.
- **Distill To:** engineering.md, ui-ux.md
