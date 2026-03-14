package usecase

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/MarioFronza/media-tui/internal/domain"
)

func TestSearch_SingleRepo_ReturnsResults(t *testing.T) {
	want := []domain.MediaItem{
		{ID: 1, Title: "Inception", Type: domain.MediaTypeMovie},
	}
	repo := &mockRepo{
		searchFn: func(term string) ([]domain.MediaItem, error) {
			return want, nil
		},
	}

	uc := NewSearchUseCase(repo)
	got := uc.Execute("inception")

	if len(got) != len(want) {
		t.Fatalf("expected %d results, got %d", len(want), len(got))
	}
	if got[0].Title != want[0].Title {
		t.Errorf("expected title %q, got %q", want[0].Title, got[0].Title)
	}
}

func TestSearch_MultipleRepos_AggregatesResults(t *testing.T) {
	repo1 := &mockRepo{
		searchFn: func(term string) ([]domain.MediaItem, error) {
			return []domain.MediaItem{{ID: 1, Title: "Movie A", Type: domain.MediaTypeMovie}}, nil
		},
	}
	repo2 := &mockRepo{
		searchFn: func(term string) ([]domain.MediaItem, error) {
			return []domain.MediaItem{{ID: 2, Title: "Series B", Type: domain.MediaTypeSeries}}, nil
		},
	}

	uc := NewSearchUseCase(repo1, repo2)
	got := uc.Execute("test")

	if len(got) != 2 {
		t.Fatalf("expected 2 aggregated results, got %d", len(got))
	}
}

func TestSearch_OneRepoErrors_SkipsIt(t *testing.T) {
	errRepo := &mockRepo{
		searchFn: func(term string) ([]domain.MediaItem, error) {
			return nil, errors.New("timeout")
		},
	}
	okRepo := &mockRepo{
		searchFn: func(term string) ([]domain.MediaItem, error) {
			return []domain.MediaItem{{ID: 1, Title: "Good Result"}}, nil
		},
	}

	uc := NewSearchUseCase(errRepo, okRepo)
	got := uc.Execute("test")

	if len(got) != 1 {
		t.Fatalf("expected 1 result from healthy repo, got %d", len(got))
	}
	if got[0].Title != "Good Result" {
		t.Errorf("unexpected title: %q", got[0].Title)
	}
}

func TestSearch_NoRepos_ReturnsEmpty(t *testing.T) {
	uc := NewSearchUseCase()
	got := uc.Execute("anything")

	if len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}
}

func TestSearch_Concurrent_AllReposQueried(t *testing.T) {
	var callCount int64

	makeRepo := func() *mockRepo {
		return &mockRepo{
			searchFn: func(term string) ([]domain.MediaItem, error) {
				atomic.AddInt64(&callCount, 1)
				return nil, nil
			},
		}
	}

	repos := []domain.MediaRepository{makeRepo(), makeRepo(), makeRepo()}
	uc := NewSearchUseCase(repos...)
	uc.Execute("test")

	if callCount != 3 {
		t.Errorf("expected 3 repos queried, got %d", callCount)
	}
}
