<script lang="ts">
  import {
    Home,
    Clock,
    MessageCircle,
    User,
    Plus,
    ChevronsUpDown,
    LogOut,
    Sun,
    Moon,
    Type
  } from "lucide-svelte";
  import { page } from "$app/stores";
  import { browser } from "$app/environment";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import * as Popover from "$lib/components/ui/popover";
  import * as Avatar from "$lib/components/ui/avatar";
  import logoIconDark from "$lib/assets/logo/logo_dark.svg";
  import logoIconLight from "$lib/assets/logo/logo_light.svg";

  interface Props {
    addMenuOpen: boolean;
  }

  let { addMenuOpen = $bindable() }: Props = $props();

  // Track current theme state for logo + theme toggle
  let isDarkTheme = $state(true);
  let currentLogoIcon = $derived(isDarkTheme ? logoIconLight : logoIconDark);

  // Track current font theme for font toggle label
  let currentFont = $state<"classic" | "modern">("classic");
  let currentFontLabel = $derived(currentFont === "classic" ? "Classic" : "Modern");

  $effect(() => {
    if (browser) {
      const html = document.documentElement;

      // Theme observer
      const updateTheme = () => {
        const theme = html.getAttribute("data-theme");
        isDarkTheme = theme === "darkorganic" || !theme;
      };
      updateTheme();

      // Font observer
      const updateFont = () => {
        const font = html.getAttribute("data-font-theme");
        currentFont = font === "modern" ? "modern" : "classic";
      };
      updateFont();

      const observer = new MutationObserver(() => {
        updateTheme();
        updateFont();
      });
      observer.observe(html, {
        attributes: true,
        attributeFilter: ["data-theme", "data-font-theme"]
      });
      return () => observer.disconnect();
    }
  });

  function toggleTheme() {
    isDarkTheme = !isDarkTheme;
    const newTheme = isDarkTheme ? "darkorganic" : "organic";
    document.documentElement.setAttribute("data-theme", newTheme);
    localStorage.setItem("theme", newTheme);
  }

  function toggleFont() {
    currentFont = currentFont === "classic" ? "modern" : "classic";
    document.documentElement.setAttribute("data-font-theme", currentFont);
    localStorage.setItem("font-theme", currentFont);
  }

  const navItems = [
    { href: "/", label: "Home", icon: Home },
    { href: "/history", label: "History", icon: Clock },
    { href: "/chat", label: "Chat", icon: MessageCircle },
    { href: "/profile", label: "Profile", icon: User }
  ];

  // Hardcoded user (no auth yet)
  const user = {
    name: "User",
    email: "user@vitalstack.app",
    avatar: ""
  };

  function isActive(href: string, pathname: string): boolean {
    if (href === "/") return pathname === "/";
    return pathname.startsWith(href);
  }
</script>

