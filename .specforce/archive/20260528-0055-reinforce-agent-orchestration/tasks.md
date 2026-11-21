---
slug: 20260528-0055-reinforce-agent-orchestration
lens: Backend-heavy
---

# TASKS: Reinforce Agent Orchestration

## 1. MANDATED SKILLS (USE THESE SKILLS)
- task-atomic-decomposition: Senior Technical Project Manager. Specializes in atomic task decomposition, dependency sequencing, and verifiable execution roadmaps.
- tdd: Use when implementing or fixing behavior through a test-first cycle.

## 2. PROJECT RULES (MAXIMUM PRIORITY)
- DETERMINISM MANDATE: Every technical decision and implementation step MUST be concrete. Use of 'TBD', 'To be defined', or placeholders is strictly forbidden.
- Directive Logic: Action steps MUST be technical directives (e.g., "Add if guard", "Implement interface X"). Passive descriptions (e.g., "Check logic") are forbidden.
- Action Density Mandate: Every task MUST have at least 3 concrete action steps that describe the technical implementation within the target file.
- TDD-First Verification: Every task MUST include a specific terminal command or test scenario that objectively proves the task is finished.
- Strict Sequentiality: Tasks MUST be strictly sequential.
- Traceability: Every task MUST explicitly link to a requirement using the **Context:** [US-x] field.

---

### Phase 1: Planning Orchestrator Reinforcement (spec.yaml)

- [x] T1.1: Add "Coherence Gate" logic to spec.yaml template.
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [US-2]
**Action Steps:**
- Locate the `### 3. Verification & Handoff` section in `src/internal/agent/kit/commands/spec.yaml`.
- Add detailed instructions for the "Coherence Gate" check that enforces alignment between `requirements.md`, `design.md`, and `tasks.md`.
- Define the mandatory output of a `[COHERENCE_ERROR]` marker when a requirement [US-X] is found without a corresponding task.
**Acceptance Check:**
- `grep "Coherence Gate" src/internal/agent/kit/commands/spec.yaml && grep "COHERENCE_ERROR" src/internal/agent/kit/commands/spec.yaml`

### Phase 2: Implementation Mission Brief Structure (implement.yaml)

- [x] T2.1: Implement "Mission Brief" envelope in implement.yaml.
**Target:** `src/internal/agent/kit/commands/implement.yaml`
**Context:** [US-1]
**Action Steps:**
- Update the header of `src/internal/agent/kit/commands/implement.yaml` to define the MISSION BRIEF Markdown envelope.
- Add the four mandatory sections: `## 1. MANDATED SKILLS`, `## 2. PROJECT RULES (MAXIMUM PRIORITY)`, `## 3. SHARED CONTEXT`, and `--- ## 4. BASE INSTRUCTIONS`.
- Ensure the instructions mandate that `PROJECT RULES (MAXIMUM PRIORITY)` MUST always precede the `BASE INSTRUCTIONS`.
**Acceptance Check:**
- `grep -A 5 "MISSION BRIEF" src/internal/agent/kit/commands/implement.yaml`

### Phase 3: Specialized Agent Delegation (implement.yaml)

- [x] T3.1: Add role-based delegation for specforce-developer and specforce-qa.
**Target:** `src/internal/agent/kit/commands/implement.yaml`
**Context:** [US-3]
**Action Steps:**
- Update the "Task Claiming" section in `src/internal/agent/kit/commands/implement.yaml` to explicitly instruct the agent to adopt the `specforce-developer` persona for implementation tasks.
- Add a directive for the "Final QA" phase to switch the persona to `specforce-qa` for validation.
- Include instructions to log delegation transitions using the `[DELEGATION] -> {Persona}` format.
**Acceptance Check:**
- `grep -E "specforce-developer|specforce-qa|DELEGATION" src/internal/agent/kit/commands/implement.yaml`

### Phase 4: Final Validation & Integration

- [x] T4.1: Verify end-to-end template coherence and syntax.
**Target:** `src/internal/agent/kit/commands/*.yaml`
**Context:** [US-1, US-2, US-3]
**Action Steps:**
- Perform a manual verification that both YAML files are syntactically valid (proper nesting and string escaping).
- Verify that the `spec.yaml` handoff instructions correctly set the stage for the `implement.yaml` Mission Brief.
- Ensure no generic expert roles are mentioned without first prioritizing the specialized personas.
**Acceptance Check:**
- `grep -E "Coherence Gate|Mission Brief|specforce-developer" src/internal/agent/kit/commands/*.yaml`
