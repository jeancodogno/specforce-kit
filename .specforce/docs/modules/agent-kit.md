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

## 4. Technical Contracts & Integration Points
- **Packages:**
  - `src/internal/agent`: Embedded FS loader, translation engine, and kit manifest resolution.
  - `src/internal/project`: Project service, bootstrap, `AGENTS.md` sync, and `legacy.go` asset detection/cleanup.
  - `src/internal/spec`: Specification engine, batch task updates (`UpdateTaskStatus`), and deduplicated hook runner.
  - `src/internal/cli`: Cobra CLI commands supporting multi-task slice parsing (`StringSliceVar`).
- **CLI Commands:**
  - `specforce init [agents...]`: Initializes or updates project tools and performs legacy cleanup checks.
  - `specforce implementation update <slug> --task <id1,id2,...> --status <status>`: Updates one or more tasks atomically.

## 5. Operational Invariants
- All unit and integration tests must pass cleanly (`go test ./...`).
- Embedded blueprints in `kitFS` must be valid YAML without unmapped tools.
