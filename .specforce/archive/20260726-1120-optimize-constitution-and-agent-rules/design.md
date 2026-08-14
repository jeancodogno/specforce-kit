---
slug: optimize-constitution-and-agent-rules
lens: Integration
---

# Technical Design: Optimize Constitution and Agent Rules Guidance

## 1. Architecture Blueprint

```mermaid
graph TB
    AgentsMD["src/internal/project/agents_md.go"] -->|Generates AGENTS.md| ProjectRoot["AGENTS.md"]
    DiscoveryCmd["src/internal/agent/kit/commands/discovery.yaml"] -->|Layer 1 Anchor| DirectRead[".specforce/docs/*.md & modules/*.md"]
    SpecCmd["src/internal/agent/kit/commands/spec.yaml"] -->|Layer 1 Anchor| DirectRead
    ConstCmd["src/internal/agent/kit/commands/constitution.yaml"] -->|Status Discovery| ConstCLI["specforce constitution status --json"]
```

## 2. File & Component Inventory

**Backend (Go Template & Generation):**
- `src/internal/project/agents_md.go` -> Updates `agentsMDTemplate` constant under Section 4 to include explicit guidance on `.specforce/docs/modules/<slug>.md`.
- `src/internal/project/agents_md_test.go` -> Adds unit test assertion to verify `AGENTS.md` output contains module guidance.

**Agent Kit Commands (YAML Definitions):**
- `src/internal/agent/kit/commands/discovery.yaml` -> Updates `Layer 1: Constitutional Anchor` instructions to prioritize direct file reading of `.specforce/docs/` and `modules/` instead of executing `specforce constitution status --json`.
- `src/internal/agent/kit/commands/spec.yaml` -> Updates `Layer 1: Constitutional Anchor` instructions to prioritize direct file reading of `.specforce/docs/` and `modules/` instead of executing `specforce constitution status --json`.
