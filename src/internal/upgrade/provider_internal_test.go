package upgrade

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewGitHubProvider(t *testing.T) {
	p := NewGitHubProvider()
	if p == nil {
		t.Fatal("expected non-nil provider")
	}
	if p.BaseURL != "https://api.github.com" {
		t.Errorf("expected default baseURL, got %s", p.BaseURL)
	}
}

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

	testProvider(t, p, "v0.5.0")
}

func TestNewNPMProvider(t *testing.T) {
	p := NewNPMProvider()
	if p == nil {
		t.Fatal("expected non-nil provider")
	}
	if p.BaseURL != "https://registry.npmjs.org" {
		t.Errorf("expected default BaseURL, got %s", p.BaseURL)
	}
}

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

	testProvider(t, p, "1.2.3")
}

func testProvider(t *testing.T, p Provider, expected string) {
	t.Helper()
	version, err := p.GetLatestVersion(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if version != expected {
		t.Errorf("expected version %s, got %s", expected, version)
	}
}
