package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/bufbuild/connect-go"
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

// TODO ?

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

// TODO: request log
// TODO: request id
// TODO: recover
// TODO: trace
