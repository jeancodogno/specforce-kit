package core

import (
	"context"
	"testing"
)

func TestExecuteHooks(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		commands []string
		wantErr  bool
	}{
		{
			name:     "Simple echo",
			commands: []string{"echo hello"},
			wantErr:  false,
		},
		{
			name:     "Command with arguments",
			commands: []string{"ls -l"},
			wantErr:  false,
		},
		{
			name:     "Failed command",
			commands: []string{"false"},
			wantErr:  true,
		},
		{
			name:     "Empty command string",
			commands: []string{""},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ExecuteHooks(ctx, tt.commands)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteHooks() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil {
				if err.Error() != "one or more hooks failed" {
					t.Errorf("expected 'one or more hooks failed', got %v", err.Error())
				}
			}
		})
	}
}
