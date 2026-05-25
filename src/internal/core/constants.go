package core

const Version = "1.0.0-alpha.1"

// Tool directory prefixes that are considered "agent tools" and can be updated independently of the project constitution.
var ToolPrefixes = []string{
	".gemini/",
	".claude/",
	".opencode/",
	".kilocode/",
	".agent/",
	".qwen/",
	".codex/",
	".kimi/",
}
