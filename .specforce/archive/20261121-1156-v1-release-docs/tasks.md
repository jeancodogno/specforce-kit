---
slug: 20261121-1156-v1-release-docs
lens: Integration
---

# Implementation Roadmap: Documentation for the Stable 1.0.0 Release

## 1. Execution Strategy
- **Scope (small spec, no requirements file):** the goals below are the requirement references used by every task.
  - **[US-1] Leaving a pre-release channel:** the Updates section of the getting-started guide (en/pt/es) explains that installs of the already published pre-releases (up to `1.0.0-beta.3`, which still run the previous updater) move to 1.0.0 automatically; that from 1.0.0 onward npm pre-release users switch to stable by hand with `npm install -g @jeancodogno/specforce-kit@latest` and are told once per stable version; and that native installs follow stable releases only (a pre-release is installed natively only by downloading it from the releases page).
  - **[US-2] Changelog:** a `CHANGELOG.md` in Keep a Changelog format with a single `1.0.0` entry (changes since 0.2.x plus a "Migrating from 0.2.x" section); older history links to the GitHub releases page.
  - **[US-3] Security policy:** a `SECURITY.md` stating supported versions, reporting through GitHub private vulnerability reporting (no e-mail address), and how the updater and the release protect integrity.
  - **[US-4] Release notes draft:** `.github/release-notes/v1.0.0.md`, which the release pipeline publishes as the release page body (`release-pipeline-hardening` [US-6]).
- **Decisions (approved):** `CHANGELOG.md`, `SECURITY.md` and the release notes are English only (repository-convention files, outside `docs/`); the Updates section keeps translation parity in `docs/{en,pt,es}/` per `governance.md`.
- **Dependencies:** T1.1 describes the behaviour of `beta-stable-notice` and T1.5 describes provenance from `release-pipeline-hardening`; implement this spec after those two (or confirm their behaviour is still the approved one before writing).
- **Test pairing:** documentation only, no versioned verification tooling, so the project's RED/GREEN pairing does not apply (same precedent as the archived `docs-beta3-refresh` spec); every task is proven by a non-interactive shell check on the files it changes.
- **Gravity Order:** English Updates section -> its pt and es mirrors -> the three English-only repository files, which are independent of each other.

## 2. Tasks

### Phase 1: Updates section (three languages)

- [x] T1.1: [DOCS] Leaving a pre-release channel, in English
**Target:** `docs/en/getting-started.md`
**Context:** [US-1]

**Constraints:**
- Edit only the existing `## Updates` section; keep the rest of the page unchanged. Replace the native-install wording so it says native installs follow stable releases only.
- Quote the notice sentence in the same format as `beta-stable-notice` requirements [US-1] scenario 1, but with a post-1.0.0 example (stable `1.1.0` on a `1.1.0-beta.x` install), because published pre-releases up to `1.0.0-beta.3` cannot show it; state it never installs anything.

**Acceptance Check:**
`f=docs/en/getting-started.md; s=$(sed -n '/^## Updates/,/^## First Steps/p' $f); echo "$s" | grep -q 'specforce-kit@latest' && echo "$s" | grep -q 'STABLE RELEASE AVAILABLE\|(stable) is available' && echo "$s" | grep -qi 'stable releases only' && echo OK` prints `OK`.

- [x] T1.2: [DOCS] Portuguese mirror of the Updates section
**Target:** `docs/pt/getting-started.md`
**Context:** [US-1]
**Parallel With:** T1.3

**Constraints:**
- Translate from the final English text of T1.1; commands, the package name and the quoted notice sentence stay in English. The page keeps the same line count as the English page.

**Acceptance Check:**
`f=docs/pt/getting-started.md; grep -q 'specforce-kit@latest' $f && [ "$(wc -l < $f)" -eq "$(wc -l < docs/en/getting-started.md)" ] && echo OK` prints `OK`.

- [x] T1.3: [DOCS] Spanish mirror of the Updates section
**Target:** `docs/es/getting-started.md`
**Context:** [US-1]
**Parallel With:** T1.2

**Constraints:**
- Same rules as T1.2, in Spanish.

**Acceptance Check:**
`f=docs/es/getting-started.md; grep -q 'specforce-kit@latest' $f && [ "$(wc -l < $f)" -eq "$(wc -l < docs/en/getting-started.md)" ] && echo OK` prints `OK`.

### Phase 2: Repository files (English only)

- [x] T2.1: [DOCS] Changelog with the 1.0.0 entry and migration guide
**Target:** `CHANGELOG.md`
**Context:** [US-2]
**Parallel With:** T2.2, T2.3

**Constraints:**
- Heading `## [1.0.0] - Unreleased` (the maintainer sets the date when tagging), with Added / Changed / Removed groups.
- Facts to cover, verified against the code before writing: slash commands and subagent definitions replaced by native skills (`spf.discovery`, `spf.spec`, `spf.constitution`, `spf.implement`, `spf.archive`, `consultative-grill`); project memorial removed in favour of module living specs; Gemini CLI support removed, Cursor added; tiered spec sizes, module deltas and independent QA added; background auto-update added.
- "Migrating from 0.2.x": re-run `specforce init` (it detects legacy assets and asks before removing them, `[BR-KIT-03]`), then invoke skills by name instead of the old slash commands.
- End with a link to `https://github.com/jeancodogno/specforce-kit/releases` for versions before 1.0.0.

**Acceptance Check:**
`f=CHANGELOG.md; grep -q 'keepachangelog' $f && grep -q '^## \[1.0.0\]' $f && grep -q 'Migrating from 0.2' $f && grep -q 'specforce init' $f && grep -qi 'memorial' $f && grep -qi 'gemini' $f && grep -q 'releases' $f && echo OK` prints `OK`.

- [x] T2.2: [DOCS] Security policy
**Target:** `SECURITY.md`
**Context:** [US-3]
**Parallel With:** T2.1, T2.3

**Constraints:**
- Supported versions: latest `1.x` stable only; pre-releases and `0.x` unsupported.
- Reporting only through GitHub private vulnerability reporting (the Security tab's "Report a vulnerability"); no e-mail address and no public issues for vulnerabilities.
- Integrity section: native updates are verified against the release checksum before being applied, never downgrade, and keep owner-only update state (`.specforce/docs/modules/auto-update.md`); npm packages are published with provenance.

**Acceptance Check:**
`f=SECURITY.md; grep -qi 'report a vulnerability' $f && grep -qi 'checksum' $f && grep -qi 'provenance' $f && ! grep -qE '[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[a-z]{2,}' $f && echo OK` prints `OK`.

- [x] T2.3: [DOCS] Release notes draft for v1.0.0
**Target:** `.github/release-notes/v1.0.0.md`
**Context:** [US-4], [US-1]
**Parallel With:** T2.1, T2.2

**Constraints:**
- Opens with the migration block: existing pre-release installs (npm or native, up to `1.0.0-beta.3`) update to 1.0.0 automatically, with `npm install -g @jeancodogno/specforce-kit@latest` as the optional way to get it right away on npm; 0.2.x users follow the changelog migration section.
- Then a short highlights list and a link to `CHANGELOG.md`; no full duplication of the changelog.

**Acceptance Check:**
`f=.github/release-notes/v1.0.0.md; head -15 $f | grep -q 'specforce-kit@latest' && grep -q 'CHANGELOG.md' $f && echo OK` prints `OK`.
