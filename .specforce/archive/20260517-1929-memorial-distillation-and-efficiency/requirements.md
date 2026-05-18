---
slug: 20260517-1929-memorial-distillation-and-efficiency
lens: Backend-heavy
---

# Feature: Memorial Distillation and Efficiency

## 1. Context & Value
The current memorial system is "write-only," recording fragments that are never read by agents, and lacks a cleanup mechanism for accumulated files. This feature activates memorial playback to ensure cross-session technical continuity and implements distillation to maintain a lean, high-signal project memory.

## 2. Out of Scope (Anti-Goals)
- Do not implement automatic periodic distillation (must be triggered by the user).
- Do not modify the existing `ROUTING.md` structure.
- Do not build a TUI interface for distillation (CLI-only for now).

## 3. Acceptance Criteria (BDD)

### [US-1] Memorial Content Playback
**User Story:** AS AN AI Agent, I WANT TO receive the actual content of recent memorial fragments in my instructions, SO THAT I can reuse established precedents and avoid repeating errors.

**Scenarios:**
1. **[Happy Path]** GIVEN 5 existing memorial fragments WHEN I run `specforce archive instructions` THEN the output contains the full markdown content of all 5 fragments.
2. **[Edge Case]** GIVEN more than 10 fragments WHEN I run `specforce archive instructions` THEN the output contains only the 10 most recent fragments, ordered by date (newest first).

**Technical Constraints (NFR):**
- **[Performance]:** Instruction generation must remain under 200ms.
- **[Integrity]:** Fragments must be injected as raw Markdown to preserve formatting.
- **[Observability]:** The output must clearly delineate each fragment with its timestamp and title.

### [US-2] Smart Distillation Command
**User Story:** AS A Developer, I WANT TO consolidate old memory fragments into a single file using AI summarization, SO THAT the `.specforce/memorial/` directory remains clean and easy to scan.

**Scenarios:**
1. **[Happy Path]** GIVEN multiple fragments WHEN I run `specforce archive distill --slug <slugs>` THEN the agent generates a consolidated summary of those fragments and appends it to `.specforce/memorial/DISTILLED.md`.
2. **[Edge Case]** GIVEN a non-existent slug WHEN running distillation THEN the command must fail with a clear "Fragment not found" error.

**Technical Constraints (NFR):**
- **[Performance]:** Distillation (excluding agent processing) must be instantaneous.
- **[Safety & Security]:** Original fragment files MUST be deleted ONLY after the distilled summary is successfully written.
- **[Integrity]:** The `DISTILLED.md` file must follow the same frontmatter and header standards as individual fragments.

### [US-3] Legacy Memorial Cleanup
**User Story:** AS A Developer, I WANT TO remove the stale monolithic memorial file, SO THAT agents are not confused by obsolete information.

**Scenarios:**
1. **[Happy Path]** GIVEN the existence of `.specforce/docs/memorial.md` WHEN the feature is implemented THEN the file is deleted.
2. **[Edge Case]** GIVEN the file is already missing WHEN the cleanup runs THEN it should proceed silently without error.

**Technical Constraints (NFR):**
- **[Performance]:** File deletion must be instantaneous.
- **[Reliability]:** Zero-impact if the file doesn't exist.

## 4. Business Invariants
- The system MUST always preserve the `ROUTING.md` file as the primary entry point.
- Individual fragments MUST be globally unique by their slug/timestamp.
- Distillation MUST NOT lose critical architectural decisions or active lessons.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Latency for content injection < 100ms.
- **[Reliability]:** Fail-fast if the memorial directory is inaccessible.
- **[Security]:** All memorial files (including distilled) must maintain 0600 permissions.
- **[Maintainability]:** 80% coverage gate for the new distillation logic.
