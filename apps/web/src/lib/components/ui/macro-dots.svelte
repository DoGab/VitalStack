<script lang="ts">
  /**
   * MacroDots — Compact colored-dot macro display.
   *
   * Renders a row of colored dots with values:
   *   ● 150 Cal  ● 5P  ● 3F  ● 27C
   *
   * Used in ProductResultCard and IngredientEditor for a unified look.
   */

  interface MacroValues {
    calories: number;
    protein: number;
    carbs: number;
    fat: number;
  }

  interface Props {
    macros: MacroValues;
    /** Whether to show the calories dot. Default: true */
    showCalories?: boolean;
    /** Decimal places for macro values. Default: 0 */
    decimals?: number;
    /** Additional CSS classes for the container */
    class?: string;
  }

  let { macros, showCalories = true, decimals = 0, class: className = "" }: Props = $props();

  function fmt(value: number): string {
    return value.toFixed(decimals);
  }

  /** Dot color configuration */
  const dots = [
    { key: "calories" as const, color: "bg-[#D65A31]", suffix: " Cal", caloriesOnly: true },
    { key: "protein" as const, color: "bg-red-500", suffix: "P", caloriesOnly: false },
    { key: "fat" as const, color: "bg-blue-500", suffix: "F", caloriesOnly: false },
    { key: "carbs" as const, color: "bg-amber-500", suffix: "C", caloriesOnly: false }
  ] as const;
</script>

<div class="flex items-center gap-3 flex-wrap {className}">
  {#each dots as dot (dot.key)}
    {#if !dot.caloriesOnly || showCalories}
      <span class="flex items-center gap-1 text-xs text-muted-foreground">
        <span class="size-2 rounded-full {dot.color} shrink-0"></span>
        {fmt(macros[dot.key])}{dot.suffix}
      </span>
    {/if}
  {/each}
</div>
