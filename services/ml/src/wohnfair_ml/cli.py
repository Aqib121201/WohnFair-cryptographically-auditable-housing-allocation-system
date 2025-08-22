import os
import sys
import time
import click
from fastapi import FastAPI
from fastapi.responses import PlainTextResponse
from prometheus_client import Counter, Histogram, generate_latest, CONTENT_TYPE_LATEST
import uvicorn

REQUEST_COUNTER = Counter("wohnfair_ml_requests_total", "Total API requests", ["endpoint", "method"])
REQUEST_LATENCY = Histogram("wohnfair_ml_request_duration_seconds", "Request latency", ["endpoint"])

app = FastAPI(title="WohnFair ML Service", version="0.1.0")

@app.get("/healthz", response_class=PlainTextResponse)
def healthz() -> str:
    return "ok\n"

@app.get("/metrics")
def metrics():
    return PlainTextResponse(generate_latest(), media_type=CONTENT_TYPE_LATEST)

@app.get("/status", response_class=PlainTextResponse)
def status() -> str:
    return "running\n"

@click.group()
def main():
    """WohnFair ML command line interface."""
    pass

@main.command()
@click.option("--host", default=os.getenv("SERVICE_HOST", "0.0.0.0"))
@click.option("--port", default=int(os.getenv("SERVICE_PORT", 8000)))
@click.option("--workers", default=int(os.getenv("SERVICE_WORKERS", 1)))
def serve(host: str, port: int, workers: int):
    """Start FastAPI server for inference/health/metrics."""
    uvicorn.run("wohnfair_ml.cli:app", host=host, port=port, workers=workers, log_level="info")

@main.command()
@click.option("--config", type=str, default="config/config.yaml")
def train(config: str):
    """Stub training command."""
    click.echo(f"Starting training with config: {config}")
    time.sleep(1)
    click.echo("Training complete (stub)")

@main.command()
@click.option("--model", type=str, required=False)
@click.option("--data", type=str, required=False)
def evaluate(model: str|None, data: str|None):
    """Stub evaluation command."""
    click.echo(f"Evaluating model={model} on data={data}")
    time.sleep(1)
    click.echo("Evaluation complete (stub)")

if __name__ == "__main__":
    main()