<Sidebar.Root collapsible="icon">
  <!-- Header: Logo -->
  <Sidebar.Header>
    <Sidebar.Menu>
      <Sidebar.MenuItem>
        <Sidebar.MenuButton size="lg">
          {#snippet child({ props })}
            <a href="/" {...props} class="flex items-center gap-2">
              <div class="flex aspect-square size-8 items-center justify-center shrink-0">
                <img src={currentLogoIcon} alt="VitalStack" class="size-8" />
              </div>
              <span
                class="truncate font-semibold group-data-[collapsible=icon]:hidden"
                style="font-family: var(--font-heading);">VitalStack</span
              >
            </a>
          {/snippet}
        </Sidebar.MenuButton>
      </Sidebar.MenuItem>
    </Sidebar.Menu>
  </Sidebar.Header>

  <!-- Navigation Content -->
  <Sidebar.Content>
    <Sidebar.Group>
      <Sidebar.GroupLabel>Navigation</Sidebar.GroupLabel>
      <Sidebar.GroupContent>
        <Sidebar.Menu>
          {#each navItems as item (item.href)}
            {@const Icon = item.icon}
            {@const active = isActive(item.href, $page.url.pathname)}
            <Sidebar.MenuItem>
              <Sidebar.MenuButton isActive={active} tooltipContent={item.label}>
                {#snippet child({ props })}
                  <a href={item.href} {...props}>
                    <Icon />
                    <span>{item.label}</span>
                  </a>
                {/snippet}
              </Sidebar.MenuButton>
            </Sidebar.MenuItem>
          {/each}
        </Sidebar.Menu>
      </Sidebar.GroupContent>
    </Sidebar.Group>

    <!-- Add Entry Button -->
    <Sidebar.Group>
      <Sidebar.GroupContent>
        <Sidebar.Menu>
          <Sidebar.MenuItem>
            <Sidebar.MenuButton tooltipContent="Add Entry" class="text-primary font-semibold">
              {#snippet child({ props })}
                <button {...props} onclick={() => (addMenuOpen = !addMenuOpen)}>
                  <Plus />
                  <span>Add Entry</span>
                </button>
              {/snippet}
            </Sidebar.MenuButton>
          </Sidebar.MenuItem>
        </Sidebar.Menu>
      </Sidebar.GroupContent>
    </Sidebar.Group>
  </Sidebar.Content>

  <!-- Footer: User Profile -->
  <Sidebar.Footer>
    <Sidebar.Menu>
      <Sidebar.MenuItem>
        <Popover.Root>
          <Popover.Trigger>
            {#snippet child({ props })}
              <Sidebar.MenuButton
                size="lg"
                class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
              >
                {#snippet child({ props: btnProps })}
                  <button {...props} {...btnProps} class="flex items-center gap-2 w-full">
                    <Avatar.Root class="h-8 w-8 rounded-lg">
                      <Avatar.Fallback class="rounded-lg bg-secondary text-secondary-foreground"
                        >{user.name.charAt(0).toUpperCase()}</Avatar.Fallback
                      >
                    </Avatar.Root>
                    <div
                      class="grid flex-1 text-left text-sm leading-tight group-data-[collapsible=icon]:hidden"
                    >
                      <span class="truncate font-semibold">{user.name}</span>
                      <span class="truncate text-xs text-muted-foreground">{user.email}</span>
                    </div>
                    <ChevronsUpDown class="ml-auto size-4 group-data-[collapsible=icon]:hidden" />
                  </button>
                {/snippet}
              </Sidebar.MenuButton>
            {/snippet}
          </Popover.Trigger>
          <Popover.Content
            class="w-[--bits-popover-anchor-width] min-w-56 rounded-lg p-0"
            side="top"
            align="start"
            sideOffset={4}
          >
            <!-- User info -->
            <div class="flex items-center gap-2 px-3 py-3 text-left text-sm border-b border-border">
              <Avatar.Root class="h-8 w-8 rounded-lg">
                <Avatar.Fallback class="rounded-lg bg-secondary text-secondary-foreground"
                  >{user.name.charAt(0).toUpperCase()}</Avatar.Fallback
                >
              </Avatar.Root>
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-semibold">{user.name}</span>
                <span class="truncate text-xs text-muted-foreground">{user.email}</span>
              </div>
            </div>

            <!-- Theme & Font picker (full-width buttons matching Log out style) -->
            <div class="p-2 space-y-1">
              <button
                onclick={toggleTheme}
                class="flex w-full items-center gap-2 rounded-sm px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground transition-colors"
              >
                {#if isDarkTheme}
                  <Sun class="size-4" />
                  <span>Light Mode</span>
                {:else}
                  <Moon class="size-4" />
                  <span>Dark Mode</span>
                {/if}
              </button>
              <button
                onclick={toggleFont}
                class="flex w-full items-center gap-2 rounded-sm px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground transition-colors"
              >
                <Type class="size-4" />
                <span>Font: {currentFontLabel}</span>
              </button>
            </div>

            <!-- Log out -->
            <div class="border-t border-border p-2">
              <button
                class="flex w-full items-center gap-2 rounded-sm px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground transition-colors"
              >
                <LogOut class="size-4" />
                <span>Log out</span>
              </button>
            </div>
          </Popover.Content>
        </Popover.Root>
      </Sidebar.MenuItem>
    </Sidebar.Menu>
  </Sidebar.Footer>
</Sidebar.Root>
