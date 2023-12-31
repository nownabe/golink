package backend

import (
	"context"

	"github.com/bufbuild/connect-go"
	"go.nownabe.dev/clog"
	"go.nownabe.dev/clog/errors"

	golinkv1 "github.com/nownabe/golink/backend/gen/golink/v1"
)

type debugService struct{}

func (s *debugService) Debug(
	ctx context.Context,
	req *connect.Request[golinkv1.DebugRequest],
) (*connect.Response[golinkv1.DebugResponse], error) {
	if err := debug1(ctx, req); err != nil {
		clog.Err(ctx, errors.Errorf("debug1 failed: %w", err))
	}

	return connect.NewResponse(&golinkv1.DebugResponse{}), nil
}

func debug1(ctx context.Context, req *connect.Request[golinkv1.DebugRequest]) error {
	if err := debug2(ctx, req); err != nil {
		return errors.Errorf("debug2 failed: %w", err)
	}

	return nil
}

func debug2(ctx context.Context, req *connect.Request[golinkv1.DebugRequest]) error {
	if err := debug3(ctx, req); err != nil {
		return errors.Errorf("debug3 failed: %w", err)
	}

	return nil
}

func debug3(_ context.Context, req *connect.Request[golinkv1.DebugRequest]) error {
	return errors.New("debug error")
}
