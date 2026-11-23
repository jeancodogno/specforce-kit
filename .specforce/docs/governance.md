# Governance

> **Version:** 1.1.0 (MINOR: Release Process added) · **Last amended:** 2026-11-23

## Versioning Policy (SemVer)
- **MAJOR:** Breaking constitutional or architectural process change. Requires team consensus.
- **MINOR:** New principle or standard added without breaking prior commitments.
- **PATCH:** Clarifications, typo fixes, wording, and non-functional adjustments.

## Amendment Process
1. Propose the change with a clear rationale and impact analysis.
2. Classify the SemVer impact (Major/Minor/Patch).
3. Update all affected docs in `.specforce/docs/`.
4. Record the version bump and the amendment date in the document's header.

## Release Process
1. Classify the version with SemVer: template or skill behavior changes that keep existing specs valid are MINOR; fixes only are PATCH.
2. Update user documentation in `docs/{en,pt,es}/` for any behavior the release changes, scoped with "from X.Y.Z onward".
3. Add a dated `## [X.Y.Z] - YYYY-MM-DD` entry to `CHANGELOG.md` (Keep a Changelog sections) and its link reference at the bottom.
4. Bump the version in `package.json` (package and every `optionalDependencies` entry), `src/internal/core/constants.go` and the recovery command in `AGENTS.md`, in one `chore(release): prepare vX.Y.Z` commit.
5. Commits carry no `Co-Authored-By` trailer.
6. Releases are prepared on the `X.x` release branch (for example `1.x`). Run `go test ./...` and `make lint`, then push the branch and an annotated `vX.Y.Z` tag; the tag triggers the release workflow.
7. `main` holds the current released version: open a "Release X.Y.Z" pull request from `X.x` into `main` and merge it with a merge commit (not squash or rebase) so the release tag stays on `main`'s history. The merge is done by a human.

## Ownership & Review
- **Owner:** Project Lead / Solutions Architect.
- **Review Cadence:** Quarterly or after every Major release.

## Agentic Authority & Boundaries
- **Permitted AI Autonomous Actions:**
  - Propose new Spec drafts for requested features.
  - Generate implementation code based on approved Specs.
  - Consolidate lessons learned, architectural precedents and as-built business rules into the affected module living specs (`.specforce/docs/modules/<slug>.md`) during archival.
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
- **Translation Parity:** A change to an English page MUST update its Portuguese and Spanish versions in the same change.
- **Human vs. Agent Actions:** In user documentation, every step the reader performs MUST be labeled as either a terminal command or an agent chat request. Commands the agent runs on its own appear only as explanation (prose or a collapsed "behind the scenes" block), never as a step.
- **Chat First:** User documentation MUST NOT ask the reader to run a terminal check for state the agent already reports (task status, delta state, QA results) or to choose a value the agent decides (such as the spec size).
- **Concrete Examples:** Example chat requests MUST be concrete (a named feature, bug or spec slug), not bare placeholders.
- **Published-Version Accuracy:** Update, migration and release documentation MUST describe the behavior of the versions users actually have installed, verified against the released git tags, and scope any behavior that exists only in unreleased code with the version it ships in (for example "from 1.0.0 onward").
