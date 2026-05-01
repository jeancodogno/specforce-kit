package upgrade

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBinaryInstaller_DownloadAndVerify(t *testing.T) {
	binaryContent := "dummy binary content"
	h := sha256.New()
	h.Write([]byte(binaryContent))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	assetName := fmt.Sprintf("specforce-kit_%s_%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		assetName += ".exe"
	}
	checksumContent := fmt.Sprintf("%s  %s\n", expectedHash, assetName)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case fmt.Sprintf("/repos/jeancodogno/specforce-kit/releases/download/v1.0.0/%s", assetName):
			_, _ = w.Write([]byte(binaryContent))
		case "/repos/jeancodogno/specforce-kit/releases/download/v1.0.0/specforce-kit_1.0.0_checksums.txt":
			_, _ = w.Write([]byte(checksumContent))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	installer := &BinaryInstaller{
		Client: ts.Client(),
	}

	tmpFile, err := installer.DownloadAndVerify(context.Background(), "v1.0.0", ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != binaryContent {
		t.Errorf("expected content %q, got %q", binaryContent, string(content))
	}
}

func TestBinaryInstaller_ReplaceAt(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "replace-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	targetPath := filepath.Join(tempDir, "specforce")
	newPath := filepath.Join(tempDir, "specforce-new")

	if err := os.WriteFile(targetPath, []byte("old content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newPath, []byte("new content"), 0644); err != nil {
		t.Fatal(err)
	}

	installer := &BinaryInstaller{}
	if err := installer.ReplaceAt(newPath, targetPath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify target has new content
	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "new content" {
		t.Errorf("expected new content, got %q", string(content))
	}

	// Verify backup is gone
	if _, err := os.Stat(targetPath + ".old"); !os.IsNotExist(err) {
		t.Error("expected backup file to be removed")
	}
}

func TestBinaryInstaller_ReplaceAt_Rollback(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "rollback-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	targetPath := filepath.Join(tempDir, "specforce")
	newPath := "/non/existent/path" // This will cause moveFile to fail

	if err := os.WriteFile(targetPath, []byte("original content"), 0644); err != nil {
		t.Fatal(err)
	}

	installer := &BinaryInstaller{}
	err = installer.ReplaceAt(newPath, targetPath)
	if err == nil {
		t.Error("expected error during replacement, got nil")
	}

	// Verify target still has original content (rollback)
	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "original content" {
		t.Errorf("expected original content after rollback, got %q", string(content))
	}
}
