# Proposal: Project-Level Customizable Workflow Schemas

**Suggested type/size:** feature / complex (cross-cutting; touches artifact registry, QA rules, config). Start with `/spf:discovery`.
**Origin:** "Pacote C3" of the consolidated plan after archived spec `.specforce/archive/20261020-2222-independent-qa-and-free-batch-log` (commit 9c6eab8 on `1.x`). Inspired by OpenSpec (github.com/Fission-AI/OpenSpec) project config and `schemas/spec-driven`.

## Context: what OpenSpec does
OpenSpec allows customization via profiles, project config and schemas (a `schemas/spec-driven` folder in the repo; community schemas as standalone repos). (Per the README we read; the schema format itself was not described on the page; read its Customize docs before designing.)

## Today in Specforce
- Artifact definitions (description, instruction, template, dependency) come from the embedded kit via `specforce spec artifact <name> --json`, with sizing categories small/medium/large/complex deciding which artifacts are required (topological order, architecture principle 7).
- Project overrides are limited to `instructions:` per artifact type and `hooks:` in `.specforce/config.yaml` (architecture principle 4 and 5). Embedded kits "MAY be overridden by local versions for development" (`--dev`).
- Skill texts (spf-implement QA tiers, narration, role charters, batch budget) are fixed in the embedded blueprints.

## Problem
Teams cannot adapt which artifacts exist, their templates, or the QA protocol (tiers, correction-cycle limit, batch budget, roles) without forking the kit.

## Ideas to evaluate (not decided)
- A project-level schema file in `.specforce/` describing artifacts, dependencies and sizing matrix overrides, validated by strict Go types (architecture principle 6 "Strict Typing").
- QA configuration surface: enable/disable tiers, `max_correction_cycles` (currently fixed at 2), batch budget; injected into the instructions via CLI the same way `specforce archive instructions` injects context (engineering.md "Instruction-Driven Agent Pattern").
- Precedence: embedded defaults < project schema < `--dev` local kit.

## Constraints
- Skills should not carry complex hardcoded logic; dynamic instructions come from the CLI (engineering.md).
- Constitution documents may not be modified by AI without human approval (governance.md).
- Deterministic generation across supported agents; backward compatible for projects with no schema.

## Open questions
1. How much of `spf-implement` can move to CLI-injected instructions without breaking the 221-line/token budget and tests?
2. Schema validation errors: fail fast at `specforce init`/`spec init` or at use time?
3. Is this a prerequisite for, or a follow-up of, `workflow-profiles` and `delta-living-specs`? Recommended order to decide in discovery: profiles → delta specs → schemas.

## Acceptance hints
A project with a custom schema produces different required artifacts and a different QA tier set; invalid schemas fail with a clear message; projects without a schema behave exactly as before; `go test ./...` passes with new service-level tests.

## How to use this proposal in a new conversation
Run `/spf:discovery` first and ask for a recommended order against the other two Pacote C proposals (`20261023-1159-workflow-profiles`, `20261023-1159-delta-living-specs`), then `/spf:spec` at size `complex` (full adversarial grill).
