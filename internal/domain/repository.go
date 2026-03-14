package domain

type MediaRepository interface {
	Search(term string) ([]MediaItem, error)
	Add(item MediaItem) error
	GetLibrary() ([]LibraryItem, error)
	GetQueue() ([]QueueItem, error)
	MediaType() MediaType
}
