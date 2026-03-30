<script lang="ts">
  import Camera from "lucide-svelte/icons/camera";
  import Upload from "lucide-svelte/icons/upload";
  import Loader2 from "lucide-svelte/icons/loader-2";
  import RotateCcw from "lucide-svelte/icons/rotate-ccw";
  import Scale from "lucide-svelte/icons/scale";
  import Check from "lucide-svelte/icons/check";
  import { api } from "$lib/api/client";
  import { useCamera } from "$lib/hooks/use-camera.svelte";
  import type { components } from "$lib/api/schema";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Drawer from "$lib/components/ui/drawer";
  import { Badge } from "$lib/components/ui/badge";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import ScanResultsDisplay from "./ScanResultsDisplay.svelte";
  import { IsMobile } from "$lib/hooks/is-mobile.svelte.js";
  import type { EditableIngredient } from "../dashboard/IngredientEditor.svelte";

  type ScanResult = components["schemas"]["ScanOutputBody"];

  interface Props {
    open: boolean;
    mode: "camera" | "upload";
  }

  let { open = $bindable(), mode = "camera" }: Props = $props();

  const isMobile = new IsMobile();

  // State
  let imagePreview: string | null = $state(null);
  let imageBase64: string | null = $state(null);
  let description = $state("");
  let isScanning = $state(false);
  let scanResult: ScanResult | null = $state(null);
  let error: string | null = $state(null);

  // States for 'Log to Diary' action
  let isLogging = $state(false);
  let hasLogged = $state(false);
  let logError = $state<string | null>(null);

  // Editable Ingredients state
  let editableIngredients = $state<EditableIngredient[]>([]);
  let logDatetime = $state<string>(new Date().toISOString().slice(0, 16)); // "YYYY-MM-DDTHH:mm"

  // Computed Macros
  let computedMacros = $derived.by(() => {
    if (!scanResult) return null;
    if (editableIngredients.length === 0) return scanResult.macros;

    let total = { calories: 0, protein: 0, carbs: 0, fat: 0, fiber: 0 };
    for (const item of editableIngredients) {
      if (item.selected) {
        total.calories += item.macros.calories;
        total.protein += item.macros.protein;
        total.carbs += item.macros.carbs;
        total.fat += item.macros.fat;
      }
    }
    return total;
  });

  // Camera state
  const camera = useCamera();
  let videoElement: HTMLVideoElement | null = $state(null);
  let showCamera = $state(false);

  // File input ref
  let fileInput: HTMLInputElement | null = $state(null);

  // Dynamic title / description
  let title = $derived.by(() => {
    if (scanResult) return scanResult.food_name || "Scan Food";
    return "Scan Food";
  });
  let subtitle = $derived.by(() => {
    if (scanResult) return "Here's what we found in your meal";
    if (showCamera) return "Point your camera at the food";
    if (imagePreview) return "Review and scan your photo";
    return "Upload a photo of your meal";
  });

  // Initialize camera when opening in camera mode
  $effect(() => {
    if (open && mode === "camera" && !imagePreview && !scanResult) {
      startCamera();
    }
  });

  // Cleanup on close
  $effect(() => {
    if (!open) {
      stopCamera();
      resetState();
    }
  });

  function resetState() {
    imagePreview = null;
    imageBase64 = null;
    description = "";
    scanResult = null;
    editableIngredients = [];
    error = null;
    logError = null;
    showCamera = false;
    isLogging = false;
    hasLogged = false;
    logDatetime = new Date().toISOString().slice(0, 16);
  }

  async function startCamera() {
    error = null;
    await camera.start();
    if (camera.error) {
      error = camera.error;
    } else {
      showCamera = true;

      // Wait for video element to be available
      await new Promise((resolve) => setTimeout(resolve, 100));
      if (videoElement && camera.stream) {
        videoElement.srcObject = camera.stream;
      }
    }
  }

  function stopCamera() {
    camera.stop();
    showCamera = false;
  }

  function capturePhoto() {
    if (!videoElement) return;

    const canvas = document.createElement("canvas");
    canvas.width = videoElement.videoWidth;
    canvas.height = videoElement.videoHeight;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    ctx.drawImage(videoElement, 0, 0);

    imagePreview = canvas.toDataURL("image/jpeg", 0.8);
    imageBase64 = imagePreview.split(",")[1];

    stopCamera();
  }

  function handleFileSelect(event: Event) {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      const result = e.target?.result as string;
      imagePreview = result;
      imageBase64 = result.split(",")[1];
    };
    reader.readAsDataURL(file);
  }

  function triggerFileInput() {
    fileInput?.click();
  }

  async function scanFood() {
    if (!imageBase64) return;

    isScanning = true;
    error = null;

    try {
      const { data, error: apiError } = await api.POST("/api/nutrition/scan", {
        body: {
          image_base64: imageBase64,
          description: description || undefined
        }
      });

      if (apiError) {
        error =
          (apiError.errors as Array<{ message: string }> | undefined)?.[0]?.message ||
          apiError.detail ||
          "Failed to scan food. Please try again.";
        return;
      }

      if (data) {
        scanResult = data;
        if (data.ingredients) {
          editableIngredients = data.ingredients.map((ing) => ({
            ...ing,
            selected: true,
            base_quantity: ing.serving_quantity || 1,
            base_macros: { ...ing.macros }
          }));
        } else {
          editableIngredients = [];
        }
      }
    } catch (err) {
      error = "Network error. Please check your connection.";
      console.error("Scan error:", err);
    } finally {
      isScanning = false;
    }
  }

  function retake() {
    imagePreview = null;
    imageBase64 = null;
    scanResult = null;
    error = null;
    if (mode === "camera") {
      startCamera();
    }
  }

  async function handleLogFood() {
    if (!scanResult) return;
    isLogging = true;
    logError = null;

    try {
      const { data, error: logApiError } = await api.POST("/api/nutrition/log", {
        body: {
          food_name: scanResult.food_name,
          confidence: scanResult.confidence,
          macros: computedMacros!,
          ingredients: editableIngredients
            .filter((i) => i.selected)
            .map((i) => ({
              name: i.name,
              serving_quantity: i.serving_quantity,
              serving_size: i.serving_size,
              serving_unit: i.serving_unit,
              macros: i.macros
            }))
          // Will use created_at from user's explicit logDatetime (wait, backend needs to support overriding created_at, but we'll include it tentatively if possible)
          // If the backend doesn't support created_at explicitly yet, it'll at least be prepared.
          // user_id will be handled by auth layer later or hardcoded
        }
      });

      if (logApiError) {
        logError = "Failed to log food. Please try again.";
        console.error("Log error", logApiError);
      } else if (data?.success) {
        hasLogged = true;
        window.dispatchEvent(new CustomEvent("app:food-logged"));
        // Optionally auto-close the modal after a short delay
        setTimeout(() => {
          open = false;
        }, 1500);
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

<!-- Shared modal body -->
{#snippet modalBody()}
  <div class="px-4 pb-4 overflow-y-auto flex-1">
    {#if scanResult && computedMacros}
      <!-- Results View -->
      <ScanResultsDisplay result={scanResult} {computedMacros} bind:editableIngredients />
    {:else if showCamera}
      <!-- Camera View -->
      <div class="space-y-4">
        <div class="relative aspect-[4/3] bg-black rounded-xl overflow-hidden">
          <video
            bind:this={videoElement}
            autoplay
            playsinline
            muted
            class="w-full h-full object-cover"
          ></video>
        </div>
        <Button class="w-full" onclick={capturePhoto}>
          <Camera class="h-5 w-5 mr-2" />
          Take Photo
        </Button>
      </div>
    {:else if imagePreview}
      <!-- Preview View -->
      <div class="space-y-4">
        <div class="relative aspect-[4/3] bg-muted rounded-xl overflow-hidden">
          <img src={imagePreview} alt="Food preview" class="w-full h-full object-cover" />
        </div>

        <!-- Optional Description -->
        <div class="space-y-2">
          <label for="food-description" class="text-sm font-medium text-foreground">
            Description (optional)
          </label>
          <Input
            id="food-description"
            type="text"
            placeholder="e.g., Lunch salad with grilled chicken"
            bind:value={description}
          />
        </div>

        {#if error}
          <div class="rounded-lg bg-destructive/10 text-destructive p-3 text-sm">
            {error}
          </div>
        {/if}

        <div class="flex gap-2">
          <Button variant="outline" class="flex-1" onclick={retake} disabled={isScanning}>
            <RotateCcw class="h-4 w-4 mr-2" />
            Retake
          </Button>
          <Button class="flex-1" onclick={scanFood} disabled={isScanning}>
            {#if isScanning}
              <Loader2 class="h-4 w-4 mr-2 animate-spin" />
              Scanning...
            {:else}
              Scan Food
            {/if}
          </Button>
        </div>
      </div>
    {:else}
      <!-- Upload View -->
      <div class="space-y-4">
        {#if error}
          <div class="rounded-lg bg-destructive/10 text-destructive p-3 text-sm">
            {error}
          </div>
        {/if}

        <input
          bind:this={fileInput}
          type="file"
          accept="image/*"
          class="hidden"
          onchange={handleFileSelect}
        />

        <button
          class="w-full border-2 border-dashed border-border rounded-xl p-8 text-center hover:border-primary/50 transition-colors cursor-pointer"
          onclick={triggerFileInput}
        >
          <Upload class="h-12 w-12 mx-auto text-muted-foreground mb-3" />
          <p class="font-medium text-foreground">Click to upload an image</p>
          <p class="text-sm text-muted-foreground mt-1">JPG, PNG up to 10MB</p>
        </button>

        {#if mode === "upload"}
          <Button variant="ghost" class="w-full" onclick={() => startCamera()}>
            <Camera class="h-5 w-5 mr-2" />
            Use Camera Instead
          </Button>
        {/if}
      </div>
    {/if}
  </div>
{/snippet}

<!-- Shared footer for scan results -->
{#snippet modalFooter()}
  {#if scanResult}
    <div class="flex flex-col w-full gap-2 pt-2">
      {#if logError}
        <p class="text-sm text-destructive font-medium text-center">{logError}</p>
      {/if}
      <div class="flex w-full gap-2">
        <Button variant="outline" class="flex-1 px-2 sm:px-4" onclick={() => (open = false)}
          >Close</Button
        >
        <Button variant="outline" class="flex-1 px-2 sm:px-4" onclick={retake}>Rescan</Button>
        <Button
          class="flex-1 px-2 sm:px-4"
          disabled={isLogging || hasLogged}
          onclick={handleLogFood}
        >
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
    </div>
  {/if}
{/snippet}

{#if isMobile.current}
  <!-- Mobile: Drawer -->
  <Drawer.Root bind:open>
    <Drawer.Content class="max-h-[90vh]">
      <Drawer.Header class={scanResult ? "text-center pb-2" : ""}>
        <Drawer.Title class={scanResult ? "text-2xl font-bold font-heading" : ""}
          >{title}</Drawer.Title
        >
        {#if scanResult}
          <div class="flex flex-wrap items-center justify-center gap-2 mt-2">
            <input
              type="datetime-local"
              bind:value={logDatetime}
              class="h-7 w-[160px] text-xs rounded-full bg-secondary/10 border border-border px-3 text-foreground font-medium outline-none cursor-pointer hover:bg-secondary/20 transition-colors"
              style="appearance: none; -webkit-appearance: none;"
            />
            <Badge
              variant="outline"
              class="gap-1.5 rounded-full border-green-500/20 bg-green-500/10 text-green-500 hover:bg-green-500/20 font-medium whitespace-nowrap"
            >
              <Check class="size-3.5" />
              {Math.round(scanResult.confidence * 100)}% Match
            </Badge>
            <Badge
              variant="outline"
              class="gap-1.5 rounded-full text-muted-foreground font-medium bg-secondary/50 whitespace-nowrap"
            >
              <Scale class="size-3.5" />
              ~{scanResult.serving_size}
            </Badge>
          </div>
        {:else}
          <Drawer.Description>{subtitle}</Drawer.Description>
        {/if}
      </Drawer.Header>
      {@render modalBody()}
      {#if scanResult}
        <Drawer.Footer>
          {@render modalFooter()}
        </Drawer.Footer>
      {/if}
    </Drawer.Content>
  </Drawer.Root>
{:else}
  <!-- Desktop: Dialog -->
  <Dialog.Root bind:open>
    <Dialog.Content class="sm:max-w-[700px] w-full max-h-[90vh] overflow-y-auto">
      <Dialog.Header class={scanResult ? "text-center pb-2" : ""}>
        <Dialog.Title class={scanResult ? "text-2xl font-bold font-heading text-center" : ""}
          >{title}</Dialog.Title
        >
        {#if scanResult}
          <div class="flex flex-wrap items-center justify-center gap-2 mt-2">
            <input
              type="datetime-local"
              bind:value={logDatetime}
              class="h-7 w-[160px] text-xs rounded-full bg-secondary/10 border border-border px-3 text-foreground font-medium outline-none cursor-pointer hover:bg-secondary/20 transition-colors"
              style="appearance: none; -webkit-appearance: none;"
            />
            <Badge
              variant="outline"
              class="gap-1.5 rounded-full border-green-500/20 bg-green-500/10 text-green-500 hover:bg-green-500/20 font-medium whitespace-nowrap"
            >
              <Check class="size-3.5" />
              {Math.round(scanResult.confidence * 100)}% Match
            </Badge>
            <Badge
              variant="outline"
              class="gap-1.5 rounded-full text-muted-foreground font-medium bg-secondary/50 whitespace-nowrap"
            >
              <Scale class="size-3.5" />
              ~{scanResult.serving_size}
            </Badge>
          </div>
        {:else}
          <Dialog.Description>{subtitle}</Dialog.Description>
        {/if}
      </Dialog.Header>
      {@render modalBody()}
      {#if scanResult}
        <Dialog.Footer>
          {@render modalFooter()}
        </Dialog.Footer>
      {/if}
    </Dialog.Content>
  </Dialog.Root>
{/if}
