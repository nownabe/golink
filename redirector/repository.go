package redirector

import (
	"context"
	"net/url"
)

type Repository interface {
	GetURLAndUpdateStats(ctx context.Context, name string) (*url.URL, error)
}

func NewRepository() (Repository, error) {
	return &repository{}, nil
}

type repository struct{}

func (r *repository) GetURLAndUpdateStats(ctx context.Context, name string) (*url.URL, error) {
	// TODO: wrap error
	return url.Parse("https://example.com/redirected/" + name)
}
