package backend

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/nownabe/golink/go/clog"
	"github.com/nownabe/golink/go/errors"
)

const (
	readHeaderTimeoutSeconds = 10
	shutdownTimeoutSeconds   = 120
)

type App interface {
	Run(ctx context.Context) error
}

// New returns a new backend app.
func New(
	port string,
	consolePrefix string,
	firestoreClient *firestore.Client,
) App {
	repo := &repository{firestoreClient}

	return &app{
		port: port,
		redirectHandler: &redirectHandler{
			consolePrefix: consolePrefix,
			repo:          repo,
		},
	}
}

type app struct {
	port            string
	redirectHandler http.Handler
}

func (a *app) Run(ctx context.Context) error {
	return a.serve(ctx)
}

func (a *app) serve(ctx context.Context) error {
	s := &http.Server{
		Addr:              ":" + a.port,
		Handler:           a.redirectHandler,
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
