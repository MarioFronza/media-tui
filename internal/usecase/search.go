package usecase

import (
	"sync"

	"github.com/MarioFronza/media-tui/internal/domain"
)

type SearchUseCase struct {
	repos []domain.MediaRepository
}

func NewSearchUseCase(repos ...domain.MediaRepository) *SearchUseCase {
	return &SearchUseCase{repos: repos}
}

func (uc *SearchUseCase) Execute(term string) []domain.MediaItem {
	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		results []domain.MediaItem
	)

	for _, repo := range uc.repos {
		wg.Add(1)
		go func(r domain.MediaRepository) {
			defer wg.Done()
			items, err := r.Search(term)
			if err != nil {
				return
			}
			mu.Lock()
			results = append(results, items...)
			mu.Unlock()
		}(repo)
	}

	wg.Wait()
	return results
}
