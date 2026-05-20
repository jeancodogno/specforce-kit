---
slug: 20260520-1116-fix-distributed-memorial-path
lens: Bugfix
---

# Bugfix: Fix Distributed Memorial Path and Generation

## 1. Issue Description
The current implementation of the Constitution Registry hardcodes the storage path for all artifacts as `.specforce/docs/`. This contradicts the Distributed Memorial architecture, where the memorial entry point is `.specforce/memorial/ROUTING.md`. Additionally, the initial generation of `ROUTING.md` uses a hardcoded Go string instead of the rich YAML-based template defined in the artifacts registry.

### Observed Logs/Errors
- `specforce constitution status` shows `memorial.md` as "Missing" even if `.specforce/memorial/ROUTING.md` exists.
- `specforce init` generates a minimal `ROUTING.md` that lacks the complete "Rules of Engagement" and "Memorial Structure" sections defined in `memorial.yaml`.

## 2. Reproduction Steps
1. Initialize a new project using `specforce init`.
2. Observe the content of `.specforce/memorial/ROUTING.md`.
   - **Current:** Basic Markdown content from `src/internal/project/memorial.go` (missing "Memorial Structure" section).
   - **Expected:** Content matching the template in `src/internal/agent/artifacts/constitution/memorial.yaml`.
3. Run `specforce constitution status`.
   - **Current:** Reports `memorial.md` missing in `.specforce/docs/`.
   - **Expected:** Reports `.specforce/memorial/ROUTING.md` as the memorial artifact and correctly identifies its existence.

## 3. Technical Requirements
- **Registry Path Overloading:** Modify `src/internal/constitution/registry.go` (specifically `loadArtifact`) to handle the `memorial` slug as a special case, mapping its path to `.specforce/memorial/ROUTING.md` instead of the default `.specforce/docs/memorial.md`.
- **Template Injection:** Update `src/internal/project/memorial.go` to accept an `fs.FS` (for artifacts) and use the `memorial.yaml` template during `Initialize`.
- **Service Decoupling:** Ensure the `memorialService` can access the template without creating circular dependencies between `project` and `constitution` packages. It might be better to load the template in the `Service` layer and pass it down.

## 4. Acceptance Criteria (BDD)

### Scenario 1: Correct Status Mapping
**GIVEN** a project with a valid `.specforce/memorial/ROUTING.md` file
**WHEN** I run `specforce constitution status`
**THEN** the "memorial" artifact SHOULD be reported as **Present**
**AND** the path displayed MUST be `.specforce/memorial/ROUTING.md`
**AND** the status check MUST NOT fail if `.specforce/docs/memorial.md` is absent.

### Scenario 2: Rich Template Generation
**GIVEN** a new project being initialized
**WHEN** the system creates the memorial directory
**THEN** the `ROUTING.md` file MUST be populated using the template from `src/internal/agent/artifacts/constitution/memorial.yaml`
**AND** it MUST include the verbatim "FOR AI AGENTS: RULES OF ENGAGEMENT" section.

### Scenario 3: Legacy Path Handling
**GIVEN** a project with a legacy `.specforce/docs/memorial.md`
**WHEN** `specforce init` is executed
**THEN** the legacy file MUST be deleted (existing logic)
**AND** the Registry MUST NOT attempt to track the legacy path anymore.

### Scenario 4: No Regressions for Standard Artifacts
**GIVEN** standard artifacts like `architecture.md` or `principles.md`
**WHEN** I run `specforce constitution status`
**THEN** they MUST still be expected and tracked in `.specforce/docs/`.
