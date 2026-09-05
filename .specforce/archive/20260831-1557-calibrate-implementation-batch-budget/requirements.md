---
slug: 20260831-1557-calibrate-implementation-batch-budget
lens: Integration
---

# Feature: Calibrate Implementation Batch Budget and Global Subagent Sizing

## 1. Context & Value
Fragmenting implementation tasks into arbitrary micro-groups causes an excessive number of subagents to be spawned, increasing latency and orchestration complexity. By establishing a holistic roadmap partitioning protocol with a target budget of 2 to 3 task batches (maximum 4) for the entire implementation lifecycle alongside a dedicated final QA specialist, the system achieves balanced throughput, predictable subagent allocation, and optimal context containment.

## 2. Out of Scope (Anti-Goals)
- Changing the within-batch feedback loop or error correction rules.
- Removing or altering the dedicated Quality Assurance Specialist subagent at the completion stage.
- Modifying CLI parsing code for task state mutations.

## 3. Acceptance Criteria (BDD)

### [US-1] Global Implementation Batch Budgeting
**User Story:** AS A developer overseeing autonomous implementation, I WANT the orchestrator to partition the implementation roadmap into a global budget of 2 to 3 batches (maximum 4), SO THAT the total number of spawned workers remains lean and proportionate to the feature complexity.

**Scenarios:**
1. **[Happy Path - Standard Feature Roadmap]** GIVEN an implementation roadmap containing 4 to 12 pending tasks across multiple phases WHEN the orchestrator maps the execution plan THEN it partitions the entire roadmap into 2 to 3 cohesive batches (sweet spot) based on phases and domain layers.
2. **[Happy Path - Small Roadmap]** GIVEN a small implementation roadmap containing 1 to 3 tasks WHEN the orchestrator maps the execution plan THEN it groups the tasks into a single batch (1 subagent).
3. **[Edge Case - Complex Multi-Domain Feature]** GIVEN a large or complex feature roadmap with multiple independent architectural domains WHEN the orchestrator maps the execution plan THEN it partitions the roadmap into no more than 4 batches in total.

**UI/UX Specifics:**
- **View/Component:** Orchestrator CLI console stream.
- **Feedback Logic:** Batch mapping summary displaying total batches within the 1-4 range.
- **Keybindings:** Non-interactive console stream.

**Technical Constraints (NFR):**
- **[Performance]:** Bounded batch creation strictly between 1 and 4 batches per feature implementation.
- **[Integrity]:** Balanced task distribution avoiding arbitrary micro-fragmentation.
- **[Observability]:** Global batch count explicitly visible from the initial delegation step.

### [US-2] Global Worker Lifecycle and Sizing Alignment
**User Story:** AS A developer, I WANT the sizing guidance to specify subagent count for the entire implementation lifecycle rather than per batch, SO THAT subagent allocation matches the global batch budget.

**Scenarios:**
1. **[Happy Path - Dedicated Worker Allocation]** GIVEN a mapped roadmap of $B$ batches ($1 \le B \le 4$) WHEN execution begins THEN exactly 1 fresh worker subagent is dedicated to each batch in sequence, followed by 1 dedicated Quality Assurance Specialist subagent for global validation.
2. **[Edge Case - Parallel Batch Tasks]** GIVEN a batch containing explicitly parallelizable tasks targeting disjoint files WHEN execution occurs THEN multiple concurrent workers may only be spawned if parallel directives are declared without file target conflicts.

**UI/UX Specifics:**
- **View/Component:** Orchestrator status output and Mission Brief logs.
- **Feedback Logic:** Clear indication of 1 worker per batch and global total worker count.
- **Keybindings:** Streamed output.

**Technical Constraints (NFR):**
- **[Performance]:** Predictable subagent process spawning limited to $B + 1$ (where $B \le 4$).
- **[Integrity]:** Clear separation between implementation worker budget and QA validation worker.
- **[Observability]:** Standardized worker sizing guidance across all blueprint mappings.

## 4. Business Invariants
- Total implementation batches for any single feature specification must not exceed 4 batches (excluding the final QA validation step).
- Standard features must target a recommended sweet spot of 2 to 3 batches across the full roadmap.
- Every mapped batch is assigned to exactly 1 dedicated implementation subagent worker session unless tasks have explicit non-conflicting parallel declarations.

## 5. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Instantaneous roadmap partitioning with zero extra LLM roundtrips.
- **[Reliability]:** Elimination of excessive subagent spawning and associated token churn.
- **[Maintainability]:** Complete alignment between `implement.yaml` blueprint, unit tests, and `.agents/skills/spf-implement/SKILL.md`.
