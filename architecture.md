# VitalStack Architecture

This document provides a high-level overview of the VitalStack system architecture.

## System Overview

VitalStack is a **monorepo** containing a Go backend API and a SvelteKit frontend, designed to analyze food images and return nutritional macro information using AI.

```mermaid
C4Context
  title System Context diagram for VitalStack

  Person(user, "User", "A user of the VitalStack app tracking their food intake.")
  
  System(vitalstack, "VitalStack", "Allows users to log food, analyze nutrition via AI, and track macros.")
  
  System_Ext(supabase, "Supabase", "Stores user data, authentication, and food logs.")
  System_Ext(llm, "LLM Provider", "Gemini/OpenAI for image analysis and macro estimation.")
  System_Ext(off, "Open Food Facts", "Primary open food database.")
  System_Ext(fsvo, "Swiss FSVO", "Swiss food composition database (~1,100 generic foods).")
  System_Ext(usda, "USDA FoodData Central", "Secondary authoritative food database.")

  Rel(user, vitalstack, "Uses")
  Rel(vitalstack, supabase, "Reads from and writes to")
  Rel(vitalstack, llm, "Analyzes food images using")
  Rel(vitalstack, off, "Searches product barcodes")
  Rel(vitalstack, fsvo, "Looks up generic Swiss foods")
  Rel(vitalstack, usda, "Falls back for product data")
```

### Container Level Overview

```mermaid
C4Container
  title Container diagram for VitalStack

  Person(user, "User", "A user of the VitalStack app tracking their food intake.")

  System_Boundary(c1, "VitalStack") {
    Container(web, "Web Application", "SvelteKit", "Delivers the PWA, handles UI and client-side logic.")
    Container(api, "API Application", "Go, Gin, Huma", "Provides REST APIs, orchestrates AI and datasources.")
    ContainerDb(meili, "Local Cache", "Meilisearch", "Typo-tolerant product cache for fast search.")
  }

  System_Ext(supabase, "Supabase", "Stores user data, authentication, and food logs.")
  System_Ext(llm, "LLM Provider", "Gemini/OpenAI for image analysis.")
  System_Ext(off, "Open Food Facts", "Primary open food database.")
  System_Ext(usda, "USDA FoodData Central", "Secondary authoritative food database.")

  Rel(user, web, "Uses", "HTTPS")
  Rel(web, api, "Makes API calls to", "JSON/HTTPS")
  Rel(api, supabase, "Reads/Writes user data", "pgx")
  Rel(api, meili, "Reads/Writes product cache", "HTTP")
  Rel(api, llm, "Analyzes images", "gRPC/HTTP")
  Rel(api, off, "Fetches product data", "HTTPS")
  Rel(api, usda, "Fetches product data", "HTTPS")
```

---

## Application Architecture

### Frontend (`apps/web`)

| Layer | Technology | Purpose |
|-------|------------|---------|
| Framework | SvelteKit | Routing, SSR, file-based routing |
| UI | Svelte 5 | Reactive components with Runes |
| Styling | TailwindCSS v4 + Shadcn-svelte | Utility CSS + component library |
| Theme | NutriFresh | Custom emerald/orange/cyan palette |
| PWA | Service Worker | Offline support, installable app |

📄 **Detailed docs:** [apps/web/architecture.md](apps/web/architecture.md)

---

### Backend (`apps/api-go`)

| Layer | Technology | Purpose |
|-------|------------|---------|
| CLI | Cobra + Viper | Commands, config, flags |
| HTTP | Gin + Huma v2 | Router, OpenAPI 3.1 generation |
| Controller | Huma handlers | Request validation, DTO conversion |
| Service | Business logic | AI/Genkit integration (future) |
| AI | Genkit | LLM orchestration (Gemini/OpenAI) |

📄 **Detailed docs:** [apps/api-go/architecture.md](apps/api-go/architecture.md)

---

## Request Flows

### AI Food Scan

```mermaid
sequenceDiagram
    participant U as User (Browser)
    participant S as SvelteKit (Frontend)
    participant A as Go API (Controller)
    participant G as Genkit (AI/LLM)
    
    U->>S: Upload Image
    S->>A: POST /api/nutrition/scan
    A->>G: Analyze Image
    G-->>A: Macro Data
    A-->>S: JSON Response
    S-->>U: Display Macros
```

### Product Search & Barcode Scan

```mermaid
sequenceDiagram
    participant U as User (Browser)
    participant S as SvelteKit (Frontend)
    participant A as Go API (Controller)
    participant DS as Datasources (OFF/FSVO/USDA)
    
    U->>S: Search or Scan Barcode
    S->>A: GET /api/products/search or /barcode/{ean}
    A->>DS: Waterfall lookup (Cache → OFF → FSVO → USDA)
    DS-->>A: Product Data
    A-->>S: JSON Response
    S-->>U: Display Product Card
    U->>S: Select Serving & Log
    S->>A: POST /api/nutrition/log
    A-->>S: Log Confirmation
```

---

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/health` | Health check |
| `POST` | `/api/nutrition/scan` | Analyze food image |
| `POST` | `/api/nutrition/log` | Log a food entry (AI or product-based) |
| `GET` | `/api/products/search` | Full-text product search (waterfall) |
| `GET` | `/api/products/barcode/{ean}` | Lookup product by barcode |
| `GET` | `/docs` | OpenAPI documentation |
| `GET` | `/openapi.json` | OpenAPI 3.1 spec |

---

## Development Setup

```bash
# Start both frontend and backend
make dev

# Start individually
make dev-web   # Frontend on :5173
make dev-api   # Backend on :8080
```

---

## Project Structure

```
VitalStack/
├── apps/
│   ├── api-go/              # Go backend
│   │   ├── architecture.md  # Backend architecture docs
│   │   ├── cmd/             # CLI commands
│   │   ├── internal/        # Private packages
│   │   └── pkg/             # Public packages
│   └── web/                 # SvelteKit frontend
│       ├── architecture.md  # Frontend architecture docs
│       ├── src/             # Source code
│       └── static/          # Static assets
├── Makefile                 # Dev orchestration
└── README.md                # This file links here
```
