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

## Container Images

Pre-built container images are available on GitHub Container Registry:

```bash
# Latest from main branch
docker pull ghcr.io/mangelajo/ovh-dns-updater:main

# Specific version
docker pull ghcr.io/mangelajo/ovh-dns-updater:v1.0.0
```

The following architectures are supported:
- linux/amd64
- linux/arm64

## CI/CD Pipeline

This repository includes a GitHub Actions workflow that:
1. Runs tests on every PR and commit
2. Builds and pushes container images to GitHub Container Registry on:
   - Every push to main (tagged as `:main`)
   - Every tag (tagged as `:v1.0.0`, `:1.0`, etc.)
   - Every PR (tagged with PR number)

## Deployment

1. Create a Kubernetes secret with your OVH credentials:
   ```bash
   # Method 1: Create the secret directly (recommended)
   kubectl create secret generic ovh-dns-updater-secret \
     --from-literal=OVH_APPLICATION_KEY='your-app-key' \
     --from-literal=OVH_APPLICATION_SECRET='your-app-secret' \
     --from-literal=OVH_CONSUMER_KEY='your-consumer-key'

   # Method 2: Using YAML
   cat <<EOF | kubectl apply -f -
   apiVersion: v1
   kind: Secret
   metadata:
     name: ovh-dns-updater-secret
   type: Opaque
   stringData:
     OVH_APPLICATION_KEY: your-app-key
     OVH_APPLICATION_SECRET: your-app-secret
     OVH_CONSUMER_KEY: your-consumer-key
   EOF
   ```

2. Create a ConfigMap for your DNS configuration:
   ```bash
   kubectl create configmap ovh-dns-updater-config \
     --from-literal=DNS_ZONE='example.com' \
     --from-literal=DNS_RECORD='home' \
     --from-literal=CHECK_INTERVAL='5m'
   ```

3. Apply the Kubernetes deployment:
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