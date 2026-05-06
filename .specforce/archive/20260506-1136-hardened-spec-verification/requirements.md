---
slug: 20260506-1136-hardened-spec-verification
lens: Backend-heavy
---

# Feature: Hardened Spec Verification

## 1. Context & Value
Currently, AI agents may conclude a specification phase without formally verifying the integrity and completeness of the generated artifacts. This feature hardens the SDD workflow by making a final status check mandatory, ensuring all required files exist and are valid before the agent marks the planning task as finished.

## 2. Out of Scope (Anti-Goals)
- Modifying the CLI `status` command logic itself.
- Changing the `implement.yaml` or `discovery.yaml` workflows.
- Implementing automatic fix loops for validation errors.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Mandatory Verification Step
**User Story:** AS A developer, I WANT the AI agent to always run a final status check, SO THAT I can be confident the specification is complete and valid.

**Scenarios:**
1. **[Happy Path]** GIVEN the agent has written all required artifacts WHEN the pipeline reaches the handoff phase THEN the agent MUST execute `specforce spec status <slug> --json`.
2. **[Edge Case]** GIVEN the status command returns a `"progress" < 100` or `"is_valid": false` WHEN the agent performs the mandatory check THEN the agent MUST NOT provide the final summary and MUST attempt to fix the identified issues.

### [REQ-2] Verified Summary Handoff
**User Story:** AS A developer, I WANT the final summary to reflect the actual state of the artifacts, SO THAT I have a reliable handoff.

**Scenarios:**
1. **[Happy Path]** GIVEN a successful mandatory verification THEN the agent MUST output a summary including the file paths and progress confirmation.

## 4. Business Invariants
- No specification session can be concluded with a "success" message without a preceding `status --json` call in the same turn or session.
- The final summary MUST NOT be output if any artifact is missing or invalid according to the CLI status output.
