package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/nownabe/golink/go/clog"
	"github.com/nownabe/golink/go/errors"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/nownabe/golink/api/gen/golink/v1/golinkv1connect"
)

const (
	readHeaderTimeoutSeconds = 10
	shutdownTimeoutSeconds   = 120
)

type API interface {
	Run(ctx context.Context) error
}

func New(
	golinkSvc golinkv1connect.GolinkServiceHandler,
	port, pathPrefix string,
	allowedOrigins []string,
	interceptors []connect.Interceptor,
	debug bool,
) API {
	return &api{
		golinkSvc:      golinkSvc,
		port:           port,
		pathPrefix:     pathPrefix,
		allowedOrigins: allowedOrigins,
		interceptors:   interceptors,
		debug:          debug,
	}
}

type api struct {
	golinkSvc golinkv1connect.GolinkServiceHandler

	port           string
	pathPrefix     string
	allowedOrigins []string
	interceptors   []connect.Interceptor
	debug          bool
}

func (a *api) Run(ctx context.Context) error {
	return a.serve(ctx)
}

func (a *api) buildServer() *http.Server {
	interceptors := connect.WithInterceptors(a.interceptors...)

	grpcHandler := http.NewServeMux()
	grpcHandler.Handle(golinkv1connect.NewGolinkServiceHandler(a.golinkSvc, interceptors))

	if a.debug {
		grpcHandler.Handle(golinkv1connect.NewDebugServiceHandler(&debugService{}, interceptors))
	}

	mux := http.NewServeMux()
	// https://connectrpc.com/docs/go/routing#prefixing-routes
	mux.Handle(a.pathPrefix+"/", http.StripPrefix(a.pathPrefix, grpcHandler))
	mux.HandleFunc(a.pathPrefix+"/healthz", a.healthz)
	mux.HandleFunc("/", http.NotFound)

	h2s := &http2.Server{}
	h1s := &http.Server{
		Addr:              ":" + a.port,
		Handler:           a.cors(h2c.NewHandler(mux, h2s)),
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

func (a *api) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
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
