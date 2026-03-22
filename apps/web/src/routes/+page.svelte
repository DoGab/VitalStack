<script lang="ts">
  import TodaysSummary from "$lib/components/dashboard/TodaysSummary.svelte";
  import TodaysMeals from "$lib/components/dashboard/TodaysMeals.svelte";
  import { onMount } from "svelte";
  import { api } from "$lib/api/client";
  import type { components } from "$lib/api/schema";
  import { nutritionState } from "$lib/state/nutrition.svelte";
  import { Button } from "$lib/components/ui/button/index.js";
  import { SvelteDate } from "svelte/reactivity";
  import ChevronLeft from "lucide-svelte/icons/chevron-left";
  import ChevronRight from "lucide-svelte/icons/chevron-right";

  type DailyData = components["schemas"]["DailyIntakeOutputBody"];

  let loading = $state(true);
  let error = $state<string | null>(null);
  let dailyData = $state<DailyData | null>(null);

  // Start with today's date stored as a timestamp to avoid Svelte Date reactivity warnings
  let currentDateTs = $state(Date.now());

  // Derived formatted date (e.g. "Today", "Yesterday", "Oct 24, 2024")
  let formattedDate = $derived.by(() => {
    const d = new SvelteDate(currentDateTs);
    const today = new SvelteDate();
    const isToday = d.toDateString() === today.toDateString();

    // Check if yesterday
    const yesterday = new SvelteDate();
    yesterday.setDate(today.getDate() - 1);
    const isYesterday = d.toDateString() === yesterday.toDateString();

    if (isToday) return "Today";
    if (isYesterday) return "Yesterday";

    return new Intl.DateTimeFormat("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric"
    }).format(d);
  });

  let isToday = $derived(
    new SvelteDate(currentDateTs).toDateString() === new SvelteDate().toDateString()
  );

  $effect(() => {
    // Reload when date changes
    loadData(currentDateTs);
  });

  async function loadData(dateTs: number) {
    loading = true;
    error = null;
    try {
      const date = new SvelteDate(dateTs);
      const tzOffset = new SvelteDate().getTimezoneOffset();

      // Format as YYYY-MM-DD local time
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, "0");
      const day = String(date.getDate()).padStart(2, "0");
      const dateStr = `${year}-${month}-${day}`;

      const { data, error: apiError } = await api.GET("/api/nutrition/daily", {
        params: {
          query: {
            tz_offset: tzOffset,
            date: dateStr
          }
        }
      });

      if (apiError) {
        error = "Failed to load dashboard data.";
        console.error("Dashboard apiError:", apiError);
      } else if (data) {
        // Hydrate global state ONLY if the date is today (for the contextual bars to be relevant)
        if (
          dateStr ===
          `${new SvelteDate().getFullYear()}-${String(new SvelteDate().getMonth() + 1).padStart(2, "0")}-${String(new SvelteDate().getDate()).padStart(2, "0")}`
        ) {
          nutritionState.dailyMacros = data.macros;
        }
        dailyData = data;
      }
    } catch (e) {
      error = "Network error loading dashboard.";
      console.error(e);
    } finally {
      loading = false;
    }
  }

  function previousDay() {
    const newDate = new SvelteDate(currentDateTs);
    newDate.setDate(newDate.getDate() - 1);
    currentDateTs = newDate.getTime();
  }

  function nextDay() {
    if (isToday) return;
    const newDate = new SvelteDate(currentDateTs);
    newDate.setDate(newDate.getDate() + 1);
    currentDateTs = newDate.getTime();
  }

  onMount(() => {
    const handleLog = () => loadData(currentDateTs);
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
  <!-- Date Navigation -->
  <div class="flex items-center justify-between mb-2">
    <Button variant="ghost" size="icon" onclick={previousDay} aria-label="Previous day">
      <ChevronLeft class="h-5 w-5" />
    </Button>
    <div class="flex flex-col items-center">
      <h2 class="text-xl font-semibold font-playfair">{formattedDate}</h2>
      {#if !isToday}
        <p class="text-xs text-muted-foreground mt-0.5">Historical log</p>
      {/if}
    </div>
    <Button
      variant="ghost"
      size="icon"
      onclick={nextDay}
      disabled={isToday}
      aria-label="Next day"
      class={isToday ? "opacity-30 pointer-events-none" : ""}
    >
      <ChevronRight class="h-5 w-5" />
    </Button>
  </div>

  {#if loading && !dailyData}
    <div class="h-64 flex flex-col items-center justify-center gap-4">
      <div class="text-muted-foreground animate-pulse">Loading {formattedDate}'s stats...</div>
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
