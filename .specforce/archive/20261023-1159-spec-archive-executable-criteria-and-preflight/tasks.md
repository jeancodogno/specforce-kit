---
slug: 20261023-1159-spec-archive-executable-criteria-and-preflight
lens: Integration
---

# Implementation Roadmap: Executable Criteria in Spec and Preflight/Closing Improvements in Archive

## 1. Execution Strategy
- **Gravity Order:** RED manifest tests first -> artifact blueprints and planning skill -> archive instructions and skill pointer -> black-box verification. Every command below runs from the repository root unless it states a different directory; Go commands run inside `src`.
- **Module delta coverage:** `agent-kit` delta items `[NEW:exec-criteria]`, `[NEW:bug-red-first]`, `[NEW:grill-waves]`, `[NEW:readonly-explore]`, `[NEW:archive-preflight]`, `[NEW:archive-closing]` and use cases `[NEW:us-exec-criteria]`, `[NEW:us-bug-red-first]`, `[NEW:us-grill-waves]`, `[NEW:us-readonly-explore]`, `[NEW:us-archive-preflight]`, `[NEW:us-archive-closing]` are covered by T2.1 to T2.5 and proven by T1.1 and Phase 3.

## 2. Tasks

### Phase 1: RED Manifest Tests

- [x] T1.1: [CODE] Add failing manifest tests for the new blueprint rules
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-1], [US-2], [US-3], [US-4], [US-5], [US-6], [US-7], [US-8], [US-9]

**Action Steps:**
- Add `TestKitManifests_ExecutableCriteriaAndArchivePreflight` with subtests `Spec`, `TaskBlueprint`, `BugBlueprint`, `Archive` and `ArchiveSkill`, reusing the existing `assertFileContains` helper.
- `Spec` subtest: assert `kit/skills/spf-spec/SKILL.yaml` contains the phrases for the read-only exploration exception, "restated for a single confirmation", "approve all recommendations", "never carries", and "one-line justification"; assert it does NOT contain the substring "`A definir`, `Dúvidas`, `N/A`"; assert its line count is at most 155.
- `TaskBlueprint` and `BugBlueprint` subtests: assert `artifacts/spec/tasks.yaml` contains "non-interactive", "Black-box Verification", "`large`", "fails before the fix"; assert `artifacts/spec/bug-requirements.yaml` contains "root cause", "RED" and a halt-on-unknown-root-cause phrase.
- `Archive` and `ArchiveSkill` subtests: assert `kit/instructions/archive.md` contains "pending task", "even if the user asks", "Suggested Commit Message", "never run" and links under `.specforce/archive/`; assert the index of the completion-refusal text is smaller than the index of "Dynamic Agent & Skill Discovery"; assert `kit/skills/spf-archive/SKILL.yaml` still contains "Suggest Next Steps & Follow-up Specs" and "/spf:discovery".

**Acceptance Check:**
Run `cd src && go test ./internal/agent/ -run TestKitManifests_ExecutableCriteriaAndArchivePreflight` and confirm it FAILS with missing-directive errors for every subtest before any blueprint is edited.

### Phase 2: Blueprint Edits

- [x] T2.1: [CODE] Update the planning skill: grill waves, proposal pre-answers, read-only helpers, omit-or-justify rule
**Target:** `src/internal/agent/kit/skills/spf-spec/SKILL.yaml`
**Context:** [US-4], [US-5], [US-6], [US-7]
**Parallel With:** T2.2, T2.3, T2.4, T2.5

**Action Steps:**
- In Critical Rule 1 replace the literal `N/A` entry with: unresolved markers (`TBD`, `TODO`, bare `N/A`) remain forbidden, and a template section that does not apply is omitted or kept with a one-line justification.
- In Critical Rule 5 keep the ban on external planning sub-threads and add the exception that `large`/`complex` grounding may delegate strictly read-only exploration that returns findings, while interview, decisions and artifact writing stay in the main session.
- In Step 1A state that decisions taken in `proposal.md` are restated for a single confirmation, only open questions and newly found gaps are asked, and a taken decision that contradicts the Constitution or code is raised as a question with a recommendation.
- In Step 1C state that questions go in waves of at most 4 with the recommendation listed first, that "approve all recommendations" confirms only the current wave and never carries over, and keep the hard gate.
- In Step 2 add one sentence: task acceptance checks must be non-interactive commands, and `bug` specs follow their artifact blueprint; keep the file at 155 lines or fewer and keep all existing asserted phrases.

