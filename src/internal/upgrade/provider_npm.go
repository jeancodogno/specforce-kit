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
	// latestStable is the "latest" dist-tag read by the last GetLatestVersion
	// call; empty when that call failed to fetch the dist-tags.
	latestStable string
}

// StableProvider is implemented by providers that also expose the stable
// version read during their last GetLatestVersion call (from the same response).
type StableProvider interface {
	// LatestStable returns the stable version of the last check, or "" if unknown.
	LatestStable() string
}

// LatestStable returns the "latest" dist-tag read by the last GetLatestVersion
// call, or "" when that call could not fetch the dist-tags.
func (p *NPMProvider) LatestStable() string {
	return p.latestStable
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
// The "latest" tag of the same response is kept for LatestStable.
func (p *NPMProvider) GetLatestVersion(ctx context.Context) (string, error) {
	p.latestStable = ""
	channel := p.Channel
	if channel == "" {
		channel = defaultChannel
	}
	url := fmt.Sprintf("%s/-/package/%s/dist-tags", p.BaseURL, npmPackage)
	tags := map[string]string{}
	if err := FetchJSON(ctx, p.Client, url, &tags); err != nil {
		return "", fmt.Errorf("failed to fetch npm dist-tags: %w", err)
	}
	p.latestStable = tags[defaultChannel]
	version, ok := tags[channel]
	if !ok || version == "" {
		return "", fmt.Errorf("%w: %s", ErrNoChannelTag, channel)
	}
	return version, nil
}
