package middleware

import (
	"net/http"
	"strings"
)

func NewLocalConsoleRedirector(consolePrefix, localConsoleURL string) Middleware {
	return func(next http.Handler) http.Handler {
		if localConsoleURL == "" {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, consolePrefix) {
				http.Redirect(w, r, localConsoleURL+r.URL.Path, http.StatusTemporaryRedirect)
			}
			next.ServeHTTP(w, r)
		})
	}
}
