---
slug: 20261020-2222-independent-qa-and-free-batch-log
lens: Backend-heavy
---

# Implementation Roadmap: Independent QA Tester, Free Batch Narration and Named Execution Roles

## 1. Execution Strategy
- **Gravity Order:** Manifest tests (RED) -> `spf-implement` blueprint rewrite (GREEN) -> command naming and archive next-step edits in `spf-spec`, `spf-archive`, `archive.md` (GREEN, parallel) -> module living spec -> full regression and line-budget verification.

## 2. Tasks

### Phase 1: Blueprint Behavior with TDD Pairs

- [x] T1.1: [TEST] [RED] Update and add kit manifest tests for the new protocol
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-1], [US-2], [US-3], [US-4], [US-5], [US-6], [US-7], [US-8], [US-9]

**Action Steps:**
- In `TestImplementBatchSubagentIsolationAndProgress`, remove the required tokens `Iniciando Lote`, `concluídas`, `Delegando para perfil de`, `BATCH` and `COMPLETED`; keep `Fresh Subagent per Batch`, `Within-Batch Feedback Loop`, `completed_tasks_count`, `total_tasks_count`; add the required token `Free-Form Batch Narration`; add a negative loop asserting `Iniciando Lote` and `Delegando para perfil de` are absent from the implement skill.
- In `TestActiveQAAndBlockerProtocols`, keep every existing token and append the required tokens `Traceability Matrix`, `NOT-VERIFIABLE`, `Correction Cycle`, `BLOCKED`, `Out-of-Scope Findings`.
- Add `TestImplementRoleChartersAndResume` that reads `kit/skills/spf-implement/SKILL.yaml`, requires the tokens `Role Charters`, `implementer`, `qa-tester`, `reviewer`, `Resume Handling`, and fails when `strings.Count(content, "\n")+1` exceeds 222.
- Add `TestLifecycleCommandNaming` that reads `spf-implement`, `spf-spec`, `spf-archive` and `kit/instructions/archive.md`; requires `/spf:archive` in implement, `/spf:implement` in spec, and `/spf:discovery` in archive skill and `archive.md`; and fails when any of those files contains the unprefixed strings `` `/archive` `` or `` `/implement` ``.
- Leave `TestImplementBlueprintGuardrails`, `TestImplementGlobalBatchBudgetAndSizing`, `TestEnhancedBlueprintsDirectives` and `TestPureSkillsBlueprintStructure` unchanged.
- Note: the Acceptance Check below negates `go test`, so exit code 0 proves the tests compile and currently fail (RED).

**Acceptance Check:**
`go vet ./src/internal/agent/... && ! go test ./src/internal/agent/ -run 'TestImplementBatchSubagentIsolationAndProgress|TestActiveQAAndBlockerProtocols|TestImplementRoleChartersAndResume|TestLifecycleCommandNaming'`

- [x] T1.2: [CODE] [GREEN] Rewrite the spf-implement blueprint
**Target:** `src/internal/agent/kit/skills/spf-implement/SKILL.yaml`
**Context:** [US-1], [US-2], [US-3], [US-4], [US-5], [US-6], [US-7], [US-9]

