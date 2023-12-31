package backend

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	"go.nownabe.dev/clog"
	"go.nownabe.dev/clog/errors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/nownabe/golink/backend/middleware"
)

const (
	readHeaderTimeoutSeconds = 10
	// https://cloud.google.com/appengine/docs/standard/how-instances-are-managed#shutdown
	shutdownTimeoutSeconds = 3
)

type LocalDevelopmentConfig struct {
	LocalConsoleURL string
	DebugEndpoint   bool
	DummyUserEmail  string
	DummyUserID     string
}

type App interface {
	Run(ctx context.Context) error
}

// New returns a new backend app.
func New(
	port string,
	allowedOrigins []string,
	tracerName string,
	apiPrefix string,
	consolePrefix string,
	firestoreClient *firestore.Client,
	ldcfg LocalDevelopmentConfig,
) App {
	repo := &repository{firestoreClient}
	h := handler(repo, apiPrefix, consolePrefix, ldcfg.DebugEndpoint)
	for _, m := range middlewares(allowedOrigins, tracerName, consolePrefix, ldcfg) {
		h = m(h)
	}

	return &app{
		port:    port,
		handler: h,
	}
}

func handler(
	repo *repository,
	apiPrefix string,
	consolePrefix string,
	debug bool,
) http.Handler {
	rh := newRedirectHandler(repo, consolePrefix)
	ah := newAPIHandler(repo, debug)
	hh := newHealthHandler()

	mux := http.NewServeMux()
	// https://connectrpc.com/docs/go/routing#prefixing-routes
	mux.Handle(apiPrefix+"/", http.StripPrefix(apiPrefix, ah))
	mux.Handle("/health", hh)
	mux.Handle("/", rh)

	h2s := &http2.Server{}
	return h2c.NewHandler(mux, h2s)
}

func middlewares(
	allowedOrigins []string,
	tracerName string,
	consolePrefix string,
	ldcfg LocalDevelopmentConfig,
) []middleware.Middleware {
	ms := []middleware.Middleware{
		// innermost
		middleware.NewLocalConsoleRedirector(consolePrefix, ldcfg.LocalConsoleURL),
		middleware.NewAuthorizer(),
		middleware.NewCORS(allowedOrigins),
		middleware.NewHTTPLogger(),
		middleware.NewRequestID(),
		middleware.NewTraceContext(tracerName),
		middleware.NewRecoverer(),
		middleware.NewDummyUser(ldcfg.DummyUserEmail, ldcfg.DummyUserID),
		// outermost
	}

	return ms
}

type app struct {
	port    string
	handler http.Handler
}

func (a *app) Run(ctx context.Context) error {
	return a.serve(ctx)
}

func (a *app) serve(ctx context.Context) error {
	s := &http.Server{
		Addr:              ":" + a.port,
		Handler:           a.handler,
		ReadHeaderTimeout: readHeaderTimeoutSeconds * time.Second,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		ch := make(chan os.Signal, 1)

		// https://cloud.google.com/appengine/docs/standard/how-instances-are-managed#shutdown
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

		sig := <-ch
		clog.Noticef(ctx, "received signal %v and started terminating gracefully", sig)

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSeconds*time.Second)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			err := errors.Errorf("failed to shutdown gracefully: %w", err)
			clog.Err(ctx, err)
		}

		clog.Notice(ctx, "completed shutdown gracefully")
		close(idleConnsClosed)
	}()

	clog.Notice(ctx, "starting to listen and serve")

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		return errors.Errorf("failed to listen and serve: %w", err)
	}

	<-idleConnsClosed

	return nil
}
