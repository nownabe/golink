package api_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/bufbuild/connect-go"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/nownabe/golink/go/golinkcontext"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nownabe/golink/api"
	golinkv1 "github.com/nownabe/golink/api/gen/golink/v1"
	"github.com/nownabe/golink/api/gen/golink/v1/golinkv1connect"
)

type dto = api.DTO

var fsClient *firestore.Client

var cmpOptions = []cmp.Option{
	cmpopts.IgnoreTypes(time.Time{}),
	cmpopts.IgnoreUnexported(golinkv1.Golink{}),
	cmpopts.IgnoreUnexported(golinkv1.CreateGolinkResponse{}),
	cmpopts.IgnoreUnexported(golinkv1.GetGolinkResponse{}),
	cmpopts.IgnoreUnexported(golinkv1.ListGolinksResponse{}),
	cmpopts.IgnoreUnexported(golinkv1.ListGolinksByUrlResponse{}),
	cmpopts.IgnoreUnexported(golinkv1.UpdateGolinkResponse{}),
	cmpopts.IgnoreUnexported(golinkv1.AddOwnerResponse{}),
	cmpopts.IgnoreUnexported(golinkv1.RemoveOwnerResponse{}),
}

func TestMain(m *testing.M) {
	clearFirestoreEmulator()

	ctx := context.Background()
	var err error
	fsClient, err = firestore.NewClient(ctx, "emulator")
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

func newService() golinkv1connect.GolinkServiceHandler {
	repo := api.NewRepository(fsClient)
	return api.NewGolinkService(repo)
}

func clearFirestoreEmulator() {
	url := fmt.Sprintf("http://%s/emulator/v1/projects/emulator/databases/(default)/documents", os.Getenv("FIRESTORE_EMULATOR_HOST"))
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		panic(err)
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
}

func createGolink(o *dto) {
	ctx := context.Background()
	col := fsClient.Collection("golinks")
	doc := col.Doc(o.ID())

	if _, err := doc.Create(ctx, o); err != nil {
		panic(err)
	}
}

func getGolink(name string) *dto {
	ctx := context.Background()
	col := fsClient.Collection("golinks")
	doc := col.Doc(name)

	snap, err := doc.Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil
	}
	if err != nil {
		panic(err)
	}

	var o dto
	if err := snap.DataTo(&o); err != nil {
		panic(err)
	}

	return &o
}

func TestService_CreateGolink_Success(t *testing.T) {
	defer clearFirestoreEmulator()

	req := &golinkv1.CreateGolinkRequest{
		Name: "link-name",
		Url:  "https://example.com",
	}

	want := &golinkv1.CreateGolinkResponse{
		Golink: &golinkv1.Golink{
			Name:   "link-name",
			Url:    "https://example.com",
			Owners: []string{"user@example.com"},
		},
	}

	ctx := golinkcontext.WithUserEmail(context.Background(), "user@example.com")
	s := newService()
	got, err := s.CreateGolink(ctx, connect.NewRequest(req))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !cmp.Equal(got.Msg, want, cmpOptions...) {
		t.Errorf("unexpected response (-want +got): %v", cmp.Diff(want, got.Msg, cmpOptions...))
	}
}

func TestService_CreateGolink_AlreadyExists(t *testing.T) {
	defer clearFirestoreEmulator()

	req := &golinkv1.CreateGolinkRequest{
		Name: "link-name",
		Url:  "https://example.com",
	}

	ctx := golinkcontext.WithUserEmail(context.Background(), "user@example.com")
	s := newService()

	_, err := s.CreateGolink(ctx, connect.NewRequest(req))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = s.CreateGolink(ctx, connect.NewRequest(req))
	if err, ok := err.(*connect.Error); !ok || err.Code() != connect.CodeAlreadyExists {
		t.Errorf("got %v, want %v", err, connect.CodeAlreadyExists)
	}
}

