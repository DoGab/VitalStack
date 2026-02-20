<script lang="ts">
  import { Sun, Moon } from "lucide-svelte";
  import { browser } from "$app/environment";
  import { Button } from "$lib/components/ui/button";

  // Dark = darkorganic, Light = organic
  let isDark = $state(true);

  // Load saved theme on mount
  $effect(() => {
    if (browser) {
      const savedTheme = localStorage.getItem("theme");
      if (savedTheme) {
        document.documentElement.setAttribute("data-theme", savedTheme);
        isDark = savedTheme === "darkorganic";
      } else {
        // Default to dark
        const currentTheme = document.documentElement.getAttribute("data-theme");
        isDark = currentTheme === "darkorganic" || !currentTheme;
      }
    }
  });

  function toggleTheme() {
    isDark = !isDark;
    const newTheme = isDark ? "darkorganic" : "organic";
    document.documentElement.setAttribute("data-theme", newTheme);
    localStorage.setItem("theme", newTheme);
  }
</script>

<Button
  variant="ghost"
  size="icon"
  onclick={toggleTheme}
  aria-label="Toggle theme"
  class="rounded-full"
>
  {#if isDark}
    <Sun class="h-5 w-5" />
  {:else}
    <Moon class="h-5 w-5" />
  {/if}
</Button>
