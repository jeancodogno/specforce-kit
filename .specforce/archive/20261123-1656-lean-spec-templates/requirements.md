---
slug: 20261123-1656-lean-spec-templates
lens: Backend-heavy
---

# Feature: Lean, Reviewable Spec Templates

## 1. Context & Value
Real specifications reach 8,000 to 16,000 words, so a human review takes close to an hour. Most of the volume is mandatory boilerplate and the same rule restated in three artifacts. The templates must keep what an agent needs to execute and let a human approve a spec by reading a short brief and compact scenarios.

## 2. Out of Scope (Anti-Goals)
- Do not change the design or bug design templates; they get their own specification after a prototype.
- Do not change how tasks are parsed, tracked or marked finished.
- Do not prescribe or forbid test-first task pairs; that stays a per-project choice.
- Do not migrate or rewrite existing specifications or archived ones.
- Do not translate the templates; they stay in English.

## 3. Acceptance Criteria (BDD)

### [US-1] Review Brief at the top of requirements
**User Story:** AS A spec reviewer, I WANT a short brief at the top of the requirements, SO THAT I can approve the spec without reading every scenario.

**Scenarios:**
1. **[Happy Path]** GIVEN a new feature specification WHEN its requirements are written THEN they open with a Review Brief holding three items: what changes, the decisions taken, and what to review closely.
2. **[Edge Case]** GIVEN a decision-heavy feature WHEN the brief is written THEN it stays within about 8 lines, each item has at most 3 short bullets, and the detail stays in the requirements below.
3. **[Edge Case]** GIVEN a feature with no decision beyond the request itself WHEN the brief is written THEN the decisions item says so in one line instead of inventing decisions.

**UI/UX Specifics:**
- **View/Component:** First section of the requirements document, before the context.
- **Feedback Logic:** Bullets name the requirement and scenario identifiers the reviewer should open.
- **Keybindings:** Not applicable to a document.

**Technical Constraints (NFR):**
- **[Performance]:** The brief is readable in under one minute (about 8 lines).
- **[Integrity]:** Every decision in the brief is also reflected in a requirement or scenario.

### [US-2] One-line acceptance scenarios with identifiers
**User Story:** AS A spec reviewer, I WANT each scenario on one line with a letter identifier, SO THAT I can scan behavior quickly and tasks can cite it.

**Scenarios:**
1. **[Happy Path]** GIVEN a requirement WHEN its scenarios are written THEN each is one bullet `a.`, `b.`, … in the form "situation → result" under a one-line user story.
2. **[Happy Path]** GIVEN two scenarios with the same result WHEN they are written THEN they are merged into one scenario.
3. **[Edge Case]** GIVEN a requirement with no real failure case WHEN its scenarios are written THEN no failure scenario is invented.
4. **[Edge Case]** GIVEN a requirement with a real failure case WHEN its scenarios are written THEN the failure is stated with its observable result.

**UI/UX Specifics:**
- **View/Component:** Requirement block: heading, one-line story, lettered scenario list.
- **Feedback Logic:** A scenario is referenced elsewhere as the requirement identifier plus its letter (for example US-3.b).
- **Keybindings:** Not applicable to a document.

**Technical Constraints (NFR):**
- **[Performance]:** A requirement block fits in about 10 lines in the common case.
- **[Integrity]:** Every scenario stays testable and keeps one observable result.

### [US-3] No per-requirement boilerplate
**User Story:** AS A spec reviewer, I WANT feature-wide rules stated once, SO THAT I do not read repeated or empty sections in every requirement.

**Scenarios:**
1. **[Happy Path]** GIVEN a feature WHEN requirements are written THEN interface, domain and quality rules that apply to the whole feature appear once in a single General Rules section.
2. **[Happy Path]** GIVEN a requirement that deviates from the general rules WHEN it is written THEN it carries a one-line constraint with only that deviation.
3. **[Edge Case]** GIVEN a feature without a user interface WHEN its requirements are written THEN the interface rules are omitted entirely.
4. **[Edge Case]** GIVEN a feature without a real, measurable quality target WHEN its requirements are written THEN no quality target is invented, and vague terms such as "optimized" remain forbidden.

**UI/UX Specifics:**
- **View/Component:** General Rules section after the requirements, grouped as Domain, Interface and Quality.
- **Feedback Logic:** No line in the document states that something does not change or does not apply.
- **Keybindings:** Not applicable to a document.

**Technical Constraints (NFR):**
- **[Performance]:** Zero mandatory per-requirement sections besides the story and scenarios.
- **[Integrity]:** Every quality line carries a concrete value.

### [US-4] Readable task layout
**User Story:** AS A spec reviewer, I WANT each task to say what it does and to look separate from the next one, SO THAT I can review the roadmap without opening every task.

**Scenarios:**
1. **[Happy Path]** GIVEN a task WHEN it is written THEN its title starts with an action verb (Create, Modify, Remove, Move or Rename) followed by the target and the outcome.
2. **[Happy Path]** GIVEN a phase WHEN it is written THEN its name states the outcome it delivers, not only the layer it touches.
3. **[Happy Path]** GIVEN a task WHEN it is written THEN each of its fields is a separate paragraph indented under the task, and consecutive tasks of a phase are separated by a horizontal rule.
4. **[Happy Path]** GIVEN a task that covers scenarios WHEN its requirement link is written THEN it cites the requirement and the scenario letters it covers, for one or several requirements.
5. **[Edge Case]** GIVEN a title written in bold or a separator without a blank line before it WHEN the template is followed THEN neither appears, because both break task tracking or rendering.

**UI/UX Specifics:**
- **View/Component:** Task card: title line, indented fields, horizontal rule.
- **Feedback Logic:** Rendered and raw Markdown both show where a task starts and ends.
- **Keybindings:** Not applicable to a document.

