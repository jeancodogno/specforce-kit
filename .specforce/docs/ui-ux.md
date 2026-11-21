# UI/UX

## 1. Visual Theme & Atmosphere
- **Framework/Base:** Bubbletea (TUI) with Lipgloss for styling.
- **Theme:** Ghost in the Machine.
- **Atmosphere:** Surgical, alien, and extremely clean. Conveys the idea of an autonomous agent working flawlessly.
- **Density Posture:** Compact (Optimized for standard terminal dimensions).
- **Design Philosophy:** Optimizes for readability and fast keyboard navigation; avoids unnecessary visual clutter.
- **Signature Moves:** The Gate mark (full on welcome surfaces, compact header `[✓] specforce` on working screens), adaptive light/dark colors, and color-coded status indicators. Minimalist separators instead of thick borders.

## 2. Color Palette & Roles
### Brand & Accent
Values are adaptive: Dark is used on dark terminals, Light on light terminals (contrast ratio vs. a white background in parentheses).

| Role | Dark | Light | Use |
|---|---|---|---|
| Primary Brand (Mint) | `#00FA9A` | `#007A4D` (5.4) | Primary borders, brand name, major headings |
| Secondary Brand (Cyan) | `#00FFFF` | `#00727F` (5.7) | Subheadings and active selection highlights |
| Accent (Ice) | `#E0FFFF` | `#3F5F6B` (6.9) | Muted subtitles and special focus elements |
| Success / Done | `#5FFF87` | `#1E7A34` (5.4) | Successful operations and completed items (`◉`) |
| Warning / In-Progress | `#FFFFAF` | `#8A6100` (5.5) | Non-critical alerts and active items (`◉`) |
| Error | `#FF5F5F` | `#C62828` (5.6) | Critical failures, validation errors, destructive actions |
| Info | `#87AFFF` | `#2F5FB3` (6.2) | Informational messages (documented role, not used by code yet) |
| Secondary Text / Metadata | `#808080` | `#5C5C5C` (6.7) | Muted information, captions, pending items |

### Text & Feedback
- **Primary Text:** `#FFFFFF` on dark; the terminal default foreground on light. Standard body text and labels.
- **Badges:** success, warning and error badges use a white label on the light fill (black label on the dark fill). The INFO badge stays neutral grey.

### Iconography & Status States
- **Done:** `◉` (Filled Circle) in **Success Green** (`#5FFF87` dark, `#1E7A34` light).
- **In-Progress:** `◉` (Filled Circle) in **Warning/Yellow** (`#FFFFAF` dark, `#8A6100` light).
- **Pending / Inactive:** `○` (Empty Circle) in **Secondary Gray** (`#808080`).
- **Description / Context:** `↳` (Right Arrow) in **Secondary Gray** (`#808080`) for indented details.
- **Timer Precision (Conditional):** To maximize high-fidelity feedback without overwhelming the interface, TUI timers MUST follow the **1Hz rule**:
    - **Short Durations (< 60m):** Show seconds (`MMm SSs`) and update every 1 second.
    - **Long Durations (>= 60m):** Hide seconds (`HHh MMm`) to maintain high density and focus on high-signal data.

### Surface & Border
- **Canvas:** terminal background (light or dark, detected). Never assume black.
- **Surface:** `#121212` | Highlighted block backgrounds (dark terminals).
- **Muted Surface:** `#080808` | Secondary block backgrounds (dark terminals).
- **Border:** `#444444` dark / `#B0B0B0` light | Standard UI borders, preferring thin or dashed lines.
- **Footer background:** `#111111` dark / `#EEEEEE` light.

### Brand
- **Symbol (Gate):** two brackets `[ ]` frame a check `✓`: a spec passing the gate. Left bracket Mint, check Ice, right bracket Cyan.
- **Full mark** (shown on `specforce init`, 9 columns, 3 lines): `┌─    ╱─┐` / `│ ╲  ╱  │` / `└─ ╲╱  ─┘`, side by side with text at width >= 48, stacked below.
- **Compact header** (report screens and console): `[✓] specforce v<version>`.
- **Wordmark:** lowercase `specforce` in terminal and logo (`spec` primary text, `force` Mint); "Specforce" only in prose.
- **Tagline:** "Spec-Driven Development for AI agents".
- **Assets:** generated SVG/PNG live in `assets/brand/`; regenerate with `make brand-assets`.

### Color environment
- **Dark fallback:** when the background cannot be determined (non-interactive output), dark is assumed.
- **`COLORFGBG`:** when set, its last field decides: 7 or 9-15 is light, any other number is dark; otherwise an interactive terminal is queried once. (The bubbletea framework issues one pre-existing background query on interactive terminals; specforce overrides its result.)
- **`NO_COLOR`:** when set to a non-empty value, all color output is disabled.

