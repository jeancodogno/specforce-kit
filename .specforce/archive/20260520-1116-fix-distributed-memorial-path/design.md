---
slug: 20260520-1116-fix-distributed-memorial-path
lens: Bugfix
---

# Technical Design: Fix Distributed Memorial Path and Generation (Fix Blueprint)

The Distributed Memorial currently uses hardcoded paths and content during initialization, which deviates from the registry-based artifact management used by other Constitution documents. This fix aligns the Memorial with the standard registry pattern, allowing its path and template to be managed via `memorial.yaml`.

## 1. Code Path Inventory

### 1.1 `src/internal/constitution/registry.go`
- **Modify `Artifact` struct**: Update the `Path` field's YAML tag to allow overrides from the artifact definition file.
  - Change `yaml:"-"` to `yaml:"path,omitempty"`.
- **Modify `loadArtifact` function**: Update logic to only apply the default path (`.specforce/docs/<slug>.md`) if a custom `Path` was not provided in the YAML.

### 1.2 `src/internal/agent/artifacts/constitution/memorial.yaml`
- **Update Metadata**: Explicitly define the target path for the Memorial Routing document.
  - Add `path: .specforce/memorial/ROUTING.md`.

### 1.3 `src/internal/project/memorial.go`
- **Update `MemorialService` Interface**: Change the `Initialize` method signature to accept an optional template string.
- **Update `Initialize` Implementation**: 
  - Remove hardcoded `routingContent`.
  - Use the provided `template` parameter to populate `ROUTING.md`.
  - Implement a minimal fallback if the template is empty (for safety/backward compatibility).

### 1.4 `src/internal/project/service.go`
- **Update `InitializeProject`**: 
  - Instantiate a temporary `constitution.Registry` using `s.artifactsFS`.
  - Retrieve the `memorial` artifact to extract its `Template`.
  - Pass the extracted template to `memSvc.Initialize(ctx, template)`.

### 1.5 `src/internal/project/memorial_test.go`
- **Update Test Suite**: Adjust all calls to `Initialize(ctx)` to match the new signature `Initialize(ctx, template)`.
- **Add Regression Test**: Verify that providing a custom template actually results in that content being written to `ROUTING.md`.

## 2. API & Data Contracts

### 2.1 Registry Schema Changes
The `Artifact` struct (representing the YAML schema) will now recognize the `path` key:

```yaml
# Example memorial.yaml
description: "..."
instruction: "..."
path: .specforce/memorial/ROUTING.md
template: |
  # Custom Memorial Template
```

### 2.2 Memorial Initialization Contract
The service signature shifts to:
- `Initialize(ctx context.Context, template string) error`

## 3. Threat Modeling (Security-First)

- **Path Traversal**: The `Path` field in YAML must be sanitized or verified to ensure it doesn't point outside the project root. Since `Artifact` paths are currently used by internal tools that target `.specforce`, we must ensure `core.SecurePath` or similar validation is applied when these paths are used for writing.
- **Input Validation**: The `template` string is trusted as it comes from embedded artifacts, but it should still be handled as a raw byte stream to avoid encoding issues.

## 4. Regression Strategy

- **Verification 1 (Clean Init)**: Run `specforce init` in a fresh directory and verify `.specforce/memorial/ROUTING.md` is created with the content from `src/internal/agent/artifacts/constitution/memorial.yaml`.
- **Verification 2 (Path Integrity)**: Verify that other constitution files (e.g., `architecture.md`) are still correctly placed in `.specforce/docs/`.
- **Verification 3 (Unit Tests)**: Ensure `go test ./src/internal/project/...` passes with the updated signatures.

## 5. Observability & Resilience

- **Failure Fallback**: If `NewRegistry` fails during `InitializeProject`, the system should log a warning but proceed with a default minimal `ROUTING.md` to avoid blocking project setup.
- **Logging**: Added logs should use structured levels to indicate if a custom path was successfully detected and used.

## 6. Side Effects & Hazards

- **Breaking Change (Internal)**: Any external package using the `MemorialService` interface directly will break due to the signature change. In this codebase, this is limited to the `project` and `cli` packages.
- **Template Inconsistency**: If the `memorial.yaml` template is updated but a project was already initialized, the `ROUTING.md` will NOT be automatically updated (current behavior), maintaining stability for existing projects.
