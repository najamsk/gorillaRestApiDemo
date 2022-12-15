package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"gorilla/internal/data"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

const name = "handlers"

type RestHandler struct {
	Repo *data.Repo
	Log  *zap.Logger
}

func NewResHandler(r *data.Repo, l *zap.Logger) *RestHandler {
	return &RestHandler{
		Repo: r,
		Log:  l,
	}
}

func (h *RestHandler) SayNameMethod(w http.ResponseWriter, r *http.Request) {
	// h.Log.Println("SayNameReal invoked")
	h.Log.Info("SayNameReal invoked")
	w.Write([]byte("my name is real slim shady\n"))
}

func (h *RestHandler) MembersHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:route GET /members members listMembers
	// Return a list of memebers from the database
	// responses:
	//	200: membersResponse

	ctx := context.Background()
	newCtx, span := otel.Tracer(name).Start(ctx, "handlers/GetMembers")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.Repo.GetMembers(newCtx))
	span.End()
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
		h.Log.Info("cant parse incoming data for member")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newMember = h.Repo.CreateMember(newMember)
	h.Log.Info("created emmber:",
		zap.Int("id", newMember.Id),
		zap.String("name", newMember.Name),
		zap.String("email", newMember.Email),
		zap.Int("teamid", newMember.TeamId),
	)
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
		h.Log.Info("cant parse incoming data for member")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	uMember, err := h.Repo.UpdateMember(newMember)
	if err != nil {
		// h.Log.Info("update member failed with error:%#v", err)
		h.Log.Error("update member failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.Log.Info("updated memeber: %#v \n",
		zap.Int("id", uMember.Id),
		zap.String("name", uMember.Name),
		zap.String("email", uMember.Email),
		zap.Int("teamid", uMember.TeamId),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(uMember)
}

func (h *RestHandler) DeleteMemberHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:route DELETE /member/{memid} members delMember
	// Deletes a member
	// responses:
	//	200: noContentResponse
	//  400: errorResponse

	ctx := context.Background()
	newCtx, span := otel.Tracer(name).Start(ctx, "handlers/DeleteMember")
	defer span.End()

	params := mux.Vars(r)
	id := params["memid"]
	mID, err := strconv.Atoi(id)
	if err != nil {
		// ... handle error
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.Log.Error("Can't parse the requested ID:", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteMember(newCtx, mID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.Log.Error("Deleting member failed with error:", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.Log.Info("Deleted memeber: %#v \n",
		zap.Int("id", mID),
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	// json.NewEncoder(w).Encode(uMember)
}

func (h *RestHandler) NewTeamHandler(w http.ResponseWriter, r *http.Request) {
	var n data.Team
	rb := json.NewDecoder(r.Body)
	err := rb.Decode(&n)
	if err != nil {
		h.Log.Info("cant parse incoming data for team")
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
	ctx := context.Background()
	newCtx, span := otel.Tracer(name).Start(ctx, "handlers/GetTeams")
	json.NewEncoder(w).Encode(h.Repo.GetTeams(newCtx))
	span.End()
}

func (h *RestHandler) StringHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("string handler invoked")
	w.Write([]byte("Gorilla!\n"))
}

func (h *RestHandler) Err501(w http.ResponseWriter, r *http.Request) {
	// swagger:route GET /501 dev err501
	// Return a not implemented error
	// responses:
	//	501: errorResponse
	h.Log.Info("Err501 handler invoked")
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

func (f *SomeStruct) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name:", f.Name)
	enc.AddString("email", f.Email)
	return nil
}

func (h *RestHandler) JsonStructHandler(w http.ResponseWriter, r *http.Request) {
	data := SomeStruct{Name: "najam", Email: "najamsk@gmail.com"}
	// h.Log.Infow("struct to json handler invoked")
	h.Log.Info("struct to json:", zap.Object("someStruct", &data))

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
		h.Log.Fatal("Error happened with json marshal", zap.Error(err))
	}
	w.Write(jsonResp)

}

func (h *RestHandler) LogHandler(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Log.Info("logging middleware start")
		next.ServeHTTP(w, r)
		h.Log.Info("logging middleware ends")
	})
}

func (h *RestHandler) StreamHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("string handler invoked")
	flusher, ok := w.(http.Flusher)
	if !ok {
		h.Log.Info("responseWriter is not really a flusher")
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
