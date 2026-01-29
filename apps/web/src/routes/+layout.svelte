<script lang="ts">
    import "../app.css";
    import {
        Camera,
        Home,
        BarChart3,
        User,
        Menu,
        Utensils,
    } from "lucide-svelte";

    let { children } = $props();
    let mobileMenuOpen = $state(false);

    const navLinks = [
        { href: "/", label: "Home", icon: Home },
        { href: "/scan", label: "Scan", icon: Camera },
        { href: "/history", label: "History", icon: BarChart3 },
        { href: "/profile", label: "Profile", icon: User },
    ];
</script>

<svelte:head>
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link
        rel="preconnect"
        href="https://fonts.gstatic.com"
        crossorigin="anonymous"
    />
    <link
        href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap"
        rel="stylesheet"
    />
    <meta name="theme-color" content="#10b981" />
    <meta name="apple-mobile-web-app-capable" content="yes" />
    <meta
        name="apple-mobile-web-app-status-bar-style"
        content="black-translucent"
    />
</svelte:head>

<div class="min-h-screen bg-base-100 flex flex-col">
    <!-- Navbar -->
    <header
        class="navbar bg-base-100/90 backdrop-blur-md sticky top-0 z-50 border-b border-base-200 shadow-sm"
    >
        <div class="navbar-start">
            <a href="/" class="btn btn-ghost text-xl font-bold gap-2">
                <Utensils class="w-6 h-6 text-primary" />
                <span
                    class="bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent"
                >
                    MacroGuard
                </span>
            </a>
        </div>

        <!-- Desktop Nav -->
        <div class="navbar-center hidden lg:flex">
            <ul class="menu menu-horizontal px-1 gap-1">
                {#each navLinks as link}
                    <li>
                        <a href={link.href} class="btn btn-ghost btn-sm gap-2">
                            <svelte:component
                                this={link.icon}
                                class="w-4 h-4"
                            />
                            {link.label}
                        </a>
                    </li>
                {/each}
            </ul>
        </div>

        <div class="navbar-end">
            <!-- Desktop: Theme toggle placeholder -->
            <div class="hidden lg:flex">
                <a href="/scan" class="btn btn-primary btn-sm gap-2">
                    <Camera class="w-4 h-4" />
                    Scan Food
                </a>
            </div>

            <!-- Mobile: Hamburger menu -->
            <div class="lg:hidden">
                <button
                    class="btn btn-ghost btn-square"
                    onclick={() => (mobileMenuOpen = !mobileMenuOpen)}
                    aria-label="Toggle menu"
                >
                    <Menu class="w-6 h-6" />
                </button>
            </div>
        </div>
    </header>

    <!-- Mobile Drawer Menu -->
    {#if mobileMenuOpen}
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div
            class="fixed inset-0 bg-black/50 z-40 lg:hidden"
            onclick={() => (mobileMenuOpen = false)}
        >
            <div
                class="absolute right-0 top-0 h-full w-64 bg-base-100 shadow-xl p-4"
                onclick={(e) => e.stopPropagation()}
            >
                <ul class="menu gap-2">
                    {#each navLinks as link}
                        <li>
                            <a
                                href={link.href}
                                class="btn btn-ghost justify-start gap-3"
                                onclick={() => (mobileMenuOpen = false)}
                            >
                                <svelte:component
                                    this={link.icon}
                                    class="w-5 h-5"
                                />
                                {link.label}
                            </a>
                        </li>
                    {/each}
                </ul>
            </div>
        </div>
    {/if}

    <!-- Main Content -->
    <main class="flex-1">
        {@render children()}
    </main>

    <!-- Mobile Bottom Nav -->
    <nav
        class="btm-nav btm-nav-md lg:hidden border-t border-base-200 bg-base-100/90 backdrop-blur-md"
    >
        {#each navLinks as link}
            <a
                href={link.href}
                class="text-base-content/70 hover:text-primary transition-colors"
            >
                <svelte:component this={link.icon} class="w-5 h-5" />
                <span class="btm-nav-label text-xs">{link.label}</span>
            </a>
        {/each}
    </nav>
</div>
