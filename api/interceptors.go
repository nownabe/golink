package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/nownabe/golink/go/clog"
	"github.com/nownabe/golink/go/errors"
)

const (
	googHeaderPrefix   = "accounts.google.com:"
	headerUserEmail    = "X-Appengine-User-Email"
	headerUserID       = "X-Appengine-User-Id"
	headerTraceContext = "X-Cloud-Trace-Context"
)

/*
"Traceparent":[]string{"00-6dc654ab15a4edc4f222de83a6b5b861-a057c88a1ebcb8e1-00"},
"X-Cloud-Trace-Context":[]string{"6dc654ab15a4edc4f222de83a6b5b861/11553923864589023457"},
*/

func NewAuthorizer() connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			email := strings.TrimPrefix(req.Header().Get(headerUserEmail), googHeaderPrefix)
			ctx = WithUserEmail(ctx, email)

			userID := strings.TrimPrefix(req.Header().Get(headerUserID), googHeaderPrefix)
			ctx = WithUserID(ctx, userID)

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

func NewRecoverer() connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (_ connect.AnyResponse, retErr error) {
			defer func() {
				if r := recover(); r != nil {
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
			return res, err
		})
	})
}

// TODO: request log
// TODO: request id
// TODO: trace
