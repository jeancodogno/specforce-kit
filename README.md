
<p align="center">
  <a href="https://github.com/jeancodogno/specforce-kit/">
    <picture>
      <img src="assets/logo.png" alt="OpenSpec logo" height="128">
    </picture>
  </a>
</p>
<p align="center">Ecosystem for AI-assisted development.</p>

# Specforce Kit

> **The Orchestration Layer for AI-Native Engineering.**

Specforce Kit is a professional-grade framework designed to coordinate AI agents (Claude Code, Cursor, KiloCode) through **Spec-Driven Development (SDD)**. It provides the standardized **Kits**, **Constitution**, and **Orchestration Tools** required to transform AI from a "vibe-based" assistant into a deterministic engineering force.

---

## 🧠 Philosophy

In the age of AI, implementation is a commodity, but design is a liability. **Implementation without a specification is strictly forbidden in Specforce.** 

Spec-Driven Development (SDD) decouples design from implementation. A specification acts as a strict contract between human designers and machine executors. By doing so, we prevent "vibe coding" - the accumulation of ad-hoc, unverified decisions - ensuring your system remains coherent, secure, and maintainable no matter which AI agent writes the code.

---

## 📊 Why use Specforce?

AI coding without strict guardrails inevitably leads to "vibe coding": vague prompts, unpredictable outputs, and a rapidly accumulating mountain of technical debt. While other tools attempt to solve this, they often force you into new IDE ecosystems, require heavy setups, or bloat your AI's context window.

Specforce is designed to be the **invisible orchestration layer**. It brings predictability, architectural control, and strict verification to your AI workflows, acting as an active state machine rather than just a collection of static templates.

### Key Strengths:

- **Segmented Constitution**: Instead of a single massive prompt, rules are broken down into specialized domains (Architecture, UI/UX, Security). The agent only loads the exact context it needs, keeping token costs low and reducing hallucinations.
- **Dynamic Verification Hooks**: Run tests, linters, or custom scripts automatically before an agent is allowed to mark a task as finished.
- **Git Worktree Support**: Natively discover and view specifications across multiple git branches simultaneously without switching contexts.
- **Context-Aware Instructions**: Inject specific rules and constraints dynamically based on the artifact or phase the agent is currently working on.
- **Tool & IDE Agnostic**: Bring Spec-Driven Development directly to your terminal. Works flawlessly with your favorite IDE and CLI-based AI agents (Gemini, Claude, Cursor, KiloCode, etc).
- **Lightning Fast & Zero Dependencies**: Distributed as a single binary, requiring no heavy Python setups or forced ecosystems to bootstrap your project.

---

## 🤖 Agent Usage (Slash Commands)

Specforce interacts seamlessly with your preferred AI agent via simple slash commands. Here is the standard workflow:

### 1. Discovery (`/spf:discovery`)
**Command:** `/spf:discovery {idea or bug report}`

The early exploration phase. Use this to brainstorm new features or investigate technical issues without modifying any files.
- **Read-Only Exploration:** The agent is strictly forbidden from writing code or artifacts. It focuses on research, root cause analysis, and technical strategy.
- **Expert Personas:** Automatically switches between "Senior Product Architect" (Brainstorming) and "Senior Systems Engineer" (Detective) based on your input.
- **Project Alignment:** Reads your project's Constitution (`.specforce/docs/`) to ensure all ideas align with your principles and architecture.

### 2. The Constitution (`/spf:constitution`)
**Command:** `/spf:constitution {project description}`

The first step in any project. This generates your project's Constitution, containing all rules, principles, UI/UX guidelines, architecture, security, and agent memory. 
- **Tool Agnostic & Segmented:** Instead of dumping rules into tool-specific files (like `.clauderc` or `.gemini/GEMINI.md`) or keeping everything in context, Specforce maintains its own specialized, segmented memory files. The agent only loads what it needs. You can switch tools (e.g., from Gemini CLI to Claude Code) and the project memory remains completely intact.
- **Interactive Setup:** Describe your project idea and stack. The agent will ask clarifying questions to help you make foundational decisions.
- **Flexible:** Run it once to set up, or rerun it anytime to update the constitution. You can use it on brand new ideas or point it at an existing project to analyze and document its current architecture.

