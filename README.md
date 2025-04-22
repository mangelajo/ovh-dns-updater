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

Example domains configuration:
```yaml
domains:
  - zone: example.com
    records:
      - home
      - office
  - zone: another.com
    records:
      - vpn
```

## Running with Docker

```bash
docker run -e DOMAINS_CONFIG='domains:
  - zone: example.com
    records:
      - home
      - office' \
           -e OVH_APPLICATION_KEY=your_app_key \
           -e OVH_APPLICATION_SECRET=your_app_secret \
           -e OVH_CONSUMER_KEY=your_consumer_key \
           ghcr.io/mangelajo/ovh-dns-updater:latest
```

## Running in Kubernetes

1. Create a secret with your OVH credentials:
```bash
kubectl create secret generic ovh-dns-updater-secret \
    --from-literal=OVH_APPLICATION_KEY=your_app_key \
    --from-literal=OVH_APPLICATION_SECRET=your_app_secret \
    --from-literal=OVH_CONSUMER_KEY=your_consumer_key
```

2. Create a ConfigMap with your domains configuration:
```bash
kubectl create configmap ovh-dns-updater-config \
    --from-literal=domains.yaml='domains:
  - zone: example.com
    records:
      - home
      - office
  - zone: another.com
    records:
      - vpn'
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