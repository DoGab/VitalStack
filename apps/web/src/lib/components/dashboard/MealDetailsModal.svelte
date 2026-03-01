<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Drawer from "$lib/components/ui/drawer";
  import { Button } from "$lib/components/ui/button";
  import { Badge } from "$lib/components/ui/badge";
  import Trash2 from "lucide-svelte/icons/trash-2";
  import Edit3 from "lucide-svelte/icons/edit-3";
  import Loader2 from "lucide-svelte/icons/loader-2";
  import Check from "lucide-svelte/icons/check";
  import Scale from "lucide-svelte/icons/scale";
  import ScanResultsDisplay from "$lib/components/food/ScanResultsDisplay.svelte";
  import type { components } from "$lib/api/schema";
  import { IsMobile } from "$lib/hooks/is-mobile.svelte.js";
  import { api } from "$lib/api/client";

  type Meal = components["schemas"]["Meal"];
  // We mock a ScanOutputBody structure that ScanResultsDisplay expects
  type ScanResult = components["schemas"]["ScanOutputBody"];

  interface Props {
    open: boolean;
    meal: Meal | null;
  }

  let { open = $bindable(), meal }: Props = $props();

  const isMobile = new IsMobile();
  let isDeleting = $state(false);
  let deleteError = $state<string | null>(null);

  // Derive a compatible object for ScanResultsDisplay
  let displayResult = $derived.by<ScanResult | null>(() => {
    if (!meal) return null;
    return {
      food_name: meal.name,
      macros: meal.macros,
      confidence: meal.confidence ?? 1.0,
      serving_size: meal.serving_size || "Unknown amount",
      ingredients: meal.ingredients
        ? meal.ingredients.map((ing: components["schemas"]["IngredientBody"]) => ({
            name: ing.name,
            macros: {
              calories: ing.macros.calories,
              protein: ing.macros.protein,
              carbs: ing.macros.carbs,
              fat: ing.macros.fat,
              fiber: ing.macros.fiber
            },
            serving_size: ing.serving_size,
            serving_quantity: ing.serving_quantity,
            serving_unit: ing.serving_unit
          }))
        : []
    };
  });

  async function handleDelete() {
    if (!meal) return;

    isDeleting = true;
    deleteError = null;

    try {
      // NOTE: Our OpenAPI schema generated `/api/nutrition/log/{id}`.
      // The path must exactly match what schema.d.ts exports.
      const { error } = await api.DELETE("/api/nutrition/log/{id}", {
        params: { path: { id: meal.id } }
      });

      if (error) {
        deleteError = "Failed to delete meal.";
        console.error("Delete error:", error);
      } else {
        // Dispatch global event so the Dashboard re-fetches
        window.dispatchEvent(new CustomEvent("app:food-logged"));
        open = false;
      }
    } catch (e) {
      deleteError = "Network error. Please try again.";
      console.error(e);
    } finally {
      isDeleting = false;
    }
  }

  function handleEdit() {
    // Placeholder for future edit feature
    console.log("Edit requested for", meal?.id);
  }
</script>

{#snippet content()}
  <div class="px-4 pb-2 pt-0 space-y-4">
    {#if displayResult}
      <!-- We reuse the beautiful layout from ScanResultsDisplay! -->
      <ScanResultsDisplay mode="details" result={displayResult} />
    {/if}

    {#if deleteError}
      <div
        class="text-sm font-medium text-destructive bg-destructive/10 p-3 rounded-md text-center"
      >
        {deleteError}
      </div>
    {/if}

    <!-- Action Buttons -->
    <div class="flex items-center gap-2 pt-4 border-t border-border mt-auto">
      <Button variant="outline" class="flex-1" onclick={() => (open = false)}>Close</Button>
      <Button variant="outline" class="flex-1 gap-2" onclick={handleEdit}>
        <Edit3 class="size-4" />
        Edit
      </Button>
      <Button
        variant="destructive"
        class="flex-1 gap-2"
        onclick={handleDelete}
        disabled={isDeleting}
      >
        {#if isDeleting}
          <Loader2 class="size-4 animate-spin" />
          Deleting...
        {:else}
          <Trash2 class="size-4" />
          Delete
        {/if}
      </Button>
    </div>
  </div>
{/snippet}

{#if isMobile.current}
  <Drawer.Root bind:open>
    <Drawer.Content>
      <Drawer.Header class="text-center flex flex-col items-center">
        <Drawer.Title class="text-xl text-center w-full">{meal?.name}</Drawer.Title>
        <Drawer.Description class="flex items-center justify-center gap-2 mt-1">
          <span>Logged at {meal?.time}</span>
          {#if meal?.confidence}
            <Badge
              variant="outline"
              class="font-normal border-green-200 bg-green-50 text-green-700"
            >
              <Check class="mr-1 size-3" />
              {Math.round(meal.confidence * 100)}% Match
            </Badge>
          {/if}
          <!-- We fallback to displaying the derived `displayResult.serving_size` for consistency if we have it -->
          {#if displayResult?.serving_size}
            <Badge
              variant="outline"
              class="font-normal text-muted-foreground outline-none border-border bg-muted/50"
            >
              <Scale class="mr-1 size-3" />
              ~{displayResult.serving_size}
            </Badge>
          {/if}
        </Drawer.Description>
      </Drawer.Header>

      <div class="max-h-[70vh] overflow-y-auto">
        {@render content()}
      </div>
    </Drawer.Content>
  </Drawer.Root>
{:else}
  <Dialog.Root bind:open>
    <Dialog.Content class="sm:max-w-md max-h-[85vh] overflow-hidden flex flex-col p-0">
      <Dialog.Header class="p-6 pb-0 shrink-0 text-center flex flex-col items-center">
        <Dialog.Title class="text-xl text-center w-full">{meal?.name}</Dialog.Title>
        <Dialog.Description class="flex items-center justify-center gap-2 mt-1">
          <span>Logged at {meal?.time}</span>
          {#if meal?.confidence}
            <Badge
              variant="outline"
              class="font-normal border-green-200 bg-green-50 text-green-700"
            >
              <Check class="mr-1 size-3" />
              {Math.round(meal.confidence * 100)}% Match
            </Badge>
          {/if}
          {#if displayResult?.serving_size}
            <Badge
              variant="outline"
              class="font-normal text-muted-foreground outline-none border-border bg-muted/50"
            >
              <Scale class="mr-1 size-3" />
              ~{displayResult.serving_size}
            </Badge>
          {/if}
        </Dialog.Description>
      </Dialog.Header>
      <div class="overflow-y-auto">
        {@render content()}
      </div>
    </Dialog.Content>
  </Dialog.Root>
{/if}
