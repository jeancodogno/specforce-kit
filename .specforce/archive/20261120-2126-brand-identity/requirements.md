---
slug: 20261120-2126-brand-identity
lens: UI-heavy
---

# Feature: Brand Identity Refresh (Gate symbol, wordmark and adaptive palette)

## 1. Context & Value
The current brand mark (a braille "droplet") does not communicate what the product does, degrades at small sizes, and the terminal palette is nearly invisible on light-background terminals (contrast between 1.04:1 and 1.39:1 on white). The product is also presented as "Specforce Kit", which reads as a derivative of a better-known spec toolkit, and the published package pages still credit the logo to an unrelated project ("OpenSpec"). This feature establishes one coherent brand identity — the "Gate" symbol (a check mark enclosed by square brackets, meaning "spec verified, work allowed through"), a lowercase `specforce` wordmark, the tagline "Spec-Driven Development for AI agents", and a palette readable on any terminal — and applies it to every brand touchpoint.

## 2. Out of Scope (Anti-Goals)
- Do NOT rename the product, the command, the package, the repository, the skill prefix (`spf.`) or the project folder (`.specforce/`). The name stays "Specforce"; "specforce-kit" remains only the repository and package identifier.
- Do NOT redesign console layout, status glyphs (`◉`/`○`), section headings, emoji usage, column alignment or report wording; that belongs to the separate terminal-redesign specification. The only console change here is replacing its brand header.
- Do NOT rewrite README or `docs/` content beyond the brand header block, the product title and logo alternative text; README restructuring, screenshots/GIFs, the duplicated license sentence, the duplicated section numbering and a documentation site belong to the separate documentation specification.
- Do NOT change the dark-terminal colors; they stay exactly as today.
- Do NOT add a user-facing setting or flag to force light or dark mode.
- Do NOT add branding to the help output or to unknown-command errors (both are produced by the command framework today); do NOT remove or rewire the currently unused help/unknown handlers.
- Do NOT add an ASCII-only fallback for terminals without Unicode support.
- Do NOT create social-media, favicon or website assets beyond the files listed in US-1.

## 3. Acceptance Criteria (BDD)

### [US-1] Official brand asset set
**User Story:** AS A project maintainer, I WANT a single official set of logo files (symbol and wordmark, for light and dark backgrounds, in scalable and raster formats), SO THAT every channel shows the same brand and I can regenerate raster files whenever the scalable masters change.

**Scenarios:**
1. **[Happy Path]** GIVEN the repository WHEN a maintainer opens the brand assets folder THEN it contains scalable masters for: the Gate symbol on dark backgrounds, the Gate symbol on light backgrounds, the horizontal lockup (symbol + `specforce` wordmark) on dark backgrounds, the horizontal lockup on light backgrounds, and an app icon (symbol on a solid dark rounded tile) readable on any background.
2. **[Happy Path]** GIVEN the single brand definition (geometry, colors and wordmark typeface) kept in the repository WHEN the maintainer runs the single documented regeneration command THEN every scalable and raster file is (re)created from it, at the sizes listed in the Technical Constraints, byte-for-byte identical across repeated runs.
3. **[Edge Case]** GIVEN a file that cannot be written during regeneration WHEN the command runs THEN it stops with a message naming the offending file, exits unsuccessfully, and every previously generated file that was not yet replaced is left untouched (no half-written files).
4. **[Edge Case]** GIVEN the app icon rendered at 16×16 pixels WHEN viewed THEN the two brackets and the check mark remain visually distinct (no strokes merge).

**UI/UX Specifics:**
- **View/Component:** Gate symbol = left bracket in Mint (`#00FA9A` dark / `#007A4D` light), right bracket in Cyan (`#00FFFF` dark / `#00727F` light), check mark in Ice (`#E0FFFF` dark / `#3F5F6B` light). Wordmark = `specforce` in lowercase monospace, "spec" in the primary text color and "force" in Mint.
- **Feedback Logic:** The legacy braille raster logos are removed so no outdated mark can be referenced.
- **Keybindings:** Not interactive.

