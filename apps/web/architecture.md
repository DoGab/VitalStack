# VitalStack Web Architecture

## Overview

The VitalStack web frontend is a **SvelteKit** application with a mobile-first, PWA-ready design. It uses TailwindCSS v4 for styling and **shadcn-svelte** (Bits UI) for UI components.

## Tech Stack

### Tooling

| Tool        | Purpose              | Notes                                        |
| ----------- | -------------------- | -------------------------------------------- |
| **pnpm**    | Package manager      | Faster installs, disk-efficient symlinks     |
| **Vite**    | Build tool / bundler | TailwindCSS v4 via `@tailwindcss/vite`       |
| **Node.js** | Runtime              | Dev server, SSR, production via adapter-node |

### Libraries

| Technology    | Version | Purpose                     |
| ------------- | ------- | --------------------------- |
| SvelteKit     | 2.x     | Full-stack framework        |
| Svelte        | 5.x     | UI framework (Runes mode)   |
| TailwindCSS   | 4.x     | Utility-first CSS           |
| shadcn-svelte | 1.x     | Component library (Bits UI) |
| Lucide Svelte | -       | Icon library                |
| adapter-node  | -       | Production Node.js adapter  |

---

## Directory Structure

```
apps/web/
├── src/
│   ├── app.css              # TailwindCSS + shadcn CSS variables + VitalStack theme
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
┌─────────────────────────────────────────────────────────────┐
│              Sidebar.Provider (+layout.svelte)              │
│ ┌──────────┐  ┌──────────────────────────────────────┐      │
│ │          │  │  Header (Sidebar.Trigger / Mobile)   │      │
│ │  App     │  ├──────────────────────────────────────┤      │
│ │  Sidebar │  │                                      │      │
│ │ (Desktop)│  │         Page Content                 │      │
│ │          │  │         (+page.svelte)               │      │
│ │          │  │                                      │      │
│ └──────────┘  └──────────────────────────────────────┘      │
│              ┌──────────────────────────────────────┐       │
│              │     Mobile Bottom Dock (lg:hidden)   │       │
│              └──────────────────────────────────────┘       │
└─────────────────────────────────────────────────────────────┘
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

## Design System: VitalStack

### Active Themes

| Theme           | Role  | Primary | Secondary | Accent  | Base    |
| --------------- | ----- | ------- | --------- | ------- | ------- |
| Organic Premium | Light | #1B3022 | #C5A059   | #7B506F | #F9F7F2 |
| Dark Organic    | Dark  | #C5A059 | #4E8056   | #D4913D | #1B3022 |

### Theme Configuration

Defined in `src/app.css` using shadcn CSS variable system with `@theme` (without `inline` for multi-theme support):

```css
@import "tailwindcss";
@import "tw-animate-css";
@custom-variant dark (&:is([data-theme="darkorganic"], [data-theme="darkorganic"] *));

:root,
[data-theme="organic"] {
  --primary: oklch(0.241 0.034 153.1); /* #1B3022 */
  --secondary: oklch(0.722 0.107 82.8); /* #C5A059 */
  --accent: oklch(0.49 0.073 337.5); /* #7B506F */
  --background: oklch(0.975 0.005 93.6); /* #F9F7F2 */
  /* ... */
}

[data-theme="darkorganic"] {
  /* dark overrides */
}

