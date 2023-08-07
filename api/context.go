package api

import "context"

type (
	contextKeyUserEmail struct{}
	contextKeyUserID    struct{}
)

func WithUserEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, contextKeyUserEmail{}, email)
}

func UserEmailFrom(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(contextKeyUserEmail{}).(string)
	return v, ok
}

func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKeyUserID{}, id)
}

func UserIDFrom(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(contextKeyUserID{}).(string)
	return v, ok
}
