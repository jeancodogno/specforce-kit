---
slug: 20261023-1159-behavioral-kit-manifest-tests
lens: Backend-heavy
---

# Implementation Roadmap: Behavioral Kit Manifest Tests

## 1. Execution Strategy
- **Gravity Order:** Section-slicing helpers (fail-fast) -> behavioral checks for the implementation blueprint -> behavioral checks for the archival instructions -> in-memory mutation self-test. Everything lives in one new file, `src/internal/agent/kit_behavior_test.go`, so all tasks are sequential. No blueprint text and no existing test is modified.
- **Conventions:** Plain Go `testing`, no new dependencies, functions up to 50 lines and cyclomatic complexity up to 10. Blueprints are read from `kit/skills/spf-implement/SKILL.yaml`, `kit/skills/spf-archive/SKILL.yaml` and `kit/instructions/archive.md` (paths relative to the package directory, as in `kit_manifests_test.go`).

## 2. Tasks

### Phase 1: Section helpers

- [x] T1.1: [RED] Tests for the section-slicing helper
**Target:** `src/internal/agent/kit_behavior_test.go`
**Context:** [US-1]

**Action Steps:**
- Create the file with `package agent` and imports `os`, `strings`, `testing`.
- Add `TestSectionHelper` with table cases: an existing heading returns the text up to the next heading of the same or higher level; the last heading returns text to end of input; a missing heading returns an error whose message contains the heading text.
- Add a case where a phrase exists only in a different section and assert the section-scoped lookup reports it as absent.
- Reference the not-yet-written helper `sectionOf(content, heading string) (string, error)` so the test fails to compile (RED).

**Acceptance Check:**
`cd src && go test ./internal/agent/ -run TestSectionHelper` exits non-zero with an "undefined: sectionOf" build error.

- [x] T1.2: [GREEN] Implement section helpers and the blueprint loader
**Target:** `src/internal/agent/kit_behavior_test.go`
**Context:** [US-1]

**Action Steps:**
- Implement `sectionOf(content, heading string)`: find the first line containing `heading` that starts (after trimming spaces) with `#`, compute its `#` level, and return the lines until the next heading line with level <= that level (or end); return `fmt.Errorf("section %q not found", heading)` when absent.
- Implement `mustRead(t *testing.T, path string) string` calling `t.Helper()` and `t.Fatalf` on read errors.
- Implement `requireIn(t *testing.T, path, content, heading string, phrases ...string)` which slices the section (failing with blueprint path and heading on error) and calls `t.Errorf` naming path, heading and phrase for each phrase missing from the slice.
- Implement `forbidIn(t *testing.T, path, content string, tokens ...string)` which calls `t.Errorf` for each token present anywhere in content.

**Acceptance Check:**
`cd src && go test ./internal/agent/ -run TestSectionHelper` exits 0.

### Phase 2: Implementation blueprint behavior

- [x] T2.1: [RED] Behavioral checks for the implementation blueprint
**Target:** `src/internal/agent/kit_behavior_test.go`
**Context:** [US-2]

**Action Steps:**
- Add `checkImplementBehavior(t *testing.T, content string)` that runs the sub-checks below against the given text (so the mutation self-test can reuse it) and reference it from `TestImplementBehavior`, which loads `kit/skills/spf-implement/SKILL.yaml`; leave the function body empty in this task; the assertions are filled in by T2.2.
- Define in the same function, via `requireIn`/`forbidIn` calls (written here, executed when the body is filled in T2.2), these rules: (a) in the `Free-Form Batch Narration` bullet the index of `Dispatch the worker FIRST` is lower than the index of `then write 1 to 3`; (b) section `#### B. Execute the 4-Tier QA Verification Protocol` contains `Maximum 2 cycles` and section `### 3. Verification & Batch Resolution` contains `up to 2 iterations`; (c) section `#### A. Spawn the qa-tester` contains `NEVER the worker reports`; (d) section `## 2. NON-NEGOTIABLE WORKER GUARDRAILS` contains `You never ask the user directly` and `BLOCKED: <question> | <options> | <recommendation>`; (e) forbidden anywhere: `[BATCH`, `COMPLETED]`, the rocket emoji, `MANDATORY LOG`, `Recommended Model Tier`, `Iniciando Lote`, `Delegando para perfil de`, `unlimited`, `indefinitely`.
- Write the test `TestImplementBehavior` as a call to `checkImplementBehavior(t, mustRead(t, implementPath))` where `implementPath` is a package-level constant.

