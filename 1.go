package main

import (
	"context"
	"gorilla/handlers"
	"gorilla/internal/data"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

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

func main() {
	//setup logger
	l := log.New(os.Stdout, "", 0)

	// Write telemetry data to a file.
	f, err := os.Create("traces.txt")
	if err != nil {
		l.Fatal(err)
	}
	defer f.Close()

	exp, err := newExporter(f)
	if err != nil {
		l.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			l.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)

	//setup database
	db := data.NewDataStore()
	repo := data.NewRepo(db)
	r := mux.NewRouter()
	/* restHandler := &handlers.RestHandler{
		Repo: repo,
	} */
	restHandler := handlers.NewResHandler(repo, l)
	// Routes consist of a path and a handler function.
	fs := http.FileServer(http.Dir("./swaggerui/"))
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
	r.HandleFunc("/stream", restHandler.StreamHandler)
	r.HandleFunc("/jsonstring", restHandler.JsonStringHandler)
	r.HandleFunc("/struct", restHandler.JsonStructHandler)
	r.HandleFunc("/501", restHandler.Err501).Methods("GET")

	// Bind to a port and pass our router in
	l.Println("server started at :8000")
	l.Fatal(http.ListenAndServe(":8000", r))

}
