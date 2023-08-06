package clog

import "golang.org/x/exp/slog"

type option func(h slog.Handler) slog.Handler

func WithServiceContext(service, version string) option {
	return option(func(h slog.Handler) slog.Handler {
		return h.WithAttrs([]slog.Attr{slog.Group("serviceContext",
			slog.String("service", service),
			slog.String("version", version),
		)})
	})
}
