package main

import (
        "fmt"
        "log"
        "os"
        "time"

        "github.com/ovh/go-ovh/ovh"
        "gopkg.in/yaml.v3"
)

type recordUpdate struct {
        Target string `json:"target"`
}

func updateDNS(client *ovh.Client, zone, record, ip string) {
        // Get the record ID
        var records []int
        err := client.Get(fmt.Sprintf("/domain/zone/%s/record?fieldType=A&subDomain=%s", zone, record), &records)
        if err != nil {
                log.Printf("Error getting DNS records for %s.%s: %v", record, zone, err)
                return
        }

        if len(records) == 0 {
                log.Printf("No DNS record found for %s.%s", record, zone)
                return
        }

        // Update the record
        update := recordUpdate{Target: ip}

        recordID := records[0]
        err = client.Put(fmt.Sprintf("/domain/zone/%s/record/%d", zone, recordID), &update, nil)
        if err != nil {
                log.Printf("Error updating DNS record %s.%s: %v", record, zone, err)
                return
        }

        // Refresh the zone
        err = client.Post(fmt.Sprintf("/domain/zone/%s/refresh", zone), nil, nil)
        if err != nil {
                log.Printf("Error refreshing DNS zone %s: %v", zone, err)
                return
        }

        log.Printf("Updated DNS record %s.%s to %s", record, zone, ip)
}

func updateAllDomains(client *ovh.Client, domains []DomainConfig, ip string) {
        for _, domain := range domains {
                for _, record := range domain.Records {
                        updateDNS(client, domain.Zone, record, ip)
                }
        }
}

func main() {
        // Load configuration from environment variables
        configYAML := os.Getenv("DOMAINS_CONFIG")
        if configYAML == "" {
                log.Fatal("DOMAINS_CONFIG environment variable is required")
        }

        var config Config
        err := yaml.Unmarshal([]byte(configYAML), &config)
        if err != nil {
                log.Fatalf("Error parsing DOMAINS_CONFIG: %v", err)
        }

        // Initialize OVH client
        client, err := ovh.NewDefaultClient()
        if err != nil {
                log.Fatalf("Error initializing OVH client: %v", err)
        }

        // Initialize IP checker
        ipChecker := NewIPChecker()

        // Get check interval from environment variable
        checkInterval := os.Getenv("CHECK_INTERVAL")
        if checkInterval == "" {
                checkInterval = "5m"
        }

        interval, err := time.ParseDuration(checkInterval)
        if err != nil {
                log.Fatalf("Invalid CHECK_INTERVAL: %v", err)
        }

        // Get initial IP
        currentIP, err := ipChecker.GetPublicIP()
        if err != nil {
                log.Printf("Error getting initial public IP: %v", err)
        } else {
                log.Printf("Initial public IP: %s", currentIP)
                updateAllDomains(client, config.Domains, currentIP)
        }

        // Start monitoring for IP changes
        ticker := time.NewTicker(interval)
        defer ticker.Stop()

        for range ticker.C {
                newIP, err := ipChecker.GetPublicIP()
                if err != nil {
                        log.Printf("Error getting public IP: %v", err)
                        continue
                }

                if newIP != currentIP {
                        log.Printf("Public IP changed from %s to %s", currentIP, newIP)
                        updateAllDomains(client, config.Domains, newIP)
                        currentIP = newIP
                }
        }
}