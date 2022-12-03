package handlers

import (
	"encoding/json"
	"fmt"
	"gorilla/internal/data"
	"log"
	"net/http"
	"time"
)

type RestHandler struct {
	Repo *data.Repo
}

func (h *RestHandler) SayNameMethod(w http.ResponseWriter, r *http.Request) {
	log.Println("SayNameReal invoked")
	w.Write([]byte("my name is real slim shady\n"))
}

func (h *RestHandler) MembersHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:route GET /members members listMembers
	// Return a list of memebers from the database
	// responses:
	//	200: membersResponse

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.Repo.GetMembers())
}

func (h *RestHandler) NewMemberHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:route POST /member members createMember
	// Return a newly created member
	// responses:
	//	200: memberResponse
	//  501: errorResponse

	var newMember data.Member
	rb := json.NewDecoder(r.Body)
	err := rb.Decode(&newMember)
	if err != nil {
		log.Println("cant parse incoming data for member")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newMember = h.Repo.CreateMember(newMember)
	log.Printf("created emmber: %#v \n", newMember)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newMember)
}

func (h *RestHandler) UpdateMemberHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:route PUT /member members updateMember
	// Return a updated member
	// responses:
	//	200: memberResponse
	//  400: errorResponse

	var newMember data.Member
	rb := json.NewDecoder(r.Body)
	err := rb.Decode(&newMember)
	if err != nil {
		log.Println("cant parse incoming data for member")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	uMember, err := h.Repo.UpdateMember(newMember)
	if err != nil {
		log.Printf("update member failed with error:%#v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("updated memeber: %#v \n", uMember)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(uMember)
}

func (h *RestHandler) NewTeamHandler(w http.ResponseWriter, r *http.Request) {
	var n data.Team
	rb := json.NewDecoder(r.Body)
	err := rb.Decode(&n)
	if err != nil {
		log.Println("cant parse incoming data for team")
		// w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	n = h.Repo.CreateTeam(n)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(n)
}

func (h *RestHandler) TeamsHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:route GET /teams teams listTeams
	// Return a list of teams from the database
	// responses:
	//	200: teamsResponse

	// repo := mystore.NewRepo()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.Repo.GetTeams())
}

func (h *RestHandler) StringHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("string handler invoked")
	w.Write([]byte("Gorilla!\n"))
}

func (h *RestHandler) Err501(w http.ResponseWriter, r *http.Request) {
	// swagger:route GET /501 dev err501
	// Return a not implemented error
	// responses:
	//	501: errorResponse
	log.Println("Err501 handler invoked")
	http.Error(w, "server failed", http.StatusNotImplemented)
}

func (h *RestHandler) JsonStringHandler(w http.ResponseWriter, r *http.Request) {
	j := `{"name": "najam awan", "email":"najamsk@gmail.com"}`
	w.Write([]byte(j))
}

type SomeStruct struct {
	Name  string
	Email string
}

func (h *RestHandler) JsonStructHandler(w http.ResponseWriter, r *http.Request) {
	data := SomeStruct{Name: "najam", Email: "najamsk@gmail.com"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h *RestHandler) JsonMapHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Status Created"
	resp["topic"] = "user/request"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

}

func (h *RestHandler) LogHandler(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("logging middleware start")
		next.ServeHTTP(w, r)
		log.Println("logging middleware ends")
	})
}

func (h *RestHandler) StreamHandler(w http.ResponseWriter, r *http.Request) {
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
