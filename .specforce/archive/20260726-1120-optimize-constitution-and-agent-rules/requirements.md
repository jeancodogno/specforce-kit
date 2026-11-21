---
slug: optimize-constitution-and-agent-rules
lens: Integration
---

# Feature: Optimize Constitution and Agent Rules Guidance

## 1. Context & Value
Currently, Specforce workflow commands (`spf.discovery` and `spf.spec`) force agents to execute `specforce constitution status --json` before reading project rules, creating unnecessary CLI execution and token overhead when constitution rules are already referenced in `AGENTS.md`. Furthermore, `AGENTS.md` lacks explicit guidance regarding domain-specific modules stored in `.specforce/docs/modules/<slug>.md`. This feature streamlines agent instructions by prioritizing direct surgical file reads of constitution and module docs in `AGENTS.md` and command specs while removing redundant CLI calls.

## 2. Out of Scope (Anti-Goals)
- Do NOT remove or modify `specforce constitution status` commands inside `spf.constitution` (`constitution.yaml`), as management/editing of the constitution still requires CLI status discovery.
- Do NOT change the behavior of the `specforce` CLI Go binary commands or alter CLI flags.
- Do NOT alter any external agent integration contracts beyond template text and command prompt definitions.

## 3. Acceptance Criteria (BDD)

### [US-1] Add Module Guidance to AGENTS.md Template
**User Story:** AS AN AI Agent, I WANT TO be informed about `.specforce/docs/modules/<slug>.md` files in `AGENTS.md`, SO THAT I can lazy-load domain-specific rules without redundant CLI discovery.

**Scenarios:**
1. **[Happy Path]** GIVEN the template generator in `src/internal/project/agents_md.go` WHEN `EnsureAgentsMD` is invoked THEN the generated `AGENTS.md` Section 4 contains clear instructions on `modules/<slug>.md` alongside global constitution files.
2. **[Edge Case]** GIVEN an existing `AGENTS.md` with Specforce markers WHEN `EnsureAgentsMD` runs THEN the Section 4 update merges cleanly between `<!-- SPECFORCE_AGENTS_START -->` and `<!-- SPECFORCE_AGENTS_END -->` without corrupting external user content.

**Technical Constraints (NFR):**
- **[Performance]:** Zero additional CLI latency during `specforce init`.
- **[Safety & Security]:** File permissions 0600 maintained for `AGENTS.md`.
- **[Integrity]:** Unit test coverage in `agents_md_test.go` verifies presence of module guidance.
- **[Observability]:** Output via `core.UI` remains clean.

### [US-2] Remove Redundant Constitution CLI Status Calls from Workflow Commands
**User Story:** AS AN AI Agent executing `spf.discovery` or `spf.spec`, I WANT TO read relevant `.specforce/docs/` files directly without being forced to run `specforce constitution status --json`, SO THAT token consumption and execution overhead are minimized.

**Scenarios:**
1. **[Happy Path]** GIVEN the command definitions `discovery.yaml` and `spec.yaml` in `src/internal/agent/kit/commands` WHEN an agent inspects Layer 1 (Constitutional Anchor) THEN the prompt instructs direct reading of `.specforce/docs/*.md` and `.specforce/docs/modules/*.md` instead of mandating `specforce constitution status --json`.
2. **[Edge Case]** GIVEN `constitution.yaml` WHEN an agent inspects it THEN `specforce constitution status --json` remains present as a necessary step for managing constitution state.

**Technical Constraints (NFR):**
- **[Performance]:** Eliminates 1 CLI roundtrip per discovery/spec initialization session.
- **[Maintainability]:** Clean YAML syntax without broken multi-line markdown strings.

## 4. Business Invariants
- `AGENTS.md` managed section must always stay bounded between `<!-- SPECFORCE_AGENTS_START -->` and `<!-- SPECFORCE_AGENTS_END -->`.
- All `spf.*` command definitions must remain valid YAML files.

## 5. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Fast test execution under `go test ./...`.
- **[Reliability]:** Fail-fast unit tests for template generation.
- **[Maintainability]:** 100% Go test pass rate.
