<script lang="ts">
  import { api } from "$lib/api/client";
  import type { components } from "$lib/api/schema";
  import { Button } from "$lib/components/ui/button/index.js";
  import * as ButtonGroup from "$lib/components/ui/button-group/index.js";
  import * as Chart from "$lib/components/ui/chart/index.js";
  import { NUTRITION_CONFIG } from "$lib/config/nutrition-config";
  import HistoryChart from "./history-chart.svelte";
  import SectionHeader from "$lib/components/ui/section-header.svelte";
  import StatCard from "$lib/components/ui/stat-card.svelte";

  type HistoryData = components["schemas"]["HistoryOutputBody"];

  let loading = $state(true);
  let error = $state<string | null>(null);
  let timeRange = $state<7 | 14 | 30 | 90>(7);
  let historyData = $state<HistoryData | null>(null);

  const caloriesChartConfig = {
    calories: { label: "Calories", color: NUTRITION_CONFIG.calories.barColor }
  } satisfies Chart.ChartConfig;

  const macrosChartConfig = {
    protein: { label: "Protein", color: NUTRITION_CONFIG.protein.barColor },
    carbs: { label: "Carbs", color: NUTRITION_CONFIG.carbs.barColor },
    fat: { label: "Fat", color: NUTRITION_CONFIG.fat.barColor }
  } satisfies Chart.ChartConfig;

  const chartData = $derived(
    historyData?.days?.map((d) => {
      const parts = d.date.split("-");
      const displayDate = `${parts[2]}-${parts[1]}`; // dd-MM
      return {
        date: d.date,
        displayDate,
        protein: d.macros.protein * 4,
        carbs: d.macros.carbs * 4,
        fat: d.macros.fat * 9,
        // Store original values for tooltip display
        rawProtein: d.macros.protein,
        rawCarbs: d.macros.carbs,
        rawFat: d.macros.fat,
        calories: d.macros.calories
      };
    }) || []
  );

  const sharedYDomain = $derived.by(() => {
    if (chartData.length === 0) return undefined;
    const maxVal = Math.max(
      ...chartData.map((d) => Math.max(d.calories, d.protein + d.carbs + d.fat))
    );
    const paddedVal = maxVal * 1.1;
    let maxY = Math.ceil(paddedVal / 50) * 50;
    if (paddedVal > 500) maxY = Math.ceil(paddedVal / 100) * 100;
    if (paddedVal > 1000) maxY = Math.ceil(paddedVal / 500) * 500;
    return [0, maxY] as [number, number];
  });

  const caloriesSeries = [
    {
      key: "calories",
      label: "Calories",
      color: caloriesChartConfig.calories.color,
      props: { rounded: "all" }
    }
  ];

  const macrosSeries = [
    {
      key: "protein",
      label: "Protein",
      color: macrosChartConfig.protein.color,
      props: { rounded: "bottom" }
    },
    {
      key: "carbs",
      label: "Carbs",
      color: macrosChartConfig.carbs.color
    },
    {
      key: "fat",
      label: "Fat",
      color: macrosChartConfig.fat.color
    }
  ];

  async function loadData(days: 7 | 14 | 30 | 90) {
    loading = true;
    error = null;
    try {
      const tzOffset = new Date().getTimezoneOffset();
      const { data, error: apiError } = await api.GET("/api/nutrition/history", {
        params: {
          query: {
            tz_offset: tzOffset,
            days: days
          }
        }
      });

      if (apiError) {
        error = "Failed to load history data.";
        console.error("History apiError:", apiError);
      } else if (data) {
        historyData = data;
      }
    } catch (e) {
      error = "Network error loading history.";
      console.error(e);
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    loadData(timeRange);
  });
</script>

<svelte:head>
  <title>VitalStack - History</title>
</svelte:head>

<div
  class="container mx-auto px-4 py-8 max-w-5xl space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500 pb-24"
>
  <!-- Header & Time Toggle -->
  <div class="space-y-6">
    <div class="flex justify-center">
      <ButtonGroup.Root>
        <Button variant={timeRange === 7 ? "default" : "outline"} onclick={() => (timeRange = 7)}>
          7 Days
        </Button>
        <Button variant={timeRange === 14 ? "default" : "outline"} onclick={() => (timeRange = 14)}>
          14 Days
        </Button>
        <Button variant={timeRange === 30 ? "default" : "outline"} onclick={() => (timeRange = 30)}>
          30 Days
        </Button>
        <Button variant={timeRange === 90 ? "default" : "outline"} onclick={() => (timeRange = 90)}>
          3 Months
        </Button>
      </ButtonGroup.Root>
    </div>
  </div>

  {#if loading && !historyData}
    <div class="h-64 flex flex-col items-center justify-center gap-4">
      <div class="text-muted-foreground animate-pulse">Loading {timeRange}-day history...</div>
    </div>
  {:else if error}
    <div class="h-64 flex flex-col items-center justify-center gap-4">
      <div class="rounded-lg bg-destructive/10 text-destructive px-6 py-4">
        {error}
      </div>
      <Button onclick={() => loadData(timeRange)}>Retry</Button>
    </div>
  {:else if historyData}
    <!-- Averages Section -->
    <div>
      <SectionHeader title="Daily Averages" />

      <div class="flex flex-wrap items-center gap-3">
        <!-- Calories -->
        <StatCard
          value={historyData.averages.calories}
          label="kcal"
          colorValue={NUTRITION_CONFIG.calories.barColor}
        />

        <!-- Protein -->
        <StatCard
          value={historyData.averages.protein}
          valueSuffix="g"
          label="Protein"
          colorValue={NUTRITION_CONFIG.protein.barColor}
        />

        <!-- Carbs -->
        <StatCard
          value={historyData.averages.carbs}
          valueSuffix="g"
          label="Carbs"
          colorValue={NUTRITION_CONFIG.carbs.barColor}
        />

        <!-- Fat -->
        <StatCard
          value={historyData.averages.fat}
          valueSuffix="g"
          label="Fat"
          colorValue={NUTRITION_CONFIG.fat.barColor}
        />
      </div>
    </div>

    <!-- Chart Section -->
    <div>
      <SectionHeader title="Weekly Trends" subtitle="Macro distribution over time" />

      <div class="space-y-4">
        <HistoryChart
          title="Calories Tracked"
          config={caloriesChartConfig}
          data={chartData}
          series={caloriesSeries}
          yDomain={sharedYDomain}
          {timeRange}
        />

        <HistoryChart
          title="Macronutrient Breakdown"
          config={macrosChartConfig}
          data={chartData}
          series={macrosSeries}
          yDomain={sharedYDomain}
          {timeRange}
        />
      </div>
    </div>
  {/if}
</div>
