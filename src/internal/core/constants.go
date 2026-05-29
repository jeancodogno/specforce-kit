package core

const Version = "1.0.0-alpha.2"

// Tool directory prefixes that are considered "agent tools" and can be updated independently of the project constitution.
var ToolPrefixes = []string{
	".gemini/",
	".claude/",
	".opencode/",
	".kilocode/",
	".agents/",
	".qwen/",
	".codex/",
	".kimi/",
}
