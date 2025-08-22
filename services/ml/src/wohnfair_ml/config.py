"""
Configuration management for WohnFair ML service.
"""

import os
from pathlib import Path
from typing import Any, Dict, List, Optional, Union

from pydantic import BaseSettings, Field, validator
from pydantic_settings import BaseSettings as PydanticBaseSettings


class DatabaseSettings(BaseSettings):
    """Database connection settings."""
    
    url: str = Field(
        default="postgresql://wohnfair:wohnfair_pass@localhost:5432/wohnfair",
        description="Database connection URL"
    )
    pool_size: int = Field(default=10, description="Connection pool size")
    max_overflow: int = Field(default=20, description="Maximum overflow connections")
    echo: bool = Field(default=False, description="Enable SQL logging")

    class Config:
        env_prefix = "DB_"


class RedisSettings(BaseSettings):
    """Redis connection settings."""
    
    url: str = Field(
        default="redis://:redis_pass@localhost:6379",
        description="Redis connection URL"
    )
    pool_size: int = Field(default=5, description="Connection pool size")
    timeout: int = Field(default=5, description="Connection timeout in seconds")

    class Config:
        env_prefix = "REDIS_"


class ClickHouseSettings(BaseSettings):
    """ClickHouse connection settings."""
    
    host: str = Field(default="localhost", description="ClickHouse host")
    port: int = Field(default=8123, description="ClickHouse HTTP port")
    database: str = Field(default="wohnfair_analytics", description="Database name")
    username: str = Field(default="wohnfair", description="Username")
    password: str = Field(default="wohnfair_pass", description="Password")
    secure: bool = Field(default=False, description="Use secure connection")

    class Config:
        env_prefix = "CLICKHOUSE_"


class MinIOSettings(BaseSettings):
    """MinIO/S3 storage settings."""
    
    endpoint: str = Field(default="localhost:9000", description="MinIO endpoint")
    access_key: str = Field(default="minioadmin", description="Access key")
    secret_key: str = Field(default="minioadmin", description="Secret key")
    bucket: str = Field(default="wohnfair-ml", description="Default bucket")
    secure: bool = Field(default=False, description="Use secure connection")

    class Config:
        env_prefix = "MINIO_"


class ModelSettings(BaseSettings):
    """Model configuration settings."""
    
    # HSBCox model settings
    hsbcx_alpha: float = Field(default=0.05, description="Cox model significance level")
    hsbcx_max_iter: int = Field(default=1000, description="Maximum iterations")
    hsbcx_tol: float = Field(default=1e-6, description="Convergence tolerance")
    
    # XGBoost settings
    xgb_n_estimators: int = Field(default=100, description="Number of estimators")
    xgb_max_depth: int = Field(default=6, description="Maximum tree depth")
    xgb_learning_rate: float = Field(default=0.1, description="Learning rate")
    xgb_subsample: float = Field(default=0.8, description="Subsample ratio")
    xgb_colsample_bytree: float = Field(default=0.8, description="Column sample ratio")
    
    # Model storage
    model_dir: str = Field(default="models", description="Model storage directory")
    artifact_dir: str = Field(default="artifacts", description="Artifact storage directory")
    
    # Model versioning
    version_format: str = Field(default="v{major}.{minor}.{patch}", description="Version format")
    auto_version: bool = Field(default=True, description="Auto-increment versions")

    class Config:
        env_prefix = "MODEL_"


