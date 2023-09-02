package clogcontext

import (
	"context"
	"sync"
)

type (
	keyRequestID struct{}
	keyLabels    struct{}
)

func WithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, keyRequestID{}, reqID)
}

func WithLabel(ctx context.Context, key, value string) context.Context {
	m, ok := ctx.Value(keyLabels{}).(*sync.Map)
	if !ok {
		m = &sync.Map{}
	}
	m.Store(key, value)
	return context.WithValue(ctx, keyLabels{}, m)
}

func RequestIDFrom(ctx context.Context) (string, bool) {
	reqID, ok := ctx.Value(keyRequestID{}).(string)
	return reqID, ok
}

func LabelFrom(ctx context.Context) (*sync.Map, bool) {
	m, ok := ctx.Value(keyLabels{}).(*sync.Map)
	return m, ok
}
