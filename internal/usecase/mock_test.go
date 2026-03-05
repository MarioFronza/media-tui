package usecase

import "github.com/MarioFronza/media-tui/internal/domain"

type mockRepo struct {
	searchFn     func(term string) ([]domain.MediaItem, error)
	addFn        func(item domain.MediaItem) error
	getLibraryFn func() ([]domain.LibraryItem, error)
	getQueueFn   func() ([]domain.QueueItem, error)
	mediaType    domain.MediaType
}

func (m *mockRepo) Search(term string) ([]domain.MediaItem, error) {
	if m.searchFn != nil {
		return m.searchFn(term)
	}
	return nil, nil
}

func (m *mockRepo) Add(item domain.MediaItem) error {
	if m.addFn != nil {
		return m.addFn(item)
	}
	return nil
}

func (m *mockRepo) GetLibrary() ([]domain.LibraryItem, error) {
	if m.getLibraryFn != nil {
		return m.getLibraryFn()
	}
	return nil, nil
}

func (m *mockRepo) GetQueue() ([]domain.QueueItem, error) {
	if m.getQueueFn != nil {
		return m.getQueueFn()
	}
	return nil, nil
}

func (m *mockRepo) MediaType() domain.MediaType {
	return m.mediaType
}
