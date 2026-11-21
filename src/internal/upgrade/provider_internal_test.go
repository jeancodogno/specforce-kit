package upgrade

import (
	"context"
	"errors"
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
	ts := newDistTagsRegistry(t, `{"latest": "1.2.3", "beta": "1.3.0-beta.1"}`, http.StatusOK)

	p := &NPMProvider{
		BaseURL: ts.URL,
		Client:  ts.Client(),
	}

	testProvider(t, p, "1.2.3")
}

// newDistTagsRegistry serves body as the dist-tags document of the Specforce package.
func newDistTagsRegistry(t *testing.T, body string, status int) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/-/package/@jeancodogno/specforce-kit/dist-tags" {
			t.Errorf("expected dist-tags path, got %s", r.URL.Path)
		}
		w.WriteHeader(status)
		_, _ = fmt.Fprint(w, body)
	}))
	t.Cleanup(ts.Close)
	return ts
}

func TestNPMProvider_ChannelTag(t *testing.T) {
	ts := newDistTagsRegistry(t, `{"latest": "0.2.2", "alpha": "1.0.0-alpha.3", "beta": "1.0.0-beta.3"}`, http.StatusOK)
	cases := map[string]string{
		"1.0.0-beta.2":  "1.0.0-beta.3",
		"1.0.0-alpha.1": "1.0.0-alpha.3",
		"0.2.1":         "0.2.2",
	}
	for running, want := range cases {
		t.Run(running, func(t *testing.T) {
			p := &NPMProvider{BaseURL: ts.URL, Client: ts.Client(), Channel: ChannelOf(running)}
			testProvider(t, p, want)
		})
	}
}

func TestNPMProvider_MissingChannelTagIsNoUpdate(t *testing.T) {
	ts := newDistTagsRegistry(t, `{"latest": "0.2.2"}`, http.StatusOK)
	p := &NPMProvider{BaseURL: ts.URL, Client: ts.Client(), Channel: "rc"}
	_, err := p.GetLatestVersion(context.Background())
	if !errors.Is(err, ErrNoChannelTag) {
		t.Errorf("expected ErrNoChannelTag, got %v", err)
	}
}

func TestNPMProvider_RegistryErrorIsNotNoUpdate(t *testing.T) {
	ts := newDistTagsRegistry(t, `oops`, http.StatusInternalServerError)
	p := &NPMProvider{BaseURL: ts.URL, Client: ts.Client(), Channel: "beta"}
	_, err := p.GetLatestVersion(context.Background())
	if err == nil || errors.Is(err, ErrNoChannelTag) {
		t.Errorf("expected a registry error distinct from ErrNoChannelTag, got %v", err)
	}
}

