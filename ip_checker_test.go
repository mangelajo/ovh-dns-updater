package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPublicIP(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("1.2.3.4"))
	}))
	defer server.Close()

	// Create IP checker with test server
	checker := &IPChecker{
		endpoints: []string{server.URL},
	}

	// Test getting IP
	ip, err := checker.GetPublicIP()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if ip != "1.2.3.4" {
		t.Errorf("Expected IP 1.2.3.4, got %s", ip)
	}
}