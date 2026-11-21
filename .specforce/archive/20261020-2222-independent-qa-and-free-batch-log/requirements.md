---
slug: 20261020-2222-independent-qa-and-free-batch-log
lens: Backend-heavy
---

# Feature: Independent QA Tester, Free Batch Narration and Named Execution Roles

## 1. Context & Value
During implementation, the developer currently sees rigid, metadata-only batch announcements and receives a final QA verdict that is not traced back to the acceptance criteria of the specification. This feature makes batch narration informative and free-form, turns the final QA into an independent tester that proves every acceptance criterion with evidence, and gives execution roles, blockers and findings a clear, consistent protocol.

## 2. Out of Scope (Anti-Goals)
- Do not introduce a mechanism for parallel execution of batches.
- Do not change the batch budget (2 to 3 batches, at most 4) or the one-fresh-worker-per-batch rule.
- Do not add new files, agent definitions or commands to the kit; only the existing lifecycle skills and the project's behavioral documentation are changed.
- Do not persist the QA report as a new specification artifact or change how specifications are archived.
- Do not introduce delta specifications, execution profiles or project-level customization of the workflow.
- Do not change the behavior of the discovery skill, nor of the spec skill beyond command naming.
- Do not document or configure per-task quality hooks.

## 3. Acceptance Criteria (BDD)

### [US-1] Informative, Free-Form Batch Narration
**User Story:** AS A developer following an implementation, I WANT TO read, right after each worker is dispatched, a short natural explanation of what that worker will do, SO THAT I understand the work in progress without decoding fixed metadata blocks.

**Scenarios:**
1. **[Happy Path]** GIVEN a batch is ready to run WHEN the orchestrator dispatches the worker THEN it first dispatches the worker and then writes an explanation of one to three sentences naming the concrete areas to be changed, the approach, and how the result will be verified.
2. **[Happy Path]** GIVEN the explanation is written WHEN the developer reads it THEN it begins with a short progress indicator showing the batch position and the number of completed tasks, and is written in the same language the developer uses.
3. **[Happy Path]** GIVEN a batch has been verified WHEN the orchestrator reports completion THEN it summarizes in natural language what was actually delivered, not only that verification passed.
4. **[Edge Case]** GIVEN the developer changes language mid-session WHEN the next batch is dispatched THEN the narration follows the developer's most recent language.
5. **[Edge Case]** GIVEN a batch contains a single trivial task WHEN narrating THEN the explanation stays within one to three sentences and does not pad with fixed headers or role metadata.

**UI/UX Specifics:**
- **View/Component:** Conversational message in the main session, one per batch dispatch and one per batch completion.
- **Feedback Logic:** The dispatch message appears after the worker is started; no fixed template, emoji header, or bracketed status chain is required.
- **Keybindings:** Not applicable to this requirement because the narration is passive text.

**Technical Constraints (NFR):**
- **[Performance]:** Narration adds no extra worker round trip and is limited to three sentences plus the progress indicator.
- **[Integrity]:** Every narration cites at least one concrete deliverable of the batch.
- **[Observability]:** The progress indicator always shows batch position and completed versus total tasks.

### [US-2] Independent QA Tester with Criterion Traceability
**User Story:** AS A developer receiving a finished implementation, I WANT TO have an independent tester prove each acceptance criterion of the specification with evidence, SO THAT I can trust the result without re-testing it myself.

**Scenarios:**
1. **[Happy Path]** GIVEN all batches are finished WHEN final validation starts THEN a fresh tester session is started that receives the specification's requirements, design and tasks together with the code changes, and does not receive the implementing workers' own reports.
2. **[Happy Path]** GIVEN the tester has finished WHEN it reports THEN the report lists every acceptance criterion with the scenario executed, the evidence collected, and a verdict of passed, failed or not verifiable.
3. **[Happy Path]** GIVEN the tester reports all criteria passed WHEN the orchestrator receives the report THEN it independently re-executes at least one key verification from the report before declaring the implementation passed.
4. **[Edge Case]** GIVEN a specification of the smallest size that has no requirements document WHEN the tester starts THEN the acceptance checks of the tasks serve as the criteria to prove.
5. **[Edge Case]** GIVEN a criterion cannot be exercised in the available environment WHEN the tester reports THEN that criterion is marked not verifiable with the reason, and is never silently marked as passed.
6. **[Edge Case]** GIVEN the tester finds a defect in production behavior WHEN it reports THEN it does not change production behavior itself; it may only create verification-only material such as additional tests or checking scripts.

