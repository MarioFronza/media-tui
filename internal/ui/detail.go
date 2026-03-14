package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/MarioFronza/media-tui/internal/domain"
)

type DetailModel struct {
	item domain.MediaItem
}

func NewDetailModel() DetailModel { return DetailModel{} }

func NewDetailModelWithItem(item domain.MediaItem) DetailModel { return DetailModel{item: item} }

func (m DetailModel) Init() tea.Cmd { return nil }

func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) { return m, nil }

func (m DetailModel) View() string { return "Detail screen — coming soon" }
