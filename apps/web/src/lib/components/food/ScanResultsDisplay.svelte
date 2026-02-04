<script lang="ts">
  import { Scale, UtensilsCrossed } from "lucide-svelte";
  import type { components } from "$lib/api/schema";
  import { NUTRITION_CONFIG, getMacroDisplayOrder } from "$lib/config/nutrition-config";

  type ScanResult = components["schemas"]["ScanOutputBody"];

  let { result }: { result: ScanResult } = $props();

  // Build macro items from centralized config
  const macroItems = $derived(() => {
    const macroValues: Record<string, number> = {
      calories: result.macros.calories,
      protein: result.macros.protein,
      carbs: result.macros.carbs,
      fat: result.macros.fat,
      fiber: result.macros.fiber
    };

    return getMacroDisplayOrder().map((key) => ({
      ...NUTRITION_CONFIG[key],
      value: macroValues[key],
      key
    }));
  });

  // Format compact macro display for ingredients
  function formatMacro(value: number): string {
    return value.toFixed(0);
  }

  // Get ingredient macro items for inline display
  function getIngredientMacros(macros: {
    calories: number;
    protein: number;
    carbs: number;
    fat: number;
  }) {
    return [
      { key: "calories", ...NUTRITION_CONFIG.calories, value: macros.calories },
      { key: "protein", ...NUTRITION_CONFIG.protein, value: macros.protein },
      { key: "carbs", ...NUTRITION_CONFIG.carbs, value: macros.carbs },
      { key: "fat", ...NUTRITION_CONFIG.fat, value: macros.fat }
    ];
  }
</script>

<div class="space-y-4">
  <!-- Food Name & Confidence -->
  <div class="text-center">
    <h3 class="text-xl font-bold text-base-content">{result.food_name}</h3>
    <div class="flex items-center justify-center gap-2 mt-2">
      <span class="badge badge-success gap-1.5 px-3 py-3 text-sm font-medium">
        ✓ {Math.round(result.confidence * 100)}% confident
      </span>
    </div>
  </div>

  <!-- Serving Size -->
  <div class="flex items-center justify-center gap-2 text-base-content/70">
    <Scale class="w-4 h-4" />
    <span class="text-sm">{result.serving_size}</span>
  </div>

  <!-- Macro Grid -->
  <div class="grid grid-cols-2 gap-3">
    {#each macroItems() as macro, i (macro.key)}
      {@const Icon = macro.icon}
      <div class="stat bg-base-200 rounded-xl p-3 {i === 0 ? 'col-span-2' : ''}">
        <div class="stat-figure {macro.color}">
          <Icon class="w-6 h-6" />
        </div>
        <div class="stat-title text-xs">{macro.label}</div>
        <div class="stat-value text-lg {macro.color}">
          {typeof macro.value === "number" ? macro.value.toFixed(0) : macro.value}
          <span class="text-xs font-normal text-base-content/50">{macro.unit}</span>
        </div>
      </div>
    {/each}
  </div>

  <!-- Collapsible Ingredient Breakdown -->
  {#if result.ingredients && result.ingredients.length > 0}
    <div class="collapse collapse-arrow bg-base-200 rounded-xl">
      <input type="checkbox" />
      <div class="collapse-title font-medium flex items-center gap-2">
        <UtensilsCrossed class="w-4 h-4" />
        Ingredient Breakdown ({result.ingredients.length} items)
      </div>
      <div class="collapse-content">
        <div class="space-y-3 pt-2">
          {#each result.ingredients as ingredient (ingredient.name)}
            <div class="bg-base-100 rounded-lg p-3">
              <div class="flex justify-between items-start">
                <div class="flex items-center gap-2">
                  <UtensilsCrossed class="w-4 h-4 text-base-content/50" />
                  <span class="font-medium text-sm">{ingredient.name}</span>
                </div>
                <span class="text-xs text-base-content/60 font-medium"
                  >{ingredient.weight_grams}g</span
                >
              </div>
              <div class="flex gap-3 mt-2 text-xs">
                {#each getIngredientMacros(ingredient.macros) as m (m.key)}
                  {@const Icon = m.icon}
                  <span class="flex items-center gap-1 {m.color}">
                    <Icon class="w-3 h-3" />
                    {formatMacro(m.value)}{m.unit !== "kcal" ? "g" : ""}
                  </span>
                {/each}
              </div>
            </div>
          {/each}
        </div>
      </div>
    </div>
  {/if}
</div>
