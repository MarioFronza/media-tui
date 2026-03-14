package ui

import tea "github.com/charmbracelet/bubbletea"

type LibraryModel struct{}

func NewLibraryModel() LibraryModel { return LibraryModel{} }

func (m LibraryModel) Init() tea.Cmd { return nil }

func (m LibraryModel) Update(msg tea.Msg) (LibraryModel, tea.Cmd) { return m, nil }

func (m LibraryModel) View() string { return "Library screen — coming soon" }
