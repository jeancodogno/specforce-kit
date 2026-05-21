---
slug: 20260520-1824-init-layout-refinement
lens: UI-heavy
---

# Feature: Init Layout Refinement

## 1. Context & Value
The `specforce init` command requires a professional, high-density TUI refinement to transform project bootstrapping into a "surgical" experience. By implementing the "Arsenal" selection and "Infrastructure Pulse" feedback systems, we provide users with immediate visual confidence in the project's structural integrity. This ensures the first touchpoint with Specforce feels autonomous, precise, and aligned with the "Ghost in the Machine" aesthetic.

## 2. Out of Scope (Anti-Goals)
- **Content Generation:** This spec MUST NOT handle the actual LLM-driven generation of documentation or code; it is strictly limited to UI layout and directory bootstrapping.
- **Git Operations:** Initializing git repositories or performing commits is out of scope.
- **Template Authoring:** Creating or editing the underlying agent kits/templates is out of scope.
- **Project Re-initialization:** Logic to handle existing directories or migrations is out of scope.

## 3. Acceptance Criteria (BDD)

### [US-1] The Arsenal: Agent Selection
**User Story:** AS A Lead Solutions Architect, I WANT TO select specialized AI agent kits through a high-density tool interface, SO THAT I can equip my project with the correct governance and execution models.

**Scenarios:**
1. **[Happy Path]** GIVEN a clean terminal session WHEN the user runs `specforce init` THEN "The Arsenal" component displays a list of available agent kits with Mint Green borders and Cyan highlights for the active selection.
2. **[Edge Case]** GIVEN a scenario where no agent kits are discovered in the local environment WHEN the command is executed THEN "The Arsenal" displays a "Depleted" state with Red borders and an instructional message on how to fetch kits.

**UI/UX Specifics:**
- **View/Component:** "The Arsenal" (A 2-column selection grid or dense list).
- **Feedback Logic:** Active selection uses a Cyan background or text highlight; non-selected items remain Silver/Ice.
- **Keybindings:** `up/down` or `j/k` to navigate, `enter` to toggle selection, `tab` to move to the confirmation button.

**Technical Constraints (NFR):**
- **[Performance]:** UI rendering and kit discovery must complete in < 50ms.
- **[Safety & Security]:** Tool selection must be validated against the available local kits before proceeding.
- **[Integrity]:** The selection state must be preserved throughout the session until the handoff.
- **[Observability]:** Selected kit IDs must be logged to the internal session state for auditing.

### [US-2] Infrastructure Pulse: Bootstrap Feedback
**User Story:** AS A Developer, I WANT TO receive real-time, surgical feedback during the directory creation process, SO THAT I can verify the project skeleton is being built correctly.

**Scenarios:**
1. **[Happy Path]** GIVEN a set of selected kits WHEN the user confirms the "Init" action THEN the "Infrastructure Pulse" displays a vertical stream of creation events (e.g., `DIR .specforce/docs [OK]`) using Success Green icons (`◉`).
2. **[Edge Case]** GIVEN a "Permission Denied" error during directory creation WHEN the bootstrap is running THEN the "Infrastructure Pulse" halts the stream at the failure point, highlighting the specific path with a Red border.

**UI/UX Specifics:**
- **View/Component:** "Infrastructure Pulse" (A scrolling log-style feed with a surgical "Scanning" animation prefix).
- **Feedback Logic:** Every successful operation uses the `↳` (Right Arrow) prefix in Secondary Gray and a Success Green `◉`.
- **Keybindings:** No user input during pulse; `ctrl+c` to abort (triggering immediate cleanup).

**Technical Constraints (NFR):**
- **[Performance]:** Pulse animations must run at 60fps; file I/O operations must be non-blocking for the TUI thread.
- **[Safety & Security]:** All created directories and files MUST have 0755 or 0600 permissions as defined by the kit.
- **[Integrity]:** Partial failures must trigger an "Incomplete Infrastructure" state, preventing the ghost handoff.
- **[Observability]:** Each pulse event must be reported to the `core.UI` status reporter.

### [US-3] The Ghost Handoff: Completion
**User Story:** AS A User, I WANT TO see a clear, finalized status screen once the infrastructure is ready, SO THAT I know the local environment is prepared for the LLM content generation phase.

**Scenarios:**
1. **[Happy Path]** GIVEN all directories and core files are created WHEN the process finishes THEN the UI transitions to "The Ghost Handoff" screen, displaying an ASCII Braille logo and the message "Infrastructure Online. Standing by for LLM Handoff."
2. **[Edge Case]** GIVEN a successful bootstrap but with non-critical warnings WHEN the process finishes THEN the completion screen displays a "Degraded Readiness" warning in Warning Yellow (`#FFFFAF`).

**UI/UX Specifics:**
- **View/Component:** "The Ghost Handoff" (Full-screen minimalist summary).
- **Feedback Logic:** Mint Green border for the final branding box; Cyan accent for the "Next Steps" text.
- **Keybindings:** `q` or `enter` to exit the TUI and return to the shell.

**Technical Constraints (NFR):**
- **[Performance]:** Final transition must be instantaneous (< 10ms).
- **[Safety & Security]:** The process must ensure the final state is written to the `config.yaml` before exiting.
- **[Integrity]:** The handoff must only occur if the mandatory `.specforce/` directory structure exists.
- **[Observability]:** Final "Success" signal must be emitted to the CLI exit code (0).

## 4. Business Invariants
- **Atomic Bootstrapping:** A project is either fully initialized or nothing is changed; no "partially initialized" projects should be left behind on fatal errors.
- **Path Isolation:** `specforce init` can only create files within the target project directory or its subdirectories.
- **Agent Presence:** A project cannot be initialized without at least one valid Agent Kit being selected.

## 5. Global UI/UX Contract (TUI Ghost Protocol)
- **Density Posture:** Compact (Optimized for 80x24 terminals).
- **Signature Moves:** Thin Mint Green borders (`+-|`), ASCII-Braille branding, and Cyan highlights.
- **Interaction Model:** Keyboard-only navigation (j/k, arrows, enter, tab, q).
- **State Behavior:**
    - **Loading:** Subtle "scanning" pulses using Silver/Ice.
    - **Error:** Red borders (`#FF5F5F`) with explicit recovery instructions.
    - **Success:** Persistent Green `◉` markers for completed tasks.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Total TUI initialization and transition time between screens must be < 100ms.
- **[Reliability]:** The TUI must handle SIGINT (Ctrl+C) gracefully, ensuring the terminal state is restored.
- **[Security]:** Path sanitization is mandatory; never allow initialization in sensitive system directories (e.g., `/`, `/etc`).
- **[Maintainability]:** Use standard Lipgloss styling functions to ensure theme consistency across all views.
