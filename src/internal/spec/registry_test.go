package spec

import (
	"strings"
	"testing"
	"testing/fstest"
)

func TestNewRegistry(t *testing.T) {
	fs := fstest.MapFS{
		"requirements.yaml": {
			Data: []byte("description: Req Desc\ninstruction: Req Inst\ntemplate: Req Temp\n"),
		},
		"design.yaml": {
			Data: []byte("description: Design Desc\ninstruction: Design Inst\ntemplate: Design Temp\ndependency: requirements\n"),
		},
		"tasks.yaml": {
			Data: []byte("description: Task Desc\ninstruction: Task Inst\ntemplate: Task Temp\ndependency: design\n"),
		},
	}

	registry, err := NewRegistry(fs)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}
	if registry == nil {
		t.Fatal("NewRegistry returned nil registry")
	}

	// Check Get
	art, ok := registry.Get("requirements")
	if !ok {
		t.Error("artifact 'requirements' not found")
	}
	if art.Name != "requirements" {
		t.Errorf("expected name 'requirements', got '%s'", art.Name)
	}
	if art.Description != "Req Desc" {
		t.Errorf("expected description 'Req Desc', got '%s'", art.Description)
	}

	art, ok = registry.Get("design")
	if !ok {
		t.Error("artifact 'design' not found")
	}
	if art.Dependency != "requirements" {
		t.Errorf("expected dependency 'requirements', got '%s'", art.Dependency)
	}

	// Check List order
	list := registry.List()
	if len(list) != 3 {
		t.Errorf("expected list length 3, got %d", len(list))
	}
	if list[0].Name != "requirements" {
		t.Errorf("expected first element 'requirements', got '%s'", list[0].Name)
	}
	if list[1].Name != "design" {
		t.Errorf("expected second element 'design', got '%s'", list[1].Name)
	}
	if list[2].Name != "tasks" {
		t.Errorf("expected third element 'tasks', got '%s'", list[2].Name)
	}
}

func TestNewRegistry_CircularDependency(t *testing.T) {
	fs := fstest.MapFS{
		"a.yaml": {
			Data: []byte("description: A\ninstruction: A\ntemplate: A\ndependency: b\n"),
		},
		"b.yaml": {
			Data: []byte("description: B\ninstruction: B\ntemplate: B\ndependency: a\n"),
		},
	}

	registry, err := NewRegistry(fs)
	if err == nil {
		t.Fatal("expected error for circular dependency, got nil")
	}
	if !strings.Contains(err.Error(), "circular dependency detected") {
		t.Errorf("expected circular dependency error message, got: %v", err)
	}
	if registry != nil {
		t.Error("expected nil registry on error")
	}
}

func TestNewRegistry_MissingFields(t *testing.T) {
	fs := fstest.MapFS{
		"invalid.yaml": {
			Data: []byte("description: Only Desc\n"),
		},
	}

	registry, err := NewRegistry(fs)
	if err == nil {
		t.Fatal("expected error for missing fields, got nil")
	}
	if !strings.Contains(err.Error(), "missing 'instruction'") {
		t.Errorf("expected missing instruction error, got: %v", err)
	}
	if registry != nil {
		t.Error("expected nil registry on error")
	}
}

func TestRegistry_TypeAwareness(t *testing.T) {
	fs := fstest.MapFS{
		"requirements.yaml": {
			Data: []byte("description: Feature Req\ninstruction: Feature Inst\ntemplate: Feature Temp\n"),
		},
		"bug-requirements.yaml": {
			Data: []byte("description: Bug Req\ninstruction: Bug Inst\ntemplate: Bug Temp\n"),
		},
		"design.yaml": {
			Data: []byte("description: Design Desc\ninstruction: Design Inst\ntemplate: Design Temp\n"),
		},
	}

	registry, err := NewRegistry(fs)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}

	// 1. Test GetForType with specific type
	art, ok := registry.GetForType("bug", "requirements")
	if !ok {
		t.Fatal("bug-requirements not found")
	}
	if art.Description != "Bug Req" {
		t.Errorf("expected 'Bug Req', got %q", art.Description)
	}
	if art.Name != "requirements" {
		t.Errorf("expected Name to be 'requirements' (normalized), got %q", art.Name)
	}

	// 2. Test GetForType fallback
	art, ok = registry.GetForType("bug", "design")
	if !ok {
		t.Fatal("design not found (fallback)")
	}
	if art.Description != "Design Desc" {
		t.Errorf("expected 'Design Desc', got %q", art.Description)
	}

	// 4. Test Get with prefix
	art, ok = registry.Get("bug-requirements")
	if !ok {
		t.Fatal("Get bug-requirements failed")
	}
	if art.Description != "Bug Req" {
		t.Errorf("expected 'Bug Req', got %q", art.Description)
	}

	// 5. Test Get with prefix fallback
	art, ok = registry.Get("bug-design")
	if !ok {
		t.Fatal("Get bug-design failed")
	}
	if art.Description != "Design Desc" {
		t.Errorf("expected 'Design Desc', got %q", art.Description)
	}
}

