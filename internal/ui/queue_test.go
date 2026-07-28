package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/MarioFronza/media-tui/internal/domain"
	"github.com/MarioFronza/media-tui/internal/usecase"
)

type fakeQueueRepo struct {
	items []domain.QueueItem
}

func (f *fakeQueueRepo) Search(_ string) ([]domain.MediaItem, error)  { return nil, nil }
func (f *fakeQueueRepo) Add(_ domain.MediaItem) error                 { return nil }
func (f *fakeQueueRepo) GetLibrary() ([]domain.LibraryItem, error)    { return nil, nil }
func (f *fakeQueueRepo) GetQueue() ([]domain.QueueItem, error)        { return f.items, nil }
func (f *fakeQueueRepo) MediaType() domain.MediaType                  { return "" }

func TestQueueModel_UpdateWithR_SetsLoading(t *testing.T) {
	m := NewQueueModel(nil)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})

	if !updated.loading {
		t.Fatalf("expected loading to be true, got false")
	}
}

func TestQueueModel_UpdateWithR_RefreshesQueueItems(t *testing.T) {
	repo := &fakeQueueRepo{items: []domain.QueueItem{{ID: 1, Title: "Inception"}}}
	uc := usecase.NewQueueUseCase(repo)
	m := NewQueueModel(uc)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})

	msg := cmd()
	batch, ok := msg.(tea.BatchMsg)
	if !ok {
		t.Fatalf("expected tea.BatchMsg, got %T", msg)
	}
	var result *queueResultMsg
	for _, c := range batch {
		if r, ok := c().(queueResultMsg); ok {
			result = &r
			break
		}
	}

	if result == nil {
		t.Fatal("no queueResultMsg found in batch")
	}
	if len(result.items) == 0 {
		t.Fatal("queueResultMsg contained no items")
	}
	if result.items[0].Title != "Inception" {
		t.Fatalf("expected title Inception, got %q", result.items[0].Title)
	}
}
