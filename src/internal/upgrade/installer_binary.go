package upgrade

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

// BinaryInstaller handles upgrades by downloading the native binary.
type BinaryInstaller struct {
	Client *http.Client
}

// NewBinaryInstaller creates a new BinaryInstaller.
func NewBinaryInstaller() *BinaryInstaller {
	return &BinaryInstaller{
		Client: NewHTTPClient(),
	}
}

// DownloadAndVerify downloads the binary and verifies its checksum.
func (i *BinaryInstaller) DownloadAndVerify(ctx context.Context, version, baseURL string) (string, error) {
	osName := runtime.GOOS
	archName := runtime.GOARCH
	
	// Format asset name: specforce-kit_linux_amd64
	assetName := fmt.Sprintf("specforce-kit_%s_%s", osName, archName)
	if osName == "windows" {
		assetName += ".exe"
	}

	// Format checksum file name: specforce-kit_0.2.2_checksums.txt
	// version might have 'v' prefix, need to strip it for checksum filename if that's the pattern
	versionClean := strings.TrimPrefix(version, "v")
	checksumName := fmt.Sprintf("specforce-kit_%s_checksums.txt", versionClean)

	// In a real scenario, we'd find the asset URL from the GitHub API response.
	// For simplicity in this implementation, we construct it assuming GitHub's download pattern
	// or we pass the download URL directly if we had it.
	// Since we are implementation, let's assume we have a helper to get the asset URL.
	
	downloadURL := fmt.Sprintf("%s/repos/jeancodogno/specforce-kit/releases/download/%s/%s", 
		baseURL, version, assetName)
	checksumURL := fmt.Sprintf("%s/repos/jeancodogno/specforce-kit/releases/download/%s/%s", 
		baseURL, version, checksumName)

	// 1. Download Binary to Temp File
	tmpBinary, err := os.CreateTemp("", "specforce-update-*")
	if err != nil {
		return "", err
	}
	defer func() { _ = tmpBinary.Close() }()

	if err := i.downloadFile(ctx, downloadURL, tmpBinary); err != nil {
		return "", fmt.Errorf("failed to download binary: %w", err)
	}

	// 2. Download Checksum
	resp, err := i.Client.Get(checksumURL)
	if err != nil {
		return "", fmt.Errorf("failed to download checksum: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("checksum file not found (status %d)", resp.StatusCode)
	}

	expectedHash, err := i.findHashInChecksums(resp.Body, assetName)
	if err != nil {
		return "", err
	}

	// 3. Verify Hash
	if err := i.verifyHash(tmpBinary.Name(), expectedHash); err != nil {
		_ = os.Remove(tmpBinary.Name())
		return "", err
	}

	return tmpBinary.Name(), nil
}

func (i *BinaryInstaller) downloadFile(ctx context.Context, url string, w io.Writer) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := i.Client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file (status %d)", resp.StatusCode)
	}

	_, err = io.Copy(w, resp.Body)
	return err
}

func (i *BinaryInstaller) findHashInChecksums(r io.Reader, assetName string) (string, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[1] == assetName {
			return parts[0], nil
		}
	}
	return "", fmt.Errorf("hash for %s not found in checksum file", assetName)
}

func (i *BinaryInstaller) verifyHash(filePath, expectedHash string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	actualHash := hex.EncodeToString(h.Sum(nil))
	if actualHash != expectedHash {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedHash, actualHash)
	}

	return nil
}

// Replace replaces the currently running binary with the new one at newPath.
func (i *BinaryInstaller) Replace(newPath string) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	return i.ReplaceAt(newPath, exePath)
}

// ReplaceAt replaces the binary at targetPath with the new one at newPath.
func (i *BinaryInstaller) ReplaceAt(newPath, targetPath string) error {
	// 1. Move current binary to backup
	backupPath := targetPath + ".old"
	if err := os.Rename(targetPath, backupPath); err != nil {
		return fmt.Errorf("failed to move current binary to backup: %w", err)
	}

	// 2. Move new binary to target path
	if err := i.moveFile(newPath, targetPath); err != nil {
		// Rollback: move backup back to targetPath
		_ = os.Rename(backupPath, targetPath)
		return fmt.Errorf("failed to move new binary to target: %w", err)
	}

	// 3. Remove backup
	_ = os.Remove(backupPath)

	// 4. Ensure new binary is executable
	// #nosec G302 - Binary must be executable by others to function globally
	return os.Chmod(targetPath, 0755)
}

func (i *BinaryInstaller) moveFile(src, dst string) error {
	// Try renaming first (same filesystem)
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}

	// Fallback to copy + delete (different filesystems)
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	_ = os.Remove(src)
	return nil
}
