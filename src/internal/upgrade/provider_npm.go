package upgrade

import (
	"context"
	"fmt"
	"net/http"
)

// NPMProvider fetches the latest version from NPM Registry.
type NPMProvider struct {
	BaseURL string
	Client  *http.Client
}

// NewNPMProvider creates a new NPMProvider with default configuration.
func NewNPMProvider() *NPMProvider {
	return &NPMProvider{
		BaseURL: "https://registry.npmjs.org",
		Client:  NewHTTPClient(),
	}
}

type npmLatest struct {
	Version string `json:"version"`
}

// GetLatestVersion fetches the latest version from NPM.
func (p *NPMProvider) GetLatestVersion(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/@jeancodogno/specforce-kit/latest", p.BaseURL)
	var latest npmLatest
	if err := FetchJSON(ctx, p.Client, url, &latest); err != nil {
		return "", err
	}
	return latest.Version, nil
}
