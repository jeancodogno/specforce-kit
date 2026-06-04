# Tasks: Module Manifests Support

### Phase 1: CLI Infra
- [x] T1: Update ConstitutionStatus struct
**Target:** `src/internal/constitution/status.go`
**Context:** [US-1]
**Action Steps:**
- Add `Modules []string `json:"modules"`` field to the `ConstitutionStatus` struct.
- Ensure the field is correctly tagged for JSON serialization.
**Acceptance Check:**
- Run `grep '"modules"' src/internal/constitution/status.go` and verify the field exists.

- [x] T2: Implement Module scanning in GetStatus
**Target:** `src/internal/constitution/status.go`
**Context:** [US-1]
**Action Steps:**
- Modify `GetStatus` to perform a directory walk of `.specforce/docs/modules/`.
- Extract slugs (filenames without extension) from `.md` files found in the directory.
- Handle missing directory gracefully by returning an empty slice.
**Acceptance Check:**
- Code review of the `GetStatus` implementation ensuring error handling for missing directory.

- [x] T3: Unit tests for GetStatus Module discovery
**Target:** `src/internal/constitution/status_test.go`
**Context:** [US-1]
**Action Steps:**
- Add test case for empty `modules/` directory.
- Add test case for multiple `.md` files correctly mapping to slugs.
- Add test case ensuring non-`.md` files are ignored.
**Acceptance Check:**
- Run `go test -v src/internal/constitution/status_test.go` and ensure all tests pass.

### Phase 2: Agent Templates
- [x] T4: Create Module Manifest Template
**Target:** `src/internal/agent/artifacts/constitution/module.yaml`
**Context:** [US-1]
**Action Steps:**
- Define the YAML structure for module manifests including `Domain Invariants`, `Technical Patterns`, and `Success Metrics`.
- Use the structure defined in the design document (Wireframe 4.1).
**Acceptance Check:**
- Run `cat src/internal/agent/artifacts/constitution/module.yaml` and verify the YAML structure.

- [x] T5: Update Archive Instructions for Harvesting
**Target:** `src/internal/agent/kit/instructions/archive.md`
**Context:** [US-3]
**Action Steps:**
- Insert a new section for "Harvesting Module Invariants".
- Add instructions for agents to check for module affinity and propose manifest updates.
**Acceptance Check:**
- Run `grep "Harvest Module Invariants" src/internal/agent/kit/instructions/archive.md`.

- [x] T6: Add Module Affinity Prompting to Discovery
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [US-2]
**Action Steps:**
- Update "Layer 1: Constitutional Anchor" to prompt for module affinity check.
- Add the "Lazy Load" protocol instruction (`read_file` if match found).
**Acceptance Check:**
- Run `grep "affinity" src/internal/agent/kit/commands/discovery.yaml`.

- [x] T7: Add Module Affinity Prompting to Spec Orchestration
**Target:** `src/internal/agent/kit/commands/spec.yaml`
**Context:** [US-2]
**Action Steps:**
- Update "Layer 1: Constitutional Anchor" with the same affinity check as discovery.
- Ensure consistency in "Lazy Load" instructions across workflow commands.
**Acceptance Check:**
- Run `grep "affinity" src/internal/agent/kit/commands/spec.yaml`.

### Phase 3: E2E Verification
- [x] T8: Manual E2E: Status Output
**Target:** `CLI`
**Context:** [US-1]
**Action Steps:**
- Create dummy module files: `mkdir -p .specforce/docs/modules && touch .specforce/docs/modules/auth.md .specforce/docs/modules/billing.md`.
- Execute `specforce constitution status --json`.
**Acceptance Check:**
- Verify JSON output contains `"modules": ["auth", "billing"]`.

- [x] T9: Manual E2E: Archive Harvesting
**Target:** `Agent Workflow`
**Context:** [US-3]
**Action Steps:**
- Simulate a feature completion in a domain with an existing manifest (e.g., `billing`).
- Execute the `spf.archive` flow.
**Acceptance Check:**
- Verify the agent identifies affinity, reads the manifest, and proposes an update.
