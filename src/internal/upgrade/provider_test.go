package upgrade

import (
	"context"
	"testing"
)

func TestMockProvider(t *testing.T) {
	expected := "v2.0.0"
	p := &MockProvider{Version: expected}

	version, err := p.GetLatestVersion(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if version != expected {
		t.Errorf("expected version %s, got %s", expected, version)
	}
}