**UI/UX Specifics:**
- **View/Component:** Final report with a traceability table (criterion, scenario, evidence, verdict) followed by the existing multi-tier evidence summary and the developer walkthrough.
- **Feedback Logic:** Overall status is passed only when no criterion is failed and every not-verifiable criterion is explicitly justified.
- **Keybindings:** Not applicable to this requirement because the report is passive text.

**Technical Constraints (NFR):**
- **[Performance]:** The orchestrator re-executes at most the key verifications needed to confirm the report, never the full set again.
- **[Integrity]:** Coverage of acceptance criteria in the report is 100%; no criterion is omitted.
- **[Safety & Security]:** The tester never modifies production behavior and never handles plaintext secrets.
- **[Observability]:** Every verdict is accompanied by reproducible evidence.

### [US-3] Bounded Correction Cycle after a Failed QA
**User Story:** AS A developer, I WANT failures found by the tester to be corrected and re-tested automatically within a limit, SO THAT trivial defects do not stop the flow while persistent ones reach me quickly.

**Scenarios:**
1. **[Happy Path]** GIVEN the tester reports at least one failed criterion WHEN the orchestrator handles the report THEN it dispatches a correction batch to a fresh worker with the failed criteria and evidence, and then starts a fresh tester session again.
2. **[Happy Path]** GIVEN a correction cycle succeeds WHEN the new tester report has no failed criterion THEN the implementation is declared passed and the cycle count is mentioned in the final summary.
3. **[Edge Case]** GIVEN two correction cycles have been used and criteria still fail WHEN the third report arrives THEN the orchestrator stops, presents the failing criteria with evidence, and asks the developer how to proceed.
4. **[Edge Case]** GIVEN a failure reveals missing requirements or a needed change of behavior or architecture WHEN the orchestrator evaluates it THEN it halts and directs the developer to update the specification before any code change instead of using a correction cycle.

**UI/UX Specifics:**
- **View/Component:** Narration message for each correction cycle, stating which criteria are being fixed.
- **Feedback Logic:** Each cycle is announced with its number out of two.
- **Keybindings:** Not applicable to this requirement because the cycle is non-interactive until escalation.

**Technical Constraints (NFR):**
- **[Performance]:** Maximum of 2 correction cycles per implementation, each with a fresh tester session.
- **[Integrity]:** The tester that judges a correction is never the same session as the worker that produced it.
- **[Observability]:** The final summary states the number of cycles consumed.

### [US-4] Named Execution Roles
**User Story:** AS A lead solutions architect, I WANT the roles used during implementation to be defined once with clear mission, permissions, inputs and outputs, SO THAT every harness and every batch uses the same expectations instead of improvised personas.

**Scenarios:**
1. **[Happy Path]** GIVEN the implementation guidance WHEN it is read THEN it defines the roles implementer, qa-tester and reviewer, each with mission, allowed actions, forbidden actions, required inputs and expected output format.
2. **[Happy Path]** GIVEN a batch of a given kind WHEN the orchestrator prepares it THEN it selects a role and a model and reasoning level from a fixed correspondence between kind of work and role.
3. **[Happy Path]** GIVEN a batch is of high complexity or touches security, concurrency or core architecture WHEN the batch is verified THEN a reviewer may be dispatched to inspect the result without modifying it.
4. **[Edge Case]** GIVEN a batch does not match any listed kind of work WHEN the orchestrator prepares it THEN it uses the implementer role and states its choice in the narration.

**UI/UX Specifics:**
- **View/Component:** Role definitions presented as a single compact table in the implementation guidance.
- **Feedback Logic:** The chosen role appears naturally in the batch narration only when relevant to the developer.
- **Keybindings:** Not applicable to this requirement because roles are guidance text.

