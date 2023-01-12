FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" .

FROM debian:latest
WORKDIR /app
COPY --from=builder /app/caddy-logger /app/caddy-logger
ENTRYPOINT ["/app/caddy-logger"]
