package agent

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// InstructionManager handles dynamic instruction fetching and variable injection.
type InstructionManager struct {
	kitFS  fs.FS
	config *core.ProjectConfig
}

// NewInstructionManager creates a new instruction manager.
func NewInstructionManager(kitFS fs.FS, config *core.ProjectConfig) *InstructionManager {
	return &InstructionManager{
		kitFS:  kitFS,
		config: config,
	}
}

// GetInstructions fetches instructions for a category and applies variable injection.
func (m *InstructionManager) GetInstructions(category string) (string, error) {
	// 1. Fetch from kitFS
	kitPath := filepath.Join("instructions", category+".md")
	data, err := fs.ReadFile(m.kitFS, kitPath)
	content := ""
	if err == nil {
		content = string(data)
	}

	// 2. Merge with config instructions (if any)
	if m.config != nil && m.config.Instructions != nil {
		if custom, ok := m.config.Instructions[category]; ok {
			if content != "" {
				content += "\n\n"
			}
			content += strings.Join(custom, "\n")
		}
	}

	// 3. Inject variables
	return m.InjectVariables(content), nil
}

// InjectVariables replaces {{key}} placeholders with values from config context.
func (m *InstructionManager) InjectVariables(content string) string {
	if m.config == nil || m.config.Context == nil {
		return content
	}

	for k, v := range m.config.Context {
		placeholder := fmt.Sprintf("{{%s}}", k)
		content = strings.ReplaceAll(content, placeholder, v)
	}

	return content
}