**Technical Constraints (NFR):**
- **[Performance]:** The role table is no longer than one screen of about 40 lines, to keep guidance cost low.
- **[Maintainability]:** Each role is defined in exactly one place in the guidance.
- **[Integrity]:** Roles are described as part of the existing lifecycle skills only and require no per-tool registration.

### [US-5] Blocked-Worker Escalation Protocol
**User Story:** AS A developer, I WANT workers that cannot proceed to report the blocker to the orchestrator, SO THAT I am asked a clear question through the main session instead of the worker stalling or guessing.

**Scenarios:**
1. **[Happy Path]** GIVEN a worker faces ambiguity, multiple valid options, a possibly incorrect test, or missing credentials WHEN it cannot proceed THEN it stops and returns a blocked report containing the question, the options and its recommendation.
2. **[Happy Path]** GIVEN the orchestrator receives a blocked report WHEN it consults the developer THEN it presents structured trade-offs with a recommendation and, after the answer, returns it to the same worker so the batch continues.
3. **[Edge Case]** GIVEN the blocker is a missing secret WHEN the orchestrator consults the developer THEN it asks the developer to provide it through the local environment configuration and never requests the secret value in conversation.
4. **[Edge Case]** GIVEN the blocker implies a change of behavior or architecture WHEN the orchestrator evaluates it THEN it halts the implementation and requires a specification update before continuing.

**UI/UX Specifics:**
- **View/Component:** Structured question with labeled options and a marked recommendation.
- **Feedback Logic:** The batch is visibly paused until the developer answers.
- **Keybindings:** Not applicable to this requirement because selection uses the environment's question tool.

**Technical Constraints (NFR):**
- **[Performance]:** A blocker is returned by the worker in a single report, with no more than one question per blocker.
- **[Safety & Security]:** Plaintext secret values are never requested, logged or echoed.
- **[Integrity]:** Workers never ask the developer directly; only the orchestrator does.

### [US-6] Out-of-Scope Findings Carried to Closing
**User Story:** AS A developer, I WANT workers and the tester to report relevant findings that were outside the current scope, SO THAT follow-up work is based on real observations instead of generic suggestions.

**Scenarios:**
1. **[Happy Path]** GIVEN a worker or the tester notices an issue outside the active scope WHEN it reports THEN it lists the finding separately and does not act on it.
2. **[Happy Path]** GIVEN findings were reported WHEN the orchestrator produces the final summary THEN it includes a consolidated, deduplicated list of findings.
3. **[Happy Path]** GIVEN the specification is archived WHEN suggested next steps are produced THEN they are derived from the consolidated findings and from criteria that could not be verified, up to three suggestions, each phrased as a ready-to-use request for discovery or planning.
4. **[Edge Case]** GIVEN no findings were reported WHEN archival suggests next steps THEN it bases them on the completed work and states that no out-of-scope findings were reported.

**UI/UX Specifics:**
- **View/Component:** A short "Out-of-scope findings" list in the final summary and a numbered "Next steps" list at archival.
- **Feedback Logic:** Each next step names its originating finding.
- **Keybindings:** Not applicable to this requirement because lists are passive text.

**Technical Constraints (NFR):**
- **[Performance]:** At most 3 next-step suggestions per archival.
- **[Integrity]:** Findings never trigger changes to code in the current implementation.
- **[Observability]:** Every suggestion cites the finding or the unverifiable criterion that originated it.

### [US-7] Safe Resumption of an Interrupted Implementation
**User Story:** AS A developer returning to an interrupted implementation, I WANT the orchestrator to recognize tasks left in progress or failed, SO THAT no work is silently skipped or repeated.

**Scenarios:**
1. **[Happy Path]** GIVEN the roadmap contains tasks left in progress WHEN the orchestrator starts THEN it reports them and includes them in the next batch as pending work, after inspecting the current state of the affected code.
2. **[Happy Path]** GIVEN the roadmap contains tasks marked as failed WHEN the orchestrator starts THEN it asks the developer whether to retry them, revise the specification, or stop.
3. **[Edge Case]** GIVEN every task is already finished WHEN the orchestrator starts THEN it skips batch execution and proceeds directly to the final independent validation.
4. **[Edge Case]** GIVEN the code state of an in-progress task cannot be reconciled with its description WHEN the orchestrator inspects it THEN it asks the developer before continuing.

