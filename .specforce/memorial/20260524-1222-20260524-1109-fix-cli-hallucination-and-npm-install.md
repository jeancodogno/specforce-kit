---
date: 2026-05-24
scope: 20260524-1109-fix-cli-hallucination-and-npm-install
author: agent
type: Lesson
---

# LLM CLI Execution Hallucinations

AI agents frequently hallucinate CLI execution paths (e.g., using `npx` or absolute paths) when the execution environment is not explicitly defined. To prevent this, all generated AI instructions (`AGENTS.md`) and command kit YAMLs must contain explicit 'CLI Execution & Environment' guardrails forbidding path guessing and providing a concrete self-healing installation command (e.g., `npm install -g`).
