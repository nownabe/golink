package middleware

import (
	"net/http"

	"go.nownabe.dev/clog"
	"go.nownabe.dev/clog/errors"
)

func NewRecoverer() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			panicked := true
			defer func() {
				if panicked {
					rcv := recover()
					if rcv == http.ErrAbortHandler {
						panic(rcv)
					}
					err, ok := rcv.(error)
					if !ok {
						err = errors.Errorf("%v", rcv)
					}
					err = errors.Errorf("recovering panic: %w", err)
					clog.AlertErr(ctx, err)

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
			panicked = false
		})
	}
}
