/**
 * Centralized nutrition configuration for consistent icons and colors
 * across the application (summary cards, ingredient breakdown, etc.)
 */
import { Flame, Drumstick, Wheat, Droplet, Apple } from "lucide-svelte";
import type { ComponentType } from "svelte";

export interface NutritionConfig {
  icon: ComponentType;
  color: string;
  label: string;
  unit: string;
}

export const NUTRITION_CONFIG: Record<string, NutritionConfig> = {
  calories: { icon: Flame, color: "text-orange-500", label: "Calories", unit: "kcal" },
  protein: { icon: Drumstick, color: "text-red-500", label: "Protein", unit: "g" },
  carbs: { icon: Wheat, color: "text-amber-500", label: "Carbs", unit: "g" },
  fat: { icon: Droplet, color: "text-blue-500", label: "Fat", unit: "g" },
  fiber: { icon: Apple, color: "text-green-500", label: "Fiber", unit: "g" }
};

/**
 * Get ordered macro items for display
 */
export function getMacroDisplayOrder(): Array<keyof typeof NUTRITION_CONFIG> {
  return ["calories", "protein", "carbs", "fat", "fiber"];
}
