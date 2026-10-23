---
prefix: KIT
---

# Module: Agent Kit & Project Integration

## 1. Domain Scope
Manages the embedded Agent Kit blueprints (skills), artifact generation schemas, and project workspace integration (`specforce init`, tool synchronization, and legacy asset cleanup).

## 2. Business Rules & Invariants
- `[BR-KIT-01]` The embedded kit MUST contain exclusively sovereign skills (`spf-spec`, `spf-constitution`, `spf-discovery`, `spf-implement`, `spf-archive`, and `consultative-grill`). Legacy commands, workflows, and decommissioned agents/tools (including Gemini CLI) are decommissioned.
- `[BR-KIT-02]` `requirements.md` MUST adhere to the Strict Zero-Technical-Specification Policy: containing exclusively business rules, end-user personas, and BDD acceptance criteria ("What" and "Why"). Technical implementation details ("How") belong strictly in `design.md`.
- `[BR-KIT-03]` Project initialization and update flows (`specforce init`, `UpdateTools`) MUST detect legacy Specforce assets and prompt the user before performing any destructive cleanup.
- `[BR-KIT-04]` Secondary workers and subagents spawned during implementation are STRICTLY FORBIDDEN from modifying `tasks.md` directly. State transitions belong exclusively to the primary orchestrator via `specforce implementation update`.
- `[BR-KIT-05]` Implementation orchestration MUST partition the roadmap into a global budget of 2 to 3 Batches (recommended sweet spot) and at most 4 Batches for the entire implementation roadmap (or 1 batch for 1-3 task roadmaps), assigning exactly 1 dedicated worker subagent per batch. The final verification worker MUST be an independent qa-tester that did not implement any task of the roadmap.
- `[BR-KIT-06]` Task status updates via `specforce implementation update` support multi-task batching (comma-separated or repeated `--task`), executing verification hooks once in deduplicated sequence and updating task states and session logs atomically.
- `[BR-KIT-07]` Implementation orchestration blueprints and Mission Brief Envelopes MUST embed non-negotiable worker guardrails (pre-coding interrogation, anti-sycophancy, scope jail, and test invariance) and an orchestrator mid/post-implementation specification gate.
- `[BR-KIT-08]` Discovery orchestration blueprints MUST guide agents to adopt a consultative stance ("Strong Opinions, Weakly Held"), proactively providing architectural opinions, trade-off comparisons with recommended picks, and edge-case mitigations.
- `[BR-KIT-09]` Implementation orchestration blueprints MUST mandate fresh subagent worker sessions per batch (isolation) and confine message-based feedback loops strictly to within-batch verification adjustments (up to 2 iterations). Progress visibility follows a free-form post-dispatch narration standard: after dispatching a batch, the orchestrator narrates in the user's language with a progress prefix (batches and tasks completed), the batch scope and the worker role (only when relevant to the developer), and closes each batch with a brief completion summary. No literal banner or completion-tag formats are mandated.
- `[BR-KIT-10]` Archival blueprints and instructions MUST append a conceptual "Suggested Next Steps & Follow-up Specs" section to feature closing summaries to maintain architectural continuity.
- `[BR-KIT-11]` Archival blueprints and instructions MUST mandate that `.specforce/docs/modules/<domain>.md` living specifications are structured exclusively as 5-part behavioral specifications (Domain Scope, Invariants, BDD Scenarios, Public Integration Surfaces, Operational Invariants), omitting low-level internal code file paths and opportunistically migrating legacy module formats upon merge, with every module declaring a unique `prefix` in its frontmatter and expressing each rule and use case as a uniformly structured block identified as `BR-<PREFIX>-<NN>` or `US-<PREFIX>-<NN>`.
- `[BR-KIT-12]` Implementation and QA orchestration blueprints MUST enforce a 4-Tier active QA verification protocol (Tier 1: Global Test Suite, Tier 2: Active Black-Box & Smoke Verification on real compiled artifacts/APIs, Tier 3: Visual & E2E UI Verification, Tier 4: Adversarial & Edge Case Testing), prohibiting completion sign-off based exclusively on unit tests when runtime artifacts or interfaces exist. The protocol is governed by a criterion traceability matrix mapping every acceptance criterion to a verdict (`PASSED`, `FAILED`, or `NOT-VERIFIABLE`). The qa-tester is read-only with respect to the codebase, and failed criteria trigger at most 2 correction cycles before escalation.
- `[BR-KIT-13]` Worker blueprints MUST enforce multimodal visual verification for user interface tasks, requiring workers to run local dev servers/previews, capture headless screenshots across responsive states via automated browser tools (e.g., Playwright, Puppeteer, or scripts), and visually inspect rendered output against design and UI/UX standards.
- `[BR-KIT-14]` Agent blueprints MUST enforce a strict Human-in-the-Loop Blocker Protocol: prohibiting plaintext secrets in conversations and requiring `.env` updates, mandating interactive consultation tools for minor ambiguities with structured trade-offs, and enforcing mandatory specification gates (`/spf:spec`) before coding any architectural or scope drift. Workers that cannot proceed without a decision MUST stop and return a BLOCKED report stating the question, the available options and a recommendation; only the orchestrator consults the user and relays the answer.
- `[BR-KIT-15]` Named role charters (for example implementer and qa-tester) MUST be defined inside the skill blueprints themselves as instructions given at delegation time; the kit MUST NOT require registering per-tool agent definitions for them.
- `[BR-KIT-16]` Workers MUST report observations outside their task scope as out-of-scope findings without acting on them, and the orchestrator MUST surface them in the final summary. When an implementation is resumed after an interruption, the orchestrator MUST inspect the current repository and task state and continue safely without redoing finished work or discarding uncommitted changes.
- `[BR-KIT-17]` Planning blueprints MUST direct `medium`, `large` and `complex` specifications to write one delta document per affected module after requirements (and design when required) and before the task roadmap, to seal it, to make the roadmap cover every delta item, and to list the delta in the hand-off summary; when no module changes, the blueprint MUST direct declaring `delta.impact: none` instead and not sealing.
- `[BR-KIT-18]` Archival blueprints and instructions MUST direct the agent to correct each delta to the as-built result, allocate final identifiers only in the module (the delta keeps its provisional key and the mapping is recorded in the specification metadata), apply the delta so module blocks equal delta blocks, rewrite the scope, public surface and operational sections as-built where the delta changes them, stop without touching the module on a base conflict, lazily migrate a prefix-less target module after the user confirms a proposed prefix, and require `specforce spec delta check <slug>` to report `applied` (or `none`/`exempt`) before archiving.
- `[BR-KIT-19]` The independent qa-tester's traceability matrix MUST contain one additional row per delta item (operation plus identifier or provisional key) using the existing verdicts; specifications without a delta or with `delta.impact: none` add no rows, and a removal row fails when the block is still present.

