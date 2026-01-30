<h1 align="center" style="border-bottom: none">
    <a href="https://github.com/DoGab/MacroGuard" target="_blank"><img alt="MacroGuard" src=""></a><br>MacroGuard
</h1>

<p align="center">An AI-powered macro and nutrition scanner application.</p>

<div align="center">

[![Go API CI](https://github.com/DoGab/MacroGuard/actions/workflows/go-ci.yml/badge.svg)](https://github.com/DoGab/MacroGuard/actions/workflows/go-ci.yml)
[![Web CI](https://github.com/DoGab/MacroGuard/actions/workflows/web-ci.yml/badge.svg)](https://github.com/DoGab/MacroGuard/actions/workflows/web-ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/DoGab/MacroGuard/apps/api-go)](https://goreportcard.com/report/github.com/DoGab/MacroGuard/apps/api-go)

</div>

## Architecture

See [architecture.md](architecture.md) for system design, component interactions, and detailed documentation links.

| App | Description | Docs |
|-----|-------------|------|
| `apps/api-go` | Go REST API (Gin + Huma + Genkit) | [architecture.md](apps/api-go/architecture.md) |
| `apps/web` | SvelteKit frontend (TailwindCSS + DaisyUI) | [architecture.md](apps/web/architecture.md) |

## Quick Start

```bash
# Start both frontend and backend
make dev
```

| Service | URL |
|---------|-----|
| Frontend | http://localhost:5173 |
| API | http://localhost:8080 |
| API Docs | http://localhost:8080/docs |

## Development Commands

Run `make help` for all available commands:

```bash
make dev         # Start both servers concurrently
make dev-api     # Start Go API only (with Air hot reload)
make dev-web     # Start SvelteKit frontend only
make install     # Install all dependencies
```

## Tech Stack

- **Backend:** Go 1.22+, Gin, Huma v2, Genkit
- **Frontend:** SvelteKit, Svelte 5, TailwindCSS v4, DaisyUI v5
- **AI:** Genkit (ready for Gemini/OpenAI integration)
