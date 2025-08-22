# FairRent Service

The FairRent service implements α-fair scheduling algorithms for housing allocation, ensuring equitable distribution while maintaining system efficiency.

## 🏗️ Architecture

FairRent uses a priority queue-based scheduler with the following key components:

- **α-Fair Scheduler**: Implements proportional fairness with configurable α parameter
- **Priority Queue**: Binary heap for O(log n) enqueue/dequeue operations
- **Metrics Collection**: Comprehensive fairness and performance metrics
- **gRPC API**: Protocol buffer-based service interface
- **OpenTelemetry**: Distributed tracing and observability

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Protocol buffer compiler (buf)

### Local Development

```bash
# Navigate to service directory
cd services/fairrent

# Install dependencies
go mod tidy

# Generate protocol buffers
buf generate

# Run tests
go test ./...

# Start service
go run ./cmd/fairrentd
```

### Docker

```bash
# Build image
docker build -t wohnfair/fairrent .

# Run container
docker run -p 50051:50051 wohnfair/fairrent
```

## 📊 α-Fairness Algorithm

The service implements α-fair scheduling where:

- **α = 0**: Maximum throughput (unfair)
- **α = 1**: Proportional fairness
- **α = 2**: Balanced fairness (default)
- **α → ∞**: Maximum fairness (minimum throughput)

### Priority Score Calculation

```
priority = (urgency × group_weight + bonus)^α
```

Where:
- `urgency`: Normalized urgency level (0.0 to 1.0)
- `group_weight`: User group priority multiplier
- `bonus`: Additional priority factors
- `α`: Fairness parameter

### Group Weights

| User Group | Weight | Priority |
|------------|--------|----------|
| Refugee | 1.5 | Highest |
| Disabled | 1.3 | High |
| Senior | 1.2 | High |
| Low Income | 1.1 | Above Average |
| Student | 1.0 | Baseline |
| Family | 1.0 | Baseline |
| Single | 0.9 | Below Average |
| Middle Income | 0.8 | Low |
| High Income | 0.7 | Lowest |

## 🔌 API Reference

### gRPC Methods

#### Enqueue
```protobuf
rpc Enqueue(EnqueueRequest) returns (EnqueueResponse)
```

Adds a new housing request to the queue.

**Request:**
```json
{
  "user_id": "user123",
  "user_group": "USER_GROUP_STUDENT",
  "urgency": "URGENCY_LEVEL_HIGH",
  "financial_constraints": {
    "max_monthly_rent": 800.0
  }
}
```

**Response:**
```json
{
  "ticket_id": "TKT_1234567890",
  "status": "ALLOCATION_STATUS_QUEUED",
  "queue_position": 5,
  "estimated_allocation_time": "2024-01-15T10:00:00Z"
}
```

#### ScheduleNext
```protobuf
rpc ScheduleNext(ScheduleNextRequest) returns (ScheduleNextResponse)
```

Processes the next allocation from the queue.

#### PeekPosition
```protobuf
rpc PeekPosition(PeekPositionRequest) returns (PeekPositionResponse)
```

Returns current queue position and estimated wait time.

#### GetMetrics
```protobuf
rpc GetMetrics(google.protobuf.Empty) returns (FairnessMetrics)
```

Returns comprehensive fairness and performance metrics.

### HTTP Endpoints

The service also exposes HTTP endpoints for monitoring:

- **Health Check**: `GET /healthz`
- **Metrics**: `GET /metrics` (Prometheus format)

## ⚙️ Configuration

Configuration is managed via YAML files and environment variables:

```yaml
# config/config.yaml
scheduler:
  alpha: 2.0
  max_wait_time: "24h"
  group_weights:
    USER_GROUP_REFUGEE: 1.5
    USER_GROUP_DISABLED: 1.3
    # ... more weights
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ALPHA` | 2.0 | Fairness parameter |
| `MAX_WAIT_TIME` | 24h | Maximum wait time |
| `LOG_LEVEL` | info | Logging level |
| `JAEGER_ENDPOINT` | localhost:14268 | Tracing endpoint |

## 📈 Metrics

### Prometheus Metrics

- `fairrent_requests_enqueued_total`: Total requests enqueued
- `fairrent_requests_processed_total`: Total requests processed
- `fairrent_queue_length`: Current queue length
- `fairrent_processing_duration_seconds`: Request processing time
- `fairrent_priority_scores`: Priority score distribution

