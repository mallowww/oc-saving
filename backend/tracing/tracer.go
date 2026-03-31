package tracing

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

/*
idea

	context -> exporter -> Jaeger
	traceProvider -> resource -> serviceName
*/
func InitTracerAbstraced(serviceName, serviceVersion, deployEnv, jaegerEndpoint string) (func(context.Context) error, error) {
	ctx := context.Background()

	exporter, err := NewJaegerExporter(ctx, jaegerEndpoint)
	if err != nil {
		return nil, err
	}

	traceProvider, err := NewTraceProvider(ctx, serviceName, serviceVersion, deployEnv, exporter)
	if err != nil {
		return nil, err
	}

	SetGlobalTraceProvider(traceProvider)
	log.Println("tracer init")

	// cleanup := func(ctx context.Context) error {
	// 	return traceProvider.Shutdown(ctx)
	// }
	cleanup := func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := traceProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("fail to shutdown tracerProvider: %w", err)
		}
		log.Println("tracer shutdown")
		return nil
	}

	return cleanup, nil
}

func NewJaegerExporter(ctx context.Context, endpoint string) (sdktrace.SpanExporter, error) {
	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("fail to create OLTP exporter: %w", err)
	}
	return exporter, nil
}

func NewTraceProvider(ctx context.Context, serviceName, serviceVersion, deployEnv string, exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
			semconv.DeploymentEnvironment(deployEnv),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("fail to create resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	return tp, nil
}

func SetGlobalTraceProvider(tp *sdktrace.TracerProvider) {
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}
