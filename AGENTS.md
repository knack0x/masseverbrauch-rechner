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

## Project Structure

```
masseverbrauch-rechner/
├── web/                 # PHP frontend (HTMX + SSR)
│   ├── index.php
│   ├── calculate.php
│   ├── cache.php
│   └── Dockerfile
├── api/                 # Go API (stdlib)
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── Makefile
├── docker-compose.yml
├── AGENTS.md
├── CONTEXT.md
└── README.md
```
