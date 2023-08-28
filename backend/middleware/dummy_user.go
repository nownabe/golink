package middleware

import "net/http"

func NewDummyUser(email, userID string) Middleware {
	return func(next http.Handler) http.Handler {
		if email == "" || userID == "" {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set(headerUserEmail, googAccountsHeaderPrefix+email)
			r.Header.Set(headerUserID, googAccountsHeaderPrefix+userID)
			next.ServeHTTP(w, r)
		})
	}
}
