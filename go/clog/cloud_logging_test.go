package clog_test

import (
	"context"
	"runtime"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestSourceHandler(t *testing.T) {
	t.Parallel()

	l, w := newLogger()
	ctx := context.Background()

	pc, _, _, _ := runtime.Caller(0) // This file must be just abeve l.Info(ctx, "text")
	l.Info(ctx, "test")

	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()

	want := expectation{
		"time":     anyStringVal{},
		"severity": "INFO",
		"message":  "test",
		"logging.googleapis.com/sourceLocation": map[string]any{
			"file":     frame.File,
			"line":     float64(frame.Line + 1), // json.Unmarshal use float64 for any type
			"function": frame.Function,
		},
	}

	w.expect(t, want)
}

func TestOtelTraceHandler(t *testing.T) {
	t.Parallel()

	l, w := newLogger()
	tracer := sdktrace.NewTracerProvider().Tracer("test")

	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "testspan")
	defer span.End()

	l.Info(ctx, "test")

	want := expectation{
		"time":                                  anyStringVal{},
		"logging.googleapis.com/sourceLocation": nonNilVal{},
		"severity":                              "INFO",
		"message":                               "test",
		"logging.googleapis.com/trace":          span.SpanContext().TraceID().String(),
		"logging.googleapis.com/spanId":         span.SpanContext().SpanID().String(),
	}

	w.expect(t, want)
}