func TestService_CreateGolink_Validations(t *testing.T) {
	tests := map[string]struct {
		name     string
		url      string
		wantCode connect.Code
	}{
		"valid name and url": {
			name:     "linkname",
			url:      "https://example.com",
			wantCode: connect.Code(0),
		},
		"valid name with underscore": {
			name:     "link_name",
			url:      "https://example.com",
			wantCode: connect.Code(0),
		},
		"valid name with dash": {
			name:     "link-name",
			url:      "https://example.com",
			wantCode: connect.Code(0),
		},
		"valid name with dot": {
			name:     "link.name",
			url:      "https://example.com",
			wantCode: connect.Code(0),
		},
		"long name": {
			name:     strings.Repeat("a", 1500),
			url:      "https://example.com",
			wantCode: connect.Code(0),
		},
		"name with space": {
			name:     "link name",
			url:      "https://example.com",
			wantCode: connect.Code(0),
		},
		"name with non-ascii characters": {
			name:     "リンク名",
			url:      "https://example.com",
			wantCode: connect.Code(0),
		},
		"empty name": {
			name:     "",
			url:      "https://example.com",
			wantCode: connect.CodeInvalidArgument,
		},
		"name with slash": {
			name:     "link/name",
			url:      "https://example.com",
			wantCode: connect.CodeInvalidArgument,
		},
		"name with only dash": {
			name:     "-",
			url:      "https://example.com",
			wantCode: connect.CodeInvalidArgument,
		},
		"name starts with -": {
			name:     "-link-name",
			url:      "https://example.com",
			wantCode: connect.CodeInvalidArgument,
		},
		"name ends with -": {
			name:     "link-name-",
			url:      "https://example.com",
			wantCode: connect.CodeInvalidArgument,
		},
		"name starts with double underscores": {
			name:     "__link-name",
			url:      "https://example.com",
			wantCode: connect.CodeInvalidArgument,
		},
		"name ends with double underscores": {
			name:     "link-name__",
			url:      "https://example.com",
			wantCode: connect.CodeInvalidArgument,
		},
		"name with only underscore": {
			name:     "_",
			url:      "https://example.com",
			wantCode: connect.CodeInvalidArgument,
		},
		"name with only underscores": {
			name:     "__",
			url:      "https://example.com",
			wantCode: connect.CodeInvalidArgument,
		},
		"name solely consist of a single period": {
			name:     ".",
			url:      "https://example.com",
			wantCode: connect.CodeInvalidArgument,
		},
		"name solely consist of double periods": {
			name:     "..",
			url:      "https://example.com",
			wantCode: connect.CodeInvalidArgument,
		},
		"too long name": {
			name:     strings.Repeat("a", 1501),
			url:      "https://example.com",
			wantCode: connect.CodeInvalidArgument,
		},
		"name is c": {
			name:     "c",
			url:      "https://example.com",
			wantCode: connect.CodeInvalidArgument,
		},
		"name is api": {
			name:     "api",
			url:      "https://example.com",
			wantCode: connect.CodeInvalidArgument,
		},
		"empty url": {
			name:     "link-name",
			url:      "",
			wantCode: connect.CodeInvalidArgument,
		},
		"invalid url": {
			name:     "link-name",
			url:      "example.com",
			wantCode: connect.CodeInvalidArgument,
		},
		"invalid url scheme": {
			name:     "link-name",
			url:      "ftp://example.com",
			wantCode: connect.CodeInvalidArgument,
		},
	}

	defer clearFirestoreEmulator()

	s := newService()
	ctx := golinkcontext.WithUserEmail(context.Background(), "user@example.com")

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			defer clearFirestoreEmulator()

			req := &golinkv1.CreateGolinkRequest{
				Name: tt.name,
				Url:  tt.url,
			}

			_, err := s.CreateGolink(ctx, connect.NewRequest(req))

			if tt.wantCode == connect.Code(0) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			} else {
				if err, ok := err.(*connect.Error); !ok || err.Code() != tt.wantCode {
					t.Errorf("got %v, want %v", err, tt.wantCode)
				}
			}
		})
	}
}

func TestService_GetGolink_Success(t *testing.T) {
	defer clearFirestoreEmulator()

	o := &dto{
		Name:   "link-name",
		URL:    "https://example.com",
		Owners: []string{"other@example.com"},
	}
	createGolink(o)

	s := newService()
	ctx := golinkcontext.WithUserEmail(context.Background(), "user@example.com")

	req := &golinkv1.GetGolinkRequest{
		Name: o.Name,
	}

	got, err := s.GetGolink(ctx, connect.NewRequest(req))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := &golinkv1.GetGolinkResponse{Golink: o.ToProto()}

	if !cmp.Equal(got.Msg, want, cmpOptions...) {
		t.Errorf("unexpected response (-want +got): %v", cmp.Diff(want, got.Msg, cmpOptions...))
	}
}

