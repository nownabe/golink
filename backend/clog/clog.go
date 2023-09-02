package clog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync/atomic"

	"github.com/nownabe/golink/backend/clog/clogcontext"
)

var defaultLogger atomic.Value

func init() {
	defaultLogger.Store(New(os.Stdout, LevelInfo))
}

func Default() *Logger {
	return defaultLogger.Load().(*Logger)
}

func SetDefault(l *Logger) {
	defaultLogger.Store(l)
}

func SetContextHandler(projectID string) {
	l := Default()
	h := clogcontext.NewHandler(l.Handler(), projectID)
	SetDefault(&Logger{slog.New(h)})
}

func New(w io.Writer, l slog.Level, opts ...option) *Logger {
	jh := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: l,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			a = replaceLevelKey(a)
			a = replaceMessageKey(a)

			return a
		},
	})
	h := slog.Handler(jh)

	for _, opt := range opts {
		h = opt(h)
	}

	return &Logger{slog.New(h)}
}

func Debug(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, LevelDebug, msg, args...)
}

func Debugf(ctx context.Context, format string, args ...any) {
	Default().log(ctx, LevelDebug, fmt.Sprintf(format, args...))
}

func Info(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, LevelInfo, msg, args...)
}

func Infof(ctx context.Context, format string, args ...any) {
	Default().log(ctx, LevelInfo, fmt.Sprintf(format, args...))
}

func Notice(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, LevelNotice, msg, args...)
}

func Noticef(ctx context.Context, format string, args ...any) {
	Default().log(ctx, LevelNotice, fmt.Sprintf(format, args...))
}

func Warning(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, LevelWarning, msg, args...)
}

func Warningf(ctx context.Context, format string, args ...any) {
	Default().log(ctx, LevelWarning, fmt.Sprintf(format, args...))
}

func Error(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, LevelError, msg, args...)
}

func Errorf(ctx context.Context, format string, args ...any) {
	Default().log(ctx, LevelError, fmt.Sprintf(format, args...))
}

func Critical(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, LevelCritical, msg, args...)
}

func Criticalf(ctx context.Context, format string, args ...any) {
	Default().log(ctx, LevelCritical, fmt.Sprintf(format, args...))
}

func Alert(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, LevelAlert, msg, args...)
}

func Alertf(ctx context.Context, format string, args ...any) {
	Default().log(ctx, LevelAlert, fmt.Sprintf(format, args...))
}

func Emergency(ctx context.Context, msg string, args ...any) {
	Default().log(ctx, LevelEmergency, msg, args...)
}

func Emergencyf(ctx context.Context, format string, args ...any) {
	Default().log(ctx, LevelEmergency, fmt.Sprintf(format, args...))
}

func Err(ctx context.Context, err error) {
	Default().err(ctx, LevelError, err)
}

func WarningErr(ctx context.Context, err error) {
	Default().err(ctx, LevelWarning, err)
}

func CriticalErr(ctx context.Context, err error) {
	Default().err(ctx, LevelCritical, err)
}

func AlertErr(ctx context.Context, err error) {
	Default().err(ctx, LevelAlert, err)
}

func EmergencyErr(ctx context.Context, err error) {
	Default().err(ctx, LevelEmergency, err)
}

func InfoHTTPRequest(ctx context.Context, msg string, req *HTTPRequest) {
	Default().log(ctx, LevelInfo, msg, httpRequestKey, req)
}