## 3. Canonical Requirements & Use Cases
### [US-KIT-01] Zero Technical Specification Enforcement in Requirements
- **Scenario:** Planning a new feature or bugfix
  - **GIVEN** the orchestrator is generating `requirements.md`
  - **WHEN** drafting user stories and BDD scenarios
  - **THEN** it uses exclusively end-user personas and domain-level outcomes without technical endpoints, status codes, SQL queries, or JSON schemas.

### [US-KIT-02] Legacy Asset Detection and Interactive Cleanup
- **Scenario:** Initializing or updating an existing workspace with legacy files
  - **GIVEN** legacy agent/skill files exist in integration directories (`.agents`, `.cursor`, `.claude`, etc.)
  - **WHEN** `specforce init` executes
  - **THEN** the system prompts the user to confirm legacy asset removal and safely removes confirmed assets if approved.

### [US-KIT-03] Batch Task Updates and Deduplicated Verification
- **Scenario:** Updating multiple tasks in a single implementation batch
  - **GIVEN** multiple task IDs passed to `specforce implementation update <slug> --task T1.1,T1.2 --status finished`
  - **WHEN** verification hooks are configured
  - **THEN** the system deduplicates all applicable hooks, runs them once, and updates all task states and session logs atomically.

### [US-KIT-04] Global Batch Budget and Subagent Sizing Protocol
- **Scenario:** Executing an implementation roadmap with spawned subagents
  - **GIVEN** an active implementation roadmap of tasks
  - **WHEN** the orchestrator partitions the roadmap and delegates tasks
  - **THEN** it partitions the roadmap into 2 to 3 Batches (max 4) across the entire implementation lifecycle, assigns 1 dedicated fresh worker subagent per batch, and retains exclusive authority over `tasks.md` state transitions.

