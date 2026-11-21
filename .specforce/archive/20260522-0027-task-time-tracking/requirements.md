---
slug: 20260522-0027-task-time-tracking
lens: Balanced full-stack
---

# Feature: Task Time Tracking

## 1. Context & Value
Development teams often lack visibility into the actual time spent on individual tasks, making it difficult to estimate future work accurately. This feature introduces automated time tracking within the Specforce workflow by capturing timestamps during status transitions. By removing the need for manual logging, we ensure high-fidelity data for productivity analysis and project planning.

## 2. Out of Scope (Anti-Goals)
- Manual time entry or manual adjustment of recorded timestamps.
- Integration with external time-tracking tools (e.g., Jira, Toggl, Clockify).
- Time tracking for Discovery or Planning phases (only Implementation tasks).
- Detailed "Break" or "Pause" button (time is strictly linked to task status).
- User-specific time tracking (Specforce is currently developer-agnostic at the local level).

## 3. Acceptance Criteria (BDD)

### [US-1] Automated Task Time Recording
**User Story:** AS A Developer, I WANT the system to automatically record the time I spend on each task, SO THAT I can track my productivity without manual overhead.

**Scenarios:**
1. **[Happy Path]** GIVEN a task in 'todo' status WHEN I change its status to 'in-progress' THEN the system MUST record the current timestamp as the session start time in the specification data.
2. **[Happy Path]** GIVEN a task in 'in-progress' status WHEN I change its status to 'finished' THEN the system MUST record the end timestamp and calculate the total duration for that task.
3. **[Edge Case - Resuming Task]** GIVEN a task that was previously 'finished' WHEN I move it back to 'in-progress' THEN the system MUST preserve the previous duration and start a new session, appending the additional time upon completion.

**UI/UX Specifics:**
- **View/Component:** Implementation Status TUI list.
- **Feedback Logic:** A "Clock" icon (⏲) appears next to tasks currently in the 'in-progress' state.
- **Keybindings:** Standard status transition keys (e.g., 'space' to toggle status).

**Technical Constraints (NFR):**
- **[Performance]:** Timestamp recording MUST happen in < 10ms to ensure no lag during status transitions.
- **[Safety & Security]:** Timestamps MUST be stored in UTC format to avoid timezone-related calculation errors.
- **[Integrity]:** Duration calculations MUST be performed server-side (within the core service) to ensure consistency across TUI sessions.
- **[Observability]:** Start and Stop events MUST be logged to the internal audit trail.

### [US-2] Aggregate Specification Time Tracking
**User Story:** AS A Project Owner, I WANT TO see the total time invested in a specific feature, SO THAT I can understand the actual cost and complexity of the implementation.

**Scenarios:**
1. **[Happy Path]** GIVEN a specification with multiple completed and in-progress tasks WHEN viewing the Specification Details THEN the system MUST display the "Total Time Invested" as the sum of all recorded task durations.
2. **[Edge Case - Empty Spec]** GIVEN a specification with no tasks defined WHEN viewing its details THEN the "Total Time" field MUST display "0m".

**UI/UX Specifics:**
- **View/Component:** Specification Detail Header / Console Dashboard.
- **Feedback Logic:** Total time is displayed in a "Human Readable" format (e.g., "2h 15m" instead of "135m").
- **Keybindings:** 'Enter' on a spec in the list view to open the detail view.

**Technical Constraints (NFR):**
- **[Performance]:** Summation of task times MUST be performed on-the-fly when the view is rendered in < 50ms.
- **[Safety & Security]:** File-level locking MUST be used when writing to `spec.yaml` to prevent corruption during simultaneous updates.
- **[Integrity]:** The total time MUST always be a derived value from individual task entries, never a standalone editable field.
- **[Observability]:** The sum operation MUST be verifiable through the implementation status logs.

### [US-3] Real-time Progress Visibility
**User Story:** AS A Developer, I WANT TO see the elapsed time for my current task in the TUI, SO THAT I have real-time awareness of my progress.

**Scenarios:**
1. **[Happy Path]** GIVEN an 'in-progress' task WHEN viewing the Implementation Status screen THEN the system MUST show a live-updating counter indicating the time elapsed since the task started.
2. **[Edge Case - Persisted Session]** GIVEN a task was marked 'in-progress' in a previous terminal session WHEN the developer re-opens Specforce THEN the counter MUST correctly calculate and display the total elapsed time from the original start point.

**UI/UX Specifics:**
- **View/Component:** Implementation Status Task Row.
- **Feedback Logic:** 
    - Real-time counter updates every second.
    - **Display Format:** 
        - If Duration < 60m: Show `MMm SSs` (e.g., `12m 45s`).
        - If Duration >= 60m: Show `HHh MMm` (e.g., `01h 15m`).
    - Elapsed time counter is displayed in Mint Green if under 1 hour, and Warning Yellow if the task exceeds 4 hours.
- **Keybindings:** N/A (Passive display).

**Technical Constraints (NFR):**
- **[Performance]:** Counter UI updates MUST happen once per second in the active view to ensure "alive" feeling.
- **[Safety & Security]:** N/A.
- **[Integrity]:** The counter MUST use the `spec.yaml` timestamp as the "Ground Truth" for its calculation.
- **[Observability]:** N/A.

## 4. Business Invariants
- A task cannot have a negative duration.
- Total specification time must always equal the sum of its individual task durations.
- A task can only have one active "in-progress" session at a time.
- Timestamps are immutable once a task is marked as "finished" (unless the task is restarted).

## 5. Global UI/UX Contract (TUI Ghost Protocol)
- **Density Posture:** Standard (80x40).
- **Signature Moves:** Mint Green accents for active time tracking, ASCII-Braille progress bars where appropriate.
- **Interaction Model:** Selection-based status updates; real-time TUI polling for counters.
- **State Behavior:** 'In-Progress' tasks show a blinking duration counter (1Hz) to indicate active measurement.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Latency for status-to-time transitions < 50ms. TUI state refresh at 1Hz (1s) for active views.
- **[Reliability]:** Automated recovery of tracking state if the TUI process is interrupted (Persistence in `spec.yaml`).
- **[Security]:** 0600 file permissions for `spec.yaml` to prevent external modification of time data.
- **[Maintainability]:** Time calculation logic MUST be unit-tested with 100% coverage of edge cases (e.g., leap years, day-light savings transitions).
