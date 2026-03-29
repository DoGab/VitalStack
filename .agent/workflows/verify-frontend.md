---
description: Verify frontend changes by starting the full stack and taking screenshots
---

## Steps

1. Check if the dev server is already running on port 5173
// turbo
2. If not running, start the full stack with `make dev` from the project root
3. Wait for the server to be ready by checking that `https://localhost:5173` responds
// turbo
4. Take a desktop screenshot (1280×800 viewport) at `https://localhost:5173` (navigate to the relevant page if needed)
// turbo
5. Take a mobile screenshot (375×812 viewport) at `https://localhost:5173` (same page)
6. Compare both screenshots and verify:
   - Layout is responsive (no overflow, no broken elements)
   - VitalStack design language is correct (Organic Premium / Dark Organic colors, fonts)
   - Custom components render correctly (CircularProgress, MacroBars, StatCard, SectionHeader)
   - Nutrition colors match `NUTRITION_CONFIG` (terracotta calories, red protein, amber carbs, blue fat, green fiber)
   - Mobile view has proper bottom navigation, thumb-friendly tap targets
7. Report findings — flag any visual issues or confirm success