## 3. Typography Rules
### Font Family
- **Primary:** Monospace (System Default).
- **Secondary:** N/A.
- **Usage Logic:** Hierarchy is strictly created via color, weight (bold), and whitespace.

### Hierarchy
| Role | Font | Size | Weight | Line Height | Letter Spacing | Notes |
|---|---|---|---|---|---|---|
| Display | Monospace | N/A | Bold | N/A | N/A | Brand/Logo |
| Heading | Monospace | N/A | Bold | N/A | N/A | Primary sections |
| Subheading | Monospace | N/A | Normal | N/A | N/A | Secondary sections |
| Body | Monospace | N/A | Normal | N/A | N/A | Content text |
| UI Label | Monospace | N/A | Bold | N/A | N/A | Input labels |
| Caption | Monospace | N/A | Normal | N/A | N/A | Muted metadata |

### Content Voice
- **Headlines:** Direct and imperative.
- **Buttons & Actions:** Clear and concise (e.g., [ Confirm ], [ Cancel ]).
- **Form Guidance:** Helpful and instructional.

## 4. Component Stylings
### Buttons
- **Primary:** High-contrast text or inverted colors (e.g., `[ Confirm ]`).
- **Secondary:** Standard text within brackets (e.g., `[ Back ]`).
- **Destructive:** Red text or background.

### Cards & Panels
- **Container Style:** Lipgloss `Normal` or custom clean separators (e.g., `-` and `|` with `+` corners). Avoid `Thick` borders.
- **State Behavior:** Borders change color based on focus or status (e.g., Mint Green for active).

### Inputs & Forms
- **Input Baseline:** Clear prefix (e.g., `> ` or `? `).
- **Focus Style:** Cursor highlight or colored prefix.
- **Error Style:** Red text below the input.

### Navigation
- **Primary Navigation:** List-based selection with cursor indicator.
- **Active State:** Cyan or Mint Green highlight.
- **Secondary Navigation:** Footer shortcuts (e.g., `esc: back • q: quit`).

## 5. Layout Principles
- **Spacing Scale:** 1-character/line baseline grid.
- **Grid Logic:** Strictly vertical stacking with 1-line internal padding.
- **Whitespace Philosophy:** Use empty lines to separate logical contexts.
- **Page Structure Bias:** Centered or Left-aligned with a fixed width for core content.
- **Navigation Distribution:** Persistent footer for global actions.

## 6. Depth & Elevation
| Level | Treatment | Use |
|---|---|---|
| Flat | No border | Backgrounds, passive surfaces |
| Raised | Single border | Cards, status blocks |
| Interactive Hover | Highlighted text | Active menu items, buttons |
| Modal / Priority | Clean solid border | Dialogs, progress bars, blocking UI |

- **Shadow Philosophy:** N/A (TUI).
- **Surface Hierarchy:** Higher priority elements use brighter border colors.

## 7. Do's and Don'ts
### Do
- Use the full Gate mark on welcome surfaces and the compact header on working screens.
- Ensure all interactive elements provide visual feedback on focus.
- Adhere to the 80-character width standard for core content.

### Don't
- Use blinking text or excessive flashing colors.
- Depend on mouse-only interactions.
- Overcrowd a single screen; use scrolling or multi-step wizards.

## 8. Responsive Behavior
| Breakpoint | Width | Key Changes |
|---|---|---|
| Mobile | < 40 chars | Minimalist mode; single-column only. |
| Tablet | < 80 chars | Compact mode; standard layout. |
| Desktop | > 80 chars | Expanded mode; sidebars allowed. |
| Large Desktop | > 120 chars | Maximum readability; additional gutters. |

- **Touch Targets:** N/A (Keyboard-centric).
- **Collapsing Strategy:** Hide non-essential metadata on small terminals.
- **Responsive Priority:** Command inputs and core status are prioritized.

## 9. Agent Prompt Guide
### Quick Reference
- **Theme Summary:** Ghost in the Machine (TUI).
- **Primary Colors:** Mint Green (#00FA9A), Cyan (#00FFFF), and Silver/Ice (#E0FFFF).
- **Typography Summary:** Monospace only; hierarchy via weight/color.
- **Key Surface Rules:** Borders define depth; spacing defines context. Prefer thin/clean lines over thick borders.

### Example Prompts
- "Build a TUI component using Lipgloss that represents a project status card with a thin Mint Green border and a green success label."
- "Build a Bubbletea form for gathering user input, ensuring the prompt prefix changes to red on validation failure."
- "When designing this TUI layout, ensure it remains readable on an 80x24 terminal by using vertical stacking and minimal padding."
