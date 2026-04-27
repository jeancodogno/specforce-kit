# Engineering

## Coding Standards
- **Backend Style & Linters:** Go 1.24+ following standard `gofmt` and `golangci-lint` (Level 8 equivalent).
- **CLI Standard:** All new commands MUST be implemented using the Cobra framework (nested under `src/internal/cli/cobra`).
- **Frontend Style & Linters:** Bubbletea/Lipgloss for TUI consistency; styles encapsulated in a shared `tui` package.
- **Typing Strictness:** Rigid Go typing. Minimize use of `interface{}` or `any`; prefer explicit types and interfaces for testing mocks.
- **Embedded Registry Pattern:** Components discovered via metadata MUST be managed by a centralized `Registry` type. This registry should be initialized once (usually with an `fs.FS`) and passed as a dependency to other services or UI components.
- **Agent-Specific Artifact Mapping:** To support diverse AI coding agents with differing directory conventions, the framework MUST use centralized `kit.yaml` (kit-level) and `mapping.yaml` (legacy/blueprint-level) files. These files translate internal Specforce artifact paths to the native hidden directory names of each supported agent (e.g., `.kimi/`, `.claude/`).
- **Automated Platform Configuration:** The framework MUST automatically set up environment-specific discovery rules for supported agents during project initialization or refresh. This includes:
    - Creating `.gemini/settings.json` with the correct `fileName` context mapping.
    - Creating relative symbolic links (e.g., `.agent/rules/AGENTS.md -> ../../AGENTS.md`) to expose global project rules to agents that use standardized rules directories.
    - Ensuring all symlinks are relative to maintain project portability.
- **Generic Wildcard Expansion:** Mapping configurations in `kit.yaml` MUST support the `*` wildcard in both `path` and `name` fields. The system MUST replace the `*` with the artifact's slug (source filename without extension).
    - **Strict Category Whitelisting:** Artifacts are only installed for an agent if their source category (e.g., `agents`, `commands`, `skills`) is explicitly defined in that agent's `mappings` within `kit.yaml`. Categories not present in the mapping MUST be ignored by default.
    - **Slug Defaulting:** If a mapping's `name` field is empty, the system MUST default to the artifact's slug.
    - **Hierarchy Preservation:** For static mappings (no `*` in path), the system MUST preserve any sub-directory hierarchy found in the source category.
    - **Agent Tool Addition Process:** To integrate a new agent tool into the Specforce ecosystem, developers MUST follow these steps:
        1. **Register Constant:** Add the tool's hidden directory name (e.g., `.newagent/`) to the `ToolPrefixes` slice in `src/internal/core/constants.go`.
        2. **Define Mapping:** Update `src/internal/agent/kit/kit.yaml` with the tool's `name`, `description`, `target` directory, and `mappings` for required categories (e.g., `commands`, `skills`).
        3. **Granular Pathing:** Utilize the mapping-level `target` override and environment variable expansion (e.g., `${VAR:-DEFAULT}`) to support global path requirements when necessary.
        4. **Enable Global Writing (Optional):** If the agent requires global system access (outside the project root), add its slug to the `globalEnabledAgents` allow-list in `src/internal/agent/translator.go`.
    - **Agent Command Flatness:** For OpenCode and KiloCode, commands MUST be stored directly in a flat `commands/` subfolder using the `spf.*` prefix to ensure correct discovery by those agents.
- **Generated Markdown Files:** Any automatically generated documentation files (e.g., `AGENTS.md`) MUST use explicit start and end markers (e.g., `<!-- SPECFORCE_AGENTS_START -->` and `<!-- SPECFORCE_AGENTS_END -->`) to delimit the framework-managed content. The update logic MUST preserve any user-written content located outside these markers.
- **Service Layer Pattern:** Domain services MUST be initialized with their dependencies (filesystems, registries, etc.) and expose UI-agnostic methods. They MUST depend on the `core.UI` interface for reporting progress and logging, ensuring they remain decoupled from specific frontend implementations.
- **Instruction Injection Pattern:** Services responsible for generating or displaying agent-facing content (e.g., Spec artifacts, Implementation reports) MUST merge global instructions from `core.ProjectConfig` into their output. This ensures that the orchestration loop respects project-specific constraints without modifying the core Specforce protocol.
- **Hook-Based Verification Gating:** Critical state transitions (e.g., marking a task as `FINISHED`) MUST be gated by project-defined external hooks. If a hook fails (non-zero exit code), the transition MUST be blocked, the original state MUST NOT be modified, and the CLI command MUST exit with a non-zero status code to support shell command chaining.
- **Interactive Confirmation Gating:** Irreversible or high-impact operations (e.g., overwriting existing tool instructions or deleting specs) MUST be gated by an explicit `Confirm` prompt via the `core.UI` interface. The system MUST NOT proceed without affirmative user consent (Sim/Yes).
- **DTO Pattern (Data Transfer Objects):** Services SHOULD use explicit structs (DTOs) for complex input configurations or output reports. This maintains a stable contract between the CLI/TUI layer and the domain logic.
- **Complexity Limits:** Maximum function length: 50 lines. Maximum cyclomatic complexity: 10 per function.
- **Refactoring Patterns:** 
    - **Surgical Helper Extraction:** When a function exceeds complexity or length limits, extract logic into private helper functions within the same file to maintain readability without changing the public contract.
    - **Explicit Error Ignoring (Tests):** In test files, if an error from a function like `os.RemoveAll` or `os.Chdir` is safe to ignore, use `_ =` to make the omission explicit and satisfy linting requirements.
