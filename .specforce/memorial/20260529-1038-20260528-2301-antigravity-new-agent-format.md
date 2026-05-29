---
date: 2026-05-29
scope: 20260528-2301-antigravity-new-agent-format
author: agent
type: Decision
---

# Antigravity Agent Migration & Native Discovery

Migrated agent configurations from legacy .agent/ to standard .agents/ directory. Replaced manual AGENTS.md symlinks with generated agent.json profiles to rely on native Antigravity CLI discovery. The migration is backward compatible via an atomic folder rename during project initialization and tool updates. Skill and workflow mappings remained functionally identical, specifically targeting the .agents/ directory.