func TestTopologicalSort(t *testing.T) {
	t.Run("Linear dependency A->B->C", func(t *testing.T) {
		fs := fstest.MapFS{
			"c.yaml": {
				Data: []byte("description: C\ninstruction: C\ntemplate: C\ndependency: b\n"),
			},
			"b.yaml": {
				Data: []byte("description: B\ninstruction: B\ntemplate: B\ndependency: a\n"),
			},
			"a.yaml": {
				Data: []byte("description: A\ninstruction: A\ntemplate: A\n"),
			},
		}

		registry, err := NewRegistry(fs)
		if err != nil {
			t.Fatalf("NewRegistry failed: %v", err)
		}

		list := registry.List()
		if len(list) != 3 {
			t.Fatalf("expected 3 artifacts, got %d", len(list))
		}

		// Topological order for C depends on B, B depends on A should be A, B, C
		if list[0].Name != "a" || list[1].Name != "b" || list[2].Name != "c" {
			t.Errorf("expected order [a, b, c], got [%s, %s, %s]", list[0].Name, list[1].Name, list[2].Name)
		}
	})

	t.Run("Circular dependency A->A", func(t *testing.T) {
		fs := fstest.MapFS{
			"a.yaml": {
				Data: []byte("description: A\ninstruction: A\ntemplate: A\ndependency: a\n"),
			},
		}

		_, err := NewRegistry(fs)
		if err == nil {
			t.Fatal("expected error for circular dependency A->A, got nil")
		}
		if !strings.Contains(err.Error(), "circular dependency detected") {
			t.Errorf("expected circular dependency error, got: %v", err)
		}
	})
}

func createTestSizeRegistry(t *testing.T) *Registry {
	t.Helper()
	fs := fstest.MapFS{
		"requirements.yaml": {
			Data: []byte("description: Feature Req\ninstruction: Feature Inst\ntemplate: Feature Temp\n"),
		},
		"design.yaml": {
			Data: []byte("description: Feature Design\ninstruction: Feature Design Inst\ntemplate: Feature Design Temp\ndependency: requirements\n"),
		},
		"tasks.yaml": {
			Data: []byte("description: Feature Tasks\ninstruction: Feature Tasks Inst\ntemplate: Feature Tasks Temp\ndependency: design\n"),
		},
		"bug-requirements.yaml": {
			Data: []byte("description: Bug Req\ninstruction: Bug Inst\ntemplate: Bug Temp\n"),
		},
		"bug-tasks.yaml": {
			Data: []byte("description: Bug Tasks\ninstruction: Bug Tasks Inst\ntemplate: Bug Tasks Temp\n"),
		},
	}

	registry, err := NewRegistry(fs)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}
	return registry
}

type expectedArtifact struct {
	name        string
	description string
}

func verifyArtifacts(t *testing.T, arts []Artifact, expected []expectedArtifact) {
	t.Helper()
	if len(arts) != len(expected) {
		t.Fatalf("expected %d artifacts, got %d", len(expected), len(arts))
	}
	for i, exp := range expected {
		if exp.name != "" && arts[i].Name != exp.name {
			t.Errorf("artifact[%d]: expected name %q, got %q", i, exp.name, arts[i].Name)
		}
		if exp.description != "" && arts[i].Description != exp.description {
			t.Errorf("artifact[%d]: expected description %q, got %q", i, exp.description, arts[i].Description)
		}
	}
}

type typeAndSizeTestCase struct {
	name     string
	specType string
	size     SpecSize
	expected []expectedArtifact
}

func listForTypeAndSizeCases() []typeAndSizeTestCase {
	return []typeAndSizeTestCase{
		{
			name:     "Small size feature spec",
			specType: "feature",
			size:     SpecSizeSmall,
			expected: []expectedArtifact{{name: "tasks", description: "Feature Tasks"}},
		},
		{
			name:     "Small size bug spec",
			specType: "bug",
			size:     SpecSizeSmall,
			expected: []expectedArtifact{{name: "tasks", description: "Bug Tasks"}},
		},
		{
			name:     "Medium size feature spec",
			specType: "feature",
			size:     SpecSizeMedium,
			expected: []expectedArtifact{{name: "requirements"}, {name: "tasks"}},
		},
		{
			name:     "Medium size bug spec",
			specType: "bug",
			size:     SpecSizeMedium,
			expected: []expectedArtifact{
				{name: "requirements", description: "Bug Req"},
				{name: "tasks", description: "Bug Tasks"},
			},
		},
		{
			name:     "Large size feature spec",
			specType: "feature",
			size:     SpecSizeLarge,
			expected: []expectedArtifact{{name: "requirements"}, {name: "design"}, {name: "tasks"}},
		},
		{
			name:     "Complex size feature spec",
			specType: "feature",
			size:     SpecSizeComplex,
			expected: []expectedArtifact{{name: "requirements"}, {name: "design"}, {name: "tasks"}},
		},
		{
			name:     "Default/Unknown fallback size returns full matrix",
			specType: "feature",
			size:     "",
			expected: []expectedArtifact{{name: "requirements"}, {name: "design"}, {name: "tasks"}},
		},
	}
}

func TestRegistryListForTypeAndSize(t *testing.T) {
	registry := createTestSizeRegistry(t)

	for _, tt := range listForTypeAndSizeCases() {
		t.Run(tt.name, func(t *testing.T) {
			arts := registry.ListForTypeAndSize(tt.specType, tt.size)
			verifyArtifacts(t, arts, tt.expected)
		})
	}
}


