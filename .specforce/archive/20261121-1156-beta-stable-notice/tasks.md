---
slug: 20261121-1156-beta-stable-notice
lens: Backend-heavy
---

# Implementation Roadmap: Stable Release Notice for Pre-release npm Users

## 1. Execution Strategy
- **Gravity Order:** Registry provider exposes the stable tag -> background check records it in the update state -> notice selection and "shown once" persistence -> rendering and wiring in the root command.
- **Architecture:** All rules live in `src/internal/upgrade` (service layer); the CLI layer (`src/internal/cli/cobra/root.go`) and the TUI renderer (`src/internal/tui`) only map the new notice kind to text, per the existing notice pattern.
- **Delta coverage:** `[NEW:stable-notice]` and `[NEW:us-stable-notice]` -> Phases 1-3; `[BR-UPD-11]` and `[US-UPD-04]` (channel without tag) -> T1.3/T1.4; `[BR-UPD-12]` (precedence, once per version, shown only after written) -> Phases 2-3.

## 2. Tasks

### Phase 1: Record the newer stable version during the check

- [x] T1.1: [RED] Provider exposes the stable tag from the same response
**Target:** `src/internal/upgrade/provider_internal_test.go`
**Context:** [US-1]

**Constraints:**
- Cover: channel `beta` with tags `{beta: 1.0.0-beta.3, latest: 1.0.0}` returns `1.0.0-beta.3` and exposes stable `1.0.0`; missing `beta` tag still returns the `ErrNoChannelTag` error AND exposes stable `1.0.0`; registry error exposes no stable version; exactly one HTTP request per check (count requests on the `httptest` server).

**Acceptance Check:**
`go test ./src/internal/upgrade/ -run TestNPMProvider_Stable` fails to compile or fails (RED state).

- [x] T1.2: [GREEN] Capture the `latest` dist-tag in the npm provider
**Target:** `src/internal/upgrade/provider_npm.go`
**Context:** [US-1]

**Constraints:**
- Keep the `Provider` interface unchanged; expose the stable tag through a separate small interface (for example a `LatestStable() string` method) that the service detects with a type assertion, so `GitHubProvider` and `MockProvider` need no change.
- The value is reset at the start of every `GetLatestVersion` call so a failed fetch never exposes a stale stable version.

**Acceptance Check:**
`go test ./src/internal/upgrade/ -run 'TestNPMProvider'` passes.

- [x] T1.3: [RED] Background check records the newer stable version
**Target:** `src/internal/upgrade/service_test.go`
**Context:** [US-1]

**Constraints:**
- Cover with an npm-install service: running `1.0.0-beta.3` + stable `1.0.0` records it and installs only the channel version (npm runner called only for `@1.0.0-beta.x`, never for `@latest`); channel tag missing + stable newer records it without any install; running `1.1.0-beta.1` + stable `1.0.2` records nothing; running stable `1.0.0` records nothing; native install never records; provider error leaves `StableVersion` and `StableNoticeShownVersion` unchanged.

**Acceptance Check:**
`go test ./src/internal/upgrade/ -run TestStableRecord` fails (RED state).

- [x] T1.4: [GREEN] Persist the stable version in the update state
**Target:** `src/internal/upgrade/service.go`, `src/internal/upgrade/state.go`
**Context:** [US-1]

**Constraints:**
- New state keys `stable_version` and `stable_notice_shown_version` (`omitempty`), so older state files load unchanged (`[BR-UPD-14]` tolerance).
- Record only when `isNPM()`, `ChannelOf(currentVer) != "latest"` and the stable version `IsNewer` than the running version; evaluate it on both the success path and the `ErrNoChannelTag` path of `PerformBackgroundUpdate`, never on the generic error path.
- Emit one `debugf` line stating whether a newer stable version was recorded.

**Acceptance Check:**
`go test ./src/internal/upgrade/` passes.

### Phase 2: Choose and remember the stable notice

- [x] T2.1: [RED] Notice precedence and once-per-version behaviour
**Target:** `src/internal/upgrade/service_test.go`
**Context:** [US-1], [US-2]

**Constraints:**
- Cover: applied beats stable; manual beats stable; stable beats available; stable shown once (after `MarkNoticeShown` the same version is no longer pending) and pending again for a newer stable version; no stable notice when the running version is equal to or newer than the recorded stable version; `MarkNoticeShown` for stable persists `StableNoticeShownVersion`.

**Acceptance Check:**
`go test ./src/internal/upgrade/ -run TestNotice_Stable` fails (RED state).

- [x] T2.2: [GREEN] `NoticeStable` kind in notice selection
**Target:** `src/internal/upgrade/notice.go`
**Context:** [US-1], [US-2]

**Constraints:**
- Append `NoticeStable` after the existing `NoticeKind` constants so existing values do not change.
- Precedence in `PendingNotice`: applied, manual, stable, available.
- `PendingNotice` does not re-check the installation method: recording is already gated by T1.4, and the root tests use a non-npm executable path.

**Acceptance Check:**
`go test ./src/internal/upgrade/` passes.

### Phase 3: Show the notice

- [x] T3.1: [RED] Rendering and root-command wiring of the stable notice
**Target:** `src/internal/tui/upgrade_notification_test.go`, `src/internal/cli/cobra/root_notice_test.go`
**Context:** [US-1], [US-2]

**Constraints:**
- Renderer test: output contains the title `STABLE RELEASE AVAILABLE` and, ignoring ANSI styling and line wrapping, the exact sentence from requirements [US-1] scenario 1 for version `1.0.0` and channel `beta`; it does not contain "automatically".
- Root test: add `stable ...` subtests inside `TestRootUpdateNotice`, reusing `newNoticeCmd` and a state seeded with `StableVersion` and no `ManualUpdateVersion` (the existing `setupNoticeTest` seed would make the manual notice win): a pending stable notice is written to stderr for an interactive non-agent command and marked shown; for agent and non-TTY commands nothing is written and `StableNoticeShownVersion` stays empty; stdout stays empty in every case.

**Acceptance Check:**
`go test ./src/internal/tui/ ./src/internal/cli/cobra/ -run 'Stable'` fails (RED state).

- [x] T3.2: [GREEN] Stable notice renderer and wiring
**Target:** `src/internal/tui/upgrade_notification.go`, `src/internal/cli/cobra/root.go`
**Context:** [US-1], [US-2]

**Constraints:**
- Same box style as `RenderManualUpdateNotification`; the channel shown is `upgrade.ChannelOf(appVersion)`; the command is `npm install -g @jeancodogno/specforce-kit@latest`.

**Acceptance Check:**
`go test ./src/internal/tui/ ./src/internal/cli/cobra/ ./src/internal/upgrade/` passes.

## 3. Black-box Scenarios

- **S1:** Pending stable notice is shown once to an interactive user
**Context:** [US-1], [US-2]
**Command:** `go test ./src/internal/cli/cobra/ -run 'TestRootUpdateNotice/stable' -v`
**Expected:** The subtests for interactive (notice written once, then nothing on the second run), agent and non-TTY commands pass.

- **S2:** Pre-release check never installs the stable channel
**Context:** [US-1]
**Command:** `go test ./src/internal/upgrade/ -run TestStableRecord -v`
**Expected:** All subtests pass, including the one asserting that the npm runner was never called with `@latest`.
