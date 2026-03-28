<script lang="ts">
  import * as Card from "$lib/components/ui/card/index.js";
  import { cn } from "$lib/utils.js";

  let {
    value,
    label,
    valueSuffix = "",
    colorValue,
    class: className
  }: {
    value: number | string;
    label: string;
    valueSuffix?: string;
    colorValue?: string;
    class?: string;
  } = $props();
</script>

<Card.Root
  class={cn("flex-1 min-w-[70px] py-1 md:py-4 gap-0 md:gap-4", className)}
  style={colorValue
    ? `background-color: color-mix(in srgb, ${colorValue} 20%, transparent); border-color: color-mix(in srgb, ${colorValue} 50%, transparent);`
    : undefined}
>
  <Card.Content class="px-2 md:px-3 flex flex-col items-center justify-center text-center">
    <p
      class="text-xl font-bold font-mono tracking-tight"
      style={colorValue ? `color: ${colorValue};` : undefined}
    >
      {value}{valueSuffix}
    </p>
    <p class="text-[10px] text-muted-foreground uppercase font-semibold">{label}</p>
  </Card.Content>
</Card.Root>
