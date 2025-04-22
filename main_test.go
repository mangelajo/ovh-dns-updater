package main

import (
        "os"
        "testing"
        "time"
)

func TestLoadConfig(t *testing.T) {
        testCases := []struct {
                name     string
                vars     map[string]string
                wantErr  bool
                checkDNS bool
        }{
                {
                        name: "with subdomain",
                        vars: map[string]string{
                                "OVH_APPLICATION_KEY":    "test-key",
                                "OVH_APPLICATION_SECRET": "test-secret",
                                "OVH_CONSUMER_KEY":      "test-consumer",
                                "DNS_ZONE":              "example.com",
                                "DNS_RECORD":            "test",
                                "CHECK_INTERVAL":        "1m",
                        },
                        checkDNS: true,
                },
                {
                        name: "without subdomain",
                        vars: map[string]string{
                                "OVH_APPLICATION_KEY":    "test-key",
                                "OVH_APPLICATION_SECRET": "test-secret",
                                "OVH_CONSUMER_KEY":      "test-consumer",
                                "DNS_ZONE":              "example.com",
                                "CHECK_INTERVAL":        "1m",
                        },
                        checkDNS: true,
                },
                {
                        name: "missing required var",
                        vars: map[string]string{
                                "OVH_APPLICATION_KEY": "test-key",
                                "DNS_ZONE":           "example.com",
                        },
                        wantErr: true,
                },
        }

        for _, tc := range testCases {
                t.Run(tc.name, func(t *testing.T) {
                        // Clear environment
                        os.Clearenv()

                        // Set environment variables
                        for k, v := range tc.vars {
                                os.Setenv(k, v)
                        }

                        config, err := loadConfig()
                        
                        if tc.wantErr {
                                if err == nil {
                                        t.Fatal("Expected error, got nil")
                                }
                                return
                        }

                        if err != nil {
                                t.Fatalf("Expected no error, got %v", err)
                        }

                        // Verify config values
                        if config.OVHApplicationKey != tc.vars["OVH_APPLICATION_KEY"] {
                                t.Errorf("Expected application key %s, got %s", tc.vars["OVH_APPLICATION_KEY"], config.OVHApplicationKey)
                        }

                        if config.OVHApplicationSecret != tc.vars["OVH_APPLICATION_SECRET"] {
                                t.Errorf("Expected application secret %s, got %s", tc.vars["OVH_APPLICATION_SECRET"], config.OVHApplicationSecret)
                        }

                        if config.OVHConsumerKey != tc.vars["OVH_CONSUMER_KEY"] {
                                t.Errorf("Expected consumer key %s, got %s", tc.vars["OVH_CONSUMER_KEY"], config.OVHConsumerKey)
                        }

                        if config.DNSZone != tc.vars["DNS_ZONE"] {
                                t.Errorf("Expected DNS zone %s, got %s", tc.vars["DNS_ZONE"], config.DNSZone)
                        }

                        if tc.checkDNS {
                                expectedRecord := tc.vars["DNS_RECORD"]
                                if expectedRecord == "" {
                                        expectedRecord = ""
                                }
                                if config.DNSRecord != expectedRecord {
                                        t.Errorf("Expected DNS record %s, got %s", expectedRecord, config.DNSRecord)
                                }
                        }

                        if interval := tc.vars["CHECK_INTERVAL"]; interval != "" {
                                expectedDuration, _ := time.ParseDuration(interval)
                                if config.CheckInterval != expectedDuration {
                                        t.Errorf("Expected check interval %v, got %v", expectedDuration, config.CheckInterval)
                                }
                        }
                })
        }
}

func TestLoadConfigMissingRequired(t *testing.T) {
        // Set test mode
        os.Setenv("TEST_MODE", "1")
        defer os.Unsetenv("TEST_MODE")
        
        // Clear all relevant environment variables
        requiredVars := []string{
                "OVH_APPLICATION_KEY",
                "OVH_APPLICATION_SECRET",
                "OVH_CONSUMER_KEY",
                "DNS_ZONE",
        }

        // First test: missing required variables should fail
        for _, v := range requiredVars {
                t.Run("missing "+v, func(t *testing.T) {
                        // Clear environment
                        os.Clearenv()
                        os.Setenv("TEST_MODE", "1")

                        // Set all variables except one
                        for _, other := range requiredVars {
                                if other != v {
                                        os.Setenv(other, "test-value")
                                }
                        }

                        _, err := loadConfig()
                        if err == nil {
                                t.Errorf("Expected error when %s is missing", v)
                        }
                })
        }

        // Second test: missing DNS_RECORD should not fail
        t.Run("missing DNS_RECORD", func(t *testing.T) {
                // Clear environment
                os.Clearenv()
                os.Setenv("TEST_MODE", "1")

                // Set all required variables
                for _, v := range requiredVars {
                        os.Setenv(v, "test-value")
                }

                _, err := loadConfig()
                if err != nil {
                        t.Errorf("Expected no error when DNS_RECORD is missing, got: %v", err)
                }
        })
}

func TestGetCurrentIP(t *testing.T) {
        ip, err := getCurrentIP()
        if err != nil {
                t.Fatalf("Expected no error getting IP, got %v", err)
        }

        // Very basic IP format validation
        if len(ip) < 7 { // minimum valid IP length (e.g., "1.1.1.1")
                t.Errorf("IP address seems too short: %s", ip)
        }
}

func TestGetEnvWithDefault(t *testing.T) {
        testCases := []struct {
                name         string
                key          string
                value        string
                defaultValue string
                expected     string
        }{
                {
                        name:         "existing value",
                        key:          "TEST_KEY_1",
                        value:        "test-value",
                        defaultValue: "default",
                        expected:     "test-value",
                },
                {
                        name:         "use default",
                        key:          "TEST_KEY_2",
                        value:        "",
                        defaultValue: "default",
                        expected:     "default",
                },
        }

        for _, tc := range testCases {
                t.Run(tc.name, func(t *testing.T) {
                        if tc.value != "" {
                                os.Setenv(tc.key, tc.value)
                                defer os.Unsetenv(tc.key)
                        }

                        result := getEnvWithDefault(tc.key, tc.defaultValue)
                        if result != tc.expected {
                                t.Errorf("Expected %s, got %s", tc.expected, result)
                        }
                })
        }
}