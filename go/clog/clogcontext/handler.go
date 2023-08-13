package clogcontext

import (
	"context"

	"golang.org/x/exp/slog"
)

func NewHandler(h slog.Handler) slog.Handler {
	return &handler{h}
}

type handler struct {
	slog.Handler
}

func (h *handler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.Handler.Enabled(ctx, l)
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	if reqID, ok := ctx.Value(keyRequestID{}).(string); ok {
		r.AddAttrs(slog.String("request_id", reqID))
	}

	return h.Handler.Handle(ctx, r)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{h.Handler.WithAttrs(attrs)}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{h.Handler.WithGroup(name)}
}
