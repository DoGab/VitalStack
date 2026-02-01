<script lang="ts">
  import { Flame, Drumstick, Wheat, Droplet, Apple, Scale } from "lucide-svelte";
  import type { components } from "$lib/api/schema";

  type ScanResult = components["schemas"]["ScanOutputBody"];

  let { result, onClose }: { result: ScanResult; onClose: () => void } = $props();

  const macroItems = $derived([
    {
      icon: Flame,
      label: "Calories",
      value: result.macros.calories,
      unit: "kcal",
      color: "text-orange-500"
    },
    {
      icon: Drumstick,
      label: "Protein",
      value: result.macros.protein,
      unit: "g",
      color: "text-red-500"
    },
    { icon: Wheat, label: "Carbs", value: result.macros.carbs, unit: "g", color: "text-amber-500" },
    { icon: Droplet, label: "Fat", value: result.macros.fat, unit: "g", color: "text-blue-500" },
    { icon: Apple, label: "Fiber", value: result.macros.fiber, unit: "g", color: "text-green-500" }
  ]);
</script>

<div class="space-y-4">
  <!-- Food Name & Confidence -->
  <div class="text-center">
    <h3 class="text-xl font-bold text-base-content">{result.food_name}</h3>
    <div class="flex items-center justify-center gap-2 mt-1">
      <div class="badge badge-success badge-sm">
        {Math.round(result.confidence * 100)}% confident
      </div>
    </div>
  </div>

  <!-- Serving Size -->
  <div class="flex items-center justify-center gap-2 text-base-content/70">
    <Scale class="w-4 h-4" />
    <span class="text-sm">{result.serving_size}</span>
  </div>

  <!-- Macro Grid -->
  <div class="grid grid-cols-2 gap-3">
    {#each macroItems as macro, i (macro.label)}
      {@const Icon = macro.icon}
      <div class="stat bg-base-200 rounded-xl p-3 {i === 0 ? 'col-span-2' : ''}">
        <div class="stat-figure {macro.color}">
          <Icon class="w-6 h-6" />
        </div>
        <div class="stat-title text-xs">{macro.label}</div>
        <div class="stat-value text-lg {macro.color}">
          {typeof macro.value === "number" ? macro.value.toFixed(0) : macro.value}
          <span class="text-xs font-normal text-base-content/50">{macro.unit}</span>
        </div>
      </div>
    {/each}
  </div>

  <!-- Actions -->
  <div class="flex gap-2 pt-2">
    <button class="btn btn-outline flex-1" onclick={onClose}> Close </button>
    <button class="btn btn-primary flex-1" disabled> Add to Log </button>
  </div>
</div>
