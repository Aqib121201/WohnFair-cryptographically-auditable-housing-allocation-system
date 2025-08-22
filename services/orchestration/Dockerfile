# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/gateway ./cmd/gateway

# Runtime
FROM alpine:3.19
RUN adduser -D -H -s /sbin/nologin wohnfair
WORKDIR /app
COPY --from=builder /out/gateway /app/gateway
USER wohnfair
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=10s --start-period=20s --retries=3 CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1
CMD ["/app/gateway"]
