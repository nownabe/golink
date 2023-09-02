package clog

import (
	"context"
	"fmt"
	"log/slog"
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) Handler() slog.Handler {
	return l.Logger.Handler()
}

func (l *Logger) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	// skip [runtime.Callers, getSource, this function, clog exported functions]
	s := getSource(4)
	l.logWithSource(ctx, level, msg, s, args...)
}

func (l *Logger) logWithSource(ctx context.Context, level slog.Level, msg string, s *slog.Source, args ...any) {
	args = append(args, sourceLocationKey, s)
	l.Logger.Log(ctx, level, msg, args...)
}

func (l *Logger) logAttrsWithSource(ctx context.Context, level slog.Level, msg string, s *slog.Source, attrs ...slog.Attr) {
	attrs = append(attrs, slog.Any(sourceLocationKey, s))
	l.Logger.LogAttrs(ctx, level, msg, attrs...)
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelDebug, msg, args...)
}

func (l *Logger) Debugf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LevelDebug, fmt.Sprintf(format, args...))
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelInfo, msg, args...)
}

func (l *Logger) Infof(ctx context.Context, format string, args ...any) {
	l.log(ctx, LevelInfo, fmt.Sprintf(format, args...))
}

func (l *Logger) Notice(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelNotice, msg, args...)
}

func (l *Logger) Noticef(ctx context.Context, format string, args ...any) {
	l.log(ctx, LevelNotice, fmt.Sprintf(format, args...))
}

func (l *Logger) Warning(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelWarning, msg, args...)
}

func (l *Logger) Warningf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LevelWarning, fmt.Sprintf(format, args...))
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelError, msg, args...)
}

func (l *Logger) Errorf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LevelError, fmt.Sprintf(format, args...))
}

func (l *Logger) Critical(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelCritical, msg, args...)
}

func (l *Logger) Criticalf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LevelCritical, fmt.Sprintf(format, args...))
}

func (l *Logger) Alert(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelAlert, msg, args...)
}

func (l *Logger) Alertf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LevelAlert, fmt.Sprintf(format, args...))
}

func (l *Logger) Emergency(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelEmergency, msg, args...)
}

func (l *Logger) Emergencyf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LevelEmergency, fmt.Sprintf(format, args...))
}

func (l *Logger) Err(ctx context.Context, err error) {
	l.err(ctx, LevelError, err)
}

func (l *Logger) WarningErr(ctx context.Context, err error) {
	l.err(ctx, LevelWarning, err)
}

func (l *Logger) CriticalErr(ctx context.Context, err error) {
	l.err(ctx, LevelCritical, err)
}

func (l *Logger) AlertErr(ctx context.Context, err error) {
	l.err(ctx, LevelAlert, err)
}

func (l *Logger) EmergencyErr(ctx context.Context, err error) {
	l.err(ctx, LevelEmergency, err)
}

func (l *Logger) err(ctx context.Context, lv slog.Level, err error) {
	var attrs []slog.Attr

	if ee, ok := err.(ErrorEvent); ok {
		attrs = append(attrs, slog.String("@type", "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"))

		if ec := ee.ErrorContext(); ec != nil {
			attrs = append(attrs, slog.Group("context", ec.LogAttr()))
		}

		if s := ee.Stack(); s != nil {
			attrs = append(attrs, slog.String("stack_trace", string(s)))
		}
	}

	// skip [runtime.Callers, getSource, this function, clog exported functions]
	s := getSource(4)
	l.logAttrsWithSource(ctx, lv, err.Error(), s, attrs...)
}

func (l *Logger) InfoHTTPRequest(ctx context.Context, msg string, req *HTTPRequest) {
	l.log(ctx, LevelInfo, msg, httpRequestKey, req)
}
