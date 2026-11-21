---
slug: 20261023-1159-delta-living-specs
lens: Backend-heavy
---

# Feature: Delta Specs and Structured Module Living Specs

## 1. Context & Value
Module living specs already hold the behavior of each capability, but they are updated at archive time by free-form agent judgment, so rules can be silently dropped or duplicated and nobody can review beforehand what a specification will change. This feature gives modules a stable block structure (Part A) and introduces a per-specification delta that declares exactly which module rules and use cases are added, modified or removed, which is reviewed during planning, reconciled and applied at archive, and verified before the specification can be archived (Part B).

## 2. Out of Scope (Anti-Goals)
- Do not make the CLI apply or rewrite module content; the agent applies the delta and the CLI only inspects and verifies.
- Do not support a rename operation; changing a rule's wording is a modification and a rule keeps its identifier.
- Do not give identifiers or deltas to the module sections for domain scope, public integration surfaces and operational invariants; they remain as-built syntheses written at archive.
- Do not add a bulk migration command for modules; legacy modules are migrated lazily when they first receive a delta.
- Do not change any Constitution document (principles, architecture, UI/UX, security, engineering, governance).
- Do not require deltas for `small` specifications.
- Do not change the artifact registry, the size matrix of required artifacts, workflow profiles, or customizable workflow schemas.
- Do not alter the content owned by the sibling spec for executable criteria and archive preflight beyond the delta-related passages.

## 3. Acceptance Criteria (BDD)

### [US-1] Module Identity and Block Structure
**User Story:** AS A Lead Solutions Architect, I WANT every module living spec to declare a unique identifier prefix and to express each business rule and each use case as a uniformly structured block with a stable identifier, SO THAT any rule can be referenced, reviewed and changed unambiguously across specifications.

**Scenarios:**
1. **[Happy Path]** GIVEN a module living spec that declares the prefix `SPEC` WHEN its rules and use cases are read THEN every rule is identified as `BR-SPEC-<number>` and every use case as `US-SPEC-<number>`, each forming one self-contained block (a rule with its sub-points, or a use case with all its scenarios).
2. **[Edge Case]** GIVEN a module whose rule or use case carries an identifier with a different prefix than the one the module declares WHEN the module is inspected THEN the inconsistency is reported as a problem naming the offending identifier.
3. **[Edge Case]** GIVEN two modules that declare the same prefix WHEN the modules are inspected THEN the duplicate prefix is reported as a problem naming both modules.
4. **[Edge Case]** GIVEN a module without a declared prefix that has no delta pointed at it WHEN any spec command runs THEN the module is ignored and no problem is reported.

**UI/UX Specifics:**
- **View/Component:** Plain terminal messages through the project's standard output channel.
- **Feedback Logic:** Each problem is one line naming the module and the identifier involved.
- **Keybindings:** Not applicable to this requirement.

**Technical Constraints (NFR):**
- **[Performance]:** Inspecting all modules of a project completes in under 50ms.
- **[Integrity]:** Identifiers are unique within a module; fenced example content inside a module is never treated as a real rule.

### [US-2] Declare the Delta During Planning
**User Story:** AS AN AI-First Senior Developer, I WANT the planning step of every `medium`, `large` or `complex` specification to declare the module impact, SO THAT reviewers see before implementation which rules and use cases the specification will change.

**Scenarios:**
1. **[Happy Path]** GIVEN a `large` specification whose requirements affect the behavior of one module WHEN planning completes THEN the specification contains a delta document for that module, written before the task roadmap, and the roadmap tasks cover every delta item.
2. **[Happy Path]** GIVEN a `medium` specification that is pure developer tooling with no module behavior change WHEN planning completes THEN the specification explicitly declares that it has no module impact instead of a delta.
3. **[Happy Path]** GIVEN a `small` specification WHEN planning completes THEN no delta and no declaration is required.
4. **[Edge Case]** GIVEN a `medium` or larger specification with neither a delta nor a no-impact declaration WHEN its delta is inspected THEN it is reported as missing and the specification cannot be archived.
5. **[Edge Case]** GIVEN a specification that touches two modules WHEN planning completes THEN there is one delta document per module and each is inspected independently.

