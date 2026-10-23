# Delta: constitution
## ADDED
- `[NEW:module-prefix]` Every module living spec declares a unique uppercase `prefix` in its frontmatter; each business rule is a block identified as `BR-<PREFIX>-<NN>` and each use case a block headed `US-<PREFIX>-<NN>`, so any block can be referenced unambiguously by deltas. Modules without a prefix are ignored until they first receive a delta, at which point they are migrated with a user-confirmed prefix.
### [NEW:us-module-blocks] Addressable Module Blocks
- **Scenario:** Referencing a rule from a specification
  - **GIVEN** a module that declares the prefix `SPEC`
  - **WHEN** a delta modifies `BR-SPEC-03`
  - **THEN** the reference identifies exactly one block in exactly one module.
- **Scenario:** Inconsistent identifiers
  - **GIVEN** a module containing an identifier whose prefix differs from the declared one, or two modules declaring the same prefix
  - **WHEN** a delta targets the module
  - **THEN** the inconsistency is reported as a problem and the delta is invalid.
