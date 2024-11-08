package jtexpress

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test_account", "test_key")
	if client.apiAccount != "test_account" {
		t.Errorf("expected apiAccount to be test_account, got %s", client.apiAccount)
	}
}

func TestWithOptions(t *testing.T) {
	customTimeout := 5 * time.Second
	client := NewClient(
		"test_account",
		"test_key",
		WithHTTPClient(&http.Client{Timeout: customTimeout}),
		WithBaseURL("https://test.api.com"),
	)

	if client.httpClient.Timeout != customTimeout {
		t.Error("custom timeout was not set correctly")
	}

	if client.baseURL != "https://test.api.com" {
		t.Error("custom base URL was not set correctly")
	}
}
