package upgrade

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitHubProvider(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/jeancodogno/specforce-kit/releases/latest" {
			t.Errorf("expected path /repos/jeancodogno/specforce-kit/releases/latest, got %s", r.URL.Path)
		}
		_, _ = fmt.Fprint(w, `{"tag_name": "v0.5.0"}`)
	}))
	defer ts.Close()

	p := &GitHubProvider{
		BaseURL: ts.URL,
		Client:  ts.Client(),
	}

	version, err := p.GetLatestVersion(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if version != "v0.5.0" {
		t.Errorf("expected version v0.5.0, got %s", version)
	}
}
