package api

import (
	"context"

	"github.com/bufbuild/connect-go"

	golinkv1 "github.com/nownabe/golink/api/gen/golink/v1"
)

type debugService struct{}

func (s *debugService) Debug(
	ctx context.Context,
	req *connect.Request[golinkv1.DebugRequest],
) (*connect.Response[golinkv1.DebugResponse], error) {
	debug1(ctx, req)
	return connect.NewResponse(&golinkv1.DebugResponse{}), nil
}

func debug1(ctx context.Context, req *connect.Request[golinkv1.DebugRequest]) {
	debug2(ctx, req)
}

func debug2(ctx context.Context, req *connect.Request[golinkv1.DebugRequest]) {
	debug3(ctx, req)
}

func debug3(ctx context.Context, req *connect.Request[golinkv1.DebugRequest]) {
}
