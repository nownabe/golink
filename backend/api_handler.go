package backend

import (
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/nownabe/golink/backend/gen/golink/v1/golinkv1connect"
	"github.com/nownabe/golink/backend/interceptor"
)

func newAPIHandler(repo *repository, debug bool) http.Handler {
	// TODO: Move interceptors to route http middlewares
	interceptors := []connect.Interceptor{
		// outermost
		interceptor.NewRecoverer(),
		interceptor.NewAuthorizer(),
		interceptor.NewLogger(),
		// innermost
	}

	interceptorsOpt := connect.WithInterceptors(interceptors...)

	grpcHandler := http.NewServeMux()

	svc := &golinkService{repo}
	grpcHandler.Handle(golinkv1connect.NewGolinkServiceHandler(svc, interceptorsOpt))

	if debug {
		grpcHandler.Handle(golinkv1connect.NewDebugServiceHandler(&debugService{}, interceptorsOpt))
	}

	// TODO: Remove this after Golink v0.0.7 is published
	grpcHandler.Handle("/healthz", newHealthHandler())

	grpcHandler.HandleFunc("/", http.NotFound)

	return grpcHandler
}
