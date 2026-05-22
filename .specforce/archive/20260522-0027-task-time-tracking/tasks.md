# Task Roadmap: Task Time Tracking

## 1. Execution Strategy
We will implement automated time tracking by extending the `spec.yaml` metadata to store session-based timestamps for each task. The core logic will be integrated into the existing task status transition flow, ensuring zero-overhead tracking for the developer. The TUI will then be updated to display these durations and provide real-time feedback.

## 2. Atomic Tasks

### Phase 1: Data Structures & Core Logic (TDD Cycle)
- [x] T1.1: [CODE] Define TimeSession, TaskTimeLog and Metadata extensions.
**Target:** `src/internal/spec/metadata.go`
**Context:** [US-1]
**Action Steps:**
- Define `TimeSession` and `TaskTimeLog` structs.
- Update `Metadata` struct with `TimeLogs` map and YAML tags.
**Verification (TDD):**
`go build ./src/internal/spec/metadata.go`

- [x] T1.2: [TEST] Create failing unit tests for session management (RED).
**Target:** `src/internal/spec/metadata_test.go`
**Context:** [US-1]
**Action Steps:**
- Implement tests for `StartSession`, `EndSession`, and `GetTaskDuration`.
- Assert expected UTC timestamps and duration calculations.
**Verification (TDD):**
- Run `go test ./src/internal/spec/...` and confirm it fails.

- [x] T1.3: [CODE] Implement Metadata methods to pass tests (GREEN).
**Target:** `src/internal/spec/metadata.go`
**Context:** [US-1]
**Action Steps:**
- Implement `StartSession` with UTC logic and session cleanup.
- Implement `EndSession` to close the active session.
- Implement `GetTaskDuration` for cumulative calculation.
**Verification (TDD):**
`go test -v src/internal/spec/metadata_test.go src/internal/spec/metadata.go`

### Phase 2: Persistence Integration (TDD Cycle)
- [x] T2.1: [TEST] Create integration test for automated status-time trigger (RED).
**Target:** `src/internal/spec/tasks_test.go`
**Context:** [US-1]
**Action Steps:**
- Write a test that simulates `updateTaskStatusFile` and checks if `spec.yaml` logs are updated.
**Verification (TDD):**
- Run test and confirm it fails to record time logs.

- [x] T2.2: [CODE] Integrate session tracking into updateTaskStatusFile (GREEN).
**Target:** `src/internal/spec/tasks.go`
**Context:** [US-1]
**Action Steps:**
- Update `updateTaskStatusFile` to call `StartSession`/`EndSession` side-effects.
- Ensure `LoadMetadata` and `SaveMetadata` are correctly utilized.
**Verification (TDD):**
`go test -v src/internal/spec/tasks_test.go`

### Phase 3: Domain & Scanner Updates
- [x] T3.1: Update ImplementationTask and ImplementationReport structs with time fields.
**Target:** `src/internal/spec/implementation.go`
**Context:** [US-2]
**Action Steps:**
- Add `TotalTime time.Duration` and `IsWorking bool` fields to `ImplementationTask`.
- Add `TotalTime time.Duration` field to `ImplementationReport` to store the aggregate feature time.
**Verification (TDD):**
`go build ./src/internal/spec/implementation.go`

- [x] T3.2: Update ParseTasks to cross-reference time logs from Metadata.
**Target:** `src/internal/spec/implementation.go`
**Context:** [US-2]
**Action Steps:**
- Update `ParseTasks` to load metadata using `LoadMetadata`.
- For each task parsed, fetch its cumulative duration and "working" status from metadata.
- Calculate the total report duration by summing all individual task durations.
**Verification (TDD):**
`go test -v src/internal/spec/implementation_test.go`

- [x] T3.3: Update StateItem in scanner.go to include duration fields.
**Target:** `src/internal/spec/scanner.go`
**Context:** [US-3]
**Action Steps:**
- Update `StateItem` struct to include `TotalTime time.Duration` and `ActiveTaskElapsed time.Duration`.
- Update `scanSingleActiveSpec` to populate these fields using the data from the `ImplementationReport`.
**Verification (TDD):**
`go test -v src/internal/spec/scanner_test.go`

### Phase 4: TUI Presentation
- [x] T4.1: Define ClockStyle and OvertimeStyle in the TUI theme.
**Target:** `src/internal/tui/theme.go`
**Context:** [US-3]
**Action Steps:**
- Add `ClockStyle` using the `brandMint` color for active tracking feedback.
- Add `OvertimeStyle` using the `warningYellow` color for tasks exceeding 4 hours.
**Verification (TDD):**
`go build ./src/internal/tui/theme.go`

- [x] T4.2: Display task duration and active counter in Implementation view.
**Target:** `src/internal/tui/implementation.go`
**Context:** [US-3]
**Action Steps:**
- Modify `renderTask` to display the formatted duration next to the task ID.
- **Display Logic:** If duration < 60m, show `MMm SSs`; if duration >= 60m, show `HHh MMm`.
- Display a blinking or Mint Green `⏲` icon for tasks where `IsWorking` is true.
- Apply `OvertimeStyle` to the duration text if it exceeds the 4-hour threshold.
**Verification (TDD):**
`go build ./src/internal/tui/implementation.go`

- [x] T4.3: Display aggregate specification time in Console view.
**Target:** `src/internal/tui/console.go`
**Context:** [US-2]
**Action Steps:**
- Update `renderItem` for `CategoryImplementations` to display the total aggregate time for the spec.
- Update the "Working on" detail line to include the real-time elapsed duration for the current task (using the same conditional seconds logic).
**Verification (TDD):**
`go build ./src/internal/tui/console.go`

- [x] T4.4: Implement 1s tick-based refresh for TUI time counters.
**Target:** `src/internal/tui/console.go`
**Context:** [US-3]
**Action Steps:**
- Update the `tick` function to use `time.Second` to ensure the UI feels alive.
- Ensure the `Update` loop correctly re-scans the project to refresh all duration counters.
**Verification (TDD):**
- Verify in the console that the counter updates every second when a task is in-progress.

## 3. Pre-emptive Mitigations
- **Timezone Inconsistency:** Use `time.Now().UTC()` for all persistence and internal logic to avoid DST and timezone shifts.
- **File Corruption:** Ensure `SaveMetadata` uses `0600` permissions and consider adding a temporary backup during write if corruption is observed.
- **TUI Performance:** The 1s tick rate ensures the UI feels alive, while the display logic (hiding seconds after 60m) keeps the interface clean for long-running tasks.
