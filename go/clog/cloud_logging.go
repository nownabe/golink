package clog

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
)

const (
	httpRequestKey    = "httpRequest"
	insertIDKey       = "logging.googleapis.com/insertId"
	labelsKey         = "logging.googleapis.com/labels"
	operationKey      = "logging.googleapis.com/operation"
	sourceLocationKey = "logging.googleapis.com/sourceLocation"
	spanIDKey         = "logging.googleapis.com/spanId"
	traceKey          = "logging.googleapis.com/trace"
	traceSampledKey   = "logging.googleapis.com/trace_sampled"
)

// HTTPRequest is a log value for HTTP request.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#HttpRequest
type HTTPRequest struct {
	RequestMethod                  string
	RequestURL                     string
	RequestSize                    string
	Status                         string
	ResponseSize                   string
	UserAgent                      string
	RemoteIP                       string
	ServerIP                       string
	Referer                        string
	Latency                        string
	CacheLookup                    bool
	CacheHit                       bool
	CacheValidatedWithOriginServer bool
	CacheFillBytes                 string
	Protocol                       string
}

func (r *HTTPRequest) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("requestMethod", r.RequestMethod),
		slog.String("requestUrl", r.RequestURL),
		slog.String("requestSize", r.RequestSize),
		slog.String("status", r.Status),
		slog.String("responseSize", r.ResponseSize),
		slog.String("userAgent", r.UserAgent),
		slog.String("remoteIp", r.RemoteIP),
		slog.String("serverIp", r.ServerIP),
		slog.String("referer", r.Referer),
		slog.String("latency", r.Latency),
		slog.Bool("cacheLookup", r.CacheLookup),
		slog.Bool("cacheHit", r.CacheHit),
		slog.Bool("cacheValidatedWithOriginServer", r.CacheValidatedWithOriginServer),
		slog.String("cacheFillBytes", r.CacheFillBytes),
		slog.String("protocol", r.Protocol),
	)
}

// TODO
/*
func WithHTTPRequest(r *HTTPRequest) *Logger {
	return Default().With(httpRequestKey, r)
}

func WithInsertID(id string) *Logger {
	return Default().With(insertIDKey, id)
}
*/

// TODO: Labels

type Operation struct {
	ID       string
	Producer string
	First    bool
	Last     bool
}

func (o *Operation) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", o.ID),
		slog.String("producer", o.Producer),
		slog.Bool("first", o.First),
		slog.Bool("last", o.Last),
	)
}

// TODO
/*
func WithOperation(o *Operation) *Logger {
	return Default().With(operationKey, o)
}
*/

type sourceHandler struct {
	slog.Handler
}

func (h *sourceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &sourceHandler{h.Handler.WithAttrs(attrs)}
}

func (h *sourceHandler) Handle(ctx context.Context, r slog.Record) error {
	pcs := make([]uintptr, 1)

	// skip [runtime.Callers, this function, slog.Logger.log, slog.Logger.Log, clog]
	n := runtime.Callers(5, pcs)
	if n == 0 {
		return nil
	}

	fs := runtime.CallersFrames(pcs)
	f, _ := fs.Next()
	r.AddAttrs(slog.Any(sourceLocationKey, slog.Source{
		File:     f.File,
		Line:     f.Line,
		Function: f.Function,
	}))

	/*
		for {
			f, more := fs.Next()
			fmt.Printf("%s %s:%d\n", f.Function, f.File, f.Line)

			if !strings.Contains(f.File, "go/clog") {
				r.AddAttrs(slog.Any(sourceLocationKey, slog.Source{
					File:     f.File,
					Line:     f.Line,
					Function: f.Function,
				}))
			}

			if !more {
				break
			}
		}
	*/

	return h.Handler.Handle(ctx, r)
}

type otelTraceHandler struct {
	slog.Handler
}

func (h *otelTraceHandler) Handle(ctx context.Context, r slog.Record) error {
	spanCtx := trace.SpanContextFromContext(ctx)

	if spanCtx.HasTraceID() {
		r.AddAttrs(slog.String(traceKey, spanCtx.TraceID().String()))
	}

	if spanCtx.HasSpanID() {
		r.AddAttrs(slog.String(spanIDKey, spanCtx.SpanID().String()))
	}

	return h.Handler.Handle(ctx, r)
}
