# Proposal: Delta Specs and Living Capability Specs

**Suggested type/size:** feature / complex (cross-cutting, CLI + templates + archive). Start with `/spf:discovery`.
**Origin:** "Pacote C1" of the consolidated plan after archived spec `.specforce/archive/20261020-2222-independent-qa-and-free-batch-log` (commit 9c6eab8 on `1.x`). Inspired by OpenSpec (github.com/Fission-AI/OpenSpec).

## Context: what OpenSpec does
OpenSpec keeps a living spec per capability and each change describes only a **delta** (ADDED / MODIFIED / REMOVED / RENAMED requirements, with WHEN/THEN scenarios). On archive the deltas are merged into the main specs and the change folder moves to a dated archive. (Sourced from the OpenSpec README; the details of MODIFIED/REMOVED handling were not fully documented on the page we read, so verify against the repo before copying semantics.)

## Today in Specforce
- Archive (see `src/internal/agent/kit/instructions/archive.md`) asks the agent to hand-merge the feature into `.specforce/docs/modules/<domain>.md` (5 sections, `[BR-xx]`, `[US-xx]` with BDD) with "non-destructive merge". It is an LLM judgment call, not a structured operation; risk of silently dropping or duplicating rules.
- Modules are behavioral living specs (BR-KIT-11), not tied to requirement IDs of individual specs.
- `requirements.md` of a spec is a full document, not a delta against the module.

## Problem
Merge fidelity between a spec's requirements and the module living spec relies on prose instructions. There is no explicit "this spec ADDS BR-KIT-15, MODIFIES BR-KIT-09, REMOVES X" artifact to review before archive or to verify afterwards.

## Ideas to evaluate (not decided)
- A delta section in `requirements.md` (or a separate `delta.md`) listing ADDED/MODIFIED/REMOVED rules/use cases against a named module.
- `specforce spec archive` validating that the module doc reflects the delta (every ADDED id exists, every REMOVED id is gone, MODIFIED text changed).
- Keep the current free-form path as fallback for specs without a delta.
- Impact on the qa-tester: the delta becomes a ready-made checklist for the traceability matrix.

## Constraints from the Constitution
- Module docs stay behavioral (no source paths), 5 sections (BR-KIT-11).
- Artifact registry has deterministic topological ordering (architecture principle 7); a new artifact must slot in.
- Machine-readable state belongs in `spec.yaml`, not in the Markdown (engineering.md).
- Zero-technical-specification policy for requirements (BR-KIT-02).

## Open questions
1. Delta as a section of requirements or a new artifact (affects size matrix for `small`/`medium`/`large`)?
2. ID scheme: how to reference module rules unambiguously (`[BR-KIT-09]` is global per module; what about multiple modules in one spec)?
3. Who writes the delta: spec skill during planning, or archive skill at the end (as-built)?
4. Migration of existing modules (they follow the current format and have no delta history).

## Acceptance hints
Go unit tests for the delta parser/validator in the spec service; archive refuses to proceed when the module does not match the delta; blueprint tests per `behavioral-kit-manifest-tests`.

## How to use this proposal in a new conversation
Begin with `/spf:discovery` and ask for opinions on the open questions; expect a recommendation and trade-offs, then `/spf:spec` (size `complex` means the full grill across 5 dimensions).
