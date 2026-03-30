<script lang="ts">
  import { Checkbox } from "$lib/components/ui/checkbox";
  import { Input } from "$lib/components/ui/input";
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

  function fmt(value: number, decimals = 0): string {
    return value.toFixed(decimals);
  }

  function fmtWeight(item: EditableIngredient): string {
    const qty = item.serving_quantity ?? 1;
    const size = item.serving_size ?? "";
    const unit = item.serving_unit ?? "";
    if (size && unit) return `${qty} × ${size}${unit}`;
    if (unit) return `${qty}${unit}`;
    return `${qty} serving`;
  }

  function handleQuantityChange(item: EditableIngredient, newQty: number) {
    if (newQty < 0) return;
    const ratio = newQty / (item.base_quantity || 1);
    item.macros = {
      calories: item.base_macros.calories * ratio,
      protein: item.base_macros.protein * ratio,
      carbs: item.base_macros.carbs * ratio,
      fat: item.base_macros.fat * ratio,
      fiber: item.base_macros.fiber * ratio
    };
    item.serving_quantity = newQty;
  }
</script>

<div class="flex flex-col divide-y divide-border/40">
  {#each ingredients as item (item.name)}
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

      <!-- Name + macros column -->
      <div class="flex-1 min-w-0 flex flex-col justify-center">
        <span class="text-sm font-semibold text-foreground truncate leading-tight">
          {item.name}
        </span>

        <!-- Macro dots row — always visible -->
        <div class="flex items-center gap-3 mt-1 flex-wrap">
          <span class="flex items-center gap-1 text-xs text-muted-foreground">
            <span class="size-2 rounded-full bg-[#D65A31] shrink-0"></span>
            {fmt(item.macros.calories)} Cal
          </span>
          <span class="flex items-center gap-1 text-xs text-muted-foreground">
            <span class="size-2 rounded-full bg-red-500 shrink-0"></span>
            {fmt(item.macros.protein)}P
          </span>
          <span class="flex items-center gap-1 text-xs text-muted-foreground">
            <span class="size-2 rounded-full bg-blue-500 shrink-0"></span>
            {fmt(item.macros.fat)}F
          </span>
          <span class="flex items-center gap-1 text-xs text-muted-foreground">
            <span class="size-2 rounded-full bg-amber-500 shrink-0"></span>
            {fmt(item.macros.carbs)}C
          </span>
        </div>
      </div>

      <!-- Weight display / editable qty -->
      {#if readonly}
        <span class="text-xs text-muted-foreground shrink-0 font-mono">
          {fmtWeight(item)}
        </span>
      {:else}
        <div class="flex items-center gap-1 shrink-0 font-mono text-xs">
          <Input
            type="number"
            step="0.1"
            min="0"
            class="w-14 h-7 text-center text-xs p-1 font-mono border-outline_variant/30 bg-surface_container_low"
            value={item.serving_quantity ?? 1}
            oninput={(e) => {
              const val = parseFloat(e.currentTarget.value);
              if (!isNaN(val)) handleQuantityChange(item, val);
            }}
          />
          <span class="text-muted-foreground"
            >× {item.serving_size ?? ""}{item.serving_unit ?? "srv"}</span
          >
        </div>
      {/if}
    </div>
  {/each}
</div>
