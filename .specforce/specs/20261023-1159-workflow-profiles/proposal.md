# Proposal: Workflow Profiles (Minimal and Expanded)

**Suggested type/size:** feature / large
**Origin:** "Pacote C2" of the consolidated plan after archived spec `.specforce/archive/20261020-2222-independent-qa-and-free-batch-log` (commit 9c6eab8 on `1.x`). Inspired by OpenSpec (github.com/Fission-AI/OpenSpec).

## Context: what OpenSpec does
OpenSpec's default profile ships only `explore` and `propose`; an expanded profile adds new/continue/ff/verify/bulk-archive/onboard. Chosen with `openspec config profile` and applied with `openspec update`. (From the README we read; confirm details in the repo.)

## Today in Specforce
- `specforce init [agents...]` installs all six embedded skills for every selected agent: `spf-spec`, `spf-constitution`, `spf-discovery`, `spf-implement`, `spf-archive`, `consultative-grill` (`[BR-KIT-01]`, kit mapping in `src/internal/agent/kit/kit.yaml`).
- Skills are discovered via metadata (architecture principle 6, "Metadata-Driven Discovery"); no hardcoded lists. Mapping per agent is centralized in `kit.yaml`.
- There is no way to install a subset.

## Problem
New users face the full lifecycle at once. The independent QA (`spf-implement` step 4) and archive are heavy; some teams only want discovery + spec, or spec + implement without archival.

## Ideas to evaluate (not decided)
- Profiles defined as metadata (e.g. a `profile` list per skill manifest, or a profiles section in `kit.yaml`): `minimal` = discovery + spec (+ grill), `standard` = + implement, `full` = + constitution + archive.
- `specforce init --profile <name>` and a way to switch (`specforce update`/re-init) that removes skills outside the profile only after user confirmation (reuse the legacy-asset prompt pattern, `[BR-KIT-03]`).
- Skills referencing others (spec hand-off to `/spf:implement`, implement hand-off to `/spf:archive`) must degrade gracefully when the next skill is not installed.
- Persist the chosen profile in `.specforce/config.yaml`.

## Constraints
- Deterministic generation across agents; atomic, no partial failure (module agent-kit section 5).
- Must not break `BR-KIT-01` (embedded kit contains exclusively sovereign skills); profiles only select among them.
- Cobra CLI and TUI wiring stay in the CLI layer; logic in a service (architecture principle 2).

## Open questions
1. Is `spf-constitution` part of every profile (other skills read `.specforce/docs/`)?
2. Behavior of hand-off lines when the target skill is absent.
3. Is the default profile the minimal one (like OpenSpec) or the full one (backward compatible)?
4. TUI selection vs flag only?

## Acceptance hints
Go tests: init with each profile installs exactly the expected skills for every supported agent; switching removes extras only after confirmation; unit tests for config persistence. `go test ./...` passes.

## How to use this proposal in a new conversation
`/spf:discovery` to settle the open questions, then `/spf:spec` (large: requirements + design + tasks), then `/spf:implement` to exercise the full flow with the new qa-tester and correction cycle (this one touches Go code, so Tier 1 and Tier 2 QA will be real).
