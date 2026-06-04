# Module: Constitution

## 1. Domain Invariants
- All constitution documents MUST reside in `.specforce/docs/`.
- Domain-specific module manifests MUST reside in `.specforce/docs/modules/`.
- The CLI MUST NOT inject module content into the global context to prevent context bloat.
- Domain-specific rules SHOULD be moved to module manifests to keep the Global Constitution focused on high-level architecture and security.

## 2. Technical Patterns
- **Lazy Loading**: Agents MUST only read module manifests when domain affinity (path, name, or feature slug) is detected.
- **Affinity-Based Discovery**: The `specforce constitution status --json` command reports only module slugs, not their full content.
- **Module Template**: Module manifests MUST follow the structure defined in `src/internal/agent/artifacts/constitution/module.yaml`.

## 3. Success Metrics (Local)
- Maintain context token overhead below 15% per session even with manifests loaded.
- Reduce agent hallucinations regarding specific business rules in legacy modules.
