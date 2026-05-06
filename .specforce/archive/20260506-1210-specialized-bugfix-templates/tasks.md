---
slug: 20260506-1210-specialized-bugfix-templates
lens: Backend-heavy
---

# Implementation Roadmap: Specialized Bugfix Templates

## 1. Execution Strategy
- **Gravity Order:** Metadata Foundation -> Registry Logic -> CLI Command Updates -> Status Integration -> Artifact Content.
- We start by enabling the system to understand `spec.yaml`, then update the registry to handle typed artifacts, and finally bridge the two in the CLI and Status reports.

## 2. Tasks

### Phase 1: Metadata Foundation

- [x] T1.1: [SCAFFOLD] Implement Metadata persistence
**Target:** `src/internal/spec/metadata.go`
**Context:** [REQ-1]

**Action Steps:**
- Create `Metadata` struct with `Slug`, `Name`, and `Type`.
- Implement `LoadMetadata` with backward compatibility (defaulting to `feature`).
- Implement `SaveMetadata` using `yaml.Marshal`.

**Verification (TDD):**
`go test ./src/internal/spec/...` with a test case that creates a directory, saves metadata, and loads it back.

### Phase 2: Registry Evolution

- [x] T2.1: [CODE] Update Registry for type-awareness
**Target:** `src/internal/spec/registry.go`
**Context:** [REQ-2]

**Action Steps:**
- Add `GetForType(typeName, artifactName)`: Checks for `{type}-{name}.yaml`, then falls back to `{name}.yaml`.
- Add `ListForType(typeName)`: Returns the ordered list of artifacts, prioritizing typed versions.

**Verification (TDD):**
Unit test in `registry_test.go` with a mock filesystem containing `requirements.yaml` and `bug-requirements.yaml`.

### Phase 3: Command Implementation

- [x] T3.1: [CODE] Update `spec init` with type support
**Target:** `src/internal/cli/spec.go`
**Context:** [REQ-3]

**Action Steps:**
- Add `--type` flag to the `init` command.
- Call `spec.SaveMetadata` after directory creation.

**Verification (TDD):**
`specforce spec init test-bug --type bug` should create `.specforce/specs/YYYYMMDD-HHMM-test-bug/spec.yaml` with `type: bug`.

### Phase 4: Status Integration

- [x] T4.1: [CODE] Integrate Metadata into Status loop
**Target:** `src/internal/spec/status.go`
**Context:** [REQ-3]

**Action Steps:**
- Update `GetStatus` to load metadata at the start.
- Use `registry.ListForType(meta.Type)` instead of the generic list.

**Verification (TDD):**
`specforce spec status <slug>` should show different artifact lists based on the `type` in `spec.yaml`.

### Phase 5: Artifact Injection

- [x] T5.1: [DOCS] Create Bugfix-specialized artifacts
**Target:** `src/internal/agent/artifacts/spec/bug-requirements.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Define bug-specific requirements template (RCA, Reproduction, Evidence).

**Verification (TDD):**
`specforce spec artifact bug-requirements` shows the correct new template.

- [x] T5.2: [DOCS] Create Bugfix-specialized design artifact
**Target:** `src/internal/agent/artifacts/spec/bug-design.yaml`
**Context:** [REQ-2]

**Action Steps:**
- Define bug-specific design template (Regression strategy, fix blueprint).

**Verification (TDD):**
`specforce spec artifact bug-design` shows the correct new template.

## 3. Pre-emptive Mitigations
- **Risk:** Existing specs won't have `spec.yaml`. -> **Mitigation:** `LoadMetadata` must gracefully handle missing files and return a default "feature" type.
- **Risk:** Typo in `--type` flag. -> **Mitigation:** Validate type against a whitelist ("feature", "bug") during `init`.
