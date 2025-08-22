use std::net::SocketAddr;
use std::sync::Arc;
use tonic::transport::Server;
use tracing::{info, warn, error};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

mod config;
mod error;
mod grpc;
mod prover;
mod verifier;
mod utils;
mod metrics;

use config::Config;
use grpc::zk_lease_service::ZkLeaseService;
use error::ZkLeaseError;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Initialize tracing
    tracing_subscriber::registry()
        .with(tracing_subscriber::EnvFilter::new(
            std::env::var("RUST_LOG").unwrap_or_else(|_| "info".into()),
        ))
        .with(tracing_subscriber::fmt::layer())
        .init();

    info!("Starting ZK-Lease service...");

    // Load configuration
    let config = Config::load()?;
    info!("Configuration loaded successfully");

    // Initialize metrics
    let metrics = Arc::new(metrics::Metrics::new());
    info!("Metrics initialized");

    // Initialize prover and verifier
    let prover = Arc::new(prover::Prover::new(&config)?);
    let verifier = Arc::new(verifier::Verifier::new(&config)?);
    info!("Prover and verifier initialized");

    // Create gRPC service
    let service = ZkLeaseService::new(prover, verifier, metrics);
    // Health reporter
    let (health_reporter, health_service) = tonic_health::server::health_reporter();
    health_reporter.set_serving::<grpc::zk_lease_service::zk_lease_service_server::ZkLeaseServiceServer<ZkLeaseService>>().await;

    // Bind address
    let addr: SocketAddr = format!("[::]:{}", config.server.port).parse()?;
    info!("ZK-Lease service listening on {}", addr);

    // Start gRPC server
    Server::builder()
        .add_service(grpc::zk_lease_service::zk_lease_service_server::ZkLeaseServiceServer::new(service))
        .add_service(health_service)
        .serve(addr)
        .await?;

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_config_loading() {
        // Test configuration loading
        let config = Config::load();
        assert!(config.is_ok());
    }
}