**Acceptance Check:**
Run `cd src && go test ./internal/agent/ -run 'TestKitManifests_ExecutableCriteriaAndArchivePreflight/Spec|TestSpecBlueprintTieredSizing|TestKitManifests_DeltaLivingSpecs|TestLifecycleCommandNaming'` and confirm all pass; run `wc -l src/internal/agent/kit/skills/spf-spec/SKILL.yaml` and confirm at most 155.

- [x] T2.2: [CODE] Make task acceptance checks executable and add black-box and bug-first rules
**Target:** `src/internal/agent/artifacts/spec/tasks.yaml`
**Context:** [US-1], [US-2], [US-3]
**Parallel With:** T2.1, T2.3, T2.4, T2.5

**Action Steps:**
- Extend instruction 3 so every Acceptance Check is a single non-interactive command or test scenario with an observable expected result, runnable in under 2 minutes, with no prompts, manual clicks or pasted secrets; criteria needing human judgment are rewritten into scriptable checks or marked as not verifiable by a tester.
- Add an instruction: for `large` and `complex` specs the last phase is `### Phase {N}: Black-box Verification`, one task per scenario tracing to `[US-x]` with command and expected observable result; `small` and `medium` omit it.
- Add an instruction: for `bug` specs task T1.1 creates the regression test and states that it fails before the fix, and no fix task precedes it.
- Add a Black-box Verification phase example to the template, marked to be rendered only for `large` and `complex`, and keep the YAML valid.

**Acceptance Check:**
Run `cd src && go test ./internal/agent/ -run 'TestKitManifests_ExecutableCriteriaAndArchivePreflight/TaskBlueprint'` and confirm it passes; run `specforce spec artifact spec.tasks --json | python3 -c "import json,sys; json.load(sys.stdin)"` and confirm exit code 0.

- [x] T2.3: [CODE] Require reproduction, root cause and RED regression link in bug requirements
**Target:** `src/internal/agent/artifacts/spec/bug-requirements.yaml`
**Context:** [US-3]
**Parallel With:** T2.1, T2.2, T2.4, T2.5

**Action Steps:**
- Strengthen the Reproduction Mandate so reproduction steps are concrete and runnable in order.
- Add a Root Cause rule: the root cause section must be based on evidence; if it is not established, halt and ask the user or recommend `/spf:discovery` instead of guessing.
- Add a rule that every `[FIX-x]` acceptance criterion maps to the first RED regression test task of the roadmap.
- Keep the template sections and the frontmatter rule unchanged.

**Acceptance Check:**
Run `cd src && go test ./internal/agent/ -run 'TestKitManifests_ExecutableCriteriaAndArchivePreflight/BugBlueprint'` and confirm it passes.

- [x] T2.4: [CODE] Add archive preflight, link summary and commit suggestion
**Target:** `src/internal/agent/kit/instructions/archive.md`
**Context:** [US-8], [US-9]
**Parallel With:** T2.1, T2.2, T2.3, T2.5

**Action Steps:**
- In Step 1 keep the status command; on progress below 100% (blocked or unfinished tasks included) stop before any retrospective, Constitution or module change, list each pending task as `T<id> - <title> - <state>`, point to `/spf:implement`, and refuse even if the user asks to archive anyway.
- In the Step 7 format add an `**Artifacts:**` list of clickable links to the existing `requirements.md`, `design.md` and `tasks.md` under `.specforce/archive/<slug>/`, keep the module line (state when no module changed), and add a `### Suggested Commit Message` code block with one conventional commit message before the next steps.
- Add a guardrail: suggest a commit message, never run it; no version-control command is executed.
- Keep every asserted phrase and keep the file under 100 lines.

**Acceptance Check:**
Run `cd src && go test ./internal/agent/ -run 'TestKitManifests_ExecutableCriteriaAndArchivePreflight/Archive|TestKitManifests_LivingSpecAndModuleReconciliation|TestKitManifests_DeltaLivingSpecs|TestLifecycleCommandNaming'` and confirm all pass; run `wc -l src/internal/agent/kit/instructions/archive.md` and confirm under 100.

