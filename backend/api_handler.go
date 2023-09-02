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
		// innermost
	}

	interceptorsOpt := connect.WithInterceptors(interceptors...)

	grpcHandler := http.NewServeMux()

	svc := &golinkService{repo}
	grpcHandler.Handle(golinkv1connect.NewGolinkServiceHandler(svc, interceptorsOpt))

	if debug {
		grpcHandler.Handle(golinkv1connect.NewDebugServiceHandler(&debugService{}, interceptorsOpt))
	}

	return grpcHandler
}
