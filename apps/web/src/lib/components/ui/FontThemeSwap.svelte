<script lang="ts">
  import { Type, Check } from "lucide-svelte";
  import { Button } from "$lib/components/ui/button";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu";

  const fontThemes = [
    { id: "classic", label: "Classic", description: "Playfair Display" },
    { id: "modern", label: "Modern", description: "Outfit" }
  ] as const;

  type FontTheme = (typeof fontThemes)[number]["id"];

  let currentFont = $state<FontTheme>("classic");

  function setFontTheme(theme: FontTheme) {
    currentFont = theme;
    if (typeof document !== "undefined") {
      document.documentElement.setAttribute("data-font-theme", theme);
      localStorage.setItem("font-theme", theme);
    }
  }

  // Load saved font theme on mount
  $effect(() => {
    if (typeof window !== "undefined") {
      const saved = localStorage.getItem("font-theme") as FontTheme | null;
      if (saved && fontThemes.some((t) => t.id === saved)) {
        currentFont = saved;
        document.documentElement.setAttribute("data-font-theme", saved);
      }
    }
  });
</script>

<DropdownMenu.Root>
  <DropdownMenu.Trigger>
    {#snippet child({ props })}
      <Button {...props} variant="ghost" size="icon" class="rounded-full" aria-label="Font Theme">
        <Type class="h-4 w-4" />
      </Button>
    {/snippet}
  </DropdownMenu.Trigger>
  <DropdownMenu.Content align="end" class="w-48">
    <DropdownMenu.Label>Font Theme</DropdownMenu.Label>
    <DropdownMenu.Separator />
    {#each fontThemes as theme (theme.id)}
      <DropdownMenu.Item onclick={() => setFontTheme(theme.id)}>
        <span class="flex items-center justify-between w-full">
          <span>
            <span class="font-medium">{theme.label}</span>
            <span class="text-xs text-muted-foreground ml-2">{theme.description}</span>
          </span>
          {#if currentFont === theme.id}
            <Check class="h-4 w-4 text-primary" />
          {/if}
        </span>
      </DropdownMenu.Item>
    {/each}
  </DropdownMenu.Content>
</DropdownMenu.Root>
