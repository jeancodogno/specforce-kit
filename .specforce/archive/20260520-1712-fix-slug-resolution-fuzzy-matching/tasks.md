---
slug: 20260520-1712-fix-slug-resolution-fuzzy-matching
lens: Backend-heavy
---

# Implementation Roadmap: Fuzzy Slug Resolution

## 1. Execution Strategy
- **Gravity Order:** Implement core logic in `slug.go` -> Update services and entry points -> Verify with unit and integration tests.

## 2. Tasks

### Phase 1: Core Resolution Logic

- [x] T1.1: [CODE] Implement ResolveSlug in `src/internal/spec/slug.go`
**Target:** `src/internal/spec/slug.go`
**Context:** [FIX-1, FIX-2, FIX-3]

**Action Steps:**
- Add `os` import.
- Implement `ResolveSlug` with exact match check.
- Add fuzzy matching for timestamped slugs in `specs/`.
- Add fuzzy matching for timestamped slugs in `archive/`.
- Ensure sub-path support (e.g., `team-a/slug`).
- Implement "Prioritize Newest" logic by sorting matches.

**Acceptance Check:**
`go test -v src/internal/spec/slug.go src/internal/spec/slug_test.go` (After T2.1)

### Phase 2: Service Integration

- [x] T2.1: [CODE] Integrate ResolveSlug into domain services
**Target:** `Global Scope`
**Context:** [FIX-1]

**Action Steps:**
- Update `GetStatus` in `src/internal/spec/status.go`.
- Update `GetImplementationStatus` in `src/internal/spec/service.go`.
- Update `UpdateTaskStatus` in `src/internal/spec/service.go`.
- Update `ArchiveSpec` in `src/internal/spec/archive.go`.

**Acceptance Check:**
Existing tests should pass and newly added resolution tests should pass.

### Phase 3: Verification & Hardening

- [x] T3.1: [TEST] Add comprehensive unit tests for ResolveSlug
**Target:** `src/internal/spec/slug_test.go`
**Context:** [FIX-1, FIX-2, FIX-3]

**Action Steps:**
- Add `TestResolveSlug` with multiple scenarios.
- Verify "Prioritize Newest" behavior.
- Verify Archive resolution.

**Acceptance Check:**
`go test -v ./src/internal/spec/...`
