<script lang="ts">
  import { NUTRITION_CONFIG } from "$lib/config/nutrition-config";
  import * as Progress from "$lib/components/ui/progress";

  interface MacroData {
    key: string;
    current: number;
    goal: number;
    added?: number;
  }

  interface Props {
    macros: MacroData[];
  }

  let { macros }: Props = $props();

  function formatMacro(value: number): string {
    return value.toFixed(0);
  }
</script>

<div class="space-y-4">
  {#each macros as macro (macro.key)}
    {@const config = NUTRITION_CONFIG[macro.key]}
    {@const Icon = config.icon}
    {@const currentPercent = Math.min((macro.current / macro.goal) * 100, 100)}
    {@const addedAmount = macro.added || 0}
    {@const addedPercent = Math.min((addedAmount / macro.goal) * 100, 100 - currentPercent)}

    <div class="space-y-1.5">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <Icon class="size-4 {config.color}" />
          <span class="text-sm font-medium">{config.label}</span>
        </div>
        <div class="text-sm">
          {#if addedAmount > 0}
            <span class="font-bold text-foreground">+{formatMacro(addedAmount)}{config.unit}</span>
            <span class="text-muted-foreground ml-1">
              ({formatMacro(macro.current + addedAmount)}{config.unit} / {macro.goal}{config.unit})
            </span>
          {:else}
            <span class="text-primary">{formatMacro(macro.current)}{config.unit}</span>
            <span class="text-muted-foreground">/ {macro.goal}{config.unit}</span>
          {/if}
        </div>
      </div>

      {#if addedAmount > 0}
        <!-- Custom Split Progress Bar for incoming items -->
        <div class="bg-muted relative h-2.5 w-full overflow-hidden rounded-full">
          <!-- Current Segment -->
          <div
            class="absolute top-0 bottom-0 left-0 transition-all duration-500 rounded-full"
            style="width: {currentPercent}%; background-color: {config.barColor}"
          ></div>
          <!-- Added Segment (Lighter Transparency) -->
          <div
            class="absolute top-0 bottom-0 transition-all duration-500 opacity-40 block {currentPercent >
            0
              ? 'rounded-r-full'
              : 'rounded-full'}"
            style="left: {currentPercent}%; width: {addedPercent}%; background-color: {config.barColor}"
          ></div>
        </div>
      {:else}
        <!-- Standard Progress Bar -->
        <Progress.Root
          value={currentPercent}
          indicatorColor={config.barColor}
          aria-label={config.label}
        />
      {/if}
    </div>
  {/each}
</div>