**UI/UX Specifics:**
- **View/Component:** A short list in the planning hand-off summary naming each module and the count of added, modified and removed items.
- **Feedback Logic:** The summary states "no module impact" explicitly when declared.
- **Keybindings:** Not applicable to this requirement.

**Technical Constraints (NFR):**
- **[Performance]:** Reading the delta of a specification completes in under 50ms.
- **[Integrity]:** The declared module impact is stored in the specification's machine-readable metadata, never only in prose.

### [US-3] Delta Operations and Content Rules
**User Story:** AS A Lead Solutions Architect, I WANT a delta to express additions, modifications and removals of module rules and use cases with complete content, SO THAT a reviewer can read exactly what the module will look like for each touched item.

**Scenarios:**
1. **[Happy Path]** GIVEN a delta that adds a new rule WHEN it is read THEN the new rule is identified by a provisional key (not a final identifier) and carries its full text.
2. **[Happy Path]** GIVEN a delta that modifies an existing rule WHEN it is read THEN it carries the identifier of the existing rule and the complete replacement text of that block, including all its scenarios when it is a use case.
3. **[Happy Path]** GIVEN a delta that removes an existing rule WHEN it is read THEN it lists only the identifier to remove.
4. **[Edge Case]** GIVEN a delta that modifies a use case and omits a scenario that exists in the current module block WHEN the delta is inspected THEN it is reported as a problem because scenarios cannot be dropped silently.
5. **[Edge Case]** GIVEN a delta where the same identifier appears under two operations WHEN the delta is inspected THEN it is reported as a conflicting entry.
6. **[Edge Case]** GIVEN a delta that targets a module that does not exist or an identifier that does not exist in the module WHEN the delta is inspected THEN each case is reported as a problem naming the module or identifier.

**UI/UX Specifics:**
- **View/Component:** Delta documents are readable Markdown with one section per operation.
- **Feedback Logic:** Problems are listed all at once, never one at a time.
- **Keybindings:** Not applicable to this requirement.

**Technical Constraints (NFR):**
- **[Performance]:** Inspecting a delta against its module completes in under 50ms.
- **[Safety & Security]:** The module name in a delta is accepted only as a lowercase hyphenated name and can never address anything outside the project's module collection.

### [US-4] Stale-Base Protection
**User Story:** AS AN AI-First Senior Developer, I WANT the delta to remember the exact text of every module rule it modifies or removes at the time it was planned, SO THAT two specifications changing the same rule in parallel are detected instead of the last one silently overwriting the other.

**Scenarios:**
1. **[Happy Path]** GIVEN a delta that modifies or removes existing rules WHEN the developer seals the delta THEN the current text fingerprint of each referenced rule is recorded in the specification's metadata.
2. **[Happy Path]** GIVEN a sealed delta and a module whose referenced rule still matches the recorded text WHEN the delta is inspected THEN that item is reported as pending application.
3. **[Edge Case]** GIVEN a sealed delta and a module whose referenced rule was changed by another specification since sealing, so that it matches neither the recorded text nor the delta text WHEN the delta is inspected THEN the item is reported as a base conflict and the developer is told to redo the delta against the current module.
4. **[Edge Case]** GIVEN a modify item whose replacement text is identical to the current module text WHEN the delta is inspected THEN it is reported as having no effect.
5. **[Edge Case]** GIVEN a delta that was never sealed and references existing rules WHEN the delta is inspected THEN it is reported as unsealed.

**UI/UX Specifics:**
- **View/Component:** One line per item with its state (pending, applied, conflict, unsealed).
- **Feedback Logic:** Conflicts carry the recommended next step in the same message.
- **Keybindings:** Not applicable to this requirement.

**Technical Constraints (NFR):**
- **[Performance]:** Sealing completes in under 50ms.
- **[Integrity]:** Sealing never modifies the module or the delta document, only the specification's metadata; text comparison ignores line-ending style and trailing whitespace.

### [US-5] Delta Inspection Before Archive
**User Story:** AS AN AI-First Senior Developer, I WANT a read-only inspection command for a specification's delta, SO THAT I can rehearse the archive and see everything that is wrong or still pending before attempting it.

