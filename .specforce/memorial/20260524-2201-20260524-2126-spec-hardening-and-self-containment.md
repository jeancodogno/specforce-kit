---
date: 2026-05-24
scope: 20260524-2126-spec-hardening-and-self-containment
author: agent
type: Lesson
---

# Self-Containment & Technical Density in Specs

Implemented a 'Self-Containment Mandate' for spec generation. Specs must now be actionable by agents with zero context history. This is enforced by hardening YAML templates (mandating edge cases, technical NFRs, and directive-based tasks) and programmatically validating 'Action Step' density in tasks.md (min 2 steps per task) within the Go source.
