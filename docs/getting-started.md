# Getting Started with Specforce Kit

Welcome to the **Specforce Kit**! This guide will help you set up your environment and take your first steps with Spec-Driven Development (SDD).

## Prerequisites

- Node.js and NPM (for the recommended installation)
- Go 1.26+ and Make (if you wish to build from source)
- Your favorite AI agent (Gemini CLI, Claude Code, KiloCode, etc.)

## Installation

The easiest way to install Specforce Kit is globally via NPM:

```bash
npm install -g @specforce/kit
```

If you prefer to install from source:

```bash
git clone https://github.com/jeancodogno/specforce-kit.git
cd specforce-kit
make build
make install
```

## First Steps

### 1. Initialize your Project

Navigate to your project directory (new or existing) and initialize Specforce:

```bash
specforce init
```

This command will prepare the necessary structure, setting up the state folders (`.specforce/`) and base files for your AI agent.

### 2. Define the Constitution

The **Constitution** is the foundational set of rules for your project (architecture, UI/UX, security). Instead of polluting your agent's context with dozens of loose rules, Specforce centralizes this in a structured way.

In your terminal, start your AI agent (e.g., `gemini`) and run:

```text
/spf:constitution {your project description and tech stack}
```

The agent will ask a few questions to align the project and will create the foundational documents in the `docs/` folder.

### 3. Create your First Specification

Every code change in Specforce requires a specification. To add a new feature, ask your agent:

```text
/spf:spec {describe the feature you want to build}
```

The agent will create the essential artifacts: `requirements.md`, `design.md`, and `tasks.md`.

---

## Workflow: The 4-Step Lifecycle

The Specforce Kit workflow is based on **Spec-Driven Development (SDD)**. In the age of AI, implementation has become a commodity, but architectural design without validation is a major risk.

SDD separates **design** from **implementation**, ensuring that no code is generated on the fly ("vibe coding") without a formal contract.

The development cycle with AI agents in Specforce follows a deterministic path:

### 1. Foundation: The Constitution (`/spf:constitution`)
Everything starts with the **Constitution**. Before writing any code, the project needs clear rules on Architecture, UI/UX, and Security.
- The agent interacts with you to understand the project's scope.
- It documents these decisions, which will serve as immutable "laws" for all future implementations.
- **Advantage:** The agent only loads the segments of the constitution it needs, saving tokens and preventing hallucinations.

### 2. The Contract: The Specification (`/spf:spec`)
Every new feature, fix, or refactor starts here.
- You explain what you need. The agent translates your needs into 3 artifacts:
  - **Requirements** (the "what")
  - **Design** (the "how", architecturally)
  - **Tasks** (the atomic execution roadmap)
- *Human Action:* You must review and approve this specification. This is your contract with the machine.

### 3. Execution: The Implementation (`/spf:implement`)
With the approved specification, the agent enters execution mode.
- The agent strictly follows the task list generated in the previous step.
- Internal tools, such as [Validation Hooks](configuration.md), ensure the agent doesn't skip tests, linters, or security steps before marking a task as "finished".
- *Human Action:* Use the `specforce console` command in another terminal to monitor the agent's progress in real-time.

### 4. Closure: Archiving (`/spf:archive`)
After the feature is fully coded and tested, the specification needs to be closed.
- The agent archives the completed specification to keep the environment clean.
- It follows standardized **Archiving Instructions** to analyze what was done and update the Global Constitution (and the project Memorial) if new systemic architectural decisions have emerged.
- You can customize these rules in `config.yaml` to include steps like updating external trackers or specific cleanup tasks.

---

## Human Role vs. Machine Role

In Specforce, roles are clear:

- **The Human Architect (Governance):** Sets the project's direction, approves specifications, and ensures the system as a whole makes business sense.
- **The AI Agent (Executor):** Executes the approved roadmap, generating high-quality code, tests, and internal documentation efficiently and rapidly, always respecting the Constitution.

---

To explore the terminal commands, see the [CLI Reference](cli.md).
Learn more about the documents generated in this process on the [Artifacts](artifacts.md) page.
