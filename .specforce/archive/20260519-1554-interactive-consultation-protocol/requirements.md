# Requirements: Interactive Consultation Protocol

## 1. Goal
Enforce a standardized, platform-agnostic interactive consultation protocol across all Specforce agents and skills. This ensures that agents (Gemini, Claude, etc.) use their native interaction tools (like `ask_user`) instead of making assumptions or stalling when user input is required.

## 2. Context
Specforce Kit is used by various AI CLIs. Each CLI has its own way of asking questions. Agents must be explicitly instructed to use these tools to maintain flow and resolve ambiguity.

## 3. User Stories
- **US-1 (Project-wide Guidance):** As a developer, I want `AGENTS.md` to contain a clear mandate about interactive consultation so that any agent reading it knows how to handle ambiguity.
- **US-2 (Agent Specificity):** As an AI agent, I want my role definition to explicitly require the use of interactive tools so that I don't alucinate decisions.
- **US-3 (Skill Consistency):** As an AI agent using a skill (like Grill or Opportunity Framing), I want the skill instructions to force me to use the native interaction tool for questions.

## 4. Acceptance Criteria (AC)
- **AC-1:** The `agentsMDTemplate` in `src/internal/project/agents_md.go` MUST include a "## 5. Interactive Consultation Protocol" section.
- **AC-2:** All agent YAML files in `src/internal/agent/kit/agents/` MUST include the protocol in their operational mandates or guardrails.
- **AC-3:** The `consultative-grill` and `opportunity-framing` skills MUST be updated to explicitly mention using native interaction tools for their interview/questioning phases.
- **AC-4:** The protocol text MUST be platform-agnostic, mentioning common tool names like `ask_user`, `ask`, or `prompt` as examples.

## 5. Out of Scope
- Implementing the actual tool calling logic in the Go binary (this is handled by the AI providers/CLIs themselves).
- Specific tool implementations for each CLI (we provide the instructions, not the tool definition).
