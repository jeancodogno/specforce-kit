---
slug: optional-native-packages
lens: Migration
---

# Feature: Optional Native Packages

## 1. Context & Value
The current distribution of Specforce Kit relies on `go-npm`, which downloads binaries from GitHub during the `postinstall` phase, creating a supply-chain risk. By migrating to "NPM Native Binaries" via `optionalDependencies`, we ensure that binaries are delivered through the NPM registry, verified by shasums, and installed natively by the package manager based on the user's OS and CPU architecture.

## 2. Out of Scope (Anti-Goals)
*Explicitly state what MUST NOT be built in this specific feature to prevent scope creep.*
- Supporting automatic binary compilation from source during installation.
- Building a custom binary downloader or updater.
- Implementing code signing/notarization for macOS.

## 3. Acceptance Criteria (BDD)
*Every requirement MUST be testable. Define success and failure.*

### [REQ-1] Platform-Specific Sub-packages
**User Story:** AS A maintainer, I WANT TO publish architecture-specific packages, SO THAT users only download the binary relevant to their system.

**Scenarios:**
1. **[Happy Path]** GIVEN a release build WHEN the publishing pipeline runs THEN it MUST generate separate NPM packages for linux-x64, linux-arm64, darwin-x64, darwin-arm64, win32-x64.
2. **[Edge Case]** GIVEN an unsupported OS/Architecture combination WHEN the installer runs THEN it MUST report a clear error that the platform is not supported by native binaries.
3. **[Ambiguity]** We assume standard naming conventions `@jeancodogno/specforce-kit-{os}-{arch}` for the sub-packages.

### [REQ-2] Zero-Script Installation
**User Story:** AS A security-conscious user, I WANT TO install Specforce Kit with `--ignore-scripts`, SO THAT I can be sure no arbitrary code runs during installation.

**Scenarios:**
1. **[Happy Path]** GIVEN the main `@jeancodogno/specforce-kit` package WHEN `npm install` is executed THEN the package manager MUST resolve and install the correct native sub-package via `optionalDependencies` without calling any `postinstall` scripts.
2. **[Edge Case]** GIVEN an environment where `optionalDependencies` are skipped WHEN the user tries to run the `specforce` command THEN the system MUST provide a helpful instruction to manually install the required native package.
3. **[Ambiguity]** We assume npm, pnpm, and yarn will all correctly handle `os` and `cpu` fields in the `package.json`.

## 4. Business Invariants
*Strict domain rules that the system must never violate, regardless of the user story. These must be technically testable.*
- All binaries must be hosted within the NPM Registry; no external network requests to GitHub are allowed during installation.
- If a binary for the current architecture is missing, the installation of the main package should succeed, but the CLI proxy must fail with a clear diagnostic message.
