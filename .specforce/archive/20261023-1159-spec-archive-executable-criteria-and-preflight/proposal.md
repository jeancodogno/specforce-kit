# Proposal: Executable Criteria in Spec and Preflight/Closing Improvements in Archive

**Suggested type/size:** feature / large (blueprints + templates + tests)
**Origin:** "Pacote B" of the consolidated plan made after the archived spec `.specforce/archive/20261020-2222-independent-qa-and-free-batch-log` (commit 9c6eab8 on `1.x`). That spec only standardized hand-off commands and made archive next steps findings-based.

## Problem
The new independent qa-tester proves each acceptance criterion with evidence, but the planning skill does not yet produce criteria that are easy to execute, and the archive skill relies entirely on the CLI instructions for verification.

## Findings (verified in `src/internal/agent/kit/skills/spf-spec/SKILL.yaml`, `spf-archive/SKILL.yaml`, `kit/instructions/archive.md`)
- spec rule 5 forbids spawning subagents; grounding for `large`/`complex` specs reads many files in the main session.
- The tiered grill hard gate demands 100% explicit confirmation, heavy for `medium` specs.
- `proposal.md` is only "checked" in `context_files`; decisions already taken in discovery are asked again.
- spec rule forbids `N/A` literally, conflicting with template sections that legitimately do not apply.
- Bug specs have no mandatory reproduction/root cause or regression-test-first task, although engineering.md says every bug must be reproduced with a failing test before the fix.
- archive: verification of completion is delegated to the CLI text; the final summary does not use clickable links like the spec skill; no suggested commit message.

## Proposed changes (to be confirmed in the grill)
**spf-spec**
1. Acceptance criteria must be executable by a tester: every task `Acceptance Check` is a non-interactive command; `large`/`complex` specs include black-box/E2E scenarios the qa-tester can run.
2. `bug` specs require reproduction + root cause in requirements and a first RED regression-test task.
3. `proposal.md` pre-answers the grill: decisions already taken are not asked again; only open questions are asked.
4. Grill in waves with a recommendation per question; accept "approve all recommendations" as confirmation.
5. Allow read-only exploration subagents during grounding for `large`/`complex` (planning and writing stay in the main session).
6. Replace the literal `N/A` ban with "omit the section or justify in one line".

**spf-archive / archive.md**
7. Explicit preflight: refuse and list pending tasks if progress is not 100%.
8. Final summary with clickable links to archived artifacts and changed module docs.
9. Suggest (never run) a conventional commit message.

## Decisions already taken
- Skills keep the existing line/size discipline; prefer folding text into existing sections.
- Hand-off commands use the `/spf:` prefix.
- Requirements stay free of technical detail (BR-KIT-02).

## Open questions
1. One spec for spec+archive, or split (spec first)?
2. Is allowing subagents in `spf-spec` compatible with the Constitution principle "Primary Orchestration Only" (engineering.md)? It only forbids LLM blackbox planning modes, but confirm.
3. Where do link/commit-suggestion rules live: the skill, or `archive.md` (CLI-injected)?

## Acceptance hints
Manifest tests assert the new rules in context (see proposal `behavioral-kit-manifest-tests`); `go test ./...` passes; a generated bug spec contains a RED regression task first.

## How to use this proposal in a new conversation
Run `/spf:spec` with this proposal as input and observe the grill: it should not re-ask decisions listed above. Then `/spf:implement` to exercise batch narration, role charters and the independent qa-tester.
