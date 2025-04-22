package main

import (
	"net"
	"testing"
)

func TestGetCurrentIPFormat(t *testing.T) {
	ip, err := getCurrentIP()
	if err != nil {
		t.Fatalf("Failed to get IP: %v", err)
	}

	// Parse the IP address
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		t.Errorf("Invalid IP address format: %s", ip)
	}

	// Check if it's an IPv4 address
	if parsedIP.To4() == nil {
		t.Errorf("Not an IPv4 address: %s", ip)
	}
}

func TestGetCurrentIPConsistency(t *testing.T) {
	// Get IP twice and compare
	ip1, err := getCurrentIP()
	if err != nil {
		t.Fatalf("Failed to get first IP: %v", err)
	}

	ip2, err := getCurrentIP()
	if err != nil {
		t.Fatalf("Failed to get second IP: %v", err)
	}

	if ip1 != ip2 {
		t.Errorf("Inconsistent IP addresses: %s != %s", ip1, ip2)
	}
}