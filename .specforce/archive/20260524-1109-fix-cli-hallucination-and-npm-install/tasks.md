---
slug: 20260524-1109-fix-cli-hallucination-and-npm-install
lens: Integration
---

# Implementation Roadmap: CLI Path Hardening & NPM Recovery Guide

## 1. Execution Strategy
- **Gravity Order:** Core Constants -> Logic (Template) -> Data (Command Kits) -> Verification. We start by defining the version source of truth, then update the templates and static configurations, and finally verify everything with tests.

## 2. Tasks

### Phase 1: Core & Infrastructure

- [x] T1.1: [CODE] Centralize Version Constant
**Target:** `src/internal/core/constants.go`
**Context:** [FIX-4]

**Action Steps:**
- Add `const Version = "1.0.0-alpha.1"` to the package.

**Acceptance Check:**
- Create/Run a temporary test: `go test -v src/internal/core/constants_test.go` (if it exists) or verify with `grep`.

### Phase 2: AGENTS.md Hardening

- [x] T2.1: [CODE] Update AGENTS.md Template
**Target:** `src/internal/project/agents_md.go`
**Context:** [FIX-1, FIX-2]

**Action Steps:**
- Update `agentsMDTemplate` to include "2. CLI Execution & Environment" and "3. Environment Recovery" sections.
- Ensure the recovery command uses the centralized `Version` constant (via `fmt.Sprintf` or similar if needed, or string concatenation).

**Acceptance Check:**
- Run `go test -v src/internal/project/agents_md_test.go`.

### Phase 3: Command Kit Guardrails

- [x] T3.1: [DATA] Add Guardrails to Discovery Command
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [FIX-3]

**Action Steps:**
- Add a `Guardrails` section with the CLI execution warning.

**Acceptance Check:**
- Inspect file: `grep "Guardrails" src/internal/agent/kit/commands/discovery.yaml`.

- [x] T3.2: [DATA] Add Guardrails to Spec Command
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [FIX-3]

**Action Steps:**
- Add a `Guardrails` section.

**Acceptance Check:**
- Inspect file: `grep "Guardrails" src/internal/agent/kit/commands/spec.yaml`.

- [x] T3.3: [DATA] Add Guardrails to Implement Command
**Target:** `src/internal/agent/kit/commands/implement.yaml`
**Context:** [FIX-3]

**Action Steps:**
- Add a `Guardrails` section.

**Acceptance Check:**
- Inspect file: `grep "Guardrails" src/internal/agent/kit/commands/implement.yaml`.

- [x] T3.4: [DATA] Add Guardrails to Constitution Command
**Target:** `src/internal/agent/kit/commands/constitution.yaml`
**Context:** [FIX-3]

**Action Steps:**
- Add a `Guardrails` section.

**Acceptance Check:**
- Inspect file: `grep "Guardrails" src/internal/agent/kit/commands/constitution.yaml`.

- [x] T3.5: [DATA] Add Guardrails to Archive Command
**Target:** `src/internal/agent/kit/commands/archive.yaml`
**Context:** [FIX-3]

**Action Steps:**
- Add a `Guardrails` section.

**Acceptance Check:**
- Inspect file: `grep "Guardrails" src/internal/agent/kit/commands/archive.yaml`.

### Phase 4: Final Verification

- [x] T4.1: [CLI] Verify Full Project Refresh
**Target:** `Global Scope`
**Context:** [FIX-1, FIX-2, FIX-3, FIX-4]

**Action Steps:**
- Run `specforce refresh` (if available) or `specforce init` in a temp directory.
- Verify the generated `AGENTS.md` matches the new template.

**Acceptance Check:**
- `cat AGENTS.md | grep "Environment Recovery"`
