# CONTEXT.md - Masseverbrauch-Rechner

## Project Overview

A web application to calculate mass consumption in the press department of tile production.
Rebuild of the existing static app at https://knack0x.github.io/masseverbrauch-rechner.github.io/

## Project Structure

```
masseverbrauch-rechner/
├── web/                 # Go web server (package main, has its own go.mod)
│   ├── main.go          # Entry point
│   ├── server.go        # Serve(), asset embed, config
│   ├── handlers.go      # HTTP handlers + cache middleware
│   ├── types.go         # Shared types
│   ├── templates.go     # Template init + custom functions
│   ├── templates/
│   │   ├── index.html       # Main page (Go html/template)
│   │   └── calculate.html   # Result dialog (Go html/template)
│   ├── assets/
│   │   ├── manifest.json    # PWA manifest
│   │   ├── favicon.ico
│   │   └── icons/           # PWA icons (48-1000px)
 │   └── Dockerfile       # TODO: create Go multi-stage build
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

## Tech Stack

- **Frontend**: Go html/templates, HTMX, HTML/CSS (mobile-first)
- **Backend API**: Go with standard library only
- **Web Server**: Go http.ServeMux (proxies to Go API)
- **JavaScript**: Minimal
- **Deployment**: Docker on Mac Studio → https://knackwurstking.com/masseverbrauch-rechner/

## Calculation Logic

- **Hauptmasse** = flow (kg/s) × runtime (seconds)
- **Flakes (per slot)** = weight after − weight before
- **Percentages** = each mass / total mass × 100
- Total = 100% of all masses combined
- Max 5 flakes slots ("Turmposition 1" through "Turmposition 5" in UI, "Tower Slot 1" through "Tower Slot 5" in API)

## Input Fields

1. Flow rate (kg/s) — input id: `main-flow`
2. Runtime in minutes — input id: `main-runtime`
3. 5× Flakes slots, each with:
   - KG before measurement — input id: `tower-slot-{n}-before`
   - KG after measurement — input id: `tower-slot-{n}-after`

## API Contract

**Endpoint**: `POST /api/calculate`

**Request JSON:**
```json
{
  "flow": 1.5,
  "runtime_minutes": 10,
  "slots": [
    {"before": 100.0, "after": 85.0},
    {"before": 50.0, "after": 42.0},
    ...
  ]
}
```

**Response JSON:**
```json
{
  "hauptmasse_kg": 900.0,
  "hauptmasse_percent": 88.67,
  "slots": [
    {"name": "Tower Slot 1", "kg": 15.0, "percent": 1.48},  (translated to "Turmposition 1" in UI)
    ...
  ],
  "total_kg": 1015.0
}
```

## What's Done

- [x] Created `CONTEXT.md`
- [x] Created `AGENTS.md` with context file maintenance instructions
- [x] Created Go API (`api/` dir) with `/api/calculate` endpoint
- [x] Both services (`api/`, `web/`) have independent go.mod files
- [x] Converted PHP frontend to Go html/templates:
  - [x] `web/main.go` — entry point, `Serve()`, asset embed, config
  - [x] `web/handlers.go` — all HTTP handlers + cache middleware
  - [x] `web/types.go` — shared type definitions
  - [x] `web/templates.go` — template init, embed, custom functions
  - [x] `web/templates/index.html` — main page (Go template with range loop, version display, URL query param sync)
  - [x] `web/templates/calculate.html` — result dialog (Go template with German formatting)
  - [x] German number formatting (comma decimal, dot thousands separator)
  - [x] German translation for slot names (Turmposition 1-5)
  - [x] `manifest.json` moved to `web/assets/manifest.json`
  - [x] HTMX integration, version proxy, static assets serving
  - [x] Form inputs sync to URL query params via `history.replaceState` for reload persistence
  - [x] Fixed slot index bug: send all 5 slots to API instead of skipping empty ones
  - [x] Fixed absolute paths for subpath proxy: form action and version fetch use relative URLs
- [x] Docker setup for both services
  - [x] `api/Dockerfile` — Go API container
  - [x] `web/Dockerfile` — PHP + Apache container
  - [x] `docker-compose.yml` — orchestrates both services
- [x] Created `Makefile` for build/install commands
- [x] Added local development targets (dev-api, dev-web) to Makefile
- [x] Created project root `README.md` with Makefile documentation
- [x] Restructured project: moved `web/` and `api/` to root, removed `php/` directory
- [x] Added basic request logging to Go API (incoming requests, request data, calculation results)
- [x] Fixed slot calculation bug in Go API (inverted before/after subtraction causing zero results)
- [x] Nginx config for path `knackwurstking.com/masseverbrauch-rechner/`
- [x] Added reusable cache mechanism (cache.php) with ETag support to index.php and calculate.php
- [x] Fixed web/Dockerfile COPY instruction (build context already in web/, changed `COPY web/` to `COPY .`)
- [x] Changed exposed ports to avoid conflicts: API 50570:8080, Web 50571:80
- [x] Updated calculate.php error handling to show HTTP status code and message
- [x] Fixed API URL in calculate.php to use port 50570
- [x] Updated README.md and Makefile (dev-web) to reflect new ports
- [x] Fixed Makefile to force rebuild on `make up` (added --build flag)
- [x] Fixed Makefile build-web target to use --no-cache to avoid Docker layer caching
- [x] Fixed calculate.php API URL to use Docker service name (api:8080) instead of localhost
- [x] Added API_URL environment variable to support both Docker and local dev environments
- [x] Fixed local dev API URL to use port 8080 (Go default) instead of Docker-mapped port 50570
- [x] Added VERSION file (1.0.0 → 1.1.0) and version endpoints to API (/api/version) and web UI
- [x] Updated Dockerfiles to copy VERSION file into containers
- [x] Fixed Docker build context to use project root so VERSION file is accessible
- [x] Updated docker-compose.yml and Makefile to use project root as build context
- [x] Fixed VERSION file paths in both API and web UI to use correct relative paths
- [x] Added version.php proxy to fetch API version through PHP (fixes browser CORS issue)
- [x] Removed deprecated curl_close() from version.php
- [x] Fixed VERSION file reading to suppress warnings in dev mode (check file_exists first)
- [x] Added PWA icon images (48x48 to 1000x1000) to web/assets/icons/
- [x] Created manifest.json for PWA installation
- [x] Updated index.php with PWA meta tags, apple-touch-icon, and manifest link
- [x] Added structured log prefixes (`[req]`, `[info]`, `[warn]`, `[error]`, `[debug]`, `[fatal]`) to all web server logs

## What's Left

- [x] Test complete Go stack (API + Web server)

## Reference

- Existing app: https://knack0x.github.io/masseverbrauch-rechner.github.io/
- Source: https://github.com/knack0x/masseverbrauch-rechner.github.io
