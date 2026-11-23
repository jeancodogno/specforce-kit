---
slug: 20261123-1656-lean-spec-templates
lens: Backend-heavy
---

# Implementation Roadmap: Lean, Reviewable Spec Templates

## 1. Execution Strategy

- Reference layout: `.specforce/scratch/proto-mesa-sender-avatars/requirements.md` and `tasks.md` (Portuguese). The templates are written in English from it. The scratch folder is disposable, so the rules below are the source of truth.
- Section titles and field labels read by tools stay exact: `## 2. Out of Scope`, `## 1. Execution Strategy`, `## N. Pre-emptive Mitigations`, `## N. Black-box Scenarios`, `**Target:**`, `**Context:**`, `**Constraints:**`, `**Acceptance Check:**`, `**Parallel With:**`.
- The task parser needs `- [ ] T{n}.{n}:` with no bold, and the acceptance check text on the line right after its label with no blank line inside. Indented fields and `---` separators were proven valid with `specforce spec status` and `specforce implementation update`.
- Requirement links take the form `**Context:** [US-1] a, b · [US-2] c`. Until T1.2 lands, the coverage audit only sees the first requirement of a link.
- New kit assertions live in one new test function, `TestKitManifests_LeanReviewableTemplates`. The existing assertions on removed text are updated in the same task that adds the new ones.

## 2. Tasks

### Phase 1: Coverage audit recognizes every linked requirement

- [x] T1.1: Create coverage-audit tests for links that cite several requirements [RED]

  **Target:** `src/internal/spec/auditor_test.go`

  **Context:** [US-8] a, b, c, d

  **Acceptance Check:**
  `go test ./src/internal/spec/ -run TestDeterministicCheck` fails on the new cases for a second requirement in a link and for a requirement followed by scenario letters.

---

- [x] T1.2: Modify the coverage audit to collect every requirement cited in each requirement link

  **Target:** `src/internal/spec/auditor.go`

  **Context:** [US-8]

  **Constraints:**
  - Only `**Context:**` lines count. Matching any `[US-x]` in the roadmap would mark requirements covered by a title or constraint mention (US-8.d).

  **Acceptance Check:**
  `go test ./src/internal/spec/ -run TestDeterministicCheck` passes, including the T1.1 cases.

### Phase 2: Requirements blueprints produce a brief and compact scenarios

- [x] T2.1: Create kit assertions for the lean requirements and bug requirements blueprints [RED]

  **Target:** `src/internal/agent/kit_manifests_test.go`

  **Context:** [US-1] a, b, c · [US-2] a, b, c · [US-3] a, b, c, d · [US-7] a, b, c

  **Acceptance Check:**
  `go test ./src/internal/agent/ -run TestKitManifests_LeanReviewableTemplates` fails: the blueprints still lack the Review Brief, the "situation → result" scenarios and General Rules, and still contain `UI/UX Specifics` and the mandatory `[Performance]` tag.

---

- [x] T2.2: Modify the requirements blueprint to open with a Review Brief, use lettered one-line scenarios and state feature-wide rules once

  **Target:** `src/internal/agent/artifacts/spec/requirements.yaml`

  **Context:** [US-1] · [US-2] · [US-3]

  **Constraints:**
  - Keep the zero-technical-specification rule and the `lens` frontmatter. Removing them lets implementation details back into requirements (BR-KIT-02).

  **Acceptance Check:**
  `go test ./src/internal/agent/ -run 'TestKitManifests_LeanReviewableTemplates/Requirements'` passes.

---

- [x] T2.3: Modify the bug requirements blueprint to add the Review Brief and lettered regression scenarios

  **Target:** `src/internal/agent/artifacts/spec/bug-requirements.yaml`

  **Context:** [US-7]

  **Constraints:**
  - Keep the words "root cause", "RED" and "halt". `TestKitManifests_ExecutableCriteriaAndArchivePreflight/BugBlueprint` asserts them.

  **Acceptance Check:**
  `go test ./src/internal/agent/ -run 'TestKitManifests_LeanReviewableTemplates/BugRequirements|TestKitManifests_ExecutableCriteriaAndArchivePreflight'` passes.

### Phase 3: Tasks blueprint produces readable, lean roadmaps

- [x] T3.1: Create kit assertions for the readable task layout, lean constraints and test-first neutrality [RED]

  **Target:** `src/internal/agent/kit_manifests_test.go`

  **Context:** [US-4] a, b, c, d, e · [US-5] a, b, c, d, e · [US-6] a, b, c

  **Acceptance Check:**
  `go test ./src/internal/agent/ -run 'TestKitManifests_LeanReviewableTemplates/Tasks'` fails: the blueprint still has `Gravity Order` and `[SCAFFOLD/CLI]`, and lacks the action verbs, the `---` separator rule and the constraint exclusion list. The planning skill does not yet say that the `[RED]` token ends the title.

---

- [x] T3.2: Modify the tasks blueprint for verb-first titles, indented fields, task separators and lean constraints

  **Target:** `src/internal/agent/artifacts/spec/tasks.yaml`

  **Context:** [US-4] · [US-5] · [US-6]

  **Constraints:**
  - Keep the rules for black-box scenarios by size, verification-only tasks and bug-first regression word for word where tests assert them (`TestKitManifests_ExecutableCriteriaAndArchivePreflight/TaskBlueprint`, `testLeanTasksBlueprint`).
  - The template example shows `**Constraints:**` on the second task only. The example on the first task is what agents copy into every task today.

  **Acceptance Check:**
  `go test ./src/internal/agent/ -run 'TestKitManifests_LeanReviewableTemplates|TestKitManifests_ExecutableCriteriaAndArchivePreflight|TestLeanBlueprintsManifest'` passes.

---

- [x] T3.3: Modify the planning skill so a RED task title ends with the token after its action verb

  **Target:** `src/internal/agent/kit/skills/spf-spec/SKILL.yaml`

  **Context:** [US-6] b

  **Constraints:**
  - Change only the sentence about the `[RED]` token. `TestKitManifests_QAPersistenceAndRedMarker` asserts that the token stays in this file.

  **Acceptance Check:**
  `go test ./src/internal/agent/ -run 'TestKitManifests_QAPersistenceAndRedMarker|TestKitManifests_LeanReviewableTemplates'` passes.
