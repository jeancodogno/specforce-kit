---
slug: 20260813-2259-remove-memorial-and-enhance-modules
lens: Backend-heavy
---

# Implementation Roadmap: Remove Memorial System and Enhance Module Living Specs

## 1. Execution Strategy
- **Gravity Order:** Go Core & CLI Cleanup -> Constitution Registry & Unit Tests -> Constitution & Module Living Spec Templates -> Agent Kit Workflow Instructions -> Documentation & Configuration.

## 2. Tasks

### Phase 1: Go Core & CLI Cleanup

- [x] T1.1: [CLEANUP] Remove Memorial Service and Associated Unit Tests
**Target:** `src/internal/project/memorial.go`
**Context:** [US-1]
**Parallel With:** T1.2

**Action Steps:**
- Delete `src/internal/project/memorial.go` and `src/internal/project/memorial_test.go`.
- Remove `MemorialService` interface, `Fragment` struct, and related methods (`Record`, `Consolidate`, `Distill`).
- Verify no remaining direct references to `MemorialService` in `src/internal/project/`.

**Acceptance Check:**
`go build ./src/internal/project/...` succeeds without references to deleted memorial structures.

- [x] T1.2: [CODE] Clean CLI Archive Handlers and Cobra Commands
**Target:** `src/internal/cli/archive.go`
**Context:** [US-1]
**Parallel With:** T1.1

**Action Steps:**
- In `src/internal/cli/archive.go`, remove subcommands `memorial` and `distill` from `HandleArchive`.
- Delete `HandleArchiveMemorial` and `HandleArchiveDistill` methods.
- In `HandleArchiveInstructions`, remove `NewMemorialService` call, fragment counting, and `MEMORIAL_FRAGMENTS` context injection.
- In `printArchiveInstructions`, remove the "Archival Scopes & Command Separation" section and fragment count logs.
- In `src/internal/cli/cobra/archive.go`, remove `archiveMemorialCmd`, `archiveDistillCmd`, and their corresponding flag variables.

**Acceptance Check:**
`go test ./src/internal/cli/...` compiles and runs successfully.

- [x] T1.3: [CODE] Update Bootstrapper, Service, Registry, and Scanner
**Target:** `src/internal/project/bootstrapper.go`
**Context:** [US-1]

**Action Steps:**
- In `src/internal/project/bootstrapper.go`, remove `".specforce/memorial"` from the `dirs` slice in `BootstrapProject`.
- In `src/internal/project/service.go`, remove the memorial initialization block from `InitializeProject`.
- In `src/internal/constitution/registry.go`, remove `"memorial"` from the default ordering slice in `Registry.List`.
- In `src/internal/spec/scanner.go`, clean up the legacy `entry.Name() == "memorial.md"` exclusion.

**Acceptance Check:**
`go test ./src/internal/constitution/... ./src/internal/project/... ./src/internal/spec/...`

- [x] T1.4: [TEST] Update Unit Test Suites for Constitution and Service
**Target:** `src/internal/constitution/registry_test.go`
**Context:** [US-1]

**Action Steps:**
- In `src/internal/constitution/registry_test.go`, remove the test block asserting the `memorial` artifact path.
- In `src/internal/constitution/status_test.go`, adjust the expected core artifacts count from 7 to 6.
- In `src/internal/project/service_test.go`, remove memorial template setup and `.specforce/memorial/ROUTING.md` assertions.

**Acceptance Check:**
`go test ./src/...` passes 100% with zero test failures.

---

### Phase 2: Constitution Artifacts & Module Living Spec Enhancement

- [x] T2.1: [CLEANUP] Remove Memorial Artifact Template
**Target:** `src/internal/agent/artifacts/constitution/memorial.yaml`
**Context:** [US-1]
**Parallel With:** T2.2

**Action Steps:**
- Delete `src/internal/agent/artifacts/constitution/memorial.yaml`.
- In `src/internal/agent/artifacts/constitution/governance.yaml`, remove the `"and propose Memorial updates"` text from AI agent permissions.

**Acceptance Check:**
File `src/internal/agent/artifacts/constitution/memorial.yaml` does not exist on disk.

