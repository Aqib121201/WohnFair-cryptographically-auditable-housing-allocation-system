use serde::{Deserialize, Serialize};
use std::env;
use std::path::Path;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Config {
    pub server: ServerConfig,
    pub database: DatabaseConfig,
    pub redis: RedisConfig,
    pub prover: ProverConfig,
    pub verifier: VerifierConfig,
    pub metrics: MetricsConfig,
    pub tracing: TracingConfig,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ServerConfig {
    pub host: String,
    pub port: u16,
    pub workers: usize,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DatabaseConfig {
    pub url: String,
    pub max_connections: u32,
    pub min_connections: u32,
    pub connect_timeout: u64,
    pub idle_timeout: u64,
    pub max_lifetime: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RedisConfig {
    pub url: String,
    pub pool_size: usize,
    pub timeout: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProverConfig {
    pub circuit_path: String,
    pub proving_key_path: String,
    pub max_proof_size: usize,
    pub timeout: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct VerifierConfig {
    pub verifying_key_path: String,
    pub max_verification_time: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MetricsConfig {
    pub enabled: bool,
    pub port: u16,
    pub path: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TracingConfig {
    pub enabled: bool,
    pub jaeger_endpoint: String,
    pub sample_rate: f64,
}

impl Config {
    pub fn load() -> Result<Self, config::ConfigError> {
        let config_dir = env::var("CONFIG_DIR").unwrap_or_else(|_| "config".to_string());
        let config_path = Path::new(&config_dir).join("config.yaml");

        let mut builder = config::Config::builder()
            .set_default("server.host", "0.0.0.0")?
            .set_default("server.port", 50052)?
            .set_default("server.workers", 4)?
            .set_default("database.max_connections", 10)?
            .set_default("database.min_connections", 2)?
            .set_default("database.connect_timeout", 30)?
            .set_default("database.idle_timeout", 300)?
            .set_default("database.max_lifetime", 3600)?
            .set_default("redis.pool_size", 5)?
            .set_default("redis.timeout", 5)?
            .set_default("prover.max_proof_size", 1024 * 1024)?
            .set_default("prover.timeout", 300)?
            .set_default("verifier.max_verification_time", 60)?
            .set_default("metrics.enabled", true)?
            .set_default("metrics.port", 9091)?
            .set_default("metrics.path", "/metrics")?
            .set_default("tracing.enabled", true)?
            .set_default("tracing.sample_rate", 0.1)?;

        // Load from config file if it exists
        if config_path.exists() {
            builder = builder.add_source(config::File::from(config_path));
        }

        // Load from environment variables
        builder = builder.add_source(config::Environment::default().separator("_"));

        // Build and deserialize
        let config = builder.build()?;
        let config: Config = config.try_deserialize()?;

        Ok(config)
    }

    pub fn from_env() -> Self {
        Self {
            server: ServerConfig {
                host: env::var("ZK_LEASE_HOST").unwrap_or_else(|_| "0.0.0.0".to_string()),
                port: env::var("ZK_LEASE_PORT")
                    .unwrap_or_else(|_| "50052".to_string())
                    .parse()
                    .unwrap_or(50052),
                workers: env::var("ZK_LEASE_WORKERS")
                    .unwrap_or_else(|_| "4".to_string())
                    .parse()
                    .unwrap_or(4),
            },
            database: DatabaseConfig {
                url: env::var("DATABASE_URL").unwrap_or_else(|_| {
                    "postgresql://wohnfair:wohnfair_pass@localhost:5432/wohnfair".to_string()
                }),
                max_connections: env::var("DB_MAX_CONNECTIONS")
                    .unwrap_or_else(|_| "10".to_string())
                    .parse()
                    .unwrap_or(10),
                min_connections: env::var("DB_MIN_CONNECTIONS")
                    .unwrap_or_else(|_| "2".to_string())
                    .parse()
                    .unwrap_or(2),
                connect_timeout: env::var("DB_CONNECT_TIMEOUT")
                    .unwrap_or_else(|_| "30".to_string())
                    .parse()
                    .unwrap_or(30),
                idle_timeout: env::var("DB_IDLE_TIMEOUT")
                    .unwrap_or_else(|_| "300".to_string())
                    .parse()
                    .unwrap_or(300),
                max_lifetime: env::var("DB_MAX_LIFETIME")
                    .unwrap_or_else(|_| "3600".to_string())
                    .parse()
                    .unwrap_or(3600),
            },
            redis: RedisConfig {
                url: env::var("REDIS_URL").unwrap_or_else(|_| {
                    "redis://:redis_pass@localhost:6379".to_string()
                }),
                pool_size: env::var("REDIS_POOL_SIZE")
                    .unwrap_or_else(|_| "5".to_string())
                    .parse()
                    .unwrap_or(5),
                timeout: env::var("REDIS_TIMEOUT")
                    .unwrap_or_else(|_| "5".to_string())
                    .parse()
                    .unwrap_or(5),
            },
            prover: ProverConfig {
                circuit_path: env::var("CIRCUIT_PATH").unwrap_or_else(|_| "circuits".to_string()),
                proving_key_path: env::var("PROVING_KEY_PATH")
                    .unwrap_or_else(|_| "keys/proving.key".to_string()),
                max_proof_size: env::var("MAX_PROOF_SIZE")
                    .unwrap_or_else(|_| "1048576".to_string())
                    .parse()
                    .unwrap_or(1024 * 1024),
                timeout: env::var("PROVER_TIMEOUT")
                    .unwrap_or_else(|_| "300".to_string())
                    .parse()
                    .unwrap_or(300),
            },
            verifier: VerifierConfig {
                verifying_key_path: env::var("VERIFYING_KEY_PATH")
                    .unwrap_or_else(|_| "keys/verifying.key".to_string()),
                max_verification_time: env::var("MAX_VERIFICATION_TIME")
                    .unwrap_or_else(|_| "60".to_string())
                    .parse()
                    .unwrap_or(60),
            },
            metrics: MetricsConfig {
                enabled: env::var("METRICS_ENABLED")
                    .unwrap_or_else(|_| "true".to_string())
                    .parse()
                    .unwrap_or(true),
                port: env::var("METRICS_PORT")
                    .unwrap_or_else(|_| "9091".to_string())
                    .parse()
                    .unwrap_or(9091),
                path: env::var("METRICS_PATH")
                    .unwrap_or_else(|_| "/metrics".to_string()),
            },
            tracing: TracingConfig {
                enabled: env::var("TRACING_ENABLED")
                    .unwrap_or_else(|_| "true".to_string())
                    .parse()
                    .unwrap_or(true),
                jaeger_endpoint: env::var("JAEGER_ENDPOINT")
                    .unwrap_or_else(|_| "http://localhost:14268/api/traces".to_string()),
                sample_rate: env::var("TRACING_SAMPLE_RATE")
                    .unwrap_or_else(|_| "0.1".to_string())
                    .parse()
                    .unwrap_or(0.1),
            },
        }
    }
}

impl Default for Config {
    fn default() -> Self {
        Self::from_env()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_default_config() {
        let config = Config::default();
        assert_eq!(config.server.port, 50052);
        assert_eq!(config.server.workers, 4);
        assert_eq!(config.metrics.enabled, true);
        assert_eq!(config.tracing.enabled, true);
    }

    #[test]
    fn test_config_from_env() {
        env::set_var("ZK_LEASE_PORT", "50053");
        env::set_var("ZK_LEASE_WORKERS", "8");
        
        let config = Config::from_env();
        assert_eq!(config.server.port, 50053);
        assert_eq!(config.server.workers, 8);
        
        // Clean up
        env::remove_var("ZK_LEASE_PORT");
        env::remove_var("ZK_LEASE_WORKERS");
    }
}
