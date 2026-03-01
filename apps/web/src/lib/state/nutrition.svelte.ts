import type { components } from "$lib/api/schema";

type MacroData = components["schemas"]["MacroData"];

class NutritionState {
  dailyMacros = $state<MacroData | null>(null);

  // Optional helpers to return cleanly 0'd data if null
  get safeMacros(): MacroData {
    return (
      this.dailyMacros || {
        calories: 0,
        protein: 0,
        carbs: 0,
        fat: 0,
        fiber: 0
      }
    );
  }
}

export const nutritionState = new NutritionState();
