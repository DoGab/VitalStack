<script lang="ts">
  import { Camera, Upload, X, Loader2, RotateCcw } from "lucide-svelte";
  import { api } from "$lib/api/client";
  import type { components } from "$lib/api/schema";
  import ScanResultsDisplay from "./ScanResultsDisplay.svelte";

  type ScanResult = components["schemas"]["ScanOutputBody"];

  let { open = $bindable(), mode = "camera" }: { open: boolean; mode: "camera" | "upload" } =
    $props();

  // Dialog element ref
  let dialogElement: HTMLDialogElement | null = $state(null);

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

  // Sync dialog open state with `open` prop
  $effect(() => {
    if (dialogElement) {
      if (open && !dialogElement.open) {
        dialogElement.showModal();
      } else if (!open && dialogElement.open) {
        dialogElement.close();
      }
    }
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

  // Handle dialog close event (ESC key, backdrop click via form)
  function onDialogClose() {
    open = false;
  }
</script>

<dialog
  bind:this={dialogElement}
  class="modal modal-bottom sm:modal-middle"
  onclose={onDialogClose}
>
  <div class="modal-box max-w-md max-h-[90vh] flex flex-col p-0">
    <!-- Header -->
    <div class="flex justify-between items-center p-4 border-b border-base-200 shrink-0">
      <h3 class="text-lg font-bold">
        {#if scanResult}
          Scan Results
        {:else}
          Scan Food
        {/if}
      </h3>
      <form method="dialog">
        <button class="btn btn-ghost btn-sm btn-circle">
          <X class="w-5 h-5" />
        </button>
      </form>
    </div>

    <!-- Content -->
    <div class="p-4 overflow-y-auto flex-1">
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

    <!-- Fixed Footer (only for scan results) -->
    {#if scanResult}
      <div class="p-4 border-t border-base-200 flex gap-2 shrink-0">
        <form method="dialog" class="flex-1">
          <button class="btn btn-outline w-full">Close</button>
        </form>
        <button class="btn btn-primary flex-1" disabled>Add to Log</button>
      </div>
    {/if}
  </div>

  <!-- Backdrop - clicking closes the modal -->
  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
