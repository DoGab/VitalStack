/**
 * Centralized nutrition configuration for consistent icons and colors
 * across the application (summary cards, ingredient breakdown, etc.)
 */
import { Flame, Drumstick, Wheat, Droplet, Apple } from "lucide-svelte";
import type { ComponentType } from "svelte";

export interface NutritionConfig {
  icon: ComponentType;
  color: string;
  barColor: string; // CSS color value for progress bar fills
  label: string;
  unit: string;
}

export const NUTRITION_CONFIG: Record<string, NutritionConfig> = {
  calories: {
    icon: Flame,
    color: "text-[#D65A31]",
    barColor: "#D65A31",
    label: "Calories",
    unit: "kcal"
  },
  protein: {
    icon: Drumstick,
    color: "text-red-500",
    barColor: "oklch(63.7% 0.237 25.331)",
    label: "Protein",
    unit: "g"
  },
  carbs: {
    icon: Wheat,
    color: "text-amber-500",
    barColor: "oklch(76.9% 0.188 70.08)",
    label: "Carbs",
    unit: "g"
  },
  fat: {
    icon: Droplet,
    color: "text-blue-500",
    barColor: "oklch(62.3% 0.214 259.815)",
    label: "Fat",
    unit: "g"
  },
  fiber: {
    icon: Apple,
    color: "text-green-500",
    barColor: "oklch(72.3% 0.219 149.579)",
    label: "Fiber",
    unit: "g"
  }
};

/**
 * Get ordered macro items for display
 */
export function getMacroDisplayOrder(): Array<keyof typeof NUTRITION_CONFIG> {
  return ["calories", "protein", "carbs", "fat", "fiber"];
}