**Scenarios:**
1. **[Happy Path]** GIVEN a specification whose delta items are all pending WHEN the inspection runs THEN the overall state is pending and each item is listed.
2. **[Happy Path]** GIVEN a specification whose module already reflects every delta item WHEN the inspection runs THEN the overall state is applied.
3. **[Happy Path]** GIVEN additions in a delta WHEN the inspection runs THEN it reports, per module, the next free identifier for rules and for use cases so that final identifiers can be allocated without collision.
4. **[Edge Case]** GIVEN several unrelated problems in one delta WHEN the inspection runs THEN all problems are reported in one run and the command ends with a failure.
5. **[Edge Case]** GIVEN a `small` specification WHEN the inspection runs THEN the state is exempt and no problem is reported.

**UI/UX Specifics:**
- **View/Component:** A human-readable list by default and a machine-readable report on request.
- **Feedback Logic:** Overall states are none, exempt, missing, pending, applied and invalid; success feedback uses the project's standard success style and problems the standard error style.
- **Keybindings:** Not applicable to this requirement.

**Technical Constraints (NFR):**
- **[Performance]:** Inspection completes in under 50ms.
- **[Integrity]:** The command never creates, changes or deletes any file.

### [US-6] Archive Reconciliation and Application
**User Story:** AS A Lead Solutions Architect, I WANT archiving to reconcile the delta with what was actually built and then apply it to the module, SO THAT module living specs stay faithful to reality and keep a stable identifier history.

**Scenarios:**
1. **[Happy Path]** GIVEN a finished specification with a delta WHEN archiving starts THEN the delta is first corrected to match the as-built behavior and then applied, and the module blocks equal the delta blocks.
2. **[Happy Path]** GIVEN added items with provisional keys WHEN the delta is applied THEN each receives the next free identifier of its module (rules and use cases numbered independently, in order of appearance) and the key-to-identifier mapping is recorded in the specification's metadata.
3. **[Happy Path]** GIVEN an applied delta that changes what the module's scope, public surfaces or operational invariants describe WHEN archiving completes THEN those sections are refreshed as an as-built synthesis.
4. **[Happy Path]** GIVEN a target module that has no declared prefix WHEN archiving applies its first delta THEN a prefix is proposed to and confirmed by the user, the module is normalized to the structured format, and no existing rule or use case is lost.
5. **[Edge Case]** GIVEN a base conflict detected before application WHEN archiving runs THEN no module content is changed and the developer is told to redo the delta.
6. **[Edge Case]** GIVEN a specification declared as having no module impact WHEN archiving runs THEN no module is changed.

**UI/UX Specifics:**
- **View/Component:** The final archive summary lists each module updated with counts of added, modified and removed items.
- **Feedback Logic:** The summary links to each updated module.
- **Keybindings:** Not applicable to this requirement.

**Technical Constraints (NFR):**
- **[Performance]:** Verification steps required by archiving each complete in under 50ms.
- **[Integrity]:** Identifiers of modified or untouched rules are never renumbered; removed identifiers are never reused.

### [US-7] Archive Gate
**User Story:** AS A Lead Solutions Architect, I WANT archiving to refuse a specification whose module does not reflect its delta, SO THAT a specification can never be closed with silently lost or duplicated rules.

**Scenarios:**
1. **[Happy Path]** GIVEN a specification whose delta state is applied, none or exempt WHEN archiving is requested THEN archiving proceeds as before.
2. **[Edge Case]** GIVEN a delta state of pending, missing, invalid or conflict WHEN archiving is requested THEN archiving is refused, the full list of problems is shown, and the specification stays active.
3. **[Edge Case]** GIVEN pending tasks and the existing override for archiving with pending tasks WHEN archiving is requested with that override THEN the delta verification still applies; the override never bypasses it.
4. **[Edge Case]** GIVEN a `small` specification resized to `medium` or larger WHEN archiving is requested without a delta or a no-impact declaration THEN archiving is refused as missing.
5. **[Edge Case]** GIVEN an in-flight `medium` or larger specification created before this feature with no delta information WHEN archiving is requested THEN it is refused as missing until the author adds a delta or a no-impact declaration.