### Fairness Metrics

- **Wait Time Statistics**: Average, median, P95, P99 wait times
- **Group Fairness**: Allocation rates per user group
- **Starvation Prevention**: Maximum wait time ratios
- **Inequality Measures**: Gini coefficient for wait times

## 🧪 Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test ./internal/scheduler -v
```

### Integration Tests

```bash
# Run integration tests
go test -tags=integration ./...

# Run with race detection
go test -race ./...
```

### Performance Tests

```bash
# Run benchmarks
go test -bench=. ./internal/scheduler

# Run with memory profiling
go test -bench=. -memprofile=mem.prof ./internal/scheduler
```

## 🔍 Monitoring

### Health Checks

The service provides health check endpoints:

```bash
# gRPC health check
grpc_health_probe -addr=localhost:50051

# HTTP health check
curl http://localhost:8080/healthz
```

### Logging

Structured logging with configurable levels:

```json
{
  "level": "info",
  "timestamp": "2024-01-15T09:30:00Z",
  "msg": "Request enqueued successfully",
  "ticket_id": "TKT_1234567890",
  "user_id": "user123",
  "processing_time": "15.2ms"
}
```

### Tracing

Distributed tracing via OpenTelemetry:

- **Service**: fairrent
- **Version**: 0.1.0
- **Exporter**: Jaeger
- **Sample Rate**: 100%

## 🚀 Deployment

### Docker Compose

```yaml
# docker-compose.yml
services:
  fairrent:
    build: ./services/fairrent
    ports:
      - "50051:50051"
    environment:
      - ALPHA=2.0
      - LOG_LEVEL=info
    healthcheck:
      test: ["CMD", "grpc_health_probe", "-addr=localhost:50051"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### Kubernetes

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fairrent
spec:
  replicas: 3
  selector:
    matchLabels:
      app: fairrent
  template:
    metadata:
      labels:
        app: fairrent
    spec:
      containers:
      - name: fairrent
        image: wohnfair/fairrent:latest
        ports:
        - containerPort: 50051
        env:
        - name: ALPHA
          value: "2.0"
```

## 🔧 Development

### Project Structure

```
services/fairrent/
├── cmd/fairrentd/          # Main application entry point
├── internal/               # Private application code
│   ├── scheduler/          # α-fair scheduling logic
│   ├── queue/              # Priority queue implementation
│   └── telemetry/          # OpenTelemetry setup
├── api/                    # gRPC server implementation
├── config/                 # Configuration files
├── Dockerfile              # Container definition
└── README.md               # This file
```

### Adding New Features

1. **Protocol Buffers**: Update `services/proto/wohnfair/fairrent/v1/fairrent.proto`
2. **Scheduler**: Implement logic in `internal/scheduler/`
3. **API**: Add gRPC methods in `api/server.go`
4. **Tests**: Write unit and integration tests
5. **Documentation**: Update this README

### Code Style

- **Go**: Follow standard Go formatting (`go fmt`)
- **Protocol Buffers**: Use snake_case for field names
- **Logging**: Structured logging with appropriate levels
- **Error Handling**: Wrap errors with context
- **Testing**: Aim for >90% test coverage

## 🐛 Troubleshooting

### Common Issues

#### Service Won't Start

```bash
# Check logs
docker logs fairrent

# Verify port availability
netstat -tulpn | grep 50051

# Check configuration
docker exec fairrent cat /app/config/config.yaml
```

#### High Memory Usage

```bash
# Monitor memory usage
docker stats fairrent

# Check queue size
curl -s http://localhost:8080/metrics | grep queue_length

# Analyze heap profile
go tool pprof -http=:8080 mem.prof
```

#### Poor Performance

```bash
# Check metrics
curl -s http://localhost:8080/metrics | grep processing_duration

# Monitor queue throughput
curl -s http://localhost:8080/metrics | grep requests_processed

# Analyze traces in Jaeger
# http://localhost:16686
```

## 📚 References

- [α-Fairness Paper](https://ieeexplore.ieee.org/document/123456)
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers)
- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
- [Prometheus Metrics](https://prometheus.io/docs/concepts/metric_types/)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This service is part of the WohnFair project and is licensed under the MIT License.
