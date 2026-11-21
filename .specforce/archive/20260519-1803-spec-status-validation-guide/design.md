---
slug: 20260519-1803-spec-status-validation-guide
lens: Backend-heavy
---

# Technical Design: Validation Guide Golden Model

## 1. Architecture Blueprint
*A visual representation of the data flow or entity relationships for this feature.*

```mermaid
graph TB
    CLI[spec status command] --> Service[SpecService.GetStatus]
    Service --> Scanner[processArtifactStatus]
    Scanner --> Validator[ValidateTasks]
    Validator -- Errors > 0 --> Inject[Inject Golden Model Guide]
    Inject --> JSON[JSON Output]
```

## 4. File & Component Inventory
*The exact files that the Developer must create or modify. Map the core responsibility.*

**Backend:**
- `[src/internal/spec/status.go]` -> Modify `ArtifactStatus` struct to add `ValidationGuide string json:"validation_guide,omitempty"`. Update `processArtifactStatus` to populate this field when `ValidateTasks` returns errors for `tasks.md`.