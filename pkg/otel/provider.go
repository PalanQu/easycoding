package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	traceName = "easycoding-trace"
)

func NewTracer() (trace.Tracer, func() error, error) {
	exporter, err := newStdoutExporter(false)
	if err != nil {
		return nil, nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("easycoding"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	tpShutdown := func() error {
		if err := tp.Shutdown(context.Background()); err != nil {
			return err
		}
		return nil
	}
	return tp.Tracer(traceName), tpShutdown, nil
}

func newStdoutExporter(pretty bool) (sdktrace.SpanExporter, error) {
	if pretty {
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	}
	return stdouttrace.New()
}

func newJaegerExporter() (sdktrace.SpanExporter, error) {
	// TODO(qujiabao): mv the hard code trace address into config
	url := "http://localhost:14268/api/traces"
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
}
