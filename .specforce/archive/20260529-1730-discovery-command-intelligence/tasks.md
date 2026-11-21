# Implementation Roadmap: Discovery Command Intelligence Refactor

### Phase 1: Command Refactor
- [x] T1.1: Refactor `spf.discovery` Command Intelligence
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [REQ-1, REQ-2]
**Action Steps:**
- Update the command content to integrate the Dual-Mode (Scout & Detective) heuristics.
- Implement the 3-layer Consultative Funnel (Constitution -> Empirical -> Brainstorming).
- Ensure the "NON-MUTATION COVENANT" is prominently displayed as the primary guardrail.
- Add Ghost Protocol standards for ASCII wireframing.
**Acceptance Check:**
- `cat src/internal/agent/kit/commands/discovery.yaml | grep "BUG_DETECTIVE"` and `grep "FEATURE_SCOUT"`.

### Phase 2: Verification
- [x] T2.1: Verify Discovery Command Integrity
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [AC-1, AC-3]
**Action Steps:**
- Perform a manual audit of the command to ensure no delegation logic or subagent references remain.
- Verify that the instructions do not mention `write_file` or `replace` except to forbid them.
- Run `make build` and verify that `.agents/workflows/spf-discovery.md` is updated automatically.
**Acceptance Check:**
- `grep -L "MISSION BRIEF" src/internal/agent/kit/commands/discovery.yaml` (Should return the filename).
- `cat .agents/workflows/spf-discovery.md | grep "BUG_DETECTIVE"`
