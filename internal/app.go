package internal

import (
	"context"
	"log/slog"

	"github.com/novychok/authasvs/internal/handler/authapiv1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type App struct {
	l *slog.Logger

	authApiV1 *authapiv1.Server
}

func (a *App) StartAuthApiV1(ctx context.Context) error {
	return a.authApiV1.Start(ctx)
}

func startTracing(ctx context.Context) (*trace.TracerProvider, error) {
	serviceName := "auth"

	traceExporter, err := otlptrace.New(ctx, otlptracehttp.NewClient())
	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			),
		),
	)

	otel.SetTracerProvider(tracerProvider)

	return tracerProvider, nil
}

func New(
	ctx context.Context,
	l *slog.Logger,
	authApiV1 *authapiv1.Server,
) (*App, func(), error) {
	app := &App{
		l:         l,
		authApiV1: authApiV1,
	}

	tracer, err := startTracing(ctx)
	if err != nil {
		return nil, func() {}, err
	}

	cleanup := func() {
		_ = tracer.Shutdown(ctx)
	}

	return app, cleanup, nil
}
