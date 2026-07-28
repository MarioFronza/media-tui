package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type listLifecycle struct {
	spinner spinner.Model
	loading bool
}

func newListLifecycle() listLifecycle {
	return listLifecycle{spinner: spinner.New()}
}

func (lc listLifecycle) startFetch() (listLifecycle, tea.Cmd) {
	lc.loading = true
	return lc, lc.spinner.Tick
}
