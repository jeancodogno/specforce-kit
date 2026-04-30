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

---

## Critical Now
- [2026-04-19] [Tasks] Hierarchical task organization (Phases H3 / Tasks H4) implemented. Agents should now use this for all new feature roadmaps.
- [2026-04-28] [Autonomous Workflows] Proactive Mandate in `AGENTS.md` and dual-installation (Commands as Skills) are now live. Agents should automatically trigger `/spec` or `/implement`.
- [2026-04-29] [Discovery] Conversational `/discovery` mode (Specforce Scout) is live. It's purely read-only and designed for brainstorming or bug-hunting before `/spec`.

## Last Actions
- **Date:** 2026-04-29
- **Scope:** Discovery Command
- **Completed:** Implemented `spf.discovery` (Specforce Scout), updated `AGENTS.md` template, and enriched prompt with unique branding and diagnostic workflows.
- **Next:** Monitor user engagement with Discovery mode to refine 'Thinking Partner' prompts.
- **Relevant Files:** src/internal/agent/kit/commands/discovery.yaml, src/internal/project/agents_md.go, README.md

## Active Lessons & Anti-Patterns
- **First Seen:** 2026-04-29
- **Last Seen:** 2026-04-29
- **Scope:** Agent / Kit Mapping
- **Symptom:** AI Agents failing to distinguish between commands because multiple `SKILL.md` files have the same `name: SKILL` in their headers.
- **Avoid:** Using the mapping-level `name` as the primary identity for generated frontmatter if that name is a generic placeholder like "SKILL".
- **Do Instead:** Resolve identity hierarchically: Metadata Name > Blueprint Slug (prefixed with `spf.` for commands) > Mapping Name.
- **Recurrence Count:** 1
- **Status:** Resolved
- **Distill To:** engineering.md

- **First Seen:** 2026-04-29
- **Last Seen:** 2026-04-29
- **Scope:** Artifact Generation / Security
- **Symptom:** Multiple blueprints resulting in the same header name, causing non-deterministic discovery.
- **Avoid:** Blindly generating artifacts without cross-checking for global identity collisions.
- **Do Instead:** Use a session-level `nameTracker` to validate that every generated header `name` is unique across all artifacts for a given agent.
- **Recurrence Count:** 1
- **Status:** Resolved
- **Distill To:** engineering.md

- **First Seen:** 2026-04-29
- **Last Seen:** 2026-04-29
- **Scope:** Discovery / Implementation
- **Symptom:** Ambiguity between unstructured exploration and formal SDD lifecycle.
- **Avoid:** Letting discovery sessions wander into ad-hoc implementation without specification.
- **Do Instead:** Enforce a strict read-only prompt for Discovery and mandate a handoff summary that recommends `/spec` to formalize insights.
- **Recurrence Count:** 1
- **Status:** Active
- **Distill To:** governance.md

## Pending Decisions (Need Distillation)
- **Date:** 2026-04-28
- **Scope:** Architecture / Kit
- **Decision:** Support `MappingConfigs` as a slice in `kit.yaml` to allow multiple destinations for the same category.
- **Why:** Enables "Commands as Skills" without logic duplication.
- **Validate By:** `src/internal/agent/dual_install_test.go`
- **Distill To:** architecture.md
