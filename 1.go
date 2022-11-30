package main

import (
	"encoding/json"
	"fmt"
	mydata "gorilla/internal/data"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type SomeStruct struct {
	Name  string
	Email string
}

func StreamHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("string handler invoked")
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Println("responseWriter is not really a flusher")
		return
	}
	//this header had no effect
	w.Header().Set("Connection", "Keep-Alive")
	//these two headers are needed to get the http chunk incremently
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	for i := 0; i < 20; i++ {
		// w.Write([]byte("Gorilla! \n"))
		fmt.Println(i)
		fmt.Fprintf(w, "Gorilla! %v \n", i)
		flusher.Flush()
		time.Sleep(1 * time.Second)
		// time.Sleep(1 * time.Second)
	}
	fmt.Println("done")
}

func StringHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("string handler invoked")
	w.Write([]byte("Gorilla!\n"))
}

func JsonStringHandler(w http.ResponseWriter, r *http.Request) {
	j := `{"name": "najam awan", "email":"najamsk@gmail.com"}`
	w.Write([]byte(j))
}

func JsonStructHandler(w http.ResponseWriter, r *http.Request) {
	data := SomeStruct{Name: "najam", Email: "najamsk@gmail.com"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func MembersHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:route GET /members members listMembers
	// Return a list of memebers from the database
	// responses:
	//	200: membersResponse

	repo := mydata.NewRepo()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repo.GetMembers())
}

func TeamsHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:route GET /teams teams listTeams
	// Return a list of teams from the database
	// responses:
	//	200: teamsResponse

	// repo := mystore.NewRepo()
	repo := mydata.NewRepo()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repo.GetTeams())
}

func NewMemberHandler(w http.ResponseWriter, r *http.Request) {
	var newTeam mydata.Team
	rb := json.NewDecoder(r.Body)
	err := rb.Decode(&newTeam)
	if err != nil {
		log.Println("cant parse incoming data for member")
		// w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("New member handler invoked")
	repo := mydata.NewRepo()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repo.GetMembers())
}

func JsonMapHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Status Created"
	resp["topic"] = "user/request"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

}

func logHandler(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("logging middleware start")
		next.ServeHTTP(w, r)
		log.Println("logging middleware ends")
	})
}

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	fs := http.FileServer(http.Dir("./swaggerui/"))
	r.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))
	sf := http.HandlerFunc(StringHandler)
	r.HandleFunc("/", logHandler(sf)).Methods("GET")
	r.HandleFunc("/member", NewMemberHandler).Methods("POST")
	r.HandleFunc("/jsonstring", JsonStringHandler)
	r.HandleFunc("/struct", JsonStructHandler)
	r.HandleFunc("/map", JsonMapHandler)
	r.HandleFunc("/stream", StreamHandler)
	r.HandleFunc("/teams", TeamsHandler)
	r.HandleFunc("/members", MembersHandler)

	// Bind to a port and pass our router in
	log.Println("server started at :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
