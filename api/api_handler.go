package api

import (
	"net/http"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/nownabe/golink/api/gen/golink/v1/golinkv1connect"

	"github.com/nownabe/golink/api/interceptor"
)

type APIConfig struct {
	Prefix    string
	Debug     bool
	DummyUser string
}

func newAPIHandler(cfg *APIConfig, repo *repository) http.Handler {
	// TODO: Move interceptors to route http middlewares
	interceptors := []connect.Interceptor{
		// outermost
		interceptor.NewRecoverer(),
		interceptor.WithTracer(),
		interceptor.NewRequestID(),
		interceptor.NewAuthorizer(),
		interceptor.NewLogger(),
		// innermost
	}

	if cfg.DummyUser != "" {
		u := strings.Split(cfg.DummyUser, ":")
		interceptors = append([]connect.Interceptor{interceptor.NewDummyUser(u[0], u[1])}, interceptors...)
	}

	interceptorsOpt := connect.WithInterceptors(interceptors...)

	grpcHandler := http.NewServeMux()

	svc := &golinkService{repo}
	grpcHandler.Handle(golinkv1connect.NewGolinkServiceHandler(svc, interceptorsOpt))

	if cfg.Debug {
		grpcHandler.Handle(golinkv1connect.NewDebugServiceHandler(&debugService{}, interceptorsOpt))
	}

	// TODO: Move /api/healthz to /healthz
	grpcHandler.HandleFunc("/healthz", healthz)

	grpcHandler.HandleFunc("/", http.NotFound)

	// https://connectrpc.com/docs/go/routing#prefixing-routes
	prefixedMux := http.NewServeMux()
	prefixedMux.Handle(cfg.Prefix+"/", http.StripPrefix(cfg.Prefix, grpcHandler))

	return prefixedMux
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
