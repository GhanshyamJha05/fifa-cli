# ⚽ FIFA World Cup 2026 CLI

A production-quality terminal dashboard for the 2026 FIFA World Cup — built with Go, Cobra, Bubble Tea, Lip Gloss, and Bubbles.

```
⚽ FIFA WORLD CUP 2026 ⚽
━━━━━━━━━━━━━━━━━━━━━━━━━
  USA · CANADA · MEXICO
```

## Features

- **Interactive dashboard** — launch with `fifa` for keyboard-driven navigation
- **Teams & squads** — browse all nations, view squads with positions, numbers, clubs, and captains
- **Matches** — fixtures, results, today's games, live scores, and upcoming matches
- **Standings & bracket** — group tables and knockout stage layout
- **Statistics** — top scorers, assists, and more
- **Search** — find teams, players, and matches
- **Offline cache** — BoltDB caching with stale fallback when API is unavailable
- **Themes** — dark, light, and FIFA color schemes

## Quick Start

### Prerequisites

- Go 1.22+ (tested on Go 1.25)

### Install

```bash
git clone <repo-url>
cd worldcupCLI
go build -o fifa ./cmd/fifa
```

### Run

```bash
# Interactive dashboard
./fifa

# Commands
./fifa teams
./fifa team Brazil
./fifa squad Argentina
./fifa player "Lionel Messi"
./fifa matches
./fifa matches today
./fifa next
./fifa results
./fifa standings
./fifa bracket
./fifa stats
./fifa search Messi
```

## Configuration

Copy `config.example.yaml` to `~/.fifa-cli/config.yaml`:

```yaml
api_key: "your-api-football-key"
use_mock: false
theme: fifa
cache_ttl: 15m
```

Or set environment variables:

```bash
export FIFA_API_KEY="your-key"
export FIFA_USE_MOCK=false
export FIFA_THEME=fifa
```

### Live Data

This app integrates with [API-Football](https://www.api-football.com/) (league ID `1`, season `2026`). Without an API key, rich demo data is used automatically.

## Interactive Mode

| Key | Action |
|-----|--------|
| `1`–`5` | Jump to Home / Teams / Matches / Standings / Stats |
| `↑/↓` or `j/k` | Navigate lists |
| `←/→` or `h/l` | Switch tabs |
| `Enter` | View details |
| `/` | Search |
| `?` | Help |
| `q` | Quit |

## Architecture

```
cmd/fifa/           Entry point
internal/
  api/              Provider interface, API-Football client, mock data
  cache/            BoltDB response cache
  config/           Viper configuration
  domain/           Domain models
  service/          Business logic layer
  cmd/              Cobra commands
  ui/
    styles/         Lip Gloss themes
    render/         CLI output rendering
    tui/            Bubble Tea interactive app
```

## Testing

```bash
go test ./...
```

## Tech Stack

| Component | Library |
|-----------|---------|
| CLI | Cobra |
| TUI | Bubble Tea + Bubbles |
| Styling | Lip Gloss |
| Config | Viper |
| HTTP | Resty |
| Cache | BoltDB |
| Logging | slog |

## License

MIT
