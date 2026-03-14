package domain

type MediaType string

const (
	MediaTypeMovie  MediaType = "movie"
	MediaTypeSeries MediaType = "series"
	MediaTypeArtist MediaType = "artist"
	MediaTypeBook   MediaType = "book"
)

type MediaItem struct {
	ID       int
	Title    string
	Year     int
	Overview string
	Type     MediaType
	Added    bool
}

type QueueItem struct {
	ID        int
	Title     string
	Status    string
	TimeLeft  string
	MediaType MediaType
}

type LibraryItem struct {
	ID        int
	Title     string
	Year      int
	HasFile   bool
	MediaType MediaType
}
