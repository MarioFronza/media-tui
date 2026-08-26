package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/MarioFronza/media-tui/internal/domain"
	"github.com/MarioFronza/media-tui/internal/usecase"
)

func TestLibraryModel_UpdateWithTab_SetsLoading(t *testing.T) {
	uc := usecase.NewLibraryUseCase(nil)
	m := NewLibraryModel(uc, uc)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})

	if !updated.list.loading {
		t.Fatalf("expected loading to be true, got false")
	}
}

func TestLibraryModel_UpdateWithTab_LoadsActiveServiceItems(t *testing.T) {
	repo := &fakeRepo{getLibraryFn: func() ([]domain.LibraryItem, error) {
		return []domain.LibraryItem{{ID: 1, Title: "Dune"}}, nil
	}}
	uc := usecase.NewLibraryUseCase(repo)
	m := NewLibraryModel(uc, uc)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyTab})

	msg := cmd()
	batch, ok := msg.(tea.BatchMsg)
	if !ok {
		t.Fatalf("expected tea.BatchMsg, got %T", msg)
	}
	var result *libraryResultMsg
	for _, c := range batch {
		if r, ok := c().(libraryResultMsg); ok {
			result = &r
			break
		}
	}

	if result == nil {
		t.Fatal("no libraryResultMsg found in batch")
	}
	if len(result.items) == 0 {
		t.Fatal("libraryResultMsg contained no items")
	}
	if result.items[0].Title != "Dune" {
		t.Fatalf("expected title Dune, got %q", result.items[0].Title)
	}
}
