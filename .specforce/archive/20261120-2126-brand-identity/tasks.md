---
slug: 20261120-2126-brand-identity
lens: UI-heavy
---

# Implementation Roadmap: Brand Identity Refresh (Gate symbol, wordmark and adaptive palette)

## 1. Execution Strategy
- **Gravity Order:** Asset pipeline (isolated Go module, produces the files every other surface points to) -> color environment (adaptive palette and start-up configuration, which every renderer depends on) -> terminal brand forms (full mark, compact header, footer, call sites) -> READMEs, package metadata and constitution. Each code phase starts with a `[RED]` test task.

## 2. Tasks

### Phase 1: Brand asset pipeline

- [x] T1.1: [RED] Scaffold the brandgen module with failing tests
**Target:** `scripts/brandgen/go.mod`, `scripts/brandgen/brandgen_test.go`
**Context:** [US-1]

**Constraints:**
- Separate module per the design inventory entry for `scripts/brandgen/go.mod`; do not touch the root `go.mod`.
- Tests cover: the exact 27-file output list of design table 2.5; two consecutive builds produce byte-identical outputs; the 16 px app icon has background-colored pixels between the left bracket's inner edge and the check on row 7 (index (16−1)/2), where a pixel counts as background when every channel is ≤ 0x0B+0x30; a write into a read-only directory exits non-zero, names the file, and leaves a pre-existing file in that directory unchanged; generated SVGs contain no `<text`, `stroke` or `Gradient`.

**Acceptance Check:**
`cd scripts/brandgen && go test ./...` exits non-zero and every listed test is reported (compile failure on undefined generator symbols is acceptable at this stage).

- [x] T1.2: [CODE] Brand definition, wordmark outlines and SVG serialization
**Target:** `scripts/brandgen/brand.go`, `scripts/brandgen/wordmark.go`, `scripts/brandgen/svg.go`
**Context:** [US-1]

**Constraints:**
- Geometry, palettes and output table exactly per design sections 2.1–2.5; outline the wordmark per 2.3 (no SVG `<text>`).

**Acceptance Check:**
`cd scripts/brandgen && go test -run 'SVG|OutputList|Deterministic' ./...` passes.

- [x] T1.3: [CODE] Rasterizer, PNG encoding and generator entry point
**Target:** `scripts/brandgen/raster.go`, `scripts/brandgen/main.go`
**Context:** [US-1]

**Constraints:**
- Build-everything-in-memory, then temp-file-plus-rename writes, per design 2.6.

**Acceptance Check:**
`cd scripts/brandgen && go test ./...` passes (all tests from T1.1 green).

- [x] T1.4: [CLI] Add the brand-assets target, generate the asset set and remove legacy logos
**Target:** `Makefile`, `assets/brand/`, `assets/logo.png` + `assets/logo-header.png` (deleted)
**Context:** [US-1]

**Acceptance Check:**
`make brand-assets && test "$(ls assets/brand | wc -l)" -eq 27 && test ! -e assets/logo.png && test ! -e assets/logo-header.png` succeeds.

### Phase 2: Color environment

- [x] T2.1: [RED] Palette and color-environment tests
**Target:** `src/internal/tui/theme_test.go`, `src/internal/tui/colorenv_test.go`
**Context:** [US-4], [US-5]

**Constraints:**
- Cover: every palette variable of design table 4.1 has the listed Light/Dark values; every light text color has WCAG contrast ≥ 4.5:1 on `#FFFFFF`; `COLORFGBG` precedence table of design 4.2 (including `0;15` light, `15;0` dark, `default;default` falls through, three-field form); non-interactive returns dark without invoking the query; `NO_COLOR=1` sets the Ascii profile and never invokes the query; empty `NO_COLOR` keeps colors; with a light background and the TrueColor profile, primary text renders with no foreground color sequence; no `lipgloss.Color(` literal outside `theme.go` in non-test files of `src/internal/tui`.
- Inject the query function; tests must never send a real terminal query. Restore the global lipgloss profile/background after each test.

**Acceptance Check:**
`go test ./src/internal/tui -run 'Palette|Contrast|ConfigureColor|DetectDarkBackground|NoStrayColor'` exits non-zero with those tests failing or not compiling.

- [x] T2.2: [CODE] Adaptive palette
**Target:** `src/internal/tui/theme.go`, `src/internal/tui/components.go`
**Context:** [US-4]

