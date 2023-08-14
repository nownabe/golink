package main

import (
	"context"
	"os"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/firestore"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/bufbuild/connect-go"
	"github.com/nownabe/golink/go/clog"
	"github.com/nownabe/golink/go/errors"
	"github.com/nownabe/golink/go/interceptors"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"github.com/nownabe/golink/api"
)

func main() {
	clog.SetDefault(clog.New(os.Stdout, clog.LevelInfo))

	ctx := context.Background()

	a, err := buildAPI(ctx)
	if err != nil {
		clog.AlertErr(ctx, errors.Wrap(err, "failed to build API"))
		os.Exit(1)
	}

	if err := a.Run(ctx); err != nil {
		clog.AlertErr(ctx, errors.Wrap(err, "failed to run server"))
		os.Exit(1)
	}
}

func buildAPI(ctx context.Context) (api.API, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	projectID, err := getProjectID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project ID")
	}
	clog.Infof(ctx, "Project ID: %s", projectID)
	clog.SetContextHandler(projectID)

	origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	clog.Infof(ctx, "Allowed origins: %v", origins)

	fsClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Firestore client")
	}

	if err := setOtel(ctx, projectID); err != nil {
		return nil, errors.Wrap(err, "failed to get tracer")
	}

	repo := api.NewRepository(fsClient)
	svc := api.NewGolinkService(repo)
	apiInterceptors := []connect.Interceptor{
		// outermost
		interceptors.NewRecoverer(),
		interceptors.WithTracer(),
		interceptors.NewRequestID(),
		interceptors.NewAuthorizer(),
		interceptors.NewLogger(),
		// innermost
	}

	debug := false
	if isDebug := strings.ToLower(os.Getenv("DEBUG")); isDebug == "true" {
		debug = true
	}

	if user := os.Getenv("USE_DUMMY_USER"); user != "" {
		u := strings.Split(user, ":")
		apiInterceptors = append([]connect.Interceptor{interceptors.NewDummyUser(u[0], u[1])}, apiInterceptors...)
	}

	return api.New(svc, port, "/api", origins, apiInterceptors, debug), nil
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

func setOtel(ctx context.Context, projectID string) error {
	exporter, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		return errors.Wrapf(err, "failed to create exporter with project ID %s", projectID)
	}

	res, err := resource.New(ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("golink-api"),
		),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create resource")
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return nil
}
