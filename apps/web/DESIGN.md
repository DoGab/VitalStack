# Design System: VitalStack

**Project ID:** 9071666792879869019

## 1. Visual Theme & Atmosphere

"Organic Premium" — a warm, grounded, nature-inspired wellness aesthetic. The UI feels like a premium health journal: calm, trustworthy, and sophisticated. Heavy emphasis on whitespace, soft shadows, and organic textures. The mood is relaxed luxury — Whole Foods meets a private health club.

Light mode uses warm cream paper backgrounds with deep arboretum green text. Dark mode inverts to a "Forest at Night" feel — deep greens with burnished gold accents.

The app is a mobile-first progressive web app. All components use shadcn-svelte (Bits UI library) with Tailwind CSS v4. Forms use shadcn Input, Select, and Button. Modals use shadcn Dialog (desktop) and Drawer (mobile). Cards use shadcn Card. All interactive elements must be 48px minimum tap targets.

## 2. Color Palette & Roles

### Light Theme (Organic Premium)

- **Deep Arboretum** (#1B3022) — Primary brand color. Used for main navigation, primary buttons, and headings. Communicates depth, nature, and stability.
- **Burnished Gold** (#C5A059) — Secondary color. Used for progress bars, AI suggestions, and premium feature highlights. Adds a touch of luxury.
- **Regal Plum** (#7B506F) — Accent color. Used sparingly for badges, special highlights, and call-to-action elements.
- **Cream Paper** (#F9F7F2) — Base background. A warm off-white that makes the app feel like a physical journal.
- **Warm Beige** (#EBE7DE) — Card backgrounds. Slightly darker than base for layered surfaces.
- **Soft Taupe** (#D8D3C8) — Input fields and borders.
- **Charcoal** (#2C2F2D) — Primary text. A soft, warm off-black for body text.
- **Terracotta** (#D65A31) — Calories indicator color. Warm orange for calorie-related displays.

### Macro-Specific Colors

- Calories: Terracotta (#D65A31)
- Protein: Warm Red (oklch 63.7% — red-500)
- Carbs: Golden Amber (oklch 76.9% — amber-500)
- Fat: Ocean Blue (oklch 62.3% — blue-500)

### Dark Theme (Dark Organic)

- **Burnished Gold** (#C5A059) — Primary. Premium CTAs with excellent contrast.
- **Deep Arboretum** (#1B3022) — Base background. The same anchor green as light theme.
- **Leaf Green** (#4E8056) — Secondary actions.
- **Warm Amber** (#D4913D) — Accent highlights.

## 3. Typography Rules

### Classic Theme (Default)

- **Display Font**: Playfair Display — Elegant serif for headlines, h1/h2. Conveys premium quality.
- **Body Font**: Inter — Clean sans-serif for all body text, buttons, labels, navigation.
- **Mono Font**: JetBrains Mono — Used exclusively for numbers, stats, and macro data. Provides visual authority to numerical values.

### Modern Theme (Alternative)

- **Display Font**: Outfit — Geometric, modern replacement for Playfair Display.
- Same Body and Mono fonts as Classic.

## 4. Component Stylings

- **Buttons**: Generously rounded corners (rounded-md, ~8px). Primary buttons use Deep Arboretum fill with cream text. Tactile feedback on press (active:scale-95). All buttons minimum 48px height for mobile.
- **Cards/Containers**: Subtly rounded corners (rounded-xl, ~12px). Light whisper-soft shadow (shadow-sm). Cream paper background in light mode. Border uses soft taupe (#D8D3C8).
- **Inputs/Forms**: Rounded corners (rounded-md). Soft taupe border with cream background. Focused state uses primary green ring. Number inputs use JetBrains Mono font.
- **Badges**: Pill-shaped (rounded-full). Light tinted backgrounds with matching text color. Used for confidence scores, serving sizes, timestamps.
- **Modals**: On mobile — bottom Drawer that slides up. On desktop — centered Dialog with max-width ~md (28rem). Both with smooth transitions.
- **Checkboxes**: Rounded (rounded-sm). Primary green check when selected. Used for ingredient selection.

## 5. Layout Principles

- **Mobile-first**: All layouts designed for 375px width first, then expanded for desktop.
- **Generous whitespace**: Minimum 16px padding on containers, 12px gaps between elements.
- **Bottom navigation**: Fixed bottom dock on mobile with 5 navigation items.
- **Overscroll-behavior: none**: For native app feel.
- **Content within 48rem (max-w-3xl)**: Centered with auto margins on desktop.
- **Ingredient lists**: Full-width items spanning the modal width for optimal mobile touch targets.

## 6. Design System Notes for Stitch Generation

Use these colors when generating any VitalStack screen:

- Background: Warm Cream (#F9F7F2)
- Card Surface: Warm Beige (#EBE7DE) or white
- Primary Button: Deep Arboretum (#1B3022)
- Primary Text: Charcoal (#2C2F2D)
- Secondary Text: Medium Gray (#6B7280)
- Accent Gold: Burnished Gold (#C5A059)
- Accent Plum: Regal Plum (#7B506F)
- Calories: Terracotta (#D65A31)
- Protein: Red (#EF4444)
- Carbs: Amber (#F59E0B)
- Fat: Blue (#3B82F6)

Visually mimic the shadcn-svelte component library. All buttons, inputs, dialogs, and cards should visually follow shadcn-svelte aesthetics (Tailwind styling, rounded-md, soft shadows). Do not attempt to output Svelte 5 logic; simply generate the precise visual layout matching the shadcn-svelte appearance. Your goal is to create a high-fidelity visual mockup. Inter font for body text. JetBrains Mono for numerical values.