**Action Steps:**
- Step 1: add a `Resume Handling` block (report `in-progress` tasks and re-batch them as pending after inspecting affected code; ask the user via the interactive question tool for `failed` tasks with options retry, revise spec via `/spf:spec`, stop; skip straight to Step 4 when no pending tasks remain; keep the single status call).
- Step 2: delete duplicated sizing and task-boundary text so each is stated once; keep the headings `Fresh Subagent per Batch`, `Within-Batch Feedback Loop`, `Universal Harness Compatibility`, `Subagent Sizing Guidance`, `Global Batch Budget`, `1 Dedicated Subagent per Batch`, `Global Implementation Sizing`, `2 to 3 Batches (Recommended Sweet Spot)`, `at most 4 Batches`, `2 to 3 Subagents (Recommended Sweet Spot)` and `Adaptive Model & Effort Routing`.
- Step 2.B: replace the persona prose with a `Role Charters` table for `implementer`, `qa-tester` and `reviewer` (Mission, Allowed, Forbidden, Inputs, Output) and the kind-of-work to role/tier correspondence table from design.md section 1; mark `reviewer` as optional and read-only for high-complexity batches.
- Step 2.C: remove the `Recommended Model Tier` and `Recommended Reasoning Effort` lines from the Mission Brief; add to the worker guardrails a `BLOCKED Protocol` (workers never ask the user; they stop and return `BLOCKED: <question> | <options> | <recommendation>`) and an `Out-of-Scope Findings` return rule; keep every guardrail token asserted by `TestImplementBlueprintGuardrails` and `TestActiveQAAndBlockerProtocols` (including `Active Smoke Testing & Execution`, `Multimodal UI Visual Verification`, `Automated Screenshot Capture`, `Multimodal Visual Inspection`).
- Step 2.C log: replace the `MANDATORY LOG` template with `Free-Form Batch Narration`: dispatch the worker first, then write 1 to 3 sentences naming concrete areas, approach and verification, prefixed by a progress indicator using `{batch_index}/{total_batches}` and `{completed_tasks_count}/{total_tasks_count}`, in the user's most recent language, with no fixed emoji header or bracketed chains.
- Step 3: replace the `[BATCH ... COMPLETED]` line with a natural completion summary of what was delivered; route `BLOCKED` reports through the existing `Human-in-the-Loop & Blocker Protocol` and reply to the SAME worker; fix the duplicated step number "3.".
- Step 4: rename the persona to the `qa-tester` charter, set its inputs to `requirements.md`, `design.md`, `tasks.md` and the diff (never worker reports; Acceptance Checks of `tasks.md` for small specs), keep the `4-Tier QA Verification Protocol` and the four tier names verbatim, add the `Traceability Matrix` (criterion, scenario, evidence, verdict `PASSED | FAILED | NOT-VERIFIABLE`), forbid production edits by the tester, require the orchestrator to re-execute at least one key verification, and add the `Correction Cycle` (fresh correction worker plus fresh qa-tester, maximum 2 cycles, then escalate; behavior or architecture change halts for `/spf:spec`).
- Step 4.D: keep `4-Tier QA Verification Evidence` and `Developer Manual Walkthrough`; add `### Traceability Matrix`, `### Correction Cycles Used`, `### Out-of-Scope Findings`; end with `> Run \`/spf:archive\` to close the lifecycle.`
- Guardrails: delete the two bullets duplicated earlier (`Human-in-the-Loop Blocker Protocol`, `Task Boundary Safety`); keep `Mid/Post-Implementation Spec Gate`, `Terminal Truth`, `Scope Containment`, `No Interactive Prompts`, `CLI Execution`.
- Verify the resulting file has at most 222 lines and parses as YAML with `name` and `content`.

**Acceptance Check:**
`go test ./src/internal/agent/ -run 'TestImplementBatchSubagentIsolationAndProgress|TestActiveQAAndBlockerProtocols|TestImplementRoleChartersAndResume|TestImplementBlueprintGuardrails|TestImplementGlobalBatchBudgetAndSizing|TestEnhancedBlueprintsDirectives|TestPureSkillsBlueprintStructure'`

- [x] T1.3: [CODE] [GREEN] Standardize the spec hand-off command
**Target:** `src/internal/agent/kit/skills/spf-spec/SKILL.yaml`
**Context:** [US-8]
**Parallel With:** T1.4, T1.5

**Action Steps:**
- Replace the closing line `> Run \`/implement\` to start implementing.` with `> Run \`/spf:implement\` to start implementing.`
- Search the file for any other unprefixed lifecycle command (`/archive`, `/implement`) and prefix it with `spf:`.
- Do not change any other line of the file.

**Acceptance Check:**
`grep -c '/spf:implement' src/internal/agent/kit/skills/spf-spec/SKILL.yaml && ! grep -E '`/(archive|implement)`' src/internal/agent/kit/skills/spf-spec/SKILL.yaml`

- [x] T1.4: [CODE] [GREEN] Derive archive next steps from findings
**Target:** `src/internal/agent/kit/skills/spf-archive/SKILL.yaml`
**Context:** [US-6], [US-8]
**Parallel With:** T1.3, T1.5

