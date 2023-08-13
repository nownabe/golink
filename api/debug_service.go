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
	debug1(ctx)
	return connect.NewResponse(&golinkv1.DebugResponse{}), nil
}

func debug1(ctx context.Context) {
	debug2(ctx)
}

func debug2(ctx context.Context) {
	debug3(ctx)
}

func debug3(ctx context.Context) {
	panic("in debug")
}
