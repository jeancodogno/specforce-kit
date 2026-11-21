[ [English](getting-started.md) | [Português](../pt/getting-started.md) | [Español](../es/getting-started.md) ]

# Getting Started with Specforce Kit

Welcome to the **Specforce Kit**! This guide will help you set up your environment and take your first steps with Spec-Driven Development (SDD).

## Prerequisites

- Node.js and NPM (for the recommended installation)
- Go 1.26+ and Make (if you wish to build from source)
- A [supported AI agent](supported-tools.md) (Claude Code, Qwen, Kimi Code, OpenCode, KiloCode, Codex, Antigravity or Cursor)

## Installation

The easiest way to install Specforce Kit is globally via NPM:

```bash
npm i -g @jeancodogno/specforce-kit
```

If you prefer to install from source:

```bash
git clone https://github.com/jeancodogno/specforce-kit.git
cd specforce-kit
make build
sudo make install   # installs to /usr/local/bin; use "make install PREFIX=$HOME/.local" to avoid sudo
```

### Troubleshooting

If you encounter a `command not found: specforce` error after installation, it usually means your system's `PATH` does not include the directory where binaries are stored.

#### 1. Fix for macOS/Linux (NPM)
If you installed via `npm i -g`, add the following to your shell profile (e.g., `~/.zshrc` or `~/.bashrc`):

```bash
# Add npm global bin to PATH
export PATH="$(npm config get prefix)/bin:$PATH"
```

Then reload your profile: `source ~/.zshrc`.

#### 2. Fix for Windows (NPM)
Ensure that your npm global directory (usually `%AppData%\npm`) is in your system's Environment Variables.

#### 3. Fix for Source/Git Installation
If you installed from source or Git and the command is missing, ensure you have built the native binary:

```bash
cd specforce-kit
make build
# This creates a 'specforce' binary in the root
```

If the automated build failed during installation, check if you have `go` and `make` installed on your system.

#### 4. Automatic Diagnostics
If the `specforce` command works but fails to find its internal "brain" (the native binary), it will automatically run a diagnostic check and print the exact command you need to run to fix your `PATH`.

## Updates

Specforce keeps itself up to date in the background, following the channel you installed it from:

