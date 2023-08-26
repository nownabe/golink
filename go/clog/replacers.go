package clog

import "log/slog"

// https://cloud.google.com/logging/docs/structured-logging#special-payload-fields
func replaceLevelKey(a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		l := a.Value.Any().(slog.Level)

		a.Key = "severity"
		a.Value = slog.StringValue(levelString(l))
	}

	return a
}

func replaceMessageKey(a slog.Attr) slog.Attr {
	if a.Key == slog.MessageKey {
		a.Key = "message"
	}

	return a
}
