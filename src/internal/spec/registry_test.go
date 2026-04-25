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
