<script lang="ts">
  import type { components } from "$lib/api/schema";
  import { api } from "$lib/api/client";
  import * as Drawer from "$lib/components/ui/drawer";
  import * as Dialog from "$lib/components/ui/dialog";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Badge } from "$lib/components/ui/badge";
  import { IsMobile } from "$lib/hooks/is-mobile.svelte.js";
  import Loader2 from "lucide-svelte/icons/loader-2";
  import Check from "lucide-svelte/icons/check";
  import Scale from "lucide-svelte/icons/scale";
  import ShoppingBasket from "lucide-svelte/icons/shopping-basket";
  import ArrowLeft from "lucide-svelte/icons/arrow-left";
  import { MODAL_TRANSITION_DELAY_MS } from "$lib/constants";

  type Product = components["schemas"]["ProductOutput"];

  interface Props {
    open: boolean;
    product: Product | null;
    onlogged?: () => void;
    onback?: () => void;
  }

  let { open = $bindable(), product, onlogged, onback }: Props = $props();

  const isMobile = new IsMobile();

  // Serving mode
  type ServingMode = "100g" | "serving" | "custom";
  let servingMode: ServingMode = $state("100g");
  let customWeight = $state(100);

  // Logging state
  let isLogging = $state(false);
  let hasLogged = $state(false);
  let logError = $state<string | null>(null);

  // Reset when product changes or drawer opens
  $effect(() => {
    if (open && product) {
      servingMode = "100g";
      customWeight = 100;
      isLogging = false;
      hasLogged = false;
      logError = null;
    }
  });

  // Computed weight based on serving mode
  let selectedWeight = $derived.by(() => {
    if (!product) return 100;
    if (servingMode === "100g") return 100;
    if (servingMode === "serving") return product.serving_quantity || 100;
    return customWeight;
  });

  // Scaled macros
  let scaledMacros = $derived.by(() => {
    if (!product) return { calories: 0, protein: 0, carbs: 0, fat: 0, fiber: 0 };
    const factor = selectedWeight / 100;
    return {
      calories: Math.round(product.macros.calories * factor),
      protein: +(product.macros.protein * factor).toFixed(1),
      carbs: +(product.macros.carbs * factor).toFixed(1),
      fat: +(product.macros.fat * factor).toFixed(1),
      fiber: +(product.macros.fiber * factor).toFixed(1)
    };
  });

  let servingLabel = $derived.by(() => {
    if (servingMode === "100g") return "100g";
    if (servingMode === "serving" && product?.serving_size) return product.serving_size;
    return `${customWeight}g`;
  });

  async function handleLog() {
    if (!product) return;
    isLogging = true;
    logError = null;

    try {
      const { data, error: logApiError } = await api.POST("/api/nutrition/log", {
        body: {
          food_name: product.name,
          product_id: product.id,
          macros: scaledMacros
        }
      });

      if (logApiError) {
        logError = "Failed to log product. Please try again.";
        console.error("Product log error:", logApiError);
      } else if (data?.success) {
        hasLogged = true;
        window.dispatchEvent(new CustomEvent("app:food-logged"));
        setTimeout(() => {
          open = false;
          onlogged?.();
        }, 1200);
      } else {
        logError = "Unexpected response from server.";
      }
    } catch (e) {
      logError = "A network error occurred.";
      console.error(e);
    } finally {
      isLogging = false;
    }
  }
</script>

{#snippet drawerBody()}
  {#if product}
    <div class="px-4 pb-4 space-y-4">
      <!-- Product Header -->
      <div class="flex items-center gap-3">
        <div
          class="h-14 w-14 shrink-0 rounded-lg bg-muted overflow-hidden flex items-center justify-center"
        >
          {#if product.image_url}
            <img src={product.image_url} alt={product.name} class="h-full w-full object-cover" />
          {:else}
            <ShoppingBasket class="h-5 w-5 text-muted-foreground" />
          {/if}
        </div>
        <div class="min-w-0">
          <p class="font-semibold text-foreground truncate">{product.name}</p>
          {#if product.brand}
            <p class="text-xs text-muted-foreground">{product.brand}</p>
          {/if}
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

          <!-- Product serving (if available) -->
          {#if product.serving_size && product.serving_quantity}
            <button
              type="button"
              class="rounded-lg border px-3 py-2 text-sm font-medium transition-colors truncate
                {servingMode === 'serving'
                ? 'border-primary bg-primary/10 text-primary'
                : 'border-border text-foreground hover:bg-accent/50'}"
              onclick={() => (servingMode = "serving")}
              title={product.serving_size}
            >
              {product.serving_size}
            </button>
          {:else}
            <div></div>
          {/if}

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
              id="custom-weight-input"
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

        <!-- Macro Bars -->
        <div class="text-center">
          <p class="text-3xl font-bold text-foreground">{scaledMacros.calories}</p>
          <p class="text-xs text-muted-foreground">calories</p>
        </div>

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

      <!-- Error -->
      {#if logError}
        <div class="rounded-lg bg-destructive/10 text-destructive p-3 text-sm">
          {logError}
        </div>
      {/if}
    </div>
  {/if}
{/snippet}

{#snippet drawerFooter()}
  <div class="flex gap-2 w-full">
    <Button variant="outline" class="flex-1" onclick={() => (open = false)}>Cancel</Button>
    <Button class="flex-1" disabled={isLogging || hasLogged} onclick={handleLog}>
      {#if isLogging}
        <Loader2 class="mr-2 h-4 w-4 animate-spin" />
        Saving...
      {:else if hasLogged}
        <Check class="mr-2 h-4 w-4" />
        Done
      {:else}
        Add to Log
      {/if}
    </Button>
  </div>
{/snippet}

{#if isMobile.current}
  <Drawer.Root bind:open>
    <Drawer.Content>
      <Drawer.Header>
        <div class="flex items-center gap-2">
          {#if onback}
            <Button
              variant="ghost"
              size="icon"
              class="h-8 w-8 shrink-0 -ml-1"
              onclick={() => {
                open = false;
                setTimeout(() => onback?.(), MODAL_TRANSITION_DELAY_MS);
              }}
            >
              <ArrowLeft class="h-4 w-4" />
            </Button>
          {/if}
          <div>
            <Drawer.Title>Log Product</Drawer.Title>
            <Drawer.Description>Select serving size and add to your diary</Drawer.Description>
          </div>
        </div>
      </Drawer.Header>
      {@render drawerBody()}
      <Drawer.Footer>
        {@render drawerFooter()}
      </Drawer.Footer>
    </Drawer.Content>
  </Drawer.Root>
{:else}
  <Dialog.Root bind:open>
    <Dialog.Content class="sm:max-w-md">
      <Dialog.Header>
        <div class="flex items-center gap-2">
          {#if onback}
            <Button
              variant="ghost"
              size="icon"
              class="h-8 w-8 shrink-0 -ml-1"
              onclick={() => {
                open = false;
                setTimeout(() => onback?.(), MODAL_TRANSITION_DELAY_MS);
              }}
            >
              <ArrowLeft class="h-4 w-4" />
            </Button>
          {/if}
          <div>
            <Dialog.Title>Log Product</Dialog.Title>
            <Dialog.Description>Select serving size and add to your diary</Dialog.Description>
          </div>
        </div>
      </Dialog.Header>
      {@render drawerBody()}
      <Dialog.Footer>
        {@render drawerFooter()}
      </Dialog.Footer>
    </Dialog.Content>
  </Dialog.Root>
{/if}
