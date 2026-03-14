package api

import (
	"fmt"
	"net/url"

	"github.com/MarioFronza/media-tui/internal/domain"
)

type RadarrAdapter struct {
	client *Client
}

func NewRadarrAdapter(baseURL, apiKey string) *RadarrAdapter {
	return &RadarrAdapter{client: NewClient(baseURL, apiKey)}
}

func (a *RadarrAdapter) MediaType() domain.MediaType {
	return domain.MediaTypeMovie
}

func (a *RadarrAdapter) Search(term string) ([]domain.MediaItem, error) {
	var results []struct {
		TmdbID   int    `json:"tmdbId"`
		Title    string `json:"title"`
		Year     int    `json:"year"`
		Overview string `json:"overview"`
		Added    string `json:"added"`
	}
	if err := a.client.Get("/api/v3/movie/lookup?term="+url.QueryEscape(term), &results); err != nil {
		return nil, fmt.Errorf("radarr search: %w", err)
	}

	items := make([]domain.MediaItem, 0, len(results))
	for _, r := range results {
		items = append(items, domain.MediaItem{
			ID:       r.TmdbID,
			Title:    r.Title,
			Year:     r.Year,
			Overview: r.Overview,
			Type:     domain.MediaTypeMovie,
			Added:    r.Added != "" && r.Added != "0001-01-01T00:00:00Z",
		})
	}
	return items, nil
}

func (a *RadarrAdapter) Add(item domain.MediaItem) error {
	body := map[string]any{
		"tmdbId":            item.ID,
		"title":             item.Title,
		"qualityProfileId":  1,
		"rootFolderPath":    "/movies",
		"monitored":         true,
		"addOptions":        map[string]bool{"searchForMovie": true},
	}
	if err := a.client.Post("/api/v3/movie", body); err != nil {
		return fmt.Errorf("radarr add: %w", err)
	}
	return nil
}

func (a *RadarrAdapter) GetLibrary() ([]domain.LibraryItem, error) {
	var results []struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Year    int    `json:"year"`
		HasFile bool   `json:"hasFile"`
	}
	if err := a.client.Get("/api/v3/movie", &results); err != nil {
		return nil, fmt.Errorf("radarr library: %w", err)
	}

	items := make([]domain.LibraryItem, 0, len(results))
	for _, r := range results {
		items = append(items, domain.LibraryItem{
			ID:        r.ID,
			Title:     r.Title,
			Year:      r.Year,
			HasFile:   r.HasFile,
			MediaType: domain.MediaTypeMovie,
		})
	}
	return items, nil
}

func (a *RadarrAdapter) GetQueue() ([]domain.QueueItem, error) {
	var result struct {
		Records []struct {
			ID               int    `json:"id"`
			Title            string `json:"title"`
			Status           string `json:"status"`
			TimeLeft         string `json:"timeleft"`
		} `json:"records"`
	}
	if err := a.client.Get("/api/v3/queue", &result); err != nil {
		return nil, fmt.Errorf("radarr queue: %w", err)
	}

	items := make([]domain.QueueItem, 0, len(result.Records))
	for _, r := range result.Records {
		items = append(items, domain.QueueItem{
			ID:        r.ID,
			Title:     r.Title,
			Status:    r.Status,
			TimeLeft:  r.TimeLeft,
			MediaType: domain.MediaTypeMovie,
		})
	}
	return items, nil
}
