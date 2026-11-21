---
slug: 20260509-1132-archive-error-prevention
lens: Backend-heavy
---

# Technical Design: Archive Error Prevention

## 1. Architecture Blueprint
*A visual representation of the updated `spf.archive` flow.*

```mermaid
graph TB
    Agent[Principal Architect Agent] -->|Reads| TasksMD(tasks.md)
    Agent -->|Reads| ReqMD(requirements.md)
    Agent -->|Reads| DesignMD(design.md)
    TasksMD -->|Analyzes Challenges| Eval[Constitution Impact Analysis]
    ReqMD -->|Identifies Precedents| Eval
    DesignMD -->|Identifies Precedents| Eval
    Eval -->|Formulates Rule| PromptUser[Information Gathering]
    PromptUser -->|If Approved| Constitution[(Project Constitution)]
    PromptUser -->|If Declined/None| Archive[Archive Execution]
```

## 4. File & Component Inventory
*The exact files that the Developer must create or modify. Map the core responsibility.*

**Backend:**
- `src/internal/agent/kit/instructions/archive.md` -> Update markdown instructions to mandate analyzing challenges from `tasks.md`, checking if errors could be prevented by Constitution updates, and explicitly formulating a prompt that includes these challenges.