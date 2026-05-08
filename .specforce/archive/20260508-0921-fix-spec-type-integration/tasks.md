---
slug: 20260508-0921-fix-spec-type-integration
lens: Backend-heavy
---

# Implementation Roadmap: Spec Type Integration Fix

## 1. Execution Strategy
Foundation-up: Core logic first (Registry/Status), then CLI exposure, and finally the Agent Kit instructions.

## 2. Tasks

### Phase 1: Core Domain (Registry & Metadata)

- [x] T1.1: [CODE] Implement Prefixed Lookup in Registry
**Target:** `src/internal/spec/registry.go`
**Context:** [REQ-4]

**Action Steps:**
- Update `Registry.Get(name string)` to check for prefix (bug-, feature-).
- If prefix exists, use `GetForType(prefix, base)`.
- Fallback to base if not found.

**Verification (TDD):**
`go test ./src/internal/spec/registry_test.go` (Add a test case for prefixed lookup).

- [x] T1.2: [CODE] Harden Metadata Defaulting
**Target:** `src/internal/spec/metadata.go`
**Context:** [REQ-3]

**Action Steps:**
- Ensure `LoadMetadata` always returns `Type: "feature"` if the field is missing or file is absent.

**Verification (TDD):**
`go test ./src/internal/spec/metadata_test.go`.

### Phase 2: Status & CLI Handlers

- [x] T2.1: [CODE] Update SpecStatus and ArtifactStatus
**Target:** `src/internal/spec/status.go`
**Context:** [REQ-2]

**Action Steps:**
- Add `Type` field to `SpecStatus`.
- Update `processArtifactStatus` to prefix `ArtifactStatus.Name` with the spec type.
- Ensure `Path` remains the physical base path.

**Verification (TDD):**
`go test ./src/internal/spec/status_test.go`.

- [x] T2.2: [CODE] Update CLI spec init and status
**Target:** `src/internal/cli/spec.go`
**Context:** [REQ-1]

**Action Steps:**
- Add type validation to `HandleSpecInit`.
- Include type in the success message.

**Verification (TDD):**
`go run src/cmd/specforce/main.go spec init test-bug --type bug --json`.

### Phase 3: Orchestration Kit

- [x] T3.1: [CONFIG] Update Orchestrator (spec.yaml)
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [REQ-5]

**Action Steps:**
- Update initialization instructions to ask for type.
- Update artifact fetch logic to use the names from status JSON.

**Verification (TDD):**
Manual inspection of the YAML content.

### Phase 4: Final Verification

- [x] T4.1: [TEST] E2E Verification
**Target:** `Global Scope`
**Context:** [REQ-1, REQ-2, REQ-4]

**Action Steps:**
- Init a bug spec.
- Check status --json (verify prefixed names).
- Run `spec artifact bug-requirements --json`.

**Verification (TDD):**
All commands return expected type-aware outputs.
