package clog

import "golang.org/x/exp/slog"

type option func(h slog.Handler) slog.Handler

func WithServiceContext(service, version string) option {
	return option(func(h slog.Handler) slog.Handler {
		h = h.WithAttrs(slog.Group("serviceContext",
			slog.String("service", service),
			slog.String("version", version),
		))
	})
}
