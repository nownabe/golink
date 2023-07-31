// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: golink/v1/golink.proto

package golinkv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/nownabe/golink/api/gen/golink/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// GolinkServiceName is the fully-qualified name of the GolinkService service.
	GolinkServiceName = "golink.v1.GolinkService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// GolinkServiceCreateGolinkProcedure is the fully-qualified name of the GolinkService's
	// CreateGolink RPC.
	GolinkServiceCreateGolinkProcedure = "/golink.v1.GolinkService/CreateGolink"
)

// GolinkServiceClient is a client for the golink.v1.GolinkService service.
type GolinkServiceClient interface {
	CreateGolink(context.Context, *connect_go.Request[v1.CreateGolinkRequest]) (*connect_go.Response[v1.CreateGolinkResponse], error)
}

// NewGolinkServiceClient constructs a client for the golink.v1.GolinkService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewGolinkServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) GolinkServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &golinkServiceClient{
		createGolink: connect_go.NewClient[v1.CreateGolinkRequest, v1.CreateGolinkResponse](
			httpClient,
			baseURL+GolinkServiceCreateGolinkProcedure,
			opts...,
		),
	}
}

// golinkServiceClient implements GolinkServiceClient.
type golinkServiceClient struct {
	createGolink *connect_go.Client[v1.CreateGolinkRequest, v1.CreateGolinkResponse]
}

// CreateGolink calls golink.v1.GolinkService.CreateGolink.
func (c *golinkServiceClient) CreateGolink(ctx context.Context, req *connect_go.Request[v1.CreateGolinkRequest]) (*connect_go.Response[v1.CreateGolinkResponse], error) {
	return c.createGolink.CallUnary(ctx, req)
}

// GolinkServiceHandler is an implementation of the golink.v1.GolinkService service.
type GolinkServiceHandler interface {
	CreateGolink(context.Context, *connect_go.Request[v1.CreateGolinkRequest]) (*connect_go.Response[v1.CreateGolinkResponse], error)
}

// NewGolinkServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewGolinkServiceHandler(svc GolinkServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	golinkServiceCreateGolinkHandler := connect_go.NewUnaryHandler(
		GolinkServiceCreateGolinkProcedure,
		svc.CreateGolink,
		opts...,
	)
	return "/golink.v1.GolinkService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case GolinkServiceCreateGolinkProcedure:
			golinkServiceCreateGolinkHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedGolinkServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedGolinkServiceHandler struct{}

func (UnimplementedGolinkServiceHandler) CreateGolink(context.Context, *connect_go.Request[v1.CreateGolinkRequest]) (*connect_go.Response[v1.CreateGolinkResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("golink.v1.GolinkService.CreateGolink is not implemented"))
}
