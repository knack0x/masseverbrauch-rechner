# Changelog

All notable changes to this project will be documented in this file.

## [2.0.0] - 2026-05-08

### Changed
- Complete rewrite of frontend: PHP replaced with Go html/templates
- Web server is now a Go binary (http.ServeMux) instead of PHP + Apache
- Both services (`api/`, `web/`) now have independent go.mod files
- `web/web.go` split into `main.go`, `handlers.go`, `types.go`, `templates.go`
- Form inputs sync to URL query params via `history.replaceState` for reload persistence

### Fixed
- Slot index mismatch in form parsing: all 5 slots are now sent to the API instead of skipping empty ones, preserving correct Tower Slot naming

### Removed
- PHP files (`index.php`, `calculate.php`, `version.php`, `cache.php`)
- PHP-based Dockerfile replaced with Go multi-stage build

### Technical Details
- **Frontend**: Go html/templates + HTMX + HTML/CSS (mobile-first)
- **Backend API**: Go 1.21 (standard library)
- **Web Server**: Go http.ServeMux (proxies to API)
- **Deployment**: Docker + Docker Compose
- **Exposed Ports**: API on 50570, Web on 50571 (Docker); API on 8080, Web on 8081 (local dev)

## [1.1.0] - 2026-05-03

### Added
- PWA (Progressive Web App) support with manifest.json
- PWA icon images (48x48 to 1000x1000) in web/assets/icons/
- PWA meta tags in index.php (apple-touch-icon, manifest link)
- Web app manifest for installability and offline capability

## [1.0.0] - 2026-05-02

### Added
- Initial release of Masseverbrauch-Rechner
- PHP frontend with HTMX for dynamic form submission
- Go API with `/api/calculate` endpoint for mass consumption calculations
- Docker support with docker-compose orchestration
- Local development targets (`make dev-api`, `make dev-web`)
- Version display in web UI (Web version + API version)
- Version endpoints (`/api/version` for API, `version.php` for web)
- Cache mechanism with ETag support (cache.php)
- German translations for slot names (Turmposition 1-5)
- API_URL environment variable to support both Docker and local development
- Makefile targets for build, deploy, test, and development

### Fixed
- Docker build cache issues (added `--no-cache` and `--build` flags)
- API URL configuration for Docker networking (using service names)
- VERSION file path resolution in both API and web UI
- PHP warning suppression for missing VERSION file in dev mode
- Removed deprecated `curl_close()` from version.php

### Technical Details
- **Frontend**: PHP 8.2 + HTMX + HTML/CSS (mobile-first design)
- **Backend**: Go 1.21 (standard library only)
- **Deployment**: Docker + Docker Compose
- **Exposed Ports**: API on 50570, Web on 50571 (Docker); API on 8080 (local dev)
