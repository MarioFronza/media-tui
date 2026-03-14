package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MarioFronza/media-tui/internal/adapter/api"
	"github.com/MarioFronza/media-tui/internal/domain"
)

func newServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

func encode(t *testing.T, w http.ResponseWriter, v any) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Errorf("encode response: %v", err)
	}
}

// --- Radarr ---

func TestRadarrAdapter_MediaType(t *testing.T) {
	a := api.NewRadarrAdapter("http://localhost", "key")
	if a.MediaType() != domain.MediaTypeMovie {
		t.Errorf("expected movie, got %s", a.MediaType())
	}
}

func TestRadarrAdapter_Search(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		encode(t, w, []map[string]any{
			{"tmdbId": 27205, "title": "Inception", "year": 2010, "overview": "A thief...", "added": ""},
		})
	})

	results, err := api.NewRadarrAdapter(srv.URL, "key").Search("inception")
	if err != nil {
		t.Fatalf("Search() error: %v", err)
	}
	if len(results) != 1 || results[0].Title != "Inception" || results[0].ID != 27205 {
		t.Errorf("unexpected results: %+v", results)
	}
}

func TestRadarrAdapter_GetLibrary(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		encode(t, w, []map[string]any{
			{"id": 1, "title": "Inception", "year": 2010, "hasFile": true},
		})
	})

	items, err := api.NewRadarrAdapter(srv.URL, "key").GetLibrary()
	if err != nil {
		t.Fatalf("GetLibrary() error: %v", err)
	}
	if len(items) != 1 || !items[0].HasFile || items[0].MediaType != domain.MediaTypeMovie {
		t.Errorf("unexpected items: %+v", items)
	}
}

func TestRadarrAdapter_GetQueue(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		encode(t, w, map[string]any{
			"records": []map[string]any{
				{"id": 1, "title": "Inception", "status": "downloading", "timeleft": "00:10:00"},
			},
		})
	})

	items, err := api.NewRadarrAdapter(srv.URL, "key").GetQueue()
	if err != nil {
		t.Fatalf("GetQueue() error: %v", err)
	}
	if len(items) != 1 || items[0].Status != "downloading" || items[0].TimeLeft != "00:10:00" {
		t.Errorf("unexpected items: %+v", items)
	}
}

// --- Sonarr ---

func TestSonarrAdapter_MediaType(t *testing.T) {
	a := api.NewSonarrAdapter("http://localhost", "key")
	if a.MediaType() != domain.MediaTypeSeries {
		t.Errorf("expected series, got %s", a.MediaType())
	}
}

func TestSonarrAdapter_Search(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		encode(t, w, []map[string]any{
			{"tvdbId": 81189, "title": "Breaking Bad", "year": 2008, "overview": "A chemistry teacher...", "added": ""},
		})
	})

	results, err := api.NewSonarrAdapter(srv.URL, "key").Search("breaking bad")
	if err != nil {
		t.Fatalf("Search() error: %v", err)
	}
	if len(results) != 1 || results[0].Title != "Breaking Bad" || results[0].ID != 81189 {
		t.Errorf("unexpected results: %+v", results)
	}
}

func TestSonarrAdapter_GetLibrary(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		encode(t, w, []map[string]any{
			{"id": 1, "title": "Breaking Bad", "year": 2008, "episodeFileCount": 62},
		})
	})

	items, err := api.NewSonarrAdapter(srv.URL, "key").GetLibrary()
	if err != nil {
		t.Fatalf("GetLibrary() error: %v", err)
	}
	if len(items) != 1 || !items[0].HasFile || items[0].MediaType != domain.MediaTypeSeries {
		t.Errorf("unexpected items: %+v", items)
	}
}

func TestSonarrAdapter_GetQueue(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		encode(t, w, map[string]any{
			"records": []map[string]any{
				{"id": 2, "title": "Breaking Bad S01E01", "status": "queued", "timeleft": "00:30:00"},
			},
		})
	})

	items, err := api.NewSonarrAdapter(srv.URL, "key").GetQueue()
	if err != nil {
		t.Fatalf("GetQueue() error: %v", err)
	}
	if len(items) != 1 || items[0].Status != "queued" || items[0].MediaType != domain.MediaTypeSeries {
		t.Errorf("unexpected items: %+v", items)
	}
}

