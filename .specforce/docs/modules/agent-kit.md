# Module: Agent Kit & Project Integration

## 1. Domain Scope
Manages the embedded Agent Kit blueprints (commands, skills), artifact generation schemas, and project workspace integration (`specforce init`, tool synchronization, and legacy asset cleanup).

## 2. Business Rules & Invariants
- `[BR-KIT-01]` The embedded kit MUST contain exclusively sovereign workflow commands (`spec`, `constitution`, `discovery`, `implement`, `archive`) and the `consultative-grill` skill. Legacy subagents and deprecated skills are decommissioned.
- `[BR-KIT-02]` `requirements.md` MUST adhere to the Strict Zero-Technical-Specification Policy: containing exclusively business rules, end-user personas, and BDD acceptance criteria ("What" and "Why"). Technical implementation details ("How") belong strictly in `design.md`.
- `[BR-KIT-03]` Project initialization and update flows (`specforce init`, `UpdateTools`) MUST detect legacy Specforce assets and prompt the user before performing any destructive cleanup.
- `[BR-KIT-04]` Secondary workers and subagents spawned during implementation are STRICTLY FORBIDDEN from modifying `tasks.md` directly. State transitions belong exclusively to the primary orchestrator via `specforce implementation update`.
- `[BR-KIT-05]` Implementation orchestration MUST scale subagent allocation dynamically (1 subagent for small/self-contained tasks, 2-3 subagents sweet spot for standard tasks, up to 4 for complex tasks) and mandate subagent usage whenever supported by the environment.
- `[BR-KIT-06]` Task status updates via `specforce implementation update` support multi-task batching (comma-separated or repeated `--task`), executing verification hooks once in deduplicated sequence and updating task states and session logs atomically.
- `[BR-KIT-07]` Implementation orchestration blueprints and Mission Brief Envelopes MUST embed non-negotiable worker guardrails (pre-coding interrogation, anti-sycophancy, scope jail, and test invariance) and an orchestrator mid/post-implementation specification gate.
- `[BR-KIT-08]` Discovery orchestration blueprints MUST guide agents to adopt a consultative stance ("Strong Opinions, Weakly Held"), proactively providing architectural opinions, trade-off comparisons with recommended picks, and edge-case mitigations.
- `[BR-KIT-09]` Implementation orchestration blueprints MUST mandate worker continuity (returning verification feedback to the same subagent session for up to 2 iterations before escalating), universal harness compatibility (explicit worker declarations across Claude Code, OpenCode, Antigravity, etc.), and adaptive model tier / reasoning effort routing scaled to batch complexity.
- `[BR-KIT-10]` Archival blueprints and instructions MUST append a conceptual "Suggested Next Steps & Follow-up Specs" section to feature closing summaries to maintain architectural continuity.
- `[BR-KIT-11]` Archival blueprints and instructions MUST mandate that `.specforce/docs/modules/<domain>.md` living specifications are structured exclusively as 5-part behavioral specifications (Domain Scope, Invariants, BDD Scenarios, Public Integration Surfaces, Operational Invariants), omitting low-level internal code file paths and opportunistically migrating legacy module formats upon merge.

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

### [US-KIT-04] Subagent Delegation Boundaries and Sizing Protocol
- **Scenario:** Executing an implementation roadmap with spawned subagents
  - **GIVEN** an active roadmap with batches of varying complexity
  - **WHEN** the orchestrator delegates tasks
  - **THEN** it sizes subagents between 1 and 4 (2-3 ideal), delegates code modification, and retains exclusive authority over `tasks.md` state transitions.

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

### [US-KIT-07] Worker Continuity, Poly-Harness Compatibility, and Adaptive Effort in Implementation
- **Scenario:** Executing task batches with worker subagents across various agent harnesses
  - **GIVEN** task batches with varying risk and complexity
  - **WHEN** delegating and verifying implementation tasks
  - **THEN** the orchestrator routes model tier and reasoning effort appropriately, declares explicit worker roles across harnesses, and loops verification errors back to the same worker session.

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


## 4. Public Integration Surfaces & Contracts
- **Public CLI Commands:**
  - `specforce init [agents...] [--non-interactive]`: Initializes or updates project agent tools and performs interactive or automated legacy asset cleanup.
  - `specforce implementation update <slug> --task <id1,id2,...> --status <in-progress|finished|pending>`: Atomically updates one or more tasks and executes deduplicated verification hooks.
  - `specforce implementation status <slug> [--json]`: Inspects current task execution progress and roadmap completion.
- **Cross-Module Contracts & Dependencies:**
  - Consumes specification metadata and task roadmaps from `spec-management`.
  - Integrates with `constitution` for living spec template reconciliation during archival.
  - Generates standardized prompt and command envelopes for supported harnesses (Claude Code, OpenCode, Antigravity, Cursor, Kimi, Kilo, Qwen).

## 5. Operational & Quality Invariants
- All agent blueprint definitions must be syntactically valid YAML.
- Tool and command generation must be deterministic across all supported agent environments.
- Multi-task status transitions must execute atomically without partial failure or corrupted task files.

