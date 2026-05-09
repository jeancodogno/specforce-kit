---
date: 2026-05-09
scope: 20260509-0023-replace-req-with-us-prefix
author: gemini-cli
type: Decision
---

# Decision: Standardize Feature Requirement Prefix to US (User Story)

## Context
Previously, Specforce used the generic prefix `[REQ-x]` for functional requirements in feature specifications. While functional, this didn't align perfectly with the "User Story" field already present in the templates, nor with industry-standard Agile/Scrum terminology.

## Decision
We have standardized the prefix for functional requirements in non-bugfix specifications to `[US-x]`.

## Consequences
- **Requirements Template:** `requirements.yaml` now mandates and uses `[US-x]`.
- **Tasks Template:** `tasks.yaml` now uses `**Context:** [US-X]` for linking tasks to requirements.
- **Bugfixes:** Bugfix requirements continue to use `[FIX-x]` to maintain a clear distinction between new value and repairs.
- **Traceability:** Traceability between the "User Story" field and the requirement ID is now explicit.
- **Backward Compatibility:** The Go task validator remains prefix-agnostic, ensuring archived specifications using `[REQ-x]` remain valid.
