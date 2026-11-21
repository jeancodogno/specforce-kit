# Feature: Discovery Command Intelligence Refactor

## [REQ-1] Integrated Discovery Intelligence (`spf.discovery`)
**User Story:** AS A Developer, I WANT the discovery command to act as a specialized thought partner directly within my session, SO THAT I can brainstorm features or debug issues without the context loss or token overhead of subagent delegation.

- **Persona:** Senior Discovery Specialist (The Scout & Detective).
- **Mindset:** Dual-mode (Product Scout / Technical Detective) integrated into the command instructions.
- **Constraint:** STRICT Non-Mutation Covenant. The discovery phase remains strictly read-only.
- **Integrated Heuristics:** ASCII wireframing (Ghost Protocol), Trace Analysis, and Thread-Driven Exploration.

## [REQ-2] Intent-Driven Discovery Funnel
**User Story:** AS A Developer, I WANT the discovery process to automatically adapt its lens based on my request, SO THAT I get the right expert heuristics for either new ideas or bug investigations.

- **Layer 1: Constitutional Anchor:** Immediate sync with project principles.
- **Layer 2: Empirical Grounding:** Direct codebase scan via `grep` and `read`.
- **Layer 3: Consultative Brainstorming:** Generation of 1-3 distinct technical paths.
- **Automatic Classification:** The agent adopts the `BUG_DETECTIVE` lens for issues or `FEATURE_SCOUT` for ideas based on prompt analysis.

## [REQ-3] Token & Context Efficiency
**User Story:** AS A Developer, I WANT the discovery process to be highly token-efficient, SO THAT I don't pay for redundant context reads or delegation headers.

- **Inline Execution:** No subagents called; the orchestrator performs all research and brainstorming.
- **Context Persistence:** All insights from Discovery remain in the primary session memory for subsequent Planning (`/spf:spec`).

## Acceptance Criteria

### [AC-1] Refactored Discovery Command
- GIVEN the `src/internal/agent/kit/commands/discovery.yaml` file.
- WHEN I check the instructions.
- THEN they must contain the dual-mode heuristics (Scout/Detective) and the 3-layer funnel.

### [AC-2] Unified Session Context
- GIVEN a discovery session has analyzed the codebase.
- WHEN the user pivots to `/spf:spec`.
- THEN the agent must demonstrate it still holds the technical insights gained during discovery without re-reading the same files.

### [AC-3] Read-Only Enforcement
- GIVEN any discovery session.
- WHEN the agent is active.
- THEN it must be strictly forbidden from using `write_file`, `replace`, or other mutation tools.
