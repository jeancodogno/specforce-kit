---
slug: 20260508-2213-distributed-memorial-architecture
lens: Backend-heavy
---

# Implementation Roadmap: Distributed Memorial Architecture

## 1. Execution Strategy
- **Gravity Order:** Domain Logic (Memorial Service) -> CLI Integration -> Template Updates -> Migration Script.

## 2. Tasks

### Phase 1: Memorial Service & Domain Logic

- [x] T1.1: [CODE] Create Memorial Service for distributed fragment management
**Target:** `src/internal/project/memorial.go`
**Context:** [REQ-1, REQ-2]

**Action Steps:**
- Define `Fragment` struct and `MemorialService` interface.
- Implement `Record` to write fragments as unique timestamped files in `.specforce/memorial/`.
- Implement `Consolidate` to aggregate `ROUTING.md` and the most recent N fragments into a single string.
- Ensure 0600 permissions for files and 0750 for the directory.

**Verification (TDD):**
`go test ./src/internal/project/...` - Create unit tests for fragment creation and consolidation ordering.

### Phase 2: Workflow & CLI Integration

- [x] T2.1: [CODE] Update Bootstrapper and AGENTS.md generation
**Target:** `src/internal/project/bootstrapper.go`
**Context:** [REQ-4]

**Action Steps:**
- Update `BootstrapProject` to create the `.specforce/memorial/` directory.
- Update `src/internal/project/agents_md.go` to change the `memorial.md` reference to the new distributed structure.

**Verification (TDD):**
`go test ./src/internal/project/bootstrapper_test.go` and manual `specforce init` verification.

- [x] T2.2: [CODE] Integrate Memorial Service into CLI HandleInit
**Target:** `src/internal/cli/cli.go`
**Context:** [REQ-4]

**Action Steps:**
- Update `HandleInit` to call `MemorialService.Initialize()` during project setup.

**Verification (TDD):**
Manual verification by initializing a new project and checking directory structure.

### Phase 3: Artifacts & Migration

- [x] T3.1: [CODE] Update Memorial Template for Agent adaptation
**Target:** `src/internal/agent/artifacts/constitution/memorial.yaml`
**Context:** [REQ-4]

**Action Steps:**
- Update the template to explain the new distributed structure to agents.
- Add `ROUTING.md` boilerplate to the template.

**Verification (TDD):**
Run `specforce init` and verify the content of `.specforce/memorial/ROUTING.md`.

- [x] T3.2: [CODE] Update Archival Skill instructions
**Target:** `src/internal/agent/kit/instructions/archive.md`
**Context:** [REQ-3]

**Action Steps:**
- Modify the "Knowledge Harvesting" section to instruct agents to use the new fragment recording logic.

**Verification (TDD):**
Manual verification by running an archival process and checking if a fragment is created.

- [x] T3.3: [SCAFFOLD] Implement Legacy Migration Logic
**Target:** `src/internal/project/memorial.go`
**Context:** [REQ-1]

**Action Steps:**
- Add logic to check for legacy `memorial.md` and move it to `.specforce/memorial/legacy.md` if it exists.

**Verification (TDD):**
`go test ./src/internal/project/memorial_test.go` with a test case for migration from a monolithic file.

### Phase 4: Documentation Updates

- [x] T4.1: [DOCS] Update Project Constitution (Engineering & Governance)
**Target:** `Global Scope`
**Context:** [REQ-5]

**Action Steps:**
- Replace references to `.specforce/docs/memorial.md` with `.specforce/memorial/` in `.specforce/docs/engineering.md` and `.specforce/docs/governance.md`.

**Verification (TDD):**
Manual inspection of the updated Markdown files.

- [x] T4.2: [DOCS] Update Multi-language Artifacts and Getting Started guides
**Target:** `Global Scope`
**Context:** [REQ-5]

**Action Steps:**
- Update `docs/{en,pt,es}/artifacts.md` to describe the new directory-based memorial.
- Update `docs/{en,pt,es}/getting-started.md` to reflect the new archival/memory harvesting process.

**Verification (TDD):**
Manual inspection of the updated Markdown files across all languages.
