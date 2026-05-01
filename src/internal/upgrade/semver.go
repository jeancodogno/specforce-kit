package upgrade

import (
	"strings"

	"golang.org/x/mod/semver"
)

// CompareVersions compares two version strings and returns:
// -1 if v1 < v2
//  0 if v1 == v2
//  1 if v1 > v2
func CompareVersions(v1, v2 string) int {
	v1 = normalizeVersion(v1)
	v2 = normalizeVersion(v2)
	return semver.Compare(v1, v2)
}

// IsNewer returns true if latest is greater than current.
func IsNewer(current, latest string) bool {
	return CompareVersions(current, latest) == -1
}

// normalizeVersion ensures the version string has a 'v' prefix for semver.Compare.
func normalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if !strings.HasPrefix(v, "v") {
		return "v" + v
	}
	return v
}
