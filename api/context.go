package api

import "context"

type contextKey[T any] struct{}

func withValue[T any](ctx context.Context, val T) context.Context {
	return context.WithValue(ctx, contextKey[T]{}, val)
}

func valueFrom[T any](ctx context.Context) (T, bool) {
	v, ok := ctx.Value(contextKey[T]{}).(T)
	return v, ok
}
