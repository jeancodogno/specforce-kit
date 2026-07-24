---
slug: 20260723-2301-consolidation-of-module-behavior-on-archival
lens: Backend-heavy
---

# Tasks / Implementation Roadmap: Consolidation of Module Behavior on Archival

## Overview
This roadmap details the steps required to update Specforce agent instructions in `src/internal/agent/kit/instructions/archive.md` and `src/internal/agent/kit/commands/archive.yaml` to mandate module behavior consolidation during feature archival.

---

### Phase 1: Archival Instruction & Protocol Update

- [x] T1.1: Update `archive.md` Instruction Protocol for Mandatory Module Behavior Consolidation
**Target:** `src/internal/agent/kit/instructions/archive.md`  
**Context:** [US-1]  
**Action Steps:**
- Revise Step 5 of `archive.md` from an optional check to a mandatory requirement for Knowledge Harvesting & Module Behavior Consolidation.
- Add detailed instructions requiring the agent to read `requirements.md` and `design.md` of the feature being archived, extract new domain rules, and merge/update `.specforce/docs/modules/<slug>.md`.
- Add clear instructions for handling cases where no prior module document exists (requesting user approval to create `.specforce/docs/modules/<slug>.md`).
- Update Step 9 (Verification & Handoff) in `archive.md` to include a mandatory line in the output summary template for `**Module Behavior Consolidated:** [Yes/No/NA - List updated module files]`.
**Acceptance Check:**
- Run `grep -i "Harvest Module Invariants" src/internal/agent/kit/instructions/archive.md && grep -i "Module Behavior Consolidated" src/internal/agent/kit/instructions/archive.md`

- [x] T1.2: Verify and Sync `archive.yaml` Command Definition
**Target:** `src/internal/agent/kit/commands/archive.yaml`  
**Context:** [US-2]  
**Action Steps:**
- Review `src/internal/agent/kit/commands/archive.yaml` to ensure its description and instructions trigger reference `specforce archive instructions` cleanly.
- Verify that command execution protocol correctly aligns with the updated `archive.md` rules without conflicts.
- Run Go unit tests for agent translator and compatibility to verify that kit artifacts build and translate properly across supported agents.
**Acceptance Check:**
- Run `go test ./src/internal/agent/... -v`
