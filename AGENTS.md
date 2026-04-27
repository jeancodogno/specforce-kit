<!-- SPECFORCE_AGENTS_START -->
# AI Agent Collaboration Guide

This project uses **Specforce** for Spec-Driven Development (SDD). As an AI agent, you MUST adhere to the following rules:

## 1. Spec-Driven Development (SDD) Protocol
- **Specs First:** Never write implementation code until a corresponding Specification (requirements.md, design.md, tasks.md) exists and is approved. You MUST consult spec artifacts using the Specforce CLI (e.g., specforce spec list, specforce spec status <slug> --json, and specforce spec artifact <name> --json) rather than reading the raw markdown files directly.
- **Atomic Tasks:** Follow the exact sequence of the tasks.md roadmap. Mark tasks as [DONE] or [FINISHED] sequentially.
- **Verification:** Execute the exact verification/TDD steps defined in the tasks.md file before marking a task as complete.

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








