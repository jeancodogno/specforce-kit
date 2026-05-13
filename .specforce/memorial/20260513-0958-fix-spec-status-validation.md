---
date: 2026-05-13 09:58
scope: 20260513-0935-fix-spec-status-validation
author: senior-orchestrator
type: Lesson
---

# Lesson: Conditional Validation for Spec Artifacts

## Context
During the initial planning of a feature, artifacts like `tasks.md` might be created before `requirements.md` or `design.md`. If the `spec status` command performs deep validation (structural checks) eagerly, it can lead to confusing "No valid Phase" or sequence errors while the user is still brainstorming.

## Decision
Implemented conditional validation in `src/internal/spec/status.go`. Deep validation is now only triggered when `foundCount == totalCount` (progress 100%).

## Lessons Learned
1. **UX vs. Rigor:** Rigorous validation is essential for implementation but can be a blocker during creative planning.
2. **Self-Healing Tests:** Adding conditional logic requires updating existing test suites that might rely on eager validation. Always ensure test mocks create the full artifact set if deep validation is the target of the test.
3. **Typo Resilience:** Missing artifacts due to typos (e.g., `requiments.md`) correctly trigger the validation suppression, making the progress report the primary feedback loop until the naming is fixed.
