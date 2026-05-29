---
name: task-atomic-decomposition
description: Use when converting design into executable tasks. Focuses on 5-20 minute task slices, explicit dependencies, verification, rollback, and safe parallel tracks.
version: 1.0
priority: HIGH
---

# Skill: Task Atomic Decomposition & Dependency Mapping
## MANDATORY: TDD-First Verification Every task MUST define success before coding. Every task description MUST begin with or include the specific terminal command or test state that confirms success. If a machine cannot verify the task, it is considered too vague.
## When to use this skill Use when turning design into tasks, splitting large work, sequencing dependencies, or defining verify steps.
## How to use it 1. Break work into the smallest meaningful units with one clear outcome (5-20 min slices). 2. Define dependency order and verify command per task. 3. Confirm the roadmap can be executed without hidden assumptions.
## Checklist - [ ] Every task has a TDD-ready verification command or state defined. - [ ] Tasks are atomic enough to prove complete in 5-20 minutes. - [ ] Dependencies and sequential tracks are explicit.
## Artifact Guidance - **For design.md consumption**: use the class/file inventory and roadmap as the raw material for decomposition. - **For tasks.md**: you MUST produce 5-20 minute slices. Every single task MUST have a TDD-ready verification command. Placeholder verification like "Verify logic" is STRICTLY FORBIDDEN. - **For implementation handoff**: ensure every phase ends in a still-working state with passing tests.
## 1. The Micro-Atomic Standard You must break down technical designs into micro-tasks. A task is only valid if it targets a single logical change that can be verified in isolation. - Rule of One: One file, one function, or one configuration entry per task. - Verify-First: The verification command is the "Definition of Done".
## 2. The Gravity Rule (Sequencing) Order tasks to build a stable foundation: 1. Infrastructure/Schema: The "soil" for the code. 2. Contracts/Interfaces: The "promises" between layers. 3. Core Logic: The "brain" of the feature. 4. Implementation/UI: The "skin" and interaction.
## 3. Resilience & Context - Explicit Blockers: Clearly mark tasks that cannot start until another is finished. - Contextual Intelligence: Every task must carry a fragment of the "Why" (referencing [REQ-X]) so the Executor understands the business value.