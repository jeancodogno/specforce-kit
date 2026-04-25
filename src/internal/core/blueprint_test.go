package core

import (
	"testing"
)

func TestParseBlueprint(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		data    string
		wantErr bool
	}{
		{
			name: "Valid blueprint",
			id:   "test-bp",
			data: `
name: Test Blueprint
description: A test blueprint
mapping:
  agent1:
    path: path/to/agent1
    name: name1
    ext: .md
content: |
  This is the content.
`,
			wantErr: false,
		},
		{
			name:    "Invalid YAML",
			id:      "invalid-yaml",
			data:    `name: : invalid`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBlueprint(tt.id, []byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBlueprint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.ID != tt.id {
					t.Errorf("ParseBlueprint() ID = %v, want %v", got.ID, tt.id)
				}
				if got.Metadata.Name != "Test Blueprint" {
					t.Errorf("ParseBlueprint() Name = %v, want %v", got.Metadata.Name, "Test Blueprint")
				}
				if got.Content != "This is the content." {
					t.Errorf("ParseBlueprint() Content = %v, want %v", got.Content, "This is the content.")
				}
			}
		})
	}
}
