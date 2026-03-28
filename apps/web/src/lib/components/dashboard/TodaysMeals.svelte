<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import { Badge } from "$lib/components/ui/badge";
  import Flame from "lucide-svelte/icons/flame";
  import Clock from "lucide-svelte/icons/clock";
  import UtensilsCrossed from "lucide-svelte/icons/utensils-crossed";
  import type { components } from "$lib/api/schema";
  import MealDetailsModal from "./MealDetailsModal.svelte";
  import SectionHeader from "$lib/components/ui/section-header.svelte";

  type Meal = components["schemas"]["Meal"];

  interface Props {
    meals?: Meal[];
  }

  let { meals = [] }: Props = $props();

  let hasMeals = $derived(meals.length > 0);

  // Modal State
  let selectedMeal = $state<Meal | null>(null);
  let modalOpen = $state(false);

  function openMealDetails(meal: Meal) {
    selectedMeal = meal;
    modalOpen = true;
  }
</script>

<div>
  <SectionHeader
    title="Today's Meals"
    actionLabel={hasMeals ? "View All" : undefined}
    href={hasMeals ? "/history" : undefined}
  />

  {#if hasMeals}
    <!-- Desktop: horizontal cards -->
    <div class="hidden md:grid md:grid-cols-3 gap-4">
      {#each meals as meal (meal.name)}
        <Card.Root
          class="overflow-hidden hover:shadow-lg transition-all cursor-pointer hover:border-primary/50 active:scale-[0.98]"
          onclick={() => openMealDetails(meal)}
        >
          <!-- Emoji thumbnail area -->
          <div
            class="h-32 bg-muted flex items-center justify-center transition-colors group-hover:bg-primary/5"
          >
            <span class="text-5xl drop-shadow-sm">{meal.emoji}</span>
          </div>
          <Card.Content class="p-3 space-y-2">
            <div class="flex items-center justify-between">
              <h3 class="text-sm font-semibold truncate">{meal.name}</h3>
            </div>
            <div class="flex items-center gap-3 text-xs text-muted-foreground">
              <span class="flex items-center gap-1">
                <Clock class="size-3" />
                {meal.time}
              </span>
              <Badge variant="secondary" class="gap-1 text-[10px] px-1.5 py-0.5">
                <Flame class="size-3" />
                {meal.calories} kcal
              </Badge>
            </div>
            {#if meal.tag}
              <Badge variant="outline" class="text-[10px]">{meal.tag}</Badge>
            {/if}
          </Card.Content>
        </Card.Root>
      {/each}
    </div>

    <!-- Mobile: vertical list -->
    <div class="md:hidden space-y-2">
      {#each meals as meal (meal.name)}
        <Card.Root
          class="cursor-pointer hover:bg-muted/50 active:scale-[0.98] transition-transform"
          onclick={() => openMealDetails(meal)}
        >
          <Card.Content class="flex items-center gap-3 p-3">
            <div class="size-12 rounded-lg bg-muted flex items-center justify-center shrink-0">
              <span class="text-2xl drop-shadow-sm">{meal.emoji}</span>
            </div>
            <div class="flex-1 min-w-0">
              <h3 class="text-sm font-semibold truncate">{meal.name}</h3>
              <div class="flex items-center gap-2 text-xs text-muted-foreground mt-0.5">
                <span class="flex items-center gap-1">
                  <Clock class="size-3" />
                  {meal.time}
                </span>
                <span>·</span>
                <span class="flex items-center gap-1">
                  <Flame class="size-3" />
                  {meal.calories} kcal
                </span>
              </div>
            </div>
          </Card.Content>
        </Card.Root>
      {/each}
    </div>
  {:else}
    <!-- Empty state -->
    <Card.Root>
      <Card.Content class="flex flex-col items-center py-8 text-center">
        <div class="size-12 rounded-full bg-muted flex items-center justify-center mb-3">
          <UtensilsCrossed class="size-6 text-muted-foreground" />
        </div>
        <p class="text-sm text-muted-foreground">
          No meals logged yet today.<br />
          Scan or add your first meal!
        </p>
      </Card.Content>
    </Card.Root>
  {/if}
</div>

<MealDetailsModal bind:open={modalOpen} meal={selectedMeal} />
