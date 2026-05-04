# Governance

## Versioning Policy (SemVer)
- **MAJOR:** Breaking constitutional or architectural process change. Requires team consensus.
- **MINOR:** New principle or standard added without breaking prior commitments.
- **PATCH:** Clarifications, typo fixes, wording, and non-functional adjustments.

## Amendment Process
1. Propose the change with a clear rationale and impact analysis.
2. Classify the SemVer impact (Major/Minor/Patch).
3. Update all affected docs in `.specforce/docs/`.
4. Record the version bump and the amendment date in the document's header.

## Ownership & Review
- **Owner:** Project Lead / Solutions Architect.
- **Review Cadence:** Quarterly or after every Major release.

## Agentic Authority & Boundaries
- **Permitted AI Autonomous Actions:**
  - Propose new Spec drafts for requested features.
  - Generate implementation code based on approved Specs.
  - Update the Project Memorial (`.specforce/docs/memorial.md`) with lessons learned and architectural precedents during archival or implementation.
  - Refactor existing code within the scope of a feature.
- **Prohibited AI Actions (Requires Human Approval):**
  - Modify the project Constitution (`.specforce/docs/`).
  - Approve Specs for implementation.
  - Change project-wide architectural patterns or technology stack.
  - Approve Pull Requests or perform final merges.
- **Conflict Resolution Protocol:**
  - If a user prompt or a draft Spec contradicts the Constitution, the AI MUST refuse the prompt and cite the violated Principle/Rule.

## Documentation & Internationalization
- **Primary Language:** English (Source of Truth).
- **Supported Regions:** Portuguese (PT) and Spanish (ES).
- **Structure:** Documentation must be organized in `docs/{en,pt,es}/` subdirectories.
- **Consistency:** Filenames must remain identical across language folders.
- **Navigation:** Root `README.md` and regional versions MUST contain a language selector at the top.
- **Verification:** All internal links in documentation MUST be verified and relative to the current language folder during any restructuring.
