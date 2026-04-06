<script lang="ts">
  import type { components } from "$lib/api/schema";
  import { api } from "$lib/api/client";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Drawer from "$lib/components/ui/drawer";
  import { Input } from "$lib/components/ui/input";
  import { Skeleton } from "$lib/components/ui/skeleton";
  import { IsMobile } from "$lib/hooks/is-mobile.svelte.js";
  import ProductResultCard from "./ProductResultCard.svelte";
  import Search from "lucide-svelte/icons/search";
  import PackageOpen from "lucide-svelte/icons/package-open";
  import { tick } from "svelte";

  type Product = components["schemas"]["ProductOutput"];

  interface Props {
    open: boolean;
    onproductselected?: (product: Product) => void;
  }

  let { open = $bindable(), onproductselected }: Props = $props();

  const isMobile = new IsMobile();

  // Search state
  let query = $state("");
  let isSearching = $state(false);
  let products: Product[] = $state([]);
  let hasSearched = $state(false);
  let searchError = $state<string | null>(null);

  // Input ref for programmatic focus
  let searchInputEl = $state<HTMLInputElement | null>(null);

  // Debounce timer
  let debounceTimer: ReturnType<typeof setTimeout> | null = null;

  // Get user language
  function getUserLang(): string {
    if (typeof navigator !== "undefined") {
      return navigator.language.split("-")[0] || "en";
    }
    return "en";
  }

  // Debounced search
  function handleInput() {
    if (debounceTimer) clearTimeout(debounceTimer);

    if (query.trim().length < 2) {
      products = [];
      hasSearched = false;
      searchError = null;
      return;
    }

    debounceTimer = setTimeout(() => {
      searchProducts();
    }, 300);
  }

  async function searchProducts() {
    const trimmed = query.trim();
    if (trimmed.length < 2) return;

    isSearching = true;
    searchError = null;

    try {
      const { data, error: apiError } = await api.GET("/api/products/search", {
        params: {
          query: {
            query: trimmed,
            limit: 15,
            lang: getUserLang()
          }
        }
      });

      if (apiError) {
        searchError = "Search failed. Please try again.";
        console.error("Search error:", apiError);
        return;
      }

      products = data?.products || [];
      hasSearched = true;
    } catch (e) {
      searchError = "Network error. Please check your connection.";
      console.error("Search error:", e);
    } finally {
      isSearching = false;
    }
  }

  // When true, the next close will NOT reset state (used for back-navigation)
  let skipNextReset = false;

  function selectProduct(product: Product) {
    // Preserve results so the user can come back via the back button
    skipNextReset = true;
    onproductselected?.(product);
  }

  function resetState() {
    query = "";
    products = [];
    hasSearched = false;
    searchError = null;
    isSearching = false;
    if (debounceTimer) clearTimeout(debounceTimer);
  }

  // Auto-focus the search input when the modal opens
  $effect(() => {
    if (open) {
      // Wait for DOM to render, then focus
      tick().then(() => {
        setTimeout(() => {
          searchInputEl?.focus();
        }, 100);
      });
    } else {
      // Only reset state if the close was a dismiss/cancel, not a product selection
      if (skipNextReset) {
        skipNextReset = false;
      } else {
        resetState();
      }
    }
  });
</script>

{#snippet searchBody()}
  <div class="px-4 pb-4 space-y-3">
    <!-- Search Input -->
    <div class="relative">
      <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
      <Input
        type="text"
        placeholder="Search products (e.g. Nutella, Banana...)"
        bind:value={query}
        bind:ref={searchInputEl}
        oninput={handleInput}
        class="pl-10"
        id="product-search-input"
      />
    </div>

    <!-- Results Area -->
    <div class="max-h-[55vh] overflow-y-auto space-y-2 -mx-1 px-1">
      {#if isSearching}
        <!-- Loading Skeletons -->
        {#each [0, 1, 2, 3] as i (i)}
          <div class="flex gap-3 p-3 rounded-xl border border-border bg-card">
            <Skeleton class="h-16 w-16 rounded-lg shrink-0" />
            <div class="flex-1 space-y-2">
              <Skeleton class="h-4 w-3/4" />
              <Skeleton class="h-3 w-1/2" />
              <Skeleton class="h-3 w-full" />
            </div>
          </div>
        {/each}
      {:else if searchError}
        <div class="rounded-lg bg-destructive/10 text-destructive p-3 text-sm">
          {searchError}
        </div>
      {:else if hasSearched && products.length === 0}
        <!-- Empty State -->
        <div class="flex flex-col items-center justify-center py-8 text-center">
          <PackageOpen class="h-12 w-12 text-muted-foreground/40 mb-3" />
          <p class="font-medium text-foreground">No products found</p>
          <p class="text-sm text-muted-foreground mt-1">
            Try different keywords or check the spelling
          </p>
        </div>
      {:else if products.length > 0}
        <!-- Product Results -->
        {#each products as product (product.id)}
          <ProductResultCard {product} onclick={() => selectProduct(product)} />
        {/each}
      {:else}
        <!-- Initial state -->
        <div class="flex flex-col items-center justify-center py-8 text-center">
          <Search class="h-10 w-10 text-muted-foreground/30 mb-3" />
          <p class="text-sm text-muted-foreground">Type to search for products</p>
        </div>
      {/if}
    </div>
  </div>
{/snippet}

{#if isMobile.current}
  <Drawer.Root bind:open>
    <Drawer.Content class="max-h-[90vh]">
      <Drawer.Header>
        <Drawer.Title>Search Products</Drawer.Title>
        <Drawer.Description>Find products by name or brand</Drawer.Description>
      </Drawer.Header>
      {@render searchBody()}
    </Drawer.Content>
  </Drawer.Root>
{:else}
  <Dialog.Root bind:open>
    <Dialog.Content class="sm:max-w-lg max-h-[85vh] overflow-hidden flex flex-col">
      <Dialog.Header>
        <Dialog.Title>Search Products</Dialog.Title>
        <Dialog.Description>Find products by name or brand</Dialog.Description>
      </Dialog.Header>
      {@render searchBody()}
    </Dialog.Content>
  </Dialog.Root>
{/if}
