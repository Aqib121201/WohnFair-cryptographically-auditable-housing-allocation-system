use std::net::SocketAddr;
use tonic::{transport::Server, Request, Response, Status};
use tracing::{info};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

pub mod gen {
    pub mod wohnfair {
        pub mod policy {
            pub mod v1 {
                include!(concat!(env!("CARGO_MANIFEST_DIR"), "/src/gen/wohnfair.policy.v1.rs"));
            }
        }
        pub mod common {
            pub mod v1 {
                include!(concat!(env!("CARGO_MANIFEST_DIR"), "/src/gen/wohnfair.common.v1.rs"));
            }
        }
    }
}

use gen::wohnfair::policy::v1 as policyv1;
use policyv1::policy_service_server::{PolicyService, PolicyServiceServer};

#[derive(Default, Clone)]
struct PolicyServer;

#[tonic::async_trait]
impl PolicyService for PolicyServer {
    async fn validate_eligibility(
        &self,
        _request: Request<policyv1::ValidateEligibilityRequest>,
    ) -> Result<Response<policyv1::ValidateEligibilityResponse>, Status> {
        Ok(Response::new(policyv1::ValidateEligibilityResponse { eligible: true, reasons: vec![] }))
    }

    async fn evaluate_quota(
        &self,
        _request: Request<policyv1::EvaluateQuotaRequest>,
    ) -> Result<Response<policyv1::EvaluateQuotaResponse>, Status> {
        Ok(Response::new(policyv1::EvaluateQuotaResponse { approved: true, quota_used: 1.0, details: vec![] }))
    }

    async fn compile_policy(
        &self,
        request: Request<policyv1::CompilePolicyRequest>,
    ) -> Result<Response<policyv1::CompilePolicyResponse>, Status> {
        let src = request.into_inner().source;
        // Minimal stub: wrap DSL as Rego policy comment
        let rego = format!("package wohnfair.policy\n\n# compiled from DSL\n# ---\n# {}\n\ndefault allow = true\n", src.replace("\n", "\n# "));
        Ok(Response::new(policyv1::CompilePolicyResponse { rego, messages: vec![] }))
    }

    async fn get_policy_version(
        &self,
        _request: Request<policyv1::GetPolicyVersionRequest>,
    ) -> Result<Response<policyv1::GetPolicyVersionResponse>, Status> {
        Ok(Response::new(policyv1::GetPolicyVersionResponse { version: "0.1.0".into() }))
    }

    async fn health(
        &self,
        _request: Request<policyv1::HealthRequest>,
    ) -> Result<Response<policyv1::HealthResponse>, Status> {
        Ok(Response::new(policyv1::HealthResponse { status: "SERVING".into() }))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    tracing_subscriber::registry()
        .with(tracing_subscriber::EnvFilter::new(std::env::var("RUST_LOG").unwrap_or_else(|_| "info".into())))
        .with(tracing_subscriber::fmt::layer())
        .init();

    let addr: SocketAddr = "0.0.0.0:50053".parse()?;
    info!("Policy DSL service listening on {}", addr);

    let svc = PolicyServer::default();
    // gRPC health service
    let (health_reporter, health_service) = tonic_health::server::health_reporter();
    health_reporter.set_serving::<PolicyServiceServer<PolicyServer>>().await;

    Server::builder()
        .add_service(PolicyServiceServer::new(svc))
        .add_service(health_service)
        .serve(addr)
        .await?;

    Ok(())
}
