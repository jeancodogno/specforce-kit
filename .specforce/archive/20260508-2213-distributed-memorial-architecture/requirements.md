---
slug: 20260508-2213-distributed-memorial-architecture
lens: Backend-heavy
---

# Feature: Distributed Memorial Architecture

## 1. Context & Value
The current `memorial.md` is a single file that causes frequent merge conflicts in multi-developer environments because it's updated at the end of every feature lifecycle. This feature replaces the monolithic file with a distributed directory of memory fragments (`.specforce/memorial/*.md`), enabling parallel contributions while maintaining the AI agent's ability to consume consolidated project memory.

## 2. Out of Scope (Anti-Goals)
- Do not implement automatic "distillation" (moving fragments to the Constitution) in this spec.
- Do not build a TUI for browsing fragments yet.
- Do not implement historical versioning beyond Git for the fragments.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Directory-Based Memory Storage
**User Story:** AS A developer, I WANT the system to store memory fragments in individual files, SO THAT I can merge my changes without conflicts with other team members.

**Scenarios:**
1. **[Happy Path]** GIVEN a project is initialized WHEN a new memory is recorded THEN it is saved as a unique file in `.specforce/memorial/` using a timestamp-slug format.
2. **[Migration Path]** GIVEN an existing `memorial.md` WHEN the system is updated THEN the legacy file is moved to `.specforce/memorial/legacy.md` or similar to preserve history.

**Technical Constraints (NFR):**
- **[Performance]:** Fragment creation latency < 10ms.
- **[Safety & Security]:** Fragment files MUST use 0600 permissions.
- **[Integrity]:** Filenames MUST be deterministic and unique (e.g., `YYYYMMDD-HHMM-slug.md`).

### [REQ-2] Consolidated AI Reading Mode
**User Story:** AS AN AI Agent, I WANT to read a single consolidated view of project memory, SO THAT I can quickly synchronize with the project context without scanning dozens of files.

**Scenarios:**
1. **[Happy Path]** GIVEN multiple fragment files WHEN the agent starts a session THEN the CLI provides a "Consolidated Memorial" containing the Rules of Engagement and the most recent N fragments.

**Technical Constraints (NFR):**
- **[Performance]:** Consolidation of 50 fragments MUST take < 50ms.
- **[Observability]:** The consolidated output MUST explicitly state which fragments were included.

### [REQ-3] Workflow Integration (spf.archive)
**User Story:** AS A developer, I WANT the `spf.archive` skill to automatically create a new fragment, SO THAT I don't have to manually manage memory files.

**Scenarios:**
1. **[Happy Path]** GIVEN a completed feature WHEN I run `spf.archive` THEN a new fragment is created in `.specforce/memorial/` instead of appending to a monolithic file.

**Technical Constraints (NFR):**
- **[Integrity]:** All mandatory fields (Date, Scope, Lessons) from the legacy memorial schema MUST be preserved in the fragment format.

### [REQ-4] Project Initialization & AGENTS.md
**User Story:** AS A user, I WANT `specforce init` to set up the new memorial structure, SO THAT my project starts with conflict-proof memory.

**Scenarios:**
1. **[Happy Path]** GIVEN a new project initialization WHEN `init` runs THEN the `.specforce/memorial/` directory is created and `AGENTS.md` is updated to point to the new structure.

**Technical Constraints (NFR):**
- **[Security]:** Directory MUST use 0750 permissions.

### [REQ-5] Multi-language Documentation Updates
**User Story:** AS A user, I WANT the documentation to accurately reflect the new distributed memorial structure, SO THAT I understand how memory is managed in my project.

**Scenarios:**
1. **[Happy Path]** GIVEN the new memorial architecture WHEN I read the documentation (EN/PT/ES) THEN it explains the directory-based structure instead of a single file.

**Technical Constraints (NFR):**
- **[Consistency]:** All supported languages (English, Portuguese, Spanish) MUST be updated.

## 4. Business Invariants
- The "Rules of Engagement" for AI agents MUST always be present and prioritized.
- Memory fragments MUST never be deleted by the system automatically unless explicitly distilled.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Memory operations must not block the main developer loop.
- **[Reliability]:** Corrupt fragments must be ignored with a warning, not cause a system crash.
- **[Maintainability]:** Use the existing `internal/project` and `internal/spec` abstractions for file management.
