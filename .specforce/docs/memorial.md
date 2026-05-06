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
- **Date:** 2026-05-05
- **Scope:** SDD Protocol Trigger Optimization
- **Completed:** Updated `AGENTS.md` and Go template in `agents_md.go` with explicit triggers and direct edit prohibitions. Clarified "Primary Orchestration Only" in `engineering.md`. Refined metadata keywords in source kit (`discovery.yaml`, `spec.yaml`) to improve agent auto-triggering.
- **Next:** Monitor agent behavior to ensure they don't bypass the protocol during complex bug investigations.
- **Relevant Files:** AGENTS.md, src/internal/project/agents_md.go, src/internal/agent/kit/commands/*, .specforce/docs/engineering.md

- **Date:** 2026-05-05
- **Scope:** Conditional Tool Folder Creation (Init Fix)
- **Completed:** Refactored `EnsureAgentsMD` and `ensurePlatformConfigs` in `project` package to accept `selectedAgents`. Removed premature `EnsureAgentsMD` call from `BootstrapProject`. Orchestrated final configuration in `Service.InitializeProject` after tool selection is known.
- **Next:** Monitor for any other tools that might need similar conditional logic in the future.
- **Relevant Files:** src/internal/project/agents_md.go, src/internal/project/service.go, src/internal/project/bootstrapper.go

- **Date:** 2026-05-04
- **Scope:** Agent & Skill Reorganization (v1.x)
- **Completed:** Refactored `agent.Registry` for multi-source discovery, updated `kit.yaml` schema with `defaults` and `security` metadata, and implemented `InstructionManager` for dynamic `{{variable}}` injection. Updated EN/PT/ES docs.
- **Next:** Monitor for any variable injection collisions in complex templates.
- **Relevant Files:** src/internal/agent/registry.go, src/internal/agent/translator.go, src/internal/agent/instructions.go, kit.yaml, docs/en/configuration.md

- **Date:** 2026-05-04
- **Scope:** Auto-Timestamped Spec Slugs
- **Completed:** Implemented `PrepareSlug` logic in `spec` package, integrated into `HandleSpecInit`, and added unit/integration tests. Created new standard for spec directory naming.
- **Next:** Monitor for any issues with nested directory timestamping in CI.
- **Relevant Files:** src/internal/spec/slug.go, src/internal/cli/spec.go, src/internal/cli/integration_test.go

- **Date:** 2026-05-04
- **Scope:** Multi-language Documentation Expansion
- **Completed:** Restructured `docs/` folder into `en/`, `pt/`, and `es/`. Migrated English docs, added language selector to root README, and created PT/ES localized versions of root documents.
- **Next:** Continue translating specific documentation files within the language folders.
- **Relevant Files:** README.md, README.pt.md, README.es.md, docs/{en,pt,es}/*

## Active Lessons & Anti-Patterns
- **First Seen:** 2026-05-05
- **Last Seen:** 2026-05-05
- **Scope:** Agent / SDD Protocol
- **Symptom:** Agents bypass the SDD pipeline and edit code directly when user intent is slightly vague or a bug is reported.
- **Avoid:** Providing generic "proactive" mandates without explicit trigger conditions or strict prohibitions on tool usage.
- **Do Instead:** Define explicit intent-to-workflow mappings (Discovery for vague/bugs, Spec for planning, Implement for roads) and strictly forbid mutation tools (`replace`/`write_file`) if no approved spec exists.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** engineering.md

- **First Seen:** 2026-05-05
- **Last Seen:** 2026-05-05
- **Scope:** Project / Initialization
- **Symptom:** Project root becomes cluttered with unused tool directories (.gemini, .claude) during `init`.
- **Avoid:** Unconditionally creating platform-specific directories in "bootstrapping" phases before user preferences (selections) are resolved.
- **Do Instead:** Defer platform-specific configuration and directory creation until after the selection logic is complete. Pass the list of `selectedAgents` down to the configuration logic.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** engineering.md

- **First Seen:** 2026-05-04
- **Last Seen:** 2026-05-04
- **Scope:** Agent / Infrastructure
- **Symptom:** Cognitive complexity spikes in FS-walking functions (like `scanSkills`) when adding metadata parsing logic.
- **Avoid:** Nesting heavy conditional logic inside `WalkDir` functions.
- **Do Instead:** Extract metadata loading into surgical helper functions (e.g., `loadSkillMetadata`) to keep the traversal loop clean and testable.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** engineering.md

- **First Seen:** 2026-05-04
- **Last Seen:** 2026-05-04
- **Scope:** CLI / UX
- **Symptom:** Users might not notice if a slug is transformed (timestamped) unless explicitly told.
- **Avoid:** Transforming user input silently without feedback.
- **Do Instead:** Always print the final resolved path/slug in the success message of `spec init` to maintain transparency.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** engineering.md

- **First Seen:** 2026-05-04
- **Last Seen:** 2026-05-04
- **Scope:** Documentation / Migration
- **Symptom:** Links breaking during documentation folder restructuring.
- **Avoid:** Moving documentation files without a recursive link verification step.
- **Do Instead:** Use a systematic approach to update all relative paths in `README.md` and internal `.md` files immediately after moving files to nested subdirectories.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** engineering.md

- **First Seen:** 2026-05-01
- **Last Seen:** 2026-05-01
- **Scope:** Testing / Platform
- **Symptom:** Tests relying on `os.Unsetenv("HOME")` to force `os.UserHomeDir` failure fail on macOS because of system fallback.
- **Avoid:** Assuming `os.UserHomeDir` will always fail if environment variables are removed.
- **Do Instead:** Explicitly check if the function still succeeds after `Unsetenv` and use `t.Skip` to handle platforms with persistent home directory resolution.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** engineering.md

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
