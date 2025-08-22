# Multi-stage build for WohnFair ML service
FROM python:3.11-slim as builder

# Install system dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    libpq-dev \
    libssl-dev \
    libffi-dev \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy requirements and install Python dependencies
COPY pyproject.toml ./
COPY src/ ./src/

# Install Python dependencies
RUN pip install --no-cache-dir --upgrade pip setuptools wheel
RUN pip install --no-cache-dir -e .

# Runtime stage
FROM python:3.11-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    libpq5 \
    libssl3 \
    libffi8 \
    curl \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd --create-home --shell /bin/bash wohnfair

# Set working directory
WORKDIR /app

# Copy Python packages from builder
COPY --from=builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin

# Copy source code
COPY src/ ./src/

# Copy configuration files
COPY config/ ./config/

# Create necessary directories
RUN mkdir -p /app/data /app/models /app/artifacts /app/logs /app/notebooks

# Set ownership
RUN chown -R wohnfair:wohnfair /app

# Switch to non-root user
USER wohnfair

# Set environment variables
ENV PYTHONPATH=/app/src
ENV PYTHONUNBUFFERED=1
ENV PYTHONDONTWRITEBYTECODE=1

# Expose ports
EXPOSE 8000 9092

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD ["curl", "-f", "http://localhost:8000/healthz"] || exit 1

# Default command
CMD ["python", "-m", "wohnfair_ml.cli", "serve"]
