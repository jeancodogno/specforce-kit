# Supported AI Agents and Tools

Specforce is designed to be independent of any specific tool or ecosystem (tool-agnostic). However, we offer native *Slash Commands* integration out-of-the-box for the most popular AI coding assistants:

- **Gemini CLI** (Google)
- **Claude Code** (Anthropic)
- **Qwen** (Alibaba)
- **Kimi Code** (Moonshot AI)
- **OpenCode**
- **KiloCode**
- **Codex**
- **Antigravity**

## Automated Configuration

Specforce automatically configures your environment to ensure agents can discover the project rules defined in `AGENTS.md`. When running `specforce init` or updating tools:

- **Gemini CLI**: Automatically creates `.gemini/settings.json` with the correct context mapping:
  ```json
  {
    "context": {
      "fileName": ["AGENTS.md", "GEMINI.md"]
    }
  }
  ```
- **Antigravity & Claude Code**: Automatically creates symbolic links at `.agent/rules/AGENTS.md` and `.claude/rules/AGENTS.md` pointing to the root rules file.

*Can't find your favorite agent? Submit a PR creating a Kit for it!*
