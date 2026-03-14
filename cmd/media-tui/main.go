package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/MarioFronza/media-tui/internal/ui"
	"github.com/MarioFronza/media-tui/internal/usecase"
)

func main() {
	searchUC := usecase.NewSearchUseCase()
	p := tea.NewProgram(ui.NewApp(searchUC), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
