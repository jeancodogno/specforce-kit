# Proposal: Module Manifests Support

## 1. Context & Motivation
As the project grows, global architectural rules (the Constitution) are no longer sufficient to capture the nuances of specific business domains. Currently, agents must "rediscover" business rules and module-specific patterns by reading code every time. 

The goal is to provide a place for **Module Manifests** that capture domain-specific invariants, business logic, and architectural patterns without bloating the global Constitution or the agent's context window.

## 2. Proposed Solution: Domain-Specific Manifests
We will introduce a new tier of documentation located in `.specforce/docs/modules/*.md`.

### 2.1 Artifact Template
Create a `module.yaml` in `src/internal/agent/artifacts/constitution/` to serve as a template. It will define sections like:
- **Core Intent:** What this module solves.
- **Business Invariants:** Critical rules (e.g., "Always validate inventory before checkout").
- **Module Patterns:** Domain-specific technical choices (e.g., "Uses Redis for session caching").
- **Integration Points:** How it interacts with other modules.

### 2.2 CLI & Core Integration
- **`constitution status`**: Updated to scan the `modules/` directory and include a list of detected module slugs in the JSON output.
- **`archive instructions`**: Updated to include a "Module Behavior Distillation" step, prompting the agent to harvest module-specific lessons into the corresponding manifest.

### 2.3 Agent Workflow Integration
- **Discovery (`spf.discovery`)**: During the "Constitutional Anchor" phase, the agent will check the list of available modules. If a domain match is found (e.g., working on `billing` and `billing.md` exists), it will perform a surgical read of that manifest.
- **Planning (`spf.spec`)**: Similar to Discovery, the Orchestrator will use the manifest to ground the requirements and design in existing domain rules.

## 3. Context Optimization (The "Lazy Load" Strategy)
To prevent overwhelming the LLM with unnecessary data:
1. **Discovery only**: The CLI `status` command only returns the *names* of available manifests.
2. **Surgical Retrieval**: The agent is instructed to read a manifest ONLY if it is relevant to the current task.
3. **High Density**: Manifests will be encouraged to use concise bullet points rather than verbose documentation.

## 4. Non-Blocking Nature (Legacy Support)
Manifests are treated as "Reference Guides" rather than "Source of Truth" for legacy code. If the code contradicts the manifest, the agent is instructed to flag the drift during the `archive` retrospective to ensure the manifest eventually reflects reality.

## 5. Implementation Roadmap
1. **Infra**: Add `module.yaml` template and update `status.go` / `archive.go` logic.
2. **Guidelines**: Update `archive.md`, `discovery.yaml`, and `spec.yaml` instructions.
3. **Bootstrapping**: Create the `.specforce/docs/modules/` directory with a `.gitkeep`.
