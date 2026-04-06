<script lang="ts">
  import type { components } from "$lib/api/schema";
  import { api } from "$lib/api/client";
  import { BarqodeStream, type DetectedBarcode } from "barqode";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Drawer from "$lib/components/ui/drawer";
  import { Button } from "$lib/components/ui/button";
  import { IsMobile } from "$lib/hooks/is-mobile.svelte.js";
  import ProductResultCard from "./ProductResultCard.svelte";
  import Loader2 from "lucide-svelte/icons/loader-2";
  import Search from "lucide-svelte/icons/search";
  import ScanBarcode from "lucide-svelte/icons/scan-barcode";
  import TriangleAlert from "lucide-svelte/icons/triangle-alert";
  import RotateCcw from "lucide-svelte/icons/rotate-ccw";

  type Product = components["schemas"]["ProductOutput"];

  interface Props {
    open: boolean;
    onsearchInstead?: () => void;
    onproductselected?: (product: Product) => void;
  }

  let { open = $bindable(), onsearchInstead, onproductselected }: Props = $props();

  const isMobile = new IsMobile();

  // Scanner states
  type ScanState = "scanning" | "loading" | "found" | "not-found" | "error";
  let scanState: ScanState = $state("scanning");
  let cameraLoading = $state(true);
  let scannedBarcode = $state<string | null>(null);
  let foundProduct = $state<Product | null>(null);
  let errorMessage = $state<string | null>(null);

  // Prevent duplicate lookups
  let lastLookedUp = "";

  function getUserLang(): string {
    if (typeof navigator !== "undefined") {
      return navigator.language.split("-")[0] || "en";
    }
    return "en";
  }

  function onCameraOn() {
    cameraLoading = false;
  }

  function onCameraError(error: Error) {
    errorMessage = error.message;
    scanState = "error";
  }

  async function onDetect(detectedCodes: DetectedBarcode[]) {
    const code = detectedCodes[0];
    if (!code?.rawValue) return;

    const barcode = code.rawValue;

    // Prevent re-lookup of same barcode
    if (barcode === lastLookedUp) return;
    lastLookedUp = barcode;

    scannedBarcode = barcode;
    scanState = "loading";

    try {
      const { data, error: apiError } = await api.GET("/api/products/barcode/{barcode}", {
        params: {
          path: { barcode },
          query: { lang: getUserLang() }
        }
      });

      if (apiError) {
        // Check if 404 (not found) vs other error
        if (apiError.status === 404) {
          scanState = "not-found";
        } else {
          errorMessage = "Failed to look up product. Please try again.";
          scanState = "error";
        }
        return;
      }

      if (data) {
        foundProduct = data;
        scanState = "found";
      } else {
        scanState = "not-found";
      }
    } catch (e) {
      errorMessage = "Network error. Please check your connection.";
      scanState = "error";
      console.error("Barcode lookup error:", e);
    }
  }

  function resetScanner() {
    scanState = "scanning";
    cameraLoading = true;
    scannedBarcode = null;
    foundProduct = null;
    errorMessage = null;
    lastLookedUp = "";
  }

  function handleLogProduct() {
    if (foundProduct) {
      onproductselected?.(foundProduct);
    }
  }

  function handleSearchInstead() {
    open = false;
    onsearchInstead?.();
  }

  // Track function — draw highlight rectangle over detected barcode
  function track(detectedCodes: DetectedBarcode[], ctx: CanvasRenderingContext2D) {
    for (const detectedCode of detectedCodes) {
      const [firstPoint, ...otherPoints] = detectedCode.cornerPoints;
      if (!firstPoint) continue;

      ctx.strokeStyle = "hsl(142.1, 76.2%, 36.3%)"; // green-600
      ctx.lineWidth = 3;
      ctx.lineJoin = "round";
      ctx.beginPath();
      ctx.moveTo(firstPoint.x, firstPoint.y);
      for (const { x, y } of otherPoints) {
        ctx.lineTo(x, y);
      }
      ctx.lineTo(firstPoint.x, firstPoint.y);
      ctx.closePath();
      ctx.stroke();
    }
  }

  // Reset on close
  $effect(() => {
    if (!open) {
      // Small delay so the closing animation finishes before resetting
      setTimeout(() => resetScanner(), 300);
    }
  });
</script>

