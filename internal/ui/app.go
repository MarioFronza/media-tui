package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/MarioFronza/media-tui/internal/domain"
	"github.com/MarioFronza/media-tui/internal/usecase"
)

type screen int

const (
	screenSearch screen = iota
	screenLibrary
	screenQueue
	screenDetail
)

// SwitchScreenMsg is sent by sub-models to request a screen transition.
type SwitchScreenMsg struct {
	Target screen
}

// SwitchToDetailMsg carries the selected item to the detail screen.
type SwitchToDetailMsg struct {
	Item domain.MediaItem
}

type App struct {
	current    screen
	search     SearchModel
	library    LibraryModel
	queue      QueueModel
	detail     DetailModel
	libraryUCs []*usecase.LibraryUseCase
	width      int
	height     int
}

func NewApp(searchUC *usecase.SearchUseCase, queueUC *usecase.QueueUseCase, libraryUCs ...*usecase.LibraryUseCase) App {
	return App{
		current:    screenSearch,
		search:     NewSearchModel(searchUC),
		library:    NewLibraryModel(libraryUCs...),
		queue:      NewQueueModel(queueUC),
		detail:     NewDetailModel(),
		libraryUCs: libraryUCs,
	}
}

func (a App) Init() tea.Cmd {
	return a.search.Init()
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		case "/":
			a.current = screenSearch
			return a, a.search.Init()
		case "l":
			if a.current != screenDetail {
				a.current = screenLibrary
				return a, a.library.Init()
			}
		case "q":
			if a.current != screenDetail {
				a.current = screenQueue
				return a, a.queue.Init()
			}
		}

	case SwitchScreenMsg:
		a.current = msg.Target
		switch msg.Target {
		case screenSearch:
			return a, a.search.Init()
		case screenLibrary:
			return a, a.library.Init()
		case screenQueue:
			return a, a.queue.Init()
		case screenDetail:
			return a, a.detail.Init()
		}

	case SwitchToDetailMsg:
		a.detail = NewDetailModelWithItem(msg.Item, a.ucForItem(msg.Item))
		a.current = screenDetail
		return a, a.detail.Init()
	}

	return a.updateActive(msg)
}

func (a App) updateActive(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch a.current {
	case screenSearch:
		a.search, cmd = a.search.Update(msg)
	case screenLibrary:
		a.library, cmd = a.library.Update(msg)
	case screenQueue:
		a.queue, cmd = a.queue.Update(msg)
	case screenDetail:
		a.detail, cmd = a.detail.Update(msg)
	}
	return a, cmd
}

var tabStyle = lipgloss.NewStyle().Padding(0, 2)
var activeTabStyle = tabStyle.Bold(true).Foreground(lipgloss.Color("205"))

func (a App) View() string {
	tabs := lipgloss.JoinHorizontal(lipgloss.Top,
		a.tab("Search (/)", screenSearch),
		a.tab("Library (l)", screenLibrary),
		a.tab("Queue (q)", screenQueue),
	)

	var content string
	switch a.current {
	case screenSearch:
		content = a.search.View()
	case screenLibrary:
		content = a.library.View()
	case screenQueue:
		content = a.queue.View()
	case screenDetail:
		content = a.detail.View()
	}

	return lipgloss.JoinVertical(lipgloss.Left, tabs, content)
}

func (a App) tab(label string, s screen) string {
	if a.current == s {
		return activeTabStyle.Render(label)
	}
	return tabStyle.Render(label)
}

// ucForItem returns the LibraryUseCase that matches the item's media type.
// The libraryUCs slice is ordered: [0] Radarr (movie), [1] Sonarr (series),
// [2] Lidarr (artist), [3] Readarr (book). Returns nil if none available.
func (a App) ucForItem(item domain.MediaItem) *usecase.LibraryUseCase {
	idx := map[domain.MediaType]int{
		domain.MediaTypeMovie:  0,
		domain.MediaTypeSeries: 1,
		domain.MediaTypeArtist: 2,
		domain.MediaTypeBook:   3,
	}
	if i, ok := idx[item.Type]; ok && i < len(a.libraryUCs) {
		return a.libraryUCs[i]
	}
	return nil
}
