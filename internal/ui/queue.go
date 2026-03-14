package ui

import tea "github.com/charmbracelet/bubbletea"

type QueueModel struct{}

func NewQueueModel() QueueModel { return QueueModel{} }

func (m QueueModel) Init() tea.Cmd { return nil }

func (m QueueModel) Update(msg tea.Msg) (QueueModel, tea.Cmd) { return m, nil }

func (m QueueModel) View() string { return "Queue screen — coming soon" }
