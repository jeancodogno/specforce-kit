# Module: Agent Kit & Project Integration

## 1. Domain Scope
Manages the embedded Agent Kit blueprints (commands, skills), artifact generation schemas, and project workspace integration (`specforce init`, tool synchronization, and legacy asset cleanup).

## 2. Business Rules & Invariants
- `[BR-KIT-01]` The embedded kit MUST contain exclusively sovereign workflow commands (`spec`, `constitution`, `discovery`, `implement`, `archive`) and the `consultative-grill` skill. Legacy subagents and deprecated skills are decommissioned.
- `[BR-KIT-02]` `requirements.md` MUST adhere to the Strict Zero-Technical-Specification Policy: containing exclusively business rules, end-user personas, and BDD acceptance criteria ("What" and "Why"). Technical implementation details ("How") belong strictly in `design.md`.
- `[BR-KIT-03]` Project initialization and update flows (`specforce init`, `UpdateTools`) MUST detect legacy Specforce assets and prompt the user before performing any destructive cleanup.

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

## 4. Technical Contracts & Integration Points
- **Packages:**
  - `src/internal/agent`: Embedded FS loader, translation engine, and kit manifest resolution.
  - `src/internal/project`: Project service, bootstrap, `AGENTS.md` sync, and `legacy.go` asset detection/cleanup.
- **CLI Commands:**
  - `specforce init [agents...]`: Initializes or updates project tools and performs legacy cleanup checks.

## 5. Operational Invariants
- All unit and integration tests must pass cleanly (`go test ./...`).
- Embedded blueprints in `kitFS` must be valid YAML without unmapped tools.
