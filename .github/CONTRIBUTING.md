# Contributing to media-tui

## Prerequisites

- Go 1.23+
- Access to at least one *arr service (Radarr, Sonarr, Lidarr, or Readarr)
- `golangci-lint` for linting

## Setup

```bash
git clone https://github.com/MarioFronza/media-tui
cd media-tui
cp config.yaml.example config.yaml  # fill in your API keys
go test ./...
```

## Project structure

See [CLAUDE.md](../CLAUDE.md) for the full architecture overview and implementation status.

## Workflow

1. Pick an open issue or create one before starting work
2. Create a branch: `git checkout -b feat/your-feature`
3. Write tests for your changes
4. Run `go test ./...` and `golangci-lint run` before pushing
5. Open a PR — CI must be green to merge

## Architecture rules

- `domain` has zero external dependencies
- `usecase` depends only on `domain` interfaces — never on `adapter` or `ui`
- `ui` calls usecases only — never adapters directly
- All API calls are `tea.Cmd` (non-blocking)

## Commit style

```
type: short description

feat: add library screen
fix: handle empty search results
chore: add golangci-lint config
```
