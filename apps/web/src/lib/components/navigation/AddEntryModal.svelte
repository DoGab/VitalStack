<script lang="ts">
  import { Camera, Upload, Zap } from "lucide-svelte";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Drawer from "$lib/components/ui/drawer";
  import { Button } from "$lib/components/ui/button";
  import FoodScannerModal from "$lib/components/food/FoodScannerModal.svelte";
  import { IsMobile } from "$lib/hooks/is-mobile.svelte.js";

  interface Props {
    open: boolean;
  }

  let { open = $bindable() }: Props = $props();

  // Detect mobile for Drawer vs Dialog
  const isMobile = new IsMobile();

  // Food scanner state
  let scannerOpen = $state(false);
  let scannerMode: "camera" | "upload" = $state("camera");

  const addOptions = [
    {
      icon: Camera,
      label: "Take Photo",
      description: "Snap a meal photo",
      action: "camera" as const
    },
    {
      icon: Upload,
      label: "Upload Image",
      description: "Choose from gallery",
      action: "upload" as const
    },
    { icon: Zap, label: "Quick Add", description: "Manual macro entry", action: "quick" as const }
  ];

  function handleAddOption(action: "camera" | "upload" | "quick") {
    if (action === "camera") {
      scannerMode = "camera";
      scannerOpen = true;
      open = false;
    } else if (action === "upload") {
      scannerMode = "upload";
      scannerOpen = true;
      open = false;
    } else {
      console.log("Quick add:", action);
      open = false;
    }
  }
</script>

{#snippet optionsList()}
  <div class="flex flex-col gap-2 p-4">
    {#each addOptions as option (option.action)}
      {@const Icon = option.icon}
      <Button
        variant="ghost"
        class="justify-start gap-3 h-16 px-4"
        onclick={() => handleAddOption(option.action)}
      >
        <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center shrink-0">
          <Icon class="h-5 w-5 text-primary" />
        </div>
        <div class="flex flex-col items-start">
          <span class="font-medium">{option.label}</span>
          <span class="text-xs text-muted-foreground">{option.description}</span>
        </div>
      </Button>
    {/each}
  </div>
{/snippet}

<!-- Food Scanner Modal -->
<FoodScannerModal bind:open={scannerOpen} mode={scannerMode} />

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