- [x] T2.5: [CODE] Point the archive skill to the new preflight and summary rules
**Target:** `src/internal/agent/kit/skills/spf-archive/SKILL.yaml`
**Context:** [US-8], [US-9]
**Parallel With:** T2.1, T2.2, T2.3, T2.4

**Action Steps:**
- Add one sentence in the execution protocol saying the retrieved instructions define the refusal on incomplete work and the summary format with links and commit suggestion.
- Keep the "Suggest Next Steps & Follow-up Specs" step and the `/spf:discovery` reference unchanged.
- Keep the file a thin trigger with no multi-step logic.

**Acceptance Check:**
Run `cd src && go test ./internal/agent/ -run 'TestKitManifests_ExecutableCriteriaAndArchivePreflight/ArchiveSkill|TestLifecycleCommandNaming'` and confirm they pass.

- [x] T2.6: [CODE] Align the interview skill with the wave rule and the hard gate
**Target:** `src/internal/agent/kit/skills/consultative-grill/SKILL.yaml`
**Context:** [US-5]
**Parallel With:** T2.1

**Action Steps:**
- Replace the single-question rule in Interview Rules with waves of at most 4 questions, recommendation listed first, one at a time in Adversarial Mode.
- State that approving all recommendations confirms only the current wave.
- Reword Decision Persistence so decisions are written to artifacts once they exist, never before the planning gate is confirmed.

**Acceptance Check:**
Run `cd src && go test ./internal/agent/...` and `grep -c "Wave Precision" src/internal/agent/kit/skills/consultative-grill/SKILL.yaml` and confirm tests pass and the count is 1.

### Phase 3: Black-box Verification

- [x] T3.1: [VERIFY] Full test suite and static checks
**Target:** `src` (Go module)
**Context:** [US-1], [US-2], [US-3], [US-4], [US-5], [US-6], [US-7], [US-8], [US-9]

**Action Steps:**
- Run `cd src && gofmt -l .` and confirm empty output.
- Run `cd src && go vet ./...` and confirm exit code 0.
- Run `cd src && go test ./...` and confirm all packages pass, including the new test.

**Acceptance Check:**
Run `cd src && go test ./... && go vet ./...` and confirm exit code 0.

- [x] T3.2: [VERIFY] Built binary serves the updated blueprints
**Target:** `specforce` CLI (built from the repository)
**Context:** [US-1], [US-2], [US-3], [US-8], [US-9]

**Action Steps:**
- Build the binary into the scratchpad directory (`cd src && go build -o "$TMPDIR/specforce-qa" ./cmd/specforce`) without touching the global binary.
- Run `"$TMPDIR/specforce-qa" spec artifact spec.tasks --json` and confirm the output contains `Black-box Verification` and `non-interactive`.
- Run `"$TMPDIR/specforce-qa" archive instructions` and confirm the output contains `Suggested Commit Message` and the pending-task refusal text.

**Acceptance Check:**
Run `"$TMPDIR/specforce-qa" archive instructions | grep -c -e "Suggested Commit Message" -e "even if the user asks"` and confirm the count is at least 2; run `"$TMPDIR/specforce-qa" spec artifact spec.tasks --json | grep -c "Black-box Verification"` and confirm at least 1.

- [x] T3.3: [VERIFY] Delta and spec consistency
**Target:** `specforce` CLI (global)
**Context:** [US-1], [US-8]

**Action Steps:**
- Run `specforce spec delta check 20261023-1159-spec-archive-executable-criteria-and-preflight` and confirm the status is sealed and not conflicting.
- Run `specforce implementation status 20261023-1159-spec-archive-executable-criteria-and-preflight --json` and confirm no task is pending after implementation.
- Confirm the QA traceability matrix has one row per delta item listed in the Execution Strategy.

**Acceptance Check:**
Run `specforce spec delta check 20261023-1159-spec-archive-executable-criteria-and-preflight --json` and confirm no entry in `problems`.

## 3. Pre-emptive Mitigations
- **Risk:** Existing exact-phrase assertions break when editing the skill and archive text -> **Mitigation:** edits are additive; run the full `internal/agent` tests after each of T2.1 to T2.5.
- **Risk:** The planning skill grows past its size discipline -> **Mitigation:** the 155-line assertion in T1.1 fails the build if exceeded.
