<script lang="ts">
  import type { ComponentType } from "svelte";
  import ChevronDown from "lucide-svelte/icons/chevron-down";
  import Flame from "lucide-svelte/icons/flame";
  import Plus from "lucide-svelte/icons/plus";
  import type { components } from "$lib/api/schema";
  import { NUTRITION_CONFIG, getMacroDisplayOrder } from "$lib/config/nutrition-config";
  import * as Card from "$lib/components/ui/card";
  import * as Accordion from "$lib/components/ui/accordion";
  import CircularProgress from "$lib/components/ui/circular-progress.svelte";
  import MacroBars from "$lib/components/ui/macro-bars.svelte";
  import { nutritionState } from "$lib/state/nutrition.svelte";
  import IngredientEditor, { type EditableIngredient } from "../dashboard/IngredientEditor.svelte";

  type ScanResult = components["schemas"]["ScanOutputBody"];

  interface Props {
    result: ScanResult;
    mode?: "preview" | "details"; // Determines how progress bars render against daily totals
    computedMacros?: {
      calories: number;
      protein: number;
      carbs: number;
      fat: number;
      fiber: number;
    };
    editableIngredients?: EditableIngredient[];
    readonly?: boolean;
  }

  let {
    result,
    mode = "preview",
    computedMacros,
    editableIngredients = $bindable([]),
    readonly = false
  }: Props = $props();

  // Reference daily goals for progress bars
  const NUTRITION_GOALS: Record<string, number> = {
    calories: 2200,
    protein: 150,
    carbs: 200,
    fat: 80
  };

  // For previews (scanning new food), we show daily macros as solid context and the scan as the transparent preview.
  // For details (viewing already logged food), we show ONLY the meal's macros as solid and 0 added.
  const activeMacros = $derived(computedMacros || result.macros);
  const currentIntake = $derived(mode === "preview" ? nutritionState.safeMacros : activeMacros);

  const addedMacros = $derived(
    mode === "preview" ? activeMacros : { calories: 0, protein: 0, carbs: 0, fat: 0, fiber: 0 }
  );

  // Derive totals for labels
  const totalCalories = $derived(currentIntake.calories + addedMacros.calories);

  // Derive percentages specifically for CircularProgress
  const currentCalPercent = $derived(
    Math.min((currentIntake.calories / NUTRITION_GOALS.calories) * 100, 100)
  );

  const addedCalPercent = $derived(
    Math.min((addedMacros.calories / NUTRITION_GOALS.calories) * 100, 100 - currentCalPercent)
  );

  const macroItems = $derived.by(() => {
    return getMacroDisplayOrder()
      .filter((k) => k !== "fiber" && k !== "calories")
      .map((key) => ({
        ...NUTRITION_CONFIG[key],
        current: currentIntake[key as keyof typeof currentIntake],
        added: addedMacros[key as keyof typeof addedMacros],
        goal: NUTRITION_GOALS[key],
        key
      })) as Array<{
      key: string;
      label: string;
      icon: ComponentType;
      color: string;
      barColor: string;
      unit: string;
      current: number;
      added: number;
      goal: number;
    }>;
  });

  function formatMacro(value: number): string {
    return value.toFixed(0);
  }
</script>

<div class="space-y-3">
  <!-- Meal Macros Breakdown (TodaysSummary Replica Layout) -->
  <Card.Root>
    <Card.Header class="flex flex-row items-center justify-between pb-4">
      <div class="space-y-1">
        <Card.Title class="text-base font-semibold text-foreground"
          >Nutrient Contribution</Card.Title
        >
        <Card.Description>
          {#if mode === "preview"}
            <span class="font-bold text-foreground">+{formatMacro(activeMacros.calories)}kcal</span>
            ({totalCalories.toLocaleString()} / {NUTRITION_GOALS.calories} kcal)
          {:else}
            <span class="font-bold text-foreground">{formatMacro(activeMacros.calories)}kcal</span>
            <span class="text-muted-foreground font-normal">contribution</span>
          {/if}
        </Card.Description>
      </div>

      <!-- Smaller Calorie ring (Mobile) -->
      <CircularProgress
        class="w-14 h-14 md:hidden block text-primary shrink-0"
        size={64}
        radius={28}
        strokeWidth={5}
        percent={currentCalPercent}
        addedPercent={addedCalPercent}
      >
        <Flame class="size-4" />
      </CircularProgress>
    </Card.Header>
    <Card.Content>
      <!-- Mobile Layout -->
      <div class="md:hidden">
        <MacroBars macros={macroItems} />
      </div>

      <!-- Desktop Layout -->
      <div class="hidden md:flex flex-row items-center gap-8">
        <!-- Large Calorie ring on the left -->
        <div class="flex flex-col items-center shrink-0">
          <CircularProgress
            class="w-32 h-32"
            size={120}
            radius={54}
            strokeWidth={8}
            percent={currentCalPercent}
            addedPercent={addedCalPercent}
          >
            <Flame class="size-6 text-primary mb-1" />
            <span class="text-2xl font-bold text-foreground" style="font-family: var(--font-mono);">
              {#if mode === "preview"}+{/if}{formatMacro(activeMacros.calories)}
            </span>
          </CircularProgress>
        </div>

        <!-- Macro progress bars on the right -->
        <div class="flex-1 w-full">
          <MacroBars macros={macroItems} />
        </div>
      </div>
    </Card.Content>
  </Card.Root>

  <!-- Collapsible Ingredient Breakdown using shadcn Accordion -->
  {#if editableIngredients && editableIngredients.length > 0}
    <Accordion.Root type="single" class="w-full">
      <Accordion.Item value="ingredients" class="border-b border-border/40">
        <Accordion.Trigger
          class="py-2 hover:no-underline w-full [&_[data-slot=accordion-trigger-icon]]:hidden px-1 group"
        >
          <div class="flex items-center justify-between w-full">
            <div class="flex items-center gap-2">
              <h3 class="text-base font-semibold text-foreground">Ingredients</h3>
              <ChevronDown
                class="size-4 shrink-0 transition-transform duration-200 group-aria-expanded/accordion-trigger:rotate-180 text-muted-foreground"
              />
            </div>
            {#if !readonly}
              <button
                class="text-[11px] tracking-wider font-semibold text-muted-foreground flex items-center gap-1 hover:text-foreground transition-colors z-10 p-1"
                onclick={(e) => {
                  e.stopPropagation();
                  console.log("Add ingredient clicked");
                }}
                onkeydown={(e) => {
                  if (e.key === "Enter" || e.key === " ") e.stopPropagation();
                }}
              >
                <Plus class="size-3" /> ADD INGREDIENT
              </button>
            {/if}
          </div>
        </Accordion.Trigger>
        <Accordion.Content>
          <IngredientEditor bind:ingredients={editableIngredients} {readonly} />
        </Accordion.Content>
      </Accordion.Item>
    </Accordion.Root>
  {/if}
</div>
