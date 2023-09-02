package interceptor

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/nownabe/golink/backend/clog"
	"github.com/nownabe/golink/backend/errors"
)

const (
	googHeaderPrefix   = "accounts.google.com:"
	headerUserEmail    = "X-Appengine-User-Email"
	headerUserID       = "X-Appengine-User-Id"
	headerTraceContext = "X-Cloud-Trace-Context"
)

// https://github.com/golang/go/issues/25448
func NewRecoverer() connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (_ connect.AnyResponse, retErr error) {
			panicked := true
			defer func() {
				if panicked {
					r := recover()
					if r == http.ErrAbortHandler {
						panic(r)
					}
					err, ok := r.(error)
					if !ok {
						err = errors.Errorf("%v", r)
					}
					err = errors.Wrap(err, "recovering panic")
					clog.AlertErr(ctx, err)

					retErr = connect.NewError(http.StatusInternalServerError, errors.NewWithoutStack("internal error"))
				}
			}()
			res, err := next(ctx, req)
			panicked = false
			return res, err
		})
	})
}

const (
	headerContentLength = "Content-Length"
	headerUserAgent     = "User-Agent"
	headerUserIP        = "X-Appengine-User-Ip"
	headerForwardedFor  = "X-Forwarded-For"
	headerRealIP        = "X-Real-Ip"
	headerReferer       = "Referer"
)

func getRemoteIP(req connect.AnyRequest) string {
	if ip := req.Header().Get(headerForwardedFor); ip != "" {
		return ip
	}
	if ip := req.Header().Get(headerUserIP); ip != "" {
		return ip
	}
	if ip := req.Header().Get(headerRealIP); ip != "" {
		return ip
	}

	return ""
}
