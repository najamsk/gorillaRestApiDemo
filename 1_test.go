package main

import (
	"fmt"
	"gorilla/handlers"
	"gorilla/internal/data"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func Add(x, y int) int {
	// if y == 0 {
	// 	panic("y == 0")
	// }
	return x + y
}
func FuzzHandler(f *testing.F) {
	tests := []struct {
		x    int
		y    int
		want int
	}{
		{x: 1, y: 2, want: 3},
		{x: 2, y: 6, want: 8},
	}
	for _, tt := range tests {
		f.Add(tt.x, tt.y)
	}
	f.Fuzz(func(t *testing.T, x, y int) {
		got := Add(x, y)
		fmt.Println("got:", got)
	})

}

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
	defer resp.Body.Close()
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

func TestAdd(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}
