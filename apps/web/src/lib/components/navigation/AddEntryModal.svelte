<script lang="ts">
  import { Camera, Upload, Zap, X } from "lucide-svelte";
  import FoodScannerModal from "$lib/components/food/FoodScannerModal.svelte";

  let { open = $bindable() } = $props();

  // Food scanner state
  let scannerOpen = $state(false);
  let scannerMode: "camera" | "upload" = $state("camera");

  const addOptions = [
    { icon: Camera, label: "Take Photo", action: "camera" },
    { icon: Upload, label: "Upload Image", action: "upload" },
    { icon: Zap, label: "Quick Add", action: "quick" }
  ];

  function handleAddOption(action: string) {
    if (action === "camera") {
      scannerMode = "camera";
      scannerOpen = true;
      open = false;
    } else if (action === "upload") {
      scannerMode = "upload";
      scannerOpen = true;
      open = false;
    } else {
      console.log("Add action:", action);
      open = false;
    }
  }
</script>

<!-- Food Scanner Modal -->
<FoodScannerModal bind:open={scannerOpen} mode={scannerMode} />

{#if open}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="fixed inset-0 bg-black/50 z-[60]" onclick={() => (open = false)}>
    <div
      class="fixed bottom-24 lg:top-16 lg:bottom-auto lg:right-4 left-1/2 lg:left-auto -translate-x-1/2 lg:translate-x-0 w-64 bg-base-100 rounded-2xl shadow-2xl p-4 z-[70]"
      onclick={(e) => e.stopPropagation()}
    >
      <div class="flex justify-between items-center mb-3">
        <h3 class="font-semibold">Add Entry</h3>
        <button class="btn btn-ghost btn-sm btn-circle" onclick={() => (open = false)}>
          <X class="w-4 h-4" />
        </button>
      </div>
      <div class="flex flex-col gap-2">
        {#each addOptions as option (option.action)}
          {@const Icon = option.icon}
          <button
            class="btn btn-ghost justify-start gap-3 h-14"
            onclick={() => handleAddOption(option.action)}
          >
            <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
              <Icon class="w-5 h-5 text-primary" />
            </div>
            {option.label}
          </button>
        {/each}
      </div>
    </div>
  </div>
{/if}
