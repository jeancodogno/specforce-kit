package upgrade

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHTTPClient(t *testing.T) {
	client := NewHTTPClient()
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.Timeout == 0 {
		t.Error("expected non-zero timeout")
	}
}

func TestFetchJSON_Errors(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	err := FetchJSON(context.Background(), ts.Client(), ts.URL, nil)
	if err == nil {
		t.Error("expected error for 500 status code, got nil")
	}

	err = FetchJSON(context.Background(), ts.Client(), "invalid-url", nil)
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestMockProvider(t *testing.T) {
	p_err := &MockProvider{Err: fmt.Errorf("error")}
	_, err := p_err.GetLatestVersion(context.Background())
	if err == nil {
		t.Error("expected error from MockProvider, got nil")
	}

	expected := "v2.0.0"
	p := &MockProvider{Version: expected}

	version, err := p.GetLatestVersion(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if version != expected {
		t.Errorf("expected version %s, got %s", expected, version)
	}
}
