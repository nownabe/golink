package main

import (
	"os"
	"strings"

	"github.com/nownabe/golink/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	svc := api.NewGolinkService(nil)
	if err := api.New(svc, port, "/api", origins).Run(); err != nil {
		// TODO: log: fatal
		panic(err)
	}
}
