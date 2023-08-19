package main

import (
	"context"
	"os"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/firestore"
	"github.com/nownabe/golink/go/clog"
	"github.com/nownabe/golink/go/errors"

	"github.com/nownabe/golink/backend"
)

func main() {
	clog.SetDefault(clog.New(os.Stdout, clog.LevelInfo))

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
	clog.Infof(ctx, "project ID: %s", projectID)

	fsClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		clog.AlertErr(ctx, errors.Wrap(err, "failed to create firestore client"))
		os.Exit(1)
	}

	app := backend.New(port, "/c/", fsClient)

	if err := app.Run(ctx); err != nil {
		clog.AlertErr(ctx, errors.Wrap(err, "failed to run redirector"))
	}
}

func getProjectID(ctx context.Context) (string, error) {
	projectID := os.Getenv("PROJECT_ID")
	if projectID != "" {
		return projectID, nil
	}

	// Get project ID from metadata server
	os.Setenv("GCE_METADATA_HOST", "metadata.google.internal")
	projectID, err := metadata.ProjectID()
	if err != nil {
		return "", errors.Wrap(err, "failed to get project ID from metadata server")
	}

	return projectID, nil
}