func TestService_GetGolink_NotFound(t *testing.T) {
	defer clearFirestoreEmulator()

	s := newService()
	ctx := golinkcontext.WithUserEmail(context.Background(), "user@example.com")

	req := &golinkv1.GetGolinkRequest{
		Name: "link-name",
	}

	_, err := s.GetGolink(ctx, connect.NewRequest(req))
	if err == nil {
		t.Errorf("error expected")
	}
}

func TestService_ListGolinks(t *testing.T) {
	defer clearFirestoreEmulator()

	dtoOwned1 := &dto{
		Name:   "link-owned-1",
		URL:    "https://example.com",
		Owners: []string{"user@example.com"},
	}
	dtoOwned2 := &dto{
		Name:   "link-owned-2",
		URL:    "https://example.com",
		Owners: []string{"user@example.com", "other@example.com"},
	}
	dtoNotOwned := &dto{
		Name:   "link-not-owned",
		URL:    "https://example.com",
		Owners: []string{"other@example.com"},
	}

	tests := map[string]struct {
		golinks  []*dto
		wantDTOs []*dto
	}{
		"no owned golinks": {
			golinks:  []*dto{dtoNotOwned},
			wantDTOs: []*dto{},
		},
		"all": {
			golinks:  []*dto{dtoOwned1, dtoOwned2, dtoNotOwned},
			wantDTOs: []*dto{dtoOwned1, dtoOwned2},
		},
	}

	s := newService()
	ctx := golinkcontext.WithUserEmail(context.Background(), "user@example.com")

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			defer clearFirestoreEmulator()

			for _, o := range tt.golinks {
				createGolink(o)
			}

			got, err := s.ListGolinks(ctx, connect.NewRequest(&golinkv1.ListGolinksRequest{}))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			wantGolinks := []*golinkv1.Golink{}
			for _, o := range tt.wantDTOs {
				wantGolinks = append(wantGolinks, o.ToProto())
			}
			want := &golinkv1.ListGolinksResponse{Golinks: wantGolinks}

			if !cmp.Equal(got.Msg, want, cmpOptions...) {
				t.Errorf("unexpected response (-want +got): %v", cmp.Diff(want, got.Msg, cmpOptions...))
			}
		})
	}
}

func TestService_ListGolinksByURL(t *testing.T) {
	defer clearFirestoreEmulator()

	o1 := &dto{
		Name:   "o1",
		URL:    "https://example.com/1",
		Owners: []string{"user@example.com"},
	}
	o2 := &dto{
		Name:   "o2",
		URL:    "https://example.com/2",
		Owners: []string{"user@example.com", "other@example.com"},
	}
	o3 := &dto{
		Name:   "o3",
		URL:    "https://example.com/1",
		Owners: []string{"other@example.com"},
	}
	createGolink(o1)
	createGolink(o2)
	createGolink(o3)

	tests := map[string]struct {
		url      string
		wantDTOs []*dto
	}{
		"two golinks": {
			url:      "https://example.com/1",
			wantDTOs: []*dto{o1, o3},
		},
		"no golinks": {
			url:      "https://example.com/3",
			wantDTOs: []*dto{},
		},
	}

	s := newService()
	ctx := golinkcontext.WithUserEmail(context.Background(), "user@example.com")

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			req := &golinkv1.ListGolinksByUrlRequest{Url: tt.url}
			got, err := s.ListGolinksByUrl(ctx, connect.NewRequest(req))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			wantGolinks := []*golinkv1.Golink{}
			for _, o := range tt.wantDTOs {
				wantGolinks = append(wantGolinks, o.ToProto())
			}
			want := &golinkv1.ListGolinksByUrlResponse{Golinks: wantGolinks}

			if !cmp.Equal(got.Msg, want, cmpOptions...) {
				t.Errorf("unexpected response (-want +got): %v", cmp.Diff(want, got.Msg, cmpOptions...))
			}
		})
	}
}

