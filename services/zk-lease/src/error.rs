use thiserror::Error;
use tonic::{Code, Status};

#[derive(Error, Debug)]
pub enum ZkLeaseError {
    #[error("Configuration error: {0}")]
    Config(String),

    #[error("Database error: {0}")]
    Database(#[from] sqlx::Error),

    #[error("Redis error: {0}")]
    Redis(#[from] redis::RedisError),

    #[error("Proof generation error: {0}")]
    ProofGeneration(String),

    #[error("Proof verification error: {0}")]
    ProofVerification(String),

    #[error("Circuit compilation error: {0}")]
    CircuitCompilation(String),

    #[error("Invalid input: {0}")]
    InvalidInput(String),

    #[error("Timeout error: {0}")]
    Timeout(String),

    #[error("Internal error: {0}")]
    Internal(String),

    #[error("Not found: {0}")]
    NotFound(String),

    #[error("Unauthorized: {0}")]
    Unauthorized(String),

    #[error("Validation error: {0}")]
    Validation(String),

    #[error("Serialization error: {0}")]
    Serialization(#[from] serde_json::Error),

    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),

    #[error("Cryptographic error: {0}")]
    Cryptographic(String),
}

impl From<ZkLeaseError> for Status {
    fn from(err: ZkLeaseError) -> Self {
        match err {
            ZkLeaseError::Config(msg) => Status::new(Code::Internal, msg),
            ZkLeaseError::Database(msg) => Status::new(Code::Internal, msg.to_string()),
            ZkLeaseError::Redis(msg) => Status::new(Code::Internal, msg.to_string()),
            ZkLeaseError::ProofGeneration(msg) => Status::new(Code::Internal, msg),
            ZkLeaseError::ProofVerification(msg) => Status::new(Code::Internal, msg),
            ZkLeaseError::CircuitCompilation(msg) => Status::new(Code::Internal, msg),
            ZkLeaseError::InvalidInput(msg) => Status::new(Code::InvalidArgument, msg),
            ZkLeaseError::Timeout(msg) => Status::new(Code::DeadlineExceeded, msg),
            ZkLeaseError::Internal(msg) => Status::new(Code::Internal, msg),
            ZkLeaseError::NotFound(msg) => Status::new(Code::NotFound, msg),
            ZkLeaseError::Unauthorized(msg) => Status::new(Code::PermissionDenied, msg),
            ZkLeaseError::Validation(msg) => Status::new(Code::InvalidArgument, msg),
            ZkLeaseError::Serialization(msg) => Status::new(Code::Internal, msg.to_string()),
            ZkLeaseError::Io(msg) => Status::new(Code::Internal, msg.to_string()),
            ZkLeaseError::Cryptographic(msg) => Status::new(Code::Internal, msg),
        }
    }
}

impl From<config::ConfigError> for ZkLeaseError {
    fn from(err: config::ConfigError) -> Self {
        ZkLeaseError::Config(err.to_string())
    }
}

impl From<serde_json::Error> for ZkLeaseError {
    fn from(err: serde_json::Error) -> Self {
        ZkLeaseError::Serialization(err)
    }
}

impl From<std::io::Error> for ZkLeaseError {
    fn from(err: std::io::Error) -> Self {
        ZkLeaseError::Io(err)
    }
}

impl From<sqlx::Error> for ZkLeaseError {
    fn from(err: sqlx::Error) -> Self {
        ZkLeaseError::Database(err)
    }
}

impl From<redis::RedisError> for ZkLeaseError {
    fn from(err: redis::RedisError) -> Self {
        ZkLeaseError::Redis(err)
    }
}

pub type Result<T> = std::result::Result<T, ZkLeaseError>;

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_error_conversion_to_status() {
        let config_error = ZkLeaseError::Config("test config error".to_string());
        let status: Status = config_error.into();
        assert_eq!(status.code(), Code::Internal);

        let input_error = ZkLeaseError::InvalidInput("test input error".to_string());
        let status: Status = input_error.into();
        assert_eq!(status.code(), Code::InvalidArgument);

        let not_found_error = ZkLeaseError::NotFound("test not found error".to_string());
        let status: Status = not_found_error.into();
        assert_eq!(status.code(), Code::NotFound);
    }

    #[test]
    fn test_error_from_config_error() {
        let config_error = config::ConfigError::NotFound("test".to_string());
        let zk_error: ZkLeaseError = config_error.into();
        assert!(matches!(zk_error, ZkLeaseError::Config(_)));
    }
}
