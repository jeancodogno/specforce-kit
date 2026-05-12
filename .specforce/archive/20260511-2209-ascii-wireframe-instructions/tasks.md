---
slug: 20260511-2209-ascii-wireframe-instructions
lens: UI-heavy
---

# Implementation Roadmap: ASCII Wireframe Instructions for Agent Skills and Artifacts

## 1. Execution Strategy
Re-implement updates to kit commands, agent definitions, and the design artifact template to ensure UI-neutral language.

## 2. Tasks

### Phase 1: Update Kit Instructions (Agents & Commands)

- [x] T1.1: [SCAFFOLD] Update discovery.yaml (General UI)
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [US-1]

**Action Steps:**
- Update the "Interface Wireframing (ASCII Layouts)" section to mention all UI types (Web, TUI, etc.).
- Ensure the example remains ASCII-based.

**Verification (TDD):**
`cat src/internal/agent/kit/commands/discovery.yaml` and verify the generalized language.

- [x] T1.2: [SCAFFOLD] Update spec.yaml (General UI)
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [US-4]

**Action Steps:**
- Generalize the task verification step for any UI layout.

**Verification (TDD):**
`cat src/internal/agent/kit/commands/spec.yaml` and verify the content update.

- [x] T1.3: [SCAFFOLD] Update technical-solution-architect.yaml (General UI)
**Target:** `src/internal/agent/kit/agents/technical-solution-architect.yaml`
**Context:** [US-2]

**Action Steps:**
- Update section `4. Surface Blueprint` to mandate ASCII wireframes for all UI types.

**Verification (TDD):**
`cat src/internal/agent/kit/agents/technical-solution-architect.yaml` and verify the mandate.

### Phase 2: Update Artifact Template (General UI)

- [x] T2.1: [SCAFFOLD] Update design.yaml artifact template
**Target:** `src/internal/agent/artifacts/spec/design.yaml`
**Context:** [US-3]

**Action Steps:**
- Generalize the ASCII wireframe mandate to include all UI surfaces.

**Verification (TDD):**
`cat src/internal/agent/artifacts/spec/design.yaml` and verify the instruction line.

### Phase 3: Verification

- [x] T3.1: [CLI] Verify Agent Instruction Loading
**Target:** `Global Scope`
**Context:** [US-1, US-2, US-3, US-4]

**Action Steps:**
- Run the agent test suite to ensure YAML parsing still works.

**Verification (TDD):**
Run `go test ./src/internal/agent/...`
