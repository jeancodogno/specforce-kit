<!-- SPECFORCE_AGENTS_START -->
# AI Agent Collaboration Guide

This project uses **Specforce** for Spec-Driven Development (SDD). As an AI agent, you MUST adhere to the following rules:

## 1. Spec-Driven Development (SDD) Protocol
You MUST operate exclusively through the Specforce workflow engines (commands/skills). They define your mindset and mandatory steps:

- **Planning (`/spec`):** Always activate this workflow when a new feature intent or structural change is detected. It governs requirements discovery, constitutional alignment, and task decomposition.
- **Governance (`/constitution`):** Use to ensure your proposals respect the project's architecture, security, and principles.
- **Execution (`/implement`):** Activate this engine to perform the deterministic execution cycle following the approved roadmap.
- **Archival (`/archive`):** Activate this workflow once implementation is verified to harvest lessons learned, update the Project Constitution, and clean up active specs.

### Proactive Mandate
Do NOT wait for explicit slash commands from the user. If the conversation context shifts to "planning" or "implementation," you MUST automatically invoke the corresponding workflow command/skill to proceed.

- **Specs First:** Never write implementation code until a fully approved Specification (requirements.md, design.md, tasks.md) exists.
- **Total Consistency:** If a change is required at any point (even mid-implementation), you MUST update ALL related artifacts. You are strictly forbidden from updating only tasks.md while leaving requirements.md or design.md inconsistent.
- **Atomic Execution:** Follow the exact sequence of the tasks.md roadmap. Mark tasks as [DONE] or [FINISHED] sequentially and ONLY after successful verification.

## 2. Project Constitution
Before proposing architectural changes or adding new patterns, you MUST review the relevant Constitution documents located in .specforce/docs/:
- principles.md: Core values, philosophy, and cultural/technical axioms.
- architecture.md: System boundaries, dependency direction, and persistence topology.
- ui-ux.md: Visual direction, interaction patterns, and aesthetic DNA.
- security.md: AuthZ, roles, permissions, and data protection rules.
- engineering.md: Coding standards, testing strategy, and refactoring guidelines.
- governance.md: Project lifecycle rules, ownership, and AI boundaries.
- memorial.md: Durable lessons learned and cross-session memory.

## 3. Custom Hooks Configuration
Specforce allows developers to gate state transitions (e.g., finishing a task) using custom hooks. You can configure these in the project root's config.yaml:

```yaml
# config.yaml example
hooks:
  on_task_finished:
    - "make lint"
    - "make test"
  on_phase_finished:
    - "go test ./src/internal/..."
  on_all_tasks_finished:
    - "go test ./..."
```
If a hook fails, the state transition will be blocked.

*Note: The content above is managed by Specforce. Do not edit inside these markers.*
<!-- SPECFORCE_AGENTS_END -->










