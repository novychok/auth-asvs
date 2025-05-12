package slog

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

type TracingHandler struct {
	slog.Handler
	tracer trace.Tracer
}

func (h *TracingHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		r.Add(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}

	return h.Handler.Handle(ctx, r)
}

func NewTracingHandler(tracer trace.Tracer) *TracingHandler {
	return &TracingHandler{
		Handler: slog.NewJSONHandler(os.Stdout, nil),
		tracer:  tracer,
	}
}

func New() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		// AddSource: true,
		// Level:     slog.LevelDebug,
	}))
}
