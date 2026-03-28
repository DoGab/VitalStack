<script lang="ts">
  import type { Snippet } from "svelte";
  import { cn } from "$lib/utils.js";

  let {
    title,
    subtitle,
    href,
    actionLabel,
    action,
    class: className
  }: {
    title: string;
    subtitle?: string;
    href?: string;
    actionLabel?: string;
    action?: Snippet;
    class?: string;
  } = $props();
</script>

<div class={cn("flex items-center justify-between pt-4 first:pt-0 pb-2", className)}>
  <div>
    <h2 class="text-lg font-semibold">{title}</h2>
    {#if subtitle}
      <p class="text-sm text-muted-foreground">{subtitle}</p>
    {/if}
  </div>

  {#if action}
    {@render action()}
  {:else if href && actionLabel}
    <a {href} class="text-sm text-primary hover:underline">{actionLabel}</a>
  {/if}
</div>
