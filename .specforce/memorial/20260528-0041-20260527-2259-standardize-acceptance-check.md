---
date: 2026-05-28
scope: 20260527-2259-standardize-acceptance-check
author: agent
type: Lesson
---

# Methodology-Agnostic Task Verification

Standardized the task verification label as '**Acceptance Check:**' to prevent LLM methodology confusion with TDD. Implemented a Hard Migration across the repository to enforce the new standard. Updated the Go parser to strictly validate this label while providing clear deprecation feedback for the legacy label.
