# Design: Interactive Consultation Protocol

## 1. Architectural Strategy
We will inject the "Interactive Consultation Protocol" at three levels of the Specforce ecosystem:
1. **Global Manifesto:** The `AGENTS.md` file (generated via Go template).
2. **Role Personas:** The individual agent YAML files in the Kit.
3. **Behavioral Skills:** The markdown content within Skill YAML files.

## 2. Protocol Content (The Standard)
The following text block will be used as the base for all updates:

> **Interactive Consultation Protocol:** 
> If you encounter ambiguity, require user iteration, or need to make a critical decision, you MUST actively prompt the user using the interactive question tool native to your specific AI environment (e.g., `ask_user`, `ask`, `prompt`). Do not halt execution, make blind assumptions, or output generic chat questions. Use your environment's tool to explicitly request the required input.

## 3. Component Updates

### 3.1. Go Template Update
Modify `src/internal/project/agents_md.go`:
- Append the protocol to `agentsMDTemplate` before the `<!-- SPECFORCE_AGENTS_END -->` marker.

**Injection Text:**
```markdown
## 5. Interactive Consultation Protocol
If you encounter ambiguity, require user iteration, or need to make a critical decision, you MUST actively prompt the user using the interactive question tool native to your specific AI environment (e.g., `ask_user`, `ask`, `prompt`). Do not halt execution, make blind assumptions, or output generic chat questions. Use your environment's tool to explicitly request the required input.
```

### 3.2. Agent Kit Updates
Modify the `content` field in all agent YAMLs (`product-analyst.yaml`, `technical-developer.yaml`, etc.).

**Injection Text (Append to end of content):**
```markdown
  ## 5. Interactive Consultation Protocol
  If you encounter ambiguity, require user iteration, or need to make a critical decision, you MUST actively prompt the user using the interactive question tool native to your specific AI environment (e.g., `ask_user`, `ask`, `prompt`). Do not halt execution, make blind assumptions, or output generic chat questions. Use your environment's tool to explicitly request the required input.
```

### 3.3. Skill Updates
- **Consultative Grill:** Update "Interview Rules".
- **Opportunity Framing:** Update "How to use it".

**Injection Text (Integration):**
```markdown
- **Mandatory Tool Usage:** When asking questions, you MUST use your environment's native interaction tool (`ask_user`, `ask`, etc.) to ensure the user is prompted for input.
```

## 4. Verification Plan
- **Template Check:** Run a test or manual `init` to verify `AGENTS.md` contains the new section.
- **YAML Validation:** Use `cat` or similar to verify the content field of the YAML files contains the protocol.
- **Unit Tests:** Run existing tests in `src/internal/project/agents_md_test.go` to ensure no regressions in marker handling.
