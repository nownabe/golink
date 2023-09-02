package middleware

import (
	"net/http"

	"github.com/nownabe/golink/backend/clog"
	"github.com/nownabe/golink/backend/errors"
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
					err = errors.Wrap(err, "recovering panic")
					clog.AlertErr(ctx, err)

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
			panicked = false
		})
	}
}
