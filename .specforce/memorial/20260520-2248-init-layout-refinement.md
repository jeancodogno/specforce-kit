---
date: 2026-05-20
scope: init-layout-refinement
author: agent
type: Decision
---

# Atomic Decommissioning in Init Flow

Implemented 'DecommissionAgents' in the Project Service to allow safe and atomic removal of unselected agent directories during project initialization or updates. The TUI was refactored to return both selected and to-be-removed agent IDs, enabling the CLI to orchestrate the decommissioning process before synchronizing new artifacts, preventing duplicate confirmation prompts and ensuring a clean project state.
