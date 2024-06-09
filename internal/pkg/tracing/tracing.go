package tracing

// tracing using open telemetry and jaeger
import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

// using jaeger with jaeger url eg: "http://localhost:14268/api/traces"
func InitTracer(ctx context.Context, serviceName, serviceVersion, jaegerTraceUrl string) (*trace.TracerProvider, error) {
	// Create the exporter
	client := otlptracehttp.NewClient(otlptracehttp.WithEndpoint(jaegerTraceUrl), otlptracehttp.WithInsecure(), otlptracehttp.WithCompression(otlptracehttp.NoCompression))
	traceExporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}

	// Create the resource to be traced
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		),
	)
	if err != nil {
		return nil, err
	}

	// Configure the trace provider
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(2*time.Second)),
		trace.WithResource(res),
	)

	return traceProvider, nil
}
