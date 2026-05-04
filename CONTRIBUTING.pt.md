# Contributing to Specforce Kit

First off, thank you for considering contributing to Specforce Kit! It's people like you that make Specforce a great tool for the entire AI-Engineering community.

## 🚀 Our Philosophy: Spec-First
We drink our own champagne. **All new features, bug fixes, or architectural changes in this repository MUST follow the Spec-Driven Development (SDD) lifecycle.**

1.  **Specification:** Every change starts with a Spec via your AI agent (`/spf:spec`).
2.  **Implementation:** Once approved, we implement using the Specforce roadmap (`/spf:implement`).
3.  **Verification & Archive:** Automated tests and manual TUI verification are mandatory before closing the loop (`/spf:archive`).

---

## 🐛 Reporting Issues

### Bug Reports
- Use the GitHub Issue tracker.
- Describe the expected behavior vs. actual behavior.
- Provide a minimal reproduction script or steps.
- Include your OS and Specforce version (`specforce version`).

### Feature Requests
- We love new ideas! Please open an issue to discuss the "Why" and "Value" before diving into technical details.

---

## 🛠 Pull Request Process

1.  **Fork and Branch:** Create a feature branch from `main`.
2.  **Create a Spec:** Use your preferred AI agent (Claude Code, Cursor, etc.) and run the `/spf:spec {description}` slash command to generate the `requirements.md`, `design.md`, and `tasks.md`.
3.  **Implement:** Use the `/spf:implement` command to let the agent execute the tasks.
4.  **Code Standards:**
    - Language: Go 1.26+
    - Formatting: `gofmt` is mandatory.
    - Linting: Must pass `golangci-lint`.
    - Standards: Adhere to our project Constitution.
5.  **Tests:** All new logic must have unit tests. Regression tests are required for bug fixes.
6.  **Archive & Commit:** Run `/spf:archive` to update the global Constitution. Use [Conventional Commits](https://www.conventionalcommits.org/) (e.g., `feat:`, `fix:`, `docs:`).
7.  **Submit:** Open a PR against the `main` branch. Ensure the CI pipeline passes.

---

## 🤖 Kit Development (New Agent Support)

Specforce Kit is built to be model-agnostic. We welcome contributions that add support for new AI agents or improve existing ones.

### What is a Kit?
A Kit is a set of instructions, rules, and tool mappings that allow an AI agent to interact with the Specforce framework. These are stored in hidden directories to ensure native discovery by the agents (e.g., `.claude/`, `.gemini/`, `.kilocode/`, etc.).

### Creating a New Kit
1.  **Define the Instructions:** Create the necessary markdown files (e.g., `SKILL.md`, `commands.md`) that teach the agent how to call `specforce` CLI commands via the slash commands (`/spf:spec`, `/spf:implement`, etc.).
2.  **Ensure Context Segmentation:** Make sure the new agent integration respects our **Segmented Constitution** model, loading only the necessary architecture, UI/UX, or security rules when needed, rather than dumping everything into the context window.

---

## 🏗 Development Environment

### Prerequisites
- Go 1.26+
- Make

### Setup
```bash
git clone https://github.com/jeancodogno/specforce-kit.git
cd specforce-kit
make build
make install
```

### Running Tests
```bash
make test
```

---

## 📜 Code of Conduct
By participating in this project, you agree to abide by our project's Code of Conduct and maintain a professional and respectful environment for all contributors.

---
<p align="center">Let's build the future of AI-Orchestrated Engineering together.</p>r.</p>