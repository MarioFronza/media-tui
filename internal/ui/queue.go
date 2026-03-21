package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/MarioFronza/media-tui/internal/domain"
	"github.com/MarioFronza/media-tui/internal/usecase"
)

type queueResultMsg struct {
	items []domain.QueueItem
}

type QueueModel struct {
	usecase *usecase.QueueUseCase
	spinner spinner.Model
	table   table.Model
	items   []domain.QueueItem
	loading bool
	width   int
	height  int
}

func NewQueueModel(uc *usecase.QueueUseCase) QueueModel {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	cols := []table.Column{
		{Title: "Title", Width: 40},
		{Title: "Service", Width: 10},
		{Title: "Status", Width: 12},
		{Title: "Time Left", Width: 12},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(15),
	)
	t.SetStyles(tableStyles())

	return QueueModel{
		usecase: uc,
		spinner: sp,
		table:   t,
	}
}

func (m *QueueModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	titleW := max(10, w-42)
	m.table.SetColumns([]table.Column{
		{Title: "Title", Width: titleW},
		{Title: "Service", Width: 10},
		{Title: "Status", Width: 12},
		{Title: "Time Left", Width: 12},
	})
	m.table.SetHeight(max(5, h-5))
}

func (m QueueModel) Init() tea.Cmd {
	_, cmd := m.fetchQueue()
	return cmd
}

func (m QueueModel) Update(msg tea.Msg) (QueueModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case queueResultMsg:
		return m.handleResult(msg)
	case spinner.TickMsg:
		return m.updateSpinner(msg)
	}
	return m.updateTable(msg)
}

func (m QueueModel) handleKey(msg tea.KeyMsg) (QueueModel, tea.Cmd) {
	switch msg.String() {
	case "r":
		return m.fetchQueue()
	}
	return m.updateTable(msg)
}

func (m QueueModel) handleResult(msg queueResultMsg) (QueueModel, tea.Cmd) {
	m.loading = false
	m.items = msg.items
	m.table.SetRows(toQueueRows(msg.items))
	return m, nil
}

func (m QueueModel) updateSpinner(msg spinner.TickMsg) (QueueModel, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m QueueModel) updateTable(msg tea.Msg) (QueueModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m QueueModel) fetchQueue() (QueueModel, tea.Cmd) {
	m.loading = true
	return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
		items, _ := m.usecase.Execute()
		return queueResultMsg{items: items}
	})
}

func (m QueueModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Render("Queue") + "\n\n"

	if m.loading {
		return header + m.spinner.View() + " Loading..."
	}

	if len(m.items) == 0 {
		return header + "No items in queue."
	}

	hint := lipgloss.NewStyle().Faint(true).Render("r refresh · ↑/↓ navigate")
	return header + m.table.View() + "\n" + hint
}

func toQueueRows(items []domain.QueueItem) []table.Row {
	rows := make([]table.Row, 0, len(items))
	for _, item := range items {
		rows = append(rows, table.Row{item.Title, string(item.MediaType), item.Status, item.TimeLeft})
	}
	return rows
}
