# OVH DNS Updater

A simple Go application that monitors your public IPv4 address and updates DNS records in OVH when it changes.

## Features

- Monitors public IPv4 address changes
- Updates multiple OVH DNS A records automatically
- Configurable check interval
- Kubernetes deployment ready
- Supports multiple domains and records

## Configuration

The application is configured using environment variables:

- `DOMAINS_CONFIG`: YAML configuration for domains and records (see example below)
- `CHECK_INTERVAL`: How often to check for IP changes (default: "5m")
- `OVH_APPLICATION_KEY`: Your OVH API application key
- `OVH_APPLICATION_SECRET`: Your OVH API application secret
- `OVH_CONSUMER_KEY`: Your OVH API consumer key
- `OVH_ENDPOINT`: Your OVH API endpoint (e.g., "ovh-eu", "ovh-us", "ovh-ca")

### Obtaining OVH API Keys

To use this application, you need to obtain API credentials from OVH. Follow these steps:

1. Go to the OVH API key creation page for your region:
   - Europe: https://eu.api.ovh.com/createToken/
   - US: https://api.us.ovhcloud.com/createToken/
   - Canada: https://ca.api.ovh.com/createToken/

2. Log in with your OVH account credentials.

3. Fill in the following information:
   - Application name (e.g., "OVH DNS Updater")
   - Application description (optional)
   - Validity period (choose based on your needs)

4. For API access rights, add the following permissions:
   - GET /domain/zone/*
   - PUT /domain/zone/*
   - POST /domain/zone/*

5. Click "Create keys" to generate your credentials.

6. You will receive three keys:
   - Application Key (AK)
   - Application Secret (AS)
   - Consumer Key (CK)

7. Save these keys securely as they will be used to configure the application.

8. Set the appropriate endpoint based on your OVH region:
   - Europe: "ovh-eu"
   - US: "ovh-us"
   - Canada: "ovh-ca"

Example domains configuration:
```yaml
domains:
  - zone: example.com
    records:
      - home     # Updates home.example.com
      - office   # Updates office.example.com
      - "@"      # Updates example.com (apex domain)
  - zone: another.com
    records:
      - vpn      # Updates vpn.another.com
      - ""       # Updates another.com (apex domain)
```

Note: To update the apex domain (the domain itself without any subdomain), you can use either `"@"` or an empty string `""` in the records list.

## Running with Docker

```bash
docker run -e DOMAINS_CONFIG='domains:
  - zone: example.com
    records:
      - home     # Updates home.example.com
      - office   # Updates office.example.com
      - "@"      # Updates example.com itself
  - zone: another.com
    records:
      - ""       # Updates another.com itself' \
           -e OVH_APPLICATION_KEY=your_app_key \
           -e OVH_APPLICATION_SECRET=your_app_secret \
           -e OVH_CONSUMER_KEY=your_consumer_key \
           -e OVH_ENDPOINT=ovh-eu \
           ghcr.io/mangelajo/ovh-dns-updater:latest
```

## Running in Kubernetes

1. Create a secret with your OVH credentials:
```bash
kubectl create secret generic ovh-dns-updater-secret \
    --from-literal=OVH_APPLICATION_KEY=your_app_key \
    --from-literal=OVH_APPLICATION_SECRET=your_app_secret \
    --from-literal=OVH_CONSUMER_KEY=your_consumer_key \
    --from-literal=OVH_ENDPOINT=ovh-eu
```

2. Create a ConfigMap with your domains configuration:
```bash
kubectl create configmap ovh-dns-updater-config \
    --from-literal=domains.yaml='domains:
  - zone: example.com
    records:
      - home     # Updates home.example.com
      - office   # Updates office.example.com
      - "@"      # Updates example.com itself
  - zone: another.com
    records:
      - vpn      # Updates vpn.another.com
      - ""       # Updates another.com itself'
```

3. Deploy the application:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ovh-dns-updater
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ovh-dns-updater
  template:
    metadata:
      labels:
        app: ovh-dns-updater
    spec:
      containers:
      - name: ovh-dns-updater
        image: ghcr.io/mangelajo/ovh-dns-updater:latest
        env:
        - name: CHECK_INTERVAL
          value: "5m"
        - name: DOMAINS_CONFIG
          valueFrom:
            configMapKeyRef:
              name: ovh-dns-updater-config
              key: domains.yaml
        envFrom:
        - secretRef:
            name: ovh-dns-updater-secret
```

## Development

### Building

```bash
go build
```

### Running Tests

```bash
go test ./...
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -am 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.