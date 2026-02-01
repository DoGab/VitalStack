<script lang="ts">
  import { Camera, Upload, X, Loader2, RotateCcw } from "lucide-svelte";
  import { api } from "$lib/api/client";
  import type { components } from "$lib/api/schema";
  import ScanResultsDisplay from "./ScanResultsDisplay.svelte";

  type ScanResult = components["schemas"]["ScanOutputBody"];

  let { open = $bindable(), mode = "camera" }: { open: boolean; mode: "camera" | "upload" } =
    $props();

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

  // Initialize based on mode
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

    // Get preview URL and base64
    imagePreview = canvas.toDataURL("image/jpeg", 0.8);
    imageBase64 = imagePreview.split(",")[1]; // Remove data:image/jpeg;base64, prefix

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
        error = apiError.detail || "Failed to scan food. Please try again.";
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

  function handleClose() {
    open = false;
  }
</script>

{#if open}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="fixed inset-0 bg-black/70 z-[80] flex items-center justify-center p-4"
    onclick={handleClose}
  >
    <div
      class="bg-base-100 rounded-2xl shadow-2xl w-full max-w-md max-h-[90vh] overflow-y-auto"
      onclick={(e) => e.stopPropagation()}
    >
      <!-- Header -->
      <div class="flex justify-between items-center p-4 border-b border-base-200">
        <h2 class="text-lg font-bold">
          {#if scanResult}
            Scan Results
          {:else}
            Scan Food
          {/if}
        </h2>
        <button class="btn btn-ghost btn-sm btn-circle" onclick={handleClose}>
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Content -->
      <div class="p-4">
        {#if scanResult}
          <!-- Results View -->
          <ScanResultsDisplay result={scanResult} onClose={handleClose} />
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
            <button class="btn btn-primary w-full" onclick={capturePhoto}>
              <Camera class="w-5 h-5" />
              Take Photo
            </button>
          </div>
        {:else if imagePreview}
          <!-- Preview View -->
          <div class="space-y-4">
            <div class="relative aspect-[4/3] bg-base-200 rounded-xl overflow-hidden">
              <img src={imagePreview} alt="Food preview" class="w-full h-full object-cover" />
            </div>

            <!-- Optional Description -->
            <div class="form-control">
              <label class="label" for="food-description">
                <span class="label-text">Description (optional)</span>
              </label>
              <input
                id="food-description"
                type="text"
                placeholder="e.g., Lunch salad with grilled chicken"
                class="input input-bordered w-full"
                bind:value={description}
              />
            </div>

            {#if error}
              <div class="alert alert-error text-sm">
                <span>{error}</span>
              </div>
            {/if}

            <div class="flex gap-2">
              <button class="btn btn-ghost flex-1" onclick={retake} disabled={isScanning}>
                <RotateCcw class="w-4 h-4" />
                Retake
              </button>
              <button class="btn btn-primary flex-1" onclick={scanFood} disabled={isScanning}>
                {#if isScanning}
                  <Loader2 class="w-4 h-4 animate-spin" />
                  Scanning...
                {:else}
                  Scan Food
                {/if}
              </button>
            </div>
          </div>
        {:else}
          <!-- Upload View (initial state for upload mode) -->
          <div class="space-y-4">
            {#if error}
              <div class="alert alert-error text-sm">
                <span>{error}</span>
              </div>
            {/if}

            <input
              bind:this={fileInput}
              type="file"
              accept="image/*"
              class="hidden"
              onchange={handleFileSelect}
            />

            <div
              class="border-2 border-dashed border-base-300 rounded-xl p-8 text-center hover:border-primary transition-colors cursor-pointer"
              onclick={triggerFileInput}
              role="button"
              tabindex="0"
              onkeydown={(e) => e.key === "Enter" && triggerFileInput()}
            >
              <Upload class="w-12 h-12 mx-auto text-base-content/50 mb-3" />
              <p class="font-medium">Click to upload an image</p>
              <p class="text-sm text-base-content/50 mt-1">JPG, PNG up to 10MB</p>
            </div>

            {#if mode === "upload"}
              <button class="btn btn-ghost w-full" onclick={() => startCamera()}>
                <Camera class="w-5 h-5" />
                Use Camera Instead
              </button>
            {/if}
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}
