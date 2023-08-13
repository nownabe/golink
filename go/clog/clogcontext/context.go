package clogcontext

import "context"

type (
	keyRequestID struct{}
)

func WithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, keyRequestID{}, reqID)
}
