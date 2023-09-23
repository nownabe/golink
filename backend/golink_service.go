package backend

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/bufbuild/connect-go"
	"go.nownabe.dev/clog"
	"go.nownabe.dev/clog/errors"

	golinkv1 "github.com/nownabe/golink/backend/gen/golink/v1"
	"github.com/nownabe/golink/backend/golinkcontext"
)

type golinkService struct {
	repo *repository
}

func (s *golinkService) CreateGolink(
	ctx context.Context,
	req *connect.Request[golinkv1.CreateGolinkRequest],
) (*connect.Response[golinkv1.CreateGolinkResponse], error) {
	email, ok := golinkcontext.UserEmailFrom(ctx)
	if !ok {
		err := errors.New("user email not found in context")
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	if !isValidName(req.Msg.Name) {
		return nil, errf(connect.CodeInvalidArgument, "invalid name")
	}

	if !isValidURL(req.Msg.Url) {
		return nil, errf(connect.CodeInvalidArgument, "invalid url")
	}

	o := &dto{
		Name:   req.Msg.Name,
		URL:    req.Msg.Url,
		Owners: []string{email},
	}

	err := s.repo.Transaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		e, err := s.repo.TxExists(ctx, tx, req.Msg.Name)
		if err != nil {
			return errors.Errorf("s.repo.Exists: name=%s: %w", req.Msg.Name, err)
		}

		if e {
			return errf(connect.CodeAlreadyExists, "go/%s already exists", req.Msg.Name)
		}

		if err := s.repo.TxCreate(ctx, tx, o); err != nil {
			return errors.Errorf("failed to create Golink(name=%s): %w", req.Msg.Name, err)
		}

		return nil
	})

	if connect.CodeOf(err) != connect.CodeUnknown {
		return nil, err
	}
	if err != nil {
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	res := connect.NewResponse(&golinkv1.CreateGolinkResponse{Golink: o.ToProto()})

	return res, nil
}

