<script lang="ts">
  import { Checkbox } from "$lib/components/ui/checkbox";
  import MacroDots from "$lib/components/ui/macro-dots.svelte";
  import IngredientServingDrawer from "$lib/components/food/IngredientServingDrawer.svelte";
  import ChevronRight from "lucide-svelte/icons/chevron-right";
  import type { components } from "$lib/api/schema";

  export type EditableIngredient = components["schemas"]["IngredientBody"] & {
    selected: boolean;
    base_quantity: number;
    base_macros: {
      calories: number;
      protein: number;
      carbs: number;
      fat: number;
      fiber: number;
    };
  };

  interface Props {
    ingredients: EditableIngredient[];
    readonly?: boolean;
  }

  let { ingredients = $bindable(), readonly = false }: Props = $props();

  // Serving drawer state
  let servingDrawerOpen = $state(false);
  let selectedIngredientIndex = $state<number | null>(null);
  let selectedIngredient = $derived(
    selectedIngredientIndex !== null ? ingredients[selectedIngredientIndex] : null
  );

  function fmtWeight(item: EditableIngredient): string {
    const qty = item.serving_quantity ?? 1;
    const size = item.serving_size ?? "";
    const unit = item.serving_unit ?? "";
    if (size && unit) return `${qty} × ${size}${unit}`;
    if (unit) return `${qty}${unit}`;
    return `${qty} serving`;
  }

  function openServingDrawer(index: number) {
    selectedIngredientIndex = index;
    servingDrawerOpen = true;
  }

  function handleServingUpdate(updated: EditableIngredient) {
    if (selectedIngredientIndex === null) return;
    // Update the ingredient in-place so the parent's bindable array reflects changes
    ingredients[selectedIngredientIndex] = {
      ...ingredients[selectedIngredientIndex],
      macros: updated.macros,
      serving_quantity: updated.serving_quantity
    };
  }
</script>

<div class="flex flex-col divide-y divide-border/40">
  {#each ingredients as item, i (item.name)}
    <div
      class="flex flex-row items-center gap-3 py-2.5 transition-opacity {item.selected
        ? ''
        : 'opacity-40'}"
    >
      <!-- Checkbox -->
      {#if !readonly}
        <Checkbox
          bind:checked={item.selected}
          class="size-4 shrink-0"
          id={`checkbox-${item.name}`}
        />
      {/if}

      <!-- Tappable ingredient row — opens serving drawer in edit mode -->
      {#if readonly}
        <div class="flex-1 min-w-0 flex flex-col justify-center">
          <span class="text-sm font-semibold text-foreground truncate leading-tight">
            {item.name}
          </span>
          <MacroDots macros={item.macros} class="mt-1" />
        </div>
        <span class="text-xs text-muted-foreground shrink-0 font-mono">
          {fmtWeight(item)}
        </span>
      {:else}
        <button
          type="button"
          class="flex-1 min-w-0 flex items-center gap-2 text-left rounded-lg -mx-1 px-1 py-0.5
            hover:bg-accent/50 active:scale-[0.98] transition-all"
          onclick={() => openServingDrawer(i)}
        >
          <div class="flex-1 min-w-0 flex flex-col justify-center">
            <span class="text-sm font-semibold text-foreground truncate leading-tight">
              {item.name}
            </span>
            <MacroDots macros={item.macros} class="mt-1" />
          </div>
          <span class="text-xs text-muted-foreground shrink-0 font-mono">
            {fmtWeight(item)}
          </span>
          <ChevronRight class="size-3.5 text-muted-foreground/50 shrink-0" />
        </button>
      {/if}
    </div>
  {/each}
</div>

<!-- Serving Size Drawer -->
<IngredientServingDrawer
  bind:open={servingDrawerOpen}
  ingredient={selectedIngredient}
  onupdate={handleServingUpdate}
/>
