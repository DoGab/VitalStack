<script lang="ts">
  import { Camera, Upload, Loader2, RotateCcw } from "lucide-svelte";
  import { api } from "$lib/api/client";
  import type { components } from "$lib/api/schema";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Drawer from "$lib/components/ui/drawer";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import ScanResultsDisplay from "./ScanResultsDisplay.svelte";
  import { IsMobile } from "$lib/hooks/is-mobile.svelte.js";

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

  // Camera refs
  let videoElement: HTMLVideoElement | null = $state(null);
  let streamRef: MediaStream | null = $state(null);
  let showCamera = $state(false);

  // File input ref
  let fileInput: HTMLInputElement | null = $state(null);

  // Dynamic title / description
  let title = $derived(scanResult ? "Scan Results" : "Scan Food");
  let subtitle = $derived(
    scanResult
      ? "Here's what we found in your meal"
      : showCamera
        ? "Point your camera at the food"
        : imagePreview
          ? "Review and scan your photo"
          : "Upload a photo of your meal"
  );

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
    error = null;
    showCamera = false;
  }

  async function startCamera() {
    try {
      error = null;
      const stream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: "environment" }
      });
      streamRef = stream;
      showCamera = true;

      // Wait for video element to be available
      await new Promise((resolve) => setTimeout(resolve, 100));
      if (videoElement) {
        videoElement.srcObject = stream;
      }
    } catch (err) {
      error = "Could not access camera. Please check permissions.";
      console.error("Camera error:", err);
    }
  }

  function stopCamera() {
    if (streamRef) {
      streamRef.getTracks().forEach((track) => track.stop());
      streamRef = null;
    }
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
</script>

<!-- Shared modal body -->
{#snippet modalBody()}
  <div class="px-4 pb-4 overflow-y-auto flex-1">
    {#if scanResult}
      <!-- Results View -->
      <ScanResultsDisplay result={scanResult} />
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
    <Button variant="outline" class="w-full" onclick={() => (open = false)}>Close</Button>
    <Button class="w-full" disabled>Add to Log</Button>
  {/if}
{/snippet}

{#if isMobile.current}
  <!-- Mobile: Drawer -->
  <Drawer.Root bind:open>
    <Drawer.Content class="max-h-[90vh]">
      <Drawer.Header>
        <Drawer.Title>{title}</Drawer.Title>
        <Drawer.Description>{subtitle}</Drawer.Description>
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
    <Dialog.Content class="sm:max-w-lg max-h-[90vh] overflow-y-auto">
      <Dialog.Header>
        <Dialog.Title>{title}</Dialog.Title>
        <Dialog.Description>{subtitle}</Dialog.Description>
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
