package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expands leading ~ to user home directory and environment variables in the form ${VAR:-DEFAULT}.
func ExpandPath(path string) string {
	// 1. Expand environment variables first: ${VAR:-DEFAULT}
	path = expandEnv(path)

	// 2. Expand tilde if it's now at the beginning
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			if path == "~" {
				path = home
			} else if len(path) > 1 && (path[1] == '/' || path[1] == filepath.Separator || path[1] == '.') {
				path = filepath.Join(home, path[1:])
			}
		}
	}

	return path
}

func expandEnv(path string) string {
	for {
		start := strings.Index(path, "${")
		if start == -1 {
			break
		}

		// Find the matching closing brace, handling nested ${}
		balance := 0
		end := -1
		for i := start; i < len(path); i++ {
			if strings.HasPrefix(path[i:], "${") {
				balance++
				i++ // skip '{'
			} else if path[i] == '}' {
				balance--
				if balance == 0 {
					end = i
					break
				}
			}
		}

		if end == -1 {
			break // Malformed
		}

		content := path[start+2 : end]
		replacement := ""

		if idx := strings.Index(content, ":-"); idx != -1 {
			varName := content[:idx]
			fallback := content[idx+2:]
			val := os.Getenv(varName)
			if val != "" {
				replacement = val
			} else {
				replacement = expandEnv(fallback) // Recursive fallback expansion
			}
		} else {
			replacement = os.Getenv(content)
		}

		path = path[:start] + replacement + path[end+1:]
	}
	return path
}

// SecurePath takes a root and a target path, cleans it, and ensures it's within the root.
// In Go 1.24+, we use os.OpenRoot to validate that the path is within boundaries if the root exists.
func SecurePath(root, target string) (string, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for root %s: %w", root, err)
	}

	// Clean the target path first to handle relative segments
	cleanTarget := filepath.Clean(target)
	if filepath.IsAbs(cleanTarget) {
		rel, err := filepath.Rel(absRoot, cleanTarget)
		if err != nil || strings.HasPrefix(rel, "..") {
			return "", fmt.Errorf("security: absolute path %s is outside root %s", target, absRoot)
		}
		cleanTarget = rel
	}

	r, err := os.OpenRoot(absRoot)
	if err == nil {
		defer func() { _ = r.Close() }()
		// os.OpenRoot.Stat will fail if path escapes root.
		_, err = r.Stat(cleanTarget)
		if err != nil && !os.IsNotExist(err) {
			return "", fmt.Errorf("security: path traversal attempt or invalid path: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to open root %s: %w", absRoot, err)
	}

	// Fallback/Double check: Ensure the final path still has the absRoot prefix
	finalPath := filepath.Join(absRoot, cleanTarget)
	if !strings.HasPrefix(finalPath, absRoot) {
		return "", fmt.Errorf("security: path traversal attempt detected for %s", target)
	}

	return finalPath, nil
}
