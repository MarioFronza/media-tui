package usecase

import (
	"errors"
	"testing"

	"github.com/MarioFronza/media-tui/internal/domain"
)

func TestLibrary_List_DelegatesToRepo(t *testing.T) {
	want := []domain.LibraryItem{
		{ID: 1, Title: "The Matrix", MediaType: domain.MediaTypeMovie},
	}
	repo := &mockRepo{
		getLibraryFn: func() ([]domain.LibraryItem, error) {
			return want, nil
		},
	}

	uc := NewLibraryUseCase(repo)
	got, err := uc.List()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d items, got %d", len(want), len(got))
	}
	if got[0].Title != want[0].Title {
		t.Errorf("expected title %q, got %q", want[0].Title, got[0].Title)
	}
}

func TestLibrary_List_PropagatesError(t *testing.T) {
	repo := &mockRepo{
		getLibraryFn: func() ([]domain.LibraryItem, error) {
			return nil, errors.New("connection refused")
		},
	}

	uc := NewLibraryUseCase(repo)
	_, err := uc.List()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLibrary_Add_DelegatesToRepo(t *testing.T) {
	item := domain.MediaItem{ID: 42, Title: "Dune", Type: domain.MediaTypeMovie}
	var received domain.MediaItem

	repo := &mockRepo{
		addFn: func(i domain.MediaItem) error {
			received = i
			return nil
		},
	}

	uc := NewLibraryUseCase(repo)
	err := uc.Add(item)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.ID != item.ID || received.Title != item.Title {
		t.Errorf("repo received wrong item: %+v", received)
	}
}

func TestLibrary_Add_PropagatesError(t *testing.T) {
	repo := &mockRepo{
		addFn: func(item domain.MediaItem) error {
			return errors.New("already exists")
		},
	}

	uc := NewLibraryUseCase(repo)
	err := uc.Add(domain.MediaItem{ID: 1, Title: "Dune"})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
