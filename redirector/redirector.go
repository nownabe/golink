package redirector

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nownabe/golink/go/clog"
	"github.com/nownabe/golink/go/errors"
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

func (r *Redirector) Run(ctx context.Context) error {
	return r.serve(ctx)
}

func (r Redirector) serve(ctx context.Context) error {
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
		clog.Noticef(ctx, "received signal %v and terminating", sig)

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSeconds*time.Second)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			err := errors.Wrap(err, "failed to shutdown gracefully") // TODO:
			clog.Err(ctx, err)
		}

		clog.Notice(ctx, "completed shutdown gracefully")
		close(idleConnsClosed)
	}()

	clog.Notice(ctx, "starting to listen and serve")

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		return errors.Wrap(err, "failed to listen and serve")
	}

	<-idleConnsClosed

	return nil
}
