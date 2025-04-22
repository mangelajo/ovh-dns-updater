FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ovh-dns-updater

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/ovh-dns-updater .
CMD ["./ovh-dns-updater"]