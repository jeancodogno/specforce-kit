---
slug: 20261023-1159-kit-residual-consistency-fixes
lens: Balanced full-stack
---

# Implementation Roadmap: Kit Residual Consistency Fixes

## 1. Execution Strategy
- **Gravity Order:** Strengthen the naming guard test first (RED), then fix the kit manifests it flags (GREEN), then the independent doc and sizing wording fixes. Tasks that touch the same file are sequential.
- **Constraints:** Hand-off commands use the `/spf:` prefix everywhere. `src/internal/agent/kit/skills/spf-implement/SKILL.yaml` must not exceed 221 lines by `wc -l` (currently 213); fold wording into existing lines. Vocabulary decision: "persona" is unified to "role" in `spf-constitution` and `consultative-grill`. The `--non-interactive` flag is NOT added to the CLI; the module doc is corrected instead.

## 2. Tasks

### Phase 1: Naming guard and constitution hand-off

- [x] T1.1: [RED] Extend TestLifecycleCommandNaming to scan the whole kit
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-1]

**Action Steps:**
- Keep the existing four `required` cases (spf-implement, spf-spec, spf-archive, instructions/archive.md) unchanged.
- Replace the forbidden list with `` `/spec` ``, `` `/archive` ``, `` `/implement` ``, `` `/discovery` ``, `` `/constitution` ``.
- Add a `filepath.Glob` over `kit/skills/*/SKILL.yaml` and `kit/instructions/*.md`, failing the test if the glob returns zero files, and apply the forbidden-token check to every matched file (so `spf-discovery` and `spf-constitution` are covered).
- Report each violation with file path and token via `t.Errorf`.

**Acceptance Check:**
Run `go test ./src/internal/agent/ -run TestLifecycleCommandNaming`; it MUST fail, reporting `` `/spec` `` in `kit/skills/spf-constitution/SKILL.yaml` (RED state).

- [x] T1.2: [GREEN] Fix the constitution hand-off command and stray fence
**Target:** `src/internal/agent/kit/skills/spf-constitution/SKILL.yaml`
**Context:** [US-1]

**Action Steps:**
- At the final hand-off line, replace `> Run `/spec` to start planning a new feature.` followed by the stray closing triple backtick with `> Run `/spf:spec` to start planning a new feature.`, moving the closing triple backtick onto its own line so the fenced block closes cleanly.
- Verify the surrounding indentation matches the neighbouring block lines.

**Acceptance Check:**
Run `go test ./src/internal/agent/ -run TestLifecycleCommandNaming`; it MUST pass. Run ``grep -rn '`/spec`\|`/archive`\|`/implement`' src/internal/agent/kit``; output MUST be empty.

- [x] T1.3: [GREEN] Unify "persona" to "role" in spf-constitution
**Target:** `src/internal/agent/kit/skills/spf-constitution/SKILL.yaml`
**Context:** [US-4]

**Action Steps:**
- Rename the heading `**B. Dynamic Skill & Persona Adoption:**` to `**B. Dynamic Skill & Role Adoption:**`.
- Check the rest of the file with `grep -in persona` and replace any remaining occurrence with "role".

**Acceptance Check:**
Run `grep -rniw 'persona\|personas' src/internal/agent/kit/skills/spf-constitution`; output MUST be empty. Run `go test ./src/internal/agent/...`; MUST pass.

### Phase 2: Independent wording and doc fixes

- [x] T2.1: [CODE] Unify "persona" to "role" in consultative-grill
**Target:** `src/internal/agent/kit/skills/consultative-grill/SKILL.yaml`
**Context:** [US-4]
**Parallel With:** T2.2, T2.3

**Action Steps:**
- In the "Grill me" intent line, replace "adopt a senior adversarial technical persona" with "adopt a senior adversarial technical role".
- Run `grep -in persona` on the file and replace any other occurrence.

**Acceptance Check:**
Run `grep -rniw 'persona\|personas' src/internal/agent/kit`; output MUST be empty. Run `go test ./src/internal/agent/...`; MUST pass.

- [x] T2.2: [CODE] Correct the init command in the agent-kit module doc
**Target:** `.specforce/docs/modules/agent-kit.md`
**Context:** [US-2]
**Parallel With:** T2.1, T2.3

**Action Steps:**
- At the `specforce init` bullet (around line 135), change `specforce init [agents...] [--non-interactive]` to `specforce init [agents...]`.
- Keep the rest of the description sentence intact.
- Do not modify any Go CLI file.

**Acceptance Check:**
Run `grep -rn 'non-interactive' .specforce/docs`; output MUST contain no match for `specforce init`. Run `specforce init --help`; the listed flags MUST match the documentation (no `--non-interactive`).

- [x] T2.3: [CODE] Align implement sizing and registration wording
**Target:** `src/internal/agent/kit/skills/spf-implement/SKILL.yaml`
**Context:** [US-3]
**Parallel With:** T2.1, T2.2

**Action Steps:**
- In the Universal Harness Compatibility line (around line 39), delete the sentence "If the harness requires pre-registration (e.g., `define_subagent`), register an explicit worker type with write tools enabled before invoking." so it no longer conflicts with "no per-tool agent registration required" (around line 54).
- In the Global Implementation Sizing line (around line 44), change "plus 1 fresh `qa-tester` in Step 4" to "plus 1 fresh `qa-tester` in Step 4 (and up to 2 extra correction implementers and qa-testers after `FAILED` verdicts)", editing in place without adding lines.

**Acceptance Check:**
Run `grep -c 'pre-registration' src/internal/agent/kit/skills/spf-implement/SKILL.yaml`; MUST print 0. Run `wc -l src/internal/agent/kit/skills/spf-implement/SKILL.yaml`; MUST be at most 221. Run `go test ./src/internal/agent/...`; MUST pass.

### Phase 3: Final verification

- [x] T3.1: [VERIFY] Run the full suite and repository greps
**Target:** `go test ./...`
**Context:** [US-1], [US-2], [US-3], [US-4]

**Action Steps:**
- Run `go test ./...` from the repository root.
- Run ``grep -rn '`/spec`\|`/archive`\|`/implement`' src/internal/agent/kit`` and confirm empty output.
- Run `grep -rniw 'persona\|personas' src/internal/agent/kit` and confirm empty output.

**Acceptance Check:**
`go test ./...` exits 0 and both greps print nothing.

## 3. Pre-emptive Mitigations
- **Risk:** The glob scan may flag legitimate occurrences of `/discovery` or `/constitution` inside other kit files (e.g., prose about directories). -> **Mitigation:** Forbidden tokens are matched with backticks (`` `/discovery` ``) so only inline command references are caught; if a legitimate hit appears, fix the text to use the `/spf:` form rather than loosening the test.
