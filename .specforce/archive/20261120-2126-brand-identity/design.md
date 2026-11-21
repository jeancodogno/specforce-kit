---
slug: 20261120-2126-brand-identity
lens: UI-heavy
---

# Technical Design: Brand Identity Refresh (Gate symbol, wordmark and adaptive palette)

## 1. Architecture Blueprint

```mermaid
graph TB
    subgraph Asset pipeline (separate Go module, never shipped)
        Def[scripts/brandgen brand definition<br/>geometry + palettes + Go Mono Bold] --> Build[build all outputs in memory]
        Build --> SVG[SVG writer]
        Build --> PNG[x/image/vector rasterizer + PNG encoder]
        SVG --> Atomic[temp file + rename]
        PNG --> Atomic
        Atomic --> Assets[(assets/brand/*.svg, *.png)]
    end
    Assets --> READMEs[README.md / .pt / .es<br/>absolute raw URLs]
    READMEs -->|make npm copy| NpmReadmes[npm/*/README.md]

    subgraph Shipped binary
        Main[cobra Execute] --> Configure[tui.ConfigureColor]
        Configure -->|NO_COLOR non-empty| Ascii[lipgloss.SetColorProfile Ascii]
        Configure -->|COLORFGBG / TTY query / fallback dark| Bg[lipgloss.SetHasDarkBackground]
        Main --> Cmds[commands]
        Cmds -->|init, help, unknown| Full[tui.PrintBranding: full Gate mark]
        Cmds -->|spec/constitution/implementation status, artifact lists, console| Compact[tui.PrintHeader / RenderCompactHeader]
        Full --> Theme[tui theme: lipgloss.AdaptiveColor palette]
        Compact --> Theme
    end
```

## 2. Brand Definition (single source of truth)

All coordinates use a 72×72 unit grid, origin top-left, y down. Every shape is a **filled polygon** (no SVG strokes) so the SVG and PNG outputs share identical geometry.

### 2.1 Palettes
| Role | Dark background | Light background |
|---|---|---|
| Left bracket (Mint) | `#00FA9A` | `#007A4D` |
| Right bracket (Cyan) | `#00FFFF` | `#00727F` |
| Check (Ice) | `#E0FFFF` | `#3F5F6B` |
| Wordmark "spec" | `#FFFFFF` | `#1A1A1A` |
| Wordmark "force" | `#00FA9A` | `#007A4D` |
| App icon tile | `#0B0B0B` (always, with dark-palette symbol) | — |

### 2.2 Gate symbol polygons
- **Left bracket** (one polygon, clockwise): `(6.5,6.5) (25.5,6.5) (25.5,13.5) (13.5,13.5) (13.5,58.5) (25.5,58.5) (25.5,65.5) (6.5,65.5)`.
- **Right bracket**: mirror of the left bracket with `x' = 72 − x`: `(65.5,6.5) (46.5,6.5) (46.5,13.5) (58.5,13.5) (58.5,58.5) (46.5,58.5) (46.5,65.5) (65.5,65.5)`.
- **Check**: centerline `A(24,37) → B(33,46) → C(49,27)`, thickness 7, square caps at `A` and `C`, miter join at `B`. Emitted as ONE polygon (coordinates rounded to 2 decimals, computed as the offset lines at ±3.5 with their intersections at `B`): `(19.05,37.00) (33.22,51.17) (53.93,26.58) (48.58,22.07) (32.78,40.83) (24.00,32.05)`.
  Why: two overlapping square-capped quads leave a V notch at `B`; a single mitered outline does not.

### 2.3 Lockup (symbol + wordmark)
- ViewBox `0 0 312 72`. Symbol at x 0–72 (polygons from 2.2 unchanged).
- Wordmark `specforce` drawn from the glyph outlines of **Go Mono Bold** (`golang.org/x/image/font/gofont/gomonobold`), em size 40 units, first glyph origin at x = 88, baseline y = 50; advance per glyph = font advance at that size. Glyphs 1–4 ("spec") use the "spec" color, glyphs 5–9 ("force") the "force" color.
- Glyph outlines are loaded with `sfnt.LoadGlyph` at ppem 1000 and scaled by `40/1000`; segments map 1:1 to SVG `M/L/Q/C/Z` and to `vector.Rasterizer` `MoveTo/LineTo/QuadTo/CubeTo/ClosePath`.
  Why: outlining the text removes any dependence on fonts installed on the viewer's machine.

### 2.4 App icon
- ViewBox `0 0 72 72`. Rounded square tile `x=0 y=0 w=72 h=72 r=16`, fill `#0B0B0B` (rounded corners approximated by quarter-circle cubic Béziers with k = 0.5523).
- Dark-palette Gate symbol transformed by `scale(0.72)` then `translate(10.08,10.08)` (centered).

