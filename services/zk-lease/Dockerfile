# Multi-stage build for ZK-Lease service
FROM rust:1.75-slim as builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /usr/src/zk-lease

# Copy manifest files
COPY Cargo.toml Cargo.lock ./

# Create a dummy main.rs to build dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs

# Build dependencies
RUN cargo build --release

# Remove dummy main.rs and copy source
RUN rm -rf src
COPY src ./src

# Build the application
RUN cargo build --release

# Runtime stage
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    libpq5 \
    libssl3 \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd --create-home --shell /bin/bash wohnfair

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /usr/src/zk-lease/target/release/zk-lease /app/zk-lease

# Copy configuration files
COPY config ./config

# Create necessary directories
RUN mkdir -p /app/circuits /app/keys /app/logs

# Install grpc_health_probe for healthchecks
RUN wget -qO /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.11/grpc_health_probe-linux-amd64 \
  && chmod +x /bin/grpc_health_probe

# Set ownership
RUN chown -R wohnfair:wohnfair /app

# Switch to non-root user
USER wohnfair

# Expose port
EXPOSE 50052

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD ["/bin/grpc_health_probe", "-addr=localhost:50052"]

# Run the service
CMD ["/app/zk-lease"]
