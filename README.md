# OVH DNS Updater

This service monitors your public IPv4 address and automatically updates a specified DNS record in OVH when the IP changes. It's designed to run in Kubernetes and uses environment variables for configuration.

## Features

- Monitors public IPv4 address changes
- Updates OVH DNS A records automatically
- Configurable check interval
- Kubernetes-ready deployment

## Prerequisites

- Kubernetes cluster
- OVH API credentials
- Docker (for building the image)

## Configuration

The following environment variables are required:

- `OVH_APPLICATION_KEY`: Your OVH API application key
- `OVH_APPLICATION_SECRET`: Your OVH API application secret
- `OVH_CONSUMER_KEY`: Your OVH API consumer key
- `DNS_ZONE`: Your domain name (e.g., "example.com")

Optional environment variables:

- `DNS_RECORD`: The subdomain to update (e.g., "home"). If not provided, updates the base domain's A record
- `OVH_ENDPOINT`: OVH API endpoint (default: "ovh-eu")
- `CHECK_INTERVAL`: How often to check for IP changes (default: "5m")

## Deployment

1. First, build the Docker image:
   ```bash
   docker build -t ovh-dns-updater:latest .
   ```

2. Create a secret with your OVH credentials:
   ```bash
   # Convert your credentials to base64
   echo -n "your-app-key" | base64
   echo -n "your-app-secret" | base64
   echo -n "your-consumer-key" | base64

   # Update the values in k8s/deployment.yaml with the base64 encoded values
   ```

3. Update the DNS_ZONE and DNS_RECORD values in k8s/deployment.yaml

4. Apply the Kubernetes configuration:
   ```bash
   kubectl apply -f k8s/deployment.yaml
   ```

## Getting OVH API Credentials

1. Go to https://api.ovh.com/createToken/
2. Create API credentials with the following rights:
   - GET /domain/zone/*
   - PUT /domain/zone/*
   - POST /domain/zone/*
3. Save the generated credentials

## Monitoring

You can check the logs of the running pod:
```bash
kubectl logs -f deployment/ovh-dns-updater
```