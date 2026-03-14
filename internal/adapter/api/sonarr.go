package api

import (
	"fmt"
	"net/url"

	"github.com/MarioFronza/media-tui/internal/domain"
)

type SonarrAdapter struct {
	client *Client
}

func NewSonarrAdapter(baseURL, apiKey string) *SonarrAdapter {
	return &SonarrAdapter{client: NewClient(baseURL, apiKey)}
}

func (a *SonarrAdapter) MediaType() domain.MediaType {
	return domain.MediaTypeSeries
}

func (a *SonarrAdapter) Search(term string) ([]domain.MediaItem, error) {
	var results []struct {
		TvdbID   int    `json:"tvdbId"`
		Title    string `json:"title"`
		Year     int    `json:"year"`
		Overview string `json:"overview"`
		Added    string `json:"added"`
	}
	if err := a.client.Get("/api/v3/series/lookup?term="+url.QueryEscape(term), &results); err != nil {
		return nil, fmt.Errorf("sonarr search: %w", err)
	}

	items := make([]domain.MediaItem, 0, len(results))
	for _, r := range results {
		items = append(items, domain.MediaItem{
			ID:       r.TvdbID,
			Title:    r.Title,
			Year:     r.Year,
			Overview: r.Overview,
			Type:     domain.MediaTypeSeries,
			Added:    r.Added != "" && r.Added != "0001-01-01T00:00:00Z",
		})
	}
	return items, nil
}

func (a *SonarrAdapter) Add(item domain.MediaItem) error {
	body := map[string]any{
		"tvdbId":           item.ID,
		"title":            item.Title,
		"qualityProfileId": 1,
		"rootFolderPath":   "/tv",
		"monitored":        true,
		"addOptions":       map[string]bool{"searchForMissingEpisodes": true},
	}
	if err := a.client.Post("/api/v3/series", body); err != nil {
		return fmt.Errorf("sonarr add: %w", err)
	}
	return nil
}

func (a *SonarrAdapter) GetLibrary() ([]domain.LibraryItem, error) {
	var results []struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Year        int    `json:"year"`
		EpisodeFileCount int `json:"episodeFileCount"`
	}
	if err := a.client.Get("/api/v3/series", &results); err != nil {
		return nil, fmt.Errorf("sonarr library: %w", err)
	}

	items := make([]domain.LibraryItem, 0, len(results))
	for _, r := range results {
		items = append(items, domain.LibraryItem{
			ID:        r.ID,
			Title:     r.Title,
			Year:      r.Year,
			HasFile:   r.EpisodeFileCount > 0,
			MediaType: domain.MediaTypeSeries,
		})
	}
	return items, nil
}

func (a *SonarrAdapter) GetQueue() ([]domain.QueueItem, error) {
	var result struct {
		Records []struct {
			ID       int    `json:"id"`
			Title    string `json:"title"`
			Status   string `json:"status"`
			TimeLeft string `json:"timeleft"`
		} `json:"records"`
	}
	if err := a.client.Get("/api/v3/queue", &result); err != nil {
		return nil, fmt.Errorf("sonarr queue: %w", err)
	}

	items := make([]domain.QueueItem, 0, len(result.Records))
	for _, r := range result.Records {
		items = append(items, domain.QueueItem{
			ID:        r.ID,
			Title:     r.Title,
			Status:    r.Status,
			TimeLeft:  r.TimeLeft,
			MediaType: domain.MediaTypeSeries,
		})
	}
	return items, nil
}
