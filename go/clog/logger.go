package clog

import (
	"context"
	"fmt"

	"golang.org/x/exp/slog"
)

var (
	logger *Logger
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{l.Logger.With(args...)}
}

// TODO
func (l *Logger) WithErr(err error) *Logger {
	return l.With("error", err)
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, LevelDebug, msg, args...)
}

func (l *Logger) Debugf(ctx context.Context, format string, args ...any) {
	l.Log(ctx, LevelDebug, fmt.Sprintf(format, args...))
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, LevelInfo, msg, args...)
}

func (l *Logger) Infof(ctx context.Context, format string, args ...any) {
	l.Log(ctx, LevelInfo, fmt.Sprintf(format, args...))
}

func (l *Logger) Notice(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, LevelNotice, msg, args...)
}

func (l *Logger) Noticef(ctx context.Context, format string, args ...any) {
	l.Log(ctx, LevelNotice, fmt.Sprintf(format, args...))
}

func (l *Logger) Warning(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, LevelWarning, msg, args...)
}

func (l *Logger) Warningf(ctx context.Context, format string, args ...any) {
	l.Log(ctx, LevelWarning, fmt.Sprintf(format, args...))
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, LevelError, msg, args...)
}

func (l *Logger) Errorf(ctx context.Context, format string, args ...any) {
	l.Log(ctx, LevelError, fmt.Sprintf(format, args...))
}

func (l *Logger) Critical(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, LevelCritical, msg, args...)
}

func (l *Logger) Criticalf(ctx context.Context, format string, args ...any) {
	l.Log(ctx, LevelCritical, fmt.Sprintf(format, args...))
}

func (l *Logger) Alert(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, LevelAlert, msg, args...)
}

func (l *Logger) Alertf(ctx context.Context, format string, args ...any) {
	l.Log(ctx, LevelAlert, fmt.Sprintf(format, args...))
}

func (l *Logger) Emergency(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, LevelEmergency, msg, args...)
}

func (l *Logger) Emergencyf(ctx context.Context, format string, args ...any) {
	l.Log(ctx, LevelEmergency, fmt.Sprintf(format, args...))
}

func (l *Logger) Err(ctx context.Context, err error, args ...any) {
	l.WithErr(err).Log(ctx, LevelError, err.Error(), args...)
}
