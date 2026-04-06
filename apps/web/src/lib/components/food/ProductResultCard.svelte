<script lang="ts">
  import type { components } from "$lib/api/schema";
  import { Badge } from "$lib/components/ui/badge";
  import MacroDots from "$lib/components/ui/macro-dots.svelte";
  import ShoppingBasket from "lucide-svelte/icons/shopping-basket";

  type Product = components["schemas"]["ProductOutput"];

  interface Props {
    product: Product;
    onclick?: () => void;
  }

  let { product, onclick }: Props = $props();

  /** NutriScore color mapping */
  const nutriScoreColors: Record<string, string> = {
    a: "bg-green-600 text-white",
    b: "bg-lime-500 text-white",
    c: "bg-yellow-400 text-black",
    d: "bg-orange-500 text-white",
    e: "bg-red-600 text-white"
  };

  /** Datasource display names */
  const sourceLabels: Record<string, string> = {
    off: "Open Food Facts",
    fsvo: "FSVO",
    usda: "USDA"
  };

  let nutriScoreClass = $derived(
    nutriScoreColors[product.nutri_score?.toLowerCase()] || "bg-muted text-muted-foreground"
  );

  let servingHint = $derived.by(() => {
    if (product.serving_size && product.serving_quantity) {
      return `1 serving = ${product.serving_size}`;
    }
    if (product.serving_size) {
      return product.serving_size;
    }
    return null;
  });
</script>

<button
  type="button"
  class="w-full text-left rounded-xl border border-border bg-card p-3 transition-all duration-150
    hover:border-primary/30 hover:bg-accent/50 hover:shadow-sm active:scale-[0.98]
    focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
  {onclick}
>
  <div class="flex gap-3">
    <!-- Product Image -->
    <div
      class="h-16 w-16 shrink-0 rounded-lg bg-muted overflow-hidden flex items-center justify-center"
    >
      {#if product.image_url}
        <img
          src={product.image_url}
          alt={product.name}
          class="h-full w-full object-cover"
          loading="lazy"
        />
      {:else}
        <ShoppingBasket class="h-6 w-6 text-muted-foreground" />
      {/if}
    </div>

    <!-- Product Info -->
    <div class="flex-1 min-w-0">
      <div class="flex items-start justify-between gap-2">
        <div class="min-w-0">
          <p class="font-medium text-sm text-foreground truncate leading-tight">{product.name}</p>
          {#if product.brand}
            <p class="text-xs text-muted-foreground truncate mt-0.5">{product.brand}</p>
          {/if}
        </div>

        <!-- NutriScore Badge -->
        {#if product.nutri_score && product.nutri_score !== "unknown"}
          <span
            class="inline-flex items-center justify-center h-6 w-6 rounded-full text-xs font-bold shrink-0 {nutriScoreClass}"
          >
            {product.nutri_score.toUpperCase()}
          </span>
        {/if}
      </div>

      <!-- Macros Row -->
      <div class="flex items-center justify-between mt-1.5">
        <MacroDots macros={product.macros} decimals={1} />
        <span class="text-[10px] text-muted-foreground/60 shrink-0 ml-2">per 100g</span>
      </div>

      <!-- Bottom row: Source + Serving -->
      <div class="flex items-center gap-1.5 mt-1">
        <Badge
          variant="outline"
          class="h-4 px-1.5 text-[10px] font-normal text-muted-foreground border-border/50"
        >
          {sourceLabels[product.source] || product.source}
        </Badge>
        {#if servingHint}
          <span class="text-[10px] text-muted-foreground/70 truncate">{servingHint}</span>
        {/if}
      </div>
    </div>
  </div>
</button>
