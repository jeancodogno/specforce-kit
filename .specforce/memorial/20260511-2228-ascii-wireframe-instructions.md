---
date: 20260511-2228
scope: 20260511-2209-ascii-wireframe-instructions
author: gemini-cli
type: Decision
---

# Decision: ASCII Wireframe Standard for Multi-Platform UI Design

## Context
Specforce agents previously lacked a formal mandate to provide visual representations of user interfaces during the technical design phase. While the project follows a "Ghost Protocol" (TUI-first) aesthetic, the system must also support Web and Mobile interfaces.

## Decision
We established ASCII wireframing as the universal standard for UI visualization in all technical blueprints (`design.md`). 

### Rules:
- **Mandatory for UI-Heavy Features:** Any feature with a visible interface must include an ASCII surface blueprint.
- **Universal Scope:** Applied to TUI, Web, and Mobile layouts.
- **Ghost Protocol Style:** Thin borders (`+`, `-`, `|`), 80-character width, and semantic markers (`[ ]`, `( )`, `↳`).

## Lessons Learned
- **Portability:** ASCII wireframes in Markdown ensure that UI designs are readable across all developer environments (CLI, Web, Mobile) without external graphical tools.
- **Agent Intelligence:** Explicit instructions for ASCII layouts improve the AI's ability to reason about density and hierarchy.
- **Generalization:** Initially restricted to TUI, the requirement was generalized to all UI types based on user feedback to maintain a consistent "low-fi/high-fidelity" design language.
