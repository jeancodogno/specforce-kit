---
slug: 20260825-1410-tiered-spec-sizing
lens: Balanced full-stack
---

# Implementation Roadmap: Tiered Spec Sizing & Dynamic Lifecycle Transitions

## 1. Execution Strategy
- **Gravity Order:** Core Metadata & Model Extensions -> Registry & Status Dynamic Sizing -> Task & Auditor Validation Updates -> CLI Commands (`init --size`, `resize`) -> Agent Kit Prompt Updates.

## 2. Tasks

### Phase 1: Core Metadata and Dynamic Registry Sizing

- [x] T1.1: [RED] Write unit tests for Metadata Size and Dynamic Artifact Matrix
**Target:** `src/internal/spec/metadata_test.go`
**Context:** [US-1]

**Action Steps:**
- Add unit tests verifying `SpecSize` constants (`small`, `medium`, `large`, `complex`) and default fallback behavior to `medium`.
- Add test assertions for `LoadMetadata` parsing `size` field from `spec.yaml` and defaulting when omitted.
- Add test assertions in `registry_test.go` for `ListForTypeAndSize` checking required artifacts per size category.

**Acceptance Check:**
- Run `go test -v ./src/internal/spec -run TestMetadataSize` and verify failures for unimplemented size fields.

- [x] T1.2: [GREEN] Implement Metadata Size definitions and Registry Dynamic Resolution
**Target:** `src/internal/spec/metadata.go`
**Context:** [US-1]

**Action Steps:**
- Define `SpecSize` type and constants `SpecSizeSmall`, `SpecSizeMedium`, `SpecSizeLarge`, `SpecSizeComplex` in `src/internal/spec/metadata.go`.
- Add `Size SpecSize` field to `Metadata` struct with JSON and YAML annotations.
- Update `LoadMetadata` to default empty size to `SpecSizeMedium`.
- Add `ListForTypeAndSize(specType string, size SpecSize) []Artifact` in `src/internal/spec/registry.go` returning conditional artifact slices.

**Acceptance Check:**
- Run `go test -v ./src/internal/spec -run "TestMetadataSize|TestRegistryListForTypeAndSize"` and confirm 100% pass.

### Phase 2: Sizing-Aware Status and Validation Rules

- [x] T2.1: [RED] Write unit tests for Sizing-Aware Status, Task Validation, and Auditor Checks
**Target:** `src/internal/spec/status_test.go`
**Context:** [US-3]

**Action Steps:**
- Add tests in `status_test.go` verifying that `GetStatus` on `small` specs only checks `tasks.md` and reports 100% progress.
- Add tests in `tasks_validation_test.go` verifying that `small` specs pass validation with a single action step and without `**Context:** [US-X]`.
- Add tests in `auditor_test.go` verifying `DeterministicCheck` skips missing `requirements.md` when spec size is `small`.

**Acceptance Check:**
- Run `go test -v ./src/internal/spec -run "TestStatusSizing|TestSmallTaskValidation|TestAuditorSmallSpec"` and confirm expected failures.

- [x] T2.2: [GREEN] Implement Sizing-Aware Status calculation, Task Validation, and Auditor Checks
**Target:** `src/internal/spec/status.go`
**Context:** [US-3]

**Action Steps:**
- Update `GetStatus` in `src/internal/spec/status.go` to invoke `registry.ListForTypeAndSize(meta.Type, meta.Size)`.
- Update `ValidateTasks` in `src/internal/spec/tasks.go` to load spec metadata and apply relaxed validation (1+ action step, optional US context) when size is `small`.
- Update `DeterministicCheck` in `src/internal/spec/auditor.go` to check metadata and skip `FILE_MISSING` for `requirements.md` if size is `small`.

**Acceptance Check:**
- Run `go test -v ./src/internal/spec` and verify all status, task validation, and auditor tests pass.

### Phase 3: CLI Commands for Initialization and Resizing

- [x] T3.1: [RED] Write CLI tests for spec init with --size and spec resize command
**Target:** `src/internal/cli/spec_test.go`
**Context:** [US-2]

**Action Steps:**
- Add CLI test for `specforce spec init <slug> --size <size>` verifying metadata creation.
- Add CLI test for `specforce spec resize <slug> --size <size>` testing promotion (`small` -> `large`) and demotion (`large` -> `small`).
- Add CLI test for invalid size rejection and JSON output format.

**Acceptance Check:**
- Run `go test -v ./src/internal/cli -run "TestSpecInitSize|TestSpecResize"` and observe failure before CLI command registration.

- [x] T3.2: [GREEN] Implement CLI Handlers and Cobra Command Registration for Resize and Init Size
**Target:** `src/internal/cli/spec.go`
**Context:** [US-2]

**Action Steps:**
- Add `HandleSpecResize(ctx context.Context, ui core.UI, slug string, newSize string, jsonMode bool) error` in `src/internal/cli/spec.go`.
- Update `HandleSpecInit` in `src/internal/cli/spec.go` to accept `size` parameter, validate against allowed sizes, and persist to metadata.
- Register `resize` sub-command in `src/internal/cli/cobra/spec.go` and add `--size` flag to both `spec init` and `spec resize`.
- Update `handleSpecCmd` dispatcher in `src/internal/cli/spec.go` to route `resize` sub-command.

**Acceptance Check:**
- Run `go test -v ./src/internal/cli` and verify end-to-end command execution.

### Phase 4: Agent Kit Orchestration Updates

- [x] T4.1: [RED] Validate Spec Agent Kit Prompt Rules and Command Definitions
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [US-4]

**Action Steps:**
- Inspect current `src/internal/agent/kit/commands/spec.yaml` to identify hardcoded 5-dimension grill requirements and lack of sizing steps.
- Add test or verification check in `src/internal/agent/` verifying kit commands load valid YAML structure.

**Acceptance Check:**
- Run `go test -v ./src/internal/agent/...` to establish baseline test status.

- [x] T4.2: [GREEN] Update Spec Agent Command with Tiered Discovery & Resizing Workflow
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [US-4]

**Action Steps:**
- Update `src/internal/agent/kit/commands/spec.yaml` Pre-Flight step to introduce Scope & Complexity Assessment (Small, Medium, Large, Complex).
- Add directives to skip exhaustive 5-dimension grill for `small` specs and trigger deep adversarial grill for `complex` specs.
- Add directives for dynamic resizing when mid-flight scope changes occur via `specforce spec resize`.

**Acceptance Check:**
- Run `specforce agent export` or test suite to ensure valid kit command compilation and output structure.
