package api

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/bufbuild/connect-go"
	"github.com/nownabe/golink/go/clog"
	"github.com/nownabe/golink/go/errors"

	golinkv1 "github.com/nownabe/golink/api/gen/golink/v1"
	"github.com/nownabe/golink/api/gen/golink/v1/golinkv1connect"
)

func NewGolinkService(repo Repository) golinkv1connect.GolinkServiceHandler {
	return &golinkService{
		repo: repo,
	}
}

type golinkService struct {
	repo Repository
}

func (s *golinkService) CreateGolink(
	ctx context.Context,
	req *connect.Request[golinkv1.CreateGolinkRequest],
) (*connect.Response[golinkv1.CreateGolinkResponse], error) {
	email, ok := UserEmailFrom(ctx)
	if !ok {
		err := errors.New("user email not found in context")
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	o := &dto{
		Name:   req.Msg.Name,
		URL:    req.Msg.Url,
		Owners: []string{email},
	}

	var exists bool

	err := s.repo.Transaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		e, err := s.repo.Exists(ctx, tx, req.Msg.Name)
		if err != nil {
			return errors.Wrapf(err, "s.repo.Exists: name=%s", req.Msg.Name)
		}

		if e {
			exists = true
			return nil
		}

		if err := s.repo.Create(ctx, tx, o); err != nil {
			return errors.Wrapf(err, "failed to create Golink(name=%s)", req.Msg.Name)
		}

		return nil
	})
	if err != nil {
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	if exists {
		return nil, errf(connect.CodeAlreadyExists, "go/%s already exists", req.Msg.Name)
	}

	res := connect.NewResponse(&golinkv1.CreateGolinkResponse{
		Golink: &golinkv1.Golink{
			Name:   req.Msg.Name,
			Url:    req.Msg.Url,
			Owners: []string{email},
		},
	})

	return res, nil
}

func (s *golinkService) GetGolink(
	ctx context.Context,
	req *connect.Request[golinkv1.GetGolinkRequest],
) (*connect.Response[golinkv1.GetGolinkResponse], error) {
	res := connect.NewResponse(&golinkv1.GetGolinkResponse{
		Golink: &golinkv1.Golink{
			Name:   req.Msg.Name,
			Url:    "https://example.com/",
			Owners: []string{"myself@example.com", "other1@example.com"},
		},
	})
	return res, nil
}

func (s *golinkService) ListGolinks(
	ctx context.Context,
	req *connect.Request[golinkv1.ListGolinksRequest],
) (*connect.Response[golinkv1.ListGolinksResponse], error) {
	res := connect.NewResponse(&golinkv1.ListGolinksResponse{
		Golinks: []*golinkv1.Golink{
			{
				Name:   "example1",
				Url:    "https://link1.example.com/",
				Owners: []string{"user@example.com"},
			},
			{
				Name:   "example2",
				Url:    "https://link2.example.com/",
				Owners: []string{"user@example.com"},
			},
		},
	})
	return res, nil
}

func (s *golinkService) ListGolinksByUrl(
	ctx context.Context,
	req *connect.Request[golinkv1.ListGolinksByUrlRequest],
) (*connect.Response[golinkv1.ListGolinksByUrlResponse], error) {
	res := connect.NewResponse(&golinkv1.ListGolinksByUrlResponse{
		Golinks: []*golinkv1.Golink{
			{
				Name:   "example1",
				Url:    "https://link1.example.com/",
				Owners: []string{"user@example.com"},
			},
			{
				Name:   "example2",
				Url:    "https://link2.example.com/",
				Owners: []string{"user@example.com"},
			},
		},
	})
	return res, nil
}

func (s *golinkService) UpdateGolink(
	ctx context.Context,
	req *connect.Request[golinkv1.UpdateGolinkRequest],
) (*connect.Response[golinkv1.UpdateGolinkResponse], error) {
	time.Sleep(1 * time.Second)
	res := connect.NewResponse(&golinkv1.UpdateGolinkResponse{
		Golink: &golinkv1.Golink{
			Name: req.Msg.Name,
			Url:  req.Msg.Url,
		},
	})
	return res, nil
}

func (s *golinkService) DeleteGolink(
	ctx context.Context,
	req *connect.Request[golinkv1.DeleteGolinkRequest],
) (*connect.Response[golinkv1.DeleteGolinkResponse], error) {
	res := connect.NewResponse(&golinkv1.DeleteGolinkResponse{})
	return res, nil
}

func (s *golinkService) AddOwner(
	ctx context.Context,
	req *connect.Request[golinkv1.AddOwnerRequest],
) (*connect.Response[golinkv1.AddOwnerResponse], error) {
	time.Sleep(1 * time.Second)
	res := connect.NewResponse(&golinkv1.AddOwnerResponse{
		Golink: &golinkv1.Golink{
			Name:   req.Msg.Name,
			Url:    "https://link1.example.com/",
			Owners: []string{"user@example.com"},
		},
	})
	return res, nil
}

func (s *golinkService) RemoveOwner(
	ctx context.Context,
	req *connect.Request[golinkv1.RemoveOwnerRequest],
) (*connect.Response[golinkv1.RemoveOwnerResponse], error) {
	time.Sleep(1 * time.Second)
	res := connect.NewResponse(&golinkv1.RemoveOwnerResponse{
		Golink: &golinkv1.Golink{
			Name:   req.Msg.Name,
			Url:    "https://link1.example.com/",
			Owners: []string{"user@example.com"},
		},
	})
	return res, nil
}

func errf(code connect.Code, format string, args ...any) error {
	return connect.NewError(code, fmt.Errorf(format, args...))
}
