---
slug: 20261023-1159-discovery-skill-rewrite
lens: Balanced full-stack
---

# Implementation Roadmap: Discovery Skill Rewrite (Format, Coherence, Evidence Sandbox, Proposal Skeleton)

## 1. Execution Strategy
- **Gravity Order:** Failing blueprint test -> rewritten blueprint (also fixes its description) -> failing guide test -> guide generator -> repository guide copy -> full suite. Each phase pairs its RED test with the GREEN change that satisfies it.

## 2. Tasks

### Phase 1: Discovery Blueprint

- [x] T1.1: [RED] Blueprint structure and evidence-sandbox test
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-1], [US-2], [US-3], [US-4], [US-5]

**Action Steps:**
- Add a test `TestDiscoveryBlueprintEvidenceSandbox` that reads `kit/skills/spf-discovery/SKILL.yaml`, unmarshals it, and fails if the `content` field is empty.
- Assert the raw file is not a single escaped line: the `content:` key must be followed by a block scalar indicator and the file must have more than 20 lines.
- Assert the file has at most 120 lines and that the `content` contains the strings `.specforce/scratch/`, `spf-probe`, `.gitignore`, `specforce spec init`, `--type bug`, and `Strong Opinions, Weakly Held`.
- Assert the `description` and `content` do not contain the phrase `Purely read-only`, and that the `content` contains no language-specific test-runner or extension names (`go test`, `pytest`, `jest`, `.go`, `.py`, `.ts`).
- Assert the `content` does not contain both the phrases `Micro-Turns over Monoliths` and `DO NOT give short, superficial`.

**Acceptance Check:**
`go test ./src/internal/agent/ -run TestDiscoveryBlueprintEvidenceSandbox` fails (non-zero exit) before the blueprint is rewritten.

- [x] T1.2: [CODE] Rewrite the discovery blueprint
**Target:** `src/internal/agent/kit/skills/spf-discovery/SKILL.yaml`
**Context:** [US-1], [US-2], [US-3], [US-4], [US-5]

**Action Steps:**
- Convert `content` to a `|` block scalar of at most 120 lines, keeping the "Proactive Suggestions & Technical Opinions" and "Strong Opinions, Weakly Held" directives and the `spf.discovery` name.
- Replace the "strict read-only sanctuary" paragraph with the evidence-sandbox rule: product code read-only; evidence created only as new files in `.specforce/scratch/` (created with a `.gitignore` containing `*` on first use); in-tree fallback only when the scratch area cannot prove the hypothesis, with `spf-probe` in the file name, path recorded, and removal verified through version-control status before ending, reporting leftovers; existing files never edited; scratch evidence stale in later sessions.
- Add the safe-execution rule: non-mutating commands, running the app, suites and browsers locally and capturing screenshots allowed; commands that alter tracked product files or have external or irreversible effects require asking first; broad searches may go to strictly read-only helpers.
- Resolve the length contradiction: short dense turns with an opinion; remove the "DO NOT give short, superficial" line.
- Add the formalization guidance (proposal sections: problem, verified findings with file references, decisions taken, rejected alternatives, open questions, suggested type and size, evidence references), the bug hand-off offering `specforce spec init <slug> --type bug`, the unchanged rule of no push toward specs on early and mid turns, and the decomposition recommendation for oversized ideas.
- Keep one compact example conversation and one compact wireframe; remove the others.
- Update `description` to say read-only for product code with disposable evidence allowed.

**Acceptance Check:**
`go test ./src/internal/agent/` exits 0 with the new test and the existing `TestEnhancedBlueprintsDirectives` passing.

### Phase 2: Governance Text Consistency

- [x] T2.1: [RED] Agent guide wording test
**Target:** `src/internal/project/agents_md_test.go`
**Context:** [US-6]

**Action Steps:**
- In `TestGenerateAgentsContent`, assert the generated content contains `read-only for product code` and `.specforce/scratch/`.
- Assert the generated content no longer contains `Purely read-only`.
- Assert the direct-edit prohibition sentence mentioning `replace` and `write_file` is still present.
- In `TestMergeAgentsContent`, add a case where existing managed content carries the old wording and assert the merged result contains the new wording exactly once.

**Acceptance Check:**
`go test ./src/internal/project/ -run 'TestGenerateAgentsContent|TestMergeAgentsContent'` fails (non-zero exit) before the generator changes.

- [x] T2.2: [CODE] Update the guide generator
**Target:** `src/internal/project/agents_md.go`
**Context:** [US-6]

**Action Steps:**
- Replace `Purely read-only.` in the Discovery bullet with a sentence stating discovery is read-only for product code and may create disposable evidence in `.specforce/scratch/`.
- Leave the "DIRECT EDIT PROHIBITION" paragraph unchanged in meaning for product code.
- Confirm the merge routine replaces the old managed block rather than appending a second one.

**Acceptance Check:**
`go test ./src/internal/project/` exits 0.

- [x] T2.3: [CODE] Align the repository agent guide
**Target:** `AGENTS.md`
**Context:** [US-6]

**Action Steps:**
- Replace `Purely read-only.` in the Discovery bullet with the same wording used by the generator in T2.2.
- Keep every other line, including the managed-content note, untouched.
- Verify that no occurrence of the old phrase remains.

**Acceptance Check:**
`! grep -q "Purely read-only" AGENTS.md` exits 0 and `grep -q ".specforce/scratch/" AGENTS.md` exits 0.

### Phase 3: Full Verification

- [x] T3.1: [CODE] Run the complete quality gate
**Target:** `make test`
**Context:** [US-1], [US-2], [US-3], [US-4], [US-5], [US-6]

**Action Steps:**
- Run the full test suite from the repository root.
- Fix any failing assertion caused by the rewritten blueprint or guide text, without weakening existing assertions.
- Run the linter target afterwards and resolve findings that this change introduced.

**Acceptance Check:**
`make test` exits 0.
