---
slug: 20260827-1515-living-module-specs-refinement
lens: Integration
---

# Implementation Roadmap: Standardize Module Living Specifications as Behavioral BDD Specs

## 1. Execution Strategy
- **Gravity Order:** Kit Test Assertions -> Archival Instructions & Module Blueprint -> Existing Project Module Sanitization -> Full Test Suite Verification.

## 2. Tasks

### Phase 1: Kit Blueprints & Archival Instructions

- [x] T1.1: [RED] Add Test Assertions for Archive Instructions and Module Template
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-1, US-2]

**Action Steps:**
- Add test assertion verifying that `instructions/archive.md` explicitly contains behavioral living spec guidelines (`Domain Scope`, `Business Rules & Invariants`, `Canonical Requirements & Use Cases`, `Public Integration Surfaces`).
- Add test assertion verifying that `instructions/archive.md` forbids internal code file dumps or private struct mentions.
- Add test assertion verifying that `artifacts/constitution/module.yaml` template includes the 5 canonical behavioral sections and omits internal code file references.

**Acceptance Check:**
`go test -v ./src/internal/agent -run TestKitManifests`

- [x] T1.2: [GREEN] Update Archive Instructions for Behavioral Living Spec Synthesis
**Target:** `src/internal/agent/kit/instructions/archive.md`
**Context:** [US-1]

**Action Steps:**
- Rewrite Step 5 ("Canonical Living Spec Reconciliation") to mandate structuring `.specforce/docs/modules/<domain>.md` as a high-density Behavioral Living Spec.
- Specify the exact 5 canonical sections: Domain Scope, Business Rules & Invariants (`[BR-xx]`), Canonical Requirements & Use Cases (`[US-xx]` with BDD GIVEN/WHEN/THEN for happy and failure paths), Public Integration Surfaces, and Operational Invariants.
- Add an explicit guardrail prohibiting the inclusion of internal source code file paths, internal package listings, private struct names, or low-level implementation details in module specifications.

**Acceptance Check:**
`specforce archive instructions` outputs the updated Step 5 and guardrails.

- [x] T1.3: [GREEN] Update Constitution Module Blueprint Schema and Template
**Target:** `src/internal/agent/artifacts/constitution/module.yaml`
**Context:** [US-2]

**Action Steps:**
- Update `module.yaml` instruction field to emphasize BDD-first behavior modeling, domain invariants, and public integration surfaces.
- Refactor the `template` section in `module.yaml` to ensure clean 5-section layout with numbered `[BR-01]` rules, `[US-01]` / `[UC-01]` BDD scenarios (Happy Path and Edge Cases), and Public Integration Surfaces (CLI, APIs, Events, Dependencies).
- Remove any ambiguity regarding technical contracts by clarifying that only public interfaces and external surfaces belong in Section 3.

**Acceptance Check:**
`go test -v ./src/internal/agent` passes and confirms valid YAML manifest parsing.

### Phase 2: Canonical Project Module Documentation Sanitization

- [x] T2.1: [DOCS] Sanitize Agent Kit Living Spec to Pure Behavioral BDD Standard
**Target:** `.specforce/docs/modules/agent-kit.md`
**Context:** [US-3]

**Action Steps:**
- Review `.specforce/docs/modules/agent-kit.md` and preserve all domain invariants `[BR-KIT-xx]` and use cases `[US-KIT-xx]`.
- Replace Section 4 ("Technical Contracts & Integration Points") internal Go package paths (`src/internal/...`) with public CLI contracts (`specforce init`, `specforce implementation update`, `specforce spec`) and cross-module interfaces.
- Verify that no internal source code file paths or implementation details remain in the document.

**Acceptance Check:**
`grep -E "src/internal" .specforce/docs/modules/agent-kit.md` returns 0 matches.

- [x] T2.2: [DOCS] Sanitize Constitution Living Spec to Pure Behavioral BDD Standard
**Target:** `.specforce/docs/modules/constitution.md`
**Context:** [US-3]

**Action Steps:**
- Review `.specforce/docs/modules/constitution.md` and ensure alignment with the 5 canonical sections.
- Ensure all business invariants and BDD scenarios are intact.
- Sanitize the technical surfaces section to reference public CLI interfaces (`specforce constitution status`, `specforce constitution update`) instead of internal implementation files.

**Acceptance Check:**
`grep -E "src/internal" .specforce/docs/modules/constitution.md` returns 0 matches.

- [x] T2.3: [DOCS] Sanitize Spec Management Living Spec to Pure Behavioral BDD Standard
**Target:** `.specforce/docs/modules/spec-management.md`
**Context:** [US-3]

**Action Steps:**
- Review `.specforce/docs/modules/spec-management.md` and verify complete preservation of business rules and use cases.
- Sanitize Section 4 to list public CLI contracts (`specforce spec init`, `specforce spec status`, `specforce spec archive`) and specification lifecycle transitions without internal package listings.
- Format document to perfectly match the 5 canonical sections.

**Acceptance Check:**
`grep -E "src/internal" .specforce/docs/modules/spec-management.md` returns 0 matches.

### Phase 3: Comprehensive Verification & Synchronization

- [x] T3.1: [VERIFY] Run Test Suite and Verify Kit Artifacts & CLI Instructions
**Target:** `Global Scope`
**Context:** [US-1, US-2, US-3]

**Action Steps:**
- Execute `go test ./...` to verify all internal packages and embedded manifests pass validation.
- Execute `specforce archive instructions` to verify end-to-end output rendering of the new lifecycle instructions.
- Validate that all `.specforce/docs/modules/*.md` files conform to the canonical living specification standard.

**Acceptance Check:**
`go test ./...` exits 0 with all tests passing.

### Phase 4: Opportunistic Legacy Module Migration Protocol

- [x] T4.1: [RED] Add Test Assertions for Opportunistic Legacy Module Migration in Archive Instructions
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-4]

**Action Steps:**
- Add unit test assertion in `kit_manifests_test.go` verifying that `instructions/archive.md` mandates opportunistic migration/reformatting when a legacy module format or outdated code dump is encountered during reconciliation.
- Run `go test -v ./src/internal/agent -run TestKitManifests` and ensure assertions are evaluated.

**Acceptance Check:**
`go test -v ./src/internal/agent -run TestKitManifests`

- [x] T4.2: [GREEN] Update Archive Instructions with Opportunistic Legacy Migration Protocol
**Target:** `src/internal/agent/kit/instructions/archive.md`
**Context:** [US-4]

**Action Steps:**
- Update Step 5 ("Canonical Living Spec Reconciliation") Substep 4 ("Non-Destructive Merge & Legacy Migration") in `instructions/archive.md`.
- Explicitly state: "If an existing `.specforce/docs/modules/<domain>.md` uses a legacy format or contains internal code paths/low-level details, the agent MUST opportunistically reformat and upgrade it to the 5 canonical Living Spec sections during reconciliation while preserving all accumulated domain rules and use cases."
- Update guardrails if necessary to reinforce zero legacy debt persistence.

**Acceptance Check:**
`specforce archive instructions` displays the opportunistic legacy migration directive in Step 5.

- [x] T4.3: [VERIFY] Run Full Test Suite and Linter
**Target:** `Global Scope`
**Context:** [US-4]

**Action Steps:**
- Run `make lint` to verify zero linter warnings.
- Run `go test ./...` to verify 100% test pass rate across all packages.

**Acceptance Check:**
`make lint && go test ./...` exits 0.
