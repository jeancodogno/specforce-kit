---
slug: 20260520-1712-fix-slug-resolution-fuzzy-matching
lens: Bugfix
---

# Bugfix: Fuzzy Slug Resolution for Timestamped Specs

## 1. Issue Description
When a user initializes a specification using `specforce spec init`, the system automatically prepends a timestamp to the directory name (e.g., `20260520-1000-my-feature`). However, subsequent commands like `specforce spec status my-feature` fail because they expect the exact directory name including the timestamp. This forces users to manually track and type full timestamped slugs, defeating the purpose of the auto-timestamp convenience.

## 2. Evidence & Observations
- **Symptom:** Command `specforce spec status <base-slug>` returns `feature directory not found`.
- **Trace:** The error originates in `src/internal/spec/status.go` at the `GetStatus` function, which joins the project root with the raw slug provided by the user.

## 3. Reproduction Steps
1. Run `specforce spec init reproduction-bug`.
2. Observe the created directory name (e.g., `.specforce/specs/20260520-1234-reproduction-bug`).
3. Run `specforce spec status reproduction-bug`.
4. **Expected:** Status report for the spec.
5. **Actual:** Error `feature directory not found: .specforce/specs/reproduction-bug`.

## 4. Root Cause Analysis (RCA)
The system currently treats the user-provided slug as a literal directory path relative to `.specforce/specs/`. It does not attempt to resolve a "base slug" (the part after the timestamp) to its actual timestamped directory name. While timestamps ensure uniqueness and ordering, the resolution logic was missing.

## 5. Acceptance Criteria (Regression Tests)

### [FIX-1] Base Slug Resolution
**Scenario: [Regression]**
GIVEN a specification directory exists named `20260520-1234-my-feature`
WHEN I run `specforce spec status my-feature`
THEN the system MUST successfully resolve the slug and display the status.

### [FIX-2] Newest Timestamp Priority
**Scenario: [Conflict Resolution]**
GIVEN two specification directories exist: `20260520-1000-conflict` and `20260521-1000-conflict`
WHEN I run `specforce spec status conflict`
THEN the system MUST resolve to the newest version (`20260521-1000-conflict`).

### [FIX-3] Archive Support
**Scenario: [Legacy Search]**
GIVEN a specification is archived in `.specforce/archive/20260515-0900-archived-feat`
WHEN I run `specforce spec status archived-feat`
THEN the system MUST resolve the slug even within the archive.

## 6. Technical Constraints (NFR)
- **[Performance]:** Slug resolution must be fast (sub-100ms) for local CLI responsiveness.
- **[Safety]:** Exact matches must always be prioritized over fuzzy matches.
- **[Observability]:** If multiple matches exist, the system should log or prioritize the most recent one.
