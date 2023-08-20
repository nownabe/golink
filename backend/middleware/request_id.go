package middleware

import (
	"bufio"
	"context"
	"crypto/rand"
	"io"
	"net/http"
	"sync"

	"github.com/nownabe/golink/go/clog"
	"github.com/nownabe/golink/go/clog/clogcontext"
	"github.com/nownabe/golink/go/errors"
	"go.opentelemetry.io/otel/trace"
)

const headerRequestID = "X-Request-Id"

func NewRequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			reqID := r.Header.Get(headerRequestID)
			if reqID == "" {
				sc := trace.SpanContextFromContext(ctx)
				if sc.IsValid() {
					reqID = sc.TraceID().String()
				} else {
					reqID = randomString(64)
				}
				r.Header.Set(headerRequestID, reqID)
			}
			ctx = clogcontext.WithRequestID(r.Context(), reqID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
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
	reader, ok := randomReaderPool.Get().(*bufio.Reader)
	if !ok {
		err := errors.New("failed to get random reader from pool")
		clog.AlertErr(context.Background(), err)
		panic(err)
	}
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