// --- Lidarr ---

func TestLidarrAdapter_MediaType(t *testing.T) {
	a := api.NewLidarrAdapter("http://localhost", "key")
	if a.MediaType() != domain.MediaTypeArtist {
		t.Errorf("expected artist, got %s", a.MediaType())
	}
}

func TestLidarrAdapter_Search(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		encode(t, w, []map[string]any{
			{"id": 1, "artistName": "Radiohead", "overview": "British band", "added": ""},
		})
	})

	results, err := api.NewLidarrAdapter(srv.URL, "key").Search("radiohead")
	if err != nil {
		t.Fatalf("Search() error: %v", err)
	}
	if len(results) != 1 || results[0].Title != "Radiohead" {
		t.Errorf("unexpected results: %+v", results)
	}
}

func TestLidarrAdapter_GetLibrary(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		encode(t, w, []map[string]any{
			{"id": 1, "artistName": "Radiohead", "statistics": map[string]any{"trackFileCount": 10}},
		})
	})

	items, err := api.NewLidarrAdapter(srv.URL, "key").GetLibrary()
	if err != nil {
		t.Fatalf("GetLibrary() error: %v", err)
	}
	if len(items) != 1 || !items[0].HasFile || items[0].MediaType != domain.MediaTypeArtist {
		t.Errorf("unexpected items: %+v", items)
	}
}

func TestLidarrAdapter_GetQueue(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		encode(t, w, map[string]any{
			"records": []map[string]any{
				{"id": 3, "title": "OK Computer", "status": "downloading", "timeleft": "00:05:00"},
			},
		})
	})

	items, err := api.NewLidarrAdapter(srv.URL, "key").GetQueue()
	if err != nil {
		t.Fatalf("GetQueue() error: %v", err)
	}
	if len(items) != 1 || items[0].MediaType != domain.MediaTypeArtist {
		t.Errorf("unexpected items: %+v", items)
	}
}

// --- Readarr ---

func TestReadarrAdapter_MediaType(t *testing.T) {
	a := api.NewReadarrAdapter("http://localhost", "key")
	if a.MediaType() != domain.MediaTypeBook {
		t.Errorf("expected book, got %s", a.MediaType())
	}
}

func TestReadarrAdapter_Search(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		encode(t, w, []map[string]any{
			{
				"book":   map[string]any{"id": 1, "title": "Dune", "overview": "Epic sci-fi", "added": ""},
				"author": map[string]any{"authorName": "Frank Herbert"},
			},
		})
	})

	results, err := api.NewReadarrAdapter(srv.URL, "key").Search("dune")
	if err != nil {
		t.Fatalf("Search() error: %v", err)
	}
	if len(results) != 1 || results[0].Title != "Dune — Frank Herbert" {
		t.Errorf("unexpected results: %+v", results)
	}
}

func TestReadarrAdapter_GetLibrary(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		encode(t, w, []map[string]any{
			{"id": 1, "title": "Dune", "statistics": map[string]any{"bookFileCount": 1}},
		})
	})

	items, err := api.NewReadarrAdapter(srv.URL, "key").GetLibrary()
	if err != nil {
		t.Fatalf("GetLibrary() error: %v", err)
	}
	if len(items) != 1 || !items[0].HasFile || items[0].MediaType != domain.MediaTypeBook {
		t.Errorf("unexpected items: %+v", items)
	}
}

func TestReadarrAdapter_GetQueue(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		encode(t, w, map[string]any{
			"records": []map[string]any{
				{"id": 4, "title": "Dune", "status": "completed", "timeleft": ""},
			},
		})
	})

	items, err := api.NewReadarrAdapter(srv.URL, "key").GetQueue()
	if err != nil {
		t.Fatalf("GetQueue() error: %v", err)
	}
	if len(items) != 1 || items[0].Status != "completed" || items[0].MediaType != domain.MediaTypeBook {
		t.Errorf("unexpected items: %+v", items)
	}
}
