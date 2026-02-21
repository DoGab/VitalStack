# Archive

This file is used to store old design considerations and ideas that are no longer in use but might be useful in the future again.

## Design Considerations

### Color Scheme - Bio Hacker

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

### Color Scheme - Clinical Modernist

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

### Alternative Font themes

| Style          | Display          | Body    | Best For          |
| -------------- | ---------------- | ------- | ----------------- |
| Health/Organic | Playfair Display | Lato    | Wellness brands   |
| Minimalist     | Manrope          | Manrope | Tech-forward apps |
| Friendly       | Poppins          | Inter   | Consumer apps     |

## Shadcn UI replaces Daisy UI

DaisyUI was replaced with shadcn-svelte (Bits UI) as the component library. Key decision drivers:

- **Premium over playful** — DaisyUI's opinionated styling produces a playful, consumer-app look. shadcn-svelte provides unstyled, composable primitives that allow full control over the VitalStack "Organic Premium" aesthetic
- **CSS variable theming** — shadcn-svelte's architecture maps directly to CSS custom properties (`--primary`, `--accent`, etc.), making multi-theme support (light/dark, font themes) trivial without fighting the framework
- **Accessibility built-in** — Bits UI primitives (Dialog, Popover, Drawer, etc.) ship with proper ARIA attributes, focus management, and keyboard navigation out of the box
- **Headless architecture** — Components are unstyled by default and composed via Tailwind classes, eliminating style override battles that were common with DaisyUI's pre-styled components
- **Svelte 5 native** — shadcn-svelte is built for Svelte 5 Runes, while DaisyUI is framework-agnostic CSS that doesn't leverage Svelte's reactivity model
- **Ecosystem alignment** — shadcn-svelte is the de facto standard for production SvelteKit apps, with better long-term maintenance and community support
