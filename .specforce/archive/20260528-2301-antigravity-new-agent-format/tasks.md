---
slug: 20260528-2301-antigravity-new-agent-format
lens: Integration
---

# Implementation Roadmap: Antigravity New Agent Format

## 1. Execution Strategy
- Gravity Order: Constants & Kit Config -> Migration Logic -> JSON Transformer -> Project Service Integration -> Verification.

## 2. Tasks

### Phase 1: Core Configuration & Foundations

- [x] T1.1: Update Core Tool Prefixes Constants
  **Target:** `src/internal/core/constants.go`
  **Context:** [US-1]
  **Action Steps:**
  - Locate the `ToolPrefixes` map in `src/internal/core/constants.go`.
  - Change the `Antigravity` key value from `.agent/` to `.agents/`.
  - Ensure all references to the legacy path are updated if any other constants exist.
  **Acceptance Check:**
  `grep ".agents/" src/internal/core/constants.go`

- [x] T1.2: Update Antigravity kit configuration
  **Target:** `src/internal/agent/kit/kit.yaml`
  **Context:** [US-1, US-2]
  **Action Steps:**
  - Locate the `antigravity` tool definition in `kit.yaml`.
  - Update the `target` field to `.agents/`.
  - Add a new mapping for the `agents` category targeting `agent.json`.
  **Acceptance Check:**
  `cat src/internal/agent/kit/kit.yaml | grep -E "target: .agents/|agent.json"`

- [x] T1.3: Register JSON transformer for Antigravity profiles
  **Target:** `src/internal/agent/translator.go`
  **Context:** [US-2]
  **Action Steps:**
  - Open `src/internal/agent/translator.go`.
  - Add a new entry to the `Transformers` map for the `.json` extension.
  - Implement the transformer logic to marshal the blueprint into the structured `agent.json` format (name, description, instructions, tools).
  **Acceptance Check:**
  `go test ./src/internal/agent/...`

### Phase 2: Migration & Cleanup Logic

- [x] T2.1: Implement MigrateLegacyAgents logic
  **Target:** `src/internal/project/migration.go`
  **Context:** [US-1]
  **Action Steps:**
  - Create `src/internal/project/migration.go` with the package declaration.
  - Define `MigrateLegacyAgents` function to perform `os.Rename` from `.agent` to `.agents`.
  - Add logging to report the migration status to the user.
  **Acceptance Check:**
  `ls src/internal/project/migration.go && grep "func MigrateLegacyAgents" src/internal/project/migration.go`

- [x] T2.2: Implement CleanupLegacySymlinks logic
  **Target:** `src/internal/project/agents_md.go`
  **Context:** [US-3]
  **Action Steps:**
  - Open `src/internal/project/agents_md.go`.
  - Implement `CleanupLegacySymlinks(root string) error` to iterate over `.agents/*/rules/` and remove `AGENTS.md` if it's a symlink.
  - Use `os.Lstat` to verify it's a symlink before removal.
  **Acceptance Check:**
  `go test -v ./src/internal/project/...`

- [x] T2.3: Update platform configurations mapping
  **Target:** `src/internal/project/agents_md.go`
  **Context:** [US-1, US-3]
  **Action Steps:**
  - Locate `ensurePlatformConfigs` function.
  - Remove `.agent` entry from `agentMappings`.
  - Do NOT add `.agents` for Antigravity, as it now uses native discovery via `agent.json`.
  **Acceptance Check:**
  `grep -r ".agent" src/internal/project/agents_md.go | grep agentMappings`

### Phase 3: Project Service Integration

- [x] T3.1: Integrate migration in project initialization
  **Target:** `src/internal/project/service.go`
  **Context:** [US-1, US-3]
  **Action Steps:**
  - Open `src/internal/project/service.go`.
  - In `InitializeProject`, call `MigrateLegacyAgents` and `CleanupLegacySymlinks`.
  - Ensure proper error handling and logging for these steps.
  **Acceptance Check:**
  `go test ./src/internal/project/...`

- [x] T3.2: Integrate migration in tool updates
  **Target:** `src/internal/project/service.go`
  **Context:** [US-1, US-3]
  **Action Steps:**
  - In `src/internal/project/service.go`, locate `UpdateTools`.
  - Call `MigrateLegacyAgents` and `CleanupLegacySymlinks` before syncing artifacts.
  - Verify that synchronization respects the new `.agents/` directory.
  **Acceptance Check:**
  `go test ./src/internal/project/...`

### Phase 4: Final Verification

- [x] T4.1: End-to-end migration verification
  **Target:** `tasks.md`
  **Context:** [US-1, US-2, US-3]
  **Action Steps:**
  - Create a dummy project with a legacy `.agent/` directory and some symlinks.
  - Run the compiled `specforce init` command (or use `go run src/cmd/specforce/main.go init`).
  - Assert `.agent/` is gone, `.agents/` exists, and `agent.json` files are present in agent subdirectories.
  **Acceptance Check:**
  `ls -d .agents && [ ! -d .agent ] && [ -f .agents/spf.spec/agent.json ]`
