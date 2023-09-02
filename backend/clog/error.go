package clog

import "log/slog"

// ErrorEvent is an interface for Error Reporting.
// See https://cloud.google.com/error-reporting/docs/formatting-error-messages#log-error
type ErrorEvent interface {
	ErrorContext() *ErrorContext

	// Stack must be return value of runtime.Stack or debug.Stack.
	// See https://pkg.go.dev/runtime/debug#Stack
	Stack() []byte
}

type ServiceContext struct {
	Service string
	Version string
}

func (c *ServiceContext) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("service", c.Service),
		slog.String("version", c.Version),
	)
}

type ErrorContext struct {
	HTTPRequest      *HTTPRequestContext
	User             string
	ReportLocation   *ReportLocation
	SourceReferences []*SourceReference
}

// TODO.
func (c *ErrorContext) LogValue() slog.Value {
	return slog.GroupValue()
}

// TODO.
func (c *ErrorContext) LogAttr() slog.Attr {
	return slog.Group("context", slog.String("dummy", "dummy"))
}

// TODO.
type (
	HTTPRequestContext struct{}
	ReportLocation     struct{}
	SourceReference    struct{}
)
