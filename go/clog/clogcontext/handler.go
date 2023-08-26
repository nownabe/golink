package clogcontext

import (
	"context"
	"fmt"
	"sync"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

const (
	spanIDKey = "logging.googleapis.com/spanId"
	traceKey  = "logging.googleapis.com/trace"
	labelsKey = "logging.googleapis.com/labels"
)

func NewHandler(h slog.Handler, projectID string) slog.Handler {
	return &handler{
		Handler:   h,
		projectID: projectID,
	}
}

type handler struct {
	slog.Handler
	projectID string
}

func (h *handler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.Handler.Enabled(ctx, l)
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	if reqID, ok := ctx.Value(keyRequestID{}).(string); ok {
		r.AddAttrs(slog.String("request_id", reqID))
	}

	if spanCtx := trace.SpanContextFromContext(ctx); spanCtx.IsValid() {
		r.AddAttrs(slog.String(traceKey, fmt.Sprintf("projects/%s/traces/%s", h.projectID, spanCtx.TraceID())))
		r.AddAttrs(slog.String(spanIDKey, spanCtx.SpanID().String()))
	}

	if m, ok := ctx.Value(keyLabels{}).(*sync.Map); ok {
		var attrs []slog.Attr
		m.Range(func(key, value any) bool {
			attrs = append(attrs, slog.String(key.(string), value.(string)))
			return true
		})
		r.AddAttrs(slog.Group(labelsKey, attrs))
	}

	return h.Handler.Handle(ctx, r)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{
		Handler:   h.Handler.WithAttrs(attrs),
		projectID: h.projectID,
	}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		Handler:   h.Handler.WithGroup(name),
		projectID: h.projectID,
	}
}
