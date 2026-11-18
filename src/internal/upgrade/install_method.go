package upgrade

import (
	"path/filepath"
	"regexp"
)

// npmPlatformPackageDir matches the Specforce npm platform package directory,
// e.g. node_modules/@jeancodogno/specforce-kit-linux-x64/.
var npmPlatformPackageDir = regexp.MustCompile(`(^|/)node_modules/@jeancodogno/specforce-kit-[^/]+-[^/]+/`)

// IsNPMInstallPath reports whether exePath lives inside the Specforce npm
// platform package. Any other location is a native binary installation.
func IsNPMInstallPath(exePath string) bool {
	return npmPlatformPackageDir.MatchString(filepath.ToSlash(exePath))
}

// SelectProvider returns the version source for the installation at exePath:
// the npm registry channel of currentVersion for npm installs, GitHub otherwise.
func SelectProvider(exePath, currentVersion string) Provider {
	if IsNPMInstallPath(resolvePath(exePath)) {
		p := NewNPMProvider()
		p.Channel = ChannelOf(currentVersion)
		return p
	}
	return NewGitHubProvider()
}
