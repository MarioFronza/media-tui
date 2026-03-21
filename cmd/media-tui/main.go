package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/MarioFronza/media-tui/internal/adapter/api"
	"github.com/MarioFronza/media-tui/internal/adapter/config"
	"github.com/MarioFronza/media-tui/internal/domain"
	"github.com/MarioFronza/media-tui/internal/ui"
	"github.com/MarioFronza/media-tui/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	var repos []domain.MediaRepository
	var libraryUCs []*usecase.LibraryUseCase

	if cfg.Radarr.Enabled {
		adapter := api.NewRadarrAdapter(cfg.Radarr.URL, cfg.Radarr.APIKey)
		repos = append(repos, adapter)
		libraryUCs = append(libraryUCs, usecase.NewLibraryUseCase(adapter))
	}

	if cfg.Sonarr.Enabled {
		adapter := api.NewSonarrAdapter(cfg.Sonarr.URL, cfg.Sonarr.APIKey)
		repos = append(repos, adapter)
		libraryUCs = append(libraryUCs, usecase.NewLibraryUseCase(adapter))
	}

	if cfg.Lidarr.Enabled {
		adapter := api.NewLidarrAdapter(cfg.Lidarr.URL, cfg.Lidarr.APIKey)
		repos = append(repos, adapter)
		libraryUCs = append(libraryUCs, usecase.NewLibraryUseCase(adapter))
	}

	if cfg.Readarr.Enabled {
		adapter := api.NewReadarrAdapter(cfg.Readarr.URL, cfg.Readarr.APIKey)
		repos = append(repos, adapter)
		libraryUCs = append(libraryUCs, usecase.NewLibraryUseCase(adapter))
	}

	searchUC := usecase.NewSearchUseCase(repos...)
	queueUC := usecase.NewQueueUseCase(repos...)

	p := tea.NewProgram(ui.NewApp(searchUC, queueUC, libraryUCs...), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
