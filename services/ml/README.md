# WohnFair ML Service

The WohnFair ML Service provides machine learning capabilities for housing allocation fairness prediction, demand forecasting, and risk assessment. It implements advanced statistical models including HSBCox (Cox proportional hazards with XGBoost residual boosting) and comprehensive fairness-aware evaluation metrics.

## Architecture

The service is built with a modular, microservices architecture:

- **Model Management**: Handles model training, versioning, and deployment
- **Data Processing**: Comprehensive data preprocessing and feature engineering
- **Training Pipeline**: Automated training with hyperparameter tuning
- **Evaluation Engine**: Multi-metric evaluation with fairness analysis
- **Serving API**: RESTful API for model inference and batch predictions
- **Monitoring**: Comprehensive observability with Prometheus and Jaeger

## Features

### Core ML Capabilities

- **HSBCox Model**: Cox proportional hazards model with XGBoost residual boosting
- **Fairness-Aware Training**: Built-in fairness constraints and bias detection
- **Multi-Modal Support**: Handles structured, text, and temporal data
- **AutoML**: Automated hyperparameter tuning and model selection
- **Model Interpretability**: SHAP, LIME, and custom explanation methods

### Advanced Analytics

- **Survival Analysis**: Time-to-event modeling for housing allocation
- **Demand Forecasting**: Predictive modeling for housing demand
- **Risk Assessment**: Risk scoring for housing applications
- **Fairness Metrics**: Comprehensive bias detection and fairness analysis
- **Causal Inference**: Causal modeling for policy impact assessment

### Production Features

- **Model Versioning**: Semantic versioning with rollback capabilities
- **A/B Testing**: Model comparison and traffic splitting
- **Performance Monitoring**: Real-time model performance tracking
- **Scalable Serving**: Horizontal scaling with load balancing
- **Security**: Authentication, authorization, and input validation

## Quick Start

### Prerequisites

- Python 3.9+
- PostgreSQL 15+
- Redis 7+
- ClickHouse 22+
- MinIO/S3
- Docker & Docker Compose

### Local Development

1. **Clone and setup**:
   ```bash
   cd services/ml
   python -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   pip install -e .
   ```

2. **Environment setup**:
   ```bash
   export DATABASE_URL="postgresql://wohnfair:wohnfair_pass@localhost:5432/wohnfair"
   export REDIS_URL="redis://:redis_pass@localhost:6379"
   export CLICKHOUSE_HOST="localhost"
   export MINIO_ENDPOINT="localhost:9000"
   ```

3. **Run the service**:
   ```bash
   python -m wohnfair_ml.cli serve
   ```

### Docker

```bash
# Build the image
docker build -t wohnfair-ml .

# Run the container
docker run -p 8000:8000 -p 9092:9092 \
  -e DATABASE_URL="postgresql://wohnfair:wohnfair_pass@host.docker.internal:5432/wohnfair" \
  -e REDIS_URL="redis://:redis_pass@host.docker.internal:6379" \
  wohnfair-ml
```

## API Reference

### REST API Endpoints

#### Health and Status
- `GET /healthz` - Health check
- `GET /status` - Service status
- `GET /metrics` - Prometheus metrics

#### Model Management
- `GET /models` - List available models
- `POST /models` - Upload new model
- `GET /models/{model_id}` - Get model details
- `DELETE /models/{model_id}` - Delete model
- `POST /models/{model_id}/versions` - Create new version

#### Training
- `POST /training/jobs` - Start training job
- `GET /training/jobs/{job_id}` - Get training status
- `GET /training/jobs` - List training jobs
- `DELETE /training/jobs/{job_id}` - Cancel training job

#### Inference
- `POST /predict` - Single prediction
- `POST /predict/batch` - Batch prediction
- `POST /predict/stream` - Streaming prediction

#### Evaluation
- `POST /evaluate` - Evaluate model performance
- `GET /evaluate/{model_id}` - Get evaluation results
- `POST /evaluate/fairness` - Fairness analysis

### gRPC API

The service also provides a gRPC interface for high-performance communication:

```protobuf
service WohnFairMLService {
  rpc TrainModel(TrainModelRequest) returns (TrainModelResponse);
  rpc Predict(PredictRequest) returns (PredictResponse);
  rpc EvaluateModel(EvaluateModelRequest) returns (EvaluateModelResponse);
  rpc GetModelInfo(GetModelInfoRequest) returns (GetModelInfoResponse);
}
```

