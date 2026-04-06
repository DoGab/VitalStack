<script lang="ts">
  import type { components } from "$lib/api/schema";
  import Camera from "lucide-svelte/icons/camera";
  import Upload from "lucide-svelte/icons/upload";
  import Zap from "lucide-svelte/icons/zap";
  import Search from "lucide-svelte/icons/search";
  import ScanBarcode from "lucide-svelte/icons/scan-barcode";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Drawer from "$lib/components/ui/drawer";
  import { Button } from "$lib/components/ui/button";
  import { Separator } from "$lib/components/ui/separator";
  import FoodScannerModal from "$lib/components/food/FoodScannerModal.svelte";
  import ProductSearchModal from "$lib/components/food/ProductSearchModal.svelte";
  import BarcodeScannerModal from "$lib/components/food/BarcodeScannerModal.svelte";
  import ProductLogDrawer from "$lib/components/food/ProductLogDrawer.svelte";
  import { IsMobile } from "$lib/hooks/is-mobile.svelte.js";
  import { MODAL_TRANSITION_DELAY_MS } from "$lib/constants";

  type Product = components["schemas"]["ProductOutput"];

  interface Props {
    open: boolean;
  }

  let { open = $bindable() }: Props = $props();

  // Detect mobile for Drawer vs Dialog
  const isMobile = new IsMobile();

  // Food scanner state
  let scannerOpen = $state(false);
  let scannerMode: "camera" | "upload" = $state("camera");

  // Product search & barcode scanner state
  let productSearchOpen = $state(false);
  let barcodeScannerOpen = $state(false);

  // Shared product log drawer state (lifted from search/barcode modals)
  let logDrawerOpen = $state(false);
  let selectedProduct = $state<Product | null>(null);
  let productSource: "search" | "barcode" | null = $state(null);

  const aiOptions = [
    {
      icon: Camera,
      label: "Take Photo",
      description: "AI-powered meal analysis",
      action: "camera" as const
    },
    {
      icon: Upload,
      label: "Upload Image",
      description: "Choose from gallery",
      action: "upload" as const
    }
  ];

  const productOptions = [
    {
      icon: Search,
      label: "Search Product",
      description: "Find by name or brand",
      action: "search" as const
    },
    {
      icon: ScanBarcode,
      label: "Scan Barcode",
      description: "Scan a product barcode",
      action: "barcode" as const
    }
  ];

  const otherOptions = [
    { icon: Zap, label: "Quick Add", description: "Manual macro entry", action: "quick" as const }
  ];

  type ActionType = "camera" | "upload" | "search" | "barcode" | "quick";

  function handleAction(action: ActionType) {
    open = false;

    if (action === "camera") {
      scannerMode = "camera";
      scannerOpen = true;
    } else if (action === "upload") {
      scannerMode = "upload";
      scannerOpen = true;
    } else if (action === "search") {
      productSearchOpen = true;
    } else if (action === "barcode") {
      barcodeScannerOpen = true;
    } else {
      console.log("Quick add:", action);
    }
  }

  function handleSearchInstead() {
    barcodeScannerOpen = false;
    productSearchOpen = true;
  }

  /**
   * Called when a product is selected from either the search or barcode modal.
   * Closes the source modal first, then opens the log drawer.
   */
  function handleProductSelected(product: Product) {
    selectedProduct = product;
    // Track which modal the user came from so back button can reopen it
    if (productSearchOpen) productSource = "search";
    else if (barcodeScannerOpen) productSource = "barcode";
    // Close the parent modal first to avoid stacked drawers on mobile
    productSearchOpen = false;
    barcodeScannerOpen = false;
    // Small delay to let the drawer close animation finish before opening the next one
    setTimeout(() => {
      logDrawerOpen = true;
    }, MODAL_TRANSITION_DELAY_MS);
  }

  /**
   * Called after a product is successfully logged.
   * All modals are already closed at this point.
   */
  function handleProductLogged() {
    productSource = null;
  }

  /**
   * Called when the user taps the back arrow on the log drawer.
   * Reopens the source modal (search or barcode) so they can pick a different product.
   */
  function handleLogBack() {
    if (productSource === "search") {
      productSearchOpen = true;
    } else if (productSource === "barcode") {
      barcodeScannerOpen = true;
    }
    productSource = null;
  }
</script>

{#snippet actionButton(option: {
  icon: typeof Camera;
  label: string;
  description: string;
  action: ActionType;
})}
  {@const Icon = option.icon}
  <Button
    variant="ghost"
    class="justify-start gap-3 h-14 px-4"
    onclick={() => handleAction(option.action)}
  >
    <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center shrink-0">
      <Icon class="h-5 w-5 text-primary" />
    </div>
    <div class="flex flex-col items-start">
      <span class="font-medium">{option.label}</span>
      <span class="text-xs text-muted-foreground">{option.description}</span>
    </div>
  </Button>
{/snippet}

{#snippet optionsList()}
  <div class="flex flex-col p-4 gap-0.5">
    <!-- AI Scan Section -->
    <p class="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground px-4 mb-1">
      AI Scan
    </p>
    {#each aiOptions as option (option.action)}
      {@render actionButton(option)}
    {/each}

    <Separator class="my-2" />

    <!-- Product Database Section -->
    <p class="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground px-4 mb-1">
      Product Database
    </p>
    {#each productOptions as option (option.action)}
      {@render actionButton(option)}
    {/each}

    <Separator class="my-2" />

    <!-- Other Section -->
    {#each otherOptions as option (option.action)}
      {@render actionButton(option)}
    {/each}
  </div>
{/snippet}

<!-- Sub-Modals -->
<FoodScannerModal bind:open={scannerOpen} mode={scannerMode} />
<ProductSearchModal bind:open={productSearchOpen} onproductselected={handleProductSelected} />
<BarcodeScannerModal
  bind:open={barcodeScannerOpen}
  onsearchInstead={handleSearchInstead}
  onproductselected={handleProductSelected}
/>

<!-- Shared Product Log Drawer (lifted here to avoid stacked drawers on mobile) -->
<ProductLogDrawer
  bind:open={logDrawerOpen}
  product={selectedProduct}
  onlogged={handleProductLogged}
  onback={handleLogBack}
/>

{#if isMobile.current}
  <!-- Mobile: Drawer (slides up from bottom) -->
  <Drawer.Root bind:open>
    <Drawer.Content>
      <Drawer.Header>
        <Drawer.Title>Add Entry</Drawer.Title>
        <Drawer.Description>Choose how to log your meal</Drawer.Description>
      </Drawer.Header>
      {@render optionsList()}
      <Drawer.Footer>
        <Drawer.Close>
          <Button variant="outline" class="w-full">Cancel</Button>
        </Drawer.Close>
      </Drawer.Footer>
    </Drawer.Content>
  </Drawer.Root>
{:else}
  <!-- Desktop: Centered Dialog -->
  <Dialog.Root bind:open>
    <Dialog.Content class="sm:max-w-md">
      <Dialog.Header>
        <Dialog.Title>Add Entry</Dialog.Title>
        <Dialog.Description>Choose how to log your meal</Dialog.Description>
      </Dialog.Header>
      {@render optionsList()}
      <Dialog.Footer>
        <Dialog.Close>
          <Button variant="outline" class="w-full">Cancel</Button>
        </Dialog.Close>
      </Dialog.Footer>
    </Dialog.Content>
  </Dialog.Root>
{/if}
