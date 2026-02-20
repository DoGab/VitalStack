# VitalStack Web Frontend

SvelteKit frontend for the VitalStack food nutrition scanner.

📄 **Architecture:** See [architecture.md](architecture.md) for design patterns, structure, and theme documentation.

## Tech Stack

- **SvelteKit** (Svelte 5 with Runes)
- **TailwindCSS v4** (CSS-first configuration)
- **shadcn-svelte** (UI components via Bits UI)
- **Lucide Svelte** (icons)

---

## Setup

### 1. Install Dependencies

```bash
pnpm install
```

### 2. Configure TailwindCSS v4

TailwindCSS v4 uses a Vite plugin (no `tailwind.config.js`). See `vite.config.ts`:

```typescript
import tailwindcss from "@tailwindcss/vite";
export default defineConfig({
  plugins: [tailwindcss(), sveltekit()]
});
```

### 3. Theme Configuration

Custom themes are defined in `src/app.css` using CSS variables with `data-theme` attributes:

- `organic` — Light theme (Organic Premium)
- `darkorganic` — Dark theme (Dark Organic, default)

Font themes are controlled via `data-font-theme`:

- `classic` — Playfair Display + Inter
- `modern` — Outfit + Inter

> **Note:** TailwindCSS v4 uses OKLCH colors. Convert with [oklch.com](https://oklch.com/).

---

## Development

```bash
# From project root
make dev-web

# Or from this directory
pnpm run dev
```

## Building

```bash
pnpm run build
pnpm run preview
```

---

## Key Files

| File                        | Purpose                                |
| --------------------------- | -------------------------------------- |
| `vite.config.ts`            | TailwindCSS v4 Vite plugin             |
| `src/app.css`               | Tailwind + VitalStack theme (CSS vars) |
| `src/app.html`              | HTML template with `data-theme`        |
| `src/routes/+layout.svelte` | App shell (Sidebar + Mobile Dock)      |
| `src/routes/+page.svelte`   | Homepage                               |
| `static/manifest.json`      | PWA manifest                           |
| `src/service-worker.ts`     | PWA offline support                    |

---

## PWA

The app is PWA-ready:

- `static/manifest.json` - App metadata
- `src/service-worker.ts` - Offline caching
- `src/app.html` - Mobile meta tags

**TODO:** Add icons `static/icon-192.png` and `static/icon-512.png`
