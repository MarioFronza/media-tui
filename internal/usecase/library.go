package usecase

import "github.com/MarioFronza/media-tui/internal/domain"

type LibraryUseCase struct {
	repo domain.MediaRepository
}

func NewLibraryUseCase(repo domain.MediaRepository) *LibraryUseCase {
	return &LibraryUseCase{repo: repo}
}

func (uc *LibraryUseCase) List() ([]domain.LibraryItem, error) {
	return uc.repo.GetLibrary()
}

func (uc *LibraryUseCase) Add(item domain.MediaItem) error {
	return uc.repo.Add(item)
}
