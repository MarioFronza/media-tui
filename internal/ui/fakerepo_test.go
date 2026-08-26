package ui

import "github.com/MarioFronza/media-tui/internal/domain"

type fakeRepo struct {
	searchFn     func(term string) ([]domain.MediaItem, error)
	getLibraryFn func() ([]domain.LibraryItem, error)
	getQueueFn   func() ([]domain.QueueItem, error)
}

func (f *fakeRepo) Search(term string) ([]domain.MediaItem, error) {
	if f.searchFn != nil {
		return f.searchFn(term)
	}
	return nil, nil
}

func (f *fakeRepo) Add(_ domain.MediaItem) error { return nil }

func (f *fakeRepo) GetLibrary() ([]domain.LibraryItem, error) {
	if f.getLibraryFn != nil {
		return f.getLibraryFn()
	}
	return nil, nil
}

func (f *fakeRepo) GetQueue() ([]domain.QueueItem, error) {
	if f.getQueueFn != nil {
		return f.getQueueFn()
	}
	return nil, nil
}

func (f *fakeRepo) MediaType() domain.MediaType { return "" }
