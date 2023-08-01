package api

import (
	"context"

	"github.com/bufbuild/connect-go"

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
	res := connect.NewResponse(&golinkv1.CreateGolinkResponse{
		Golink: &golinkv1.Golink{
			Name: req.Msg.Name,
			Url:  req.Msg.Url,
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
			Name: req.Msg.Name,
			Url:  "https://example.com/",
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
				Name: "example1",
				Url:  "https://link1.example.com/",
			},
			{
				Name: "example2",
				Url:  "https://link2.example.com/",
			},
		},
	})
	return res, nil
}

func (s *golinkService) ListGolinksByURL(
	ctx context.Context,
	req *connect.Request[golinkv1.ListGolinksByURLRequest],
) (*connect.Response[golinkv1.ListGolinksByURLResponse], error) {
	res := connect.NewResponse(&golinkv1.ListGolinksByURLResponse{
		Golinks: []*golinkv1.Golink{
			{
				Name: "example1",
				Url:  "https://link1.example.com/",
			},
			{
				Name: "example2",
				Url:  "https://link2.example.com/",
			},
		},
	})
	return res, nil
}

func (s *golinkService) UpdateGolink(
	ctx context.Context,
	req *connect.Request[golinkv1.UpdateGolinkRequest],
) (*connect.Response[golinkv1.UpdateGolinkResponse], error) {
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