func TestService_UpdateGolink_Success(t *testing.T) {
	defer clearFirestoreEmulator()

	o := &dto{
		Name:   "link-name",
		URL:    "https://example.com",
		Owners: []string{"user@example.com"},
	}
	createGolink(o)

	s := newService()
	ctx := golinkcontext.WithUserEmail(context.Background(), "user@example.com")

	req := &golinkv1.UpdateGolinkRequest{
		Name: o.Name,
		Url:  "https://example.com/updated",
	}

	got, err := s.UpdateGolink(ctx, connect.NewRequest(req))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	want := &golinkv1.UpdateGolinkResponse{Golink: o.ToProto()}
	want.Golink.Url = "https://example.com/updated"

	if !cmp.Equal(got.Msg, want, cmpOptions...) {
		t.Errorf("unexpected response (-want +got): %v", cmp.Diff(want, got.Msg, cmpOptions...))
	}

	gotSaved := getGolink(o.Name)
	if gotSaved == nil {
		t.Fatalf("golink not saved")
	}

	wantSaved := &dto{
		Name:   o.Name,
		URL:    "https://example.com/updated",
		Owners: []string{"user@example.com"},
	}
	if !cmp.Equal(wantSaved, gotSaved, cmpOptions...) {
		t.Errorf("unexpected saved golink (-want +got): %v", cmp.Diff(wantSaved, gotSaved, cmpOptions...))
	}
}

func TestService_UpdateGolink_PermissionDenied(t *testing.T) {
	defer clearFirestoreEmulator()

	o := &dto{
		Name:   "link-name",
		URL:    "https://example.com",
		Owners: []string{"other@example.com"},
	}
	createGolink(o)

	s := newService()
	ctx := golinkcontext.WithUserEmail(context.Background(), "user@example.com")

	req := &golinkv1.UpdateGolinkRequest{
		Name: o.Name,
		Url:  "https://example.com/updated",
	}

	_, err := s.UpdateGolink(ctx, connect.NewRequest(req))
	if err, ok := err.(*connect.Error); !ok || err.Code() != connect.CodePermissionDenied {
		t.Errorf("got %v, want %v", err, connect.CodePermissionDenied)
	}
}

func TestService_UpdateGolink_NotFound(t *testing.T) {
	defer clearFirestoreEmulator()

	s := newService()
	ctx := golinkcontext.WithUserEmail(context.Background(), "user@example.com")

	req := &golinkv1.UpdateGolinkRequest{
		Name: "link-name",
		Url:  "https://example.com/updated",
	}

	_, err := s.UpdateGolink(ctx, connect.NewRequest(req))
	if err, ok := err.(*connect.Error); !ok || err.Code() != connect.CodeNotFound {
		t.Errorf("got %v, want %v", err, connect.CodeNotFound)
	}
}

func TestService_DeleteGolink_Success(t *testing.T) {
	defer clearFirestoreEmulator()

	o := &dto{
		Name:   "link-name",
		URL:    "https://example.com",
		Owners: []string{"user@example.com"},
	}
	createGolink(o)

	s := newService()
	ctx := golinkcontext.WithUserEmail(context.Background(), "user@example.com")

	req := &golinkv1.DeleteGolinkRequest{Name: o.Name}

	_, err := s.DeleteGolink(ctx, connect.NewRequest(req))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	saved := getGolink(o.Name)
	if saved != nil {
		t.Errorf("golink should be deleted")
	}
}

func TestService_DeleteGolink_PermissionDenied(t *testing.T) {
	defer clearFirestoreEmulator()

	o := &dto{
		Name:   "link-name",
		URL:    "https://example.com",
		Owners: []string{"other@example.com"},
	}
	createGolink(o)

	s := newService()
	ctx := golinkcontext.WithUserEmail(context.Background(), "user@example.com")

	req := &golinkv1.DeleteGolinkRequest{Name: o.Name}

	_, err := s.DeleteGolink(ctx, connect.NewRequest(req))
	if err, ok := err.(*connect.Error); !ok || err.Code() != connect.CodePermissionDenied {
		t.Errorf("got %v, want %v", err, connect.CodePermissionDenied)
	}
}

func TestService_DeleteGolink_NotFound(t *testing.T) {
	defer clearFirestoreEmulator()

	s := newService()
	ctx := golinkcontext.WithUserEmail(context.Background(), "user@example.com")

	req := &golinkv1.DeleteGolinkRequest{Name: "link-name"}

	_, err := s.DeleteGolink(ctx, connect.NewRequest(req))
	if err, ok := err.(*connect.Error); !ok || err.Code() != connect.CodeNotFound {
		t.Errorf("got %v, want %v", err, connect.CodeNotFound)
	}
}

