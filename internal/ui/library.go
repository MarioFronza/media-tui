package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/MarioFronza/media-tui/internal/domain"
	"github.com/MarioFronza/media-tui/internal/usecase"
)

type libraryResultMsg struct {
	service string
	items   []domain.LibraryItem
}

type serviceTab struct {
	name    string
	usecase *usecase.LibraryUseCase
}

type LibraryModel struct {
	services    []serviceTab
	activeIdx   int
	spinner     spinner.Model
	table       table.Model
	items       []domain.LibraryItem
	loading     bool
}

func NewLibraryModel(usecases ...*usecase.LibraryUseCase) LibraryModel {
	services := []serviceTab{
		{name: "Radarr"},
		{name: "Sonarr"},
		{name: "Lidarr"},
		{name: "Readarr"},
	}
	for i, uc := range usecases {
		if i < len(services) {
			services[i].usecase = uc
		}
	}

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	cols := []table.Column{
		{Title: "Title", Width: 40},
		{Title: "Year", Width: 6},
		{Title: "Has File", Width: 10},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(15),
	)
	t.SetStyles(tableStyles())

	return LibraryModel{
		services: services,
		spinner:  sp,
		table:    t,
	}
}

func (m LibraryModel) Init() tea.Cmd {
	_, cmd := m.fetchLibrary()
	return cmd
}

func (m LibraryModel) Update(msg tea.Msg) (LibraryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case libraryResultMsg:
		return m.handleResult(msg)
	case spinner.TickMsg:
		return m.updateSpinner(msg)
	}
	return m.updateTable(msg)
}

func (m LibraryModel) handleKey(msg tea.KeyMsg) (LibraryModel, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.activeIdx = (m.activeIdx + 1) % len(m.services)
		return m.fetchLibrary()
	case "shift+tab":
		m.activeIdx = (m.activeIdx - 1 + len(m.services)) % len(m.services)
		return m.fetchLibrary()
	}
	return m.updateTable(msg)
}

func (m LibraryModel) handleResult(msg libraryResultMsg) (LibraryModel, tea.Cmd) {
	m.loading = false
	m.items = msg.items
	m.table.SetRows(toLibraryRows(msg.items))
	return m, nil
}

func (m LibraryModel) updateSpinner(msg spinner.TickMsg) (LibraryModel, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m LibraryModel) updateTable(msg tea.Msg) (LibraryModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m LibraryModel) fetchLibrary() (LibraryModel, tea.Cmd) {
	svc := m.services[m.activeIdx]
	if svc.usecase == nil {
		return m, func() tea.Msg {
			return libraryResultMsg{service: svc.name, items: nil}
		}
	}
	m.loading = true
	return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
		items, _ := svc.usecase.List()
		return libraryResultMsg{service: svc.name, items: items}
	})
}

func (m LibraryModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Render("Library") + "\n\n"
	tabs := m.renderTabs() + "\n\n"

	if m.loading {
		return header + tabs + m.spinner.View() + " Loading..."
	}

	if len(m.items) == 0 {
		return header + tabs + "No items in library."
	}

	hint := lipgloss.NewStyle().Faint(true).Render("↑/↓ navigate · Tab next service · Shift+Tab prev service")
	return header + tabs + m.table.View() + "\n" + hint
}

func (m LibraryModel) renderTabs() string {
	tabs := make([]string, len(m.services))
	for i, svc := range m.services {
		if i == m.activeIdx {
			tabs[i] = activeTabStyle.Render(svc.name)
		} else {
			tabs[i] = tabStyle.Render(svc.name)
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func toLibraryRows(items []domain.LibraryItem) []table.Row {
	rows := make([]table.Row, 0, len(items))
	for _, item := range items {
		hasFile := "No"
		if item.HasFile {
			hasFile = "Yes"
		}
		year := ""
		if item.Year > 0 {
			year = fmt.Sprintf("%d", item.Year)
		}
		rows = append(rows, table.Row{item.Title, year, hasFile})
	}
	return rows
}
