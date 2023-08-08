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

func newService() golinkv1connect.GolinkServiceHandler {
	ctx := context.Background()
	fsClient, err := firestore.NewClient(ctx, "expecting-emulator")
	if err != nil {
		panic(err)
	}

	repo := api.NewRepository(fsClient)

	return api.NewGolinkService(repo)
}

func clearFirestoreEmulator() {
	url := fmt.Sprintf("http://%s/emulator/v1/projects/expecting-emulator/databases/(default)/documents", os.Getenv("FIRESTORE_EMULATOR_HOST"))
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