func TestService_AddOwner(t *testing.T) {
	noErr := connect.Code(0)

	tests := map[string]struct {
		originalOwners []string
		ownerToAdd     string
		wantOwners     []string
		wantErr        connect.Code
	}{
		"add owner": {
			originalOwners: []string{"user@example.com"},
			ownerToAdd:     "other@example.com",
			wantOwners:     []string{"user@example.com", "other@example.com"},
			wantErr:        noErr,
		},
		"add owner when already exists": {
			originalOwners: []string{"user@example.com", "other@example.com"},
			ownerToAdd:     "other@example.com",
			wantOwners:     nil,
			wantErr:        connect.CodeInvalidArgument,
		},
		"permission denied": {
			originalOwners: []string{"other@example.com"},
			ownerToAdd:     "user@example.com",
			wantOwners:     nil,
			wantErr:        connect.CodePermissionDenied,
		},
	}

	s := newService()
	ctx := golinkcontext.WithUserEmail(context.Background(), "user@example.com")

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			defer clearFirestoreEmulator()

			o := &dto{
				Name:   "link-name",
				URL:    "https://example.com",
				Owners: tt.originalOwners,
			}
			createGolink(o)

			req := &golinkv1.AddOwnerRequest{Name: o.Name, Owner: tt.ownerToAdd}

			got, err := s.AddOwner(ctx, connect.NewRequest(req))

			if tt.wantErr != noErr {
				if err, ok := err.(*connect.Error); !ok || err.Code() != tt.wantErr {
					t.Errorf("err got %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			wantDTO := &dto{
				Name:   o.Name,
				URL:    o.URL,
				Owners: tt.wantOwners,
			}

			want := &golinkv1.AddOwnerResponse{Golink: wantDTO.ToProto()}

			if !cmp.Equal(want, got.Msg, cmpOptions...) {
				t.Errorf("unexpected response (-want +got): %v", cmp.Diff(want, got.Msg, cmpOptions...))
			}

			saved := getGolink(o.Name)
			if !cmp.Equal(wantDTO, saved, cmpOptions...) {
				t.Errorf("save failed (-want +got): %v", cmp.Diff(wantDTO, saved))
			}
		})
	}
}

func TestService_RemoveOwner(t *testing.T) {
	noErr := connect.Code(0)

	tests := map[string]struct {
		originalOwners []string
		ownerToRemove  string
		wantOwners     []string
		wantErr        connect.Code
	}{
		"success": {
			originalOwners: []string{"user@example.com", "other@example.com"},
			ownerToRemove:  "other@example.com",
			wantOwners:     []string{"user@example.com"},
			wantErr:        noErr,
		},
		"only one user remains": {
			originalOwners: []string{"user@example.com"},
			ownerToRemove:  "user@example.com",
			wantOwners:     nil,
			wantErr:        connect.CodeInvalidArgument,
		},
		"owner not found": {
			originalOwners: []string{"user@example.com"},
			ownerToRemove:  "other@example.com",
			wantOwners:     nil,
			wantErr:        connect.CodeInvalidArgument,
		},
		"permission denied": {
			originalOwners: []string{"other@example.com"},
			ownerToRemove:  "user@example.com",
			wantOwners:     nil,
			wantErr:        connect.CodePermissionDenied,
		},
	}

	s := newService()
	ctx := golinkcontext.WithUserEmail(context.Background(), "user@example.com")

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			defer clearFirestoreEmulator()

			o := &dto{
				Name:   "link-name",
				URL:    "https://example.com",
				Owners: tt.originalOwners,
			}
			createGolink(o)

			req := &golinkv1.RemoveOwnerRequest{Name: o.Name, Owner: tt.ownerToRemove}

			got, err := s.RemoveOwner(ctx, connect.NewRequest(req))

			if tt.wantErr != noErr {
				if err, ok := err.(*connect.Error); !ok || err.Code() != tt.wantErr {
					t.Errorf("err got %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			wantDTO := &dto{
				Name:   o.Name,
				URL:    o.URL,
				Owners: tt.wantOwners,
			}

			want := &golinkv1.RemoveOwnerResponse{Golink: wantDTO.ToProto()}

			if !cmp.Equal(want, got.Msg, cmpOptions...) {
				t.Errorf("unexpected response (-want +got): %v", cmp.Diff(want, got.Msg, cmpOptions...))
			}

			saved := getGolink(o.Name)
			if !cmp.Equal(wantDTO, saved, cmpOptions...) {
				t.Errorf("save failed (-want +got): %v", cmp.Diff(wantDTO, saved))
			}
		})
	}
}
