package api_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/bufbuild/connect-go"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/nownabe/golink/api"
	golinkv1 "github.com/nownabe/golink/api/gen/golink/v1"
	"github.com/nownabe/golink/api/gen/golink/v1/golinkv1connect"
)

type dto = api.DTO

var fsClient *firestore.Client

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

var cmpOptions = []cmp.Option{
	cmpopts.IgnoreUnexported(golinkv1.Golink{}),
	cmpopts.IgnoreUnexported(golinkv1.CreateGolinkResponse{}),
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

	ctx := api.WithUserEmail(context.Background(), "user@example.com")
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

	ctx := api.WithUserEmail(context.Background(), "user@example.com")
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
	ctx := api.WithUserEmail(context.Background(), "user@example.com")

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

func createGolink(o *dto) {
	ctx := context.Background()
	col := fsClient.Collection("golinks")
	doc := col.Doc(o.ID())

	if _, err := doc.Create(ctx, o); err != nil {
		panic(err)
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
	ctx := api.WithUserEmail(context.Background(), "user@example.com")

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

			var want []*golinkv1.Golink
			for _, o := range tt.wantDTOs {
				want = append(want, o.ToProto())
			}

			if !cmp.Equal(got.Msg.Golinks, want, cmpOptions...) {
				t.Errorf("unexpected response (-want +got): %v", cmp.Diff(want, got.Msg, cmpOptions...))
			}
		})
	}
}
