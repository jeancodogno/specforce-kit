---
slug: 20261121-1156-release-pipeline-hardening
lens: Integration
---

# Feature: Release Pipeline Hardening for the Stable 1.0.0

## 1. Context & Value
Specforce is about to publish its first stable version (1.0.0) after a series of betas. Today the version line being released (`1.x`) was never verified by continuous integration, Windows is shipped but never tested, a pre-release can become the "latest" release that stable users of the native binary auto-update to, and a partial publish can leave users with a package that points to platform binaries that do not exist. This feature makes every release verified, correctly channelled, complete or safely resumable, and traceable to its source.

## 2. Out of Scope (Anti-Goals)
- Do not change the auto-updater's behaviour or code; the stable-release notice for pre-release users is a separate specification (`beta-stable-notice`).
- Do not write the 1.0.0 documentation, changelog, security policy or release notes text; that is the separate specification `v1-release-docs` (this feature only makes the release page use a notes file when one exists).
- Do not publish, tag or merge anything; merging `1.x` into `main`, tagging `v1.0.0-rc.1` and `v1.0.0` are manual maintainer operations.
- Do not add mutation testing or security scans to the release run; they stay in continuous integration.
- Do not run the core coverage gate on Windows or macOS.
- Do not add new distribution channels (package managers other than npm, container images) or new target platforms.

## 3. Acceptance Criteria (BDD)

### [US-1] Continuous verification of every version line
**User Story:** AS A release maintainer, I WANT every change pushed or proposed to the main branch or to a version line branch to be verified automatically on all shipped operating systems, SO THAT a release is only cut from code that already passed lint, security scans, tests and mutation quality.

**Scenarios:**
1. **[Happy Path]** GIVEN a commit pushed to `main` or to a version line branch named `<major>.x` (for example `1.x`, `2.x`) WHEN continuous integration runs THEN lint, security scans, mutation quality and the unit tests run, and the unit tests run on Linux, macOS and Windows.
2. **[Happy Path]** GIVEN a change proposed for merge into `main` or into a version line branch WHEN continuous integration runs THEN the same verification runs as for a push.
3. **[Happy Path]** GIVEN continuous integration runs THEN the release tooling's own tests run once on Linux and the core coverage gate (core packages average at least 80%) runs once, on Linux only.
4. **[Edge Case]** GIVEN the unit tests fail only on Windows WHEN continuous integration finishes THEN the run is reported as failed.
5. **[Edge Case]** GIVEN a commit pushed to any other branch (for example `feat/x`) that is not proposed for merge into `main` or a version line WHEN it is pushed THEN no continuous integration run starts.

**UI/UX Specifics:**
- **View/Component:** Continuous integration checks on commits and merge proposals.
- **Feedback Logic:** Each operating system's test result is reported as its own check, so a Windows-only failure is identifiable by name.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Performance]:** The coverage gate and the release tooling tests run once per run, not once per operating system.
- **[Integrity]:** Branch matching covers any future `<major>.x` branch without editing the configuration.
- **[Observability]:** The coverage gate prints the measured core average before passing or failing.

### [US-2] Release gate before anything is published
**User Story:** AS A release maintainer, I WANT the release run to stop before publishing anything when the release tag, the package manifest version and the version built into the binary disagree, or when lint or tests fail, SO THAT no broken or mislabelled version ever reaches users.

**Scenarios:**
1. **[Happy Path]** GIVEN the tag `v1.0.0`, a package manifest at `1.0.0` and a built-in version `1.0.0` WHEN the release runs THEN the version check passes, lint and tests run, and only after both pass does building and publishing start.
2. **[Edge Case]** GIVEN the tag `v1.0.0` and a package manifest still at `1.0.0-rc.1` WHEN the release runs THEN it stops before building or publishing anything, and the failure message names the three values (tag, manifest version, built-in version).
3. **[Edge Case]** GIVEN a tag that is not a valid semantic version (for example `vnext`) WHEN the release runs THEN it stops before building or publishing, stating the tag is not a valid version.
4. **[Edge Case]** GIVEN the versions agree but a unit test fails WHEN the release runs THEN nothing is built or published.

**UI/UX Specifics:**
- **View/Component:** Release run log.
- **Feedback Logic:** A mismatch fails with one line listing `tag`, `package.json` and built-in values side by side.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Performance]:** The gate runs lint and unit tests only (no mutation testing, no security scans).
- **[Reliability]:** Fail fast: the publishing stage never starts when the gate fails.
- **[Integrity]:** The tag is compared without its leading `v`; any other difference is a mismatch.

### [US-3] Correct release channel for every version
**User Story:** AS A Specforce user, I WANT stable versions to reach only the stable channel and pre-releases to reach only their own pre-release channel, SO THAT I never receive a pre-release I did not opt into.

**Scenarios:**
1. **[Happy Path]** GIVEN the stable version `1.0.0` WHEN it is published THEN the npm packages are published under the `latest` channel and the GitHub release is marked as the latest release.
2. **[Happy Path]** GIVEN the pre-release `1.0.0-rc.1` WHEN it is published THEN the npm packages are published under the `rc` channel, the `latest` channel is unchanged, and the GitHub release is marked as a pre-release (so it is never offered to stable native installs).
3. **[Edge Case]** GIVEN a pre-release with an identifier the pipeline never saw before, such as `1.1.0-next.1` WHEN it is published THEN its channel is `next` (the first pre-release identifier), never `latest`.
4. **[Edge Case]** GIVEN a pre-release whose first identifier is numeric, such as `1.1.0-1` WHEN it is published THEN the release stops before publishing, stating that the pre-release needs a named channel.
5. **[Happy Path]** GIVEN any released binary WHEN a user asks it for its version THEN it reports exactly the released version (the tag without `v`).

