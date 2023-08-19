package backend

import (
	"net/http"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/nownabe/golink/api/gen/golink/v1/golinkv1connect"

	"github.com/nownabe/golink/backend/interceptor"
)

func newAPIHandler(repo *repository, debug bool, dummyUser string) http.Handler {
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

	// TODO: Move /api/healthz to /healthz
	grpcHandler.HandleFunc("/healthz", healthz)

	grpcHandler.HandleFunc("/", http.NotFound)

	return grpcHandler
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
