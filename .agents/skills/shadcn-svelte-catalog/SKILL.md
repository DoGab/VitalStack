---
name: shadcn-svelte-catalog
description: |
  Complete catalog of all available shadcn-svelte components for installation.
  Use when building frontend features to discover components that could help.
  Includes installation commands, category organization, and documentation links.
  Keywords: shadcn-svelte, component discovery, install, catalog, bits-ui, available components
---

# shadcn-svelte Component Catalog

> **Source:** [shadcn-svelte.com/docs/components](https://www.shadcn-svelte.com/docs/components)
> **LLM Docs:** [shadcn-svelte.com/llms.txt](https://www.shadcn-svelte.com/llms.txt)

## When to Use This Skill

Before implementing any frontend feature, consult this catalog to:
1. Identify components that could be used but aren't yet installed
2. Avoid reinventing UI patterns that already exist in shadcn-svelte
3. Find the right component for a specific interaction pattern

## Component Discovery Workflow

1. **Check installed:** `ls apps/web/src/lib/components/ui/`
2. **Cross-reference** against this catalog to find relevant uninstalled components
3. **Read component docs** before installing (links below)
4. **Install:** `pnpm dlx shadcn-svelte@latest add <component>` (run from `apps/web/`)

## Installation

```bash
# From the apps/web directory
cd apps/web
pnpm dlx shadcn-svelte@latest add <component-name>

# Examples:
pnpm dlx shadcn-svelte@latest add tabs
pnpm dlx shadcn-svelte@latest add sonner
pnpm dlx shadcn-svelte@latest add switch
```

---

## Full Component Catalog

### Form & Input

| Component | Description | Docs |
|-----------|-------------|------|
| **Button** | Displays a button or a component that looks like a button. | [docs](https://shadcn-svelte.com/docs/components/button.md) |
| **Button Group** | Container that groups related buttons together with consistent styling. | [docs](https://shadcn-svelte.com/docs/components/button-group.md) |
| **Calendar** | Calendar component that allows users to select dates. | [docs](https://shadcn-svelte.com/docs/components/calendar.md) |
| **Checkbox** | A control that allows the user to toggle between checked and not checked. | [docs](https://shadcn-svelte.com/docs/components/checkbox.md) |
| **Combobox** | Autocomplete input and command palette with a list of suggestions. | [docs](https://shadcn-svelte.com/docs/components/combobox.md) |
| **Date Picker** | A date picker component with range and presets. | [docs](https://shadcn-svelte.com/docs/components/date-picker.md) |
| **Field** | Combine labels, controls, and help text to compose accessible form fields. | [docs](https://shadcn-svelte.com/docs/components/field.md) |
| **Formsnap** | Building forms with Formsnap, Superforms, & Zod. | [docs](https://shadcn-svelte.com/docs/components/form.md) |
| **Input** | Displays a form input field or a component that looks like an input field. | [docs](https://shadcn-svelte.com/docs/components/input.md) |
| **Input Group** | Display additional information or actions to an input or textarea. | [docs](https://shadcn-svelte.com/docs/components/input-group.md) |
| **Input OTP** | Accessible one-time password component with copy paste functionality. | [docs](https://shadcn-svelte.com/docs/components/input-otp.md) |
| **Label** | Renders an accessible label associated with controls. | [docs](https://shadcn-svelte.com/docs/components/label.md) |
| **Native Select** | Styled native HTML select element with consistent design. | [docs](https://shadcn-svelte.com/docs/components/native-select.md) |
| **Radio Group** | A set of radio buttons where no more than one can be checked at a time. | [docs](https://shadcn-svelte.com/docs/components/radio-group.md) |
| **Select** | Displays a list of options for the user to pick from—triggered by a button. | [docs](https://shadcn-svelte.com/docs/components/select.md) |
| **Slider** | An input where the user selects a value from within a given range. | [docs](https://shadcn-svelte.com/docs/components/slider.md) |
| **Switch** | A control that allows the user to toggle between checked and not checked. | [docs](https://shadcn-svelte.com/docs/components/switch.md) |
| **Textarea** | Displays a form textarea or a component that looks like a textarea. | [docs](https://shadcn-svelte.com/docs/components/textarea.md) |

### Layout & Navigation

| Component | Description | Docs |
|-----------|-------------|------|
| **Accordion** | Vertically stacked interactive headings that reveal sections of content. | [docs](https://shadcn-svelte.com/docs/components/accordion.md) |
| **Breadcrumb** | Displays the path to the current resource using a hierarchy of links. | [docs](https://shadcn-svelte.com/docs/components/breadcrumb.md) |
| **Navigation Menu** | A collection of links for navigating websites. | [docs](https://shadcn-svelte.com/docs/components/navigation-menu.md) |
| **Resizable** | Accessible resizable panel groups and layouts with keyboard support. | [docs](https://shadcn-svelte.com/docs/components/resizable.md) |
| **Scroll Area** | Augments native scroll for custom, cross-browser styling. | [docs](https://shadcn-svelte.com/docs/components/scroll-area.md) |
| **Separator** | Visually or semantically separates content. | [docs](https://shadcn-svelte.com/docs/components/separator.md) |
| **Sidebar** | Composable, themeable, and customizable sidebar component. | [docs](https://shadcn-svelte.com/docs/components/sidebar.md) |
| **Tabs** | Layered sections of content displayed one at a time. | [docs](https://shadcn-svelte.com/docs/components/tabs.md) |

### Overlays & Dialogs

| Component | Description | Docs |
|-----------|-------------|------|
| **Alert Dialog** | Modal dialog that interrupts the user with important content. | [docs](https://shadcn-svelte.com/docs/components/alert-dialog.md) |
| **Command** | Fast, composable, unstyled command menu for Svelte. | [docs](https://shadcn-svelte.com/docs/components/command.md) |
| **Context Menu** | Menu triggered by right click. | [docs](https://shadcn-svelte.com/docs/components/context-menu.md) |
| **Dialog** | Window overlaid on the primary window, rendering content underneath inert. | [docs](https://shadcn-svelte.com/docs/components/dialog.md) |
| **Drawer** | A drawer component for Svelte. | [docs](https://shadcn-svelte.com/docs/components/drawer.md) |
| **Dropdown Menu** | Menu triggered by a button, for actions or functions. | [docs](https://shadcn-svelte.com/docs/components/dropdown-menu.md) |
| **Hover Card** | Preview content available behind a link (sighted users). | [docs](https://shadcn-svelte.com/docs/components/hover-card.md) |
| **Menubar** | Persistent menu for desktop apps with quick access to commands. | [docs](https://shadcn-svelte.com/docs/components/menubar.md) |
| **Popover** | Displays rich content in a portal, triggered by a button. | [docs](https://shadcn-svelte.com/docs/components/popover.md) |
| **Sheet** | Extends Dialog to display complementary content from the side. | [docs](https://shadcn-svelte.com/docs/components/sheet.md) |
| **Tooltip** | Popup showing info on hover or keyboard focus. | [docs](https://shadcn-svelte.com/docs/components/tooltip.md) |

### Feedback & Status

| Component | Description | Docs |
|-----------|-------------|------|
| **Alert** | Displays a callout for user attention. | [docs](https://shadcn-svelte.com/docs/components/alert.md) |
| **Badge** | Displays a badge or component that looks like a badge. | [docs](https://shadcn-svelte.com/docs/components/badge.md) |
| **Empty** | Empty state component to indicate when no data is present. | [docs](https://shadcn-svelte.com/docs/components/empty.md) |
| **Progress** | Indicator showing task completion progress as a bar. | [docs](https://shadcn-svelte.com/docs/components/progress.md) |
| **Skeleton** | Placeholder shown while content is loading. | [docs](https://shadcn-svelte.com/docs/components/skeleton.md) |
| **Sonner** | An opinionated toast notification component for Svelte. | [docs](https://shadcn-svelte.com/docs/components/sonner.md) |
| **Spinner** | Indicator for a loading state. | [docs](https://shadcn-svelte.com/docs/components/spinner.md) |

### Display & Media

| Component | Description | Docs |
|-----------|-------------|------|
| **Aspect Ratio** | Displays content within a desired ratio. | [docs](https://shadcn-svelte.com/docs/components/aspect-ratio.md) |
| **Avatar** | Image element with fallback for representing the user. | [docs](https://shadcn-svelte.com/docs/components/avatar.md) |
| **Card** | Displays a card with header, content, and footer. | [docs](https://shadcn-svelte.com/docs/components/card.md) |
| **Carousel** | Carousel with motion and swipe built using Embla. | [docs](https://shadcn-svelte.com/docs/components/carousel.md) |
| **Chart** | Charts built using LayerChart. Copy and paste into apps. | [docs](https://shadcn-svelte.com/docs/components/chart.md) |
| **Data Table** | Powerful table and datagrids built using TanStack Table. | [docs](https://shadcn-svelte.com/docs/components/data-table.md) |
| **Item** | Versatile component for displaying any content. | [docs](https://shadcn-svelte.com/docs/components/item.md) |
| **Kbd** | Displays textual user input from keyboard. | [docs](https://shadcn-svelte.com/docs/components/kbd.md) |
| **Table** | A responsive table component. | [docs](https://shadcn-svelte.com/docs/components/table.md) |
| **Typography** | Styles for headings, paragraphs, lists, etc. | [docs](https://shadcn-svelte.com/docs/components/typography.md) |

### Misc

| Component | Description | Docs |
|-----------|-------------|------|
| **Collapsible** | Interactive component that expands/collapses a panel. | [docs](https://shadcn-svelte.com/docs/components/collapsible.md) |
| **Pagination** | Page navigation with next and previous links. | [docs](https://shadcn-svelte.com/docs/components/pagination.md) |
| **Range Calendar** | Calendar component for selecting a range of dates. | [docs](https://shadcn-svelte.com/docs/components/range-calendar.md) |
| **Toggle** | Two-state button that can be on or off. | [docs](https://shadcn-svelte.com/docs/components/toggle.md) |
| **Toggle Group** | Set of two-state buttons that can be toggled on or off. | [docs](https://shadcn-svelte.com/docs/components/toggle-group.md) |

---

## Quick Reference: Import Patterns

```svelte
<!-- Namespace import for multi-part components -->
<script lang="ts">
  import * as Card from '$lib/components/ui/card';
  import * as Dialog from '$lib/components/ui/dialog';
  import * as Tabs from '$lib/components/ui/tabs';
</script>

<!-- Direct import for single-export components -->
<script lang="ts">
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Badge } from '$lib/components/ui/badge';
</script>
```

## Looking Up Component API Details

For detailed props, events, and usage examples for any component, fetch the LLM-friendly docs:
```
https://shadcn-svelte.com/docs/components/<component-name>.md
```

Example: `https://shadcn-svelte.com/docs/components/tabs.md`
