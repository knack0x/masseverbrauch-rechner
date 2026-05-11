# AGENTS.md - Instructions for AI Agents

## Context File Maintenance

The `CONTEXT.md` file is the single source of truth for project status. **Always keep it up to date.**

### When to Update CONTEXT.md

Update `CONTEXT.md` after:
- Creating new files or directories
- Completing any task from the "What's Left" checklist
- Adding new tasks or discovering new requirements
- Making significant changes to existing code
- Completing a commit

### How to Update

1. Read `CONTEXT.md` first
2. Update the "What's Done" section - mark completed items with `[x]`
3. Update the "What's Left" section - add new items, mark completed ones
4. Update "Project Structure" if new directories/files were added
5. Commit the updated `CONTEXT.md` with the related changes

### Checklist Format

- `[ ]` for pending tasks
- `[x]` for completed tasks

### Example Workflow

```
1. Complete a task (e.g., create Docker setup)
2. Read CONTEXT.md
3. Move item from "What's Left" to "What's Done"
4. Save CONTEXT.md
5. Commit: "Add Docker setup and update CONTEXT.md"
```

## Project Root

Always use `/Users/knackwurstking/Git/knack0x/masseverbrauch-rechner` as the working directory.

## URL Conventions

The app is deployed behind an nginx reverse proxy at `https://knackwurstking.com/masseverbrauch-rechner/`. All URLs in templates and scripts MUST be relative (not starting with `/`) so they resolve correctly both at root during local dev and behind the subpath in production.

✅ `hx-post="calculate"` — resolves to `/calculate` at root or `/masseverbrauch-rechner/calculate` behind proxy
❌ `hx-post="/calculate"` — breaks behind proxy (points to domain root)

✅ `fetch('api/version')` — resolves correctly regardless of base path
❌ `fetch('/api/version')` — breaks behind proxy

## Logging Conventions

All server log messages MUST use a bracketed prefix indicating the level or middleware type:

| Prefix | When to use |
|--------|-------------|
| `[req]` | Request logging middleware |
| `[info]` | Server startup, operational info |
| `[warn]` | Recoverable issues, unexpected states |
| `[error]` | Runtime errors (API failure, template exec, decode failure) |
| `[debug]` | Detailed debug output (not currently used) |
| `[fatal]` | Unrecoverable errors that exit the process |

Examples:
- `log.Printf("[info] Server starting on :%s", port)`
- `log.Printf("[error] api: %v", err)`
- `log.Fatalf("[fatal] Failed to parse templates: %v", err)`

## Project Structure

```
masseverbrauch-rechner/
├── web/                 # Go web server (package main, has its own go.mod)
│   ├── main.go
│   ├── handlers.go
│   ├── types.go
│   ├── templates.go
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
├── Makefile
├── docker-compose.yml
├── VERSION              # Version file
├── CHANGELOG.md         # Version history
├── AGENTS.md
├── CONTEXT.md
└── README.md
```