**Acceptance Check:**
`cd src && go test ./internal/agent/ -run TestImplementBehavior -v` shows `--- PASS: TestImplementBehavior` (stub body; confirms wiring compiles).

- [x] T2.2: [GREEN] Fill in the implementation blueprint checks and make them pass on the real text
**Target:** `src/internal/agent/kit_behavior_test.go`
**Context:** [US-1], [US-2]

**Action Steps:**
- Implement rule (a) with a helper `requireOrder(t, path, content, heading, first, second string)` that slices the section and errors unless `strings.Index(first)` is >= 0 and lower than `strings.Index(second)`; for the narration rule slice the section `### 2. Batch Delegation & Dynamic Worker Spawning` and use the two phrases.
- Implement rules (b), (c), (d) with `requireIn`, using the headings named in T2.1 for the cycle limit, the qa-tester independence and the iteration limit, and `## 2. NON-NEGOTIABLE WORKER GUARDRAILS` for the BLOCKED protocol (verify the actual heading text against the file; the `## ` headings sit inside the Mission Brief block, so slicing stops at the next `## ` heading).
- Implement rule (e) with `forbidIn`.
- Assert the line count of the file does not exceed the existing ceiling by reusing the same constant value (222) as `TestImplementRoleChartersAndResume`.

**Acceptance Check:**
`cd src && go test ./internal/agent/ -run TestImplementBehavior -v` exits 0 with `--- PASS: TestImplementBehavior`.

### Phase 3: Archival behavior

- [x] T3.1: [RED] Archival next-steps checks
**Target:** `src/internal/agent/kit_behavior_test.go`
**Context:** [US-3]

**Action Steps:**
- Add `checkArchiveBehavior(t *testing.T, content string)` containing `requireIn` calls against heading `### Suggested Next Steps & Follow-up Specs` for the phrases `from: <finding or unverifiable criterion>`, `At most 3 steps` and `if no findings were persisted`; leave it referencing a non-existent constant `archiveInstructionsPath` so the build fails (RED).
- Add `TestArchiveNextStepsBehavior` calling `checkArchiveBehavior(t, mustRead(t, archiveInstructionsPath))`.
- Add an assertion in the test that a content string lacking the heading yields a failure through `sectionOf` (error case from the requirements).

**Acceptance Check:**
`cd src && go test ./internal/agent/ -run TestArchiveNextStepsBehavior` exits non-zero with an "undefined: archiveInstructionsPath" build error.

- [x] T3.2: [GREEN] Define archive path and pass on the real instructions
**Target:** `src/internal/agent/kit_behavior_test.go`
**Context:** [US-3]

**Action Steps:**
- Declare package-level constants `implementPath = "kit/skills/spf-implement/SKILL.yaml"` and `archiveInstructionsPath = "kit/instructions/archive.md"` (declare `implementPath` here if T2.1 used a literal).
- Confirm the heading text used in `checkArchiveBehavior` matches `kit/instructions/archive.md` exactly (the heading ends with a colon in the file; slicing matches by substring so the colon is not required).
- Run the test and adjust only the test side, never the instruction file.

**Acceptance Check:**
`cd src && go test ./internal/agent/ -run TestArchiveNextStepsBehavior -v` exits 0 with `--- PASS: TestArchiveNextStepsBehavior`.

### Phase 4: Mutation self-test

- [x] T4.1: [RED] Mutation self-test that expects detection
**Target:** `src/internal/agent/kit_behavior_test.go`
**Context:** [US-4]

