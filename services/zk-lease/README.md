# ZK-Lease Service

The ZK-Lease service provides zero-knowledge proof generation and verification for housing lease verification in the WohnFair system. It implements multiple proof systems including Halo2, PLONK, and Groth16 to ensure privacy-preserving verification of housing eligibility and compliance.

## Architecture

The service is built with a modular architecture:

- **Prover Module**: Generates zero-knowledge proofs for various housing-related claims
- **Verifier Module**: Verifies proofs without revealing sensitive information
- **Circuit Management**: Handles circuit compilation and key management
- **gRPC Interface**: Provides a standardized API for proof operations
- **Metrics & Observability**: Comprehensive monitoring and tracing

## Features

- **Multiple Proof Systems**: Support for Halo2, PLONK, and Groth16
- **Circuit Compilation**: Dynamic circuit loading and compilation
- **Proof Caching**: Intelligent caching for performance optimization
- **Batch Processing**: Efficient batch proof generation and verification
- **Security**: Key rotation, proof freshness validation, and circuit allowlisting
- **Observability**: Prometheus metrics, Jaeger tracing, and structured logging

## Quick Start

### Prerequisites

- Rust 1.75+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose

### Local Development

1. **Clone and setup**:
   ```bash
   cd services/zk-lease
   cargo build
   ```

2. **Environment setup**:
   ```bash
   export DATABASE_URL="postgresql://wohnfair:wohnfair_pass@localhost:5432/wohnfair"
   export REDIS_URL="redis://:redis_pass@localhost:6379"
   export JAEGER_ENDPOINT="http://localhost:14268/api/traces"
   ```

3. **Run the service**:
   ```bash
   cargo run
   ```

### Docker

```bash
# Build the image
docker build -t zk-lease .

# Run the container
docker run -p 50052:50052 \
  -e DATABASE_URL="postgresql://wohnfair:wohnfair_pass@host.docker.internal:5432/wohnfair" \
  -e REDIS_URL="redis://:redis_pass@host.docker.internal:6379" \
  zk-lease
```

## API Reference

### gRPC Methods

#### ProveReservation
Generates a zero-knowledge proof for a housing reservation.

```protobuf
rpc ProveReservation(ProveReservationRequest) returns (ProveReservationResponse);
```

**Request**:
- `circuit_type`: Type of circuit to use (e.g., "housing_eligibility")
- `public_inputs`: Public inputs for the proof
- `private_inputs`: Private inputs (encrypted)
- `metadata`: Additional proof metadata

**Response**:
- `proof`: Generated zero-knowledge proof
- `public_outputs`: Public outputs from the proof
- `verification_key_hash`: Hash of the verification key used

#### VerifyReservation
Verifies a zero-knowledge proof without revealing private inputs.

```protobuf
rpc VerifyReservation(VerifyReservationRequest) returns (VerifyReservationResponse);
```

**Request**:
- `proof`: The zero-knowledge proof to verify
- `public_inputs`: Public inputs used in the proof
- `verification_key_hash`: Hash of the verification key

**Response**:
- `is_valid`: Whether the proof is valid
- `verification_details`: Additional verification information
- `metadata`: Proof metadata

#### GenerateProofParameters
Generates parameters for proof generation.

```protobuf
rpc GenerateProofParameters(GenerateProofParametersRequest) returns (GenerateProofParametersResponse);
```

#### GetProofStatus
Retrieves the status of a proof generation or verification request.

```protobuf
rpc GetProofStatus(GetProofStatusRequest) returns (GetProofStatusResponse);
```

### HTTP Endpoints

- `GET /healthz` - Health check
- `GET /metrics` - Prometheus metrics
- `GET /status` - Service status

## Circuit Types

### Housing Eligibility Circuit
Proves eligibility for housing without revealing personal financial information.

**Public Inputs**:
- Age category
- Disability status
- Refugee status
- City of residence

**Private Inputs**:
- Income details
- Asset information
- Employment history
- Family composition

### Income Verification Circuit
Verifies income meets requirements without revealing exact amounts.

**Public Inputs**:
- Income threshold
- Verification timestamp

**Private Inputs**:
- Actual income
- Source of income
- Tax documents

### Residency Proof Circuit
Proves residency in a specific area without revealing exact address.

**Public Inputs**:
- City/region
- Residency duration

