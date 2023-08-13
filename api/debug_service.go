package api

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/nownabe/golink/go/clog"
	"github.com/nownabe/golink/go/errors"

	golinkv1 "github.com/nownabe/golink/api/gen/golink/v1"
)

type debugService struct{}

func (s *debugService) Debug(
	ctx context.Context,
	req *connect.Request[golinkv1.DebugRequest],
) (*connect.Response[golinkv1.DebugResponse], error) {
	if err := debug1(ctx, req); err != nil {
		clog.Err(ctx, errors.Wrap(err, "debug1 failed"))
	}
	return connect.NewResponse(&golinkv1.DebugResponse{}), nil
}

func debug1(ctx context.Context, req *connect.Request[golinkv1.DebugRequest]) error {
	if err := debug2(ctx, req); err != nil {
		return errors.Wrap(err, "debug2 failed")
	}
	return nil
}

func debug2(ctx context.Context, req *connect.Request[golinkv1.DebugRequest]) error {
	if err := debug3(ctx, req); err != nil {
		return errors.Wrap(err, "debug3 failed")
	}
	return nil
}

func debug3(ctx context.Context, req *connect.Request[golinkv1.DebugRequest]) error {
	return errors.New("debug error")
}
