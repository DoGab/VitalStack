# MacroGuard Web Frontend

SvelteKit frontend for the MacroGuard food nutrition scanner.

📄 **Architecture:** See [architecture.md](architecture.md) for design patterns, structure, and theme documentation.

## Tech Stack

- **SvelteKit** (Svelte 5 with Runes)
- **TailwindCSS v4** (CSS-first configuration)
- **DaisyUI v5** (via `@plugin` directive)
- **Lucide Svelte** (icons)
- **Bun** (package manager)

---

## Setup

### 1. Install Dependencies

```bash
bun install
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

Custom themes are defined in `src/app.css` using CSS variables:

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

> **Note:** TailwindCSS v4 uses OKLCH colors. Convert with [oklch.com](https://oklch.com/).

---

## Development

```bash
# From project root
make dev-web

# Or from this directory
bun run dev
```

## Building

```bash
bun run build
bun run preview
```

---

## Key Files

| File                        | Purpose                               |
| --------------------------- | ------------------------------------- |
| `vite.config.ts`            | TailwindCSS v4 Vite plugin            |
| `src/app.css`               | Tailwind + DaisyUI + NutriFresh theme |
| `src/app.html`              | HTML template with `data-theme`       |
| `src/routes/+layout.svelte` | App shell (Navbar)                    |
| `src/routes/+page.svelte`   | Homepage                              |
| `static/manifest.json`      | PWA manifest                          |
| `src/service-worker.ts`     | PWA offline support                   |

---

## PWA

The app is PWA-ready:

- `static/manifest.json` - App metadata
- `src/service-worker.ts` - Offline caching
- `src/app.html` - Mobile meta tags

**TODO:** Add icons `static/icon-192.png` and `static/icon-512.png`
