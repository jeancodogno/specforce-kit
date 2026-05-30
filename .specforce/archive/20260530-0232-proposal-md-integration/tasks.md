# Implementation Roadmap: Proposal.md Integration

## 1. Execution Strategy
- **Gravity Order:** Core Logic (Go) -> Agent Instructions (YAML) -> Verification.

## 2. Tasks

### Phase 1: Core Logic (Go Implementation)

- [x] T1.1: [CODE] Update SpecStatus struct
**Target:** `src/internal/spec/status.go`
**Context:** [US-2]
**Action Steps:**
- Locate the `SpecStatus` struct definition.
- Add a new field: `ContextFiles []string `json:"context_files,omitempty"``.
- Ensure the field is exported and properly tagged for JSON serialization.

**Acceptance Check:**
- Run `go build ./...` to ensure no compilation errors.

- [x] T1.2: [CODE] Implement proposal.md detection in GetStatus
**Target:** `src/internal/spec/status.go`
**Context:** [US-2], [US-3]
**Action Steps:**
- In the `GetStatus` function, after resolving `specDir`, check for the existence of `proposal.md` using `os.Stat(filepath.Join(specDir, "proposal.md"))`.
- If the file exists, calculate its relative path from the `projectRoot` using `filepath.Rel`.
- Append the relative path to the `status.ContextFiles` slice.
- Ensure the logic correctly handles both active specs and archived specs.

**Acceptance Check:**
- Create a dummy spec with a `proposal.md` and run `go test ./src/internal/spec/...` (or a manual check with a temporary main if tests aren't available).

- [x] T1.3: [CLI] Verify Metadata Exposure (ContextFiles)
**Target:** `CLI Output`
**Context:** [US-2]
**Action Steps:**
- Initialize a test spec: `specforce spec init test-proposal --type feature`.
- Create a dummy file: `touch .specforce/specs/test-proposal/proposal.md`.
- Execute: `specforce spec status test-proposal --json`.

**Acceptance Check:**
- Assert the JSON output contains `"context_files": [".specforce/specs/test-proposal/proposal.md"]`.

### Phase 2: Agent Instructions & Protocol

- [x] T2.1: [DOCS] Update Discovery Agent Covenant & Protocol
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [US-1], [US-4], [AC 1]
**Action Steps:**
- Locate "THE NON-MUTATION COVENANT" section.
- Add an exception for `specforce spec init` and writing to `proposal.md` ONLY after user confirmation.
- Add a new section "PROPOSAL PROTOCOL" with the steps: 1. Finalize intelligence, 2. Use `ask_user` to prompt for proposal creation, 3. Execute `specforce spec init` only if confirmed, 4. Write findings to `proposal.md`.

**Acceptance Check:**
- Verify the `discovery.yaml` file content contains the strict interactive confirmation gate in the protocol section.

- [x] T2.2: [DOCS] Update Planning Agent Context Awareness
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [US-3]
**Action Steps:**
- Locate "Layer 1: Constitutional Anchor" in the Execution Protocol.
- Add an instruction to check the `context_files` metadata from the spec status for a `proposal.md` file.
- Add an instruction to read the proposal if it exists to gain context from the discovery phase.

**Acceptance Check:**
- Verify the `spec.yaml` file content includes the instructions to check `context_files`.

- [x] T2.3: [QA] End-to-End Protocol Verification
**Target:** `Integration`
**Context:** [US-1], [US-4]
**Action Steps:**
- Mock a discovery scenario.
- Test that the agent can execute `specforce spec init` successfully in a simulated discovery environment.
- Verify that writing to `.specforce/specs/<slug>/proposal.md` is permitted while other writes remain blocked.

**Acceptance Check:**
- Manual verification of agent behavior in a controlled test environment.