{#snippet scannerBody()}
  <div class="px-4 pb-4 space-y-4">
    {#if scanState === "scanning"}
      <!-- Camera View -->
      <div class="relative aspect-[4/3] rounded-xl overflow-hidden bg-black">
        <BarqodeStream
          formats={["ean_13", "ean_8", "upc_a", "upc_e"]}
          {onCameraOn}
          {onDetect}
          onError={onCameraError}
          {track}
        >
          {#if cameraLoading}
            <div
              class="absolute inset-0 flex flex-col items-center justify-center bg-black/60 text-white gap-2"
            >
              <Loader2 class="h-8 w-8 animate-spin" />
              <p class="text-sm">Starting camera...</p>
            </div>
          {:else}
            <!-- Scan overlay guide -->
            <div class="absolute inset-0 pointer-events-none">
              <div class="absolute inset-0 border-[3px] border-white/20 rounded-xl"></div>
              <!-- Center scanning area indicator -->
              <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-64 h-28">
                <div
                  class="absolute top-0 left-0 w-6 h-6 border-t-2 border-l-2 border-white rounded-tl-lg"
                ></div>
                <div
                  class="absolute top-0 right-0 w-6 h-6 border-t-2 border-r-2 border-white rounded-tr-lg"
                ></div>
                <div
                  class="absolute bottom-0 left-0 w-6 h-6 border-b-2 border-l-2 border-white rounded-bl-lg"
                ></div>
                <div
                  class="absolute bottom-0 right-0 w-6 h-6 border-b-2 border-r-2 border-white rounded-br-lg"
                ></div>
              </div>
              <p class="absolute bottom-4 left-1/2 -translate-x-1/2 text-white/80 text-xs">
                Point at a barcode
              </p>
            </div>
          {/if}
        </BarqodeStream>
      </div>
    {:else if scanState === "loading"}
      <!-- Loading Product -->
      <div class="flex flex-col items-center justify-center py-12 gap-3">
        <Loader2 class="h-10 w-10 text-primary animate-spin" />
        <p class="font-medium text-foreground">Looking up product...</p>
        {#if scannedBarcode}
          <p class="text-xs text-muted-foreground font-mono">{scannedBarcode}</p>
        {/if}
      </div>
    {:else if scanState === "found" && foundProduct}
      <!-- Found Product -->
      <ProductResultCard product={foundProduct} onclick={handleLogProduct} />
      <div class="flex gap-2">
        <Button variant="outline" class="flex-1" onclick={resetScanner}>
          <RotateCcw class="h-4 w-4 mr-2" />
          Scan Another
        </Button>
        <Button class="flex-1" onclick={handleLogProduct}>Add to Log</Button>
      </div>
    {:else if scanState === "not-found"}
      <!-- Not Found -->
      <div class="flex flex-col items-center justify-center py-8 text-center gap-3">
        <div class="h-14 w-14 rounded-full bg-amber-500/10 flex items-center justify-center">
          <ScanBarcode class="h-7 w-7 text-amber-500" />
        </div>
        <div>
          <p class="font-medium text-foreground">Product not found</p>
          <p class="text-sm text-muted-foreground mt-1">This barcode isn't in our database yet</p>
        </div>
        {#if scannedBarcode}
          <p class="text-xs text-muted-foreground font-mono">{scannedBarcode}</p>
        {/if}
        <div class="flex gap-2 w-full mt-2">
          <Button variant="outline" class="flex-1" onclick={resetScanner}>
            <RotateCcw class="h-4 w-4 mr-2" />
            Try Again
          </Button>
          <Button variant="secondary" class="flex-1" onclick={handleSearchInstead}>
            <Search class="h-4 w-4 mr-2" />
            Search Instead
          </Button>
        </div>
      </div>
    {:else if scanState === "error"}
      <!-- Error -->
      <div class="flex flex-col items-center justify-center py-8 text-center gap-3">
        <div class="h-14 w-14 rounded-full bg-destructive/10 flex items-center justify-center">
          <TriangleAlert class="h-7 w-7 text-destructive" />
        </div>
        <div>
          <p class="font-medium text-foreground">Something went wrong</p>
          <p class="text-sm text-muted-foreground mt-1">
            {errorMessage || "Camera access may be blocked"}
          </p>
        </div>
        <Button variant="outline" onclick={resetScanner}>
          <RotateCcw class="h-4 w-4 mr-2" />
          Try Again
        </Button>
      </div>
    {/if}
  </div>
{/snippet}

{#if isMobile.current}
  <Drawer.Root bind:open>
    <Drawer.Content class="max-h-[90vh]">
      <Drawer.Header>
        <Drawer.Title>Scan Barcode</Drawer.Title>
        <Drawer.Description>
          {#if scanState === "scanning"}
            Point your camera at a product barcode
          {:else if scanState === "found"}
            Product found!
          {:else}
            Barcode scanner
          {/if}
        </Drawer.Description>
      </Drawer.Header>
      {@render scannerBody()}
    </Drawer.Content>
  </Drawer.Root>
{:else}
  <Dialog.Root bind:open>
    <Dialog.Content class="sm:max-w-lg">
      <Dialog.Header>
        <Dialog.Title>Scan Barcode</Dialog.Title>
        <Dialog.Description>
          {#if scanState === "scanning"}
            Point your camera at a product barcode
          {:else if scanState === "found"}
            Product found!
          {:else}
            Barcode scanner
          {/if}
        </Dialog.Description>
      </Dialog.Header>
      {@render scannerBody()}
    </Dialog.Content>
  </Dialog.Root>
{/if}
