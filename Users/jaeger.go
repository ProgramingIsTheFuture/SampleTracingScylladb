package main

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func initJaeger(serviceName, jaegerAgentEndpoint string) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(fmt.Sprintf("%s/api/traces", jaegerAgentEndpoint))))
	if err != nil {
		return
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
				attribute.String("environment", "production"),
				attribute.String("ID", "1"),
			),
		),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

}
