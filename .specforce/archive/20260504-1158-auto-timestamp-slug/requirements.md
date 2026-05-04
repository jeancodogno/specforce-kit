---
slug: 20260504-1158-auto-timestamp-slug
lens: Backend-heavy
---

# Feature: Auto-Timestamped Spec Slugs

## 1. Context & Value
Currently, when initializing a new specification with `specforce spec init`, users must manually provide the full slug, which can lead to collisions and inconsistent naming. By automatically prepending a timestamp (YYYYMMDD-HHMM) to the provided slug, we ensure unique, chronologically ordered specifications without extra manual effort.

## 2. Out of Scope (Anti-Goals)
*Explicitly state what MUST NOT be built in this specific feature to prevent scope creep.*
- Do not modify existing specifications or their folder names.
- Do not implement a feature to customize the timestamp format in this iteration.
- Do not build a UI for slug generation; this is strictly for the CLI command.
- Do not attempt to resolve collisions by incrementing the timestamp; if the slug still exists, the standard error handling for existing specs should apply.

## 3. Acceptance Criteria (BDD)
*Every requirement MUST be testable. Define success and failure.*

### [REQ-1] Automatic Timestamp Prepending
**User Story:** AS A Developer, I WANT THE system to automatically add a timestamp to my provided slug, SO THAT I don't have to manually manage chronological ordering and unique naming.

**Scenarios:**
1. **[Happy Path]** GIVEN the current date is 2026-05-04 and time is 11:58 WHEN I run `specforce spec init my-feature` THEN a new specification directory is created with the slug `20260504-1158-my-feature`.
2. **[Edge Case]** GIVEN I provide a slug that already starts with a timestamp in the `YYYYMMDD-HHMM-` format WHEN I run `specforce spec init 20260504-1158-my-feature` THEN the system MUST NOT prepend another timestamp and should use the provided slug as-is.
3. **[Ambiguity]** If the user provides only a slug that could be interpreted as a timestamp but lacks the separator (e.g., `202605041158`), the system ASSUMES it is a regular slug and PREPENDS the current timestamp anyway.

### [REQ-2] Preservation of Sub-directory structure
**User Story:** AS A Lead Architect, I WANT the timestamp to be applied to the final slug segment even if a sub-path is provided, SO THAT my nested organization structure remains intact.

**Scenarios:**
1. **[Happy Path]** GIVEN the current timestamp is `20260504-1158` WHEN I run `specforce spec init team-a/new-api` THEN the directory `.specforce/specs/team-a/20260504-1158-new-api` is created.

## 4. Business Invariants
*Strict domain rules that the system must never violate, regardless of the user story. These must be technically testable.*
- The generated timestamp MUST follow the `YYYYMMDD-HHMM` format using the system's local time.
- The timestamp MUST be joined to the slug with a single hyphen `-`.
- A slug MUST NOT contain double hyphens resulting from the timestamp prepending (e.g., if slug is `-feature`, result should be `{TS}-feature` not `{TS}--feature`).
