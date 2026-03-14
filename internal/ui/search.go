package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/MarioFronza/media-tui/internal/domain"
	"github.com/MarioFronza/media-tui/internal/usecase"
)

type searchResultMsg struct {
	items []domain.MediaItem
}

type SearchModel struct {
	input    textinput.Model
	spinner  spinner.Model
	table    table.Model
	results  []domain.MediaItem
	loading  bool
	searched bool
	usecase  *usecase.SearchUseCase
}

func NewSearchModel(uc *usecase.SearchUseCase) SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Search movies, series, artists, books..."
	ti.Focus()

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	cols := []table.Column{
		{Title: "Title", Width: 40},
		{Title: "Year", Width: 6},
		{Title: "Type", Width: 10},
		{Title: "Added", Width: 7},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(false),
		table.WithHeight(15),
	)
	t.SetStyles(tableStyles())

	return SearchModel{
		input:   ti,
		spinner: sp,
		table:   t,
		usecase: uc,
	}
}

func (m SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case searchResultMsg:
		return m.handleSearchResult(msg)
	case spinner.TickMsg:
		return m.updateSpinner(msg)
	}
	return m.updateFocusedComponent(msg)
}

func (m SearchModel) handleKey(msg tea.KeyMsg) (SearchModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if m.table.Focused() {
			return m.selectResult()
		}
		return m.startSearch()
	case "esc":
		return m.focusInput()
	}
	return m.updateFocusedComponent(msg)
}

func (m SearchModel) startSearch() (SearchModel, tea.Cmd) {
	term := m.input.Value()
	if term == "" || m.loading {
		return m, nil
	}
	m.loading = true
	m.searched = true
	m.input.Blur()
	return m, tea.Batch(m.spinner.Tick, m.runSearch(term))
}

func (m SearchModel) selectResult() (SearchModel, tea.Cmd) {
	idx := m.table.Cursor()
	if idx < 0 || idx >= len(m.results) {
		return m, nil
	}
	item := m.results[idx]
	return m, func() tea.Msg { return SwitchToDetailMsg{Item: item} }
}

func (m SearchModel) focusInput() (SearchModel, tea.Cmd) {
	m.table.Blur()
	m.input.Focus()
	return m, nil
}

func (m SearchModel) handleSearchResult(msg searchResultMsg) (SearchModel, tea.Cmd) {
	m.loading = false
	m.results = msg.items
	m.table.SetRows(toTableRows(msg.items))
	m.table.Focus()
	return m, nil
}

func (m SearchModel) updateSpinner(msg spinner.TickMsg) (SearchModel, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m SearchModel) updateFocusedComponent(msg tea.Msg) (SearchModel, tea.Cmd) {
	var cmd tea.Cmd
	switch {
	case m.loading:
		m.spinner, cmd = m.spinner.Update(msg)
	case m.table.Focused():
		m.table, cmd = m.table.Update(msg)
	default:
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

func (m SearchModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Render("Search") + "\n\n"
	inputView := m.input.View() + "\n\n"

	if m.loading {
		return header + inputView + m.spinner.View() + " Searching..."
	}

	if m.searched && len(m.results) == 0 {
		return header + inputView + "No results found."
	}

	if len(m.results) > 0 {
		hint := lipgloss.NewStyle().Faint(true).Render("↑/↓ navigate · Enter select · Esc back to input")
		return header + inputView + m.table.View() + "\n" + hint
	}

	return header + inputView
}

func (m SearchModel) runSearch(term string) tea.Cmd {
	return func() tea.Msg {
		items := m.usecase.Execute(term)
		return searchResultMsg{items: items}
	}
}

func toTableRows(items []domain.MediaItem) []table.Row {
	rows := make([]table.Row, 0, len(items))
	for _, item := range items {
		added := "No"
		if item.Added {
			added = "Yes"
		}
		year := ""
		if item.Year > 0 {
			year = fmt.Sprintf("%d", item.Year)
		}
		rows = append(rows, table.Row{item.Title, year, string(item.Type), added})
	}
	return rows
}

func tableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	return s
}
