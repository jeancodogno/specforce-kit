---
slug: 20260509-1132-archive-error-prevention
lens: Backend-heavy
---

# Feature: Archive Error Prevention

## 1. Context & Value
The current archiving workflow focuses on harvesting new architectural precedents, but lacks a mandatory step for analyzing challenges and bugs encountered during implementation. Enhancing the `spf.archive` instructions to explicitly ask the agent to review the `tasks.md` and failure history ensures that lessons learned from bugs are captured and translated into Constitution updates, preventing the same errors from repeating in future development cycles.

## 2. Out of Scope (Anti-Goals)
- Do not alter the actual Go codebase or CLI logic for the `archive` command; this is strictly an update to the agent's markdown instructions.
- Do not modify the structure of the distributed memorial; we are only changing what the agent *looks for* before deciding to create a memorial or update the constitution.
- Do not introduce interactive prompts beyond the existing capability (the user confirmation step).

## 3. Acceptance Criteria (BDD)

### [REQ-1] Enhanced Specification Retrospective
**User Story:** AS A Principal Architect Agent, I WANT TO analyze the completed feature's challenges and bugs, SO THAT I can identify recurring issues and lessons learned.

**Scenarios:**
1. **[Happy Path]** GIVEN an agent is executing the archive instructions WHEN they reach the Specification Retrospective step THEN they are explicitly instructed to scan `requirements.md`, `design.md`, and `tasks.md` to analyze challenges and bugs, in addition to precedents.
2. **[Edge Case]** GIVEN an agent is executing the archive instructions WHEN they find no bugs or challenges in the `tasks.md` THEN they proceed to the next step without proposing error-prevention updates.

**UI/UX Specifics:**
- **View/Component:** N/A (Agent Instruction).
- **Feedback Logic:** N/A.
- **Keybindings:** N/A.

**Technical Constraints (NFR):**
- **[Performance]:** Minimal impact, adding text to a markdown file.
- **[Safety & Security]:** Must not break the existing YAML structure or markdown formatting of the instruction file.
- **[Integrity]:** The instruction file must remain grammatically correct and logically sequential.
- **[Observability]:** N/A.

### [REQ-2] Error Prevention in Constitution Analysis
**User Story:** AS A Principal Architect Agent, I WANT TO evaluate if encountered bugs indicate a missing rule in the Constitution, SO THAT I can propose a new global standard to prevent future occurrences.

**Scenarios:**
1. **[Happy Path]** GIVEN the agent has identified a bug or challenge WHEN they perform the Constitution Impact Analysis THEN they are instructed to evaluate if a missing rule or lack of clarity caused the issue and formulate a rule to prevent it.
2. **[Edge Case]** GIVEN the agent identifies a bug that is already covered by a recent Constitution update WHEN they evaluate the impact THEN they do not propose a redundant update.

**UI/UX Specifics:**
- **View/Component:** N/A (Agent Instruction).
- **Feedback Logic:** N/A.
- **Keybindings:** N/A.

**Technical Constraints (NFR):**
- **[Performance]:** N/A.
- **[Safety & Security]:** N/A.
- **[Integrity]:** The instructions must clearly separate the concepts of "Precedents" and "Error Prevention".
- **[Observability]:** N/A.

### [REQ-3] Updated Information Gathering Prompt
**User Story:** AS A Principal Architect Agent, I WANT TO explicitly ask the user about updating the constitution to prevent repeated errors, SO THAT the user understands the context of the proposed change.

**Scenarios:**
1. **[Happy Path]** GIVEN the agent has formulated an error-prevention rule WHEN they reach the Information Gathering step THEN the suggested prompt to the user explicitly mentions the encountered challenge and asks to update the Constitution to prevent it from repeating.
2. **[Edge Case]** GIVEN the user denies the constitution update WHEN prompted THEN the agent proceeds to the archival execution step as normal.

**UI/UX Specifics:**
- **View/Component:** N/A (Agent Instruction).
- **Feedback Logic:** N/A.
- **Keybindings:** N/A.

**Technical Constraints (NFR):**
- **[Performance]:** N/A.
- **[Safety & Security]:** N/A.
- **[Integrity]:** The example prompt string must be clear and encompass both precedents and challenges.
- **[Observability]:** N/A.

## 4. Business Invariants
- The agent must always confirm with the user before mutating Constitution documents.
- The archive process must remain the final step of the feature lifecycle.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Changes are text-based and should not affect CLI execution time.
- **[Reliability]:** The modified instructions must be clear enough for the LLM to follow reliably without confusion.
- **[Security]:** N/A.
- **[Maintainability]:** The instruction file must follow the established formatting conventions (numbered lists, bold emphasis).