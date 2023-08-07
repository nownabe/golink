package api_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
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

// TODO: TestService_CreateGolink_Validations
