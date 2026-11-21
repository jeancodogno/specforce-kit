---
slug: 20260524-1109-fix-cli-hallucination-and-npm-install
lens: Bugfix
---

# Bugfix: CLI Path Hardening & NPM Recovery Guide

## 1. Issue Description
AI agents (LLMs) frequently hallucinate the execution path for the `specforce` CLI. Instead of using the globally available binary, they attempt to use `npx specforce`, relative paths like `./bin/specforce`, or absolute paths from their training data. This leads to command execution failures, broken workflows, and user frustration as the agent loses momentum trying to find the tool.

## 2. Evidence & Observations
- **Symptom:** Command fails with `command not found` or `no such file or directory` when the agent attempts to run `specforce` using incorrect prefixes.
- **Observations:** LLMs default to standard patterns (like `npx` for Node-based tools) when instructions are ambiguous about environment setup.
- **Trace:** Manual observation of agent tool calls shows `npx specforce spec status` or `./bin/specforce spec init` failing in clean environments.

## 3. Reproduction Steps
1. Initialize a new project in a fresh environment where `specforce` is not yet in the global PATH.
2. Provide the `AGENTS.md` to an LLM (e.g., Gemini or Claude) and ask it to initialize a spec or check status.
3. Observe the LLM's tool call: it will likely attempt `npx specforce` or `./bin/specforce` because the current documentation doesn't explicitly forbid it or define the correct execution path.
4. Expected outcome: The agent should use `specforce <command>` directly or be prompted to install it if missing.
5. Actual outcome: The agent guesses the path and fails.

## 4. Root Cause Analysis (RCA)
- **Ambiguity in AGENTS.md:** The current AI collaboration guide lacks a strict "Execution Environment" section that defines how the CLI should be invoked.
- **Lack of Recovery Path:** There are no clear, machine-readable instructions for the agent to follow if the binary is not found (e.g., a specific `npm install` command).
- **Missing Guardrails:** Command kit instructions (YAML files used by some agents) do not explicitly prohibit path guessing.
- **Missing Versioning:** There is no central, programmatic source of truth for the tool's version, making it harder for agents to verify their environment.

## 5. Acceptance Criteria (Regression Tests)

### [FIX-1] CLI Execution Hardening
**Scenario: [Regression]**
GIVEN an AI agent is reading `AGENTS.md`
WHEN it needs to execute a `specforce` command
THEN it MUST find a "CLI Execution & Environment" section that explicitly prohibits `npx`, absolute paths, and relative paths (e.g., `./bin/`).

### [FIX-2] Environment Recovery Path
**Scenario: [Regression]**
GIVEN the `specforce` binary is missing from the environment
WHEN the agent encounters a `command not found` error
THEN it MUST find a recovery instruction in `AGENTS.md` to run `npm i -g @jeancodogno/specforce-kit@1.0.0-alpha.1`.

### [FIX-3] Command Kit Guardrails
**Scenario: [Regression]**
GIVEN an agent is using a command kit (YAML) to execute tasks
WHEN it prepares a CLI command
THEN the kit's `Guardrails` section MUST explicitly warn against path hallucinations.

### [FIX-4] Centralized Versioning
**Scenario: [Regression]**
GIVEN the need to verify the installed version
WHEN checking the codebase
THEN a `Version` constant MUST exist in `src/internal/core/constants.go`.

## 6. Technical Constraints (NFR)
- **[Safety]:** Recovery instructions must point to a specific, vetted version of the package to avoid unintended upgrades.
- **[Observability]:** The `Version` constant must be easily accessible for future `specforce --version` implementation.
