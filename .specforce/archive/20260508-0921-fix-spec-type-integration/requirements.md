---
slug: 20260508-0921-fix-spec-type-integration
lens: Backend-heavy
---

# Feature: Spec Type Integration Fix

## 1. Context & Value
The framework currently has partial support for specification types (feature, bug), but the CLI and agent orchestration do not fully leverage it. By exposing type-prefixed artifact names in the status report, we enable AI agents to autonomously select the correct specialized templates for bugfixes or features.

## 2. Out of Scope (Anti-Goals)
- Do not add support for custom user-defined types (stick to feature/bug).
- Do not change the physical file names in .specforce/specs/ (keep requirements.md, design.md, etc.).

## 3. Acceptance Criteria (BDD)

### [REQ-1] Spec initialization with explicit type
**User Story:** AS A developer, I WANT TO specify the type of my specification during initialization, SO THAT the system creates the correct metadata.

**Scenarios:**
1. **[Happy Path]** GIVEN the user runs `specforce spec init my-bug --type bug` WHEN the initialization completes THEN the file `.specforce/specs/my-bug/spec.yaml` must contain `type: bug` AND the command output must confirm the type as "bug" with a colored success message.
2. **[Edge Case]** GIVEN an invalid type is provided (e.g. `--type chore`) WHEN executing init THEN the system must return an error and not create the directory.

**UI/UX Specifics:**
- The CLI output MUST show `[OK] Spec directory initialized: .specforce/specs/my-bug (Type: bug)`.
- The "bug" type label should be visually distinct (e.g., Yellow or Red if supported by the TUI theme).

**Technical Constraints (NFR):**
- **[Performance]:** Initialization must be instantaneous (< 50ms).
- **[Integrity]:** The `spec.yaml` must be valid YAML.

### [REQ-2] Status report with type-prefixed artifacts
**User Story:** AS AN AI agent, I WANT TO see prefixed artifact names in the status JSON, SO THAT I can directly request the correct specialized instructions.

**Scenarios:**
1. **[Happy Path]** GIVEN a spec with `type: bug` in its `spec.yaml` WHEN the user runs `specforce spec status <slug> --json` THEN the `artifacts` array in the JSON response must contain items where the `name` field is prefixed with the type (e.g., `"name": "bug-requirements"`) AND the `path` field for these artifacts must still point to the standard filename (e.g., `.specforce/specs/<slug>/requirements.md`).

**UI/UX Specifics:**
- The JSON output MUST be valid and follow the existing `ArtifactStatus` schema.
- Prefixes MUST be lowercase and separated by a hyphen.

**Technical Constraints (NFR):**
- **[Performance]:** Metadata loading must be cached or optimized to avoid repeated disk reads.
- **[Observability]:** JSON output must be pretty-printed when requested.

### [REQ-3] Backwards compatibility for legacy specs
**User Story:** AS A user with existing specs, I WANT the system to continue working without spec.yaml, SO THAT I don't have to manually update old files.

**Scenarios:**
1. **[Happy Path]** GIVEN a spec directory that does not contain a `spec.yaml` file WHEN the user runs `specforce spec status <slug> --json` THEN the system must default the spec type to `feature` AND the `artifacts` array must use `feature-` prefixed names (e.g., `"name": "feature-requirements"`).

**UI/UX Specifics:**
- The default behavior MUST be transparent to the user in the TUI mode, showing "Feature" as the default type.

**Technical Constraints (NFR):**
- **[Reliability]:** Missing `spec.yaml` must not cause errors.

### [REQ-4] Retrieval of artifact details by prefixed name
**User Story:** AS AN AI agent, I WANT TO fetch artifact details using the prefixed name, SO THAT I get the specialized template without knowing the spec type.

**Scenarios:**
1. **[Happy Path]** GIVEN the spec type is `bug` WHEN the user runs `specforce spec artifact bug-requirements --json` THEN the system must return the template and instructions associated with the `bug` requirements AND the command must succeed even if the underlying filename is `requirements.md`.

**UI/UX Specifics:**
- If the artifact name is invalid or the prefix doesn't match a known type, the system SHOULD return a helpful error message listing available artifact names.

**Technical Constraints (NFR):**
- **[Integrity]:** Registry lookup must prioritize exact matches (prefixed) before falling back to base names.

### [REQ-5] Kit Orchestration Update
**User Story:** AS A Specforce user, I WANT the orchestrator agent to autonomously use types, SO THAT bugfixes are documented with specialized RCA and reproduction steps.

**Scenarios:**
1. **[Happy Path]** GIVEN the `spf.spec` skill configuration (`src/internal/agent/kit/commands/spec.yaml`) WHEN the orchestrator is initializing a new spec THEN it must include the `--type` flag in the `specforce spec init` command AND when processing the output of `spec status --json`, it must pass the prefixed artifact names directly to the `spec artifact` command without stripping the prefix.

**UI/UX Specifics:**
- The orchestrator instruction MUST be updated to prompt for "feature" or "bug" type during initialization.

**Technical Constraints (NFR):**
- **[Observability]:** Agent logs should clearly show which type of spec is being processed.

## 4. Business Invariants
- **Deterministic Mapping:** For any given spec type T and artifact A, the prefixed name MUST be T-A.
- **Physical Immobility:** Files on disk MUST always use the base names (requirements.md, design.md, tasks.md) to ensure tool interoperability.
- **Default Resilience:** Any spec without an explicit type MUST be treated as `feature`.

## 5. Global UI/UX Contract (TUI Ghost Protocol)
- **Density Posture:** Standard (80x40).
- **Signature Moves:** Branding box header, ASCII-Braille logos.
- **State Behavior:** Use Mint Green for success labels, Red for errors.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** `spec status --json` must remain efficient, processing type lookups in O(1) time once metadata is loaded.
- **[Reliability]:** Fail-fast with clear errors if `spec.yaml` is corrupted.
- **[Security]:** 0600 file permissions for generated artifacts.
- **[Maintainability]:** 80% coverage gate for the prefix logic.
