package main

import (
	"context"
	"gorilla/handlers"
	"gorilla/internal/data"
	"io"
	"io/fs"
	"log"
	"net/http"
	"time"

	"embed"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// these constants for openTelemetry jaeger
const (
	service     = "gorilla-demo"
	environment = "demo"
	id          = 1
)

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProvider(url string) (*trace.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := trace.NewTracerProvider(
		// Always be sure to batch in production.
		trace.WithBatcher(exp),
		// Record information about this application in a Resource.
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}

// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("GorillaDemoApi"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}

// newExporter returns a console exporter.
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

//go:embed swaggerui
var SwaggerDir embed.FS

func main() {
	//zap logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	//setup database
	db := data.NewDataStore()
	repo := data.NewRepo(db, logger)
	r := mux.NewRouter()
	restHandler := handlers.NewResHandler(repo, logger)

	//embeding swaggerui folder and serving it as a fileserver
	swagFS := fs.FS(SwaggerDir)
	swaggerContent, _ := fs.Sub(swagFS, "swaggerui")
	fs := http.FileServer(http.FS(swaggerContent))

	r.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))
	sf := http.HandlerFunc(restHandler.StringHandler)
	r.HandleFunc("/", restHandler.LogHandler(sf)).Methods("GET")
	r.HandleFunc("/member", restHandler.NewMemberHandler).Methods("POST")
	r.HandleFunc("/member", restHandler.UpdateMemberHandler).Methods("PUT")
	r.HandleFunc("/member/{memid}", restHandler.DeleteMemberHandler).Methods("DELETE")
	r.HandleFunc("/team", restHandler.NewTeamHandler).Methods("POST")
	r.HandleFunc("/real", restHandler.SayNameMethod).Methods("GET")
	r.HandleFunc("/teams", restHandler.TeamsHandler)
	r.HandleFunc("/members", restHandler.MembersHandler)
	r.HandleFunc("/map", restHandler.JsonMapHandler)
	r.HandleFunc("/{id}", handlers.MakeAPIFunc(restHandler.ResourceErrHandler))
	r.HandleFunc("/stream", restHandler.StreamHandler)
	r.HandleFunc("/jsonstring", restHandler.JsonStringHandler)
	r.HandleFunc("/struct", restHandler.JsonStructHandler)
	r.HandleFunc("/501", restHandler.Err501).Methods("GET")

	// Bind to a port and pass our router in
	logger.Info("server started at :8000")
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		logger.Fatal("listening to server failed:", zap.Error(err))
	}

}
