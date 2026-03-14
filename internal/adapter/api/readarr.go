package api

import (
	"fmt"
	"net/url"

	"github.com/MarioFronza/media-tui/internal/domain"
)

type ReadarrAdapter struct {
	client *Client
}

func NewReadarrAdapter(baseURL, apiKey string) *ReadarrAdapter {
	return &ReadarrAdapter{client: NewClient(baseURL, apiKey)}
}

func (a *ReadarrAdapter) MediaType() domain.MediaType {
	return domain.MediaTypeBook
}

func (a *ReadarrAdapter) Search(term string) ([]domain.MediaItem, error) {
	var results []struct {
		Book struct {
			ID       int    `json:"id"`
			Title    string `json:"title"`
			Overview string `json:"overview"`
			Added    string `json:"added"`
		} `json:"book"`
		Author struct {
			AuthorName string `json:"authorName"`
		} `json:"author"`
	}
	if err := a.client.Get("/api/v1/book/lookup?term="+url.QueryEscape(term), &results); err != nil {
		return nil, fmt.Errorf("readarr search: %w", err)
	}

	items := make([]domain.MediaItem, 0, len(results))
	for _, r := range results {
		title := r.Book.Title
		if r.Author.AuthorName != "" {
			title = r.Book.Title + " — " + r.Author.AuthorName
		}
		items = append(items, domain.MediaItem{
			ID:       r.Book.ID,
			Title:    title,
			Overview: r.Book.Overview,
			Type:     domain.MediaTypeBook,
			Added:    r.Book.Added != "" && r.Book.Added != "0001-01-01T00:00:00Z",
		})
	}
	return items, nil
}

func (a *ReadarrAdapter) Add(item domain.MediaItem) error {
	body := map[string]any{
		"title":            item.Title,
		"qualityProfileId": 1,
		"rootFolderPath":   "/books",
		"monitored":        true,
		"addOptions":       map[string]bool{"searchForNewBook": true},
	}
	if err := a.client.Post("/api/v1/book", body); err != nil {
		return fmt.Errorf("readarr add: %w", err)
	}
	return nil
}

func (a *ReadarrAdapter) GetLibrary() ([]domain.LibraryItem, error) {
	var results []struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Statistics struct {
			BookFileCount int `json:"bookFileCount"`
		} `json:"statistics"`
	}
	if err := a.client.Get("/api/v1/book", &results); err != nil {
		return nil, fmt.Errorf("readarr library: %w", err)
	}

	items := make([]domain.LibraryItem, 0, len(results))
	for _, r := range results {
		items = append(items, domain.LibraryItem{
			ID:        r.ID,
			Title:     r.Title,
			HasFile:   r.Statistics.BookFileCount > 0,
			MediaType: domain.MediaTypeBook,
		})
	}
	return items, nil
}

func (a *ReadarrAdapter) GetQueue() ([]domain.QueueItem, error) {
	var result struct {
		Records []struct {
			ID       int    `json:"id"`
			Title    string `json:"title"`
			Status   string `json:"status"`
			TimeLeft string `json:"timeleft"`
		} `json:"records"`
	}
	if err := a.client.Get("/api/v1/queue", &result); err != nil {
		return nil, fmt.Errorf("readarr queue: %w", err)
	}

	items := make([]domain.QueueItem, 0, len(result.Records))
	for _, r := range result.Records {
		items = append(items, domain.QueueItem{
			ID:        r.ID,
			Title:     r.Title,
			Status:    r.Status,
			TimeLeft:  r.TimeLeft,
			MediaType: domain.MediaTypeBook,
		})
	}
	return items, nil
}
