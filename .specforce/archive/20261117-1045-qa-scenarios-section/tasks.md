---
slug: 20261117-1045-qa-scenarios-section
lens: Backend-heavy
---

# Implementation Roadmap: Black-box Scenarios as QA Input

## 1. Execution Strategy
- **Gravity Order:** Guard and manifest tests first (parser invariant + failing blueprint directives) -> roadmap blueprint, implementation skill and planning skill in parallel -> regenerated installed skill copies. No task runs the full suite: that is left to the final QA and the all-tasks hooks, as this specification itself requires.

## 2. Tasks

### Phase 1: Tests

- [x] T1.1: [CODE] Guard test: scenario section is not counted as tasks
**Target:** `src/internal/spec/tasks_test.go`
**Context:** [US-1], [US-4]

**Constraints:**
- The fixture roadmap places a `## 3. Black-box Scenarios` section with `S1`/`S2` plain list items (no checkbox, no `T` identifier) after the last task; compare task count and progress against the same roadmap without the section.
- A second fixture keeps a legacy last phase of `[VERIFY]` tasks and must stay valid with those tasks counted as ordinary tasks.
- This is a characterization test of existing parser behavior: it is expected to pass without touching `tasks.go`; if it fails, stop and report instead of changing the parser (out of scope).

**Acceptance Check:**
`go test ./src/internal/spec/ -run 'Scenario|LegacyVerify' -count=1` exits 0 and the new tests are listed with `-v`.

- [x] T1.2: [RED] Manifest tests for scenario section, no verification-only tasks and QA scenario rows
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-1], [US-2], [US-3]

**Constraints:**
- Replace the `"Black-box Verification"` directive in the TaskBlueprint subtest with `"Black-box Scenarios"`, and add an assertion that `artifacts/spec/tasks.yaml` no longer contains `[VERIFY]` nor the phase heading `Phase {N}: Black-box Verification`.
- Assert the roadmap blueprint forbids verification-only tasks, the implementation skill directs one matrix row per scenario, and the planning skill's coherence audit removes verification-only tasks; pick short, stable phrases and reuse them verbatim in T2.1-T2.3.

**Acceptance Check:**
`go test ./src/internal/agent/ -run TestKitManifests -count=1` exits non-zero, naming the missing scenario, verification-only and matrix-row directives.

### Phase 2: Blueprints

- [x] T2.1: [GREEN] Roadmap blueprint: scenario section and verification-only ban
**Target:** `src/internal/agent/artifacts/spec/tasks.yaml`
**Context:** [US-1], [US-2], [US-4]
**Parallel With:** T2.2, T2.3

**Constraints:**
- Rewrite instruction rule 9 per the modified `BR-KIT-20` in `deltas/agent-kit.md`: mandatory for `large`/`complex`, optional for `medium` with a runnable artifact, omitted for `small`; scenarios use `S{n}`, never a checkbox or task identifier.
- Add the verification-only task ban as its own instruction rule.
- In the template, replace the verification phase with a `## 3. Black-box Scenarios` section after the task phases and renumber Pre-emptive Mitigations to `## 4.`; keep the YAML valid.

**Acceptance Check:**
`go test ./src/internal/agent/ -run 'TestKitManifests_ExecutableCriteriaAndArchivePreflight/TaskBlueprint' -count=1` exits 0.

- [x] T2.2: [GREEN] Implementation skill: qa-tester runs every scenario
**Target:** `src/internal/agent/kit/skills/spf-implement/SKILL.yaml`
**Context:** [US-3]
**Parallel With:** T2.1, T2.3

**Constraints:**
- Edit only Step 4 (qa-tester inputs, Tier 2 and the Traceability Matrix paragraph) per the added `NEW:qa-runs-scenarios` item in `deltas/agent-kit.md`; the QA record JSON shape is unchanged.
- Roadmaps without the section keep deriving black-box checks from the requirements.

**Acceptance Check:**
`go test ./src/internal/agent/ -run 'TestKitManifests' -count=1 -v 2>&1 | grep -E '^(--- FAIL|FAIL).*(Implement|implement)'` prints nothing.

- [x] T2.3: [GREEN] Planning skill: coherence audit removes verification-only tasks
**Target:** `src/internal/agent/kit/skills/spf-spec/SKILL.yaml`
**Context:** [US-1], [US-2]
**Parallel With:** T2.1, T2.2

**Constraints:**
- Extend Step 3 (In-Session Quality Gate) so the audit removes verification-only tasks, converts runnable end-to-end checks into scenarios when the size allows the section, checks each scenario links to a requirement, and names the removals in the planning summary.

**Acceptance Check:**
`go test ./src/internal/agent/ -run TestKitManifests -count=1` exits 0.

### Phase 3: Installed Copies

- [x] T3.1: [GREEN] Refresh installed skill copies
**Target:** `.claude/skills`, `.agents/skills`
**Context:** [US-1], [US-2], [US-3]

**Constraints:**
- Regenerate through the repository's existing init flow using a freshly built binary, never by hand; only the planning and implementation skills may change.

**Acceptance Check:**
`grep -c "Black-box Scenarios" .claude/skills/spf-spec/SKILL.md .agents/skills/spf-spec/SKILL.md` reports at least 1 per file.

### Phase 4: Mitigations Heading Recognition

- [x] T4.1: [RED] Mitigations extracted regardless of section number
**Target:** `src/internal/spec/implementation_test.go`
**Context:** [US-5]

**Constraints:**
- Cover the three cases of [US-5] and the `NEW:mitigations-by-title` delta item in `deltas/spec-management.md`: a roadmap with a scenario section then `## 4. Pre-emptive Mitigations`, the legacy `## 3.` numbering (the existing test must keep passing untouched), and a roadmap with mitigations followed by another `## ` section whose text must not leak in. Test names must match `Mitigations`.

**Acceptance Check:**
`go test ./src/internal/spec/ -run Mitigations -count=1` exits non-zero, failing only on the `## 4.` case.

- [x] T4.2: [GREEN] Identify the mitigations section by title
**Target:** `src/internal/spec/implementation.go`
**Context:** [US-5]

**Constraints:**
- Change only the mitigations section match so any section number is accepted; the execution-strategy match and the task parser stay untouched.

**Acceptance Check:**
`go test ./src/internal/spec/ -run Mitigations -count=1` exits 0.

## 3. Black-box Scenarios
- **S1** [US-1], [US-2] The compiled CLI serves the new roadmap blueprint.
  - **Command:** `make build && ./specforce spec artifact feature-tasks --json | grep -c "Black-box Scenarios"`
  - **Expected:** exit 0 and a count of at least 1; the same output piped to `grep -c "\[VERIFY\]"` prints `0`.
- **S2** [US-1], [US-4] This very roadmap, which declares a scenario section, is valid and counts only its tasks.
  - **Command:** `./specforce spec status 20261117-1045-qa-scenarios-section --json && ./specforce implementation status 20261117-1045-qa-scenarios-section --json`
  - **Expected:** the first reports `"is_valid": true`; the second lists exactly 8 tasks (T1.1, T1.2, T2.1, T2.2, T2.3, T3.1, T4.1, T4.2) and no scenario identifier among them.
- **S3** [US-5] A roadmap built from the new template keeps its mitigations visible.
  - **Command:** `./specforce implementation status <slug> --json` on a scratch specification whose `tasks.md` has `## 3. Black-box Scenarios` followed by `## 4. Pre-emptive Mitigations` with one risk line.
  - **Expected:** `preemptive_mitigations` contains exactly that risk line and no scenario text.
