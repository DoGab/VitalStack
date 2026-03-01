<script lang="ts">
  import TodaysSummary from "$lib/components/dashboard/TodaysSummary.svelte";
  import TodaysMeals from "$lib/components/dashboard/TodaysMeals.svelte";
  import { onMount } from "svelte";
  import { api } from "$lib/api/client";
  import type { components } from "$lib/api/schema";
  import { nutritionState } from "$lib/state/nutrition.svelte";

  type DailyData = components["schemas"]["DailyIntakeOutputBody"];

  let loading = $state(true);
  let error = $state<string | null>(null);
  let dailyData = $state<DailyData | null>(null);

  async function loadData() {
    loading = true;
    error = null;
    try {
      const tzOffset = new Date().getTimezoneOffset();
      const { data, error: apiError } = await api.GET("/api/nutrition/daily", {
        params: { query: { tz_offset: tzOffset } }
      });

      if (apiError) {
        error = "Failed to load dashboard data.";
        console.error("Dashboard apiError:", apiError);
      } else if (data) {
        // Hydrate global state for contextual macro progress bars
        nutritionState.dailyMacros = data.macros;
        dailyData = data;
      }
    } catch (e) {
      error = "Network error loading dashboard.";
      console.error(e);
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    loadData();

    const handleLog = () => loadData();
    window.addEventListener("app:food-logged", handleLog);

    return () => {
      window.removeEventListener("app:food-logged", handleLog);
    };
  });
</script>

<svelte:head>
  <title>VitalStack - Dashboard</title>
  <meta
    name="description"
    content="Track your daily nutrition with AI-powered food scanning. View calories, macros, and meal history at a glance."
  />
</svelte:head>

<div class="container mx-auto px-4 py-6 max-w-5xl space-y-6">
  {#if loading}
    <div class="h-64 flex items-center justify-center">
      <div class="text-muted-foreground animate-pulse">Loading today's stats...</div>
    </div>
  {:else if error}
    <div class="h-64 flex flex-col items-center justify-center gap-4">
      <div class="rounded-lg bg-destructive/10 text-destructive p-4">
        {error}
      </div>
      <button class="text-primary hover:underline" onclick={() => location.reload()}>Retry</button>
    </div>
  {:else if dailyData}
    <TodaysSummary
      caloriesConsumed={dailyData.macros.calories}
      macros={[
        { key: "protein", current: dailyData.macros.protein, goal: 150 },
        { key: "carbs", current: dailyData.macros.carbs, goal: 200 },
        { key: "fat", current: dailyData.macros.fat, goal: 80 }
      ]}
    />
    <TodaysMeals meals={dailyData.meals || []} />
  {/if}
</div>
