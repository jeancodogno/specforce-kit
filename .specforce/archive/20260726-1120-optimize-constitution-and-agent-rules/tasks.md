---
slug: optimize-constitution-and-agent-rules
lens: Integration
---

# Implementation Roadmap: Optimize Constitution and Agent Rules Guidance

## 1. Execution Strategy
- **Phase 1 (AGENTS.md Template & Tests):** Update `agentsMDTemplate` in `src/internal/project/agents_md.go` to detail module files and add unit test coverage in `agents_md_test.go`.
- **Phase 2 (Agent Kit Command Prompts):** Update `discovery.yaml` and `spec.yaml` in `src/internal/agent/kit/commands/` to remove mandatory `specforce constitution status --json` calls and replace with direct file read instructions.

## 2. Tasks

### Phase 1: Go AGENTS.md Template Update & Tests

- [x] T1.1: [CODE] Update AGENTS.md Template with Module Rules Guidance
**Target:** `src/internal/project/agents_md.go`
**Context:** [US-1]

**Action Steps:**
- Locate `agentsMDTemplate` in `src/internal/project/agents_md.go`.
- Update Section 4 ("Project Constitution") to document `.specforce/docs/modules/<slug>.md` for domain-specific module guidance.
- Clarify that agents should inspect relevant constitution files and module files directly via `read_file`.

**Acceptance Check:**
`go build ./...` passes without compilation errors.

- [x] T1.2: [TEST] Add Unit Test for Module Guidance in AGENTS.md
**Target:** `src/internal/project/agents_md_test.go`
**Context:** [US-1]

**Action Steps:**
- Update or add a test in `agents_md_test.go` to assert that `generateAgentsContent()` contains the `modules/<slug>.md` reference.
- Run `go test ./src/internal/project/...` to ensure test passes.

**Acceptance Check:**
`go test -v ./src/internal/project/...` succeeds.

### Phase 2: Agent Kit YAML Commands Update

- [x] T2.1: [CODE] Update discovery.yaml Command Definition
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [US-2]

**Action Steps:**
- Edit `Layer 1: Constitutional Anchor` in `discovery.yaml`.
- Replace the mandatory `specforce constitution status --json` execution line with direct instruction to perform `read_file` on relevant `.specforce/docs/*.md` and `.specforce/docs/modules/<slug>.md` files.

**Acceptance Check:**
`go test ./src/internal/agent/...` succeeds and YAML parses cleanly.

- [x] T2.2: [CODE] Update spec.yaml Command Definition
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [US-2]

**Action Steps:**
- Edit `Layer 1: Constitutional Anchor` in `spec.yaml`.
- Replace the mandatory `specforce constitution status --json` execution line with direct instruction to perform `read_file` on relevant `.specforce/docs/*.md` and `.specforce/docs/modules/<slug>.md` files.

**Acceptance Check:**
`go test ./src/internal/agent/...` succeeds and YAML parses cleanly.