## Model Types

### HSBCox Model

The HSBCox model combines Cox proportional hazards with XGBoost residual boosting:

**Features**:
- Survival analysis for time-to-housing allocation
- XGBoost residual boosting for non-linear effects
- Fairness constraints for protected groups
- Interpretable feature importance

**Use Cases**:
- Housing wait time prediction
- Risk assessment for housing applications
- Fairness-aware allocation modeling

### Demand Forecasting Models

Multiple models for housing demand prediction:

- **Time Series Models**: ARIMA, Prophet, LSTM
- **Ensemble Methods**: Random Forest, XGBoost, LightGBM
- **Deep Learning**: Neural networks with attention mechanisms

### Fairness Models

Specialized models for fairness-aware prediction:

- **Adversarial Debiasing**: Remove bias during training
- **Preprocessing Methods**: Reweighting and resampling
- **Post-processing**: Calibration and threshold adjustment

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | PostgreSQL connection string | Database connection |
| `REDIS_URL` | Redis connection string | Redis connection |
| `CLICKHOUSE_HOST` | `localhost` | ClickHouse host |
| `MINIO_ENDPOINT` | `localhost:9000` | MinIO endpoint |
| `MLFLOW_TRACKING_URI` | `http://localhost:5000` | MLflow tracking URI |
| `JAEGER_ENDPOINT` | Jaeger endpoint | Tracing endpoint |

### Configuration File

The service supports YAML configuration files. See `config/config.yaml` for the complete configuration schema.

## Training Pipeline

### Data Pipeline

1. **Data Ingestion**: Load data from multiple sources
2. **Data Validation**: Schema validation and quality checks
3. **Feature Engineering**: Automated feature creation and selection
4. **Data Splitting**: Stratified train/validation/test splits

### Training Process

1. **Model Selection**: Choose appropriate model architecture
2. **Hyperparameter Tuning**: Automated optimization with Optuna
3. **Cross-Validation**: K-fold cross-validation with stratification
4. **Early Stopping**: Prevent overfitting with validation monitoring
5. **Model Persistence**: Save models and artifacts

### Evaluation

1. **Performance Metrics**: C-index, Brier score, calibration error
2. **Fairness Analysis**: Demographic parity, equalized odds
3. **Statistical Testing**: Significance tests for model comparison
4. **Bias Detection**: Identify and quantify model bias

## Model Serving

### REST API

FastAPI-based REST API with automatic documentation:

- **OpenAPI/Swagger**: Interactive API documentation
- **Request Validation**: Pydantic-based input validation
- **Response Caching**: Redis-based response caching
- **Rate Limiting**: Configurable rate limiting

### gRPC Interface

High-performance gRPC interface for:

- **Batch Processing**: Efficient batch predictions
- **Streaming**: Real-time prediction streams
- **Model Management**: Model operations and metadata

### Model Registry

MLflow-based model registry with:

- **Version Control**: Semantic versioning for models
- **Artifact Storage**: Model files and metadata
- **Experiment Tracking**: Training experiments and metrics
- **Model Lineage**: Track model dependencies and changes

## Monitoring and Observability

### Metrics

Prometheus metrics for:

- **Model Performance**: Prediction accuracy and latency
- **System Health**: Resource usage and availability
- **Business Metrics**: Prediction volume and success rates

### Tracing

OpenTelemetry tracing for:

- **Request Flow**: End-to-end request tracing
- **Model Inference**: Prediction pipeline tracing
- **Training Jobs**: Training process monitoring

### Logging

Structured logging with:

- **JSON Format**: Machine-readable log format
- **Log Levels**: Configurable logging levels
- **Log Rotation**: Automatic log file management

## Performance

### Benchmarks

| Model Type | Training Time | Inference Latency | Memory Usage |
|------------|---------------|-------------------|--------------|
| HSBCox | ~5 min | ~10ms | ~500MB |
| XGBoost | ~2 min | ~5ms | ~200MB |
| Neural Network | ~15 min | ~20ms | ~1GB |

### Optimization Features

- **Model Quantization**: Reduce model size and latency
- **Batch Processing**: Efficient batch predictions
- **Caching**: Redis-based prediction caching
- **Parallel Processing**: Multi-worker inference

## Security

### Authentication

- **JWT Tokens**: Stateless authentication
- **API Keys**: Simple API key authentication
- **OAuth2**: Integration with external providers

