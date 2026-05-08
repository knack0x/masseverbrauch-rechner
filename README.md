# Masseverbrauch-Rechner

A web application to calculate mass consumption in the press department of tile production.
Rebuild of the existing static app with Go + HTMX architecture.

## Project Structure

```
masseverbrauch-rechner/
├── web/                 # Go web server (go.mod + main package + handlers + templates)
│   ├── main.go          # Entry point, Serve(), asset embed, config
│   ├── handlers.go      # HTTP handlers + cache middleware
│   ├── types.go         # Shared types
│   ├── templates.go     # Template init + custom functions
│   ├── go.mod
│   ├── templates/
│   │   ├── index.html
│   │   └── calculate.html
│   ├── assets/
│   │   ├── manifest.json
│   │   ├── favicon.ico
│   │   └── icons/
│   └── Dockerfile
├── api/                 # Go API (stdlib)
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── Makefile             # Build and deploy commands
├── docker-compose.yml   # Docker orchestration
├── VERSION              # Version file
├── CHANGELOG.md         # Version history
├── AGENTS.md
├── CONTEXT.md
└── README.md
```

## Quick Start

### Prerequisites
- Docker & Docker Compose
- GNU Make
- Go 1.21+ (for local dev)

### Build & Run

```bash
# Build all Docker images
make build

# Start all services
make up

# API: http://localhost:50570
# Web: http://localhost:50571
```

### Local Development (without Docker)

```bash
# Start Go API locally (port 8080)
make dev-api

# Start Go web server locally (port 8081)
make dev-web
```

The web frontend uses `API_URL` environment variable to connect to the API. In Docker, it's set automatically. For local dev, `make dev-web` sets it to `http://localhost:8080/api/calculate`.

### Makefile Targets

| Target | Description |
|--------|-------------|
| `make help` | Show available commands |
| `make build` | Build all Docker images (API + Web) |
| `make build-api` | Build Go API binary |
| `make build-web` | Build Go web Docker image |
| `make up` | Start all services with docker-compose |
| `make down` | Stop all services |
| `make test` | Test the API endpoint |
| `make test-web` | Test the full stack |
| `make logs` | View logs from all services |
| `make clean` | Remove containers and images |
| `make deploy` | Deploy to Mac Studio (requires DEPLOY_HOST) |
| `make dev-api` | Run Go API locally without Docker |
| `make dev-web` | Run Go web server locally without Docker |
| `make restart` | Restart all services |

### Deploy to Mac Studio

```bash
make deploy DEPLOY_HOST=your-mac-studio.local
```

## API Endpoint

**POST** `/api/calculate`

```json
{
  "flow": 1.5,
  "runtime_minutes": 10,
  "slots": [
    {"before": 100.0, "after": 85.0},
    {"before": 50.0, "after": 42.0}
  ]
}
```

Response:
```json
{
  "hauptmasse_kg": 900.0,
  "hauptmasse_percent": 88.67,
  "slots": [
    {"name": "Tower Slot 1", "kg": 15.0, "percent": 1.48}
  ],
  "total_kg": 1015.0
}
```

## Tech Stack

- **Frontend**: Go html/templates + HTMX + HTML/CSS (mobile-first)
- **Backend API**: Go 1.21 (standard library)
- **Web Server**: Go http.ServeMux (proxies to API)
- **Deployment**: Docker + Docker Compose

## Version

Current version: **1.1.0** (see [CHANGELOG.md](CHANGELOG.md) for details)
