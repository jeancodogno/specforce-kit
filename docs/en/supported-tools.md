[ [English](supported-tools.md) | [Português](../pt/supported-tools.md) | [Español](../es/supported-tools.md) ]

# Supported AI Agents and Tools

Specforce is designed to be independent of any specific tool or ecosystem (tool-agnostic). However, we offer native *Skills* integration out-of-the-box for the most popular AI coding assistants:

- **Claude Code** (Anthropic)
- **Qwen** (Alibaba)
- **Kimi Code** (Moonshot AI)
- **OpenCode**
- **KiloCode**
- **Codex**
- **Antigravity**
- **Cursor** (Cursor AI) - Uses standard `.md` files in `.cursor/skills/`. Natively supports the root `AGENTS.md`.

## Automated Configuration

Specforce automatically configures your environment to ensure agents can discover the project rules defined in `AGENTS.md` and their native skills. All tools use standard skills directories (e.g., `.agents/skills/`, `.claude/skills/`, `.cursor/skills/`, `.opencode/skills/`, etc.). When running `specforce init` or updating tools:

- **Antigravity & Claude Code**: Automatically creates symbolic links at `.agent/rules/AGENTS.md` and `.claude/rules/AGENTS.md` pointing to the root rules file.
- **Cursor**: Automatically creates the `.cursor/skills/` directory structure. It natively reads the root `AGENTS.md` for project context.

*Can't find your favorite agent? Submit a PR creating a Kit for it!*
