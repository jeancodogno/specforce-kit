---
slug: 20260525-2251-spec-generation-optimization
lens: Backend-heavy
---

# Implementation Roadmap: Spec Generation Optimization

## 1. Execution Strategy
- **Gravity Order:** Registry Logic (Dependency Sorting) -> Service Logic (Instruction Prepending) -> Orchestrator Logic (Mission Brief) -> Skill Hardening.
- We start with the core registry to ensure that any subsequent validation or listing follows the new deterministic order.

## 2. Tasks

### Phase 1: Registry Logic (Topological Sort)

- [x] T1.1: [CODE] Implement Topological Sort in `Registry`
**Target:** `src/internal/spec/registry.go`
**Context:** [US-2]

**Action Steps:**
- Add `order []string` field to `Registry` struct.
- Implement `topologicalSort()` method using Kahn's Algorithm or DFS with cycle detection.
- Update `NewRegistry` to call `topologicalSort()` and fail if a cycle is detected.

**Acceptance Check:**
Run `go test ./src/internal/spec -run TestTopologicalSort` (Must create this test to verify sorting with A->B->C and failure on A->A).

- [x] T1.2: [CODE] Update Registry listing methods to use sorted order
**Target:** `src/internal/spec/registry.go`
**Context:** [US-2]

**Action Steps:**
- Modify `List()` to iterate over `r.order` instead of the hardcoded `order` slice.
- Modify `ListForType()` to respect the topological order while filtering by type.
- Ensure any artifacts not in the dependency graph are appended at the end.

**Acceptance Check:**
Run `specforce spec artifact --json` and verify that the order of the list follows Requirements -> Design -> Tasks.

### Phase 2: Service & Instruction Manager Logic (Prepend Rules)

- [x] T2.1: [CODE] Update `Service.GetArtifact` to prepend instructions
**Target:** `src/internal/spec/service.go`
**Context:** [US-1]

**Action Steps:**
- Modify `GetArtifact` logic to prepend the `custom` instructions string to `art.Instruction` instead of appending it.
- Ensure proper spacing (`\n\n`) between the custom rules and the base instructions.

**Acceptance Check:**
Run `go test ./src/internal/spec -run TestGetArtifactInstructions` to verify that project-specific rules appear BEFORE base instructions.

- [x] T2.2: [CODE] Refactor `InstructionManager` for prepending
**Target:** `src/internal/agent/instructions.go`
**Context:** [US-1]

**Action Steps:**
- Update `GetInstructions` to merge config instructions at the TOP of the content.
- Update `InjectVariables` to ensure it works correctly with the prepended content.

**Acceptance Check:**
Run `go test ./src/internal/agent -run TestInstructionMerging`.

### Phase 3: Orchestrator Logic (Mission Brief Prompt)

- [x] T3.1: [CODE] Update `spf.spec` skill with Mission Brief, Context Synthesis, and Skill Mandates
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [US-3]

**Action Steps:**
- Modify the `Artifact Processing` section to include a "Context Synthesis" step.
- Implement the skill mapping logic to identify the correct skill for each artifact type.
- Update the delegation prompt template to include sections for MANDATED SKILLS, PROJECT RULES, SHARED CONTEXT, and BASE INSTRUCTIONS.
- Ensure subagents are explicitly told which skills to activate for the task.

**Acceptance Check:**
Execute a mock `/spf:spec` run and inspect the constructed subagent prompt in logs or debug output.

### Phase 4: Skill Hardening & Final Verification

- [x] T4.1: [CODE] Harden TDD Mandates in Task Skill
**Target:** `src/internal/agent/kit/skills/task-atomic-decomposition/SKILL.yaml`
**Context:** [US-3]

**Action Steps:**
- Move "TDD-First Verification" rules from Best Practices to the top of the content.
- Update the `Checklist` to include a mandatory check for TDD-ready verification steps.

**Acceptance Check:**
Run `specforce spec artifact task-atomic-decomposition --json` and verify the new instruction weight.

- [x] T4.2: [CLI] Final System Integration Test
**Target:** `Global Scope`
**Context:** [US-1, US-2, US-3]

**Action Steps:**
- Add a custom instruction `tasks: ["Always use TDD approach"]` to `config.yaml`.
- Run `specforce spec status spec-generation-optimization`.
- Run `specforce spec artifact feature-tasks --json` and verify "Always use TDD approach" is at the VERY TOP.

**Acceptance Check:**
Confirm that all artifacts in `spec status` are listed in the order: requirements, design, tasks.