@theme {
  --color-primary: var(--primary);
  /* ... mapped to Tailwind */
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
  "name": "VitalStack",
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

## Color Scheme

Color tokens used by the shadcn-svelte component system:

| Category | Color Name | Usage                                                                                                                                 |
| -------- | ---------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| Brand    | Primary    | The main color of your app. Used for primary buttons, active links, and key focus points.                                             |
|          | Secondary  | A complementary color. Used for secondary actions or visual interest that shouldn't compete with the Primary.                         |
|          | Accent     | A "pop" color. Used sparingly for highlights, special badges, or call-to-action elements.                                             |
| State    | Info       | Blue/Cyan usually. Used for neutral alerts, help icons, or informational tooltips.                                                    |
|          | Success    | Green usually. Used for "Task Complete" toasts, confirm buttons, or positive trends.                                                  |
|          | Warning    | Yellow/Orange usually. Used for "Are you sure?" modals, pending states, or caution alerts.                                            |
|          | Error      | Red usually. Used for destructive actions (Delete), form validation errors, or critical bugs.                                         |
| Surface  | Neutral    | A dark, desaturated color (often near-black or deep gray). Used for text on light backgrounds or dark structural elements (sidebars). |
|          | Base-100   | The "floor" of your app. This is white (in light mode) or dark gray/black (in dark mode).                                             |
|          | Base-200   | Slightly darker/lighter than Base-100. Used for card backgrounds or distinct sections sitting on Base-100.                            |
|          | Base-300   | Even more distinct. Used for inputs, sidebars, or deeply nested elements.                                                             |

### Dark Organic (Default Dark Theme)

**Best for**: Premium wellness, warm dark mode, users who want a sophisticated dark experience without the cyberpunk aesthetic.
**Vibe**: "Forest at Night" — A private health club after hours.

This is the dark mode companion to the Organic Premium theme. It inverts the color scheme while maintaining the same warm, premium feel.

- Primary (#C5A059): Burnished Gold. Luxury CTAs with excellent contrast on dark backgrounds.
- Base (#1B3022): Deep Arboretum. The same anchor green as Organic Premium, now as the background.
- Secondary (#4E8056): Leaf Green. Natural, calming secondary actions.
- Accent (#D4913D): Warm Amber. A honey-toned highlight for badges and FABs.

| Category | Color Name   | Hex Code | Explanation                                                 |
| -------- | ------------ | -------- | ----------------------------------------------------------- |
| Brand    | Primary      | #C5A059  | Burnished Gold. Premium CTAs, active states.                |
|          | Secondary    | #4E8056  | Leaf Green. Natural secondary actions.                      |
|          | Accent       | #D4913D  | Warm Amber. Honey-toned highlights, warm and inviting.      |
| State    | Info         | #5F8D8B  | Sage Blue. Calm informational states.                       |
|          | Success      | #6A9E6E  | Soft Green. Natural success indication.                     |
|          | Warning      | #EDB654  | Harvest Yellow. Warm warning color.                         |
|          | Error        | #B93632  | Brick Red. Earthy, attention-grabbing.                      |
| Surface  | Neutral      | #8B9A8E  | Muted Sage. Secondary text, inactive icons.                 |
|          | Base-100     | #1B3022  | Deep Arboretum. Main background.                            |
|          | Base-200     | #243B2C  | Darker Forest. Card backgrounds.                            |
|          | Base-300     | #2D4A37  | Muted Evergreen. Input fields.                              |
|          | Base-Content | #F9F7F2  | Cream. Warm off-white text, easier on eyes than pure white. |

- **Pros**: Premium, warm, calming dark mode. Pairs perfectly with Organic Premium for light/dark toggle.
- **Cons**: Less aggressive/sporty than Bio-Hacker; may not appeal to gym-focused users.

---

### Suggestion 1 - Bio Hacker

**Best for**: A performance-focused audience tracking macros, gym stats, and hydration. **Vibe**: Nike Training Club meets Cyberpunk. High energy, modern, and distinctively "AI."

This palette uses your requested Neon Lime and Obsidian. This is designed primarily for a Dark Mode interface (which is very popular in fitness tracking apps to save battery during workouts).

- Primary (#D4FF00): This "Volt" color is synonymous with energy and electricity. It screams "Action."
- Base (#121212): Obsidian. We use a deep, rich black/gray to make the neon pop.
- Secondary (#6366F1): Indigo/Violet. This represents the "AI" intelligence aspect. It provides a cool temperature contrast to the aggressive lime.

| Category | Color Name | Hex Code | Explanation                                                                        |
| -------- | ---------- | -------- | ---------------------------------------------------------------------------------- |
| Brand    | Primary    | #D4FF00  | Volt Lime. High visibility. Note: Use black text on top of this for accessibility. |
|          | Secondary  | #6366F1  | AI Indigo. Represents the intelligence engine (scanning/processing).               |
|          | Accent     | #FFFFFF  | Pure White. Used for maximum contrast on data points or iconography.               |
| State    | Info       | #38BDF8  | Electric Sky. High-saturation blue for visibility against dark backgrounds.        |
|          | Success    | #22C55E  | Bright Green. A distinct green that differs enough from the Lime Primary.          |
|          | Warning    | #FACC15  | Hazard Yellow. Standard warning color, highly legible on dark mode.                |
|          | Error      | #EF4444  | Red Alert. High-saturation red for immediate attention.                            |
| Surface  | Neutral    | #A1A1AA  | Zinc Grey. Used for secondary text or inactive icons.                              |
|          | Base-100   | #121212  | Obsidian. The deep background foundation.                                          |
|          | Base-200   | #1E1E1E  | Graphite. Slightly lighter cards floating on the obsidian background.              |
|          | Base-300   | #2D2D2D  | Jet. Input fields and sidebar backgrounds.                                         |

- **Pros**: Extremely trendy, high energy, looks great on OLED screens.
- **Cons**: The neon lime can be fatiguing if overused; requires strict accessibility checks.

### Suggestion 2 - Organic Premium

**Best for**: A lifestyle-focused brand that wants to feel established, trustworthy, and high-end.
**Vibe**: Whole Foods meets a private health club. It feels grounded and serious about longevity.

This palette uses your requested Deep Arboretum and Burnished Gold. It avoids the "clinical" look of hospitals and instead feels like a premium wellness journal.

- Primary (#1B3022): A deep, forest green. It communicates deep health, vegetables, and nature, but darker—implying stability and seriousness.
- Secondary (#C5A059): Used for premium features (like "AI Insights" or "Goal Reached"). It adds a touch of luxury.
- Accent (#7B506F): Plum. A sophisticated, premium pop color for badges and highlights.

| Category | Color Name | Hex Code | Explanation                                                                             |
| -------- | ---------- | -------- | --------------------------------------------------------------------------------------- |
| Brand    | Primary    | #1B3022  | Deep Arboretum. The core brand anchor. Used for main navigation and ""Scan"" buttons.   |
|          | Secondary  | #C5A059  | Burnished Gold. Adds a premium feel to progress bars and AI suggestions.                |
|          | Accent     | #7B506F  | Plum. A sophisticated, premium pop color for badges and highlights.                     |
| State    | Info       | #5F8D8B  | Sage Blue. A muted, dusty teal that provides information without looking too ""techy."" |
|          | Success    | #4E8056  | Leaf Green. Natural and reassuring for ""Meal Logged"" messages.                        |
|          | Warning    | #EDB654  | Harvest Yellow. A soft yellow that warns without inducing panic.                        |
|          | Error      | #B93632  | Brick Red. Earthy rather than neon; used for ""Unknown Barcode"" or critical errors.    |
| Surface  | Neutral    | #2C2F2D  | Charcoal. A soft, warm off-black for text. Less harsh than pure black.                  |
|          | Base-100   | #F9F7F2  | Cream Paper. A warm off-white. Makes the app feel like a physical journal.              |
|          | Base-200   | #EBE7DE  | Beige. Slightly darker for card backgrounds.                                            |
|          | Base-300   | #D8D3C8  | Taupe. Used for inputs and borders.                                                     |

- **Pros**: Feels premium, natural, and calming. The gold adds a touch of luxury without being gaudy.
- **Cons**: The "Cream Paper" base might feel slightly off-white to users expecting pure white.

### Suggestion 3 - Clinical Modernist

**Best for**: Mass market appeal. It feels clean, medical (but not sterile), and very easy to read.
**Vibe**: Clean, scientific, and fresh.

It uses a Teal/Mint approach which is the industry standard for "Digital Health."

- Primary (#0D9488): Teal. It combines the trust of blue (tech) with the renewal of green (health).
- Secondary (#F472B6): Soft Pink/Grapefruit. A friendly complementary color for "Food" and "Humanity."
- Base (#FFFFFF): Pure, clean white for a sterile, accurate data feel.

| Category | Color Name | Hex Code | Explanation                                                         |
| -------- | ---------- | -------- | ------------------------------------------------------------------- |
| Brand    | Primary    | #0D9488  | Vital Teal. Clinical yet fresh. High trust factor.                  |
|          | Secondary  | #F472B6  | Grapefruit. Adds a human, friendly touch to the tech interface.     |
|          | Accent     | #0F172A  | Midnight. Deep navy for strong Call-to-Actions (CTAs) and contrast. |
| State    | Info       | #60A5FA  | Soft Blue. Friendly informational states.                           |
|          | Success    | #34D399  | Mint. Fresh and clean success states.                               |
|          | Warning    | #FBBF24  | Amber. Warm warning color.                                          |
|          | Error      | #F87171  | Soft Red. Clear error indication without being aggressive.          |
| Surface  | Neutral    | #334155  | Slate. A cool-toned dark grey for text, sharper than standard grey. |
|          | Base-100   | #FFFFFF  | Pure White. The cleanest possible canvas for data visualization.    |
|          | Base-200   | #F1F5F9  | Cool Grey. Very subtle separation for cards.                        |
|          | Base-300   | #E2E8F0  | Cloud. Borders and dividers.                                        |

- **Pros**: Safe, highly legible, familiar to users of apps like MyFitnessPal or Apple Health.
- **Cons**: Less unique; might blend in with competitors.

## Fonts

VitalStack supports two font themes that can be switched at runtime via the font theme toggle in the navigation bar.

### Font Theme: Classic (Default)

Elegant, premium feel with a serif display font matching the logo.

| Role    | Font                 | Usage                                  |
| ------- | -------------------- | -------------------------------------- |
| Display | **Playfair Display** | Headlines, hero sections, h1/h2        |
| Body    | **Inter**            | Body text, buttons, labels, navigation |
| Mono    | **JetBrains Mono**   | Macro numbers, stats, data             |

### Font Theme: Modern

Clean, geometric, modern tech aesthetic.

| Role    | Font               | Usage                                  |
| ------- | ------------------ | -------------------------------------- |
| Display | **Outfit**         | Headlines, hero sections, h1/h2        |
| Body    | **Inter**          | Body text, buttons, labels, navigation |
| Mono    | **JetBrains Mono** | Macro numbers, stats, data             |

### CSS Custom Properties

```css
:root {
  --font-display: "Playfair Display", Georgia, serif;
  --font-body: "Inter", system-ui, sans-serif;
  --font-mono: "JetBrains Mono", ui-monospace, monospace;
}

[data-font-theme="modern"] {
  --font-display: "Outfit", system-ui, sans-serif;
}
```

### Usage

- `font-display` class for headlines
- `font-mono` class for numbers/stats
- Body text uses `--font-body` by default

### Alternative Pairings (Future Consideration)

| Style          | Display          | Body    | Best For          |
| -------------- | ---------------- | ------- | ----------------- |
| Health/Organic | Playfair Display | Lato    | Wellness brands   |
| Minimalist     | Manrope          | Manrope | Tech-forward apps |
| Friendly       | Poppins          | Inter   | Consumer apps     |
