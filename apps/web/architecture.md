# MacroGuard Web Architecture

## Overview

The MacroGuard web frontend is a **SvelteKit** application with a mobile-first, PWA-ready design. It uses TailwindCSS v4 for styling and DaisyUI v5 for UI components.

## Tech Stack

| Technology    | Version | Purpose                   |
| ------------- | ------- | ------------------------- |
| SvelteKit     | 2.x     | Full-stack framework      |
| Svelte        | 5.x     | UI framework (Runes mode) |
| TailwindCSS   | 4.x     | Utility-first CSS         |
| DaisyUI       | 5.x     | Component library         |
| Lucide Svelte | -       | Icon library              |
| Bun           | -       | Package manager           |

---

## Directory Structure

```
apps/web/
├── src/
│   ├── app.css              # TailwindCSS + DaisyUI + NutriFresh theme
│   ├── app.html             # HTML template with PWA meta tags
│   ├── app.d.ts             # TypeScript declarations
│   ├── service-worker.ts    # PWA offline support
│   ├── lib/                 # Shared utilities and components
│   │   └── assets/          # Static assets (favicon, etc.)
│   └── routes/              # SvelteKit file-based routing
│       ├── +layout.svelte   # App shell (Navbar, bottom nav)
│       └── +page.svelte     # Homepage with upload UI
├── static/
│   └── manifest.json        # PWA manifest
├── vite.config.ts           # Vite + TailwindCSS plugin
├── svelte.config.js         # SvelteKit configuration
└── tsconfig.json            # TypeScript configuration
```

---

## Architecture Patterns

### Component Structure

```
┌─────────────────────────────────────────────────────────┐
│                   +layout.svelte                        │
│   ┌─────────────────────────────────────────────────┐   │
│   │              Navbar (Desktop)                   │   │
│   └─────────────────────────────────────────────────┘   │
│   ┌─────────────────────────────────────────────────┐   │
│   │                                                 │   │
│   │              Page Content                       │   │
│   │              (+page.svelte)                     │   │
│   │                                                 │   │
│   └─────────────────────────────────────────────────┘   │
│   ┌─────────────────────────────────────────────────┐   │
│   │           Bottom Nav (Mobile only)              │   │
│   └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

### State Management

Using **Svelte 5 Runes** for reactive state:

```svelte
<script lang="ts">
	let count = $state(0); // Reactive state
	let doubled = $derived(count * 2); // Derived value
	let { data } = $props(); // Component props
</script>
```

---

## Design System: NutriFresh

### Color Palette

| Role      | Color   | OKLCH                     | Usage               |
| --------- | ------- | ------------------------- | ------------------- |
| Primary   | Emerald | `oklch(62.8% 0.21 142.5)` | CTAs, active states |
| Secondary | Orange  | `oklch(70.5% 0.21 41.3)`  | Macros, energy      |
| Accent    | Cyan    | `oklch(70% 0.15 195)`     | Water, hydration    |
| Base      | White   | `oklch(100% 0 0)`         | Content background  |

### Theme Configuration

Defined in `src/app.css` using TailwindCSS v4 CSS-first approach:

```css
@import "tailwindcss";
@plugin "daisyui" {
	themes:
		nutrifresh --default,
		dark --prefersdark;
}

[data-theme="nutrifresh"] {
	--color-primary: oklch(62.8% 0.21 142.5);
	/* ... */
}
```

### Mobile-First Principles

- **48px minimum tap targets** for thumb-friendly interactions
- **Bottom navigation** on mobile, top navbar on desktop
- **`overscroll-behavior: none`** for native app feel
- **Tactile feedback** with `active:scale-95` on buttons

---

## PWA Configuration

### Manifest (`static/manifest.json`)

```json
{
	"name": "MacroGuard",
	"display": "standalone",
	"theme_color": "#22c55e"
}
```

### Service Worker (`src/service-worker.ts`)

- **Cache-first** for static assets
- **Network-first** for API calls
- **Offline fallback** when network unavailable

---

## API Integration

### Backend Communication

```
Frontend (localhost:5173)  ──HTTP──▶  Backend (localhost:8080)
         │                                      │
         │    POST /api/nutrition/scan          │
         │◀─────── JSON Response ───────────────│
```

### Planned API Client

Future: Generate TypeScript client from OpenAPI spec:

```bash
# From apps/web
npx openapi-typescript http://localhost:8080/openapi.json -o src/lib/api/schema.d.ts
```

---

## Key Routes

| Route      | Component      | Description             |
| ---------- | -------------- | ----------------------- |
| `/`        | `+page.svelte` | Homepage with upload UI |
| `/scan`    | (planned)      | Camera/upload flow      |
| `/history` | (planned)      | Scan history            |
| `/profile` | (planned)      | User settings           |
