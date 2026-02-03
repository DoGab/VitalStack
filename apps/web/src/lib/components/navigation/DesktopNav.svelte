<script lang="ts">
  import { Home, Clock, MessageCircle, User, Plus, Utensils } from "lucide-svelte";
  import ThemeSwap from "$lib/components/ui/ThemeSwap.svelte";

  let { addMenuOpen = $bindable() } = $props();

  const navLinks = [
    { href: "/", label: "Home", icon: Home },
    { href: "/history", label: "History", icon: Clock },
    { href: "/chat", label: "Chat", icon: MessageCircle },
    { href: "/profile", label: "Profile", icon: User }
  ];
</script>

<header
  class="navbar bg-base-100/90 backdrop-blur-md sticky top-0 z-50 border-b border-base-200 shadow-sm"
>
  <div class="navbar-start">
    <a href="/" class="btn btn-ghost text-xl font-bold gap-2">
      <Utensils class="w-6 h-6 text-primary" />
      <span class="bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">
        MacroGuard
      </span>
    </a>
  </div>

  <!-- Desktop Nav Links -->
  <div class="navbar-center hidden lg:flex">
    <ul class="menu menu-horizontal px-1 gap-1">
      {#each navLinks as link (link.href)}
        {@const Icon = link.icon}
        <li>
          <a href={link.href} class="btn btn-ghost btn-sm gap-2">
            <Icon class="w-4 h-4" />
            {link.label}
          </a>
        </li>
      {/each}
    </ul>
  </div>

  <div class="navbar-end gap-2">
    <!-- Theme Swap -->
    <ThemeSwap />

    <!-- Desktop: Add button -->
    <div class="hidden lg:flex">
      <button class="btn btn-primary btn-sm gap-2" onclick={() => (addMenuOpen = !addMenuOpen)}>
        <Plus class="w-4 h-4" />
        Add
      </button>
    </div>
  </div>
</header>
