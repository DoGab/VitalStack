<script lang="ts">
  /**
   * IngredientServingDrawer — 3-mode serving size picker for ingredients.
   *
   * Reuses the same visual pattern as ProductLogDrawer:
   *  - Per 100g (derived from detected macros)
   *  - Serving (original AI-detected quantity)
   *  - Custom (manual weight entry)
   *
   * Works for both AI-detected ingredients and product-based ingredients.
   */
  import * as Drawer from "$lib/components/ui/drawer";
  import * as Dialog from "$lib/components/ui/dialog";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Badge } from "$lib/components/ui/badge";
  import MacroDots from "$lib/components/ui/macro-dots.svelte";
  import { IsMobile } from "$lib/hooks/is-mobile.svelte.js";
  import Scale from "lucide-svelte/icons/scale";
  import type { EditableIngredient } from "$lib/components/dashboard/IngredientEditor.svelte";

  interface Props {
    open: boolean;
    ingredient: EditableIngredient | null;
    /** Called when the user confirms a new serving size */
    onupdate?: (ingredient: EditableIngredient) => void;
  }

  let { open = $bindable(), ingredient, onupdate }: Props = $props();

  const isMobile = new IsMobile();

  // Serving mode
  type ServingMode = "100g" | "serving" | "custom";
  let servingMode: ServingMode = $state("serving");
  let customWeight = $state(100);

  // Compute per-100g base macros from the ingredient's detected amount.
  // e.g., if the AI says "40g oats = 150 cal", then per 100g = (150/40)*100 = 375 cal
  let per100gMacros = $derived.by(() => {
    if (!ingredient) return { calories: 0, protein: 0, carbs: 0, fat: 0, fiber: 0 };
    const detectedGrams = (ingredient.serving_quantity ?? 1) * (ingredient.serving_size ?? 100);
    const factor = 100 / detectedGrams;
    return {
      calories: ingredient.base_macros.calories * factor,
      protein: ingredient.base_macros.protein * factor,
      carbs: ingredient.base_macros.carbs * factor,
      fat: ingredient.base_macros.fat * factor,
      fiber: ingredient.base_macros.fiber * factor
    };
  });

  // Detected serving in grams
  let detectedGrams = $derived.by(() => {
    if (!ingredient) return 100;
    return (ingredient.serving_quantity ?? 1) * (ingredient.serving_size ?? 100);
  });

  // Selected weight based on mode
  let selectedWeight = $derived.by(() => {
    if (servingMode === "100g") return 100;
    if (servingMode === "serving") return detectedGrams;
    return customWeight;
  });

  // Scaled macros based on selected weight
  let scaledMacros = $derived.by(() => {
    const factor = selectedWeight / 100;
    return {
      calories: Math.round(per100gMacros.calories * factor),
      protein: +(per100gMacros.protein * factor).toFixed(1),
      carbs: +(per100gMacros.carbs * factor).toFixed(1),
      fat: +(per100gMacros.fat * factor).toFixed(1),
      fiber: +(per100gMacros.fiber * factor).toFixed(1)
    };
  });

  // Serving label for display
  let servingLabel = $derived.by(() => {
    if (servingMode === "100g") return "100g";
    if (servingMode === "serving") {
      const unit = ingredient?.serving_unit ?? "g";
      return `${detectedGrams}${unit}`;
    }
    return `${customWeight}g`;
  });

  // Original serving label
  let originalServingLabel = $derived.by(() => {
    if (!ingredient) return "Serving";
    const qty = ingredient.serving_quantity ?? 1;
    const size = ingredient.serving_size ?? "";
    const unit = ingredient.serving_unit ?? "g";
    if (size && unit) return `${qty} × ${size}${unit}`;
    return `${qty} serving`;
  });

  // Reset when drawer opens with a new ingredient
  $effect(() => {
    if (open && ingredient) {
      servingMode = "serving";
      customWeight = detectedGrams;
    }
  });

  function handleUpdate() {
    if (!ingredient) return;
    // Compute the new serving quantity based on selected weight
    const newQuantity = selectedWeight / (ingredient.serving_size ?? 100);
    const updated: EditableIngredient = {
      ...ingredient,
      macros: scaledMacros,
      serving_quantity: newQuantity
    };
    onupdate?.(updated);
    open = false;
  }
</script>

