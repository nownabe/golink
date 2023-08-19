package api

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
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	readHeaderTimeoutSeconds = 10
	shutdownTimeoutSeconds   = 120
)

type API interface {
	Run(ctx context.Context) error
}

func New(
	port string,
	allowedOrigins []string,
	cfg *APIConfig,
	fs *firestore.Client,
) API {
	repo := &repository{fs}
	return &api{
		port:           port,
		allowedOrigins: allowedOrigins,
		cfg:            cfg,
		repo:           repo,
	}
}

type api struct {
	port           string
	allowedOrigins []string
	cfg            *APIConfig
	repo           *repository
}

func (a *api) Run(ctx context.Context) error {
	return a.serve(ctx)
}

func (a *api) buildServer() *http.Server {
	apiHandler := newAPIHandler(a.cfg, a.repo)

	h2s := &http2.Server{}
	h1s := &http.Server{
		Addr:              ":" + a.port,
		Handler:           a.cors(h2c.NewHandler(apiHandler, h2s)),
		ReadHeaderTimeout: readHeaderTimeoutSeconds * time.Second,
	}

	return h1s
}

func (a *api) serve(ctx context.Context) error {
	s := a.buildServer()

	idleConnsClosed := make(chan struct{})

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

		sig := <-ch
		clog.Noticef(ctx, "received signal %s and terminating", sig)

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSeconds*time.Second)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			err := errors.Wrap(err, "failed to shutdown gracefully")
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

func (a *api) cors(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   a.allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})

	return c.Handler(h)
}
