package clog

import (
	"fmt"
	"log/slog"
	"runtime"
	"strconv"
	"time"
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
	Status                         int
	ResponseSize                   int64
	UserAgent                      string
	RemoteIP                       string
	ServerIP                       string
	Referer                        string
	Latency                        time.Duration
	CacheLookup                    bool
	CacheHit                       bool
	CacheValidatedWithOriginServer bool
	CacheFillBytes                 int64
	Protocol                       string
}

func (r *HTTPRequest) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("requestMethod", r.RequestMethod),
		slog.String("requestUrl", r.RequestURL),
		slog.String("requestSize", r.RequestSize),
		slog.Int("status", r.Status),
		slog.String("responseSize", strconv.FormatInt(r.ResponseSize, 10)),
		slog.String("userAgent", r.UserAgent),
		slog.String("remoteIp", r.RemoteIP),
		slog.String("serverIp", r.ServerIP),
		slog.String("referer", r.Referer),
		slog.String("latency", fmt.Sprintf("%.9fs", r.Latency.Seconds())),
		slog.Bool("cacheLookup", r.CacheLookup),
		slog.Bool("cacheHit", r.CacheHit),
		slog.Bool("cacheValidatedWithOriginServer", r.CacheValidatedWithOriginServer),
		slog.String("cacheFillBytes", strconv.FormatInt(r.CacheFillBytes, 10)),
		slog.String("protocol", r.Protocol),
	)
}

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

func getSource(skip int) *slog.Source {
	pcs := make([]uintptr, 1)

	n := runtime.Callers(skip, pcs)
	if n == 0 {
		return nil
	}

	fs := runtime.CallersFrames(pcs)
	f, _ := fs.Next()

	return &slog.Source{
		File:     f.File,
		Line:     f.Line,
		Function: f.Function,
	}
}
