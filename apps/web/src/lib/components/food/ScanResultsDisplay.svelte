<script lang="ts">
  import { Scale, UtensilsCrossed, ChevronDown } from "lucide-svelte";
  import type { components } from "$lib/api/schema";
  import { NUTRITION_CONFIG, getMacroDisplayOrder } from "$lib/config/nutrition-config";
  import { Badge } from "$lib/components/ui/badge";
  import * as Card from "$lib/components/ui/card";
  import * as Collapsible from "$lib/components/ui/collapsible";
  import { Separator } from "$lib/components/ui/separator";
  import { Button } from "$lib/components/ui/button";

  type ScanResult = components["schemas"]["ScanOutputBody"];

  let { result }: { result: ScanResult } = $props();

  let ingredientsOpen = $state(false);

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
    <h3 class="text-xl font-bold text-foreground" style="font-family: var(--font-heading);">
      {result.food_name}
    </h3>
    <div class="flex items-center justify-center gap-2 mt-2">
      <Badge variant="secondary" class="gap-1.5 px-3 py-1 text-sm">
        ✓ {Math.round(result.confidence * 100)}% confident
      </Badge>
    </div>
  </div>

  <!-- Serving Size -->
  <div class="flex items-center justify-center gap-2 text-muted-foreground">
    <Scale class="h-4 w-4" />
    <span class="text-sm">{result.serving_size}</span>
  </div>

  <Separator />

  <!-- Macro Grid -->
  <div class="grid grid-cols-2 gap-3">
    {#each macroItems() as macro, i (macro.key)}
      {@const Icon = macro.icon}
      <Card.Root class={i === 0 ? "col-span-2" : ""}>
        <Card.Content class="flex items-center justify-between p-3">
          <div>
            <p class="text-xs text-muted-foreground">{macro.label}</p>
            <p class="text-lg font-bold {macro.color}">
              {typeof macro.value === "number" ? macro.value.toFixed(0) : macro.value}
              <span class="text-xs font-normal text-muted-foreground">{macro.unit}</span>
            </p>
          </div>
          <div class={macro.color}>
            <Icon class="h-6 w-6" />
          </div>
        </Card.Content>
      </Card.Root>
    {/each}
  </div>

  <!-- Collapsible Ingredient Breakdown -->
  {#if result.ingredients && result.ingredients.length > 0}
    <Collapsible.Root bind:open={ingredientsOpen}>
      <Collapsible.Trigger>
        {#snippet child({ props })}
          <Button {...props} variant="ghost" class="w-full justify-between px-4 py-3 h-auto">
            <span class="flex items-center gap-2 font-medium">
              <UtensilsCrossed class="h-4 w-4" />
              Ingredient Breakdown ({result.ingredients!.length} items)
            </span>
            <ChevronDown
              class="h-4 w-4 transition-transform {ingredientsOpen ? 'rotate-180' : ''}"
            />
          </Button>
        {/snippet}
      </Collapsible.Trigger>
      <Collapsible.Content>
        <div class="space-y-3 px-4 pb-4">
          {#each result.ingredients as ingredient (ingredient.name)}
            <Card.Root>
              <Card.Content class="p-3">
                <div class="flex justify-between items-start">
                  <div class="flex items-center gap-2">
                    <UtensilsCrossed class="h-4 w-4 text-muted-foreground" />
                    <span class="font-medium text-sm">{ingredient.name}</span>
                  </div>
                  <span class="text-xs text-muted-foreground font-medium">
                    {ingredient.weight_grams}g
                  </span>
                </div>
                <div class="flex gap-3 mt-2 text-xs">
                  {#each getIngredientMacros(ingredient.macros) as m (m.key)}
                    {@const Icon = m.icon}
                    <span class="flex items-center gap-1 {m.color}">
                      <Icon class="h-3 w-3" />
                      {formatMacro(m.value)}{m.unit !== "kcal" ? "g" : ""}
                    </span>
                  {/each}
                </div>
              </Card.Content>
            </Card.Root>
          {/each}
        </div>
      </Collapsible.Content>
    </Collapsible.Root>
  {/if}
</div>
