package main

import (
	"log"

	"github.com/nownabe/golink/redirector"
)

func main() {
	// TODO: Configure logger
	repo, err := redirector.NewRepository()
	if err != nil {
		// TODO: log
		log.Fatal(err)
	}

	handler := redirector.NewHandler(repo)

	if err := redirector.New("8080", handler).Run(); err != nil {
		// TODO: log
		log.Fatal(err)
	}
}
