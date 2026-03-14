package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MarioFronza/media-tui/internal/adapter/api"
)

func TestGet_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") != "test-key" {
			t.Errorf("expected X-Api-Key header, got %q", r.Header.Get("X-Api-Key"))
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"title": "Inception"})
	}))
	defer srv.Close()

	client := api.NewClient(srv.URL, "test-key")

	var result map[string]string
	if err := client.Get("/api/v3/movie", &result); err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if result["title"] != "Inception" {
		t.Errorf("unexpected title: %s", result["title"])
	}
}

func TestGet_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := api.NewClient(srv.URL, "bad-key")

	var result any
	if err := client.Get("/api/v3/movie", &result); err == nil {
		t.Fatal("expected error on 401 response")
	}
}

func TestPost_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") != "test-key" {
			t.Errorf("expected X-Api-Key header, got %q", r.Header.Get("X-Api-Key"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", r.Header.Get("Content-Type"))
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	client := api.NewClient(srv.URL, "test-key")

	if err := client.Post("/api/v3/movie", map[string]string{"title": "Inception"}); err != nil {
		t.Fatalf("Post() error: %v", err)
	}
}

func TestPost_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := api.NewClient(srv.URL, "test-key")

	if err := client.Post("/api/v3/movie", map[string]string{}); err == nil {
		t.Fatal("expected error on 500 response")
	}
}