- [x] T2.2: [CODE] Redesign Module Manifest Template as Canonical Living Spec
**Target:** `src/internal/agent/artifacts/constitution/module.yaml`
**Context:** [US-2]
**Parallel With:** T2.1

**Action Steps:**
- Update `description` and `instruction` to establish the module manifest as an OpenSpec-style Canonical Living Spec.
- Update `template` to include structured sections:
  1. `Domain Scope`
  2. `## 1. Business Rules & Invariants`
  3. `## 2. Canonical Requirements & Use Cases` (BDD GIVEN/WHEN/THEN and Edge Cases)
  4. `## 3. Technical Contracts & Integration Points`
  5. `## 4. Operational & Quality Invariants`

**Acceptance Check:**
Verify YAML syntax and validate that the template is pure Markdown ready for rendering.

---

### Phase 3: Agent Kit & Archival Protocol Enhancement

- [x] T3.1: [CODE] Upgrade Archival Instructions for As-Built Module Living Spec Consolidation
**Target:** `src/internal/agent/kit/instructions/archive.md`
**Context:** [US-1, US-3]
**Parallel With:** T3.2

**Action Steps:**
- In Step 5 (*Knowledge Harvesting & Module Behavior Consolidation*), establish the As-Built Living Spec reconciliation protocol: require comparing original requirements with actual implemented code, tests, and user modifications.
- Remove the *Memorial Update* subsection and CLI commands (`specforce archive memorial`).
- Completely remove Step 6 (*Memory Distillation*).
- In Step 8 (*Archival Execution*), remove the "CRITICAL DUAL-STEP REQUIREMENT" warning regarding memory separation.
- In Step 9 (*Verification & Handoff*), remove `**Memorial Updated & Distilled:**` from the summary markdown template.
- In *Guardrails*, remove `Mandatory Distillation`.

**Acceptance Check:**
Verify that `archive.md` contains no references to `specforce archive memorial`, `specforce archive distill`, or `{{MEMORIAL_FRAGMENTS}}`.

- [x] T3.2: [CODE] Refactor Command Definitions and Terminology
**Target:** `src/internal/agent/kit/commands/archive.yaml`
**Context:** [US-1, US-4]
**Parallel With:** T3.1

**Action Steps:**
- In `src/internal/agent/kit/commands/archive.yaml`, remove `"distill old memory"` from the description.
- In `src/internal/agent/kit/commands/constitution.yaml`, rename "Memory Check" to "Context Reuse Check".
- In `src/internal/agent/kit/commands/implement.yaml`, rename "Memory Check" to "Context Reuse Check".
- In `src/internal/agent/kit/commands/spec.yaml`, rename "Memory Check" to "Context Reuse Check".

**Acceptance Check:**
Run `git diff src/internal/agent/kit/commands/` to confirm terminology updates across all command files.

---

### Phase 4: Documentation & Project Configuration Alignment

- [x] T4.1: [CODE] Update AGENTS.md Template and Generator
**Target:** `src/internal/project/agents_md.go`
**Context:** [US-1]
**Parallel With:** T4.2

**Action Steps:**
- In `src/internal/project/agents_md.go` (`agentsMDTemplate`), remove `memorial/: Distributed cross-session memory...` from Section 4.
- In Section 1, update Archival description to remove `update Memorial`.
- Regenerate root `AGENTS.md` using the updated template.

**Acceptance Check:**
`grep -i "memorial" AGENTS.md` returns zero matches.

- [x] T4.2: [DOCS] Clean Documentation and Core Config
**Target:** `src/internal/core/config.go`
**Context:** [US-1]
**Parallel With:** T4.1

**Action Steps:**
- In `src/internal/core/config.go`, remove the commented memorial instruction `# - "Always update the project memorial with lessons learned"`.
- Update `README.md` to remove references to agent memory.
- Update documentation files (`docs/{en,es,pt}/artifacts.md`, `cli.md`, `configuration.md`, `getting-started.md`) removing `memorial/` references.
- Remove local `.specforce/memorial/` directory.

**Acceptance Check:**
Run `go test ./...` and `git grep -i "memorial" src/` to confirm complete eradication of memorial code.
