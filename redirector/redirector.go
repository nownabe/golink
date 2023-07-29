package redirector

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	readHeaderTimeoutSeconds = 10
	shutdownTimeoutSeconds   = 120
)

type Redirector struct {
	port    string
	handler http.Handler
}

func New(port string, h http.Handler) *Redirector {
	return &Redirector{
		port:    port,
		handler: h,
	}
}

func (r *Redirector) Run() error {
	return r.serve()
}

func (r Redirector) serve() error {
	s := &http.Server{
		Addr:              ":" + r.port,
		Handler:           r.handler,
		ReadHeaderTimeout: readHeaderTimeoutSeconds * time.Second,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

		sig := <-ch
		log.Print(sig)

		// TODO: log: info: received signal and terminating

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSeconds*time.Second)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			// TODO: log: err: failed to shutdown gracefully
			log.Print(err)
		}

		// TODO: log: info: completed shutdown gracefully

		close(idleConnsClosed)
	}()

	// TODO: log: info: started

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		// TODO: log: fatal: failed to listen and serve
		return err
	}

	<-idleConnsClosed

	// TODO: log: info: bye
	return nil
}
