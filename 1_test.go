package main

import (
	"gorilla/handlers"
	"gorilla/internal/data"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestRootStringHandlerEndpoint(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db := data.NewDataStore()
	repo := data.NewRepo(db, logger)
	restHandler := handlers.NewResHandler(repo, logger)

	url := "/"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	want := "Gorilla!\n"

	//Act
	restHandler.StringHandler(w, req)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	got := string(body)

	//Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("mismatched status code with url:%s, want: %v, got:%v \n", url, http.StatusOK, resp.StatusCode)
	}

	if got != want {
		t.Errorf("mismatch body with url:%s, want:%s, got:%s \n", url, want, got)
	}
}
