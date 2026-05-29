---
date: 2026-05-29
scope: 20260529-1730-discovery-command-intelligence
author: agent
type: Lesson
---

# Inline Discovery Intelligence vs Subagent Delegation

Implemented the 'Inline Persona' model for the Discovery phase to maximize token efficiency and prevent context loss. By integrating specialized Scout (Feature) and Detective (Bug) heuristics directly into the orchestrator command instead of delegating to subagents, we avoid redundant context reads and delegation overhead. This ensures discovery insights remain persistent in the primary session history for the planning phase.
