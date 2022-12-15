package main

import (
	"gorilla/handlers"
	"gorilla/internal/data"
	"io/fs"
	"log"
	"net/http"

	"embed"

	"github.com/gorilla/mux"
)

//go:embed swaggerui
var SwaggerDir embed.FS

func main() {
	//setup database
	db := data.NewDataStore()
	repo := data.NewRepo(db)
	r := mux.NewRouter()
	restHandler := &handlers.RestHandler{
		Repo: repo,
	}
	// Routes consist of a path and a handler function.
	// fs := http.FileServer(http.Dir("./swaggerui/"))

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
	r.HandleFunc("/stream", restHandler.StreamHandler)
	r.HandleFunc("/jsonstring", restHandler.JsonStringHandler)
	r.HandleFunc("/struct", restHandler.JsonStructHandler)
	r.HandleFunc("/501", restHandler.Err501).Methods("GET")

	// Bind to a port and pass our router in
	log.Println("server started at :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
