<script lang="ts">
  import { Camera, Upload, Zap, X } from "lucide-svelte";
  import FoodScannerModal from "$lib/components/food/FoodScannerModal.svelte";

  let { open = $bindable() } = $props();

  // Dialog element ref
  let dialogElement: HTMLDialogElement | null = $state(null);

  // Food scanner state
  let scannerOpen = $state(false);
  let scannerMode: "camera" | "upload" = $state("camera");

  const addOptions = [
    { icon: Camera, label: "Take Photo", action: "camera" },
    { icon: Upload, label: "Upload Image", action: "upload" },
    { icon: Zap, label: "Quick Add", action: "quick" }
  ];

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

  // Handle dialog close event (ESC key, backdrop click via form)
  function onDialogClose() {
    open = false;
  }
</script>

<!-- Food Scanner Modal -->
<FoodScannerModal bind:open={scannerOpen} mode={scannerMode} />

<dialog
  bind:this={dialogElement}
  class="modal modal-bottom sm:modal-middle"
  onclose={onDialogClose}
>
  <div
    class="modal-box w-full max-w-full sm:max-w-sm p-4 rounded-t-2xl rounded-b-none sm:rounded-2xl"
  >
    <div class="flex justify-between items-center mb-3">
      <h3 class="font-semibold text-lg">Add Entry</h3>
      <form method="dialog">
        <button class="btn btn-ghost btn-sm btn-circle">
          <X class="w-4 h-4" />
        </button>
      </form>
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

  <!-- Backdrop - clicking closes the modal -->
  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
