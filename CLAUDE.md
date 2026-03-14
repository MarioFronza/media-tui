# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`media-tui` is a terminal UI application for managing the \*arr media stack (Radarr, Sonarr, Lidarr, Readarr) from a single interface. Users can search, add, and manage media across all four services without switching between web UIs.

## Tech Stack

- **Language**: Go
- **TUI Framework**: [Bubbletea](https://github.com/charmbracelet/bubbletea) (Elm architecture: Model → Update → View)
- **UI Components**: [Bubbles](https://github.com/charmbracelet/bubbles) (tables, inputs, spinners, lists)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **HTTP**: Standard `net/http` + `encoding/json`
- **Config**: YAML via `gopkg.in/yaml.v3`

## Commands

```bash
# Run the app
go run ./cmd/media-tui

# Build
go build -o media-tui ./cmd/media-tui

# Run tests
go test ./...

# Run a single test
go test ./internal/usecase/... -run TestSearch

# Lint (requires golangci-lint)
golangci-lint run
```

## Architecture

The app follows Clean Architecture. Dependencies point inward: `ui` → `usecase` → `domain` ← `adapter`. The domain layer has zero external dependencies.

```
cmd/media-tui/          # Entrypoint: wires dependencies, starts Bubbletea
  main.go               # TODO: wire config, adapters, usecases, start tea.Program
internal/
  domain/               # ✅ DONE — Core entities and repository interfaces (no external deps)
    media.go            # MediaItem, QueueItem, LibraryItem structs + MediaType constants
    repository.go       # MediaRepository interface: Search, Add, GetQueue, GetLibrary, MediaType
  usecase/              # ✅ DONE — Application business logic; depends only on domain interfaces
    search.go           # SearchUseCase: fan-out search across all repos concurrently (goroutines + sync.WaitGroup)
    library.go          # LibraryUseCase: List() and Add() for a single service repo
    queue.go            # QueueUseCase: aggregates GetQueue() across all repos (errors skipped)
    mock_test.go        # mockRepo: shared test double implementing domain.MediaRepository
    search_test.go      # Unit tests for SearchUseCase
    library_test.go     # Unit tests for LibraryUseCase
    queue_test.go       # Unit tests for QueueUseCase
  adapter/              # TODO
    api/                # HTTP adapters implementing domain.MediaRepository
      client.go         # Shared HTTP logic, auth headers, error handling
      radarr.go         # Radarr (/api/v3/movie/*)
      sonarr.go         # Sonarr (/api/v3/series/*)
      lidarr.go         # Lidarr (/api/v1/artist/*)
      readarr.go        # Readarr (/api/v1/book/*)
    config/
      config.go         # Loads config.yaml, maps to domain config structs
  ui/                   # TODO — Bubbletea TUI; depends on usecases, never on adapters directly
    app.go              # Root model; owns screen state, routes tea.Msg to sub-models
    search.go           # Unified search screen
    library.go          # Library browser per service
    queue.go            # Download queue viewer
    detail.go           # Item detail / add-to-library confirmation
config.yaml             # User config: API keys, base URLs, enabled services
```

### Implementation Status

| Layer | Status | Notes |
|---|---|---|
| `domain` | ✅ Done | Entities and `MediaRepository` interface |
| `usecase` | ✅ Done | Search, Library, Queue — all unit tested |
| `adapter/config` | ⬜ Todo | Issue #6 |
| `adapter/api` | ⬜ Todo | Issues #7–#11 |
| `ui` | ⬜ Todo | Issues #12–#16 |
| `cmd/media-tui` | ⬜ Todo | Issue #17 |

### Domain Entities

```go
// domain.MediaItem — result of Search; passed to Add
type MediaItem struct {
    ID       int
    Title    string
    Year     int
    Overview string
    Type     MediaType  // "movie" | "series" | "artist" | "book"
    Added    bool
}

// domain.LibraryItem — result of GetLibrary
type LibraryItem struct {
    ID        int
    Title     string
    Year      int
    HasFile   bool
    MediaType MediaType
}

// domain.QueueItem — result of GetQueue
type QueueItem struct {
    ID        int
    Title     string
    Status    string
    TimeLeft  string
    MediaType MediaType
}
```

### Usecase API

```go
// SearchUseCase — fan-out across all repos concurrently
uc := usecase.NewSearchUseCase(repo1, repo2, ...)
results := uc.Execute(term)  // []domain.MediaItem, errors silently dropped

// LibraryUseCase — scoped to a single repo
uc := usecase.NewLibraryUseCase(repo)
items, err := uc.List()
err = uc.Add(item)

// QueueUseCase — aggregates across all repos, errors skipped per repo
uc := usecase.NewQueueUseCase(repo1, repo2, ...)
items, err := uc.Execute()
```

### Key Patterns

**Dependency rule**: `domain` defines the `MediaRepository` interface. `adapter/api` implements it. `usecase` receives it via constructor injection. `ui` calls usecases only — never adapters directly.

**Async API calls**: usecases are invoked as `tea.Cmd` (returning `tea.Msg`) so they never block the UI. Use a spinner while fetching.

**Concurrent fan-out**: the search usecase fires goroutines for each enabled service and collects results via a channel, then returns a single aggregated `tea.Msg`.

**Screen routing**: the root app model holds a `currentScreen` enum and a typed reference to each sub-model. `Update()` forwards `tea.Msg` to the active sub-model and handles screen-switch messages.

**Config**: `config.yaml` at the working directory (or `~/.config/media-tui/config.yaml`). Each service has `enabled`, `url`, and `api_key` fields. Disabled services are skipped in fan-out calls. Services are hosted on a homelab accessible via VPN using the pattern `http://homelab:<port>` (Radarr: 7878, Sonarr: 8989, Lidarr: 8686, Readarr: 8787).

### \*arr API Conventions

All four services share the same API structure (they're all forks of the same codebase):

- Auth: `X-Api-Key` header
- Radarr/Sonarr base path: `/api/v3/` — Lidarr/Readarr: `/api/v1/`
- Search: `GET /<base>/lookup?term=<query>`
- Add: `POST /<base>/<resource>` with a JSON body
- Queue: `GET /<base>/queue`
- Resources: `movie` (Radarr), `series` (Sonarr), `artist` (Lidarr), `book` (Readarr)
