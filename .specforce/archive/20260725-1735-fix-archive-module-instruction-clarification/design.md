# Design Specification: Fix Archive Module Instruction Clarification

## Architecture & Design Changes

### Affected Files
- `src/internal/agent/kit/instructions/archive.md`

### Proposed Changes

#### Update Step 5 in `archive.md`
Rewrite Step 5 ("Knowledge Harvesting & Module Behavior Consolidation (Mandatory Step)") to structure module decision logic cleanly into distinct steps:

1. **Identify Affected Domain:** Determine the business domain/context affected by this feature (e.g., `messaging`, `auth`, `billing`).
2. **Scan Invariants:** Scan the feature's `requirements.md` and `design.md` for new domain invariants, business logic, or technical contracts.
3. **Check Module Affinity & Action:**
   - Check the `modules` list in Constitution status.
   - **If a matching module EXISTS (`.specforce/docs/modules/<module>.md`):** Merge and update it with the new canonical behavior of that domain.
   - **If NO matching module exists (including when `modules` is empty):** Ask the user for approval to create `.specforce/docs/modules/<domain>.md` containing the new domain invariants.
4. **N/A Restriction:** Labeling module consolidation as `N/A` is ONLY permitted if the feature is purely cross-cutting infrastructure/tooling without any business domain context.
