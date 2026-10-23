# Delta: spec-management
## ADDED
- `[NEW:delta-required]` Every `medium`, `large` or `complex` specification declares its module impact before implementation: either one delta document per affected module (ADDED, MODIFIED and REMOVED sections carrying complete blocks, with provisional `NEW:` keys for additions) or the explicit declaration `delta.impact: none`. `small` specifications are exempt and keep the free-form archive path.
- `[NEW:delta-seal]` Sealing a delta (`specforce spec delta seal <slug>`) records a fingerprint of the current text of every module block that the delta modifies or removes, so that parallel specifications changing the same block are detected: a block matching the recorded fingerprint is pending, a block matching the delta text is applied, and any other text is a conflict that requires redoing the delta against the current module. A modified use case may not lose scenarios, and unknown delta sections are reported as problems.
- `[NEW:delta-gate]` `specforce spec archive <slug>` refuses to archive a `medium`, `large` or `complex` specification unless its delta status is `applied` or `none`; `small` specifications report `exempt`. The override flag for pending tasks never bypasses this check, and a refusal leaves the filesystem untouched.
- `[NEW:delta-ids]` Final identifiers for added blocks are allocated at archive time as the next free number per kind (rules and use cases numbered independently) and recorded in the specification metadata as a key-to-identifier mapping; identifiers of untouched or modified blocks are never renumbered and removed identifiers are never reused.
### [NEW:us-delta-declare] Delta Declaration and Sealing During Planning
- **Scenario:** Planning a specification that changes module behavior
  - **GIVEN** a `medium` or larger specification that adds, modifies or removes module rules
  - **WHEN** planning writes the delta and runs `specforce spec delta seal <slug>`
  - **THEN** metadata records the base fingerprints and `specforce spec delta check <slug>` reports every item as pending.
- **Scenario:** Planning a specification with no module impact
  - **GIVEN** a `medium` specification that changes no module behavior
  - **WHEN** it declares `delta.impact: none`
  - **THEN** `specforce spec delta check <slug>` reports `none` and sealing is not required.
- **Scenario:** Another specification changed the same block after sealing
  - **GIVEN** a sealed delta whose target block was changed by a different specification
  - **WHEN** `specforce spec delta check <slug>` runs
  - **THEN** the item is reported as a conflict with the instruction to redo the delta and seal again.
### [NEW:us-delta-gate] Archive Gate on Delta State
- **Scenario:** Archive refused while the module does not reflect the delta
  - **GIVEN** a `medium` or larger specification whose delta status is pending, missing or invalid
  - **WHEN** archiving is requested, with or without the override for pending tasks
  - **THEN** archiving is refused, every problem and pending item is listed with the next command to run, and the specification stays active.
- **Scenario:** Archive allowed when the module reflects the delta
  - **GIVEN** a specification whose delta status is applied, none or exempt
  - **WHEN** archiving is requested
  - **THEN** the specification moves to the archive.
