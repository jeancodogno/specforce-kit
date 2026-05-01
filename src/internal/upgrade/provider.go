package upgrade

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Provider defines the interface for fetching the latest version from a remote source.
type Provider interface {
	// GetLatestVersion fetches the latest available version string (e.g., "v1.2.3").
	GetLatestVersion(ctx context.Context) (string, error)
}

// FetchJSON performs a GET request and decodes the JSON response into v.
func FetchJSON(ctx context.Context, client *http.Client, url string, v interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// NewHTTPClient returns a pre-configured HTTP client for providers.
func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 2 * time.Second,
	}
}

// MockProvider is a provider used for testing.
type MockProvider struct {
	Version string
	Err     error
}

// GetLatestVersion returns the pre-configured version or error.
func (p *MockProvider) GetLatestVersion(ctx context.Context) (string, error) {
	if p.Err != nil {
		return "", p.Err
	}
	return p.Version, nil
}
