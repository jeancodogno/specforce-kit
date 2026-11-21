# Delta: agent-kit
## ADDED
- `[NEW:delta-plan]` Planning blueprints MUST direct `medium`, `large` and `complex` specifications to write one delta document per affected module after requirements (and design when required) and before the task roadmap, to seal it, to make the roadmap cover every delta item, and to list the delta in the hand-off summary; when no module changes, the blueprint MUST direct declaring `delta.impact: none` instead and not sealing.
- `[NEW:delta-archive]` Archival blueprints and instructions MUST direct the agent to correct each delta to the as-built result, allocate final identifiers only in the module (the delta keeps its provisional key and the mapping is recorded in the specification metadata), apply the delta so module blocks equal delta blocks, rewrite the scope, public surface and operational sections as-built where the delta changes them, stop without touching the module on a base conflict, lazily migrate a prefix-less target module after the user confirms a proposed prefix, and require `specforce spec delta check <slug>` to report `applied` (or `none`/`exempt`) before archiving.
- `[NEW:delta-qa]` The independent qa-tester's traceability matrix MUST contain one additional row per delta item (operation plus identifier or provisional key) using the existing verdicts; specifications without a delta or with `delta.impact: none` add no rows, and a removal row fails when the block is still present.
### [NEW:us-delta-archive] Delta-Driven Archive Reconciliation
- **Scenario:** Archiving a specification with a delta
  - **GIVEN** a finished `medium` or larger specification with a sealed delta
  - **WHEN** `/spf:archive` runs
  - **THEN** the agent corrects the delta to the as-built result, allocates final identifiers, applies the delta to the module, refreshes the scope, public surface and operational sections where affected, and archives only after the delta check reports applied.
- **Scenario:** Base conflict detected at archive
  - **GIVEN** a delta whose target block was changed by another specification
  - **WHEN** `/spf:archive` runs
  - **THEN** the module is left untouched and the developer is told to redo the delta and seal it again.
## MODIFIED
- `[BR-KIT-11]` Archival blueprints and instructions MUST mandate that `.specforce/docs/modules/<domain>.md` living specifications are structured exclusively as 5-part behavioral specifications (Domain Scope, Invariants, BDD Scenarios, Public Integration Surfaces, Operational Invariants), omitting low-level internal code file paths and opportunistically migrating legacy module formats upon merge, with every module declaring a unique `prefix` in its frontmatter and expressing each rule and use case as a uniformly structured block identified as `BR-<PREFIX>-<NN>` or `US-<PREFIX>-<NN>`.
