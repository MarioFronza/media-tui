package usecase

import "github.com/MarioFronza/media-tui/internal/domain"

type QueueUseCase struct {
	repos []domain.MediaRepository
}

func NewQueueUseCase(repos ...domain.MediaRepository) *QueueUseCase {
	return &QueueUseCase{repos: repos}
}

func (uc *QueueUseCase) Execute() ([]domain.QueueItem, error) {
	var all []domain.QueueItem
	for _, repo := range uc.repos {
		items, err := repo.GetQueue()
		if err != nil {
			continue
		}
		all = append(all, items...)
	}
	return all, nil
}
