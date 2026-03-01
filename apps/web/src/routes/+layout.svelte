<script lang="ts">
  import "../app.css";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import AppSidebar from "$lib/components/navigation/AppSidebar.svelte";
  import MobileDock from "$lib/components/navigation/MobileDock.svelte";
  import AddEntryModal from "$lib/components/navigation/AddEntryModal.svelte";
  import { page } from "$app/stores";

  let { data, children } = $props();
  let addMenuOpen = $state(false);

  let pageTitle = $derived.by(() => {
    const path = $page.url.pathname;
    if (path === "/") return "Dashboard";
    if (path.startsWith("/history")) return "History";
    if (path.startsWith("/chat")) return "Chat";
    if (path.startsWith("/profile")) return "Profile";
    return "Dashboard";
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
  <meta name="mobile-web-app-capable" content="yes" />
  <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
</svelte:head>

<Sidebar.Provider open={data.sidebarOpen}>
  <AppSidebar bind:addMenuOpen />

  <Sidebar.Inset>
    <!-- Unified header: sidebar trigger + mobile logo -->
    <header
      class="bg-background/80 backdrop-blur-md sticky top-0 z-40 flex h-14 shrink-0 items-center gap-2 px-4 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12"
    >
      <Sidebar.Trigger class="-ms-1" />

      <!-- Display current page name in header -->
      <h1 class="font-semibold text-lg ml-2" style="font-family: var(--font-heading);">
        {pageTitle}
      </h1>

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
