package main

import (
        "fmt"
        "io"
        "log"
        "net/http"
        "os"
        "time"

        "github.com/ovh/go-ovh/ovh"
)

type Config struct {
        OVHApplicationKey    string
        OVHApplicationSecret string
        OVHConsumerKey      string
        OVHEndpoint         string
        DNSZone            string
        DNSRecord          string
        CheckInterval       time.Duration
}

func loadConfig() (*Config, error) {
        interval, err := time.ParseDuration(getEnvWithDefault("CHECK_INTERVAL", "5m"))
        if err != nil {
                return nil, fmt.Errorf("invalid CHECK_INTERVAL: %v", err)
        }

        requiredVars := []string{
                "OVH_APPLICATION_KEY",
                "OVH_APPLICATION_SECRET",
                "OVH_CONSUMER_KEY",
                "DNS_ZONE",
        }

        for _, v := range requiredVars {
                if os.Getenv(v) == "" {
                        return nil, fmt.Errorf("environment variable %s is required", v)
                }
        }

        return &Config{
                OVHApplicationKey:    os.Getenv("OVH_APPLICATION_KEY"),
                OVHApplicationSecret: os.Getenv("OVH_APPLICATION_SECRET"),
                OVHConsumerKey:      os.Getenv("OVH_CONSUMER_KEY"),
                OVHEndpoint:         getEnvWithDefault("OVH_ENDPOINT", "ovh-eu"),
                DNSZone:            os.Getenv("DNS_ZONE"),
                DNSRecord:          os.Getenv("DNS_RECORD"),
                CheckInterval:       interval,
        }, nil
}

func getEnvOrFatal(key string) string {
        if value := os.Getenv(key); value != "" {
                return value
        }
        if os.Getenv("TEST_MODE") != "" {
                return ""
        }
        log.Fatalf("Environment variable %s is required", key)
        return ""
}

func getEnvWithDefault(key, defaultValue string) string {
        if value := os.Getenv(key); value != "" {
                return value
        }
        return defaultValue
}

func getCurrentIP() (string, error) {
        resp, err := http.Get("https://api.ipify.org")
        if err != nil {
                return "", err
        }
        defer resp.Body.Close()

        ip, err := io.ReadAll(resp.Body)
        if err != nil {
                return "", err
        }

        return string(ip), nil
}

func updateDNSRecord(client *ovh.Client, config *Config, currentIP string) error {
        // Get the record ID first
        var records []int
        var err error
        
        // Build the query based on whether we have a subdomain
        query := fmt.Sprintf("/domain/zone/%s/record?fieldType=A", config.DNSZone)
        if config.DNSRecord != "" {
                query = fmt.Sprintf("%s&subDomain=%s", query, config.DNSRecord)
        }
        
        err = client.Get(query, &records)
        if err != nil {
                return fmt.Errorf("failed to get DNS records: %v", err)
        }

        if len(records) == 0 {
                return fmt.Errorf("no matching DNS record found")
        }

        // Update the record
        type RecordUpdate struct {
                Target string `json:"target"`
        }
        
        update := RecordUpdate{Target: currentIP}
        var result interface{}
        err = client.Put(fmt.Sprintf("/domain/zone/%s/record/%d", 
                config.DNSZone, records[0]), &update, &result)
        if err != nil {
                return fmt.Errorf("failed to update DNS record: %v", err)
        }

        // Refresh the zone
        err = client.Post(fmt.Sprintf("/domain/zone/%s/refresh", config.DNSZone), nil, &result)
        if err != nil {
                return fmt.Errorf("failed to refresh DNS zone: %v", err)
        }

        return nil
}

func main() {
        config, err := loadConfig()
        if err != nil {
                log.Fatalf("Failed to load configuration: %v", err)
        }

        client, err := ovh.NewClient(
                config.OVHEndpoint,
                config.OVHApplicationKey,
                config.OVHApplicationSecret,
                config.OVHConsumerKey,
        )
        if err != nil {
                log.Fatalf("Failed to create OVH client: %v", err)
        }

        var lastIP string
        for {
                currentIP, err := getCurrentIP()
                if err != nil {
                        log.Printf("Failed to get current IP: %v", err)
                        time.Sleep(config.CheckInterval)
                        continue
                }

                if currentIP != lastIP {
                        log.Printf("IP changed from %s to %s", lastIP, currentIP)
                        err = updateDNSRecord(client, config, currentIP)
                        if err != nil {
                                log.Printf("Failed to update DNS record: %v", err)
                        } else {
                                lastIP = currentIP
                                log.Printf("Successfully updated DNS record")
                        }
                }

                time.Sleep(config.CheckInterval)
        }
}