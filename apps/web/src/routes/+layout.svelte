<script lang="ts">
  import "../app.css";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import { Separator } from "$lib/components/ui/separator/index.js";
  import AppSidebar from "$lib/components/navigation/AppSidebar.svelte";
  import MobileDock from "$lib/components/navigation/MobileDock.svelte";
  import AddEntryModal from "$lib/components/navigation/AddEntryModal.svelte";
  import { browser } from "$app/environment";
  import logoIconDark from "$lib/assets/logo/logo_dark.svg";
  import logoIconLight from "$lib/assets/logo/logo_light.svg";

  let { children } = $props();
  let addMenuOpen = $state(false);

  // Theme-aware logo switching for mobile header
  let isDarkTheme = $state(true);
  let mobileLogoIcon = $derived(isDarkTheme ? logoIconLight : logoIconDark);

  $effect(() => {
    if (browser) {
      const html = document.documentElement;
      const updateTheme = () => {
        const theme = html.getAttribute("data-theme");
        isDarkTheme = theme === "darkorganic" || !theme;
      };
      updateTheme();
      const observer = new MutationObserver(updateTheme);
      observer.observe(html, { attributes: true, attributeFilter: ["data-theme"] });
      return () => observer.disconnect();
    }
  });
</script>

<svelte:head>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
  <link
    href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Playfair+Display:wght@400;500;600;700&family=Outfit:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap"
    rel="stylesheet"
  />
  <meta name="theme-color" content="#1B3022" />
  <meta name="apple-mobile-web-app-capable" content="yes" />
  <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
</svelte:head>

<Sidebar.Provider>
  <AppSidebar bind:addMenuOpen />

  <Sidebar.Inset>
    <!-- Unified header: sidebar trigger + mobile logo -->
    <header
      class="bg-background/80 backdrop-blur-md sticky top-0 z-40 flex h-14 shrink-0 items-center gap-2 border-b border-border px-4 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12"
    >
      <Sidebar.Trigger class="-ms-1" />
      <Separator orientation="vertical" class="me-2 data-[orientation=vertical]:h-4" />

      <!-- Mobile logo: visible only below md -->
      <a href="/" class="flex items-center gap-2 md:hidden">
        <img src={mobileLogoIcon} alt="VitalStack" class="size-7" />
        <span class="font-semibold text-sm" style="font-family: var(--font-heading);"
          >VitalStack</span
        >
      </a>

      <div class="flex-1"></div>
    </header>

    <!-- Page Content -->
    <div class="pb-20 md:pb-0">
      {@render children()}
    </div>
  </Sidebar.Inset>
</Sidebar.Provider>

<!-- Add Entry Drawer -->
<AddEntryModal bind:open={addMenuOpen} />

<!-- Mobile Bottom Dock -->
<MobileDock bind:addMenuOpen />
