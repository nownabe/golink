package main

import (
	"context"
	"os"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/firestore"
	"github.com/nownabe/golink/go/clog"
	"github.com/nownabe/golink/go/errors"

	"github.com/nownabe/golink/api"
)

func main() {
	// TODO: Configure clog

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	projectID, err := getProjectID(ctx)
	if err != nil {
		clog.AlertErr(ctx, err)
		os.Exit(1)
	}

	origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	fsClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		clog.AlertErr(ctx, errors.Wrap(err, "failed to create Firestore client"))
		os.Exit(1)
	}

	repo := api.NewRepository(fsClient)

	svc := api.NewGolinkService(repo)
	if err := api.New(svc, port, "/api", origins).Run(); err != nil {
		clog.AlertErr(ctx, errors.Wrap(err, "failed to run server"))
		os.Exit(1)
	}
}

func getProjectID(ctx context.Context) (string, error) {
	projectID := os.Getenv("PROJECT_ID")
	if projectID != "" {
		return projectID, nil
	}

	// Get project ID from metadata server
	clog.Infof(ctx, "GCE_METADATA_HOST=%s", os.Getenv("GCE_METADATA_HOST"))
	os.Setenv("GCE_METADATA_HOST", "metadata.google.internal")
	projectID, err := metadata.ProjectID()
	if err != nil {
		return "", errors.Wrap(err, "failed to get project ID from metadata server")
	}

	return projectID, nil
}
