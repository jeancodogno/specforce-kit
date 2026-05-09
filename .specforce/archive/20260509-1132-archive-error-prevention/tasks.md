---
slug: 20260509-1132-archive-error-prevention
lens: Backend-heavy
---

# Implementation Roadmap: Archive Error Prevention

## 1. Execution Strategy
- **Gravity Order:** Update markdown instructions file `src/internal/agent/kit/instructions/archive.md`.

## 2. Tasks

### Phase 1: Update Agent Instructions

- [x] T1.1: [DOCS] Enhance Specification Retrospective
**Target:** `src/internal/agent/kit/instructions/archive.md`
**Context:** [REQ-1]

**Action Steps:**
- Update `### 3. Specification Retrospective` section.
- Add instructions to explicitly scan `tasks.md` in addition to `requirements.md` and `design.md`.
- Add a new bullet point to analyze challenges, roadblocks, and bugs encountered during implementation (e.g., tasks that took multiple attempts).

**Verification (TDD):**
- Verify the markdown file manually or via `cat` to ensure the "Specification Retrospective" section explicitly mentions `tasks.md` and bug analysis.

- [x] T1.2: [DOCS] Update Constitution Impact Analysis
**Target:** `src/internal/agent/kit/instructions/archive.md`
**Context:** [REQ-2]

**Action Steps:**
- Update `### 6. Constitution Impact Analysis` section (renumber if necessary to maintain sequential flow, likely to `4.`).
- Add instructions to evaluate if the challenges and bugs encountered indicate a missing rule or lack of clarity in the Constitution.
- Instruct the agent to formulate a rule to prevent the same error from repeating.

**Verification (TDD):**
- Verify the markdown file to ensure the "Constitution Impact Analysis" section includes checking for missing rules based on encountered bugs.

- [x] T1.3: [DOCS] Update Information Gathering Prompt and Formatting
**Target:** `src/internal/agent/kit/instructions/archive.md`
**Context:** [REQ-3]

**Action Steps:**
- Update the `### 8. Information Gathering (Tool Discovery) & Constitution Update` section (renumber if necessary, likely to `6.`).
- Modify the suggested `Ask the user:` prompt to explicitly mention the encountered challenge and ask if the project's Constitution should be updated to prevent it from repeating.
- Renumber subsequent sections (Archival Execution to `7.`, Verification & Handoff to `8.`).

**Verification (TDD):**
- Verify the entire markdown file has sequential, correct numbering.
- Verify the prompt string clearly asks about updating the constitution based on challenges/bugs.

## 3. Pre-emptive Mitigations
- **Risk:** Agent misses steps due to confusing numbering -> **Mitigation:** Ensure all section numbering in `archive.md` is strictly sequential and clean after making modifications.