### 2.5 Outputs (`assets/brand/`)
| File | Content | Raster sizes |
|---|---|---|
| `gate-dark.svg` | Gate, dark palette, transparent | `gate-dark-{16,32,64,128,256,512}.png` |
| `gate-light.svg` | Gate, light palette, transparent | `gate-light-{16,32,64,128,256,512}.png` |
| `app-icon.svg` | App icon | `app-icon-{16,32,64,128,256,512}.png` |
| `lockup-dark.svg` | Lockup, dark palette, transparent | `lockup-dark-{128,256}.png` (height; width = height × 312/72, rounded) |
| `lockup-light.svg` | Lockup, light palette, transparent | `lockup-light-{128,256}.png` |

Total: 5 SVG + 22 PNG = 27 files.

### 2.6 Generator behavior (`scripts/brandgen`)
- Command: `make brand-assets` → `cd scripts/brandgen && go run . -out ../../assets/brand`.
- Builds every output (bytes) in memory first, in the fixed order of table 2.5; then writes each file as `<name>.tmp` in the same directory followed by `os.Rename`. Prints `wrote <relative path>` per file.
- On any error: prints `brandgen: cannot write <path>: <err>` to stderr, removes its own `.tmp` file, exits 1. Files already renamed stay; files not yet reached are untouched.
- Determinism: coordinates formatted with `strconv.FormatFloat(v, 'f', 2, 64)`; SVG attributes in fixed order; PNG via `png.Encoder{CompressionLevel: png.BestCompression}` on `image.NRGBA`; no timestamps or map iteration.
- Rasterizing: `golang.org/x/image/vector.Rasterizer` sized to the target, all coordinates multiplied by `size/72` (or `height/72` for lockups); each shape drawn with `Draw(dst, bounds, image.NewUniform(color), image.Point{})` in order tile → left bracket → right bracket → check → glyphs.
  Why: pure Go with no SVG parser; the generator owns the geometry, so parsing its own SVG output is unnecessary.
- SVG format: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="…" width="…" height="…">` + one `<path fill="#RRGGBB" d="…"/>` per shape; no `<text>`, no strokes, no gradients.

## 3. Terminal Brand Surfaces

### 3.1 Full Gate mark (welcome surfaces)
Exact glyphs, each line split into left (Mint) | middle (Ice) | right (Cyan), 9 columns wide:

```
line 1:  "┌─" + "    ╱" + "─┐"      →  ┌─    ╱─┐
line 2:  "│ " + "╲  ╱ " + " │"      →  │ ╲  ╱  │
line 3:  "└─" + " ╲╱  " + "─┘"      →  └─ ╲╱  ─┘
```

Wireframe (width ≥ 48 columns: side by side, 2-space gap, 1 blank line above):
```
+------------------------------------------------+
|                                                |
| ┌─    ╱─┐  specforce v1.0.0                    |
| │ ╲  ╱  │  Spec-Driven Development for AI agents|
| └─ ╲╱  ─┘                                      |
+------------------------------------------------+
```
Width < 48 columns: mark first, then the two text lines below it (stacked, no gap column).

- Text line 1: wordmark (`spec` = primary text bold, `force` = Mint bold) + space + version (dimmed), version normalized to exactly one leading `v`.
- Text line 2: tagline `Spec-Driven Development for AI agents` (dimmed). Text line 3: empty.
- Rendered once per process and cached per width class (side-by-side / stacked).

### 3.2 Compact header (working screens)
```
[✓] specforce v1.0.0
```
`[` Mint, `✓` Ice, `]` Cyan, one space, wordmark as in 3.1, one space, dimmed version. `PrintHeader` prints only this line; the blank line after it comes from the callers, which already print `"\n" + title`.

### 3.3 Footer
`specforce v<version> | Workspace: <cwd>` on footer background (adaptive).

### 3.4 Call-site mapping
| Surface | Brand form |
|---|---|
| `HandleInit` (reached via `specforce init`); `HandleHelp`/`HandleUnknown` keep calling `PrintBranding` but have no caller today | Full mark (`PrintBranding`) |
| Spec status, spec artifact, constitution status, constitution artifact/listing, implementation status, artifact lists | Compact header (`PrintHeader`) |
| Console dashboard | Compact header (`RenderCompactHeader`) at top of rendered view |
All remain guarded by `IsTTY()` exactly as today.

## 4. Color Environment

### 4.1 Adaptive palette (`tui/theme.go`)
Every palette variable becomes `lipgloss.AdaptiveColor{Light, Dark}`:

| Variable | Light | Dark |
|---|---|---|
| `brandMint` | `#007A4D` | `#00FA9A` |
| `brandCyan` | `#00727F` | `#00FFFF` |
| `brandIce` | `#3F5F6B` | `#E0FFFF` |
| `successGreen` | `#1E7A34` | `#5FFF87` |
| `warningYellow` | `#8A6100` | `#FFFFAF` |
| `errorRed` | `#C62828` | `#FF5F5F` |
| `mutedGrey` (borders/separators) | `#B0B0B0` | `#444444` |
| `textGrey` (secondary) | `#5C5C5C` | `#808080` |
| `textPrimary` (renamed from `textWhite`) | `""` (terminal default) | `#FFFFFF` |
| `footerBg` (new) | `#EEEEEE` | `#111111` |
| `badgeOnFill` (replaces `black` in success/warning badges) | `#FFFFFF` | `#000000` |

