package backend

import (
	"net/http"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/nownabe/golink/backend/gen/golink/v1/golinkv1connect"
	"github.com/nownabe/golink/backend/interceptor"
)

func newAPIHandler(repo *repository, debug bool, dummyUser string) http.Handler {
	// TODO: Move interceptors to route http middlewares
	interceptors := []connect.Interceptor{
		// outermost
		interceptor.NewRecoverer(),
		interceptor.NewAuthorizer(),
		interceptor.NewLogger(),
		// innermost
	}

	if dummyUser != "" {
		u := strings.Split(dummyUser, ":")
		interceptors = append([]connect.Interceptor{interceptor.NewDummyUser(u[0], u[1])}, interceptors...)
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
