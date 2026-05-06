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
- **Scope:** Specialized Bugfix Templates & YAML Metadata
- **Completed:** Introduced `spec.Metadata` (`spec.yaml`) to track spec types (feature/bug) and lens. Updated `Registry` to support category-aware artifacts with `{type}-{name}.yaml` naming convention. Enhanced `spec init` with `--type` flag and `spec status` with type-aware filtering. Added specialized `bug-requirements` and `bug-design` templates.
- **Next:** Monitor for any edge cases where spec type might need manual adjustment or migration of old specs to the YAML format.
- **Relevant Files:** src/internal/spec/metadata.go, src/internal/spec/registry.go, src/internal/cli/spec.go, src/internal/agent/artifacts/spec/bug-*.yaml

- **Date:** 2026-05-06
- **Scope:** Hardened Spec Verification (Planning Phase)

## Active Lessons & Anti-Patterns
- **First Seen:** 2026-05-06
- **Last Seen:** 2026-05-06
- **Scope:** Architecture / Metadata
- **Symptom:** Difficulty in managing specification-wide state (type, lens, domain) when metadata is buried in Markdown frontmatter.
- **Avoid:** Relying solely on semi-structured Markdown headers for critical system-level decisions like template selection.
- **Do Instead:** Decouple metadata into a dedicated machine-readable file (`spec.yaml`) within the specification directory. This enables robust, type-aware behavior across all CLI commands and agent prompts.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** architecture.md

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
