package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/MarioFronza/media-tui/internal/domain"
	"github.com/MarioFronza/media-tui/internal/usecase"
)

type addResultMsg struct{ err error }

var (
	labelStyle   = lipgloss.NewStyle().Bold(true)
	dividerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	hintStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	addingStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
)

type DetailModel struct {
	item    domain.MediaItem
	usecase *usecase.LibraryUseCase
	status  string // "", "adding", "added", "error: ..."
	width   int
}

func NewDetailModel() DetailModel { return DetailModel{} }

func NewDetailModelWithItem(item domain.MediaItem, uc *usecase.LibraryUseCase) DetailModel {
	return DetailModel{item: item, usecase: uc}
}

func (m *DetailModel) SetSize(w, h int) {
	m.width = w
}

func (m DetailModel) Init() tea.Cmd { return nil }

func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg {
				return SwitchScreenMsg{Target: screenSearch}
			}
		case "a", "enter":
			if m.usecase == nil || m.status == "adding" || m.status == "added" || m.item.Added {
				return m, nil
			}
			m.status = "adding"
			item := m.item
			uc := m.usecase
			return m, func() tea.Msg {
				err := uc.Add(item)
				return addResultMsg{err: err}
			}
		}

	case addResultMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("error: %s", msg.err.Error())
		} else {
			m.status = "added"
		}
	}

	return m, nil
}

func (m DetailModel) View() string {
	w := m.width
	if w <= 0 {
		w = 80
	}
	divider := dividerStyle.Render(lipgloss.NewStyle().Width(w - 2).Render(""))
	overviewWidth := w - 12 // "Overview: " label width
	if overviewWidth < 20 {
		overviewWidth = 20
	}

	title := fmt.Sprintf("%s %s", labelStyle.Render("Title:"), m.item.Title)
	year := fmt.Sprintf("%s %d", labelStyle.Render("Year:"), m.item.Year)
	mediaType := fmt.Sprintf("%s %s", labelStyle.Render("Type:"), string(m.item.Type))
	overview := fmt.Sprintf("%s %s", labelStyle.Render("Overview:"),
		lipgloss.NewStyle().Width(overviewWidth).Render(m.item.Overview))

	var statusLine string
	switch {
	case m.item.Added || m.status == "added":
		if m.status == "added" {
			statusLine = successStyle.Render("Added successfully!")
		} else {
			statusLine = successStyle.Render("Already in library")
		}
	case m.status == "adding":
		statusLine = addingStyle.Render("Adding...")
	case len(m.status) > 6 && m.status[:6] == "error:":
		statusLine = errorStyle.Render(m.status)
	default:
		statusLine = hintStyle.Render("a/Enter add to library · Esc back")
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("Detail"),
		divider,
		title,
		year,
		mediaType,
		overview,
		"",
		statusLine,
	)
}
