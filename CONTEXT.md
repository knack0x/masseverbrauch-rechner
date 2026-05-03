# CONTEXT.md - Masseverbrauch-Rechner

## Project Overview

A web application to calculate mass consumption in the press department of tile production.
Rebuild of the existing static app at https://knack0x.github.io/masseverbrauch-rechner.github.io/

## Project Structure

```
masseverbrauch-rechner/
├── web/                 # PHP frontend (HTMX + SSR)
│   ├── index.php
│   ├── calculate.php
│   ├── version.php
│   ├── cache.php
│   └── Dockerfile
├── api/                 # Go API (stdlib)
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── Makefile             # Build and deploy commands
├── docker-compose.yml   # Docker orchestration
├── VERSION              # Version file (1.0.0)
├── CHANGELOG.md         # Version history
├── AGENTS.md
├── CONTEXT.md
└── README.md
```

## Tech Stack

- **Frontend**: Plain PHP, HTMX, HTML/CSS (mobile-first, improved from existing app)
- **Backend API**: Go with standard library only
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
- [x] Created PHP frontend (`web/` dir) with:
  - [x] `index.php` — main page with form (mobile-first styling)
  - [x] HTMX integration to call Go API via PHP proxy
  - [x] Results display (dialog/modal, centered on screen)
  - [x] `calculate.php` — PHP proxy to Go API
  - [x] German translations for slot names (Turmposition 1-5)
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
- [x] Added VERSION file (1.0.0) and version endpoints to API (/api/version) and web UI
- [x] Updated Dockerfiles to copy VERSION file into containers
- [x] Fixed Docker build context to use project root so VERSION file is accessible
- [x] Updated docker-compose.yml and Makefile to use project root as build context
- [x] Fixed VERSION file paths in both API and web UI to use correct relative paths
- [x] Added version.php proxy to fetch API version through PHP (fixes browser CORS issue)
- [x] Removed deprecated curl_close() from version.php
- [x] Fixed VERSION file reading to suppress warnings in dev mode (check file_exists first)

## What's Left

- [x] Test the complete setup (Go API + PHP frontend) - Verified working 2026-05-02
- [x] Deploy to Mac Studio
- [ ] Add icon images to web/ for Android and iOS support
- [ ] Add manifest.json for PWA installation (no service worker, caching only)

## Reference

- Existing app: https://knack0x.github.io/masseverbrauch-rechner.github.io/
- Source: https://github.com/knack0x/masseverbrauch-rechner.github.io
