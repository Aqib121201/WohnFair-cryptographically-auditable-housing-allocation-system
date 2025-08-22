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

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notifier ./cmd/notifier

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata curl

# Create non-root user
RUN addgroup -g 1001 wohnfair && adduser -u 1001 -G wohnfair -s /bin/sh -D wohnfair

# Copy binary from builder stage
COPY --from=builder /app/notifier /app/notifier

# Set ownership
RUN chown -R wohnfair:wohnfair /app

# Switch to non-root user
USER wohnfair

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD ["curl", "-f", "http://localhost:8080/healthz"] || exit 1

# Run the service
CMD ["./notifier"]
