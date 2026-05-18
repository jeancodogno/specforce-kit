---
date: 2026-05-17
scope: 20260517-1929-memorial-distillation-and-efficiency
author: agent
type: Decision
---

# Activation of Memorial Playback and Smart Distillation

Standardized a two-phase memorial lifecycle. 1) Playback: The CLI now injects the actual content of the last 10 fragments into agent instructions via memSvc.Consolidate. 2) Distillation: Introduced a 'smart' distillation process where agents summarize old fragments and the CLI atomically appends the summary to DISTILLED.md and removes individual files. Also implemented automatic cleanup of the legacy monolithic memorial.md.
