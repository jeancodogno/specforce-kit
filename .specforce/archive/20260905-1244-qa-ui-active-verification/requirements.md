---
slug: 20260905-1244-qa-ui-active-verification
lens: Integration
---

# Feature: Active QA, Multimodal UI Verification & Blocker Protocol

## 1. Context & Value
Enhances the Specforce implementation engine and agent blueprints by eliminating the "green test illusion" through mandatory active black-box verification, introducing multimodal visual validation for user interface tasks, and enforcing a structured protocol for managing credentials, external dependencies, and technical ambiguities. This guarantees that finished features work reliably in real execution environments without compromising system security.

## 2. Out of Scope (Anti-Goals)
- Building custom browser automation frameworks or visual regression engines inside the Specforce core binary.
- Storing or transmitting secret keys or credentials in plain text across communication channels.
- Modifying the core CLI task execution state machine outside of existing batch commands.

## 3. Acceptance Criteria (BDD)

### [US-1] Multi-Tier Active QA Verification Protocol
**User Story:** AS A software stakeholder, I WANT the implementation workflow to actively execute and verify the complete feature in a realistic environment, SO THAT features are only marked as completed when real-world execution succeeds.

**Scenarios:**
1. **[Happy Path]** GIVEN an implementation roadmap where all batch tasks have completed WHEN the final quality assurance evaluation executes THEN it executes automated regression checks, performs active black-box verification against real compiled artifacts, validates error handling, and generates a structured verification report with an interactive developer walkthrough.
2. **[Edge Case - Execution Failure]** GIVEN a feature where automated unit tests pass but real-world black-box execution produces unexpected behavior WHEN the quality assurance evaluation detects the discrepancy THEN it marks the quality status as failed, documents the specific behavioral deviation, and halts completion without marking the feature ready for archival.

**UI/UX Specifics:**
- **View/Component:** Quality Assurance Summary & Walkthrough Panel.
- **Feedback Logic:** Displays distinct operational tiers with real output evidence and step-by-step reproduction instructions.
- **Keybindings:** Standard terminal output format.

**Technical Constraints (NFR):**
- **[Performance]:** Quality assurance execution completes deterministically within the standard task timeout budget.
- **[Safety & Security]:** Verification operations must execute within project boundaries without destructive side-effects.
- **[Integrity]:** Quality status must strictly require all applicable verification tiers to pass before issuing completion sign-off.
- **[Observability]:** Output must provide explicit evidence of real execution results rather than generic pass/fail assertions.

### [US-2] Multimodal Visual Verification for User Interface Tasks
**User Story:** AS A product designer, I WANT worker agents implementing user interfaces to visually inspect their rendered output, SO THAT layout, visual styling, responsive states, and interaction behavior strictly adhere to design guidelines.

**Scenarios:**
1. **[Happy Path]** GIVEN an implementation task that creates or modifies user interface components WHEN the worker agent completes code changes THEN it captures rendered visual snapshots across relevant states (e.g., standard, responsive, error), inspects the visual output, and verifies alignment against design specifications.
2. **[Edge Case - Visual Discrepancy or Non-Multimodal Environment]** GIVEN a worker operating in an environment where visual rendering differs from design or where image inspection is restricted WHEN evaluating the visual output THEN the worker adjusts the styling until visual parity is achieved or falls back to structured document structure verification while noting the visual inspection status.

**UI/UX Specifics:**
- **View/Component:** Worker Visual Verification Log.
- **Feedback Logic:** Reports visual confirmation status alongside standard terminal verification results.
- **Keybindings:** Non-interactive execution.

**Technical Constraints (NFR):**
- **[Performance]:** Snapshot capture and inspection must introduce minimal overhead into the task execution loop.
- **[Safety & Security]:** Local development servers and preview environments must bind strictly to local interfaces.
- **[Integrity]:** Visual discrepancies must be treated as verification failures requiring remediation before claiming batch completion.
- **[Observability]:** Verification logs must document whether visual inspection was successfully executed.

### [US-3] Secure Blocker, Credential & Decision Management Protocol
**User Story:** AS A project security officer, I WANT agents to follow strict security and communication rules when encountering missing credentials or ambiguous technical decisions, SO THAT secret data is never exposed and scope changes are properly governed.

**Scenarios:**
1. **[Happy Path - Missing Configuration]** GIVEN a worker or orchestrator encountering a missing environment variable or external service credential WHEN executing verification steps THEN it updates configuration templates with placeholder keys, prompts the user to provide local environment configuration, and securely pauses until the configuration is available.
2. **[Happy Path - Technical Ambiguity]** GIVEN a minor implementation ambiguity within the approved scope WHEN a worker identifies multiple valid implementation options THEN it presents a structured decision inquiry to the user with clear trade-offs and recommendations rather than guessing silently.
3. **[Edge Case - Scope or Architectural Drift]** GIVEN an unanticipated requirement or architectural pivot discovered during implementation WHEN evaluating the change THEN the agent immediately halts code modification and activates the specification gate to formalize the requirement before proceeding.

**UI/UX Specifics:**
- **View/Component:** Interactive Consultation & Action Gate Notification.
- **Feedback Logic:** Clear action items with structured multiple-choice options or environment instructions.
- **Keybindings:** Standard interactive selection keys.

**Technical Constraints (NFR):**
- **[Performance]:** Blocker detection must occur immediately upon prerequisite evaluation without polling loops.
- **[Safety & Security]:** Zero plaintext secrets or sensitive tokens may appear in conversation logs or messages.
- **[Integrity]:** Scope alterations require explicit specification updates prior to source code modification.
- **[Observability]:** All user decisions and blocker resolutions are recorded in project audit artifacts.

## 4. Business Invariants
- Quality Assurance status must never be reported as passed based exclusively on automated unit tests when runtime artifacts or interfaces are involved.
- Secret credentials, API keys, and sensitive tokens must never be solicited, stored, or displayed in conversational transcripts.
- Architectural pivots and out-of-scope requirements discovered during implementation must never be coded without updating formal specification artifacts.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Blueprint generation and execution guidance must remain token-efficient and deterministic.
- **[Reliability]:** Verification and blocker protocols must prevent infinite loops, silent failures, and unguided stops.
- **[Security]:** Strict zero-secret-leakage policy in conversational context and blueprints.
- **[Maintainability]:** Blueprints must adhere strictly to Specforce embedded skill standards and schema specifications.