func TestChannelOf(t *testing.T) {
	cases := map[string]string{
		"1.0.0-beta.2":  "beta",
		"v1.0.0-beta.2": "beta",
		"1.0.0-alpha.3": "alpha",
		"1.0.0":         "latest",
		"v2.1.0":        "latest",
	}
	for in, want := range cases {
		if got := ChannelOf(in); got != want {
			t.Errorf("ChannelOf(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDetectNPMInstall(t *testing.T) {
	npmPaths := []string{
		"/usr/lib/node_modules/@jeancodogno/specforce-kit-linux-x64/bin/specforce",
		"/home/u/.nvm/versions/node/v20/lib/node_modules/@jeancodogno/specforce-kit-darwin-arm64/bin/specforce",
	}
	for _, p := range npmPaths {
		if !IsNPMInstallPath(p) {
			t.Errorf("expected %q to be detected as npm", p)
		}
	}
	binaryPaths := []string{
		"/usr/local/bin/specforce",
		"/home/node/bin/specforce",
		"/opt/npm-tools/specforce",
		"/usr/lib/node_modules/other-package/bin/specforce",
		"/usr/lib/node_modules/@jeancodogno/specforce-kit/bin/specforce",
	}
	for _, p := range binaryPaths {
		if IsNPMInstallPath(p) {
			t.Errorf("expected %q to be detected as a binary install", p)
		}
	}
}

func TestSelectProvider(t *testing.T) {
	p := SelectProvider("/usr/lib/node_modules/@jeancodogno/specforce-kit-linux-x64/bin/specforce", "1.0.0-beta.2")
	np, ok := p.(*NPMProvider)
	if !ok {
		t.Fatalf("expected *NPMProvider for an npm package path, got %T", p)
	}
	if np.Channel != "beta" {
		t.Errorf("expected channel beta, got %q", np.Channel)
	}
	for _, other := range []string{"/usr/local/bin/specforce", "/home/node/npm/specforce"} {
		if _, ok := SelectProvider(other, "1.0.0-beta.2").(*GitHubProvider); !ok {
			t.Errorf("expected *GitHubProvider for %q", other)
		}
	}
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

// newCountingDistTagsRegistry is newDistTagsRegistry that also counts requests.
func newCountingDistTagsRegistry(t *testing.T, body string, status int) (*httptest.Server, *int) {
	t.Helper()
	requests := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Path != "/-/package/@jeancodogno/specforce-kit/dist-tags" {
			t.Errorf("expected dist-tags path, got %s", r.URL.Path)
		}
		w.WriteHeader(status)
		_, _ = fmt.Fprint(w, body)
	}))
	t.Cleanup(ts.Close)
	return ts, &requests
}

func TestNPMProvider_StableExposedWithChannelVersion(t *testing.T) {
	ts, requests := newCountingDistTagsRegistry(t, `{"beta": "1.0.0-beta.3", "latest": "1.0.0"}`, http.StatusOK)
	p := &NPMProvider{BaseURL: ts.URL, Client: ts.Client(), Channel: "beta"}
	got, err := p.GetLatestVersion(context.Background())
	if err != nil {
		t.Fatalf("GetLatestVersion failed: %v", err)
	}
	if got != "1.0.0-beta.3" {
		t.Errorf("expected channel version 1.0.0-beta.3, got %q", got)
	}
	if stable := p.LatestStable(); stable != "1.0.0" {
		t.Errorf("expected stable 1.0.0, got %q", stable)
	}
	if *requests != 1 {
		t.Errorf("expected exactly one registry request, got %d", *requests)
	}
}

func TestNPMProvider_StableExposedWhenChannelTagMissing(t *testing.T) {
	ts, requests := newCountingDistTagsRegistry(t, `{"latest": "1.0.0"}`, http.StatusOK)
	p := &NPMProvider{BaseURL: ts.URL, Client: ts.Client(), Channel: "beta"}
	_, err := p.GetLatestVersion(context.Background())
	if !errors.Is(err, ErrNoChannelTag) {
		t.Errorf("expected ErrNoChannelTag, got %v", err)
	}
	if stable := p.LatestStable(); stable != "1.0.0" {
		t.Errorf("expected stable 1.0.0, got %q", stable)
	}
	if *requests != 1 {
		t.Errorf("expected exactly one registry request, got %d", *requests)
	}
}

func TestNPMProvider_StableNotExposedOnRegistryError(t *testing.T) {
	good, _ := newCountingDistTagsRegistry(t, `{"beta": "1.0.0-beta.3", "latest": "1.0.0"}`, http.StatusOK)
	bad, requests := newCountingDistTagsRegistry(t, `oops`, http.StatusInternalServerError)
	p := &NPMProvider{BaseURL: good.URL, Client: good.Client(), Channel: "beta"}
	if _, err := p.GetLatestVersion(context.Background()); err != nil {
		t.Fatalf("first check failed: %v", err)
	}

	// A later failed check must not expose the stale stable version.
	p.BaseURL, p.Client = bad.URL, bad.Client()
	if _, err := p.GetLatestVersion(context.Background()); err == nil || errors.Is(err, ErrNoChannelTag) {
		t.Fatalf("expected a registry error, got %v", err)
	}
	if stable := p.LatestStable(); stable != "" {
		t.Errorf("expected no stable version after a registry error, got %q", stable)
	}
	if *requests != 1 {
		t.Errorf("expected exactly one registry request, got %d", *requests)
	}
}
