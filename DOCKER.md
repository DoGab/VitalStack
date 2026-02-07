# Docker Deployment Guide

This document covers how to build and deploy VitalStack using Docker.

## Quick Start (Production)

Pull and run pre-built images from GitHub Container Registry:

```bash
# Set your API key
export GEMINI_API_KEY=your_api_key_here

# Pull and run
docker compose up
```

Visit:
- **Frontend**: http://localhost:3000
- **API**: http://localhost:8080

---

## Environment Variables

### API (`apps/api-go`)

| Variable | Description | Default |
|----------|-------------|---------|
| `LOGGING_LEVEL` | Log level (debug/info/warn/error) | `info` |
| `LOGGING_ENCODING` | Log format (json/logfmt) | `json` |
| `SERVER_ADDR` | Server listen address | `0.0.0.0:8080` |
| `SERVER_ORIGIN` | Allowed CORS origins (comma-separated) | `http://localhost:3000` |
| `GEMINI_API_KEY` | Gemini AI API key for Genkit | *required* |
| `DEV_MODE_ENABLED` | Enable dev mode (allows all CORS origins) | `false` |
| `DEV_MOCKS_NUTRITION_SERVICE` | Use mock AI service | `false` |

### Frontend (`apps/web`)

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `3000` |
| `ORIGIN` | Public URL (for CORS) | `http://localhost:3000` |
| `NODE_ENV` | Environment mode | `production` |
| `PUBLIC_API_URL` | API URL for browser requests | `""` (empty, uses /api proxy) |

---

## Building Individual Images

### API

```bash
cd apps/api-go
docker build -t vitalstack-api .

# Run
docker run -p 8080:8080 \
  -e GEMINI_API_KEY=your_key \
  -e SERVER_ADDR=0.0.0.0:8080 \
  vitalstack-api
```

### Frontend

```bash
cd apps/web
docker build -t vitalstack-web .

# Run with API URL configured at runtime:
docker run -p 3000:3000 \
  -e ORIGIN=https://vs.dgit.ch \
  -e PUBLIC_API_URL=https://vs-api.dgit.ch \
  vitalstack-web
```

---

## Production Deployment

For production at `vs.dgit.ch` and `vs-api.dgit.ch`:

### 1. Create `.env` file

```bash
# .env
GEMINI_API_KEY=your_production_api_key
```

### 2. Update docker-compose.yml

```yaml
services:
  api:
    environment:
      LOGGING_LEVEL: info
      # ... other env vars
  
  web:
    environment:
      ORIGIN: https://vs.dgit.ch
```

### 3. Deploy with reverse proxy

Use Traefik, Nginx, or Caddy as a reverse proxy:

```
vs.dgit.ch      → localhost:3000 (web)
vs-api.dgit.ch  → localhost:8080 (api)
```

### Example Caddy configuration

```caddyfile
vs.dgit.ch {
    reverse_proxy localhost:3000
}

vs-api.dgit.ch {
    reverse_proxy localhost:8080
}
```

---

## Local Testing

To build and test the containers locally (from source):

### Quick Start (Mock Mode)

```bash
# Build and run with local settings (no API key needed)
docker compose -f docker-compose.yml -f docker-compose.build.yml -f docker-compose.local.yml up --build
```

Visit:
- **Frontend**: http://localhost:3000
- **API**: http://localhost:8080

### With Real AI

```bash
# Set your API key
export GEMINI_API_KEY=your_key

# Build and run
docker compose -f docker-compose.yml -f docker-compose.build.yml -f docker-compose.local.yml up --build
```

The override files provide:
- `docker-compose.build.yml` — Builds from source (not GHCR)
- `docker-compose.local.yml` — Dev settings: mock AI, debug logging, localhost URLs

---

## Troubleshooting

### Container won't start

Check logs:
```bash
docker compose logs api
docker compose logs web
```

### Health check failing

Ensure the health endpoints are accessible:
```bash
curl http://localhost:8080/health
curl http://localhost:3000
```

### CORS issues

Make sure `ORIGIN` is set to your frontend's public URL.