**UI/UX Specifics:**
- **View/Component:** Refusal message in the project's standard error style followed by the problem list.
- **Feedback Logic:** The message names the exact next command or action to unblock archiving.
- **Keybindings:** Not applicable to this requirement.

**Technical Constraints (NFR):**
- **[Performance]:** The gate adds under 50ms to the archive command.
- **[Reliability]:** A refused archive leaves the filesystem untouched.

### [US-8] Delta-Driven Verification Checklist
**User Story:** AS A Lead Solutions Architect, I WANT the independent verification role to receive one checklist row per delta item, SO THAT every rule the specification adds, modifies or removes is proven before archive.

**Scenarios:**
1. **[Happy Path]** GIVEN a specification with a delta WHEN independent verification runs THEN the traceability matrix contains one row for every added, modified and removed item in addition to the rows for acceptance criteria, each with a verdict.
2. **[Edge Case]** GIVEN a specification with no delta or a no-impact declaration WHEN independent verification runs THEN no delta rows are added and the matrix is built from acceptance criteria only.
3. **[Edge Case]** GIVEN a removal item whose rule is still present in the module after archive-time application WHEN verification of that row runs THEN the verdict is failed.

**UI/UX Specifics:**
- **View/Component:** Matrix rows labeled by operation and identifier or provisional key.
- **Feedback Logic:** Existing verdict vocabulary is reused; no new verdicts are introduced.
- **Keybindings:** Not applicable to this requirement.

**Technical Constraints (NFR):**
- **[Performance]:** Reading the delta for the checklist completes in under 50ms.
- **[Maintainability]:** The checklist rule lives in the shared verification role charter and is covered by the kit's automated blueprint tests.

### [US-9] Migration of This Project's Modules and Backward Compatibility
**User Story:** AS A Lead Solutions Architect, I WANT this project's existing modules to adopt the new structure without losing any rule or use case, and specifications that predate the feature to remain workable, SO THAT the project can use deltas immediately.

**Scenarios:**
1. **[Happy Path]** GIVEN the project's three existing modules WHEN this feature is delivered THEN each declares a unique prefix and every existing rule and use case already follows the block convention with no content loss.
2. **[Happy Path]** GIVEN a `small` specification or a specification that declares no module impact WHEN archiving runs THEN the previous free-form archive path remains available unchanged.
3. **[Edge Case]** GIVEN a module in a user project whose rules use unprefixed identifiers WHEN it is not the target of any delta THEN nothing is reported and nothing is changed.

**UI/UX Specifics:**
- **View/Component:** Module header declaring the prefix.
- **Feedback Logic:** No new messages for untouched modules.
- **Keybindings:** Not applicable to this requirement.

**Technical Constraints (NFR):**
- **[Performance]:** Existing spec commands keep their current response times.
- **[Integrity]:** Module migration preserves every existing identifier and scenario verbatim.

## 4. Business Invariants
- A rule or use case identifier is unique within its module, belongs to exactly one module prefix, and is never reused after removal.
- Modified and removed items must reference identifiers that exist in the module at the time of sealing.
- A `medium`, `large` or `complex` specification can be archived only when its delta is applied, or it declares no module impact.
- Delta verification is read-only and never changes the module, the delta, or any other file.
- The archive override for pending tasks never bypasses delta verification.
- Requirements, deltas and modules stay behavioral: no source paths, internal structures or implementation details.
- Domain scope, public integration surfaces and operational invariants of a module are rewritten as-built at archive and are not delta items.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Every delta-related command and gate completes in under 50ms for projects with up to 50 modules and specifications.
- **[Reliability]:** Fail-fast with clear errors listing every problem; a refused operation leaves the filesystem untouched.
- **[Security]:** Module names are restricted to lowercase hyphenated names and cannot address locations outside the project's module collection; metadata files keep owner-only permissions.
- **[Maintainability]:** New logic has unit tests covering happy paths and edge cases, and the kit's blueprints are covered by automated manifest tests.
