package clog

import "log/slog"

// Levels describe the severity of the log.
// They follow definitions of Cloud Logging.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
const (
	// LevelDefault means the log entry has no assigned severity level.
	LevelDefault = slog.Level(0)

	// LevelDebug means debug or trace information.
	LevelDebug = slog.Level(100)

	// LevelInfo means routine information, such as ongoing status or performance.
	LevelInfo = slog.Level(200)

	// LevelNotice means normal but significant events, such as start up, shut down, or a configuration change.
	LevelNotice = slog.Level(300)

	// LevelWarning means warning events might cause problems.
	LevelWarning = slog.Level(400)

	// LevelError means error events are likely to cause problems.
	LevelError = slog.Level(500)

	// LevelCritical means critical events cause more severe problems or outages.
	LevelCritical = slog.Level(600)

	// LevelAlert means a person must take an action immediately.
	LevelAlert = slog.Level(700)

	// LevelEmergency means one or more systems are unusable.
	LevelEmergency = slog.Level(800)
)

func levelString(l slog.Level) string {
	switch l {
	case LevelDefault:
		return "DEFAULT"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelNotice:
		return "NOTICE"
	case LevelWarning:
		return "WARNING"
	case LevelError:
		return "ERROR"
	case LevelCritical:
		return "CRITICAL"
	case LevelAlert:
		return "ALERT"
	case LevelEmergency:
		return "EMERGENCY"
	}

	return "DEFAULT"
}
