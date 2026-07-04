package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer(serviceName string, otlpEndpoint string, environment string) (*trace.TracerProvider, func(), error) {
	exporterOptions := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithInsecure(),
	}

	traceExporter, exportError := otlptracehttp.New(context.Background(), exporterOptions...)
	if exportError != nil {
		return nil, nil, exportError
	}

	traceResource, resourceError := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("deployment.environment", environment),
		),
	)
	if resourceError != nil {
		return nil, nil, resourceError
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(5*time.Second)),
		trace.WithResource(traceResource),
	)

	otel.SetTracerProvider(tracerProvider)

	shutdown := func() {
		shutdownContext, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelShutdown()
		tracerProvider.Shutdown(shutdownContext)
		traceExporter.Shutdown(shutdownContext)
	}

	return tracerProvider, shutdown, nil
}
