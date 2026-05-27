---
slug: 20260525-2251-spec-generation-optimization
lens: Backend-heavy
---

# Feature: Spec Generation Optimization

## 1. Context & Value
AI agents often fail to follow project-specific instructions (e.g., "Use TDD") because they are appended at the end of long system prompts, where they lose prominence. Additionally, the current hardcoded generation order for spec artifacts does not scale for custom types. This feature optimizes instruction visibility and ensures a deterministic, dependency-aware generation pipeline to maximize agent compliance and system scalability.

## 2. Out of Scope (Anti-Goals)
- Implementing a full UI for spec status beyond the existing TUI/CLI output.
- Modifying the underlying AI model logic for following instructions.
- Supporting non-YAML artifact definitions.

## 3. Acceptance Criteria (BDD)

### [US-1] Top-Level Instruction Injection
**User Story:** AS AN AI Orchestrator, I WANT TO receive project-specific instructions at the beginning of my prompt, SO THAT I can prioritize them over default patterns.

**Scenarios:**
1. **[Happy Path]** GIVEN a `config.yaml` with `instructions.tasks` set to "Use TDD" WHEN I run `specforce spec artifact tasks --json` THEN the "Project Specific Instructions" section must be prepended to the `instruction` field.
2. **[Edge Case]** GIVEN no custom instructions exist WHEN the artifact is fetched THEN the `instruction` field must remain unchanged and contain only the base instructions.

**Technical Constraints (NFR):**
- **[Performance]:** Instruction merging must add < 5ms to the command execution.
- **[Safety & Security]:** Injected variables must be properly escaped to prevent prompt injection.

### [US-2] Dependency-Aware Generation Order
**User Story:** AS A Developer, I WANT Specforce to generate artifacts in their logical dependency order, SO THAT each document has the necessary context from its upstream dependency.

**Scenarios:**
1. **[Happy Path]** GIVEN a registry with `tasks` depending on `design` and `design` depending on `requirements` WHEN `spec status` lists artifacts THEN they must appear in the order: requirements, design, tasks.
2. **[Edge Case]** GIVEN a circular dependency is detected (e.g., A -> B -> A) WHEN the registry is initialized THEN the system must fail with a clear error message.

**Technical Constraints (NFR):**
- **[Performance]:** Topological sort must handle up to 100 artifacts in < 10ms.
- **[Reliability]:** The sorting algorithm must be deterministic across all platforms.

### [US-3] Mission Brief Subagent Protocol
**User Story:** AS THE Specforce Orchestrator, I WANT TO prepend a "Mission Brief" to all subagent prompts, SO THAT project rules are framed as non-negotiable constraints.

**Scenarios:**
1. **[Happy Path]** GIVEN the `spf.spec` skill is active WHEN a subagent is invoked THEN the orchestrator must prepend a header containing the artifact name, description, and project rules.
2. **[Edge Case]** GIVEN a subagent fails to follow the mission brief WHEN the output is verified THEN the orchestrator must attempt a clarification loop referencing the brief.

**Technical Constraints (NFR):**
- **[Performance]:** Prompt construction must not exceed the environment's token limits.
- **[Observability]:** The mission brief must be visible in the agent's internal reasoning logs.

## 4. Business Invariants
- Project-specific instructions MUST always override base instructions if a conflict exists.
- An artifact CANNOT be marked as valid if its dependency is missing or invalid.
- The `slug` must remain the single source of truth for file paths.

## 5. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Spec generation and validation must remain instantaneous (< 1s total).
- **[Reliability]:** Zero silent failures during artifact generation.
- **[Security]:** Zero leakage of internal file paths or system configurations in subagent prompts.
- **[Maintainability]:** Use standard Go idioms and ensure 80% code coverage for new logic.