func (s *golinkService) GetGolink(
	ctx context.Context,
	req *connect.Request[golinkv1.GetGolinkRequest],
) (*connect.Response[golinkv1.GetGolinkResponse], error) {
	o, err := s.repo.Get(ctx, req.Msg.Name)
	if err != nil {
		if errors.Is(err, errDocumentNotFound) {
			return nil, errf(connect.CodeNotFound, "go/%s not found", req.Msg.Name)
		}
		err := errors.Errorf("failed to get Golink(name=%s): %w", req.Msg.Name, err)
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	res := connect.NewResponse(&golinkv1.GetGolinkResponse{Golink: o.ToProto()})
	return res, nil
}

func (s *golinkService) ListGolinks(
	ctx context.Context,
	_ *connect.Request[golinkv1.ListGolinksRequest],
) (*connect.Response[golinkv1.ListGolinksResponse], error) {
	email, ok := golinkcontext.UserEmailFrom(ctx)
	if !ok {
		err := errors.New("user email not found in context")
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	dtos, err := s.repo.ListByOwner(ctx, email)
	if err != nil {
		err := errors.Errorf("failed to list Golinks: %w", err)
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	golinks := []*golinkv1.Golink{}
	for _, dto := range dtos {
		golinks = append(golinks, dto.ToProto())
	}

	res := connect.NewResponse(&golinkv1.ListGolinksResponse{Golinks: golinks})

	return res, nil
}

func (s *golinkService) ListGolinksByUrl(
	ctx context.Context,
	req *connect.Request[golinkv1.ListGolinksByUrlRequest],
) (*connect.Response[golinkv1.ListGolinksByUrlResponse], error) {
	dtos, err := s.repo.ListByURL(ctx, req.Msg.Url)
	if err != nil {
		err := errors.Errorf("s.repo.ListByURL(ctx, %q): %w", req.Msg.Url, err)
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	golinks := []*golinkv1.Golink{}
	for _, dto := range dtos {
		golinks = append(golinks, dto.ToProto())
	}

	res := connect.NewResponse(&golinkv1.ListGolinksByUrlResponse{Golinks: golinks})

	return res, nil
}

func (s *golinkService) ListPopularGolinks(
	ctx context.Context,
	req *connect.Request[golinkv1.ListPopularGolinksRequest],
) (*connect.Response[golinkv1.ListPopularGolinksResponse], error) {
	if req.Msg.Limit <= 0 {
		return nil, errf(connect.CodeInvalidArgument, "limit must be greater than 0")
	}
	if req.Msg.Limit > 100 {
		return nil, errf(connect.CodeInvalidArgument, "limit must be less than or equal to 100")
	}
	if req.Msg.Days != 7 && req.Msg.Days != 28 {
		return nil, errf(connect.CodeInvalidArgument, "days must be 7 or 28")
	}

	golinks, err := s.repo.ListPopularGolinks(ctx, int(req.Msg.Days), int(req.Msg.Limit))
	if err != nil {
		err := errors.Errorf("s.repo.ListPopularGolinks(ctx, %d, %d): %w", req.Msg.Days, req.Msg.Limit, err)
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	resp := &golinkv1.ListPopularGolinksResponse{}

	for _, golink := range golinks {
		resp.Golinks = append(resp.Golinks, golink.ToProto())
	}

	return connect.NewResponse(resp), nil
}

func (s *golinkService) UpdateGolink(
	ctx context.Context,
	req *connect.Request[golinkv1.UpdateGolinkRequest],
) (*connect.Response[golinkv1.UpdateGolinkResponse], error) {
	email, ok := golinkcontext.UserEmailFrom(ctx)
	if !ok {
		err := errors.New("user email not found in context")
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	var o *dto

	err := s.repo.Transaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var err error
		o, err = s.repo.TxGet(ctx, tx, req.Msg.Name)
		if err != nil {
			if errors.Is(err, errDocumentNotFound) {
				return errf(connect.CodeNotFound, "go/%s not found", req.Msg.Name)
			}
			return errors.Errorf("failed to get Golink(name=%s): %w", req.Msg.Name, err)
		}

		if !slices.Contains(o.Owners, email) {
			return errf(connect.CodePermissionDenied, "permission denied")
		}

		if !isValidURL(req.Msg.Url) {
			return errf(connect.CodeInvalidArgument, "invalid url")
		}

		o.URL = req.Msg.Url

		if err := s.repo.TxUpdate(ctx, tx, o); err != nil {
			return errors.Errorf("failed to update Golink(name=%s): %w", req.Msg.Name, err)
		}

		return nil
	})

	if connect.CodeOf(err) != connect.CodeUnknown {
		return nil, err
	}
	if err != nil {
		err := errors.Errorf("update transaction failed: Golink(name=%s): %w", req.Msg.Name, err)
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	res := connect.NewResponse(&golinkv1.UpdateGolinkResponse{Golink: o.ToProto()})

	return res, nil
}

func (s *golinkService) DeleteGolink(
	ctx context.Context,
	req *connect.Request[golinkv1.DeleteGolinkRequest],
) (*connect.Response[golinkv1.DeleteGolinkResponse], error) {
	email, ok := golinkcontext.UserEmailFrom(ctx)
	if !ok {
		err := errors.New("user email not found in context")
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	err := s.repo.Transaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		o, err := s.repo.TxGet(ctx, tx, req.Msg.Name)
		if err != nil {
			if errors.Is(err, errDocumentNotFound) {
				return errf(connect.CodeNotFound, "go/%s not found", req.Msg.Name)
			}
			return errors.Errorf("failed to get Golink(name=%s): %w", req.Msg.Name, err)
		}

		if !slices.Contains(o.Owners, email) {
			return errf(connect.CodePermissionDenied, "permission denied")
		}

		if err := s.repo.TxDelete(ctx, tx, req.Msg.Name); err != nil {
			return errors.Errorf("failed to delete Golink(name=%s): %w", req.Msg.Name, err)
		}

		return nil
	})

	if connect.CodeOf(err) != connect.CodeUnknown {
		return nil, err
	}
	if err != nil {
		err := errors.Errorf("delete transaction failed: Golink(name=%s): %w", req.Msg.Name, err)
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	res := connect.NewResponse(&golinkv1.DeleteGolinkResponse{})
	return res, nil
}

func (s *golinkService) AddOwner(
	ctx context.Context,
	req *connect.Request[golinkv1.AddOwnerRequest],
) (*connect.Response[golinkv1.AddOwnerResponse], error) {
	email, ok := golinkcontext.UserEmailFrom(ctx)
	if !ok {
		err := errors.New("user email not found in context")
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	var o *dto

	err := s.repo.Transaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var err error
		o, err = s.repo.TxGet(ctx, tx, req.Msg.Name)
		if err != nil {
			if errors.Is(err, errDocumentNotFound) {
				return errf(connect.CodeNotFound, "go/%s not found", req.Msg.Name)
			}
			return errors.Errorf("failed to get Golink(name=%s): %w", req.Msg.Name, err)
		}

		if !slices.Contains(o.Owners, email) {
			return errf(connect.CodePermissionDenied, "permission denied")
		}

		if slices.Contains(o.Owners, req.Msg.Owner) {
			return errf(connect.CodeInvalidArgument, "owner already exists")
		}

		if err := s.repo.TxAddOwner(ctx, tx, req.Msg.Name, req.Msg.Owner); err != nil {
			return errors.Errorf("failed to add owner: Golink(name=%s), owner=%s: %w", req.Msg.Name, req.Msg.Owner, err)
		}

		return nil
	})

	if connect.CodeOf(err) != connect.CodeUnknown {
		return nil, err
	}
	if err != nil {
		err := errors.Errorf("add owner transaction failed: Golink(name=%s): %w", req.Msg.Name, err)
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	o.Owners = append(o.Owners, req.Msg.Owner)
	res := connect.NewResponse(&golinkv1.AddOwnerResponse{Golink: o.ToProto()})

	return res, nil
}

func (s *golinkService) RemoveOwner(
	ctx context.Context,
	req *connect.Request[golinkv1.RemoveOwnerRequest],
) (*connect.Response[golinkv1.RemoveOwnerResponse], error) {
	email, ok := golinkcontext.UserEmailFrom(ctx)
	if !ok {
		err := errors.New("user email not found in context")
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	var o *dto

	err := s.repo.Transaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var err error
		o, err = s.repo.TxGet(ctx, tx, req.Msg.Name)
		if err != nil {
			if errors.Is(err, errDocumentNotFound) {
				return errf(connect.CodeNotFound, "go/%s not found", req.Msg.Name)
			}
			return errors.Errorf("failed to get Golink(name=%s): %w", req.Msg.Name, err)
		}

		if !slices.Contains(o.Owners, email) {
			return errf(connect.CodePermissionDenied, "permission denied")
		}

		if !slices.Contains(o.Owners, req.Msg.Owner) {
			return errf(connect.CodeInvalidArgument, "owner not found")
		}

		if len(o.Owners) == 1 {
			return errf(connect.CodeInvalidArgument, "cannot remove last owner")
		}

		if err := s.repo.TxRemoveOwner(ctx, tx, req.Msg.Name, req.Msg.Owner); err != nil {
			return errors.Errorf("failed to remove owner: Golink(name=%s), owner=%s: %w", req.Msg.Name, req.Msg.Owner, err)
		}

		return nil
	})

	if connect.CodeOf(err) != connect.CodeUnknown {
		return nil, err
	}
	if err != nil {
		err := errors.Errorf("remove owner transaction failed: Golink(name=%s): %w", req.Msg.Name, err)
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	o.Owners = slices.DeleteFunc(o.Owners, func(owner string) bool { return owner == req.Msg.Owner })
	res := connect.NewResponse(&golinkv1.RemoveOwnerResponse{Golink: o.ToProto()})

	return res, nil
}

func (s *golinkService) GetMe(
	ctx context.Context,
	_ *connect.Request[golinkv1.GetMeRequest],
) (*connect.Response[golinkv1.GetMeResponse], error) {
	email, ok := golinkcontext.UserEmailFrom(ctx)
	if !ok {
		err := errors.New("user email not found in context")
		clog.Err(ctx, err)
		return nil, errf(connect.CodeInternal, "internal error")
	}

	res := connect.NewResponse(&golinkv1.GetMeResponse{Email: email})
	return res, nil
}

func errf(code connect.Code, format string, args ...any) error {
	return connect.NewError(code, fmt.Errorf(format, args...))
}

func isValidName(name string) bool {
	// Firestore limitations
	// https://firebase.google.com/docs/firestore/quotas#collections_documents_and_fields
	if name == "" {
		return false
	}
	if len(name) > 1500 {
		return false
	}
	if strings.Contains(name, "/") {
		return false
	}
	if strings.HasPrefix(name, "__") || strings.HasSuffix(name, "__") {
		return false
	}
	if name == "." || name == ".." {
		return false
	}

	// Golink limitations
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return false
	}
	if name == "_" || name == "api" || name == "c" {
		return false
	}

	return true
}

func isValidURL(u string) bool {
	url, err := url.Parse(u)
	if err != nil {
		return false
	}

	return url.Scheme == "http" || url.Scheme == "https"
}