**Action Steps:**
- Define `type mutation struct { name, file string; apply func(string) (string, bool); check func(*testing.T, string) }` and a table `behaviorMutations` of the 6 inversions: `narrate-before-dispatch` (swap the two narration phrases), `unlimited-correction-cycles` (replace `Maximum 2 cycles` with `Unlimited cycles`), `reports-to-qa-tester` (replace `NEVER the worker reports` with `ALWAYS the worker reports`), `delete-never-ask-user-rule` (remove the sentence `You never ask the user directly. `), `readd-batch-chain` (append `[BATCH x/y COMPLETED]` to the narration bullet), `archive-not-from-findings` (replace `from: <finding or unverifiable criterion>` with an empty string); each `apply` returns false when its target text is absent.
- Add `TestBehaviorMutationsAreDetected` that, for each mutation, runs `t.Run(m.name, ...)`: reads the file, applies the mutation to the in-memory copy, fails the subtest with "mutation target not found" when `apply` returns false, then runs `m.check` on a recording substitute for `*testing.T`.
- Implement the recording with a helper `failsOn(check func(reporter), ...) bool` that uses a small `recorder` struct implementing a minimal `reporter` interface (`Errorf`, `Fatalf`, `Helper`), and change `requireIn`, `requireOrder`, `forbidIn`, `checkImplementBehavior` and `checkArchiveBehavior` to accept that interface instead of `*testing.T` (an interface with those three methods, which `*testing.T` satisfies).

**Acceptance Check:**
`cd src && go test ./internal/agent/ -run TestBehaviorMutationsAreDetected` exits non-zero with a build error about the missing `reporter` or `recorder` type.

- [x] T4.2: [GREEN] Implement the recorder and prove all 6 inversions are detected
**Target:** `src/internal/agent/kit_behavior_test.go`
**Context:** [US-4]

**Action Steps:**
- Define `type reporter interface { Helper(); Errorf(string, ...any); Fatalf(string, ...any) }` and `type recorder struct { failed bool }` with methods that set `failed = true` (`Fatalf` also sets it and returns; checks must return after calling it).
- Make `failsOn` build a `recorder`, call the check with the mutated text and return `recorder.failed`; make the subtest error with `mutation %q was not detected` when it returns false.
- Make `sectionOf` errors flow through `reporter.Fatalf` in `requireIn` and `requireOrder` so a missing heading is a recorded failure.
- Add a final subtest `unmutated-text-passes` that runs both checks on the unmodified content with a recorder and requires `failed == false`, guarding against checks that always fail.
- Ensure nothing writes to disk: only `os.ReadFile` is used.

**Acceptance Check:**
`cd src && go test ./internal/agent/ -run TestBehaviorMutationsAreDetected -v` exits 0 and lists 6 `PASS` subtests plus `unmutated-text-passes`.

### Phase 5: Package-wide confirmation

- [x] T5.1: [VERIFY] Whole package, formatting and unchanged blueprints
**Target:** `src/internal/agent/`
**Context:** [US-1], [US-2], [US-3], [US-4]

**Action Steps:**
- Run `gofmt -l src/internal/agent/` and fix formatting in the new file only.
- Run the whole agent package tests and confirm all pre-existing tests still pass.
- Run `git status --short src/internal/agent/kit` and confirm no blueprint or instruction file is modified.
- Check `wc -l src/internal/agent/kit/skills/spf-implement/SKILL.yaml` is unchanged at 213.

**Acceptance Check:**
`cd src && test -z "$(gofmt -l internal/agent/)" && go test ./internal/agent/... && test -z "$(git status --short internal/agent/kit)"` exits 0.

## 3. Pre-emptive Mitigations
- **Risk:** A heading in a blueprint is renamed later and the checks silently stop covering a rule -> **Mitigation:** `sectionOf` fails with the missing heading's name, and the mutation self-test fails when a mutation target no longer exists.
- **Risk:** Checks that always fail would make the self-test pass vacuously -> **Mitigation:** the `unmutated-text-passes` subtest requires the real text to pass the same checks.
