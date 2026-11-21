---
slug: 20260602-1744-specforce-spec-reviewer
lens: Backend-heavy
---

# Implementation Roadmap: Specforce Spec Reviewer

## 1. Execution Strategy
- **Gravity Order:** Foundation (Metadata/Interfaces) -> Logic (Auditor/Service) -> Integration (CLI/TUI) -> Agent Kit (Reviewer/Orchestrator).

## 2. Tasks

### Phase 1: Foundation & Contracts
Establishing the core data structures and interfaces for the refinement loop.

- [x] T1.1: [CODE] Refinement Metadata
**Target:** `src/internal/spec/metadata.go`
**Context:** [US-2]
**Action Steps:**
- Add `Refinement` struct to `Metadata` with `IterationCount` (int), `LastAuditAt` (time.Time), `IsValid` (bool), and `Errors` ([]string).
- Update the `Metadata` struct definition to include the new field with json/yaml tags.
- Update `metadata_test.go` to assert that the refinement fields are correctly marshaled and unmarshaled.

**Acceptance Check:**
`go test ./src/internal/spec/metadata_test.go -v`

- [x] T1.2: [CODE] Auditor Interface
**Target:** `src/internal/spec/auditor.go`
**Context:** [US-1]
**Action Steps:**
- Define the `CoherenceError` struct with `Artifact`, `Code`, `Message`, and `Context` fields.
- Define the `Auditor` interface with `Audit` and `DeterministicCheck` methods.
- Ensure the package name is `spec` and all necessary imports are included.

**Acceptance Check:**
`go build ./src/internal/spec/auditor.go`

- [x] T1.3: [SCAFFOLD] Reviewer Agent Profile
**Target:** `src/internal/agent/kit/agents/specforce-spec-reviewer.yaml`
**Context:** [US-1]
**Action Steps:**
- Create a new YAML file for the `specforce-spec-reviewer` agent.
- Define a system prompt that mandates cross-artifact logical audits.
- Include instructions for the agent to output standardized `[COHERENCE_ERROR]` blocks.

**Acceptance Check:**
`ls src/internal/agent/kit/agents/specforce-spec-reviewer.yaml && grep "specforce-spec-reviewer" src/internal/agent/kit/agents/specforce-spec-reviewer.yaml`

- [x] T1.4: [CODE] CLI Audit Command
**Target:** `src/internal/cli/spec.go`
**Context:** [US-5]
**Action Steps:**
- Implement `handleSpecAuditCmd` in `spec.go`.
- Map flags `--error`, `--clear`, and `--iteration` to `Metadata` updates.
- Call `spec.SaveMetadata` to persist changes to `spec.yaml`.

**Acceptance Check:**
`specforce spec audit specforce-spec-reviewer --error "Test Error" && grep "Test Error" .specforce/specs/specforce-spec-reviewer/spec.yaml`

### Phase 2: Logic - The Auditor
Implementing the deterministic and AI-driven validation logic.

- [x] T2.1: [CODE] Deterministic US-X Validator
**Target:** `src/internal/spec/auditor.go`
**Context:** [US-1]
**Action Steps:**
- Implement the `DeterministicCheck` function logic to scan `requirements.md` for `[US-X]` tags.
- Cross-reference these tags against `tasks.md` to ensure every requirement has a corresponding `**Context:** [US-X]` entry.
- Implement a helper to exclude "Out of Scope" requirements from the check.

**Acceptance Check:**
`go test -v ./src/internal/spec/auditor_test.go`

- [x] T2.2: [CODE] AI Audit Integration
**Target:** `src/internal/spec/auditor.go`
**Context:** [US-1]
**Action Steps:**
- Implement the `Audit` method to invoke the `specforce-spec-reviewer` agent profile.
- Construct the audit prompt by injecting the content of requirements, design, and tasks.
- Implement a parser to extract `[COHERENCE_ERROR]` markers from the agent's response.

**Acceptance Check:**
`go test -v ./src/internal/spec/auditor_test.go` (Mocking agent output)

### Phase 3: Logic - The Refinement Loop
Orchestrating the automatic refinement process in the service layer.

- [x] T3.1: [CODE] Refinement Loop Orchestration
**Target:** `src/internal/spec/service.go`
**Context:** [US-2]
**Action Steps:**
- Add `RefineSpec(ctx context.Context, projectRoot, slug string) error` to the `Service` struct.
- Implement the loop with a hardcoded `MAX_ITERATIONS = 3` limit.
- Update `Metadata` with the `IterationCount` and `IsValid` status in `spec.yaml` on each loop pass.

**Acceptance Check:**
`go test ./src/internal/spec/service_test.go`

- [x] T3.2: [CODE] Surgical Correction Context
**Target:** `src/internal/spec/service.go`
**Context:** [US-3]
**Action Steps:**
- Implement a `buildCorrectionPayload` helper to identify the relevant requirement snippet for a `CoherenceError`.
- Ensure the prompt for the refiner agent contains the exact error and the relevant source snippet.
- Implement the logic to re-invoke the correct agent based on the artifact target.

**Acceptance Check:**
`go test ./src/internal/spec/service_test.go`

### Phase 4: Integration & UX
Exposing the refinement loop through the CLI and TUI.

- [x] T4.1: [CODE] CLI Command Integration
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [US-2]
**Action Steps:**
- Locate the Phase 3 (Tasks) orchestration and append Phase 4 (Refinement).
- Define the orchestration logic that calls `RefineSpec` after artifacts are generated.
- Ensure the `specforce-spec-reviewer` is utilized in this new phase.

**Acceptance Check:**
`grep -C 5 "Phase 4" src/internal/agent/kit/commands/spec.yaml`

- [x] T4.2: [CODE] TUI Progress & Error Display
**Target:** `src/internal/tui/spec.go`
**Context:** [US-4]
**Action Steps:**
- Update `SpecStatus` to include `RefinementCount` and `CoherenceErrors`.
- Modify the TUI rendering to display "Iteration X/3" and a list of active errors.
- Ensure a clear "Escalation" message if the loop fails after 3 attempts.

**Acceptance Check:**
`specforce spec status spec-refinement-loop` (with mock errors)
