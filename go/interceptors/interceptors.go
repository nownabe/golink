package interceptors

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bufbuild/connect-go"
	"go.opentelemetry.io/otel/trace"

	"github.com/nownabe/golink/go/clog"
	"github.com/nownabe/golink/go/clog/clogcontext"
	"github.com/nownabe/golink/go/errors"
	"github.com/nownabe/golink/go/golinkcontext"
)

const (
	googHeaderPrefix   = "accounts.google.com:"
	headerUserEmail    = "X-Appengine-User-Email"
	headerUserID       = "X-Appengine-User-Id"
	headerTraceContext = "X-Cloud-Trace-Context"
	headerRequestID    = "X-Request-Id"
)

func NewAuthorizer() connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			email := strings.TrimPrefix(req.Header().Get(headerUserEmail), googHeaderPrefix)
			ctx = golinkcontext.WithUserEmail(ctx, email)

			userID := strings.TrimPrefix(req.Header().Get(headerUserID), googHeaderPrefix)
			ctx = golinkcontext.WithUserID(ctx, userID)

			return next(ctx, req)
		})
	})
}

func NewDummyUser(email, userID string) connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			req.Header().Set(headerUserEmail, fmt.Sprintf("accounts.google.com:%s", email))
			req.Header().Set(headerUserID, fmt.Sprintf("accounts.google.com:%s", userID))
			return next(ctx, req)
		})
	})
}

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

func NewRequestID() connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			reqID := req.Header().Get(headerRequestID)
			if reqID == "" {
				reqID = randomString(64)
			}
			req.Header().Set(headerRequestID, reqID)
			return next(clogcontext.WithRequestID(ctx, reqID), req)
		})
	})
}

func NewLogger() connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			res, err := next(ctx, req)

			r := &clog.HTTPRequest{
				RequestMethod:                  req.HTTPMethod(),
				RequestURL:                     req.Spec().Procedure,
				RequestSize:                    req.Header().Get(headerContentLength),
				Status:                         "",
				ResponseSize:                   "",
				UserAgent:                      req.Header().Get(headerUserAgent),
				RemoteIP:                       getRemoteIP(req),
				ServerIP:                       "",
				Referer:                        req.Header().Get(headerReferer),
				Latency:                        time.Since(start),
				CacheLookup:                    false,
				CacheHit:                       false,
				CacheValidatedWithOriginServer: false,
				CacheFillBytes:                 0,
				Protocol:                       req.Peer().Protocol,
			}
			clog.InfoHTTPRequest(ctx, req.Spec().Procedure, r)

			return res, err
		})
	})
}

func WithTracer(tracer trace.Tracer) connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			ctx, span := tracer.Start(ctx, req.Spec().Procedure)
			defer span.End()
			return next(ctx, req)
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

var randomReaderPool = sync.Pool{New: func() interface{} {
	return bufio.NewReader(rand.Reader)
}}

const (
	randomStringAlphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomStringAlphabetLen = 62
	randomStringMaxByte     = 255 - (256 % randomStringAlphabetLen)
)

func randomString(length uint8) string {
	reader := randomReaderPool.Get().(*bufio.Reader)
	defer randomReaderPool.Put(reader)

	b := make([]byte, length)
	r := make([]byte, length+(length/4))
	var i uint8 = 0

	for {
		_, err := io.ReadFull(reader, r)
		if err != nil {
			err := errors.Wrap(err, "failed to read random bytes")
			clog.AlertErr(context.Background(), err)
			panic("unexpected error in randomString")
		}
		for _, rb := range r {
			if rb > randomStringMaxByte {
				continue
			}
			b[i] = randomStringAlphabet[rb%randomStringAlphabetLen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}
