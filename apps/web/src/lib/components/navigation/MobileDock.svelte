<script lang="ts">
  import { Home, Clock, MessageCircle, User, Plus } from "lucide-svelte";
  import { page } from "$app/stores";

  interface Props {
    addMenuOpen: boolean;
  }

  let { addMenuOpen = $bindable() }: Props = $props();

  const navItems = [
    { href: "/", label: "Home", icon: Home },
    { href: "/history", label: "History", icon: Clock },
    { href: "/chat", label: "Chat", icon: MessageCircle },
    { href: "/profile", label: "Profile", icon: User }
  ];

  function isActive(href: string, pathname: string): boolean {
    if (href === "/") return pathname === "/";
    return pathname.startsWith(href);
  }
</script>

<!-- Mobile Bottom Navigation -->
<nav
  class="fixed bottom-0 left-0 right-0 z-50 flex items-center justify-around bg-background/90 backdrop-blur-md border-t border-border px-2 pb-safe md:hidden"
  style="padding-bottom: max(env(safe-area-inset-bottom), 0.5rem);"
>
  {#each navItems.slice(0, 2) as item (item.href)}
    {@const Icon = item.icon}
    {@const active = isActive(item.href, $page.url.pathname)}
    <a
      href={item.href}
      class="flex flex-col items-center justify-center gap-1 py-2 px-3 min-h-[48px] transition-colors {active
        ? 'text-primary'
        : 'text-muted-foreground'}"
    >
      <Icon class="h-5 w-5" />
      <span class="text-[10px] font-medium">{item.label}</span>
    </a>
  {/each}

  <!-- Center Add Button (elevated) -->
  <button
    class="flex items-center justify-center w-14 h-14 -translate-y-3 rounded-full bg-primary text-primary-foreground shadow-lg active:scale-95 transition-transform"
    onclick={() => (addMenuOpen = !addMenuOpen)}
    aria-label="Add entry"
  >
    <Plus class="h-7 w-7" />
  </button>

  {#each navItems.slice(2) as item (item.href)}
    {@const Icon = item.icon}
    {@const active = isActive(item.href, $page.url.pathname)}
    <a
      href={item.href}
      class="flex flex-col items-center justify-center gap-1 py-2 px-3 min-h-[48px] transition-colors {active
        ? 'text-primary'
        : 'text-muted-foreground'}"
    >
      <Icon class="h-5 w-5" />
      <span class="text-[10px] font-medium">{item.label}</span>
    </a>
  {/each}
</nav>
