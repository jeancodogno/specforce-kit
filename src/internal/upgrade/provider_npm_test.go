package upgrade

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNPMProvider(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/@jeancodogno/specforce-kit/latest" {
			t.Errorf("expected path /@jeancodogno/specforce-kit/latest, got %s", r.URL.Path)
		}
		_, _ = fmt.Fprint(w, `{"version": "1.2.3"}`)
	}))
	defer ts.Close()

	p := &NPMProvider{
		BaseURL: ts.URL,
		Client:  ts.Client(),
	}

	version, err := p.GetLatestVersion(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if version != "1.2.3" {
		t.Errorf("expected version 1.2.3, got %s", version)
	}
}
