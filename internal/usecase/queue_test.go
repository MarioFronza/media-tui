package usecase

import (
	"errors"
	"testing"

	"github.com/MarioFronza/media-tui/internal/domain"
)

func TestQueue_MultipleRepos_AggregatesItems(t *testing.T) {
	repo1 := &mockRepo{
		getQueueFn: func() ([]domain.QueueItem, error) {
			return []domain.QueueItem{{ID: 1, Title: "Movie Download", MediaType: domain.MediaTypeMovie}}, nil
		},
	}
	repo2 := &mockRepo{
		getQueueFn: func() ([]domain.QueueItem, error) {
			return []domain.QueueItem{{ID: 2, Title: "Series Download", MediaType: domain.MediaTypeSeries}}, nil
		},
	}

	uc := NewQueueUseCase(repo1, repo2)
	got, err := uc.Execute()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 queue items, got %d", len(got))
	}
}

func TestQueue_OneRepoErrors_ContinuesOthers(t *testing.T) {
	errRepo := &mockRepo{
		getQueueFn: func() ([]domain.QueueItem, error) {
			return nil, errors.New("service unavailable")
		},
	}
	okRepo := &mockRepo{
		getQueueFn: func() ([]domain.QueueItem, error) {
			return []domain.QueueItem{{ID: 1, Title: "Active Download"}}, nil
		},
	}

	uc := NewQueueUseCase(errRepo, okRepo)
	got, err := uc.Execute()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 item from healthy repo, got %d", len(got))
	}
}

func TestQueue_NoRepos_ReturnsEmptyNilError(t *testing.T) {
	uc := NewQueueUseCase()
	got, err := uc.Execute()

	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}
}