- **When it checks:** at most once every 6 hours, only when you run an interactive command in a terminal. The check runs detached and never slows the command down; being offline does not cause repeated checks.
- **Where it looks:** npm installs follow the npm tag of the running channel (for example `beta` for a `1.1.0-beta.x` version, `latest` for stable versions); native installs follow stable releases only. From 1.0.0 onward, GitHub pre-releases are marked as pre-release and are never offered automatically; to try one natively, download it from the [releases page](https://github.com/jeancodogno/specforce-kit/releases). Specforce never downgrades.
- **Leaving a pre-release channel:** if you run a published pre-release (`1.0.0-alpha.x`, or `1.0.0-beta.1` to `1.0.0-beta.3`), npm or native, Specforce moves you to 1.0.0 automatically in the background once 1.0.0 is published; nothing to do (on npm, to get it right away, run the command below; native installs can also download it from the [releases page](https://github.com/jeancodogno/specforce-kit/releases)). From 1.0.0 onward, an npm pre-release (for example `1.1.0-beta.1`) keeps receiving the updates of its own channel and never switches to stable on its own. When a newer stable version is published, the next interactive command shows a **STABLE RELEASE AVAILABLE** notice, once per stable version: `Specforce 1.1.0 (stable) is available. You are on the beta channel; to switch, run: npm install -g @jeancodogno/specforce-kit@latest`. The notice never installs anything and never appears for agent commands. Switch whenever you choose, in your terminal:

```bash
npm install -g @jeancodogno/specforce-kit@latest
```

- **When it applies:** npm installs install the exact new version; native binaries are downloaded, verified against the release checksum and swapped in on the next run. The next interactive command shows `(Specforce updated to <version>)` once.
- **When it cannot:** if the install folder is not writable, or an install or verification fails, nothing is retried in a loop and a **MANUAL UPDATE REQUIRED** notice is shown once for that version, with the exact command:

```bash
# npm installs
npm install -g @jeancodogno/specforce-kit@<version>

# native installs: download the release from GitHub, or reinstall from source
sudo make install
```

- **Agents are never disturbed:** commands run by agents (`constitution`, `spec`, `implementation`) never start a check or show a notice, and notices are written only to standard error, so agent output (including `--json`) is never polluted.
- **Diagnostics:** run any command with `DEBUG=1` (for example `DEBUG=1 specforce console`) to print each updater decision to standard error. The per-user updater state lives in `~/.specforce/state.json`.

## First Steps

Skills are referred to by name below (for example `spf.constitution`). How to call them explicitly depends on your tool: see [Invoking skills](supported-tools.md#invoking-skills). A plain-language request that names the skill works everywhere.

### 1. Initialize your Project

Navigate to your project directory (new or existing) and initialize Specforce:

```bash
specforce init
```

Pick your AI tools in the interactive list, or pass their keys directly (for example `specforce init claude cursor`). This creates the `.specforce/` folder, the root `AGENTS.md` rules and the skills of each selected tool.

### 2. Define the Constitution

The **Constitution** is the foundational set of rules for your project (principles, architecture, UI/UX, security, engineering, governance). Instead of polluting your agent's context with dozens of loose rules, Specforce centralizes them in a structured way.

Start your AI agent and ask, in its chat:

```text
Use spf.constitution to set up the constitution: <your project description and tech stack>
```

The agent will ask a few questions to align the project and will create the foundational documents in the `.specforce/docs/` folder.

### 3. Create your First Specification

Every code change in Specforce requires a specification. To add a new feature, ask your agent:

```text
Use spf.spec to plan: <describe the feature you want to build>
```

The agent creates the artifacts required by the spec size; a `medium` spec, the default size, needs `requirements.md` and `tasks.md` (see [Artifacts](artifacts.md)).

---

## Workflow: The Lifecycle

The Specforce Kit workflow is based on **Spec-Driven Development (SDD)**. In the age of AI, implementation has become a commodity, but architectural design without validation is a major risk.

SDD separates **design** from **implementation**, ensuring that no code is generated on the fly ("vibe coding") without a formal contract. The development cycle follows a deterministic path. For ready-made workflows, tips and what to do when a step is blocked, see [Workflows & Tips](workflows.md).

### 1. Foundation: The Constitution (`spf.constitution`)
Everything starts with the **Constitution**. Before writing any code, the project needs clear rules on architecture, UI/UX, security and engineering.
- The agent interacts with you to understand the project's scope, or reads the existing code of a project you are adopting.
- It documents these decisions in `.specforce/docs/`, which serve as the "laws" for all future implementations.
- **Advantage:** The agent only loads the documents it needs, saving tokens and preventing hallucinations.

### 2. Exploration: Discovery (`spf.discovery`)
Use discovery to brainstorm an idea, investigate a bug or find a root cause before committing to a plan.
- Product code stays read-only; disposable evidence (tests, scripts, screenshots) lives in `.specforce/scratch/`.
- When you are ready, ask the agent to formalize the result into a `proposal.md` that planning reuses.

### 3. The Contract: The Specification (`spf.spec`)
Every new feature, fix, or refactor starts here.
- You explain what you need. The agent writes the artifacts required by the spec size:
  - **Requirements** (the "what")
  - **Design** (the "how", architecturally)
  - **Tasks** (the atomic execution roadmap)
  - **Module deltas** (how the change affects the module living specs)
- *Human Action:* You must review and approve this specification. This is your contract with the machine.

### 4. Execution: The Implementation (`spf.implement`)
With the approved specification, the agent enters execution mode.
- The agent strictly follows the task list generated in the previous step, then runs an independent QA step that records a verdict for every acceptance criterion.
- [Validation Hooks](configuration.md) ensure the agent doesn't skip tests, linters, or security steps before marking a task as "finished".
- *Human Action:* Use the `specforce console` command in another terminal to monitor the agent's progress in real-time.

### 5. Closure: Archiving (`spf.archive`)
After the feature is fully coded and verified, the specification needs to be closed.
- The agent applies the module deltas to the living specs in `.specforce/docs/modules/`, updates the Constitution if new systemic decisions emerged, and moves the spec to `.specforce/archive/`.
- Archiving is refused while tasks are pending, the delta is not applied or the QA verdict is FAILED.
- You can customize these rules in `config.yaml` to include steps like updating external trackers or specific cleanup tasks.

---

## Human Role vs. Machine Role

In Specforce, roles are clear:

- **The Human Architect (Governance):** Sets the project's direction, approves specifications, and ensures the system as a whole makes business sense.
- **The AI Agent (Executor):** Executes the approved roadmap, generating high-quality code, tests, and internal documentation efficiently and rapidly, always respecting the Constitution.

---

To explore the terminal commands, see the [CLI Reference](cli.md).
Learn more about the documents generated in this process on the [Artifacts](artifacts.md) page.
Follow step-by-step workflows and tips in [Workflows & Tips](workflows.md).
