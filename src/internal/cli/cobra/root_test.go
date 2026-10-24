package cobra

import (
	"testing"
)

func TestCommandAnnotations(t *testing.T) {
	agentCommands := []string{"spec", "implementation", "constitution"}

	for _, name := range agentCommands {
		cmd, _, err := rootCmd.Find([]string{name})
		if err != nil {
			t.Errorf("failed to find command %s: %v", name, err)
			continue
		}

		if cmd.Annotations["IsAgentCommand"] != "true" {
			t.Errorf("expected IsAgentCommand: true annotation for command %s, got %v", name, cmd.Annotations)
		}
	}
}

func TestRootCommands(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"implementation", "qa", "record"})
	if err != nil || cmd == nil || cmd.Name() != "record" {
		t.Fatalf("implementation qa record not registered: cmd=%v err=%v", cmd, err)
	}
}
