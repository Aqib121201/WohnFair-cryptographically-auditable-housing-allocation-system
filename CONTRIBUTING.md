# Contributing to WohnFair

Thank you for your interest in contributing to WohnFair! This document provides guidelines for contributing to the project.

## 🚀 Quick Start

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Run tests: `make test`
5. Commit with conventional format: `git commit -m "feat: add amazing feature"`
6. Push to your fork: `git push origin feature/amazing-feature`
7. Open a Pull Request

## 📝 Commit Convention

We use [Conventional Commits](https://www.conventionalcommits.org/) for commit messages:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, missing semicolons, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples
```
feat(fairrent): implement α-fair scheduling algorithm
fix(gateway): resolve CORS issue in development
docs: update API documentation with examples
test(zk-lease): add benchmark tests for proof generation
```

## 🛠️ Development Setup

### Prerequisites
- Go 1.21+
- Rust 1.70+
- Node.js 18+
- Python 3.11+
- Docker & Docker Compose

### Local Development
```bash
# Clone and setup
git clone https://github.com/wohnfair/wohnfair.git
cd wohnfair

# Install dependencies
make deps

# Generate protocol buffers
make proto

# Run tests
make test

# Start services locally
make dev

# Build and run with Docker
make compose-up
```

### Service-Specific Development

#### Go Services
```bash
cd services/fairrent
go mod tidy
go test ./...
go run cmd/fairrentd/main.go
```

#### Rust Services
```bash
cd services/zk-lease
cargo build
cargo test
cargo run
```

#### Python ML
```bash
cd services/ml
pip install -e .
pytest
python -m wohnfair_ml.cli train --help
```

#### Frontend
```bash
cd frontend
npm install
npm run dev
npm run test
```

## 🧪 Testing

### Running Tests
```bash
# All tests
make test

# Service-specific tests
make test-go
make test-rust
make test-python
make test-js

# Integration tests
make test-integration

# Performance tests
make test-bench
```

### Test Coverage
```bash
# Generate coverage reports
make coverage

# View coverage in browser
make coverage-html
```

## 🔍 Code Quality

### Linting
```bash
# All linting
make lint

# Service-specific linting
make lint-go
make lint-rust
make lint-python
make lint-js
```

### Formatting
```bash
# Format all code
make format

# Check formatting
make format-check
```

## 📋 Pull Request Process

1. **Update Documentation**: Ensure all public APIs are documented
2. **Add Tests**: Include tests for new functionality
3. **Update CHANGELOG**: Add entry to `CHANGELOG.md`
4. **Check CI**: Ensure all CI checks pass
5. **Request Review**: Assign appropriate reviewers

### PR Checklist
- [ ] Tests pass locally
- [ ] Code follows style guidelines
- [ ] Documentation is updated
- [ ] CHANGELOG entry added
- [ ] No breaking changes (or documented as such)

## 🐛 Bug Reports

When reporting bugs, please include:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, versions)
- Relevant logs or error messages

## 💡 Feature Requests

For feature requests:
- Describe the use case
- Explain why it's needed
- Suggest implementation approach if possible
- Consider impact on existing functionality

## 📚 Additional Resources

- [Architecture Documentation](docs/architecture.md)
- [API Reference](docs/api.md)
- [Development Guidelines](docs/development.md)
- [Testing Strategy](docs/testing.md)

## 🤝 Community

- **Discussions**: Use GitHub Discussions for questions and ideas
- **Issues**: Report bugs and request features via GitHub Issues
- **Security**: Report security issues to security@wohnfair.org

## 📄 License

By contributing to WohnFair, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to fair housing allocation! 🏠✨
