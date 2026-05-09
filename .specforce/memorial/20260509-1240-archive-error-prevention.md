---
date: 20260509-1240
scope: 20260509-1132-archive-error-prevention
author: gemini-cli
type: Action
---

# Memorial: Archive Error Prevention

## Context
During the lifecycle of the "Archive Error Prevention" feature, it was identified that the `spf.archive` workflow lacked a mandatory step to analyze challenges and bugs encountered during implementation. This gap meant that valuable lessons learned from implementation friction were not being systematically converted into global standards or Constitution updates.

## Decision & Action
- Updated the `src/internal/agent/kit/instructions/archive.md` file.
- Mandated that agents scan `tasks.md` and failure history during the "Specification Retrospective" phase.
- Added a explicit requirement for "Error Prevention" analysis in the "Constitution Impact Analysis" step.
- Enhanced the user prompt in the "Information Gathering" step to specifically mention encountered challenges when proposing Constitution updates.
- Corrected the numbering sequence in the archive instructions for better agent compliance.

## Lessons Learned
- Agents should be explicitly prompted to reflect on *how* they built something, not just *what* they built.
- Analyzing the delta between the initial plan (`tasks.md`) and the actual implementation effort is a high-signal source for identifying missing project standards.
- Clear, sequential numbering in agent instructions is critical for reliable execution of complex multi-step workflows.
