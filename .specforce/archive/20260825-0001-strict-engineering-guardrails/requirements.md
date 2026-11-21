---
slug: 20260825-0001-strict-engineering-guardrails
lens: Integration
---

# Feature: Strict Engineering Guardrails and Test Invariance

## 1. Context & Value
Automated implementation workflows currently lack exhaustive, non-negotiable coding, testing, and alignment principles in execution prompts, brief envelopes, and constitutional baselines. This feature embeds the complete protocol across all phases:
- **Worker Level:** Pre-coding clarification, anti-sycophancy, scope jail, and absolute test invariance.
- **Orchestrator Level:** Mid-flight and post-implementation specification gating (prohibiting direct ad-hoc code edits when new behaviors or changes are requested, mandating a return to `/spf:spec`).

## 2. Out of Scope (Anti-Goals)
- Dynamic AST-based source code analyzers or static linters implemented as external Go daemon processes.
- Modifying specification templates outside implementation orchestration, execution envelopes, and constitutional engineering baselines.
- Interactive terminal UI redesign or alterations to CLI flag parsers.

## 3. Acceptance Criteria (BDD)

### [US-1] Pre-Coding Protocol (Adversarial Alignment & Zero Assumptions)
**User Story:** AS A software engineer overseeing autonomous execution, I WANT the worker and orchestrator to enforce an exhaustive pre-coding interrogation protocol, SO THAT execution never proceeds under ambiguity, uncertainty, or unexamined technical flaws.

**Scenarios:**
1. **[Happy Path - Uncertainty & Ambiguity Resolution]** GIVEN an assigned task with uncertain or ambiguous requirements WHEN the agent evaluates the task before coding THEN the agent halts and actively asks clarifying questions.
2. **[Happy Path - Multiple Interpretations Handling]** GIVEN a requirement with multiple possible architectural or technical interpretations WHEN evaluated THEN the agent is strictly forbidden from picking one silently, and MUST present all options to the user for explicit confirmation.
3. **[Edge Case - Anti-Sycophancy & Direct Challenge]** GIVEN a planned approach or design that is technically flawed or suboptimal WHEN analyzed THEN the agent MUST disagree, reject sycophancy/flattery, and actively challenge the flaw with technical arguments.

**Technical Constraints (NFR):**
- **[Performance]:** Prompt resolution and validation overhead < 5ms.
- **[Safety & Security]:** Strict ban on unprompted file system modifications when ambiguities are detected.
- **[Integrity]:** All escalations must present concrete options rather than silent assumptions.
- **[Observability]:** Blocked states must clearly output the trade-off options.

### [US-2] During-Coding Protocol (Scope Jail & Test Invariance)
**User Story:** AS A technical lead, I WANT execution workers to follow strict coding boundaries and absolute test invariance, SO THAT production code conforms strictly to tests without unauthorized changes or artificial passes.

**Scenarios:**
1. **[Happy Path - Minimal Conforming Scope]** GIVEN an active implementation batch WHEN writing code THEN the worker MUST NOT implement unrequested features, MUST NOT refactor code that was not requested, MUST NOT implement tests for impossible scenarios, and MUST NOT remove pre-existing dead code unless explicitly requested.
2. **[Edge Case - Absolute Test Invariance & Anti-Tampering]** GIVEN failing verification tests WHEN the worker seeks to resolve failures THEN:
   - The worker MUST NEVER remove tests to reduce the failure count.
   - The worker MUST NEVER weaken existing assertions or types to make them pass.
   - The worker MUST NEVER use mechanisms that skip, ignore, or bypass tests (e.g. `skip`, `ignore`, dummy mocks, tautological asserts).
3. **[Edge Case - Genuinely Defective Test Protocol]** GIVEN an existing test that is genuinely obsolete or incorrect WHEN discovered THEN the worker MUST stop and confirm with the user before modifying the test.
4. **[Happy Path - Tests as the Definitive Specification]** GIVEN the relationship between tests and code WHEN implementing THEN tests are treated as the immutable specification: the implementation code conforms to the tests, NEVER the reverse.

**Technical Constraints (NFR):**
- **[Performance]:** Instantaneous constraint adherence within worker memory.
- **[Safety & Security]:** Full preservation of existing test suites and assertions.
- **[Integrity]:** Verification requires exit code 0 under unaltered test assertions.
- **[Observability]:** Explicit logging of worker compliance directives in the execution envelope.

### [US-3] Orchestrator Mid/Post-Implementation Spec Gate
**User Story:** AS A software architect, I WANT the Primary Implementation Orchestrator to block ad-hoc code edits when the user requests new features or behavioral modifications mid-flight or after completion, SO THAT any behavioral change is formally updated in the specification first.

**Scenarios:**
1. **[Happy Path - User Modification Request Gating]** GIVEN an ongoing or completed implementation session WHEN the user requests a new feature, a change in business logic, or altered behavior not present in the active specification THEN the Orchestrator is strictly forbidden from directly editing code or dispatching workers, and MUST instruct the user to activate `/spf:spec` to update `requirements.md`, `design.md`, and `tasks.md` first.
2. **[Edge Case - Unplanned Architectural Pivot]** GIVEN an implementation blocker that requires changing API contracts or data models WHEN encountered THEN the Orchestrator stops implementation and routes the user back to the planning pipeline.

**Technical Constraints (NFR):**
- **[Performance]:** Verification runs in < 2s for local tests.
- **[Safety & Security]:** Zero unauthorized state mutations outside approved specifications.
- **[Integrity]:** Total consistency across `requirements.md`, `design.md`, and code.
- **[Observability]:** Actionable handoff instructions to `/spf:spec` or `/spf:archive`.

## 4. Business Invariants
- `[BR-GUARD-01]` **Uncertainty & Ambiguity:** If anything is uncertain or ambiguous, ask.
- `[BR-GUARD-02]` **No Silent Choices:** If multiple interpretations exist, present options and ask; never pick silently.
- `[BR-GUARD-03]` **Anti-Sycophancy:** If an approach appears flawed, disagree, avoid flattery, and challenge it.
- `[BR-GUARD-04]` **Scope Jail (Worker):** Do NOT code unrequested features; do NOT refactor unrequested code; do NOT write tests for impossible scenarios; do NOT remove pre-existing dead code.
- `[BR-GUARD-05]` **Test Invariance (Worker):** NEVER delete tests to reduce failures; NEVER weaken tests to make them pass; NEVER use skips/ignores/fake bypasses; if a test is genuinely wrong, HALT and confirm with the user. Tests are the specification.
- `[BR-GUARD-06]` **Orchestrator Specification Gate:** When the user requests new behaviors or feature additions (during or after implementation), the Orchestrator MUST NOT write code directly; it MUST mandate updating the specification via `/spf:spec` first.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Blueprint translation overhead < 5ms.
- **[Reliability]:** Fail-fast behavior on detected ambiguities or test degradation attempts.
- **[Maintainability]:** Clean separation between Orchestrator Governance and Worker Directives.
