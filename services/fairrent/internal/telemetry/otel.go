package telemetry

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// InitTracer initializes OpenTelemetry tracing
func InitTracer(serviceName, serviceVersion string) error {
	// Get Jaeger endpoint from environment
	jaegerEndpoint := os.Getenv("JAEGER_ENDPOINT")
	if jaegerEndpoint == "" {
		jaegerEndpoint = "http://localhost:14268/api/traces"
	}

	// Create Jaeger exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		return fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	// Create resource with service information
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
			semconv.DeploymentEnvironment("development"),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(5*time.Second),
			sdktrace.WithMaxExportBatchSize(100),
		),
		sdktrace.WithResource(res),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)

	return nil
}

// GetTracer returns a tracer for the given service
func GetTracer(serviceName string) trace.Tracer {
	return otel.Tracer(serviceName)
}

// StartSpan starts a new span with the given name and options
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	tracer := GetTracer("fairrent")
	return tracer.Start(ctx, name, opts...)
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(ctx context.Context, name string, attrs ...trace.SpanStartOption) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.AddEvent(name)
	}
}

// SetSpanAttributes sets attributes on the current span
func SetSpanAttributes(ctx context.Context, attrs ...trace.SpanStartOption) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		// This would need to be implemented based on the specific attributes
		// For now, we'll just add a generic event
		span.AddEvent("attributes_set")
	}
}

// RecordError records an error on the current span
func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
	}
}

// SetSpanStatus sets the status on the current span
func SetSpanStatus(ctx context.Context, code trace.StatusCode, description string) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.SetStatus(code, description)
	}
}
