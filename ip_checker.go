package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type IPChecker struct {
	endpoints []string
}

func NewIPChecker() *IPChecker {
	return &IPChecker{
		endpoints: []string{
			"https://api.ipify.org",
			"https://ifconfig.me/ip",
			"https://icanhazip.com",
		},
	}
}

func (c *IPChecker) GetPublicIP() (string, error) {
	var lastErr error
	for _, endpoint := range c.endpoints {
		ip, err := c.getIP(endpoint)
		if err != nil {
			lastErr = err
			continue
		}
		return ip, nil
	}
	return "", fmt.Errorf("failed to get public IP from any endpoint: %v", lastErr)
}

func (c *IPChecker) getIP(endpoint string) (string, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := strings.TrimSpace(string(body))
	return ip, nil
}