### Authorization

- **Role-Based Access**: Fine-grained permission control
- **Model Access Control**: Control access to specific models
- **Data Privacy**: GDPR-compliant data handling

### Input Validation

- **Schema Validation**: Strict input schema validation
- **Size Limits**: Configurable input size limits
- **Sanitization**: Input sanitization and cleaning

## Development

### Project Structure

```
src/wohnfair_ml/
├── __init__.py           # Package initialization
├── config.py             # Configuration management
├── models/               # Model implementations
│   ├── hsbcx.py         # HSBCox model
│   ├── forecasting.py   # Demand forecasting models
│   └── fairness.py      # Fairness-aware models
├── preprocessing/        # Data preprocessing
│   ├── features.py      # Feature engineering
│   ├── validation.py    # Data validation
│   └── cleaning.py      # Data cleaning
├── training/            # Training pipeline
│   ├── pipeline.py      # Training pipeline
│   ├── tuning.py        # Hyperparameter tuning
│   └── validation.py    # Model validation
├── evaluation/          # Model evaluation
│   ├── metrics.py       # Evaluation metrics
│   ├── fairness.py      # Fairness metrics
│   └── comparison.py    # Model comparison
├── utils/               # Utility functions
│   ├── database.py      # Database utilities
│   ├── storage.py       # Storage utilities
│   └── monitoring.py    # Monitoring utilities
└── cli.py               # Command-line interface
```

### Adding New Models

1. **Implement the model** in `src/models/`
2. **Add model configuration** in `config/config.yaml`
3. **Update the training pipeline** to handle the new model
4. **Add evaluation metrics** specific to the model
5. **Write tests** for the new model
6. **Update documentation**

### Testing

```bash
# Run all tests
pytest

# Run specific test category
pytest tests/unit/
pytest tests/integration/

# Run with coverage
pytest --cov=wohnfair_ml

# Run performance tests
pytest tests/performance/ -m "slow"
```

### Code Quality

```bash
# Format code
black src/
isort src/

# Lint code
flake8 src/
mypy src/

# Security checks
bandit -r src/
safety check
```

## Deployment

### Docker Compose

The service is included in the main `docker-compose.yml`:

```yaml
ml:
  build:
    context: ./services/ml
    dockerfile: Dockerfile
  ports:
    - "8000:8000"
    - "9092:9092"
  environment:
    - DATABASE_URL=postgresql://wohnfair:wohnfair_pass@postgres:5432/wohnfair
    - REDIS_URL=redis://:redis_pass@redis:6379
    - CLICKHOUSE_HOST=clickhouse
    - MINIO_ENDPOINT=minio:9000
  networks:
    - wohnfair-backend
```

### Kubernetes

See `infra/k8s/` for Kubernetes deployment manifests.

### Production Considerations

- **Resource Limits**: Set appropriate CPU/memory limits
- **Scaling**: Use horizontal pod autoscaling
- **Monitoring**: Enable comprehensive monitoring
- **Backup**: Regular backup of models and configurations
- **Security**: Use secrets management for sensitive data

## Troubleshooting

### Common Issues

1. **Model Loading Failures**
   - Check model file paths
   - Verify model format compatibility
   - Check available memory

2. **Training Failures**
   - Check data quality and format
   - Verify hyperparameter ranges
   - Check system resources

3. **Inference Errors**
   - Validate input data format
   - Check model compatibility
   - Verify model state

### Debug Mode

Enable debug logging:

```bash
export LOG_LEVEL=DEBUG
python -m wohnfair_ml.cli serve
```

### Performance Profiling

Enable profiling:

```bash
export PROFILING_ENABLED=true
python -m wohnfair_ml.cli serve
```

## Contributing

See [CONTRIBUTING.md](../../CONTRIBUTING.md) for contribution guidelines.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Implement changes with tests
4. Run the full test suite
5. Submit a pull request

## License

This project is licensed under the MIT License - see [LICENSE](../../LICENSE) for details.

## References

- [Cox Proportional Hazards Model](https://en.wikipedia.org/wiki/Proportional_hazards_model)
- [XGBoost Documentation](https://xgboost.readthedocs.io/)
- [Fairness in Machine Learning](https://fairmlbook.org/)
- [Survival Analysis](https://en.wikipedia.org/wiki/Survival_analysis)
- [WohnFair Architecture](../../README.md#architecture)
