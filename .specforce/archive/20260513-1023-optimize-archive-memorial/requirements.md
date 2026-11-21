---
slug: 20260513-1023-optimize-archive-memorial
lens: Backend-heavy
---

# Feature: Optimize Archive Memorial

## 1. Context & Value
Current archival workflow requires agents to manually create memorial fragments, leading to redundant 'date' and 'ls' shell calls. This feature automates fragment creation via a dedicated CLI command, reducing turn count and ensuring metadata consistency.

## 2. Out of Scope (Anti-Goals)
- Do not modify the 'spec archive' core move logic.
- Do not implement automatic Constitution updates in this spec.
- Do not build a TUI for memorial recording.

## 3. User Stories & Acceptance Criteria

### [US-1] Memorial Recording CLI Command
As an AI Agent, I want a dedicated command to record memorial entries so that I don't have to manually manage file paths, timestamps, or formatting.

**Acceptance Criteria:**
- **AC 1.1:** The command `specforce archive memorial <slug> --type <type> --title <title> --content <content>` must be implemented.
- **AC 1.2:** The command must automatically generate a timestamped markdown file in `.specforce/memorial/`.
- **AC 1.3:** The file name must follow the pattern `YYYYMMDD-HHmm-<slug>.md`.
- **AC 1.4:** The content must be formatted as a standardized memorial fragment including title, date, and content sections.
- **AC 1.5:** Supported types must include at least `lesson`, `decision`, and `context`.

**Non-Functional Requirements:**
- **NFR 1.1:** Command execution must complete in less than 500ms.
- **NFR 1.2:** Command must provide clear feedback upon successful file creation.

### [US-2] Enforce Memorial Command in Archival Blueprint
As a Product Owner, I want the archival process to strictly use the memorial command so that the workflow is standardized and efficient.

**Acceptance Criteria:**
- **AC 2.1:** The `archive.yaml` blueprint (agent instruction) must be updated to include the memorial recording step using the new CLI command.
- **AC 2.2:** The blueprint must explicitly forbid manual file creation in the `.specforce/memorial/` directory.
- **AC 2.3:** Instructions to run `date` or `ls` for the purpose of creating memorial entries must be removed from the blueprint.

**Non-Functional Requirements:**
- **NFR 2.1:** The updated blueprint must pass agent instruction validation (if applicable).

## 4. Key Entities & Data Boundaries

### Memorial Fragment
- **ID:** Auto-generated filename (`YYYYMMDD-HHmm-<slug>.md`).
- **Type:** enum {lesson, decision, context}.
- **Title:** String (required).
- **Content:** Markdown text (required).
- **Timestamp:** Auto-generated during command execution.

## 5. Success Metrics
- **Business Metric:** Reduce the number of agent turns during the `archive` phase by at least 2 turns (eliminating `date` and `ls`).
- **Performance Target:** Memorial entry creation time < 1s.
- **UX Efficiency:** Agent can record a memorial entry with a single atomic command.

## 6. Functional & Safety Guards
- **Validation:** The command must fail if the specified `<slug>` does not exist in `.specforce/archive/` (optional, depends on implementation details but good for safety).
- **Idempotency:** The command allows multiple memorial entries for the same slug. Uniqueness is ensured by the minute-level timestamp in the filename (YYYYMMDD-HHmm).
- **Multiple Entries:** AI agents can record multiple findings (e.g., a lesson and a decision) during the same archival phase without collision.
