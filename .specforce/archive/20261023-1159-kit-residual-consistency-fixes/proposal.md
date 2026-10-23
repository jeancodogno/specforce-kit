# Proposal: Kit Residual Consistency Fixes

**Suggested type/size:** bug / small (initialized as such; the spec will only need `tasks.md`)
**Origin:** follow-up 2 of archived spec `.specforce/archive/20261020-2222-independent-qa-and-free-batch-log` (commit 9c6eab8 on `1.x`).

## Problem
Small inconsistencies were found by independent QA outside the scope of the previous spec:
1. `src/internal/agent/kit/skills/spf-constitution/SKILL.yaml` (around lines 115-116): the hand-off still says "Run `/spec` to start planning a new feature." without the `/spf:` prefix, and a stray closing triple backtick sits at the end of line 115.
2. `TestLifecycleCommandNaming` (in `src/internal/agent/kit_manifests_test.go`) only forbids `` `/archive` `` and `` `/implement` ``; it never checks `` `/spec` ``, and it does not cover `spf-discovery` or `spf-constitution`. The archived design.md had planned discovery coverage.
3. `.specforce/docs/modules/agent-kit.md` documents `specforce init [agents...] [--non-interactive]`, but the binary rejects the flag (`unknown flag: --non-interactive`). Verified by running the built binary. `specforce init --help` lists only `-h` and the global `--dev`.
4. `spf-implement/SKILL.yaml` line ~44 sizing text says "plus 1 fresh `qa-tester`" and does not count the up to 2 extra correction implementers and qa-testers; line ~39 still says to pre-register a worker type on harnesses that require it, slightly in tension with line ~54 "no per-tool agent registration required".
5. Cosmetic: `spf-constitution` and `consultative-grill` still use the word "persona" while the implement blueprint standardized on "role" (only decide if the vocabulary should be unified).

## Decisions already taken
- Hand-off commands use the `/spf:` prefix everywhere (US-8 of the previous spec).
- Do not grow `spf-implement/SKILL.yaml` beyond 221 lines by `wc -l`; fold wording into existing lines.

## Open questions
1. For item 3: change the module doc, or add the flag to the CLI? (Doc-only is the small path; the flag is a CLI feature and would need its own spec.)
2. For item 5: unify vocabulary now or leave it?

## Acceptance hints
`grep -rn '`/spec`\|`/archive`\|`/implement`' src/internal/agent/kit` is empty; the naming test is extended to all skills and passes; `go test ./...` passes.

## How to use this proposal in a new conversation
Run `/spf:spec`, point to this proposal, expect a `small` spec with only `tasks.md`. Then `/spf:implement` to exercise the new free-form narration and the independent qa-tester on a tiny change.