**Constraints:**
- Per design 4.1, including the renames (`textWhite` → `textPrimary`, `black` → `badgeOnFill`, `inProgressOrange` → `inProgressYellow`) and the new `footerBg` used by `FooterStyle`.

**Acceptance Check:**
`go test ./src/internal/tui -run 'Palette|Contrast'` passes and `go build ./...` succeeds.

- [x] T2.3: [CODE] Map upgrade-screen colors to palette roles
**Target:** `src/internal/tui/upgrade_notification.go`, `src/internal/tui/upgrade_prompt.go`, `src/internal/tui/upgrade_progress.go`
**Context:** [US-4]

**Acceptance Check:**
`go test ./src/internal/tui -run 'NoStrayColor|Upgrade'` passes.

- [x] T2.4: [CODE] Start-up color configuration
**Target:** `src/internal/tui/colorenv.go`, `src/internal/cli/cobra/root.go`
**Context:** [US-4], [US-5]

**Constraints:**
- Call order per design 4.2: configuration runs in `Execute` before the command tree executes, so it precedes every output including upgrade notices.

**Acceptance Check:**
`go test ./src/internal/tui ./src/internal/cli/... -run 'ConfigureColor|DetectDarkBackground|Root'` passes.

### Phase 3: Terminal brand forms

- [x] T3.1: [RED] Full mark, compact header and footer tests
**Target:** `src/internal/tui/logo_test.go`, `src/internal/tui/components_test.go`
**Context:** [US-2], [US-3]

**Constraints:**
- Force the Ascii profile when asserting glyphs. Cover: the three exact mark lines of design 3.1; side-by-side layout at width 80 and stacked layout at width 40; version normalized to one leading `v` for inputs `1.0.0` and `v1.0.0`; compact header equals `[✓] specforce v1.0.0`; footer starts with `specforce v`; no output contains `SPECFORCE` or `Specforce` with an uppercase S.

**Acceptance Check:**
`go test ./src/internal/tui -run 'FullMark|CompactHeader|Footer'` exits non-zero with those tests failing or not compiling.

- [x] T3.2: [CODE] Gate mark rendering, terminal width and console header
**Target:** `src/internal/tui/logo.go`, `src/internal/tui/terminal.go`, `src/internal/tui/console.go`
**Context:** [US-2], [US-3]

**Constraints:**
- Remove the braille logo symbols listed in the design inventory entry for `logo.go`; keep `PrintBranding` as the full-mark entry point.

**Acceptance Check:**
`go test ./src/internal/tui -run 'FullMark|CompactHeader|Console'` passes.

- [x] T3.3: [CODE] Compact header and footer on shared TUI listings
**Target:** `src/internal/tui/components.go`, `src/internal/tui/implementation.go`
**Context:** [US-3]

**Acceptance Check:**
`go test ./src/internal/tui` passes (all tui tests green, including `Footer`).

- [x] T3.4: [CODE] Switch CLI report screens to the compact header and update short help
**Target:** `src/internal/cli/spec.go`, `src/internal/cli/constitution.go`, `src/internal/cli/cobra/root.go`
**Context:** [US-3], [US-6]

**Constraints:**
- Only the branding calls listed in design 3.4 change; `HandleInit`, `HandleHelp` and `HandleUnknown` keep the full mark.

**Acceptance Check:**
`go run ./src/cmd/specforce --help | head -1 | grep -q 'Spec-Driven Development for AI agents' && ! grep -n 'PrintBranding' src/internal/cli/spec.go src/internal/cli/constitution.go` succeeds.

### Phase 4: Project pages, metadata and constitution

- [x] T4.1: [DOCS] New brand header in the three READMEs
**Target:** `README.md`, `README.pt.md`, `README.es.md`
**Context:** [US-6]

**Constraints:**
- Replace only the header block and the `# Specforce Kit` title, per the design inventory entry for the READMEs; all other lines stay byte-identical.

**Acceptance Check:**
`for f in README.md README.pt.md README.es.md; do grep -q 'assets/brand/app-icon-256.png' $f && grep -q 'prefers-color-scheme: dark' $f && grep -qx '# Specforce' $f && ! grep -q 'Specforce Kit$' <(sed -n 1,15p $f) || exit 1; done; test "$(git diff --numstat README.md | awk '{print $1+$2}')" -le 16` succeeds.

