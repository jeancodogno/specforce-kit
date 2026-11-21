---
slug: 20261023-1159-behavioral-kit-manifest-tests
lens: Backend-heavy
---

# Feature: Behavioral Kit Manifest Tests

## 1. Context & Value
The automated checks guarding the implementation and archival skill blueprints only confirm that certain words still exist, so a change that inverts a rule while keeping its words goes unnoticed (an independent QA proved 6 such inversions all passed). This feature makes the checks verify what each rule actually says and where it applies, and proves that the checks themselves detect rule inversions, so that kit maintainers can trust a green result.

## 2. Out of Scope (Anti-Goals)
- Do not change the wording, structure or length of any skill blueprint or instruction text; the implementation blueprint keeps its current line ceiling.
- Do not add behavioral checks for the planning (`spf-spec`), discovery or constitution blueprints; they are left for a later spec.
- Do not attempt to execute an agent or simulate model behavior; only the text of the blueprints is inspected.
- Do not introduce new third-party dependencies.
- Do not remove or weaken any existing check.

## 3. Acceptance Criteria (BDD)

### [US-1] Rules are verified in the place where they apply
**User Story:** AS A kit maintainer, I WANT each critical rule to be verified inside the specific section of the blueprint it belongs to, SO THAT a rule moved, deleted or contradicted elsewhere is detected even if its words still appear in another section.

**Scenarios:**
1. **[Happy Path]** GIVEN the current blueprints WHEN the checks run THEN every anchored rule is found inside its expected section and the checks pass.
2. **[Edge Case]** GIVEN a section heading that no longer exists in the blueprint WHEN the checks run THEN the check fails with a message naming the missing section, never silently passing or skipping.
3. **[Edge Case]** GIVEN a rule whose words appear only in a different section than expected WHEN the checks run THEN the check fails.

**UI/UX Specifics:**
- **View/Component:** Test failure messages shown in the maintainer's terminal.
- **Feedback Logic:** Each failure names the blueprint, the section and the missing or forbidden rule.
- **Keybindings:** Not applicable (non-interactive).

**Technical Constraints (NFR):**
- **[Performance]:** Full suite for the agent kit completes in under 10 seconds.
- **[Integrity]:** Section lookup is by heading text and fails fast on a missing heading.
- **[Observability]:** Failure messages contain blueprint path, section name and rule text.

### [US-2] Behavioral rules of the implementation blueprint are pinned
**User Story:** AS A kit maintainer, I WANT the implementation blueprint's behavioral rules to be checked for order, limits, independence and prohibitions, SO THAT inverting any of them breaks the checks.

**Scenarios:**
1. **[Happy Path]** GIVEN the current implementation blueprint WHEN the checks run THEN all of the following hold: the worker is dispatched first and the progress narration comes after it; correction cycles are limited to a maximum of 2; the independent tester never receives the worker reports; workers never ask the user directly and must return the structured blocked report with question, options and recommendation; and no fixed bracketed status template or banner appears.
2. **[Edge Case]** GIVEN the blueprint is altered so narration comes before dispatching the worker WHEN the checks run THEN at least one check fails.
3. **[Edge Case]** GIVEN the blueprint is altered so correction cycles have no maximum (for example "unlimited" or "indefinitely") WHEN the checks run THEN at least one check fails.
4. **[Edge Case]** GIVEN the blueprint is altered so the independent tester receives the worker reports WHEN the checks run THEN at least one check fails.
5. **[Edge Case]** GIVEN the rule "workers never ask the user directly" is deleted while the blocked report keyword remains elsewhere WHEN the checks run THEN at least one check fails.
6. **[Edge Case]** GIVEN a bracketed batch-completed status chain is re-added WHEN the checks run THEN at least one check fails.

**UI/UX Specifics:**
- **View/Component:** Test failure messages shown in the maintainer's terminal.
- **Feedback Logic:** The failing check identifies which behavioral rule was violated.
- **Keybindings:** Not applicable (non-interactive).

**Technical Constraints (NFR):**
- **[Performance]:** Each check completes in under 100ms.
- **[Integrity]:** Forbidden-template checks cover the bracketed batch markers, the rocket emoji banner, the "mandatory log" header and the "recommended model tier" header, in addition to the two already forbidden fixed log phrases.
- **[Maintainability]:** Checks are plain standard-library tests with no new dependencies.

### [US-3] Archival next steps stay derived from findings
**User Story:** AS A kit maintainer, I WANT the archival instructions to be verified for deriving suggested next steps from reported findings with a bounded count, SO THAT archival cannot silently stop using findings.

**Scenarios:**
1. **[Happy Path]** GIVEN the current archival blueprint and instructions WHEN the checks run THEN the next-steps block requires each step to cite the finding or unverifiable criterion it comes from, limits the list to at most 3 steps, and states the fallback when no findings exist.
2. **[Edge Case]** GIVEN the "from finding" requirement is removed from the next-steps block WHEN the checks run THEN at least one check fails.
3. **[Edge Case]** GIVEN the limit of at most 3 steps is removed WHEN the checks run THEN at least one check fails.

**UI/UX Specifics:**
- **View/Component:** Test failure messages shown in the maintainer's terminal.
- **Feedback Logic:** The failure names the archival file and the missing derivation rule.
- **Keybindings:** Not applicable (non-interactive).

**Technical Constraints (NFR):**
- **[Performance]:** Checks complete in under 100ms.
- **[Integrity]:** The next-steps block is located by its heading, failing fast if absent.

### [US-4] The checks are proven to detect rule inversions
**User Story:** AS A kit maintainer, I WANT a self-test that applies known rule inversions to an in-memory copy of the blueprints and requires the checks to fail, SO THAT the checks cannot regress into word-presence checks again.

**Scenarios:**
1. **[Happy Path]** GIVEN the 6 known inversions (narrate before dispatching; unlimited correction cycles; worker reports handed to the independent tester; "workers never ask the user directly" deleted; bracketed batch-completed chain re-added; archival next steps no longer derived from findings) WHEN each is applied in memory to the blueprint text THEN the behavioral checks report at least one failure for each inversion.
2. **[Edge Case]** GIVEN an inversion whose target text can no longer be found in the blueprint WHEN the self-test runs THEN the self-test fails stating the inversion cannot be applied, instead of passing vacuously.
3. **[Edge Case]** GIVEN any inversion is applied in memory WHEN the self-test finishes THEN the files on disk are unchanged.

**UI/UX Specifics:**
- **View/Component:** Test results shown in the maintainer's terminal, one named result per inversion.
- **Feedback Logic:** An inversion not detected is reported by its name.
- **Keybindings:** Not applicable (non-interactive).

**Technical Constraints (NFR):**
- **[Performance]:** Self-test completes in under 2 seconds.
- **[Safety & Security]:** Blueprints are never written to disk by the self-test.
- **[Integrity]:** Every inversion is applied to a fresh copy, so inversions do not affect each other.

## 4. Business Invariants
- Blueprint texts remain unchanged by this feature; all existing tests continue to pass.
- A behavioral check never passes when its section or target text is missing.
- Each of the 6 known inversions is detected by at least one check.
- The implementation blueprint stays within its current line ceiling.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** The whole agent kit test package runs in under 10 seconds.
- **[Reliability]:** Fail-fast with clear messages; no check is skipped silently.
- **[Security]:** No network access and no writes outside temporary test locations.
- **[Maintainability]:** Standard Go formatting, functions within the project complexity limits, shared helpers reused by the checks and the self-test, new code kept in its own test file to avoid growing the existing manifest test file.
