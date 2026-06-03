package core

// Version is the current tool version, synced from package.json and optionally overwritten at build time.
var Version = "1.0.0-alpha.3"

// Tool directory prefixes that are considered "agent tools" and can be updated independently of the project constitution.
var ToolPrefixes = []string{
	".gemini/",
	".claude/",
	".cursor/",
	".opencode/",
	".kilocode/",
	".agents/",
	".qwen/",
	".codex/",
	".kimi/",
}
