package cobra

import (
	"testing"
)

func TestImplementationUpdateCmd_Flags(t *testing.T) {
	t.Run("comma-separated tasks", func(t *testing.T) {
		taskIds = nil
		cmd := implementationUpdateCmd
		err := cmd.ParseFlags([]string{"--task", "T1.1,T1.2", "--status", "finished"})
		if err != nil {
			t.Fatalf("ParseFlags failed: %v", err)
		}
		tasks, err := cmd.Flags().GetStringSlice("task")
		if err != nil {
			t.Fatalf("GetStringSlice failed: %v", err)
		}
		if len(tasks) != 2 || tasks[0] != "T1.1" || tasks[1] != "T1.2" {
			t.Errorf("expected [T1.1 T1.2], got %v", tasks)
		}
	})

	t.Run("repeated task flags", func(t *testing.T) {
		taskIds = nil
		cmd := implementationUpdateCmd
		err := cmd.ParseFlags([]string{"--task", "T1.1", "--task", "T1.2", "--status", "finished"})
		if err != nil {
			t.Fatalf("ParseFlags failed: %v", err)
		}
		tasks, err := cmd.Flags().GetStringSlice("task")
		if err != nil {
			t.Fatalf("GetStringSlice failed: %v", err)
		}
		if len(tasks) != 2 || tasks[0] != "T1.1" || tasks[1] != "T1.2" {
			t.Errorf("expected [T1.1 T1.2], got %v", tasks)
		}
	})
}
