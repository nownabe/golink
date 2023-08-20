package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func NewCORS(allowedOrigins []string) Middleware {
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})

	return func(next http.Handler) http.Handler {
		return c.Handler(next)
	}
}