{#snippet body()}
  {#if ingredient}
    <div class="px-4 pb-4 space-y-4">
      <!-- Ingredient Header -->
      <div class="flex items-center gap-3">
        <div class="min-w-0">
          <p class="font-semibold text-foreground truncate">{ingredient.name}</p>
          <MacroDots macros={scaledMacros} class="mt-1" />
        </div>
      </div>

      <!-- Serving Size Selection -->
      <div class="space-y-2">
        <p class="text-sm font-medium text-foreground">Serving Size</p>
        <div class="grid grid-cols-3 gap-2">
          <!-- 100g -->
          <button
            type="button"
            class="rounded-lg border px-3 py-2 text-sm font-medium transition-colors
              {servingMode === '100g'
              ? 'border-primary bg-primary/10 text-primary'
              : 'border-border text-foreground hover:bg-accent/50'}"
            onclick={() => (servingMode = "100g")}
          >
            100g
          </button>

          <!-- Original serving -->
          <button
            type="button"
            class="rounded-lg border px-3 py-2 text-sm font-medium transition-colors truncate
              {servingMode === 'serving'
              ? 'border-primary bg-primary/10 text-primary'
              : 'border-border text-foreground hover:bg-accent/50'}"
            onclick={() => (servingMode = "serving")}
            title={originalServingLabel}
          >
            {originalServingLabel}
          </button>

          <!-- Custom -->
          <button
            type="button"
            class="rounded-lg border px-3 py-2 text-sm font-medium transition-colors
              {servingMode === 'custom'
              ? 'border-primary bg-primary/10 text-primary'
              : 'border-border text-foreground hover:bg-accent/50'}"
            onclick={() => (servingMode = "custom")}
          >
            Custom
          </button>
        </div>

        <!-- Custom weight input -->
        {#if servingMode === "custom"}
          <div class="flex items-center gap-2">
            <Input
              type="number"
              min="1"
              max="5000"
              bind:value={customWeight}
              class="w-24 text-center"
              id="ingredient-custom-weight"
            />
            <span class="text-sm text-muted-foreground">grams</span>
          </div>
        {/if}
      </div>

      <!-- Live Macro Preview -->
      <div class="rounded-xl border border-border bg-muted/30 p-4 space-y-3">
        <div class="flex items-center justify-between">
          <p class="text-sm font-medium text-muted-foreground">Nutrition for {servingLabel}</p>
          <Badge variant="outline" class="gap-1 text-xs">
            <Scale class="h-3 w-3" />
            {selectedWeight}g
          </Badge>
        </div>

        <!-- Calorie display -->
        <div class="text-center">
          <p class="text-3xl font-bold text-foreground">{scaledMacros.calories}</p>
          <p class="text-xs text-muted-foreground">calories</p>
        </div>

        <!-- Macro breakdown -->
        <div class="grid grid-cols-3 gap-3 text-center">
          <div>
            <p class="text-lg font-semibold text-blue-500">{scaledMacros.protein}g</p>
            <p class="text-[10px] text-muted-foreground uppercase tracking-wider">Protein</p>
          </div>
          <div>
            <p class="text-lg font-semibold text-amber-500">{scaledMacros.carbs}g</p>
            <p class="text-[10px] text-muted-foreground uppercase tracking-wider">Carbs</p>
          </div>
          <div>
            <p class="text-lg font-semibold text-rose-500">{scaledMacros.fat}g</p>
            <p class="text-[10px] text-muted-foreground uppercase tracking-wider">Fat</p>
          </div>
        </div>
      </div>
    </div>
  {/if}
{/snippet}

{#snippet footer()}
  <div class="flex gap-2 w-full">
    <Button variant="outline" class="flex-1" onclick={() => (open = false)}>Cancel</Button>
    <Button class="flex-1" onclick={handleUpdate}>Update Serving</Button>
  </div>
{/snippet}

{#if isMobile.current}
  <Drawer.Root bind:open>
    <Drawer.Content>
      <Drawer.Header>
        <Drawer.Title>Adjust Serving</Drawer.Title>
        <Drawer.Description
          >Change the portion size for {ingredient?.name ?? "this ingredient"}</Drawer.Description
        >
      </Drawer.Header>
      {@render body()}
      <Drawer.Footer>
        {@render footer()}
      </Drawer.Footer>
    </Drawer.Content>
  </Drawer.Root>
{:else}
  <Dialog.Root bind:open>
    <Dialog.Content class="sm:max-w-md">
      <Dialog.Header>
        <Dialog.Title>Adjust Serving</Dialog.Title>
        <Dialog.Description
          >Change the portion size for {ingredient?.name ?? "this ingredient"}</Dialog.Description
        >
      </Dialog.Header>
      {@render body()}
      <Dialog.Footer>
        {@render footer()}
      </Dialog.Footer>
    </Dialog.Content>
  </Dialog.Root>
{/if}