**UI/UX Specifics:**
- **View/Component:** One initial message listing the resumed and failed tasks.
- **Feedback Logic:** Resumed tasks are labeled as such in the batch narration.
- **Keybindings:** Not applicable to this requirement because the interaction uses the environment's question tool.

**Technical Constraints (NFR):**
- **[Performance]:** The roadmap status is requested once per session, as before.
- **[Integrity]:** No task is marked finished without its verification passing.
- **[Observability]:** Resumed and failed tasks are named explicitly in the initial message.

### [US-8] Consistent Next-Step Command Naming
**User Story:** AS A developer moving across the lifecycle, I WANT every suggested next command to use the same naming convention, SO THAT I can copy it and run it without translating between variants.

**Scenarios:**
1. **[Happy Path]** GIVEN any lifecycle skill suggests a follow-up WHEN it names the next command THEN it uses the common `spf` prefixed naming for planning, implementation and archival.
2. **[Edge Case]** GIVEN a skill refers to another skill in its guardrails WHEN it names it THEN it uses the same prefixed naming rather than an unprefixed or alternative form.

**UI/UX Specifics:**
- **View/Component:** Closing hand-off line of each lifecycle skill.
- **Feedback Logic:** The hand-off line names exactly one next command.
- **Keybindings:** Not applicable to this requirement because the hand-off is passive text.

**Technical Constraints (NFR):**
- **[Performance]:** Command naming changes add no extra guidance lines.
- **[Maintainability]:** Zero occurrences of unprefixed lifecycle command names remain in the four lifecycle skills.

### [US-9] Non-Redundant, Non-Contradictory Implementation Guidance
**User Story:** AS A lead solutions architect, I WANT the implementation guidance to state each rule exactly once and without contradictions, SO THAT agents spend fewer tokens and follow a single unambiguous protocol.

**Scenarios:**
1. **[Happy Path]** GIVEN the revised guidance WHEN it is read end to end THEN task-state authority, blocker handling and worker sizing are each stated in exactly one place.
2. **[Happy Path]** GIVEN the revised guidance WHEN the worker instructions are read THEN they contain no instruction that requires a worker to talk directly to the developer.
3. **[Edge Case]** GIVEN any rule that previously appeared in more than one place WHEN the guidance is consolidated THEN no behavioral rule is lost, and the previously defined guardrails (scope containment, test invariance, active smoke testing, visual verification, secret handling, specification gating) remain in force.

**UI/UX Specifics:**
- **View/Component:** The implementation guidance as a single readable document.
- **Feedback Logic:** Not applicable because this requirement concerns document quality.
- **Keybindings:** Not applicable because this requirement concerns document quality.

**Technical Constraints (NFR):**
- **[Performance]:** The revised guidance is not longer than the current one, measured in lines, despite the added behavior.
- **[Maintainability]:** Each behavioral rule appears once.
- **[Integrity]:** Every guardrail present before this change is still present after it.

## 4. Business Invariants
- The orchestrator never edits production code; all code changes come from dispatched workers.
- Only the orchestrator changes task states.
- Each batch is executed by a fresh worker, and the final validation is executed by a fresh tester that is never an author of the code it judges.
- The tester never modifies production behavior.
- Workers never address the developer directly; blockers always flow through the orchestrator.
- An implementation is declared passed only when every acceptance criterion is passed, or explicitly justified as not verifiable.
- No plaintext secret is ever requested, logged or echoed.
- Behavioral or architectural change discovered during execution always requires a specification update first.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Net guidance length does not grow; narration is bounded to three sentences; at most 2 correction cycles; at most 3 next-step suggestions.
- **[Reliability]:** Every failure state ends either in a defined recovery cycle or in an explicit question to the developer; there are no silent failures.
- **[Security]:** Plaintext secrets are never requested or displayed; the tester never alters production behavior.
- **[Maintainability]:** Each rule is defined in a single place; the automated test suite of the kit verifies the presence of the new protocol and the absence of the removed fixed log template.