- `ErrorBadgeStyle` foreground stays white in both modes (`lipgloss.Color("#FFFFFF")` wrapped as `AdaptiveColor{"#FFFFFF","#FFFFFF"}`).
- Alias `inProgressOrange` renamed to `inProgressYellow` (same value as `warningYellow`).
- Off-palette literals in upgrade screens are mapped to roles: `#00FFFF` → `brandCyan`, `#FF00FF` → `brandMint`, `#FF0000` → `errorRed`, `#00FF00` → `successGreen`.
- Rule enforced by test: outside `theme.go`, no `lipgloss.Color(` literal remains in non-test files of `src/internal/tui`.

### 4.2 Configuration at start-up (`tui/colorenv.go`)
`ConfigureColor(getenv func(string) string, interactive bool, query func() bool)` runs once in `cobra.Execute` before `rootCmd.ExecuteContext`, called as `tui.ConfigureColor(os.Getenv, tui.IsTTY(), lipgloss.HasDarkBackground)`:
1. If `getenv("NO_COLOR") != ""` → `lipgloss.SetColorProfile(termenv.Ascii)` and return (no background detection; nothing to color).
2. `dark := detectDarkBackground(getenv, interactive, query)` → `lipgloss.SetHasDarkBackground(dark)`.

`detectDarkBackground` precedence:
1. `COLORFGBG` set: take the last `;`-separated field as an integer `bg`. `bg == 7` or `9 ≤ bg ≤ 15` → light (false); any other integer → dark (true). Unparseable → continue to step 2.
2. `!interactive` → dark (true), no query.
3. Otherwise `query()` (termenv OSC 11 query through lipgloss's default renderer; termenv caps it at 5 s and returns dark on failure).
   Why: `SetHasDarkBackground` marks the background as explicit, so lipgloss never queries again later in the process.
   Note: bubbletea v1's package `init()` already calls `lipgloss.HasDarkBackground()` once on interactive terminals before `main`; that pre-existing query is out of scope and its result is overridden here.

### 4.3 Terminal width (`tui/terminal.go`)
`GetTerminalWidth()` returns the width from `github.com/charmbracelet/x/term.GetSize(os.Stdout.Fd())`; on error or width ≤ 0 returns 80. (`charmbracelet/x/term` is already in the module graph; it moves from indirect to direct.)

## 5. File & Component Inventory

**Asset pipeline (new, separate module):**
- `scripts/brandgen/go.mod` -> module `github.com/jeancodogno/specforce-kit/scripts/brandgen`, Go 1.26 (required by `golang.org/x/image` v0.48), requires `golang.org/x/image`. Not part of the root module, so root `go test ./...` and the shipped binary never include it.
- `scripts/brandgen/brand.go` -> palettes (2.1), Gate polygons (2.2), check quad construction, app-icon tile and transform (2.4), output table (2.5).
- `scripts/brandgen/wordmark.go` -> Go Mono Bold outline loading and scaling (2.3).
- `scripts/brandgen/svg.go` -> deterministic SVG serialization of shapes.
- `scripts/brandgen/raster.go` -> rasterization and PNG encoding.
- `scripts/brandgen/main.go` -> flag `-out`, in-memory build, atomic writes, logging and exit codes (2.6).
- `scripts/brandgen/brandgen_test.go` -> determinism, file list, 16 px legibility, write-failure behavior.
- `Makefile` -> new `brand-assets` target (2.6).
- `assets/brand/*` -> 27 generated files (2.5), committed.
- `assets/logo.png`, `assets/logo-header.png` -> deleted.

**Terminal layer:**
- `src/internal/tui/theme.go` -> palette vars and styles per 4.1; `black` and `textWhite` removed.
- `src/internal/tui/colorenv.go` (new) -> `ConfigureColor`, `detectDarkBackground`.
- `src/internal/tui/logo.go` -> remove `brailleLines`, `gradientColors`, `GenerateLogo`, `PrintLogo`; add `Tagline` const, `formatVersion`, `renderWordmark`, `renderFullMark(version string, width int)`, `RenderCompactHeader(version string)`, `PrintHeader()`; `PrintBranding()` prints `renderFullMark(core.Version, GetTerminalWidth())`. Caches: `cachedLogo`/`cachedBranding` replaced by `cachedHeader` and `cachedFullMark` keyed by layout (side-by-side vs stacked).
- `src/internal/tui/terminal.go` -> `GetTerminalWidth` per 4.3.
- `src/internal/tui/components.go` -> `FooterStyle` background `footerBg`; `PrintFooter` text per 3.3 using `formatVersion`; `PrintArtifactList` → `PrintHeader()`; `RenderBadge` default style uses `textPrimary`/`mutedGrey`.
- `src/internal/tui/console.go` -> `renderLogo` returns `RenderCompactHeader(core.Version) + "\n"`.
- `src/internal/tui/implementation.go` -> `RenderImplementationStatus` → `PrintHeader()`.
- `src/internal/tui/upgrade_notification.go`, `upgrade_prompt.go`, `upgrade_progress.go` -> literals mapped to roles per 4.1.
- `src/internal/cli/spec.go` -> `renderSpecStatusTUI`, `renderSpecArtifactTUI` → `tui.PrintHeader()`.
- `src/internal/cli/constitution.go` -> `HandleConstitutionStatus`, `HandleConstitutionArtifact` TUI branches → `tui.PrintHeader()`.
- `src/internal/cli/cobra/root.go` -> `Execute` calls `tui.ConfigureColor` first; `rootCmd.Short` = `tui.Tagline`; `rootCmd.Long` is prefixed with `tui.Tagline` + blank line (cobra `--help` prints `Long`).
- Tests (new/extended): `src/internal/tui/logo_test.go`, `colorenv_test.go`, `theme_test.go` (adaptive completeness, light contrast ≥ 4.5:1 computed with the WCAG relative-luminance formula, no stray literals), `components_test.go` (footer text). Tests force `lipgloss.SetColorProfile(termenv.Ascii)` when asserting plain glyphs.

**Docs & metadata:**
- `README.md`, `README.pt.md`, `README.es.md` -> header block replaced by:
  ```html
  <p align="center">
    <a href="https://github.com/jeancodogno/specforce-kit/">
      <picture>
        <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/jeancodogno/specforce-kit/main/assets/brand/lockup-dark-256.png">
        <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/jeancodogno/specforce-kit/main/assets/brand/lockup-light-256.png">
        <img src="https://raw.githubusercontent.com/jeancodogno/specforce-kit/main/assets/brand/app-icon-256.png" alt="{localized alt}" height="128">
      </picture>
    </a>
  </p>
  <p align="center">{localized tagline}</p>
  ```
  Localized alt: en `Specforce logo`, pt `Logo do Specforce`, es `Logotipo de Specforce`. Localized tagline: en `Spec-Driven Development for AI agents.`, pt `Spec-Driven Development para agentes de IA.`, es `Spec-Driven Development para agentes de IA.`. Title `# Specforce Kit` → `# Specforce`. Nothing else in the files changes.
  Why: absolute URLs keep images working on the npm registry, which copies the README without the `assets/` folder; the images resolve once this change reaches `main`.
- `npm/{darwin-arm64,darwin-x64,linux-arm64,linux-x64,win32-x64}/README.md` -> overwritten with the new `README.md` (same copy the `Makefile` npm target performs).
- `package.json` -> `description`: `Specforce — Spec-Driven Development for AI agents`.
- `.specforce/docs/ui-ux.md` -> Section 1: theme stays "Ghost in the Machine"; Signature Moves replace "ASCII art logos (Braille-based)" with the Gate symbol and the two brand forms. Section 2: palette table gains a Light column (values of 4.1 with contrast ratios: Mint 5.4, Cyan 5.7, Ice 6.9, Success 5.4, Warning 5.5, Error 5.6, Info `#2F5FB3` 6.2, Secondary 6.7), Canvas line becomes "terminal background (light or dark, detected)", "Warning/Orange" becomes "Warning/Yellow". New subsection "Brand" (symbol meaning and drawing from 3.1/3.2, lowercase wordmark rule with "Specforce" only in prose, tagline, asset folder and `make brand-assets`). New subsection "Color environment" (dark fallback, `COLORFGBG`, `NO_COLOR`). Section 7 "Do": "consistent ASCII art branding" → "full Gate mark on welcome surfaces, compact header on working screens".
