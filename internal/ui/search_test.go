package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/MarioFronza/media-tui/internal/domain"
	"github.com/MarioFronza/media-tui/internal/usecase"
)

func TestSearchModel_UpdateWithEnter_SetsLoading(t *testing.T) {
	m := NewSearchModel(nil)
	m.input.SetValue("inception")
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !updated.list.loading {
		t.Fatalf("expected loading to be true, got false")
	}
}

func TestSearchModel_UpdateWithEnter_SearchesForTerm(t *testing.T) {
	repo := &fakeRepo{searchFn: func(_ string) ([]domain.MediaItem, error) {
		return []domain.MediaItem{{ID: 1, Title: "Inception"}}, nil
	}}
	m := NewSearchModel(usecase.NewSearchUseCase(repo))
	m.input.SetValue("inception")

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	msg := cmd()
	batch, ok := msg.(tea.BatchMsg)
	if !ok {
		t.Fatalf("expected tea.BatchMsg, got %T", msg)
	}
	var result *searchResultMsg
	for _, c := range batch {
		if r, ok := c().(searchResultMsg); ok {
			result = &r
			break
		}
	}

	if result == nil {
		t.Fatal("no searchResultMsg found in batch")
	}
	if len(result.items) == 0 {
		t.Fatal("searchResultMsg contained no items")
	}
	if result.items[0].Title != "Inception" {
		t.Fatalf("expected title Inception, got %q", result.items[0].Title)
	}
}
