# Design System: VitalStack

**Project ID:** 9071666792879869019

## 1. Visual Theme & Atmosphere

"Organic Premium" — a warm, grounded, nature-inspired wellness aesthetic. The UI feels like a premium health journal: calm, trustworthy, and sophisticated. Heavy emphasis on whitespace, soft shadows, and organic textures. The mood is relaxed luxury — Whole Foods meets a private health club.

Light mode uses warm cream paper backgrounds with deep arboretum green text. Dark mode inverts to a "Forest at Night" feel — deep greens with burnished gold accents.

The app is a mobile-first progressive web app. All components use shadcn-svelte (Bits UI library) with Tailwind CSS v4. Forms use shadcn Input, Select, and Button. Modals use shadcn Dialog (desktop) and Drawer (mobile). Cards use shadcn Card. All interactive elements must be 48px minimum tap targets.

## 2. Color Palette & Roles

### Light Theme (Organic Premium)

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

### Macro-Specific Colors

- Calories: Terracotta (#D65A31)
- Protein: Warm Red (oklch 63.7% — red-500)
- Carbs: Golden Amber (oklch 76.9% — amber-500)
- Fat: Ocean Blue (oklch 62.3% — blue-500)
- Fiber: Lively Green (oklch 72.3% — green-500)

### Dark Theme (Dark Organic)

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
- Fiber: Green (#22C55E)

Visually mimic the shadcn-svelte component library. All buttons, inputs, dialogs, and cards should visually follow shadcn-svelte aesthetics (Tailwind styling, rounded-md, soft shadows). Do not attempt to output Svelte 5 logic; simply generate the precise visual layout matching the shadcn-svelte appearance. Your goal is to create a high-fidelity visual mockup. Inter font for body text. JetBrains Mono for numerical values.
