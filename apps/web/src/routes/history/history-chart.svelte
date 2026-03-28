<script lang="ts">
  import * as Card from "$lib/components/ui/card/index.js";
  import * as Chart from "$lib/components/ui/chart/index.js";
  import { scaleBand } from "d3-scale";
  import { BarChart, Highlight, type ChartContextValue } from "layerchart";
  import { cubicInOut } from "svelte/easing";

  let {
    title,
    description = "",
    config,
    data,
    series,
    timeRange,
    yDomain,
    yAxisFormatter = (d: number) => {
      if (d >= 1000) return (d / 1000).toFixed(1) + "k";
      return d.toString();
    }
  }: {
    title: string;
    description?: string;
    config: Chart.ChartConfig;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    data: any[];
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    series: any[];
    timeRange: 7 | 14 | 30 | 90;
    yDomain?: [number, number];
    yAxisFormatter?: (d: number) => string;
  } = $props();

  let chartContext = $state<ChartContextValue>();
</script>

<Card.Root class="border-border/50 shadow-sm bg-card">
  <Card.Header>
    <Card.Title>{title}</Card.Title>
    {#if description}
      <Card.Description>{description}</Card.Description>
    {/if}
  </Card.Header>
  <Card.Content class="pt-6">
    {#if data.length > 0}
      <Chart.Container {config} class="h-64 w-full">
        <BarChart
          bind:context={chartContext}
          {data}
          {yDomain}
          xScale={scaleBand().padding(0.25)}
          x="displayDate"
          axis={true}
          rule={false}
          {series}
          seriesLayout="stack"
          props={{
            bars: {
              stroke: "none",
              initialY: chartContext?.height,
              initialHeight: 0,
              motion: {
                y: { type: "tween", duration: 500, easing: cubicInOut },
                height: { type: "tween", duration: 500, easing: cubicInOut }
              }
            },
            highlight: { area: false },
            xAxis: {
              format: (d: string) => {
                if (timeRange === 7) return d;
                if (timeRange === 14 && data.findIndex((item) => item.displayDate === d) % 2 === 0)
                  return d;
                if (timeRange === 30 && data.findIndex((item) => item.displayDate === d) % 4 === 0)
                  return d;
                if (timeRange === 90 && data.findIndex((item) => item.displayDate === d) % 14 === 0)
                  return d;
                return "";
              }
            },
            yAxis: {
              format: yAxisFormatter
            }
          }}
          legend
        >
          {#snippet belowMarks()}
            <Highlight area={{ class: "fill-muted" }} />
          {/snippet}

          {#snippet tooltip()}
            <Chart.Tooltip>
              {#snippet formatter({ name, item })}
                {@const keyMap: Record<string, string> = {
                  protein: "rawProtein",
                  carbs: "rawCarbs",
                  fat: "rawFat",
                  calories: "calories"
                }}
                {@const configKey = name.toLowerCase() as keyof typeof config}
                {@const originalKey = keyMap[configKey]}
                {@const rawVal = (item.payload as Record<string, number>)?.[originalKey] ?? 0}
                {@const indicatorColor = config[configKey]?.color || item.color}
                <div class="flex w-full items-center justify-between gap-8 flex-1">
                  <div class="flex items-center gap-2">
                    <div
                      class="size-2.5 shrink-0 rounded-[2px] border-(--color-border) bg-(--color-bg)"
                      style="--color-bg: {indicatorColor}; --color-border: {indicatorColor};"
                    ></div>
                    <span class="text-muted-foreground">{config[configKey]?.label || name}</span>
                  </div>
                  <span class="text-foreground font-mono font-medium tabular-nums text-right">
                    {rawVal}{configKey === "calories" ? "kcal" : "g"}
                  </span>
                </div>
              {/snippet}
            </Chart.Tooltip>
          {/snippet}
        </BarChart>
      </Chart.Container>
    {/if}
  </Card.Content>
</Card.Root>
