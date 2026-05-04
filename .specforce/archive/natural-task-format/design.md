---
slug: natural-task-format
lens: Backend-heavy
---

# Technical Design: Natural LLM Task Format

## 1. Architecture Blueprint

```mermaid
graph TB
    Input[tasks.md Content] --> Parser[findTaskBlock / ParseTasks]
    Parser --> MatchHeader{Match `####`?}
    MatchHeader -- Yes --> ExtractBlock1[Extract Block]
    MatchHeader -- No --> MatchList{Match `- [ ]`?}
    MatchList -- Yes --> ExtractBlock2[Extract Block]
    
    ExtractBlock1 --> StateUpdater[updateTaskStatusFile]
    ExtractBlock2 --> StateUpdater
    
    StateUpdater --> DetectFormat{Has `**State:**`?}
    DetectFormat -- Yes --> UpdateTag[Update Tag]
    DetectFormat -- No --> UpdateCheckbox[Update Checkbox]
    
    ExtractBlock1 --> ConsoleScanner[Scanner / ParseTasks]
    ExtractBlock2 --> ConsoleScanner
    ConsoleScanner --> MapState{Map Checkbox/Tag to State}
    MapState --> UI[Console Progress Bar]
```

## 2. API & Interfaces (The Contract)
- **Internal Contract:** `findTaskBlock(content, taskID string) (int, int, error)` in `tasks.go` must support hybrid header/list starts.
- **Internal Contract:** `ParseTasks(ctx, projectRoot, slug)` in `implementation.go` must use a hybrid regex for task headers to ensure the `StateItem` in `scanner.go` receives correct counts.
- **State Mapping:**
    - `[x]` -> `FINISHED`
    - `[/]` -> `IN-PROGRESS`
    - `[ ]` -> `PENDING`

## 3. File & Component Inventory

**Backend:**
- `[src/internal/spec/tasks.go]` -> Hybrid parsing for block identification and state updates.
- `[src/internal/spec/implementation.go]` -> Hybrid parsing for implementation reports (Scanner/Console).
- `[src/internal/spec/tasks_test.go]` -> Unit tests for `tasks.go`.
- `[src/internal/spec/implementation_test.go]` -> Unit tests for `implementation.go` hybrid parsing.
- `[src/internal/agent/artifacts/spec/tasks.yaml]` -> New `- [ ]` default template and simplified instructions.
- `[.specforce/docs/memorial.md]` -> Record the decision to favor natural formats.
