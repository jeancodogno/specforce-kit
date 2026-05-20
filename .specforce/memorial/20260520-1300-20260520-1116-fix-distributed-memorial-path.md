---
date: 2026-05-20
scope: 20260520-1116-fix-distributed-memorial-path
author: agent
type: Decision
---

# Registry-Driven Path Overrides

Implemented support for custom artifact paths in the Constitution Registry. This allows artifacts like the Distributed Memorial's ROUTING.md to be correctly tracked outside the default .specforce/docs directory. Also refactored the Memorial Service to use templates from the registry during initialization, eliminating hardcoded Markdown strings in the Go source.
