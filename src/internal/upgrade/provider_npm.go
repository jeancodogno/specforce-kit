package upgrade

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/mod/semver"
)

// ErrNoChannelTag is returned (wrapped) by NPMProvider when the registry has no
// dist-tag for the running version's channel. It means "no update", not a failure.
var ErrNoChannelTag = errors.New("no npm dist-tag for release channel")

// defaultChannel is the npm dist-tag of stable releases.
const defaultChannel = "latest"

// NPMProvider fetches the latest version of a release channel from the NPM Registry.
type NPMProvider struct {
	BaseURL string
	Client  *http.Client
	// Channel is the npm dist-tag to follow (e.g. "beta"); empty means "latest".
	Channel string
}

// NewNPMProvider creates a new NPMProvider with default configuration.
func NewNPMProvider() *NPMProvider {
	return &NPMProvider{
		BaseURL: "https://registry.npmjs.org",
		Client:  NewHTTPClient(),
	}
}

// ChannelOf returns the npm dist-tag of a version's release channel: the first
// pre-release identifier (1.0.0-beta.2 -> beta), or "latest" for stable releases.
func ChannelOf(version string) string {
	pre := strings.TrimPrefix(semver.Prerelease(normalizeVersion(version)), "-")
	if pre == "" {
		return defaultChannel
	}
	return strings.SplitN(pre, ".", 2)[0]
}

// GetLatestVersion reads the registry dist-tags and returns the version of the
// provider's channel. A missing tag returns an error wrapping ErrNoChannelTag.
func (p *NPMProvider) GetLatestVersion(ctx context.Context) (string, error) {
	channel := p.Channel
	if channel == "" {
		channel = defaultChannel
	}
	url := fmt.Sprintf("%s/-/package/%s/dist-tags", p.BaseURL, npmPackage)
	tags := map[string]string{}
	if err := FetchJSON(ctx, p.Client, url, &tags); err != nil {
		return "", fmt.Errorf("failed to fetch npm dist-tags: %w", err)
	}
	version, ok := tags[channel]
	if !ok || version == "" {
		return "", fmt.Errorf("%w: %s", ErrNoChannelTag, channel)
	}
	return version, nil
}
