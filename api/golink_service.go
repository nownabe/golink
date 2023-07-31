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

type Repository interface{}
