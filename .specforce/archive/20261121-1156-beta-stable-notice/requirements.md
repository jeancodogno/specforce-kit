---
slug: 20261121-1156-beta-stable-notice
lens: Backend-heavy
---

# Feature: Stable Release Notice for Pre-release npm Users

## 1. Context & Value
npm installations follow the release channel of the version they run, so a user on `1.0.0-beta.3` keeps receiving betas and never moves to the stable `1.0.0` on their own. Switching channels stays a deliberate manual action (decided policy), but today these users get no signal at all that a stable version exists. This feature tells them once per stable version, with the exact command to switch, without installing anything.

## 2. Out of Scope (Anti-Goals)
- Do not switch channels automatically and do not install the stable version; the user always switches by hand.
- Do not change how the running channel itself is updated (pre-release users keep receiving their channel's updates as today).
- Do not show this notice to native binary installations; they already follow stable releases only.
- Do not add a new command, flag or setting to disable or snooze the notice.
- Do not change the release pipeline or the documentation (separate specifications `release-pipeline-hardening` and `v1-release-docs`).

## 3. Acceptance Criteria (BDD)

### [US-1] Learn that a stable version exists while on a pre-release channel
**User Story:** AS A Specforce user who installed a pre-release through npm, I WANT to be told once when a stable version newer than mine is published, with the exact command to switch, SO THAT I can leave the pre-release channel when I choose to.

**Scenarios:**
1. **[Happy Path]** GIVEN an npm installation running `1.0.0-beta.3` and the stable channel at `1.0.0` WHEN a background check has run and the user then runs an interactive, non-agent command THEN the standard error shows, once: `Specforce 1.0.0 (stable) is available. You are on the beta channel; to switch, run: npm install -g @jeancodogno/specforce-kit@latest`.
2. **[Happy Path]** GIVEN the notice for `1.0.0` was already shown WHEN later interactive commands run THEN it is not shown again.
3. **[Happy Path]** GIVEN the notice for `1.0.0` was shown and the stable channel later moves to `1.0.1` WHEN the next check and interactive command run THEN the notice is shown once for `1.0.1`.
4. **[Edge Case]** GIVEN an npm installation running `1.1.0-beta.1` and the stable channel at `1.0.2` (older than the running version) WHEN checks and commands run THEN no stable notice is ever shown.
5. **[Edge Case]** GIVEN the running pre-release channel no longer exists in the registry but the stable channel is newer than the running version WHEN the check runs THEN nothing is installed and the stable notice is still shown once.
6. **[Edge Case]** GIVEN the registry cannot be reached during a check WHEN the check finishes THEN the recorded stable information is left exactly as it was (a pending notice stays pending, a shown one stays shown).
7. **[Edge Case]** GIVEN the user already switched and now runs the stable `1.0.0` WHEN interactive commands run THEN no stable notice is shown for `1.0.0` or any older version.
8. **[Edge Case]** GIVEN a native binary installation, or an npm installation already on the stable channel WHEN checks and commands run THEN the stable notice is never recorded nor shown.

**UI/UX Specifics:**
- **View/Component:** The same bordered notice box used by the other update notices, titled `STABLE RELEASE AVAILABLE`.
- **Feedback Logic:** The text names the stable version, the running channel and the exact switch command; it never says the update will happen automatically.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Performance]:** No extra network request: the stable version is read from the same registry response the channel check already fetches.
- **[Integrity]:** "Newer" uses the same semantic version ordering as the rest of the updater (`1.0.0` is newer than `1.0.0-beta.3`).
- **[Safety & Security]:** The recorded stable version lives in the existing owner-only (`0600`) per-user update state; older state files without it load unchanged.
- **[Observability]:** With `DEBUG=1`, the check writes one line stating whether a newer stable version was recorded.

### [US-2] One notice at a time, in a predictable order
**User Story:** AS A Specforce user, I WANT at most one update notice per command, the most important first, SO THAT notices stay readable and urgent ones are never hidden.

**Scenarios:**
1. **[Happy Path]** GIVEN an "updated to X" notice and a stable notice are both pending WHEN an interactive command runs THEN only the "updated to X" notice is shown, and the stable notice is shown by a later command.
2. **[Happy Path]** GIVEN a "manual update required" notice and a stable notice are both pending WHEN an interactive command runs THEN only the manual notice is shown.
3. **[Happy Path]** GIVEN a stable notice is pending and a channel update is also available WHEN an interactive command runs THEN the stable notice is shown instead of the "update available" notice.
4. **[Edge Case]** GIVEN a stable notice is pending WHEN a command is run by an agent, or with standard output not being a terminal THEN no notice is written and the stable notice stays pending (it is not marked as shown).

**UI/UX Specifics:**
- **View/Component:** Update notice on standard error.
- **Feedback Logic:** Precedence is: updated > manual update required > stable available > update available.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Performance]:** Choosing the notice reads only the local update state.
- **[Integrity]:** A notice is marked as shown only after it was actually written.
- **[Reliability]:** Standard output (including `--json` output) is never written to.

## 4. Business Invariants
- The stable notice never causes any installation, download or channel change.
- The stable notice is shown at most once per stable version and only while the running version is an npm pre-release older than that stable version.
- At most one update notice is written per command, and never for agent commands or non-terminal output.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** No additional network requests and no added latency on the foreground command beyond reading the local state.
- **[Reliability]:** Any failure in recording or showing the notice is silent and never fails the user's command.
- **[Security]:** Update state keeps its owner-only permissions; unknown or missing keys are tolerated.
- **[Maintainability]:** Core packages keep at least 80% average coverage; new behaviour is covered by service-level and command-level tests.