**Technical Constraints (NFR):**
- **[Performance]:** Reading the task titles alone tells what every task creates, modifies or removes.
- **[Integrity]:** Task status updates keep working on the new layout.

### [US-5] Lean task constraints and execution strategy
**User Story:** AS A planning agent, I WANT clear rules for what is and is not a constraint, SO THAT tasks only carry what the implementer cannot deduce.

**Scenarios:**
1. **[Happy Path]** GIVEN a task WHEN its constraints are written THEN each one states a pitfall, a convention deviation or a decision, and what goes wrong if it is ignored.
2. **[Edge Case]** GIVEN a candidate constraint that repeats a requirement, repeats another task's constraint, describes the obvious implementation, or defines a data contract WHEN the task is written THEN it is left out of the task.
3. **[Happy Path]** GIVEN technical anchors or a decision that applies to several tasks WHEN the roadmap is written THEN it is stated once in the execution strategy and not repeated in tasks.
4. **[Edge Case]** GIVEN a risk already handled by a task constraint WHEN the roadmap is written THEN it is not repeated as a mitigation.
5. **[Happy Path]** GIVEN a roadmap WHEN the execution strategy is written THEN it no longer requires a narrative of the task order, which the phases already show.

**UI/UX Specifics:**
- **View/Component:** Execution strategy section holding anchors and cross-task decisions.
- **Feedback Logic:** The template example shows the constraints block on one task only, marked optional.
- **Keybindings:** Not applicable to a document.

**Technical Constraints (NFR):**
- **[Performance]:** A task without non-deducible information has no constraints block.
- **[Integrity]:** Black-box scenario rules by size remain unchanged.

### [US-6] Test-first neutrality
**User Story:** AS A project maintainer, I WANT the templates to stay neutral about test-first task pairs, SO THAT each project decides whether to use them.

**Scenarios:**
1. **[Happy Path]** GIVEN a project with no test-first instruction WHEN a roadmap is written THEN tasks are not split into test and code pairs because of the template.
2. **[Happy Path]** GIVEN a project whose instructions ask for test-first pairs WHEN a roadmap is written THEN each test task title keeps the action verb and ends with the `[RED]` token.
3. **[Edge Case]** GIVEN a bug specification WHEN its roadmap is written THEN its first task still creates the regression test that fails before the fix.

**UI/UX Specifics:**
- **View/Component:** Task title.
- **Feedback Logic:** The `[RED]` token is the last element of the title.
- **Keybindings:** Not applicable to a document.

**Technical Constraints (NFR):**
- **[Performance]:** Finishing a `[RED]` task keeps skipping the per-task hooks.
- **[Integrity]:** Templates contain no instruction that splits test and code tasks for features.

### [US-7] Lean bug requirements
**User Story:** AS A spec reviewer, I WANT bug specifications to follow the same brief and scenario style, SO THAT bug reviews are as quick as feature reviews.

**Scenarios:**
1. **[Happy Path]** GIVEN a bug specification WHEN its requirements are written THEN they open with a Review Brief stating the fix, the decisions and what to review closely.
2. **[Happy Path]** GIVEN a fix criterion WHEN it is written THEN its regression scenarios are lettered one-line "situation → result" bullets.
3. **[Edge Case]** GIVEN a bug WHEN its requirements are written THEN evidence, reproduction steps and the root cause remain mandatory.

**UI/UX Specifics:**
- **View/Component:** Bug requirements document.
- **Feedback Logic:** Same brief layout as feature requirements.
- **Keybindings:** Not applicable to a document.

**Technical Constraints (NFR):**
- **[Performance]:** No mandatory quality tag in the bug constraints section.
- **[Integrity]:** Each fix criterion still maps to the first regression test task.

### [US-8] Coverage audit sees every linked requirement
**User Story:** AS A planning agent, I WANT the coverage audit to recognize every requirement cited in a task, SO THAT a task covering several requirements does not raise a false missing-task error.

**Scenarios:**
1. **[Happy Path]** GIVEN a task whose requirement link cites US-1 and US-2 WHEN the coverage audit runs THEN neither requirement is reported as missing.
2. **[Happy Path]** GIVEN a requirement link that also cites scenario letters WHEN the coverage audit runs THEN the requirement is recognized.
3. **[Edge Case]** GIVEN a requirement that no task cites WHEN the coverage audit runs THEN it is still reported as missing.
4. **[Edge Case]** GIVEN a requirement mentioned only in a task's constraints or title WHEN the coverage audit runs THEN it does not count as covered.

**UI/UX Specifics:**
- **View/Component:** Specification status output.
- **Feedback Logic:** Missing-task errors name the uncovered requirement, as today.
- **Keybindings:** No change to commands.

**Technical Constraints (NFR):**
- **[Performance]:** The audit stays a single pass over the roadmap.
- **[Integrity]:** Requirements listed under Out of Scope remain excluded from the audit.

## 4. Business Invariants
- Existing specifications, including archived ones, stay valid without edits.
- Section names that tools read keep their exact titles: Out of Scope, Execution Strategy, Pre-emptive Mitigations and Black-box Scenarios.
- Task fields that tools read keep their exact labels: Target, Context, Constraints, Acceptance Check and Parallel With.
- Requirements keep containing no implementation details.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** A medium specification written from the new templates is at least 40% shorter than the same specification written from the current ones, measured on the mesa-sender-avatars prototype (2,454 → about 1,050 words).
- **[Reliability]:** Specification status stays valid for specs written in both the old and new layouts.
- **[Maintainability]:** Kit validation tests assert the new template rules and reject the removed ones.