- **Permission Hardening:** All I/O operations MUST use the most restrictive permissions possible:
    - **Directories:** Use `0750` for all directory creation (`os.MkdirAll`).
    - **Files:** Use `0600` for all sensitive or system-internal files (e.g., config, tasks, implementation reports). Use `0644` only where public read access by external tools is strictly required.

- **Instruction-Driven Agent Pattern:** Agent command definitions (`.yaml`) SHOULD NOT contain complex, hardcoded logic or multi-step instructions. Instead, they MUST trigger a dedicated Specforce CLI command (e.g., `specforce archive instructions`) to retrieve a dynamic instruction set. This allows the framework to inject global context (Constitution), core kit rules, and project-specific overrides from `config.yaml` without updating the agent's definition.
- **Knowledge-First Archival:** The feature archival process MUST include a "Knowledge Harvesting" phase. Before a specification is archived, the agent or developer MUST update the project's `.specforce/docs/memorial.md` with lessons learned, established precedents, and critical architectural decisions discovered during the implementation.

## Agent Orchestration Protocol
1. **Mandatory Discovery:** Agents MUST call `specforce spec list` and `specforce constitution status --json` before initializing a new specification to ensure they have the full project context.
2. **Interrogation over Hallucination:** Agents MUST NOT guess user intent. They must use interrogation tools to ask clarifying questions until a clear feature idea is reached.
3. **Primary Orchestration Only:** Agents MUST remain in their primary orchestration loop and NOT use specialized "plan" or "design" modes for requirement discovery.
4. **Mandatory Skill Header Injection:** For agents using "Skills" (e.g., Kimi Code), all generated `SKILL.md` files (from skills or commands) MUST include a YAML frontmatter header with `name` and `description` to ensure native discovery.
5. **Binary-Path Standard:** All agent-triggered CLI commands MUST use the global binary `specforce` (not `./specforce`) to ensure cross-platform consistency.

## Context Handling
- **Mandatory Propagation:** All functions performing I/O, recursion, or loops MUST accept `context.Context` as their first argument.
- **Strict Separation:** `context.Context` MUST NOT be stored in a struct. It must be passed explicitly through the call stack to maintain visibility of the operation's lifecycle. 
    - *Exception:* UI models (e.g., Bubbletea models) may store a context for managing background "tick" or "refresh" operations, provided it is derived from the root lifecycle.
- **Cancellation Checks:** Any function that iterates over a collection or performs recursive calls MUST check for context cancellation (e.g., `if err := ctx.Err(); err != nil { return err }`) at the start of each iteration or recursion level.

## Testing Strategy
- **Frameworks:** Standard Go `testing` package; optional use of `stretchr/testify` for assertions.
- **Coverage Gate:** 80% minimum coverage for core domain packages (`src/internal/core`, `src/internal/project`, `src/internal/spec`).
- **Mutation Testing:** Projects SHOULD use `Gremlins` for mutation testing to verify test suite efficacy. Core packages MUST target a mutation score > 0 (all survivors investigated).
- **Unit Testing Rules:** All core logic and state transitions MUST have automated unit tests.
- **Regression Rules:** Any reported bug MUST be reproduced with a failing test case before being fixed.

## Delivery Standards & DoD
- **Commit Standard:** Conventional Commits (e.g., `feat:`, `fix:`, `docs:`, `chore:`).
- **Continuous Integration:** PRs MUST trigger a multi-job parallel pipeline (Lint, Security, Test, Quality) for rapid feedback.
- **Security Scans:** Security checks MUST use `gosec` and `govulncheck`. Findings MUST be reported via SARIF to the GitHub Security tab.
- **Definition of Done (DoD):**
  - [ ] Code passes `golangci-lint` and `go test`.
  - [ ] Coverage gate is met or maintained.
  - [ ] Security scans pass (zero critical findings).
  - [ ] TUI views are manually verified for layout consistency.
  - [ ] Feature Specification status is updated to "Implemented".
  - [ ] The Memorial is updated with relevant "Last Actions".

## Error Handling Standards
- **Sentinel Errors:** All domain-specific errors MUST be defined as package-level sentinel variables using `errors.New()` in `src/internal/core/errors.go`. No ad-hoc error strings in business logic.
- **Error Wrapping:** All I/O and service-boundary errors MUST be wrapped with context using `fmt.Errorf("failed to <operation>: %w", err)`. Never return raw OS/library errors.
- **Domain + OS Wrapping:** When combining a sentinel domain error with an originating OS error, use `errors.Join`: `fmt.Errorf("context: %w", errors.Join(core.ErrXxx, osErr))`. This preserves both for `errors.Is` inspection.
- **CLI Handler Pattern:** Command handlers in `src/internal/cli/cli.go` MUST use an `errors.Is` switch to intercept known domain errors, emit a user-friendly message via `core.UI`, and return `nil`  -  preventing cobra from printing raw error strings to the user.
- **Idempotency Guards:** Operations that mutate persistent state (e.g., `BootstrapProject`) MUST guard against re-execution by checking existing state and returning the appropriate domain sentinel error if the operation would be redundant.

## Specification & Task Organization
- **Hierarchical Task Structure:** Implementation roadmaps (`tasks.md`) MUST use a hierarchical structure to group atomic tasks. 
    - **Phases (H3):** Use `### Phase {N}: {Title}` for logical groupings.
    - **Tasks (H4):** Use `#### T{Phase}.{Task}: {Title}` for individual implementation steps.
- **DTO Pattern for Flat Access:** Domain models that implement nested structures (like `ImplementationReport` with `Phases`) MUST provide a `Tasks()` helper method that returns a flattened slice of all tasks. This ensures backward compatibility for scanners and progress calculators that expect a sequential list.
