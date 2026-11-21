# Proposal: Auto-update Reliability

**Suggested type/size:** bug / medium
**Origin:** `/spf:discovery` session on 2026-11-17 (user report: "auto-update does not seem to work"). Every finding below marked **[REPRODUCED]** was observed end to end against real builds (`go build -ldflags "-X core.Version=..."`) with an isolated `HOME`; probes were disposable and have been removed.

## Problem
The background updater downloads and stages new releases correctly, but the apply step (swap) and the npm path are unsafe or broken. Users installed with `sudo make install` (binary in root-owned `/usr/local/bin`) never get updated, are told it happens "automatically", and the CLI retries the failing swap on every invocation. Other paths can downgrade the CLI or corrupt agent JSON output.

## How it works today
```
PersistentPreRunE (TTY && !agent cmd)
  ↳ CheckForUpdate: if now - last_checked_at >= 6h → spawn detached `specforce --internal-upgrade-check`
       ↳ PerformBackgroundUpdate: GitHub /releases/latest → if newer: stageUpdate
            ↳ binary: download asset + checksums → ~/.specforce/upgrade/staged/specforce_next
            ↳ npm (isNPM): `npm install -g @jeancodogno/specforce-kit`
            ↳ state.json: staged_version, update_ready=true
main.checkForSwap (EVERY invocation, before Cobra, TTY or not)
  ↳ PerformAtomicSwapAt(os.Executable()) → print "(Specforce updated to X)" to stdout → syscall.Exec
PersistentPostRun (TTY && !agent cmd)
  ↳ IsUpdateAvailable → box "...downloaded in the background and applied automatically."
```

## Verified findings

### F1. Swap fails silently on non-writable install dir, loops forever **[REPRODUCED]**
- `src/internal/upgrade/service.go:171` renames `exe → exe.old`; on the user's machine `strace` showed `renameat("/usr/local/bin/specforce", ".../specforce.old") = -1 EACCES`.
- `src/cmd/specforce/main.go:59-62` swallows the error; `update_ready` stays `true`, so every run retries.
- Sandbox repro (install dir `chmod 555`, current beta.2, staged beta.3, 3 TTY runs): version stays beta.2, `update_ready:true` every run, notification printed every run claiming it will be "applied automatically" (`src/internal/tui/upgrade_notification.go:31`).
- `Makefile` installs to `PREFIX ?= /usr/local` → this is the default install path for source users.

### F2. npm path installs the wrong dist-tag (downgrade to 0.2.2) **[REPRODUCED]**
- `isNPM()` (`service.go:285-294`) matches any exe path containing `node`/`npm`; the npm binary lives under `.../node_modules/@jeancodogno/specforce-kit-<os>-<arch>/bin/`.
- `stageUpdate` → `NPMInstaller.Install` (`installer_npm.go:36`) runs `npm install -g @jeancodogno/specforce-kit` with no version.
- Registry dist-tags today: `latest=0.2.2`, `beta=1.0.0-beta.3`, `alpha=1.0.0-alpha.3`. Version source is GitHub (`root.go:35`, `main.go:77`, both still `TODO`), which reports `v1.0.0-beta.3`.
- Repro with a fake `npm` on PATH logging its args: logged `npm install -g @jeancodogno/specforce-kit`; next run printed `(Specforce updated to v1.0.0-beta.3)` while `--version` still said beta.2. With real npm this installs 0.2.2.
- npm errors are also invisible to the user (stderr/stdout discarded in `RealCommandExecutor`).

### F3. Swap applies a staged version that is older or equal **[REPRODUCED]**
- `checkForSwap`/`PerformAtomicSwapAt` never compare `staged_version` with `core.Version`.
- Repro: installed beta.3, staged beta.2 → after one run `--version` = beta.2.
- Realistic trigger: user upgrades manually (reinstall/`make install`) while an older update is staged.

### F4. Swap notice leaks into stdout, breaks agent `--json` **[REPRODUCED]**
- `main.go:65` uses `fmt.Printf` to stdout, runs before Cobra, ignores TTY and `IsAgentCommand`.
- Repro: `specforce spec list --json 2>/dev/null` at swap time → first stdout lines are `\n\x1b[38;5;245m(Specforce updated to v1.0.0-beta.2)\x1b[0m`; `json.load` fails (`JSONDecodeError`).

### F5. 6h throttle bypassed when the check fails **[REPRODUCED]**
- `PerformBackgroundUpdate` returns before saving `last_checked_at` if `GetLatestVersion` fails (`service.go:84-87`).
- Repro (pty via `script`, 3 runs): online → detached checks spawned `1,0,0`; offline (dead `HTTPS_PROXY`) → `1,1,1`, `state.json` never created. Every interactive command offline spawns a new process doing a 10s-timeout HTTP call.

### F6. Update state is global, not bound to the executable it was staged for **[REPRODUCED, accidentally]**
- `state.json` and `staged/specforce_next` live in `~/.specforce`, keyed by user only; `checkForSwap` swaps into whatever `os.Executable()` is running.
- During discovery, running a locally built probe binary (`go build`, beta.3) without isolating `HOME` consumed the real staged update meant for `/usr/local/bin/specforce` and overwrote itself (sha256 `33e26bc1…` = official `specforce-kit_linux_amd64`). Any dev build (`./specforce` in the repo, `go run`) can steal and overwrite itself with a release binary.

