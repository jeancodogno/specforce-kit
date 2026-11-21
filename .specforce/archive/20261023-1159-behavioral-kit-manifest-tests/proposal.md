# Proposal: Behavioral Kit Manifest Tests

**Suggested type/size:** feature / medium (already initialized as such; resize if needed)
**Origin:** follow-up 1 of archived spec `.specforce/archive/20261020-2222-independent-qa-and-free-batch-log` (commit 9c6eab8 on branch `1.x`).

## Problem
The manifest tests in `src/internal/agent/kit_manifests_test.go` only assert that tokens are present in the skill blueprints. An independent QA ran 6 behavior-inverting mutations on a scratch copy of the `spf-implement` blueprint and all passed the tests:
- narrate BEFORE dispatching the worker instead of after;
- unlimited correction cycles instead of at most 2;
- hand the worker reports to the QA tester (breaking its independence);
- delete the rule "workers never ask the user directly" (the word BLOCKED still appears elsewhere);
- re-add a `[BATCH x/y COMPLETED]` bracketed chain;
- archive stops deriving next steps from findings.
Only wording that disappears is caught; semantic inversion is not.

## Findings (verified)
- Current assertions: `TestImplementBatchSubagentIsolationAndProgress`, `TestActiveQAAndBlockerProtocols`, `TestImplementRoleChartersAndResume`, `TestLifecycleCommandNaming`, all token/line-count based.
- The blueprint is a prompt (`src/internal/agent/kit/skills/spf-implement/SKILL.yaml`, 213 lines, hard ceiling 221 by `wc -l`), so tests can only inspect text, not run an agent.
- Existing negative checks cover only `Iniciando Lote` and `Delegando para perfil de`.

## Decisions already taken
- Keep tests as plain Go `testing` (engineering.md), no new dependencies.
- The line ceiling of the implement blueprint stays; do not grow the skill to satisfy tests.

## Ideas to evaluate (not decided)
- Context-anchored assertions: assert a phrase appears inside a specific section (e.g. "NEVER the worker reports" within the Step 4 qa-tester block), by slicing the content between headings.
- Ordering assertions: "Dispatch the worker FIRST" must precede "then write 1 to 3" in the same sentence/block.
- Negative assertions for forbidden templates: `[BATCH`, `COMPLETED]`, `🚀`, `MANDATORY LOG`, `Recommended Model Tier`.
- Numeric invariants: "Maximum 2 cycles" and "2 iterations" present, and no occurrence of "unlimited"/"indefinitely".
- A tiny mutation harness test: load the blueprint, apply a string mutation in memory, re-run the assertion helpers and require that they FAIL (tests the tests).

## Open questions
1. Slice by markdown heading or by fixed markers? (Headings are stable today but could be renamed.)
2. Should the mutation harness live in the same file or a dedicated `_test.go`?
3. Apply the same style to `spf-spec`, `spf-archive` and `archive.md` in this spec, or only to `spf-implement`?

## Acceptance hints
`go test ./src/internal/agent/...` passes; each of the 6 mutations above makes at least one test fail; no change to the skill texts is required.

## How to use this proposal in a new conversation
Start with `/spf:spec` and say: "Formalize `.specforce/specs/20261023-1159-behavioral-kit-manifest-tests/proposal.md`". Observe whether the new `spf-implement` behavior shows up later: free-form narration after dispatch, role charters, the independent qa-tester matrix, the BLOCKED protocol.
