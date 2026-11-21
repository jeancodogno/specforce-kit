# Contributing to Specforce Kit

First off, thank you for considering contributing to Specforce Kit! It's people like you that make Specforce a great tool for the entire AI-Engineering community.

## 🚀 Our Philosophy: Spec-First
We drink our own champagne. **All new features, bug fixes, or architectural changes in this repository MUST follow the Spec-Driven Development (SDD) lifecycle.** Skills are referred to by name below; see [Invoking skills](docs/en/supported-tools.md#invoking-skills) for the syntax of your AI tool.

1.  **Specification:** Every change starts with a Spec via your AI agent (`spf.spec`).
2.  **Implementation:** Once approved, we implement using the Specforce roadmap (`spf.implement`).
3.  **Verification & Archive:** Automated tests and manual TUI verification are mandatory before closing the loop (`spf.archive`).

---

## 🐛 Reporting Issues

### Bug Reports
- Use the GitHub Issue tracker.
- Describe the expected behavior vs. actual behavior.
- Provide a minimal reproduction script or steps.
- Include your OS and Specforce version (`specforce --version`).

### Feature Requests
- We love new ideas! Please open an issue to discuss the "Why" and "Value" before diving into technical details.

---

## 🛠 Pull Request Process

1.  **Fork and Branch:** Create a feature branch from `main`.
2.  **Create a Spec:** Use your preferred AI agent (Claude Code, Cursor, etc.) and ask it, for example, "Use spf.spec to plan: {description}". It generates the artifacts required by the spec size (`requirements.md`, `design.md`, `tasks.md` and module deltas). Use `spf.discovery` first when the problem or root cause is not clear yet.
3.  **Implement:** Ask "Use spf.implement to implement {spec}" to let the agent execute the tasks and run the independent QA step.
4.  **Code Standards:**
    - Language: Go 1.26+
    - Formatting: `gofmt` is mandatory.
    - Linting: Must pass `golangci-lint`.
    - Standards: Adhere to our project Constitution (`.specforce/docs/`).
5.  **Tests:** All new logic must have unit tests. Regression tests are required for bug fixes.
6.  **Archive & Commit:** Ask "Use spf.archive to archive {spec}" to update the module living specs and the Constitution. Use [Conventional Commits](https://www.conventionalcommits.org/) (e.g., `feat:`, `fix:`, `docs:`).
7.  **Submit:** Open a PR against the `main` branch. Ensure the CI pipeline passes.

---

## 🤖 Kit Development (New Agent Support)

Specforce Kit is built to be model-agnostic. We welcome contributions that add support for new AI agents or improve existing ones.

### What is a Kit?
The Kit is the set of skills that teach an AI agent how to work with the Specforce framework, plus the mapping that installs them in each tool's hidden directory for native discovery (e.g., `.claude/skills/`, `.cursor/skills/`, `.kilocode/skills/`). It is embedded in the binary:

- **Skill blueprints:** `src/internal/agent/kit/skills/<id>/SKILL.yaml` (for example `spf-spec/SKILL.yaml`). Each blueprint is translated into the native `SKILL.md` of every supported tool at `specforce init`.
- **Tool manifest:** `src/internal/agent/kit/kit.yaml` declares every supported tool: its `name`, `description`, `target` directory and the `mappings` of the categories it receives (e.g., `skills`).

### Adding a New Tool
1.  **Register the directory:** Add the tool's hidden directory (e.g., `.newagent/`) to the `ToolPrefixes` slice in `src/internal/core/constants.go`.
2.  **Define the mapping:** Add a tool entry to `src/internal/agent/kit/kit.yaml` with its `name`, `description`, `target` and `mappings`. Use the mapping-level `target` override and environment variable expansion (e.g., `${VAR:-DEFAULT}`) when the tool needs files outside its default folder.
3.  **Global writes (optional):** If the tool requires files outside the project root, declare `security: { global_write: true }` in its `kit.yaml` entry (see the `codex` entry) and document that location.
4.  **Ensure Context Segmentation:** Make sure the new integration respects our **Segmented Constitution** model, loading only the necessary architecture, UI/UX, or security rules when needed, rather than dumping everything into the context window.

To change what the skills say, edit the `SKILL.yaml` blueprints; never edit the generated files inside tool folders. While developing, run the CLI with `--dev` from the repository root to load the kit from `src/internal/agent/kit/` instead of the copy embedded in the binary (see the [CLI Reference](docs/en/cli.md)).

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
<p align="center">Let's build the future of AI-Orchestrated Engineering together.</p>
