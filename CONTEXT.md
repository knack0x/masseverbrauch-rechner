# CONTEXT.md - Masseverbrauch-Rechner

## Project Overview

A web application to calculate mass consumption in the press department of tile production.
Rebuild of the existing static app at https://knack0x.github.io/masseverbrauch-rechner.github.io/

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

## Tech Stack

- **Frontend**: Plain PHP, HTMX, HTML/CSS (mobile-first, improved from existing app)
- **Backend API**: Go with standard library only
- **JavaScript**: Minimal
- **Deployment**: Docker on Mac Studio → https://knackwurstking.com/masseverbrauch-rechner/

## Calculation Logic

- **Hauptmasse** = flow (kg/s) × runtime (seconds)
- **Flakes (per slot)** = weight before − weight after
- **Percentages** = each mass / total mass × 100
- Total = 100% of all masses combined
- Max 5 flakes slots ("Tower Slot 1" through "Tower Slot 5")

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
    {"name": "Tower Slot 1", "kg": 15.0, "percent": 1.48},
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
  - [x] Results display (dialog/modal)
  - [x] `calculate.php` — PHP proxy to Go API
- [x] Docker setup for both services
  - [x] `api/Dockerfile` — Go API container
  - [x] `web/Dockerfile` — PHP + Apache container
  - [x] `docker-compose.yml` — orchestrates both services
- [x] Created `Makefile` for build/install commands
- [x] Added local development targets (dev-api, dev-web) to Makefile
- [x] Created project root `README.md` with Makefile documentation
- [x] Restructured project: moved `web/` and `api/` to root, removed `php/` directory

## What's Left

- [x] Nginx config for path `knackwurstking.com/masseverbrauch-rechner/`
- [x] Create Makefile for build/install commands
- [x] Create project root README.md with Makefile documentation
- [ ] Test the complete setup (Go API + PHP frontend)
- [ ] Deploy to Mac Studio

## Reference

- Existing app: https://knack0x.github.io/masseverbrauch-rechner.github.io/
- Source: https://github.com/knack0x/masseverbrauch-rechner.github.io