**Technical Constraints (NFR):**
- **[Performance]:** Regenerating the full raster set completes in < 10 seconds on a standard developer laptop.
- **[Integrity]:** Raster sizes: symbols and app icon at 16, 32, 64, 128, 256 and 512 px square; lockups at 128 px and 256 px height. Square corners and flat colors only (no gradients, shadows or glow).
- **[Safety & Security]:** The regeneration tooling MUST NOT add any dependency to the shipped command-line binary.
- **[Observability]:** The regeneration command prints one line per file written.

### [US-2] Full brand mark on the welcome surfaces
**User Story:** AS A developer running specforce in a terminal, I WANT to see the full Gate mark with the wordmark, version and tagline when I initialize a project, SO THAT the first impression is recognizable and tells me what the tool is.

**Scenarios:**
1. **[Happy Path]** GIVEN an interactive terminal WHEN the developer runs project initialization (or any other surface that today prints the full branding block) THEN a 3-line Gate mark is shown to the left of three text lines: `specforce v<version>` (wordmark colored as in US-1), the tagline `Spec-Driven Development for AI agents`, and an empty line kept for alignment.
2. **[Happy Path]** GIVEN a version string with or without a leading "v" WHEN the mark renders THEN the version is always shown with exactly one leading "v".
3. **[Edge Case]** GIVEN output redirected to a file or pipe (not an interactive terminal) WHEN the same commands run THEN no brand mark is printed, as today.
4. **[Edge Case]** GIVEN a terminal narrower than the mark plus text WHEN the mark renders THEN the text block is printed below the mark instead of being cut or wrapped mid-word.

**UI/UX Specifics:**
- **View/Component:** Full Brand Mark (exact drawing defined in design). Total width ≤ 50 columns, height exactly 3 lines plus 1 blank line of top padding.
- **Feedback Logic:** The wordmark is never written in uppercase anywhere in terminal output.
- **Keybindings:** Not interactive.

**Technical Constraints (NFR):**
- **[Performance]:** Rendering the mark adds < 5 ms to command start-up and is computed at most once per process.
- **[Integrity]:** Uses only Unicode box-drawing characters and the check mark already supported by terminals that render the current braille logo.

### [US-3] Compact brand header on working screens
**User Story:** AS A developer checking status reports or the console, I WANT a single-line brand header instead of a large logo, SO THAT the information I asked for is visible without scrolling.

**Scenarios:**
1. **[Happy Path]** GIVEN an interactive terminal WHEN the developer opens the console, or prints a specification status, a specification artifact, the constitution status or listing, or the implementation status THEN exactly one header line is printed: `[✓] specforce v<version>` with the symbol colored as in US-1 and the version dimmed.
2. **[Happy Path]** GIVEN the information footer shown after listings WHEN it renders THEN it reads `specforce v<version> | Workspace: <path>` (lowercase wordmark).
3. **[Edge Case]** GIVEN output redirected to a file or pipe WHEN those commands run THEN no header is printed and the machine-readable output is unchanged.
4. **[Edge Case]** GIVEN the console refreshes its content periodically WHEN a refresh happens THEN the header does not flicker or duplicate.

**UI/UX Specifics:**
- **View/Component:** Compact Brand Header, followed by one blank line.
- **Feedback Logic:** Identical header on every listed screen.
- **Keybindings:** Console keybindings unchanged.

**Technical Constraints (NFR):**
- **[Performance]:** Header rendering < 1 ms; computed at most once per process.
- **[Integrity]:** Header width ≤ 34 columns for versions up to 20 characters (`[✓] specforce ` takes 14).

### [US-4] Palette readable on light and dark terminals
**User Story:** AS A developer whose terminal uses a light background, I WANT every brand and status color to be readable, SO THAT I can use specforce without changing my terminal theme.

