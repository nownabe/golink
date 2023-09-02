package clog

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nownabe/golink/backend/clog/clogcontext"
	"go.opentelemetry.io/otel/trace"
)

func newContextHandler(h slog.Handler, projectID string) slog.Handler {
	return &contextHandler{
		Handler:   h,
		projectID: projectID,
	}
}

type contextHandler struct {
	slog.Handler
	projectID string
}

func (h *contextHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.Handler.Enabled(ctx, l)
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if reqID, ok := clogcontext.RequestIDFrom(ctx); ok {
		r.AddAttrs(slog.String("request_id", reqID))
	}

	if spanCtx := trace.SpanContextFromContext(ctx); spanCtx.IsValid() {
		r.AddAttrs(slog.String(traceKey, fmt.Sprintf("projects/%s/traces/%s", h.projectID, spanCtx.TraceID())))
		r.AddAttrs(slog.String(spanIDKey, spanCtx.SpanID().String()))
	}

	if m, ok := clogcontext.LabelFrom(ctx); ok {
		var attrs []any
		m.Range(func(key, value any) bool {
			attrs = append(attrs, slog.String(key.(string), value.(string)))
			return true
		})
		r.AddAttrs(slog.Group(labelsKey, attrs...))
	}

	return h.Handler.Handle(ctx, r)
}

func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextHandler{
		Handler:   h.Handler.WithAttrs(attrs),
		projectID: h.projectID,
	}
}

func (h *contextHandler) WithGroup(name string) slog.Handler {
	return &contextHandler{
		Handler:   h.Handler.WithGroup(name),
		projectID: h.projectID,
	}
}