**Private Inputs**:
- Exact address
- Utility bills
- Rental agreements

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ZK_LEASE_HOST` | `0.0.0.0` | Service host |
| `ZK_LEASE_PORT` | `50052` | Service port |
| `ZK_LEASE_WORKERS` | `4` | Number of worker threads |
| `DATABASE_URL` | PostgreSQL connection string | Database connection |
| `REDIS_URL` | Redis connection string | Redis connection |
| `CIRCUIT_PATH` | `/app/circuits` | Path to circuit files |
| `PROVING_KEY_PATH` | `/app/keys/proving.key` | Path to proving key |
| `VERIFYING_KEY_PATH` | `/app/keys/verifying.key` | Path to verifying key |
| `JAEGER_ENDPOINT` | Jaeger endpoint | Tracing endpoint |

### Configuration File

The service supports YAML configuration files. See `config/config.yaml` for the complete configuration schema.

## Performance

### Benchmarks

| Circuit | Constraints | Proof Generation | Verification |
|---------|-------------|------------------|--------------|
| Housing Eligibility | 10,000 | ~2.5s | ~0.1s |
| Income Verification | 5,000 | ~1.2s | ~0.05s |
| Residency Proof | 3,000 | ~0.8s | ~0.03s |

### Optimization Features

- **Proof Caching**: Caches generated proofs for reuse
- **Batch Processing**: Processes multiple proofs simultaneously
- **Circuit Optimization**: Optimized constraint systems
- **Parallel Verification**: Parallel proof verification

## Security

### Key Management

- **Automatic Rotation**: Keys are automatically rotated
- **Secure Storage**: Keys are stored securely with access controls
- **Audit Logging**: All key operations are logged

### Proof Security

- **Freshness Validation**: Proofs have time-based validity
- **Circuit Allowlisting**: Only approved circuits can be used
- **Input Validation**: Strict validation of all inputs
- **Cryptographic Verification**: All proofs are cryptographically verified

## Monitoring

### Metrics

The service exposes Prometheus metrics:

- `zk_lease_proofs_generated_total`: Total proofs generated
- `zk_lease_proofs_verified_total`: Total proofs verified
- `zk_lease_proof_generation_duration`: Proof generation time
- `zk_lease_proof_verification_duration`: Proof verification time
- `zk_lease_circuit_compilation_duration`: Circuit compilation time

### Tracing

OpenTelemetry tracing is supported with Jaeger integration:

- Proof generation spans
- Circuit compilation spans
- Database operation spans
- gRPC request spans

### Health Checks

- Database connectivity
- Redis connectivity
- Circuit availability
- Key availability

## Development

### Project Structure

```
src/
├── main.rs              # Service entry point
├── config.rs            # Configuration management
├── error.rs             # Error handling
├── grpc/                # gRPC service implementation
├── prover/              # Proof generation logic
├── verifier/            # Proof verification logic
├── circuits/            # Circuit implementations
├── utils/               # Utility functions
└── metrics.rs           # Metrics and observability
```

### Adding New Circuits

1. **Implement the circuit** in `src/circuits/`
2. **Add circuit configuration** in `config/config.yaml`
3. **Update the prover** to handle the new circuit type
4. **Add tests** for the new circuit
5. **Update documentation**

### Testing

```bash
# Run all tests
cargo test

# Run specific test
cargo test test_proof_generation

# Run benchmarks
cargo bench

# Run with coverage
cargo tarpaulin
```

### Code Quality

```bash
# Format code
cargo fmt

# Lint code
cargo clippy

# Check for security issues
cargo audit
```

## Deployment

### Docker Compose

The service is included in the main `docker-compose.yml`:

```yaml
zk-lease:
  build:
    context: ./services/zk-lease
    dockerfile: Dockerfile
  ports:
    - "50052:50052"
  environment:
    - JAEGER_ENDPOINT=http://jaeger:14268/api/traces
    - PROMETHEUS_ENDPOINT=http://prometheus:9090
  networks:
    - wohnfair-backend
```

### Kubernetes

See `infra/k8s/` for Kubernetes deployment manifests.

### Production Considerations

- **Resource Limits**: Set appropriate CPU/memory limits
- **Scaling**: Use horizontal pod autoscaling
- **Monitoring**: Enable comprehensive monitoring
- **Backup**: Regular backup of keys and configurations
- **Security**: Use secrets management for sensitive data

## Troubleshooting

### Common Issues

1. **Circuit Compilation Failures**
   - Check circuit file syntax
   - Verify constraint count limits
   - Check available memory

2. **Proof Generation Timeouts**
   - Increase timeout configuration
   - Check system resources
   - Verify circuit complexity

3. **Verification Failures**
   - Check proof format
   - Verify verification key
   - Check input validation

### Debug Mode

Enable debug logging:

```bash
export RUST_LOG=debug
cargo run
```

### Performance Profiling

Enable profiling:

```bash
cargo build --release
perf record --call-graph=dwarf ./target/release/zk-lease
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

- [Halo2 Documentation](https://zcash.github.io/halo2/)
- [PLONK Protocol](https://eprint.iacr.org/2019/953)
- [Groth16 Protocol](https://eprint.iacr.org/2016/260)
- [Zero-Knowledge Proofs](https://en.wikipedia.org/wiki/Zero-knowledge_proof)
- [WohnFair Architecture](../../README.md#architecture)