class TrainingSettings(BaseSettings):
    """Training configuration settings."""
    
    # Data settings
    train_split: float = Field(default=0.8, description="Training data split ratio")
    val_split: float = Field(default=0.1, description="Validation data split ratio")
    test_split: float = Field(default=0.1, description="Test data split ratio")
    random_state: int = Field(default=42, description="Random seed")
    
    # Cross-validation
    cv_folds: int = Field(default=5, description="Cross-validation folds")
    cv_strategy: str = Field(default="stratified", description="CV strategy")
    
    # Hyperparameter tuning
    tune_enabled: bool = Field(default=True, description="Enable hyperparameter tuning")
    tune_trials: int = Field(default=100, description="Number of tuning trials")
    tune_timeout: int = Field(default=3600, description="Tuning timeout in seconds")
    
    # Early stopping
    early_stopping: bool = Field(default=True, description="Enable early stopping")
    patience: int = Field(default=10, description="Early stopping patience")
    
    # Distributed training
    distributed: bool = Field(default=False, description="Enable distributed training")
    num_workers: int = Field(default=4, description="Number of worker processes")

    class Config:
        env_prefix = "TRAINING_"


class EvaluationSettings(BaseSettings):
    """Evaluation configuration settings."""
    
    # Metrics
    primary_metric: str = Field(default="c_index", description="Primary evaluation metric")
    secondary_metrics: List[str] = Field(
        default=["brier_score", "calibration_error", "discrimination_index"],
        description="Secondary evaluation metrics"
    )
    
    # Thresholds
    min_c_index: float = Field(default=0.7, description="Minimum acceptable C-index")
    max_brier_score: float = Field(default=0.25, description="Maximum acceptable Brier score")
    
    # Statistical tests
    statistical_tests: bool = Field(default=True, description="Enable statistical significance tests")
    confidence_level: float = Field(default=0.95, description="Confidence level for tests")
    
    # Fairness metrics
    fairness_metrics: List[str] = Field(
        default=["demographic_parity", "equalized_odds", "calibration_by_group"],
        description="Fairness metrics to compute"
    )
    protected_attributes: List[str] = Field(
        default=["age_group", "gender", "ethnicity", "disability_status"],
        description="Protected attributes for fairness analysis"
    )

    class Config:
        env_prefix = "EVALUATION_"


class MonitoringSettings(BaseSettings):
    """Monitoring and observability settings."""
    
    # Prometheus
    prometheus_enabled: bool = Field(default=True, description="Enable Prometheus metrics")
    prometheus_port: int = Field(default=9092, description="Prometheus metrics port")
    prometheus_path: str = Field(default="/metrics", description="Metrics endpoint path")
    
    # Tracing
    tracing_enabled: bool = Field(default=True, description="Enable OpenTelemetry tracing")
    jaeger_endpoint: str = Field(
        default="http://localhost:14268/api/traces",
        description="Jaeger tracing endpoint"
    )
    sample_rate: float = Field(default=0.1, description="Tracing sample rate")
    
    # Logging
    log_level: str = Field(default="INFO", description="Logging level")
    log_format: str = Field(default="json", description="Log format")
    log_file: Optional[str] = Field(default=None, description="Log file path")

    class Config:
        env_prefix = "MONITORING_"


class MLFlowSettings(BaseSettings):
    """MLflow tracking settings."""
    
    enabled: bool = Field(default=True, description="Enable MLflow tracking")
    tracking_uri: str = Field(default="http://localhost:5000", description="MLflow tracking URI")
    experiment_name: str = Field(default="wohnfair-housing", description="MLflow experiment name")
    artifact_location: Optional[str] = Field(default=None, description="Artifact storage location")
    
    # Model registry
    registry_enabled: bool = Field(default=True, description="Enable model registry")
    registry_uri: Optional[str] = Field(default=None, description="Model registry URI")
    
    # Auto-logging
    autolog: bool = Field(default=True, description="Enable auto-logging")
    log_models: bool = Field(default=True, description="Log trained models")
    log_artifacts: bool = Field(default=True, description="Log training artifacts")

    class Config:
        env_prefix = "MLFLOW_"


