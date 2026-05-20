---
date: 2026-05-20
scope: 20260520-1005-layered-instruction-mapping
author: agent
type: Decision
---

# Layered Instruction Injection

Implemented a dual-layer instruction mapping system that merges generic base-type rules (e.g., 'requirements') with specific prefixed rules (e.g., 'feature-requirements') in config.yaml. This uses a Right-to-Left keyword matching heuristic to infer the base type from any artifact name, ensuring consistent instruction delivery even when artifacts are prefixed for specialized workflows.