**Action Steps:**
- In item 4 of the execution protocol, keep the heading text `Suggest Next Steps & Follow-up Specs` and add the directive to derive suggestions from out-of-scope findings and `NOT-VERIFIABLE` criteria present in the conversation.
- Add the constraints: at most 3 suggestions, each phrased as a ready-to-use `/spf:discovery` or `/spf:spec` request that cites its originating finding.
- Add the fallback: when no findings or unverifiable criteria exist, state that explicitly and base suggestions on the completed work.
- Keep the file valid YAML with `name`, `description` and `content`.

**Acceptance Check:**
`go test ./src/internal/agent/ -run 'TestEnhancedBlueprintsDirectives|TestPureSkillsBlueprintStructure'`

- [x] T1.5: [CODE] [GREEN] Align archive instructions with findings-based next steps
**Target:** `src/internal/agent/kit/instructions/archive.md`
**Context:** [US-6], [US-8]
**Parallel With:** T1.3, T1.4

**Action Steps:**
- Keep the heading `### Suggested Next Steps & Follow-up Specs:` unchanged.
- Replace the two example bullets with `- [Next Step 1, from: <finding or unverifiable criterion>]: <ready-to-use /spf:discovery or /spf:spec request>` and the equivalent second bullet.
- Add the line `At most 3 steps; if no findings were reported, state that and base the steps on the completed work.` below the bullets.

**Acceptance Check:**
`go test ./src/internal/agent/ -run 'TestEnhancedBlueprintsDirectives|TestLifecycleCommandNaming'`

### Phase 2: Living Spec and Final Verification

- [x] T2.1: [CODE] Update the agent-kit module living spec
**Target:** `.specforce/docs/modules/agent-kit.md`
**Context:** [US-1], [US-2], [US-3], [US-4], [US-5], [US-6], [US-7]

**Action Steps:**
- Rewrite `[BR-KIT-05]` so the final worker is an independent qa-tester; rewrite `[BR-KIT-09]` to the free-form post-dispatch narration standard (progress prefix, user language, completion summary) without the literal `🚀 **Iniciando Lote X/Y**` and `[BATCH X/Y COMPLETED]` formats; rewrite `[BR-KIT-12]` so the 4-tier protocol is governed by a criterion traceability matrix with a read-only tester and at most 2 correction cycles; extend `[BR-KIT-14]` with the BLOCKED protocol.
- Add `[BR-KIT-15]` (named role charters inside skills, no per-tool registration) and `[BR-KIT-16]` (out-of-scope findings and safe resumption of interrupted implementations).
- Update `[US-KIT-07]` and `[US-KIT-11]` to the new behavior; add `[US-KIT-14]` (bounded correction cycle) and `[US-KIT-15]` (BLOCKED escalation) in GIVEN/WHEN/THEN form.
- Keep the five-part module structure and include no internal source file paths.

**Acceptance Check:**
`grep -c 'BR-KIT-16' .specforce/docs/modules/agent-kit.md && ! grep -F 'Iniciando Lote' .specforce/docs/modules/agent-kit.md`

- [x] T2.2: [VERIFY] Run full regression and budget checks
**Target:** `Global Scope`
**Context:** [US-8], [US-9]

**Action Steps:**
- Run `go test ./...` and confirm zero failures across the repository.
- Run `wc -l src/internal/agent/kit/skills/spf-implement/SKILL.yaml` and confirm at most 222 lines.
- Run `grep -rn '`/archive`\|`/implement`' src/internal/agent/kit` and confirm no output.
- Run `go build ./...` and execute `specforce --dev init --help` to confirm the binary still starts with the local kit.

**Acceptance Check:**
`go test ./... && test "$(wc -l < src/internal/agent/kit/skills/spf-implement/SKILL.yaml)" -le 222 && ! grep -rn '`/archive`\|`/implement`' src/internal/agent/kit`

## 3. Pre-emptive Mitigations
- **Risk:** Consolidating duplicated guidance in `spf-implement` accidentally drops a guardrail asserted by existing tests. -> **Mitigation:** `TestImplementBlueprintGuardrails`, `TestImplementGlobalBatchBudgetAndSizing` and `TestEnhancedBlueprintsDirectives` stay unchanged and gate T1.2.
- **Risk:** The 222-line ceiling forces the new protocol to be cramped or cryptic. -> **Mitigation:** Delete duplicated sections first (sizing, task boundary, human-in-the-loop) to free space before adding new blocks; the ceiling is enforced by `TestImplementRoleChartersAndResume`.
