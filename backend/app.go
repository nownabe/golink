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
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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
	allowedOrigins []string,
	apiPrefix string,
	consolePrefix string,
	firestoreClient *firestore.Client,
	debug bool,
	dummyUser string,
) App {
	repo := &repository{firestoreClient}

	return &app{
		port:           port,
		allowedOrigins: allowedOrigins,
		apiPrefix:      apiPrefix,
		redirectHandler: &redirectHandler{
			consolePrefix: consolePrefix,
			repo:          repo,
		},
		apiHandler: newAPIHandler(repo, debug, dummyUser),
	}
}

type app struct {
	port            string
	allowedOrigins  []string
	apiPrefix       string
	redirectHandler http.Handler
	apiHandler      http.Handler
}

func (a *app) Run(ctx context.Context) error {
	return a.serve(ctx)
}

func (a *app) serve(ctx context.Context) error {
	mux := http.NewServeMux()
	// https://connectrpc.com/docs/go/routing#prefixing-routes
	mux.Handle(a.apiPrefix+"/", http.StripPrefix(a.apiPrefix, a.apiHandler))
	mux.Handle("/", a.redirectHandler)

	h2s := &http2.Server{}
	s := &http.Server{
		Addr:              ":" + a.port,
		Handler:           a.cors(h2c.NewHandler(mux, h2s)),
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

func (a *app) cors(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   a.allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})

	return c.Handler(h)
}
