# Masseverbrauch-Rechner

A web application to calculate mass consumption in the press department of tile production.
Rebuild of the existing static app with PHP + Go API architecture.

## Project Structure

```
test/
├── web/                 # PHP frontend (HTMX + SSR)
│   ├── index.php
│   ├── calculate.php
│   └── Dockerfile
├── api/                 # Go API (stdlib)
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── Makefile             # Build and deploy commands
├── docker-compose.yml   # Docker orchestration
├── AGENTS.md
├── CONTEXT.md
└── README.md
```

## Quick Start

### Prerequisites
- Docker & Docker Compose
- GNU Make

### Build & Run

```bash
# Build all Docker images
make build

# Start all services
make up

# API: http://localhost:8080
# Web: http://localhost:8081
```

### Makefile Targets

| Target | Description |
|--------|-------------|
| `make help` | Show available commands |
| `make build` | Build all Docker images (API + Web) |
| `make build-api` | Build Go API binary |
| `make build-web` | Build PHP Docker image |
| `make up` | Start all services with docker-compose |
| `make down` | Stop all services |
| `make test` | Test the API endpoint |
| `make test-web` | Test the full stack |
| `make logs` | View logs from all services |
| `make clean` | Remove containers and images |
| `make deploy` | Deploy to Mac Studio (requires DEPLOY_HOST) |
| `make dev-api` | Run Go API locally without Docker |
| `make dev-web` | Run PHP web locally without Docker |
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

- **Frontend**: PHP 8.2 + HTMX + HTML/CSS (mobile-first)
- **Backend API**: Go 1.21 (standard library)
- **Deployment**: Docker + Docker Compose