### [US-KIT-05] Strict Worker Guardrails and Orchestrator Spec Gating
- **Scenario:** Executing task implementation batches and handling ad-hoc change requests
  - **GIVEN** a worker subagent executing an implementation batch
  - **WHEN** coding and verifying tasks
  - **THEN** the worker operates under non-negotiable guardrails (pre-coding interrogation, anti-sycophancy, scope jail, and test invariance), and the orchestrator mandates updating specifications via `/spf:spec` before applying any behavioral drift.

### [US-KIT-06] Proactive Architectural Suggestions in Discovery
- **Scenario:** Brainstorming architecture or investigating bugs during discovery
  - **GIVEN** a developer discussing a technical design or issue
  - **WHEN** the discovery workflow runs
  - **THEN** the agent provides concrete opinions, suggests clean architectural patterns, recommends preferred trade-offs, and highlights edge cases.

### [US-KIT-07] Subagent Batch Isolation, Progress Visibility, and Within-Batch Continuity
- **Scenario:** Executing task batches with worker subagents across various agent harnesses
  - **GIVEN** task batches with varying risk and complexity
  - **WHEN** delegating and verifying implementation tasks
  - **THEN** the orchestrator spawns a fresh subagent for each batch, narrates each dispatched batch in free-form language (progress, scope and, when relevant, the worker role) with a brief completion summary, and restricts message feedback loops strictly to the active subagent for fixing verification failures.

### [US-KIT-08] Strategic Next Steps Handoff in Archival
- **Scenario:** Completing feature archival
  - **GIVEN** a completed feature whose living spec has been reconciled
  - **WHEN** `specforce spec archive <slug>` finishes
  - **THEN** the output includes conceptual next steps and candidate follow-up specifications.

### [US-KIT-09] Behavioral Living Spec Synthesis and Opportunistic Legacy Migration
- **Scenario:** Reconciling domain living specs during archival
  - **GIVEN** a completed feature and an existing or new domain module document
  - **WHEN** the archival lifecycle reconciles domain behavior into `.specforce/docs/modules/<domain>.md`
  - **THEN** it structures the document using the 5 canonical behavioral sections, omits internal code file dumps, and opportunistically upgrades legacy module formats to the new canonical standard.

### [US-KIT-10] Pure Skills Standard and Tool Decommissioning
- **Scenario:** Initializing or updating agent tools across supported AI coding harnesses
  - **GIVEN** an active Specforce installation for any supported agent
  - **WHEN** tool synchronization runs
  - **THEN** all capabilities are generated strictly as native skills in `<tool-target>/skills/<slug>/SKILL.md`, omitting all legacy command and workflow paths, and safely detecting/removing obsolete `.agents/workflows/`, `<tool>/commands/`, and `.gemini/` directories upon user confirmation.

