package upgrade

import (
	"context"
	"fmt"
	"net/http"
)

// GitHubProvider fetches the latest version from GitHub Releases.
type GitHubProvider struct {
	BaseURL string
	Client  *http.Client
}

// NewGitHubProvider creates a new GitHubProvider with default configuration.
func NewGitHubProvider() *GitHubProvider {
	return &GitHubProvider{
		BaseURL: "https://api.github.com",
		Client:  NewHTTPClient(),
	}
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

// GetLatestVersion fetches the latest tag name from GitHub.
func (p *GitHubProvider) GetLatestVersion(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/repos/jeancodogno/specforce-kit/releases/latest", p.BaseURL)
	var release githubRelease
	if err := FetchJSON(ctx, p.Client, url, &release); err != nil {
		return "", err
	}
	return release.TagName, nil
}
