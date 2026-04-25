# UI/UX

## 1. Visual Theme & Atmosphere
- **Framework/Base:** Bubbletea (TUI) with Lipgloss for styling.
- **Theme:** Ghost in the Machine.
- **Atmosphere:** Surgical, alien, and extremely clean. Conveys the idea of an autonomous agent working flawlessly.
- **Density Posture:** Compact (Optimized for standard terminal dimensions).
- **Design Philosophy:** Optimizes for readability and fast keyboard navigation; avoids unnecessary visual clutter.
- **Signature Moves:** Consistent "Branding" boxes, ASCII art logos (Braille-based), and color-coded status indicators. Minimalist separators instead of thick borders.

## 2. Color Palette & Roles
### Brand & Accent
- **Primary Brand:** `#00FA9A` (Mint Green) | Primary borders, brand name, and major headings.
- **Secondary Brand:** `#00FFFF` (Cyan) | Subheadings and active selection highlights.
- **Accent:** `#E0FFFF` (Silver/Ice) | Muted subtitles and special focus elements.

### Text & Feedback
- **Primary Text:** `#FFFFFF` | Standard body text and labels.
- **Secondary Text / Metadata:** `#808080` | Muted information, metadata, captions, and pending items.
- **Error:** `#FF5F5F` | Critical failures, validation errors, and destructive actions.
- **Success / Done:** `#5FFF87` | Successful operations and completed items (`◉`).
- **Warning / In-Progress:** `#FFFFAF` | Non-critical alerts and active/working items (`◉`).
- **Info:** `#87AFFF` | Informational messages and help context.

### Iconography & Status States
- **Done:** `◉` (Filled Circle) in **Success Green** (`#5FFF87`).
- **In-Progress:** `◉` (Filled Circle) in **Warning/Orange** (`#FFFFAF`).
- **Pending / Inactive:** `○` (Empty Circle) in **Secondary Gray** (`#808080`).
- **Description / Context:** `↳` (Right Arrow) in **Secondary Gray** (`#808080`) for indented details.

### Surface & Border
- **Canvas:** `#000000` | Terminal background.
- **Surface:** `#121212` | Highlighted block backgrounds.
- **Muted Surface:** `#080808` | Secondary block backgrounds.
- **Border:** `#444444` | Standard UI borders, preferring thin or dashed lines.

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
- Use consistent ASCII art branding across all major views.
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
