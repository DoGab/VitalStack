<h1 align="center" style="border-bottom: none">
    <a href="https://github.com/DoGab/VitalStack" target="_blank"><img alt="VitalStack" src="./apps/web/src/lib/assets/logo/logo_text_dark.svg"></a><br>VitalStack
</h1>

<p align="center">An AI-powered macro and nutrition scanner application.</p>

<div align="center">

[![Go API CI](https://github.com/DoGab/VitalStack/actions/workflows/go-ci.yml/badge.svg)](https://github.com/DoGab/VitalStack/actions/workflows/go-ci.yml)
[![Web CI](https://github.com/DoGab/VitalStack/actions/workflows/web-ci.yml/badge.svg)](https://github.com/DoGab/VitalStack/actions/workflows/web-ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/DoGab/VitalStack/apps/api-go)](https://goreportcard.com/report/github.com/DoGab/VitalStack/apps/api-go)

</div>

## Features

- **Food Scan** — Scan food with your camera and get nutritional information (calories, protein, carbs, fat) 🚧
- **Food Log** — Log your food intake manually 📋
- **Nutrition Goals** — Set your nutrition goals (calories, protein, carbs, fat) 📋
- **Progress Tracking** — Track your progress over time 📋
- **Log Hydration** — Log your hydration intake 📋

> **Status:** ✅ Done · 🚧 In Progress · 📋 Planned

## Technical features

- **Containerized** - The application is containerized using Docker and Docker Compose. ✅
- **Configuration** - Configuration can happen via environment variables or a config file. ✅
- **Logging** - The application uses structured logging with logfmt and JSON formats. ✅
- **Health Checks** - The application has health checks for both the frontend and backend. 📋
- **Authentication** - The application has authentication using JWT. 📋

> **Status:** ✅ Done · 🚧 In Progress · 📋 Planned

## Architecture

See [architecture.md](architecture.md) for system design, component interactions, and detailed documentation links.

| App | Description | Docs |
|-----|-------------|------|
| `apps/api-go` | Go REST API (Gin + Huma + Genkit) | [architecture.md](apps/api-go/architecture.md) |
| `apps/web` | SvelteKit frontend (TailwindCSS + Shadcn-svelte) | [architecture.md](apps/web/architecture.md) |

## Quick Start

```bash
# Start both frontend and backend
make dev
```

| Service | URL |
|---------|-----|
| Frontend | https://localhost:5173 |
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
- **Frontend:** SvelteKit, Svelte 5, TailwindCSS v4, Shadcn-svelte
- **AI:** Genkit (ready for Gemini/OpenAI integration)
