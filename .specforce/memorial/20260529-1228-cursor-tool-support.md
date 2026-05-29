---
date: 2026-05-29
scope: cursor-tool-support
author: agent
type: Decision
---

# Native Cursor Tool Support with Structured Directories

Implemented native support for the Cursor AI editor. Unlike other tools that require symlinks, Cursor natively reads AGENTS.md from the root. The integration maps Specforce artifacts to a structured .cursor/ directory using standard plural categories: agents/, commands/, and skills/ (the latter using dynamic subdirectories). Introduced the 'UseSubdir' flag in MappingConfig to support this hierarchy without hardcoding tool-specific logic in the core translator.
