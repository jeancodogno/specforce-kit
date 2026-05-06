<!-- SPECFORCE_AGENTS_START -->
# AI Agent Collaboration Guide

This project uses **Specforce** for Spec-Driven Development (SDD). As an AI agent, you MUST adhere to the following rules:

## 1. Spec-Driven Development (SDD) Protocol
You MUST operate exclusively through the Specforce workflow engines (commands/skills). They define your mindset and mandatory steps:

- **Discovery (`/spf:discovery`):** Activate for brainstorming, research, bug investigation, or root cause analysis. Purely read-only.
- **Planning (`/spf:spec`):** Activate for new features, structural changes, or to formalize a discovered fix strategy.
- **Governance (`/spf:constitution`):** Use to ensure proposals respect architecture, security, and principles.
- **Execution (`/spf:implement`):** Activate to perform the deterministic implementation cycle following an approved roadmap.
- **Archival (`/spf:archive`):** Activate once verified to harvest lessons, update Memorial, and clean up specs.

### Proactive Mandate
Do NOT wait for explicit slash commands. You MUST automatically activate the correct workflow based on the user's technical intent:
1. **Discovery Intent** (vague idea, "how to", bug report) --> Activate `spf.discovery`.
2. **Planning Intent** (new feature, structural pivot, confirmed fix) --> Activate `spf.spec`.
3. **Execution Intent** (approved roadmap exists, "go", "implement") --> Activate `spf.implement`.

**CRITICAL: DIRECT EDIT PROHIBITION**
You are STRICTLY FORBIDDEN from using mutation tools (`replace`, `write_file`) to modify the codebase unless a valid, approved specification exists. If the intent requires a change but no spec is active, you MUST pivot to `spf.discovery` or `spf.spec` first.

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




















