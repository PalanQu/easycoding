package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	traceName = "easycoding-trace"
)

func NewTracer() (trace.Tracer, func() error, error) {
	// exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	exporter, err := stdouttrace.New()
	if err != nil {
		return nil, nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{}),
	)
	tpShutdown := func() error {
		if err := tp.Shutdown(context.Background()); err != nil {
			return err
		}
		return nil
	}
	return tp.Tracer(traceName), tpShutdown, nil
}
