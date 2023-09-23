package middleware

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"time"

	"go.nownabe.dev/clog"
)

const (
	headerContentLength = "Content-Length"
	headerUserIP        = "X-Appengine-User-Ip"
	headerForwardedFor  = "X-Forwarded-For"
	headerRealIP        = "X-Real-Ip"
)

type wrappedWriter interface {
	http.ResponseWriter

	status() int
	bytes() int
	Unwrap() http.ResponseWriter
}

func wrapWriter(w http.ResponseWriter, protoMajor int) wrappedWriter {
	_, fl := w.(http.Flusher)

	base := newBaseWriterWrapper(w)

	if protoMajor == 2 {
		_, ps := w.(http.Pusher)
		if fl && ps {
			return &http2FancyWriterWrapper{base}
		}
	} else {
		_, hj := w.(http.Hijacker)
		_, rf := w.(io.ReaderFrom)
		if fl && hj && rf {
			return &httpFancyWriterWrapper{base}
		}
		if fl && hj {
			return &flushHijackWriterWrapper{base}
		}
		if hj {
			return &hijackWriterWrapper{base}
		}
	}
	if fl {
		return &flushWriterWrapper{base}
	}
	return base
}

type baseWriterWrapper struct {
	http.ResponseWriter

	writtenBytes int
	statusCode   int
	wroteHeader  bool
}

func newBaseWriterWrapper(w http.ResponseWriter) *baseWriterWrapper {
	return &baseWriterWrapper{
		ResponseWriter: w,
		writtenBytes:   0,
		statusCode:     0,
		wroteHeader:    false,
	}
}

func (w *baseWriterWrapper) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *baseWriterWrapper) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.writtenBytes += n
	return n, err
}

func (w *baseWriterWrapper) WriteHeader(statusCode int) {
	if !w.wroteHeader {
		w.statusCode = statusCode
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(statusCode)
	}
}

func (w *baseWriterWrapper) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *baseWriterWrapper) status() int {
	return w.statusCode
}

func (w *baseWriterWrapper) bytes() int {
	return w.writtenBytes
}

type flushWriterWrapper struct {
	*baseWriterWrapper
}

var _ http.Flusher = &flushWriterWrapper{}

func (w *flushWriterWrapper) Flush() {
	w.wroteHeader = true
	fl := w.ResponseWriter.(http.Flusher)
	fl.Flush()
}

type hijackWriterWrapper struct {
	*baseWriterWrapper
}

var _ http.Hijacker = &hijackWriterWrapper{}

func (w *hijackWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := w.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

type flushHijackWriterWrapper struct {
	*baseWriterWrapper
}

var (
	_ http.Flusher  = &flushHijackWriterWrapper{}
	_ http.Hijacker = &flushHijackWriterWrapper{}
)

func (w *flushHijackWriterWrapper) Flush() {
	w.wroteHeader = true
	fl := w.ResponseWriter.(http.Flusher)
	fl.Flush()
}

func (w *flushHijackWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := w.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

type httpFancyWriterWrapper struct {
	*baseWriterWrapper
}

var (
	_ http.Flusher  = &httpFancyWriterWrapper{}
	_ http.Hijacker = &httpFancyWriterWrapper{}
	_ io.ReaderFrom = &httpFancyWriterWrapper{}
)

func (w *httpFancyWriterWrapper) Flush() {
	w.wroteHeader = true
	fl := w.ResponseWriter.(http.Flusher)
	fl.Flush()
}

func (w *httpFancyWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := w.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

func (w *httpFancyWriterWrapper) ReadFrom(r io.Reader) (int64, error) {
	rf := w.ResponseWriter.(io.ReaderFrom)
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err := rf.ReadFrom(r)
	w.baseWriterWrapper.writtenBytes += int(n)
	return n, err
}

func (w *httpFancyWriterWrapper) WriteHeader(statusCode int) {
	w.baseWriterWrapper.WriteHeader(statusCode)
}

type http2FancyWriterWrapper struct {
	*baseWriterWrapper
}

var (
	_ http.Flusher = &http2FancyWriterWrapper{}
	_ http.Pusher  = &http2FancyWriterWrapper{}
)

func (w *http2FancyWriterWrapper) Flush() {
	w.wroteHeader = true
	fl := w.ResponseWriter.(http.Flusher)
	fl.Flush()
}

func (w *http2FancyWriterWrapper) Push(target string, opts *http.PushOptions) error {
	return w.baseWriterWrapper.ResponseWriter.(http.Pusher).Push(target, opts)
}

func NewHTTPLogger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := wrapWriter(w, r.ProtoMajor)

			start := time.Now()
			defer func() {
				hr := &clog.HTTPRequest{
					RequestMethod: r.Method,
					RequestURL:    r.URL.String(),
					RequestSize:   r.ContentLength,
					Status:        ww.status(),
					ResponseSize:  int64(ww.bytes()),
					UserAgent:     r.UserAgent(),
					RemoteIP:      getRemoteIP(r),
					ServerIP:      r.Host,
					Referer:       r.Referer(),
					Latency:       time.Since(start),
					Protocol:      r.Proto,
				}
				clog.HTTPReq(r.Context(), hr)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

func getRemoteIP(r *http.Request) string {
	if ip := r.Header.Get(headerForwardedFor); ip != "" {
		return ip
	}
	if ip := r.Header.Get(headerUserIP); ip != "" {
		return ip
	}
	if ip := r.Header.Get(headerRealIP); ip != "" {
		return ip
	}

	return ""
}