**Scenarios:**
1. **[Happy Path]** GIVEN a terminal reporting a dark background WHEN any specforce screen renders THEN colors are exactly the current dark palette.
2. **[Happy Path]** GIVEN a terminal reporting a light background WHEN any specforce screen renders THEN every colored text uses the light palette: Mint `#007A4D`, Cyan `#00727F`, Ice `#3F5F6B`, Success `#1E7A34`, Warning `#8A6100`, Error `#C62828`, Secondary text `#5C5C5C`; primary text uses the terminal's own default foreground.
3. **[Happy Path]** GIVEN a light background WHEN a status badge (success, warning, error) renders THEN its label text is white on the light-palette fill so the badge itself stays readable; the neutral INFO badge keeps the terminal's default text on the light border grey.
4. **[Edge Case]** GIVEN a terminal whose background cannot be determined WHEN a screen renders THEN the dark palette is used.

**UI/UX Specifics:**
- **View/Component:** All styled text, borders, badges, separators and footer across every terminal screen.
- **Feedback Logic:** Borders and separators use `#444444` on dark and `#B0B0B0` on light; the footer bar background is `#111111` on dark and `#EEEEEE` on light.
- **Keybindings:** Not applicable to this requirement (no new keys).

**Technical Constraints (NFR):**
- **[Performance]:** Background detection happens at most once per process, only when output is an interactive terminal, and adds < 50 ms on terminals that answer the background query; on a terminal that never answers, it gives up after at most 5 seconds and falls back to dark. When the terminal announces its colors through its environment, specforce sends no query of its own (the single query issued at start-up by the console framework is pre-existing and out of scope; its result is overridden).
- **[Accessibility]:** Every light-palette text color has a contrast ratio ≥ 4.5:1 against white (`#FFFFFF`).
- **[Integrity]:** No style in the terminal layer may reference a single fixed color for text; every text color has both a light and a dark value.

### [US-5] Colorless output
**User Story:** AS A developer who disabled colors (the `NO_COLOR` convention) or whose output is captured by tools, I WANT specforce to emit no color codes while keeping the same layout, SO THAT logs are clean and the output still reads correctly.

**Scenarios:**
1. **[Happy Path]** GIVEN the `NO_COLOR` environment variable set to any non-empty value AND an interactive terminal WHEN any command renders THEN the output contains zero color escape sequences, and the brand mark/header keep the same shape and characters.
2. **[Edge Case]** GIVEN `NO_COLOR` set to an empty string WHEN a command renders THEN colors are used normally.
3. **[Edge Case]** GIVEN `NO_COLOR` set WHEN a status is shown THEN success, warning and error remain distinguishable by their text labels (`SUCCESS`, `WARNING`, `ERROR`, `INFO`).

**UI/UX Specifics:**
- **View/Component:** Every terminal screen.
- **Feedback Logic:** Bold and underline attributes may remain; only color is removed.
- **Keybindings:** Not interactive.

**Technical Constraints (NFR):**
- **[Performance]:** No measurable start-up overhead (< 1 ms).
- **[Integrity]:** The decision is made once at process start, before any output.

### [US-6] Brand on the project and package pages
**User Story:** AS A visitor of the project's repository page or package page, I WANT to see the new logo correctly in my page theme, the product named "Specforce" and an accurate description, SO THAT I recognize the product and trust that the page is its official home.