**UI/UX Specifics:**
- **View/Component:** npm registry channels (dist-tags) and the GitHub releases page.
- **Feedback Logic:** The release log states which channel each package is published to.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Performance]:** Channel derivation is a pure computation with no network access.
- **[Integrity]:** The channel rule is identical to the auto-updater's channel rule (first pre-release identifier, or `latest` for stable versions), so publishing and updating never diverge.
- **[Integrity]:** The version a binary reports comes from the release build itself, not only from a value committed to the repository.

### [US-4] Complete or safely resumable publishing
**User Story:** AS A release maintainer, I WANT a failed publish to stop before the main package is published and a re-run to finish the job without errors on what was already published, SO THAT users never install a main package whose platform binaries are missing.

**Scenarios:**
1. **[Happy Path]** GIVEN none of the five platform packages nor the main package exist at the release version WHEN the release publishes THEN the five platform packages are published first and the main package last.
2. **[Edge Case]** GIVEN one platform package fails to publish for a reason other than "already published" WHEN the release publishes THEN the run fails and the main package is not published.
3. **[Edge Case]** GIVEN a previous run published three of five platform packages and failed WHEN the release is run again for the same tag THEN the three existing packages are skipped with a log line each, the remaining two are published, and then the main package is published.
4. **[Edge Case]** GIVEN every package, including the main package, already exists at the release version WHEN the release is run again THEN every package is skipped and the run succeeds without publishing anything.
5. **[Edge Case]** GIVEN the registry cannot be queried to know whether a package exists WHEN the release publishes THEN the run fails without publishing that package (it is never assumed absent or present).
6. **[Edge Case]** GIVEN a previous run already created the GitHub release for the tag and uploaded some or all of its binaries and checksums WHEN the release is run again for the same tag THEN the GitHub release step succeeds by replacing the existing files instead of failing, and the run continues to the npm publishing stage.
7. **[Edge Case]** GIVEN the same commit is built twice (first run and re-run) WHEN the binaries are compared THEN they are byte-identical, so the binaries already published to npm, the replaced GitHub binaries and the checksums file all describe the same files.

**UI/UX Specifics:**
- **View/Component:** Release run log.
- **Feedback Logic:** One line per package: `published <name>@<version> (<channel>)` or `skipped <name>@<version> (already published)`.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Performance]:** One existence query per package before publishing it.
- **[Reliability]:** No publishing error is swallowed; any failure other than "already exists" fails the run.
- **[Integrity]:** Publishing order is always platform packages first, main package last.
- **[Integrity]:** Release builds are reproducible: building the same commit twice yields byte-identical binaries (no machine paths or build-time timestamps embedded).

### [US-5] Verifiable package origin
**User Story:** AS A Specforce user, I WANT every published npm package to carry a provenance attestation, SO THAT I can verify it was built from this repository and commit by the official release pipeline.

**Scenarios:**
1. **[Happy Path]** GIVEN a release run WHEN any npm package (platform or main) is published THEN it is published with a provenance attestation and public access.
2. **[Edge Case]** GIVEN the provenance attestation cannot be produced WHEN a package is published THEN publishing that package fails and the run fails; the package is never published without provenance.

**UI/UX Specifics:**
- **View/Component:** npm package page provenance badge.
- **Feedback Logic:** The registry shows the source repository and commit of each version.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Performance]:** No extra publish step; provenance is part of each publish.
- **[Safety & Security]:** The release run is granted only the additional permission needed to request an identity token; no new secrets.

### [US-6] Hand-written release notes when provided
**User Story:** AS A release maintainer, I WANT the release page to use the hand-written notes file for the tag when it exists, SO THAT users read migration instructions without me editing the release by hand.

**Scenarios:**
1. **[Happy Path]** GIVEN the file `.github/release-notes/v1.0.0.md` exists with content WHEN `v1.0.0` is released THEN the GitHub release page body is that file's content.
2. **[Edge Case]** GIVEN no notes file exists for the tag WHEN the release runs THEN the release page body is the automatically generated change list, as today.
3. **[Edge Case]** GIVEN the notes file for the tag exists but is empty or contains only whitespace WHEN the release runs THEN it is treated as absent and the generated change list is used.

**UI/UX Specifics:**
- **View/Component:** GitHub release page.
- **Feedback Logic:** The release log states whether the notes file or the generated list was used.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Performance]:** A single file lookup.
- **[Integrity]:** The notes file name is the full tag including the leading `v`.

## 4. Business Invariants
- Nothing is published unless the tag, the package manifest version and the built-in version are equal and lint and unit tests passed in the same run.
- A stable version is never published to a pre-release channel, and a pre-release is never published to `latest` nor marked as the latest GitHub release.
- The main npm package is never published while any platform package of the same version is missing from the registry.
- Every npm package version published by the pipeline carries provenance.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** The release gate adds only lint and unit tests to the release run.
- **[Reliability]:** Fail fast with a clear one-line reason; no silent failures (no ignored publishing errors).
- **[Security]:** No new secrets; the only new permission is identity-token issuance for provenance.
- **[Maintainability]:** Release decision logic (version agreement, channel derivation, publish-or-skip, notes selection) is covered by automated tests that run in continuous integration, with no new third-party dependencies.
