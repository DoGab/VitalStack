<script lang="ts">
  import { Sun, Moon } from "lucide-svelte";
  import { browser } from "$app/environment";

  // Dark = biohacker, Light = organic
  let isDark = $state(true);

  // Initialize from current theme
  $effect(() => {
    if (browser) {
      const currentTheme = document.documentElement.getAttribute("data-theme");
      isDark = currentTheme === "biohacker";
    }
  });

  function toggleTheme() {
    isDark = !isDark;
    const newTheme = isDark ? "biohacker" : "organic";
    document.documentElement.setAttribute("data-theme", newTheme);
    // Persist to localStorage
    localStorage.setItem("theme", newTheme);
  }

  // Load saved theme on mount
  $effect(() => {
    if (browser) {
      const savedTheme = localStorage.getItem("theme");
      if (savedTheme) {
        document.documentElement.setAttribute("data-theme", savedTheme);
        isDark = savedTheme === "biohacker";
      }
    }
  });
</script>

<label class="swap swap-rotate btn btn-ghost btn-circle">
  <input type="checkbox" checked={!isDark} onchange={toggleTheme} aria-label="Toggle theme" />

  <!-- Sun icon (shown when dark/biohacker is active) -->
  <Sun class="swap-off w-5 h-5 fill-current" />

  <!-- Moon icon (shown when light/organic is active) -->
  <Moon class="swap-on w-5 h-5 fill-current" />
</label>