### Code-reading findings (not reproduced)
- **F7.** Temp download leaks on failure: `DownloadAndVerify` only removes the temp file on hash mismatch, not on download/checksum HTTP errors. Observed a 0-byte `/tmp/specforce-update-2412556037` (created 11:00, before the session).
- **F8.** Checksum fetch uses `i.Client.Get` (no context) on an `http.Client{}` without `Timeout` (`installer_binary.go:63`); can hang past the 1-minute background context.
- **F9.** `IgnoredVersion` is read but never written by any command — no way to skip a version. `Service.PerformUpgrade` is dead code.
- **F10.** Staged binary observed with mode `0664`; `security.md` asks for `0600` on upgrade-internal files (staged dir itself is `0700`).
- **F11 (security note).** Checksums come from the same release as the binary: integrity only, no authenticity (no signature). Out of scope unless decided otherwise.

## Decisions taken during discovery (proposed)
1. **Bind staged updates to a target.** Persist `target_path` (resolved `os.Executable()` at stage time) in `state.json`; swap only when running from that path.
2. **Never downgrade.** Swap only if `IsNewer(core.Version, staged_version)`; otherwise discard staged file and clear `update_ready`.
3. **Pre-flight writability.** Before staging (binary mode), check the install dir is writable; if not, do not download — notification switches to a manual instruction (e.g. reinstall / `sudo make install`).
4. **Swap failure is terminal for that version.** On error: clear `update_ready`, remove staged file, record the failure, emit one notice on stderr.
5. **All updater output goes to stderr, and only on TTY and non-agent commands.** The swap notice moves out of `main.go` stdout.
6. **Always persist `last_checked_at`** (success or failure) so the 6h throttle holds offline.
7. **npm: install exact version** `@jeancodogno/specforce-kit@<version>`, detect npm installs by `node_modules/@jeancodogno/specforce-kit` in the resolved path, and surface npm failure.
8. Housekeeping: remove temp file on every failure path, context-aware checksum fetch, `0600` staged file.

## Rejected alternatives
- **Elevate with sudo on swap:** interactive privilege prompts from a background updater are unacceptable and break agents.
- **Moving the default install to `~/.local/bin`:** valid but a separate install/packaging decision; this spec only makes the updater honest about non-writable locations.
- **Using npm `latest` and re-tagging releases:** fixes F2 only by publishing discipline, leaves the code fragile.

## Open questions
1. npm version source: keep GitHub `/releases/latest` for npm installs, or query the npm registry using the dist-tag matching the current channel (alpha/beta/latest)? (Both `TODO`s in `root.go:34` and `main.go:76`.)
2. Channel policy: should a `beta` user ever be offered a stable `latest` that is semver-lower (e.g. 0.2.2)? Proposed: never offer lower versions (already guaranteed by `IsNewer`), but the npm tag choice must follow the channel.
3. Non-writable install: show the manual-update hint every run, or once per version (needs a "notified" marker)?
4. Should `IgnoredVersion` get a command (`specforce upgrade --skip`) or be removed? (Proposed: remove, out of scope.)
5. Windows (`service_windows.go`) not exercised: rename of a running `.exe` behaves differently; does it need its own acceptance test?

## Acceptance hints
Go tests in `src/internal/upgrade` and `src/cmd/specforce` (use temp `HOME` and temp install dirs):
- non-writable install dir → nothing staged, no swap attempt, `update_ready=false`, hint on stderr;
- swap error → state cleared, staged file removed, single stderr notice;
- staged ≤ current → no swap, state cleared;
- staged for another `target_path` → ignored by a different executable;
- agent command with `--json` during a pending swap → stdout is valid JSON;
- provider error → `last_checked_at` saved, second call within 6h does not spawn;
- npm installer called with `@<exact version>`;
- download failure → no temp file left.
`go test ./...` passes.

## Evidence (repro recipe, probes removed)
- Build two versions: `go build -ldflags="-X 'github.com/jeancodogno/specforce-kit/src/internal/core.Version=1.0.0-beta.2'" -o specforce-1.0.0-beta.2 src/cmd/specforce/main.go` (same for beta.3).
- Always run with `HOME=<tmp>` (see F6). TTY scenarios via `script -qec "<cmd>" /dev/null` with stdout not redirected inside the pty.
- F1: copy beta.2 to `<tmp>/install/specforce`, `chmod 555 <tmp>/install`, stage beta.3 in `$HOME/.specforce/upgrade/staged/specforce_next`, write `state.json` with `update_ready:true`, run 3 times.
- F2: place beta.2 at `<tmp>/prefix/lib/node_modules/@jeancodogno/specforce-kit-linux-x64/bin/specforce`, put a fake `npm` script logging `$*` first on `PATH`, run `--internal-upgrade-check`.
- F3/F4: install beta.3 in a writable dir, stage beta.2, run `spec list --json 2>/dev/null` and parse stdout.
- F5: `strace -f -e trace=execve` under `script`, `HTTPS_PROXY=http://127.0.0.1:9`, count `internal-upgrade-check` per run.
- Real machine: `strace -f -e trace=%file specforce --version` → `renameat(...) = -1 EACCES`.

## How to use this proposal
`/spf:spec` to produce requirements, design and tasks (settle open questions 1–3 first), then `/spf:implement`.
