# Technical Design: Parallel Task Support

This document outlines the technical architecture and implementation strategy for enabling parallel execution of tasks within the Specforce roadmap.

## 1. Threat Modeling (Security & Integrity)

### 1.1. Race Condition Prevention
**Threat:** Multiple agents or processes attempting to modify the same file simultaneously due to incorrect parallel labeling.
**Countermeasure:** Strict "Target Isolation" validation. The CLI will cross-reference the `**Target:**` field of all tasks in a parallel cluster. If any overlap is detected, the roadmap will be rejected before implementation starts.

### 1.2. Input Validation (ReDoS)
**Threat:** Maliciously crafted `tasks.md` files with complex task IDs or nested metadata.
**Countermeasure:** Use non-backtracking regex patterns for metadata extraction. Limit the number of parallel IDs to a reasonable threshold (e.g., 20 peers per task).

### 1.3. Logic Integrity (Symmetry & Phase Boundaries)
**Threat:** Asymmetric definitions (T1.1 parallel with T1.2, but T1.2 not parallel with T1.1) leading to non-deterministic execution states.
**Countermeasure:** The validation layer will enforce symmetry. If a relationship is declared in one direction, it is automatically assumed for both, and any explicit contradiction or phase boundary violation will trigger a validation error.

---

## 2. Data & Persistence

### 2.1. Go Struct Updates
The `ImplementationTask` struct in `src/internal/spec/implementation.go` will be extended to support parallelism metadata.

```go
// src/internal/spec/implementation.go

type ImplementationTask struct {
    ID           string        `json:"id"`
    Title        string        `json:"title"`
    State        string        `json:"state"`
    Target       string        `json:"target"`
    // ... existing fields ...
    ParallelWith []string      `json:"parallel_with"` // List of peer Task IDs
    IsParallel   bool          `json:"is_parallel"`   // Helper for TUI/Logic
}
```

### 2.2. Parser Logic
A new regex will be introduced to `parseTaskBlock` to extract the `**Parallel With:**` metadata.

```go
// Regex: \*\*Parallel With:\*\* (T[\d.]+(?:,\s*T[\d.]+)*)
```

---

## 3. API Contracts & Interfaces

### 3.1. JSON Schema Update (`specforce spec status --json`)
The `ImplementationReport` will now include parallelism data in the task objects.

```json
{
  "name": "20260601-1553-parallel-tasks-support",
  "phases": [
    {
      "id": "1",
      "tasks": [
        {
          "id": "T1.1",
          "target": "src/auth.go",
          "parallel_with": ["T1.2"],
          "is_parallel": true,
          "state": "FINISHED"
        },
        {
          "id": "T1.2",
          "target": "src/user.go",
          "parallel_with": ["T1.1"],
          "is_parallel": true,
          "state": "PENDING"
        }
      ]
    }
  ]
}
```

---

## 4. Validation Algorithm (Structural Integrity)

The `ValidateTasks` function in `src/internal/spec/tasks.go` will be enhanced with a post-parsing validation pass.

### 4.1. The Isolation Pass
1.  **Build Map:** Create a map of `ID -> TaskInfo`.
2.  **Verify Peer Existence:** For every `T`, ensure all `peerID` in `T.ParallelWith` exist in the same phase.
3.  **Conflict Check:**
    ```mermaid
    graph TD
        A[Start Validation] --> B[Get Task T1]
        B --> C[Get Peer T2]
        C --> D{Same Target?}
        D -- Yes --> E[FAIL: Conflict Error]
        D -- No --> F[Continue]
    ```

### 4.2. Readiness Logic (Implementation Engine)
A task `T` is considered "Ready" if:
- `CurrentPhase(T)` is the first phase with unfinished tasks.
- For all tasks `S` in the same phase:
  - If `Order(S) < Order(T)` AND `S` is NOT in `T.ParallelWith`: `State(S) == FINISHED`.

---

## 5. Surface Blueprint (CLI Error Reporting)

Errors related to parallel tasks must be highly descriptive to allow quick fixes.

```text
[VALIDATION FAILED] tasks.md:
  - T1.1: Parallel conflict with T1.2. Both target "src/internal/spec/tasks.go".
  - T1.3: References non-existent task "T1.99" in Parallel With metadata.
  - T2.1: Parallel With "T1.1" violates phase boundary (must be in same phase).
```

---

## 6. Inventory of Changes

| File | Purpose |
| :--- | :--- |
| `src/internal/spec/implementation.go` | Update `ImplementationTask` struct and `ParseTasks` logic. |
| `src/internal/spec/tasks.go` | Enhance `ValidateTasks` with parallel conflict detection. |
| `src/internal/agent/artifacts/spec/tasks.yaml` | Update agent instructions and template to support parallelism. |
| `src/internal/spec/tasks_test.go` | Add regression tests for parallel conflicts and symmetry. |
| `src/internal/spec/implementation_test.go` | Add tests for parallel metadata parsing. |

---

## 7. Observability & Resilience

### 7.1. Structured Logging
The `ValidateTasks` function will return a slice of strings containing all detected errors. Each error will include:
- The Task ID of the primary offender.
- The specific rule violated (Conflict, Missing, Boundary).
- The conflicting file path (if applicable).

### 7.2. Graceful Degradation
If the parser fails to identify a task ID in the `Parallel With` list, the validation will fail early, preventing the implementation agent from proceeding with a broken roadmap.
