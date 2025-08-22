# Build
FROM rust:1.75-slim as builder
WORKDIR /usr/src/policy-dsl
COPY Cargo.toml Cargo.lock .
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release || true
RUN rm -rf src
COPY build.rs ./
COPY src ./src
COPY ../../services/proto ./services/proto
RUN cargo build --release

# Runtime
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates wget && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /usr/src/policy-dsl/target/release/policy-dsl /app/policy-dsl
RUN wget -qO /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.11/grpc_health_probe-linux-amd64 && \
  chmod +x /bin/grpc_health_probe
USER nobody
EXPOSE 50053
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD ["/bin/grpc_health_probe", "-addr=localhost:50053"]
CMD ["/app/policy-dsl"]