### 2. Creating a Spec (`/spf:spec`)
**Command:** `/spf:spec {description of what to add or modify}`

Ready to build? The agent will clarify your requirements and generate three core documents: `requirements.md`, `design.md`, and `tasks.md`.

### 3. Implementation (`/spf:implement`)
**Command:** `/spf:implement`

After you validate the generated spec documents, execute this command. The agent will strictly follow the planned tasks in the spec to implement the feature, ensuring zero drift from the agreed design.

### 4. Archiving (`/spf:archive`)
**Command:** `/spf:archive`

Once implementation is fully complete and verified, this command analyzes what needs to be updated in the project's global Constitution based on the new implementation, and archives the completed spec. It follows standardized **Archiving Instructions** that ensure lessons learned are captured in the project memorial and any project-specific cleanup is performed.

---

## ⚡ Human Touchpoints

While Specforce powers your AI agents, humans interact primarily through two terminal commands:

### 1. Project Initialization
Set up the governance and agent kits for your project.
```bash
specforce init
```

### 2. The Console (TUI)
Launch the command center to monitor project health, spec progress, and agent implementations in real-time.
```bash
specforce console
```

---

## 🛠️ Supported AI Agents & Tools

Specforce Kit is designed to be tool-agnostic but comes with pre-built kits and native slash command support for the most popular AI coding assistants:

- **Gemini CLI** (Google)
- **Claude Code** (Anthropic)
- **Qwen** (Alibaba)
- **Kimi Code** (Moonshot AI)
- **OpenCode**
- **KiloCode**
- **Codex**
- **Antigravity**

*Don't see your favorite tool? We welcome PRs to add new kits!*

---

## 🚀 Getting Started

### Installation

**Using NPM (Recommended)**
```bash
npm i -g @jeancodogno/specforce-kit
```

**From Source (Requires Go 1.26+ & Make)**
```bash
git clone https://github.com/jeancodogno/specforce-kit.git
cd specforce-kit
make build
make install
```

---

## 📚 Documentation

For a deeper dive into Specforce Kit, check out our official documentation:

- [Getting Started & Workflow](docs/getting-started.md): Environment setup, first steps, and the recommended Spec-Driven Development lifecycle.
- [Git Worktree Support](docs/git-worktrees.md): Guide to cross-branch specification discovery, unified console, and read-only constraints.
- [Artifacts](docs/artifacts.md): Details about the generated files (Constitution, Requirements, Design, Tasks).
- [Configuration](docs/configuration.md): Guide to customize hooks and instructions.
- [CLI Reference](docs/cli.md): Terminal and slash command usage.
- [Supported Tools](docs/supported-tools.md): List of natively supported AI agents and coding assistants.

---

## 🔍 Troubleshooting

### "command not found: specforce"
If you encounter this error after installation, it usually means your npm global bin directory is not in your system's `PATH`.

**Fix for macOS/Linux:**
Add the following line to your shell profile (e.g., `~/.zshrc` or `~/.bashrc`):
```bash
export PATH="$(npm config get prefix)/bin:$PATH"
```

**Fix for Windows:**
Ensure that `%AppData%\npm` is in your environment variables.

For more detailed troubleshooting, see [Getting Started: Troubleshooting](docs/getting-started.md#troubleshooting).

---

## 🤝 Contributing
We welcome contributions! See our [CONTRIBUTING.md](CONTRIBUTING.md) to learn how to add support for new AI agents or improve existing kits.

---

## 📜 License
Specforce Kit is released under the **MIT License**. See [LICENSE](LICENSE) for details.