package middleware

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func NewTraceContext(tracerName string) Middleware {
	propagator := otel.GetTextMapPropagator()
	tracer := otel.GetTracerProvider().Tracer(tracerName)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			ctx, span := tracer.Start(ctx, fmt.Sprintf("%s %s %s", r.Method, r.URL.Path, r.Proto))
			defer span.End()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
