[ [English](supported-tools.md) | [Português](../pt/supported-tools.md) | [Español](../es/supported-tools.md) ]

# Supported AI Agents and Tools

Specforce is tool-agnostic: the Constitution, the specifications and the CLI are the same whatever AI coding assistant you use. On top of that, it installs its *Skills* natively for these tools:

- **Claude Code** (Anthropic)
- **Qwen** (Alibaba)
- **Kimi Code** (Moonshot AI)
- **OpenCode**
- **KiloCode**
- **Codex**
- **Antigravity**
- **Cursor**

## Installing the skills

`specforce init` opens an interactive picker where you choose the tools to configure. To skip the picker, pass the tool keys directly:

```bash
specforce init claude cursor
```

Accepted tool keys: `claude`, `qwen`, `kimi-code`, `open-code`, `kilo-code`, `codex`, `antigravity`, `cursor`.

Every tool receives the same skills: `spf.discovery`, `spf.constitution`, `spf.spec`, `spf.implement` and `spf.archive`, plus `consultative-grill`, an internal helper used by `spf.spec` during planning. Each skill is installed as a folder named `spf-<name>` (for example `spf-spec/SKILL.md`) whose display name is `spf.<name>`.

## Invoking skills

Every tool can pick a skill on its own when your request matches it, and every tool understands a plain-language request that names the skill, for example "Use spf.spec to plan: add CSV export to reports". The table shows where the skills are installed and the explicit invocation syntax documented by each tool, using `spf.spec` as the example.

| Tool | Install folder | Explicit invocation | Notes |
|---|---|---|---|
| Claude Code | `.claude/skills/` | `/spf.spec <request>` (the folder name, `/spf-spec`, also works) | Key `claude`. Also receives `.claude/rules/AGENTS.md` (see below). |
| Qwen | `.qwen/skills/` | `/spf.spec <request>` | Key `qwen`. |
| Kimi Code | `.kimi/skills/` | `/skill:spf.spec <request>` | Key `kimi-code`. |
| OpenCode | `.opencode/skills/` | Mention `@spf-spec` in your prompt | Key `open-code`. |
| KiloCode | `.kilocode/skills/` | Ask by name in chat (e.g. 'use spf.spec ...') | Key `kilo-code`. |
| Codex | `.codex/skills/` | `/prompts:spf-spec <request>` | Key `codex`. Also writes the prompt files `spf-*.md` **outside the project**, in `${CODEX_PROMPTS_DIR:-${CODEX_HOME:-~/.codex}/prompts}`: the folder set by `CODEX_PROMPTS_DIR`, otherwise `$CODEX_HOME/prompts`, otherwise `~/.codex/prompts`. |
| Antigravity | `.agents/skills/` | Ask by name in chat (e.g. 'use spf.spec ...') | Key `antigravity`. |
| Cursor | `.cursor/skills/` | Ask by name in chat (e.g. 'use spf.spec ...') | Key `cursor`. Reads the root `AGENTS.md` natively. |

The table uses `spf.spec` as the example. For the other skills, just swap the name: in Claude Code, `/spf.discovery` or `/spf.archive`; in Codex, `/prompts:spf-implement`. Tools marked "Ask by name in chat" have no command: write your request as usual and name the skill, for example "use spf.archive to archive csv-export".

## Automated Configuration

When you run `specforce init` or update the tools, Specforce makes sure each agent can discover the project rules and its skills:

- **Root `AGENTS.md`:** Specforce writes its rules between the `<!-- SPECFORCE_AGENTS_START -->` and `<!-- SPECFORCE_AGENTS_END -->` markers. Anything you write outside that block is preserved.
- **Claude Code:** Only Claude Code receives a link: `.claude/rules/AGENTS.md` is created as a symbolic link to the root `AGENTS.md`.
- **Cursor and the other tools:** No link is created. Cursor reads the root `AGENTS.md` natively; the other tools read the root guide or their skills folder directly.
- **Legacy links:** Older `AGENTS.md` symbolic links inside `.agents/**/rules/` folders are removed.

*Can't find your favorite agent? Submit a PR creating a Kit for it!*