**Scenarios:**
1. **[Happy Path]** GIVEN the English, Portuguese and Spanish READMEs WHEN rendered on the repository page with a dark theme THEN the dark lockup is shown; WITH a light theme THEN the light lockup is shown.
2. **[Happy Path]** GIVEN any of those READMEs WHEN rendered by a viewer that ignores theme-specific images (such as the package registry) THEN the app icon is shown, and it is readable on both light and dark pages.
3. **[Happy Path]** GIVEN any README WHEN read THEN the main title is "Specforce" (without "Kit"), the line under the logo is the tagline "Spec-Driven Development for AI agents" (translated in the Portuguese and Spanish READMEs), and the logo alternative text names Specforce in the README's language.
4. **[Happy Path]** GIVEN the published package copies of the README WHEN inspected THEN none of them mention "OpenSpec", and all show the same header as the English README.
5. **[Happy Path]** GIVEN the package metadata description and the command's short help description WHEN read THEN both read "Specforce — Spec-Driven Development for AI agents" (the short help may omit the "Specforce — " prefix).
6. **[Edge Case]** GIVEN a viewer that renders the README outside the repository (the package registry) WHEN it resolves logo images THEN the images still load, because they are referenced by absolute address to the repository's main branch.

**UI/UX Specifics:**
- **View/Component:** README header block (centered logo, tagline line, title).
- **Feedback Logic:** Logo displayed at 128 px height.
- **Keybindings:** Not interactive.

**Technical Constraints (NFR):**
- **[Performance]:** Each raster image referenced by READMEs is ≤ 100 KB.
- **[Integrity]:** Only the header block, title and descriptions change; all other README content is byte-for-byte preserved.

### [US-7] Brand rules in the project constitution
**User Story:** AS A contributor or AI agent working on the project, I WANT the UI/UX constitution to describe the Gate symbol, wordmark, tagline, adaptive palette and colorless behavior, SO THAT future screens follow the brand without guessing.

**Scenarios:**
1. **[Happy Path]** GIVEN the UI/UX constitution WHEN read THEN it keeps the "Ghost in the Machine" theme and adds: the Gate symbol meaning and drawing, the two terminal brand forms (full mark on welcome surfaces, compact header on working screens), lowercase wordmark rule (capitalized "Specforce" only in prose), the tagline, the light and dark palette table with contrast values, the dark-fallback rule, and the `NO_COLOR` rule.
2. **[Edge Case]** GIVEN statements in the constitution that contradict the new identity (braille ASCII-art logo as signature move, `#000000` canvas assumption, "Warning/Orange" naming of a yellow color) WHEN the update is applied THEN those statements are corrected, not left alongside the new ones.

**UI/UX Specifics:**
- **View/Component:** Constitution document for UI/UX.
- **Feedback Logic:** Not applicable — documentation only.
- **Keybindings:** Not interactive.

**Technical Constraints (NFR):**
- **[Performance]:** Not runtime-affecting; document stays under 200 lines.
- **[Integrity]:** Section structure of the constitution document is preserved (sections 1–8).

## 4. Business Invariants
- The product name shown to users is "Specforce" in prose and `specforce` in the wordmark and terminal; "SPECFORCE", "SpecForce" and "Specforce Kit" never appear in user-facing output or in the README titles.
- The Gate symbol is always drawn as left bracket, check mark, right bracket, in that order, with brackets never closed into a box.
- Dark-background colors are never altered by this feature.
- Non-interactive output never contains brand marks, headers or color codes added by this feature.

## 5. Global UI/UX Contract (TUI Ghost Protocol)
- **Density Posture:** Compact (80×24); the full mark is 3 lines, the compact header 1 line.
- **Signature Moves:** Gate symbol `[✓]`, lowercase `specforce` wordmark with "force" in Mint, thin box-drawing strokes, flat colors.
- **Interaction Model:** Unchanged; no new keys.
- **State Behavior:** Unchanged except brand header; loading, empty and error states keep current behavior but use the adaptive palette.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Total added start-up latency < 50 ms on responsive terminals, < 5 ms when output is not interactive (no detection is attempted).
- **[Reliability]:** Background detection failure or timeout silently falls back to the dark palette; it never errors out and never waits longer than 5 seconds.
- **[Security]:** No new network access; asset tooling reads only repository files and writes only inside the brand assets folder.
- **[Maintainability]:** Terminal-layer coverage for new brand and palette code ≥ 80%; standard Go formatting and linting pass; functions ≤ 50 lines and cyclomatic complexity ≤ 10.
