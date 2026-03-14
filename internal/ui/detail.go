package ui

import tea "github.com/charmbracelet/bubbletea"

type DetailModel struct{}

func NewDetailModel() DetailModel { return DetailModel{} }

func (m DetailModel) Init() tea.Cmd { return nil }

func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) { return m, nil }

func (m DetailModel) View() string { return "Detail screen — coming soon" }
