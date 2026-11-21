---
slug: 20260724-1509-fix-archival-instruction-gap
lens: Bugfix
---

# Bugfix: Fix Archival Instruction Gap

## 1. Issue Description
During the execution of the `spf.archive` skill, the AI agent omitted executing the specification status archival command `specforce spec archive <slug>`. The agent completed retrospective memory logging (`specforce archive memorial`) and memory distillation (`specforce archive distill`), but incorrectly assumed the specification was archived. This occurred due to instruction ambiguity in `archive.md` and lack of prominent warnings distinguishing the `specforce archive` memory namespace from the `specforce spec archive` specification namespace. Furthermore, the agent had to run `--help` to discover CLI subcommands due to this ambiguity.

## 2. Evidence & Observations
- **Symptom:** Feature specification lifecycle remains `active` after running `spf.archive` instead of transitioning to `archived`.
- **Trace:** Execution of `archive.md` instruction flow; step 8 `specforce spec archive <slug>` was omitted or missed due to namespace confusion with step 5/6 `specforce archive memorial`/`distill`.

## 3. Reproduction Steps
1. Run `specforce archive instructions` or trigger `/spf.archive`.
2. Follow the output instructions as an AI agent.
3. Observe that after completing `specforce archive memorial` and `specforce archive distill`, step 8 (`specforce spec archive <slug>`) is easily missed due to lack of explicit mandatory callout/guardrail distinguishing memory archive vs spec archive.

## 4. Root Cause Analysis (RCA)
1. **Namespace Overlap & Ambiguity:** `specforce archive` manages global memorial/distill memory fragments, while `specforce spec archive` manages individual specification lifecycle states.
2. **Instruction Weight:** `archive.md` heavily detailed memory harvesting and distillation, but presented Step 8 without prominent emphasis or explicit warnings that memory logging alone does NOT archive the spec.
3. **CLI Guidance Gap:** `specforce archive instructions` CLI output did not include a dedicated callout reminding agents of the dual-step requirement (Memory Archive vs Spec Archive).

## 5. Acceptance Criteria (Regression Tests)

### [FIX-1] Explicit Archival Dual-Step Mandate in Instruction Documentation
**Scenario: [Regression] Agent executes spf.archive instructions**
GIVEN the `archive.md` instruction file
WHEN an AI agent reads the Execution Protocol
THEN Step 8 MUST feature a prominent `CRITICAL DUAL-STEP REQUIREMENT` warning explicitly stating that memory logging (`specforce archive memorial`) does NOT archive the spec, and that executing `specforce spec archive <slug>` is strictly mandatory to close the spec lifecycle.

### [FIX-2] Guardrail Enforcement for Spec Status Archival
**Scenario: [Regression] Agent completes archival protocol**
GIVEN the Guardrails section of `archive.md`
WHEN evaluating completion criteria
THEN a mandatory guardrail MUST state that ending the archival process without running `specforce spec archive <slug>` is a protocol violation.

### [FIX-3] CLI Instructions Output Enhancement
**Scenario: [Regression] Execution of specforce archive instructions**
GIVEN the CLI executor in `src/internal/cli/archive.go`
WHEN `specforce archive instructions` is invoked
THEN the printed output MUST include a dedicated note highlighting the difference between `specforce archive ...` (Memory) and `specforce spec archive <slug>` (Spec Lifecycle).

## 6. Technical Constraints (NFR)
- **[Safety]:** Zero breaking changes to CLI command syntax or flag parameters.
- **[Observability]:** CLI instructions output must clearly demarcate namespaces for LLM agents.
