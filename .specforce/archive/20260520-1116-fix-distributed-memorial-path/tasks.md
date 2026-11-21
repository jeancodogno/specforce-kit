---
slug: 20260520-1116-fix-distributed-memorial-path
lens: Backend-heavy
---

# Implementation Roadmap: Fix Distributed Memorial Path and Generation

## 1. Execution Strategy
- **Gravity Order:** Core Registry Logic -> Artifact Configuration -> Memorial Service Refactor -> Service Integration -> Verification.
- **TDD Focus:** Each task begins with a failing test case (where applicable) ensuring the new path-overloading logic and template injection are correctly verified.

## 2. Tasks

### Phase 1: Core Registry & Config Updates

- [x] T1.1: Update `Artifact` struct for path overrides
**Target:** `src/internal/constitution/registry.go`
**Context:** [FIX-1]

**Action Steps:**
- Modify `Artifact` struct: Change `Path` field YAML tag to `yaml:"path,omitempty"`.
- Update `loadArtifact` function: Check if `art.Path` is set before applying default `.specforce/docs/` logic.

**Acceptance Check:**
`go test ./src/internal/constitution/...` (Add test case for path override).

- [x] T1.2: Configure Memorial custom path
**Target:** `src/internal/agent/artifacts/constitution/memorial.yaml`
**Context:** [FIX-1]

**Action Steps:**
- Add `path: .specforce/memorial/ROUTING.md` to the YAML frontmatter.

**Acceptance Check:**
`grep "path: .specforce/memorial/ROUTING.md" src/internal/agent/artifacts/constitution/memorial.yaml`

### Phase 2: Memorial Service Refactor

- [x] T2.1: Update `MemorialService` interface and implementation
**Target:** `src/internal/project/memorial.go`
**Context:** [FIX-2]

**Action Steps:**
- Update `Initialize` signature: `Initialize(ctx context.Context, template string) error`.
- Replace hardcoded `routingContent` with the `template` parameter.
- Add minimal fallback if `template` is empty.

**Acceptance Check:**
Compilation check (will fail until T2.2 is done).

- [x] T2.2: Update Memorial unit tests
**Target:** `src/internal/project/memorial_test.go`
**Context:** [FIX-2]

**Action Steps:**
- Update all `Initialize(ctx)` calls to `Initialize(ctx, template)`.
- Add test case verifying custom template content in `ROUTING.md`.

**Acceptance Check:**
`go test -v src/internal/project/memorial_test.go`

### Phase 3: Service Integration & Verification

- [x] T3.1: Integrate Registry template into Project Initialization
**Target:** `src/internal/project/service.go`
**Context:** [FIX-2]

**Action Steps:**
- Load `memorial` artifact from registry in `InitializeProject`.
- Pass artifact template to `memSvc.Initialize`.

**Acceptance Check:**
`go build ./src/cmd/specforce/...`

- [x] T3.2: Final Integration Verification
**Target:** `Global Scope`
**Context:** [FIX-1, FIX-2]

**Action Steps:**
- Run `specforce init` in a fresh directory.
- Verify `ROUTING.md` content and `specforce constitution status` output.

**Acceptance Check:**
Manual verification as defined in requirements.md.
