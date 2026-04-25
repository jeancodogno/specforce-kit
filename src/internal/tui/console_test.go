package tui

import (
	"context"
	"testing"
	"github.com/jeancodogno/specforce-kit/src/internal/spec"
)

func TestConsoleModel_Initialization(t *testing.T) {
	tree := spec.NewStateTree()
	tree.Categories[spec.CategoryConstitution] = []spec.StateItem{{Name: "Architecture", Category: spec.CategoryConstitution}}
	tree.Categories[spec.CategoryImplementations] = []spec.StateItem{{Name: "Feature A", Category: spec.CategoryImplementations}}

	model := NewConsoleModel(context.Background(), tree, &spec.Registry{}, ".")

	if model.StateTree == nil {
		t.Errorf("Expected StateTree to be initialized")
	}
}
