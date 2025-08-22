# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fairrentd ./cmd/fairrentd

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S wohnfair && \
    adduser -u 1001 -S wohnfair -G wohnfair

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/fairrentd .

# Copy configuration files
COPY --from=builder /app/config ./config

# Set ownership
RUN chown -R wohnfair:wohnfair /app

# Switch to non-root user
USER wohnfair

# Expose port
EXPOSE 50051

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD grpc_health_probe -addr=localhost:50051 || exit 1

# Run the binary
CMD ["./fairrentd"]
