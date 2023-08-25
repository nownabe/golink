package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nownabe/golink/backend/middleware"
	"go.opentelemetry.io/otel/trace"
)

func requestIDRecordHandler(buf *string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*buf = r.Header.Get("X-Request-Id")
	})
}

func TestRequestID(t *testing.T) {
	t.Parallel()

	traceID, err := trace.TraceIDFromHex("abcdabcdabcdabcdabcdabcdabcdabcd")
	if err != nil {
		t.Fatal(err)
	}
	spanID, err := trace.SpanIDFromHex("abcdabcdabcdabcd")
	if err != nil {
		t.Fatal(err)
	}
	cfg := trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
	}
	sc := trace.NewSpanContext(cfg)

	tests := map[string]struct {
		reqHeader string
		ctx       context.Context
		want      string
	}{
		"no header": {
			reqHeader: "",
			ctx:       context.Background(),
			want:      "(random)",
		},
		"with header": {
			reqHeader: "header",
			ctx:       context.Background(),
			want:      "header",
		},
		"with trace context": {
			reqHeader: "",
			ctx:       trace.ContextWithSpanContext(context.Background(), sc),
			want:      traceID.String(),
		},
		"with trace context and header": {
			reqHeader: "header",
			ctx:       trace.ContextWithSpanContext(context.Background(), sc),
			want:      "header",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var got string

			h := middleware.NewRequestID()(requestIDRecordHandler(&got))
			s := httptest.NewServer(h)
			defer s.Close()

			req, err := http.NewRequestWithContext(tt.ctx, http.MethodGet, s.URL, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("X-Request-Id", tt.reqHeader)

			h.ServeHTTP(httptest.NewRecorder(), req)

			if tt.want == "(random)" {
				if got == "" {
					t.Fatal("Request ID should be set")
				}
			} else {
				if got != tt.want {
					t.Errorf("Request ID should be %q, but got %q", tt.want, got)
				}
			}
		})
	}
}
