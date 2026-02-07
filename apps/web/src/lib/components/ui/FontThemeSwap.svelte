<script lang="ts">
  import { Type } from "lucide-svelte";

  // Font themes available
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

<div class="dropdown dropdown-end">
  <div tabindex="0" role="button" class="btn btn-ghost btn-sm gap-2" aria-label="Font Theme">
    <Type class="w-4 h-4" />
    <span class="hidden sm:inline text-xs">{currentFont === "classic" ? "Aa" : "Aa"}</span>
  </div>
  <ul
    tabindex="0"
    class="dropdown-content menu bg-base-200 rounded-box z-50 w-52 p-2 shadow-lg mt-2"
  >
    <li class="menu-title">
      <span>Font Theme</span>
    </li>
    {#each fontThemes as theme (theme.id)}
      <li>
        <button
          class="flex justify-between"
          class:active={currentFont === theme.id}
          onclick={() => setFontTheme(theme.id)}
        >
          <span class="font-medium">{theme.label}</span>
          <span class="text-xs opacity-60">{theme.description}</span>
        </button>
      </li>
    {/each}
  </ul>
</div>
