package api

import (
	"context"
	"strings"

	"github.com/bufbuild/connect-go"
)

const (
	headerUserEmail    = "X-Appengine-User-Email"
	headerUserID       = "X-Appengine-User-Id"
	headerTraceContext = "X-Cloud-Trace-Context"
)

/*
"Traceparent":[]string{"00-6dc654ab15a4edc4f222de83a6b5b861-a057c88a1ebcb8e1-00"},
"X-Cloud-Trace-Context":[]string{"6dc654ab15a4edc4f222de83a6b5b861/11553923864589023457"},
*/

// TODO

func newAuthorizer() connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			email := strings.Split(req.Header().Get(headerUserEmail), ":")[1]
			ctx = withValue[UserEmail](ctx, UserEmail(email))

			userID := strings.Split(req.Header().Get(headerUserID), ":")[1]
			ctx = withValue[UserID](ctx, UserID(userID))

			return next(ctx, req)
		})
	})
}
