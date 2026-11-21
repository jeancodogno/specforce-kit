---
slug: 20260520-1005-layered-instruction-mapping
lens: Integration
---

# Feature: Layered Instruction Mapping

## 1. Context & Value
The Specforce CLI currently fails to inject `config.yaml` instructions when an artifact name contains a prefix (e.g., `feature-requirements` instead of `requirements`). This feature implements a layered injection strategy that combines generic base-type instructions with specific artifact-level instructions, ensuring agents always receive both global project standards and specific context.

## 2. Out of Scope (Anti-Goals)
- Support for regex-based or fuzzy matching in `config.yaml` keys.
- Modifying the schema of `config.yaml` (uses existing `instructions` block).
- External instruction fetching (API or remote files).

## 3. Acceptance Criteria (BDD)

### [US-1] Dual-Layer Instruction Injection
**User Story:** AS AN AI agent, I WANT TO receive both generic and specific instructions based on the requested artifact name, SO THAT I can comply with both project-wide and task-specific rules.

**Scenarios:**
1. **[Happy Path] Generic and Specific Combination**
   GIVEN a `config.yaml` with `requirements` (generic) and `feature-requirements` (specific) keys
   WHEN I request an artifact named `feature-requirements`
   THEN the CLI MUST inject instructions from both keys into the agent prompt.

2. **[Happy Path] Generic Only Mapping**
   GIVEN a `config.yaml` with only the `requirements` key
   WHEN I request an artifact named `bugfix-requirements`
   THEN the CLI MUST identify `requirements` as the base type and inject its instructions.

3. **[Edge Case] Duplicate Instructions**
   GIVEN identical instruction strings in both generic and specific keys
   WHEN the CLI merges the instructions
   THEN it SHOULD deduplicate the list to minimize token usage.

4. **[Edge Case] Non-existent Keys**
   GIVEN a requested artifact name `custom-artifact` that matches no base types
   WHEN the CLI processes instructions
   THEN it MUST proceed without error, providing an empty or default instruction set.

**UI/UX Specifics:**
- **View/Component:** CLI Agent Prompt (Internal).
- **Feedback Logic:** Silent operation; affects the instructions block sent to the agent.
- **Keybindings:** N/A.

**Technical Constraints (NFR):**
- **[Performance]:** Instruction mapping and merging MUST complete in < 5ms.
- **[Safety & Security]:** Instructions MUST be treated as literal text; no execution of code within instruction strings.
- **[Integrity]:** The order of injection MUST be Generic FIRST, followed by Specific, to maintain logical precedence.
- **[Observability]:** Log matched instruction keys when debug logging is enabled.

### [US-2] Base Type Inference
**User Story:** AS A Developer, I WANT the CLI to automatically detect the base artifact type from a prefixed string, SO THAT standard rules are applied to all related artifacts.

**Scenarios:**
1. **[Happy Path] Hyphenated Prefix Detection**
   GIVEN the base types: `requirements`, `design`, `tasks`, `implementation`
   WHEN an artifact name is `frontend-design`
   THEN the CLI MUST correctly map it to the generic `design` type.

2. **[Edge Case] Multiple Base Type Keywords**
   GIVEN an artifact named `tasks-for-design-artifact`
   WHEN the CLI infers the base type
   THEN it SHOULD prioritize the first occurrence or the most relevant keyword based on a right-to-left suffix match (e.g., `design-artifact` -> `design`).

**UI/UX Specifics:**
- **View/Component:** CLI Output.
- **Feedback Logic:** Silent.
- **Keybindings:** N/A.

**Technical Constraints (NFR):**
- **[Integrity]:** Inference MUST be deterministic and based on the core Specforce artifact lifecycle stages.

## 4. Business Invariants
- Generic instructions for a base type MUST be applied to any artifact identified as that type, regardless of prefixing.
- Specific instructions mapped to the exact artifact name MUST be appended if they exist.
- Missing generic or specific keys MUST NOT result in a failure (graceful degradation).

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Instruction resolution must be negligible (<5ms).
- **[Reliability]:** Must not crash if instructions are missing in config.yaml.
- **[Maintainability]:** Clean separation of generic vs specific instruction injection logic.
