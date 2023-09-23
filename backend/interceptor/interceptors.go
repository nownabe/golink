package interceptor

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"
	"go.nownabe.dev/clog"
	"go.nownabe.dev/clog/errors"
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
					err = errors.Errorf("recovering panic: %w")
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
