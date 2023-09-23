package main

import (
	"context"
	"os"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/firestore"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.nownabe.dev/clog"
	"go.nownabe.dev/clog/errors"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/nownabe/golink/backend"
	"github.com/nownabe/golink/backend/middleware"
)

func main() {
	ctx := context.Background()

	logger := clog.New(os.Stdout, clog.SeverityInfo, true,
		clog.WithHandleFunc(middleware.RequestIDHandleFunc),
	)
	clog.SetDefault(logger)

	app, err := buildApp(ctx)
	if err != nil {
		clog.AlertErr(ctx, errors.Errorf("failed to build app: %w", err))
		os.Exit(1)
	}

	if err := app.Run(ctx); err != nil {
		clog.AlertErr(ctx, errors.Errorf("failed to run redirector: %w", err))
	}
}

func buildApp(ctx context.Context) (backend.App, error) {
	projectID, err := getProjectID()
	if err != nil {
		clog.WarningErr(ctx, err)
		projectID = "dummy-project"
	} else {
		clog.SetOptions(clog.WithTrace(projectID))
	}
	clog.Infof(ctx, "projectID=%q", projectID)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	clog.Infof(ctx, "Allowed origins: %v", origins)

	fsClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.Errorf("failed to create Firestore client: %w", err)
	}

	if err := setOtel(ctx, projectID); err != nil {
		return nil, errors.Errorf("failed to get tracer: %w", err)
	}

	ldcfg := backend.LocalDevelopmentConfig{
		LocalConsoleURL: os.Getenv("LOCAL_CONSOLE_URL"),
		DebugEndpoint:   strings.ToLower(os.Getenv("DEBUG")) == "true",
		DummyUserEmail:  os.Getenv("DUMMY_USER_EMAIL"),
		DummyUserID:     os.Getenv("DUMMY_USER_ID"),
	}

	return backend.New(port, origins, "golink-backend", "/api", "/-/", fsClient, ldcfg), nil
}

func getProjectID() (string, error) {
	projectID := os.Getenv("PROJECT_ID")
	if projectID != "" {
		return projectID, nil
	}

	// Get project ID from metadata server
	os.Setenv("GCE_METADATA_HOST", "metadata.google.internal")
	projectID, err := metadata.ProjectID()
	if err != nil {
		return "", errors.Errorf("failed to get project ID from metadata server: %w", err)
	}

	return projectID, nil
}

func setOtel(ctx context.Context, projectID string) error {
	exporter, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		return errors.Errorf("failed to create exporter with project ID %s: %w", projectID, err)
	}

	res, err := resource.New(ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("golink-api"),
		),
	)
	if err != nil {
		return errors.Errorf("failed to create resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return nil
}
