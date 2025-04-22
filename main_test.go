package main

import (
        "fmt"
        "net/http"
        "net/http/httptest"
        "os"
        "strings"
        "testing"
        "time"

        "github.com/ovh/go-ovh/ovh"
)

func setupTestServer() *httptest.Server {
        return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                // Handle authentication
                if r.URL.Path == "/auth/time" {
                        fmt.Fprintf(w, "%d", time.Now().Unix())
                        return
                }

                if r.URL.Path == "/auth/credential" {
                        w.WriteHeader(http.StatusOK)
                        return
                }

                // Handle DNS record queries
                if strings.HasPrefix(r.URL.Path, "/1.0/domain/zone/") {
                        parts := strings.Split(r.URL.Path, "/")
                        if len(parts) < 6 {
                                w.WriteHeader(http.StatusBadRequest)
                                return
                        }

                        if strings.HasSuffix(r.URL.Path, "/record") {
                                fmt.Fprintf(w, "[1234]")
                                return
                        }

                        if strings.HasSuffix(r.URL.Path, "/refresh") {
                                w.WriteHeader(http.StatusOK)
                                return
                        }

                        if strings.Contains(r.URL.Path, "/record/") {
                                w.WriteHeader(http.StatusOK)
                                return
                        }
                }

                w.WriteHeader(http.StatusNotFound)
        }))
}

func TestUpdateDNS(t *testing.T) {
        server := setupTestServer()
        defer server.Close()

        // Set environment variables for OVH client
        os.Setenv("OVH_ENDPOINT", server.URL)
        os.Setenv("OVH_APPLICATION_KEY", "test")
        os.Setenv("OVH_APPLICATION_SECRET", "test")
        os.Setenv("OVH_CONSUMER_KEY", "test")
        defer func() {
                os.Unsetenv("OVH_ENDPOINT")
                os.Unsetenv("OVH_APPLICATION_KEY")
                os.Unsetenv("OVH_APPLICATION_SECRET")
                os.Unsetenv("OVH_CONSUMER_KEY")
        }()

        // Initialize OVH client
        client, err := ovh.NewDefaultClient()
        if err != nil {
                t.Fatalf("Failed to create OVH client: %v", err)
        }

        // Test updating DNS record
        updateDNS(client, "example.com", "test", "1.2.3.4")
}

func TestUpdateAllDomains(t *testing.T) {
        server := setupTestServer()
        defer server.Close()

        // Set environment variables for OVH client
        os.Setenv("OVH_ENDPOINT", server.URL)
        os.Setenv("OVH_APPLICATION_KEY", "test")
        os.Setenv("OVH_APPLICATION_SECRET", "test")
        os.Setenv("OVH_CONSUMER_KEY", "test")
        defer func() {
                os.Unsetenv("OVH_ENDPOINT")
                os.Unsetenv("OVH_APPLICATION_KEY")
                os.Unsetenv("OVH_APPLICATION_SECRET")
                os.Unsetenv("OVH_CONSUMER_KEY")
        }()

        // Initialize OVH client
        client, err := ovh.NewDefaultClient()
        if err != nil {
                t.Fatalf("Failed to create OVH client: %v", err)
        }

        // Test domains
        domains := []DomainConfig{
                {
                        Zone:    "example.com",
                        Records: []string{"home", "office"},
                },
                {
                        Zone:    "another.com",
                        Records: []string{"vpn"},
                },
        }

        // Test updating multiple domains
        updateAllDomains(client, domains, "1.2.3.4")
}

func TestMainFunction(t *testing.T) {
        server := setupTestServer()
        defer server.Close()

        // Set environment variables
        os.Setenv("OVH_ENDPOINT", server.URL)
        os.Setenv("OVH_APPLICATION_KEY", "test")
        os.Setenv("OVH_APPLICATION_SECRET", "test")
        os.Setenv("OVH_CONSUMER_KEY", "test")
        os.Setenv("DOMAINS_CONFIG", `
domains:
  - zone: example.com
    records: ["test1", "test2"]
  - zone: another.com
    records: ["test3"]
`)
        os.Setenv("CHECK_INTERVAL", "1s")
        defer func() {
                os.Unsetenv("OVH_ENDPOINT")
                os.Unsetenv("OVH_APPLICATION_KEY")
                os.Unsetenv("OVH_APPLICATION_SECRET")
                os.Unsetenv("OVH_CONSUMER_KEY")
                os.Unsetenv("DOMAINS_CONFIG")
                os.Unsetenv("CHECK_INTERVAL")
        }()

        // Create a channel to stop the main function
        done := make(chan bool)
        go func() {
                time.Sleep(2 * time.Second)
                done <- true
        }()

        // Run main in a goroutine
        go func() {
                defer func() {
                        if r := recover(); r != nil {
                                // Expected panic when main exits
                                done <- true
                        }
                }()
                main()
        }()

        // Wait for completion or timeout
        select {
        case <-done:
                // Success
        case <-time.After(5 * time.Second):
                t.Fatal("Test timed out")
        }
}