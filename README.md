# media-tui

> A terminal UI for managing your \*arr media stack from a single interface.

Manage **Radarr, Sonarr, Lidarr, and Readarr** without switching between web UIs — search, add media, browse your library, and monitor the download queue, all from the terminal.

<!-- demo gif here -->

## Features

- **Unified search** — fan-out search across all enabled services simultaneously
- **Add media** — select a result and add it to the correct service's library
- **Library browser** — browse what's already in each service
- **Download queue** — monitor active downloads aggregated from all services

## Prerequisites

- Go 1.23+
- At least one running \*arr service (Radarr, Sonarr, Lidarr, or Readarr)

## Installation

```bash
go install github.com/MarioFronza/media-tui/cmd/media-tui@latest
```

Or build from source:

```bash
git clone https://github.com/MarioFronza/media-tui
cd media-tui
go build -o media-tui ./cmd/media-tui
```

## Configuration

Copy the example config and fill in your API keys:

```bash
cp config.yaml.example config.yaml
```

```yaml
radarr:
  enabled: true
  url: http://localhost:7878
  api_key: your_api_key_here

sonarr:
  enabled: true
  url: http://localhost:8989
  api_key: your_api_key_here

lidarr:
  enabled: false
  url: http://localhost:8686
  api_key: your_api_key_here

readarr:
  enabled: false
  url: http://localhost:8787
  api_key: your_api_key_here
```

Config is loaded from `./config.yaml` or `~/.config/media-tui/config.yaml`. Disabled services are skipped entirely.

## Usage

```bash
./media-tui
```

| Key | Action |
|---|---|
| `/` | Search |
| `l` | Library |
| `q` | Queue |
| `Enter` | Select / confirm |
| `Esc` | Go back |
| `r` | Refresh |
| `ctrl+c` | Quit |

## Architecture

media-tui follows Clean Architecture — dependencies point inward:

```
ui → usecase → domain ← adapter
```

See [CLAUDE.md](CLAUDE.md) for the full architecture overview, implementation status, and development patterns.

## Contributing

See [CONTRIBUTING.md](.github/CONTRIBUTING.md).

## License

[MIT](LICENSE)
