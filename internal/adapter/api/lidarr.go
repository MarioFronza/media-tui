package api

import (
	"fmt"
	"net/url"

	"github.com/MarioFronza/media-tui/internal/domain"
)

type LidarrAdapter struct {
	client *Client
}

func NewLidarrAdapter(baseURL, apiKey string) *LidarrAdapter {
	return &LidarrAdapter{client: NewClient(baseURL, apiKey)}
}

func (a *LidarrAdapter) MediaType() domain.MediaType {
	return domain.MediaTypeArtist
}

func (a *LidarrAdapter) Search(term string) ([]domain.MediaItem, error) {
	var results []struct {
		ID       int    `json:"id"`
		ArtistName string `json:"artistName"`
		Overview string `json:"overview"`
		Added    string `json:"added"`
	}
	if err := a.client.Get("/api/v1/artist/lookup?term="+url.QueryEscape(term), &results); err != nil {
		return nil, fmt.Errorf("lidarr search: %w", err)
	}

	items := make([]domain.MediaItem, 0, len(results))
	for _, r := range results {
		items = append(items, domain.MediaItem{
			ID:       r.ID,
			Title:    r.ArtistName,
			Overview: r.Overview,
			Type:     domain.MediaTypeArtist,
			Added:    r.Added != "" && r.Added != "0001-01-01T00:00:00Z",
		})
	}
	return items, nil
}

func (a *LidarrAdapter) Add(item domain.MediaItem) error {
	body := map[string]any{
		"artistName":       item.Title,
		"qualityProfileId": 1,
		"metadataProfileId": 1,
		"rootFolderPath":   "/music",
		"monitored":        true,
		"addOptions":       map[string]bool{"searchForMissingAlbums": true},
	}
	if err := a.client.Post("/api/v1/artist", body); err != nil {
		return fmt.Errorf("lidarr add: %w", err)
	}
	return nil
}

func (a *LidarrAdapter) GetLibrary() ([]domain.LibraryItem, error) {
	var results []struct {
		ID         int    `json:"id"`
		ArtistName string `json:"artistName"`
		Statistics struct {
			TrackFileCount int `json:"trackFileCount"`
		} `json:"statistics"`
	}
	if err := a.client.Get("/api/v1/artist", &results); err != nil {
		return nil, fmt.Errorf("lidarr library: %w", err)
	}

	items := make([]domain.LibraryItem, 0, len(results))
	for _, r := range results {
		items = append(items, domain.LibraryItem{
			ID:        r.ID,
			Title:     r.ArtistName,
			HasFile:   r.Statistics.TrackFileCount > 0,
			MediaType: domain.MediaTypeArtist,
		})
	}
	return items, nil
}

func (a *LidarrAdapter) GetQueue() ([]domain.QueueItem, error) {
	var result struct {
		Records []struct {
			ID       int    `json:"id"`
			Title    string `json:"title"`
			Status   string `json:"status"`
			TimeLeft string `json:"timeleft"`
		} `json:"records"`
	}
	if err := a.client.Get("/api/v1/queue", &result); err != nil {
		return nil, fmt.Errorf("lidarr queue: %w", err)
	}

	items := make([]domain.QueueItem, 0, len(result.Records))
	for _, r := range result.Records {
		items = append(items, domain.QueueItem{
			ID:        r.ID,
			Title:     r.Title,
			Status:    r.Status,
			TimeLeft:  r.TimeLeft,
			MediaType: domain.MediaTypeArtist,
		})
	}
	return items, nil
}
