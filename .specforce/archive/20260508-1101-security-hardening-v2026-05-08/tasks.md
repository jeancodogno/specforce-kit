---
slug: 20260508-1101-security-hardening-v2026-05-08
lens: Backend-heavy
---

# Implementation Roadmap: Security Hardening (v2026-05-08)

## 1. Execution Strategy
Foundation-up hardening: Secure the runtime and CI environment first, then enforce strict file permissions across the core services, and finally transition from global security scan exclusions to granular, source-level justifications.

## 2. Tasks

### Phase 1: Toolchain & Infrastructure (REQ-1)

- [x] T1.1: [INFRA] Update Go version in go.mod
**Target:** `go.mod`
**Context:** [REQ-1]

**Action Steps:**
- Update `go` directive to `1.26.3`.
- Run `go mod tidy`.

**Verification (TDD):**
`go version` and `go mod verify` confirm the toolchain update.

- [x] T1.2: [INFRA] Update CI/Release Workflows
**Target:** `.github/workflows/ci.yml`, `.github/workflows/release.yml`
**Context:** [REQ-1]

**Action Steps:**
- Update `go-version` to `1.26.3` in all workflow steps.

**Verification (TDD):**
Manual inspection of `.github/workflows/` files.

### Phase 2: File Permission Hardening (REQ-2)

- [x] T2.1: [CODE] Tighten Upgrade State Permissions
**Target:** `src/internal/upgrade/state.go`
**Context:** [REQ-2]

**Action Steps:**
- Change `os.WriteFile` permission parameter from `0644` to `0600`.
- Ensure directory creation for state uses `0750`.

**Verification (TDD):**
Test case asserting file mode is `0600` after state write.

- [x] T2.2: [CODE] Tighten Spec Metadata Permissions
**Target:** `src/internal/spec/metadata.go`
**Context:** [REQ-2]

**Action Steps:**
- Update `SaveMetadata` to use `0600` for file creation.
- Ensure `os.MkdirAll` uses `0750`.

**Verification (TDD):**
`go test ./src/internal/spec/metadata_test.go` with file permission assertions.

- [x] T2.3: [CODE] Tighten AGENTS.md Permissions
**Target:** `src/internal/project/agents_md.go`
**Context:** [REQ-2]

**Action Steps:**
- Update `os.WriteFile` for `AGENTS.md` to use `0600`.

**Verification (TDD):**
Verification of `AGENTS.md` permissions in project initialization tests.

### Phase 3: Static Analysis & Security Refinement (REQ-3)

- [x] T3.1: [CONFIG] Refactor Security Scanning Exclusions
**Target:** `Makefile`
**Context:** [REQ-3]

**Action Steps:**
- Remove `-exclude=G304,G306,G703` from the `security` target.

**Verification (TDD):**
`make security` fails (confirming scanning is active and strict).

- [x] T3.2: [CODE] Apply Granular Justifications in Upgrade Service
**Target:** `src/internal/upgrade/service.go`, `src/internal/upgrade/installer_binary.go`
**Context:** [REQ-3]

**Action Steps:**
- Add `// #nosec G304` with "Path validated by core.SecurePath" justification to `os.Open` and `os.Create` sites.

**Verification (TDD):**
`make security` passes for the modified files.

- [x] T3.3: [CODE] Apply Granular Justifications in Spec Service
**Target:** `src/internal/spec/tasks.go`, `src/internal/spec/metadata.go`
**Context:** [REQ-3]

**Action Steps:**
- Add `// #nosec G304` with justification to `os.ReadFile` and `os.Open` calls.

**Verification (TDD):**
`make security` passes for the modified files.

- [x] T3.4: [CODE] Apply Granular Justifications in Agent & Project
**Target:** `src/internal/agent/registry.go`, `src/internal/project/agents_md.go`
**Context:** [REQ-3]

**Action Steps:**
- Add `#nosec G304` in Agent Registry and `#nosec G703` in Project Service with appropriate justifications.

**Verification (TDD):**
`make security` returns zero findings across the codebase.

### Phase 4: Final Verification

- [x] T4.1: [TEST] Security & Compliance Audit
**Target:** Global
**Context:** [All REQs]

**Action Steps:**
- Run `make security`.
- Run `go mod verify`.
- Run `go test ./...`.

**Verification (TDD):**
All tests and scans pass cleanly on the new Go 1.26.3 baseline.