class Settings(PydanticBaseSettings):
    """Main application settings."""
    
    # Service settings
    service_name: str = Field(default="wohnfair-ml", description="Service name")
    service_version: str = Field(default="0.1.0", description="Service version")
    debug: bool = Field(default=False, description="Enable debug mode")
    
    # Database settings
    database: DatabaseSettings = Field(default_factory=DatabaseSettings)
    
    # Storage settings
    redis: RedisSettings = Field(default_factory=RedisSettings)
    clickhouse: ClickHouseSettings = Field(default_factory=ClickHouseSettings)
    minio: MinIOSettings = Field(default_factory=MinIOSettings)
    
    # ML-specific settings
    model: ModelSettings = Field(default_factory=ModelSettings)
    training: TrainingSettings = Field(default_factory=TrainingSettings)
    evaluation: EvaluationSettings = Field(default_factory=EvaluationSettings)
    
    # Observability settings
    monitoring: MonitoringSettings = Field(default_factory=MonitoringSettings)
    
    # MLflow settings
    mlflow: MLFlowSettings = Field(default_factory=MLFlowSettings)
    
    # File paths
    base_dir: Path = Field(default=Path.cwd(), description="Base directory")
    data_dir: Path = Field(default=Path("data"), description="Data directory")
    model_dir: Path = Field(default=Path("models"), description="Model directory")
    artifact_dir: Path = Field(default=Path("artifacts"), description="Artifact directory")
    log_dir: Path = Field(default=Path("logs"), description="Log directory")
    
    @validator("base_dir", "data_dir", "model_dir", "artifact_dir", "log_dir", pre=True)
    def validate_paths(cls, v: Union[str, Path]) -> Path:
        """Convert string paths to Path objects and ensure they exist."""
        path = Path(v) if isinstance(v, str) else v
        if not path.is_absolute():
            path = Path.cwd() / path
        path.mkdir(parents=True, exist_ok=True)
        return path
    
    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"
        case_sensitive = False


# Global settings instance
_settings: Optional[Settings] = None


def get_settings() -> Settings:
    """Get the global settings instance."""
    global _settings
    if _settings is None:
        _settings = Settings()
    return _settings


def reload_settings() -> Settings:
    """Reload settings from environment."""
    global _settings
    _settings = Settings()
    return _settings


# Convenience functions for common settings
def get_database_url() -> str:
    """Get database connection URL."""
    return get_settings().database.url


def get_redis_url() -> str:
    """Get Redis connection URL."""
    return get_settings().redis.url


def get_clickhouse_config() -> Dict[str, Any]:
    """Get ClickHouse configuration."""
    settings = get_settings().clickhouse
    return {
        "host": settings.host,
        "port": settings.port,
        "database": settings.database,
        "username": settings.username,
        "password": settings.password,
        "secure": settings.secure,
    }


def get_minio_config() -> Dict[str, Any]:
    """Get MinIO configuration."""
    settings = get_settings().minio
    return {
        "endpoint": settings.endpoint,
        "access_key": settings.access_key,
        "secret_key": settings.secret_key,
        "bucket": settings.bucket,
        "secure": settings.secure,
    }


def get_model_dir() -> Path:
    """Get model storage directory."""
    return get_settings().model_dir


def get_artifact_dir() -> Path:
    """Get artifact storage directory."""
    return get_settings().artifact_dir


def get_log_dir() -> Path:
    """Get log directory."""
    return get_settings().log_dir


def is_debug_mode() -> bool:
    """Check if debug mode is enabled."""
    return get_settings().debug


def is_distributed_training() -> bool:
    """Check if distributed training is enabled."""
    return get_settings().training.distributed


def is_mlflow_enabled() -> bool:
    """Check if MLflow tracking is enabled."""
    return get_settings().mlflow.enabled


def is_prometheus_enabled() -> bool:
    """Check if Prometheus metrics are enabled."""
    return get_settings().monitoring.prometheus_enabled


def is_tracing_enabled() -> bool:
    """Check if OpenTelemetry tracing is enabled."""
    return get_settings().monitoring.tracing_enabled