### [US-KIT-11] Active Black-Box and UI Verification in QA
- **Scenario:** Conducting final quality assurance verification after implementation batches
  - **GIVEN** all batch tasks in an implementation roadmap have completed
  - **WHEN** the independent qa-tester executes the 4-tier verification protocol
  - **THEN** it builds a criterion traceability matrix, executes automated regression tests, performs active black-box execution against real compiled artifacts/APIs, verifies UI rendering and adversarial edge cases without modifying the codebase, and generates a structured QA report with a manual developer walkthrough.

### [US-KIT-12] Multimodal Screenshot Inspection in UI Tasks
- **Scenario:** Implementing or modifying user interface components
  - **GIVEN** a worker subagent assigned a task modifying visual/UI components
  - **WHEN** completing code modifications
  - **THEN** it launches the local preview server, captures screenshots across relevant responsive viewports using headless browser tools, visually inspects the captures via multimodal vision, and ensures fidelity with design specifications.

### [US-KIT-13] Zero-Secret Blocker and Spec Gate Protocol
- **Scenario:** Encountering missing credentials, technical ambiguity, or scope drift during execution
  - **GIVEN** an orchestrator or worker facing blockers during task implementation
  - **WHEN** resolving the condition
  - **THEN** it guides secret placement into `.env` without echoing values, uses interactive consultation tools to present structured trade-offs for ambiguities, and halts for formal spec updates via `/spf:spec` when scope or architecture diverges.

### [US-KIT-14] Bounded Correction Cycle After Failed QA
- **Scenario:** The qa-tester reports failed acceptance criteria
  - **GIVEN** the traceability matrix contains criteria with a failed verdict
  - **WHEN** the orchestrator handles the failures
  - **THEN** it delegates fixes to an implementer and re-runs a fresh qa-tester for at most 2 correction cycles, and escalates to the user if criteria still fail afterwards.

### [US-KIT-15] BLOCKED Escalation from Workers
- **Scenario:** A worker cannot proceed without a decision
  - **GIVEN** a worker facing an ambiguity or conflict it cannot resolve within its scope
  - **WHEN** it stops work
  - **THEN** it returns a BLOCKED report with the question, the options and a recommendation, and the orchestrator consults the user and relays the answer before work continues.

### [US-KIT-16] Delta-Driven Archive Reconciliation
- **Scenario:** Archiving a specification with a delta
  - **GIVEN** a finished `medium` or larger specification with a sealed delta
  - **WHEN** `/spf:archive` runs
  - **THEN** the agent corrects the delta to the as-built result, allocates final identifiers, applies the delta to the module, refreshes the scope, public surface and operational sections where affected, and archives only after the delta check reports applied.
- **Scenario:** Base conflict detected at archive
  - **GIVEN** a delta whose target block was changed by another specification
  - **WHEN** `/spf:archive` runs
  - **THEN** the module is left untouched and the developer is told to redo the delta and seal it again.


## 4. Public Integration Surfaces & Contracts
- **Public CLI Commands:**
  - `specforce init [agents...] [--non-interactive]`: Initializes or updates project agent tools and performs interactive or automated legacy asset cleanup.
  - `specforce implementation update <slug> --task <id1,id2,...> --status <in-progress|finished|pending>`: Atomically updates one or more tasks and executes deduplicated verification hooks.
  - `specforce implementation status <slug> [--json]`: Inspects current task execution progress and roadmap completion.
- **Cross-Module Contracts & Dependencies:**
  - Consumes specification metadata and task roadmaps from `spec-management`, including the module delta and its check states used by planning, archival and the independent QA matrix.
  - Integrates with `constitution` for living spec template reconciliation during archival.
  - Generates native skill structures for supported harnesses (Claude Code, OpenCode, Antigravity, Cursor, Kimi, Kilo, Qwen, Codex).

## 5. Operational & Quality Invariants
- All agent blueprint definitions must be syntactically valid YAML.
- Tool and skill generation must be deterministic across all supported agent environments.
- Multi-task status transitions must execute atomically without partial failure or corrupted task files.

