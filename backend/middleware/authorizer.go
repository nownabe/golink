package middleware

import (
	"net/http"
	"strings"

	"github.com/nownabe/golink/backend/golinkcontext"
)

const (
	googAccountsHeaderPrefix = "accounts.google.com:"
	headerUserEmail          = "X-Appengine-User-Email"
	headerUserID             = "X-Appengine-User-Id"
)

func NewAuthorizer() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			email := strings.TrimPrefix(r.Header.Get(headerUserEmail), googAccountsHeaderPrefix)
			userID := strings.TrimPrefix(r.Header.Get(headerUserID), googAccountsHeaderPrefix)

			if email == "" || userID == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = golinkcontext.WithUserEmail(ctx, email)
			ctx = golinkcontext.WithUserID(ctx, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