- [x] T4.2: [DOCS] Refresh package READMEs and package description
**Target:** `npm/*/README.md` (5 platform copies), `package.json`
**Context:** [US-6]

**Constraints:**
- Platform READMEs are exact copies of the new `README.md` (same copy the `Makefile` npm target performs at release).

**Acceptance Check:**
`! grep -rl OpenSpec npm/*/README.md && for d in npm/*/; do cmp -s README.md $d/README.md || exit 1; done && node -e "process.exit(require('./package.json').description==='Specforce — Spec-Driven Development for AI agents'?0:1)"` succeeds.

- [x] T4.3: [DOCS] Brand rules in the UI/UX constitution
**Target:** `.specforce/docs/ui-ux.md`
**Context:** [US-7]

**Constraints:**
- Edits exactly as listed in the design inventory entry for `ui-ux.md`; sections 1–8 keep their order and titles.

**Acceptance Check:**
`f=.specforce/docs/ui-ux.md; grep -q 'Ghost in the Machine' $f && grep -q '#007A4D' $f && grep -q 'NO_COLOR' $f && grep -q 'COLORFGBG' $f && grep -q '\[✓\] specforce' $f && ! grep -qi 'braille' $f && ! grep -q 'Warning/Orange' $f && test $(wc -l < $f) -lt 200` succeeds.

## 3. Black-box Scenarios

- **S1:** Full Gate mark on project initialization
**Context:** [US-2]
**Command:** `d=$(mktemp -d) && cd $d && git init -q && NO_COLOR=1 script -qc "stty cols 100; specforce init claude" /dev/null < /dev/null | head -8`
**Expected:** Output contains the three lines `┌─    ╱─┐  specforce v…`, `│ ╲  ╱  │  Spec-Driven Development for AI agents`, `└─ ╲╱  ─┘`, and no `SPECFORCE`.

- **S2:** Compact header on a report screen
**Context:** [US-3]
**Command:** `NO_COLOR=1 script -qc "stty cols 100; specforce constitution status" /dev/null | head -3`
**Expected:** The first non-empty line is `[✓] specforce v<version>`; no braille characters (`⣿`) appear anywhere in the output.

- **S3:** No branding and no colors when piped
**Context:** [US-3], [US-5]
**Command:** `specforce constitution status | grep -c -e 'specforce v' -e $'\e\['`
**Expected:** Prints `0`.

- **S4:** NO_COLOR removes every color sequence on an interactive terminal
**Context:** [US-5]
**Command:** `NO_COLOR=1 script -qc "specforce constitution status" /dev/null | grep -c $'\e\[[0-9;]*m'`
**Expected:** Prints `0` while the output still contains `[✓] specforce`.

- **S5:** Light palette when the terminal announces a light background
**Context:** [US-4]
**Command:** `COLORFGBG='0;15' COLORTERM=truecolor script -qc "specforce constitution status" /dev/null | grep -o $'\e\[[0-9;]*m' | sort -u`
**Expected:** Includes `38;2;0;121;77` (Mint `#007A4D` as emitted by termenv, which truncates one channel by one) and does not include `38;2;0;250;154` (dark Mint `#00FA9A`).

- **S6:** Dark palette unchanged
**Context:** [US-4]
**Command:** `COLORFGBG='15;0' COLORTERM=truecolor script -qc "specforce constitution status" /dev/null | grep -o $'\e\[[0-9;]*m' | sort -u`
**Expected:** Includes `38;2;0;250;154` and does not include `38;2;0;121;77`.

- **S7:** Asset regeneration is reproducible
**Context:** [US-1]
**Command:** `make brand-assets >/dev/null && git status --porcelain assets/brand`
**Expected:** Empty output on a clean checkout where the assets are committed.

## 4. Pre-emptive Mitigations
- **Risk:** The terminal background query can wait up to 5 s on emulators that never answer, and could steal keystrokes from the console. -> **Mitigation:** Query only when interactive and neither `NO_COLOR` nor `COLORFGBG` decides first; set the result once with `SetHasDarkBackground` so lipgloss never queries again (design 4.2).
- **Risk:** README images point to `main` and will not render on the `1.x` branch before merge. -> **Mitigation:** Expected per the design inventory entry for the READMEs; verify rendering after merge, not during implementation.
- **Risk:** Light primary text set to "terminal default" could render as an empty color. -> **Mitigation:** T2.1 asserts that rendering with the light background and TrueColor profile emits no foreground sequence for primary text.
