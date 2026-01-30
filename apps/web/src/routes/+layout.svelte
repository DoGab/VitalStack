<script lang="ts">
	import "../app.css";
	import {
		Home,
		Clock,
		Plus,
		MessageCircle,
		User,
		Utensils,
		Camera,
		Upload,
		Zap,
		X
	} from "lucide-svelte";

	let { children } = $props();
	let addMenuOpen = $state(false);

	const navLinks = [
		{ href: "/", label: "Home", icon: Home },
		{ href: "/history", label: "History", icon: Clock },
		{ href: "/chat", label: "Chat", icon: MessageCircle },
		{ href: "/profile", label: "Profile", icon: User }
	];

	const addOptions = [
		{ icon: Camera, label: "Take Photo", action: "camera" },
		{ icon: Upload, label: "Upload Image", action: "upload" },
		{ icon: Zap, label: "Quick Add", action: "quick" }
	];

	function handleAddOption(action: string) {
		console.log("Add action:", action);
		addMenuOpen = false;
		// Placeholder - will be implemented later
	}
</script>

<svelte:head>
	<link rel="preconnect" href="https://fonts.googleapis.com" />
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
	<link
		href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap"
		rel="stylesheet"
	/>
	<meta name="theme-color" content="#22c55e" />
	<meta name="apple-mobile-web-app-capable" content="yes" />
	<meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
</svelte:head>

<div class="min-h-screen bg-base-100 flex flex-col">
	<!-- Navbar -->
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

		<!-- Desktop Nav -->
		<div class="navbar-center hidden lg:flex">
			<ul class="menu menu-horizontal px-1 gap-1">
				{#each navLinks as link}
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

		<div class="navbar-end">
			<!-- Desktop: Add button -->
			<div class="hidden lg:flex">
				<button class="btn btn-primary btn-sm gap-2" onclick={() => (addMenuOpen = !addMenuOpen)}>
					<Plus class="w-4 h-4" />
					Add
				</button>
			</div>
		</div>
	</header>

	<!-- Add Options Modal/Dropdown -->
	{#if addMenuOpen}
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div class="fixed inset-0 bg-black/50 z-[60]" onclick={() => (addMenuOpen = false)}>
			<div
				class="fixed bottom-24 lg:top-16 lg:bottom-auto lg:right-4 left-1/2 lg:left-auto -translate-x-1/2 lg:translate-x-0 w-64 bg-base-100 rounded-2xl shadow-2xl p-4 z-[70]"
				onclick={(e) => e.stopPropagation()}
			>
				<div class="flex justify-between items-center mb-3">
					<h3 class="font-semibold">Add Entry</h3>
					<button class="btn btn-ghost btn-sm btn-circle" onclick={() => (addMenuOpen = false)}>
						<X class="w-4 h-4" />
					</button>
				</div>
				<div class="flex flex-col gap-2">
					{#each addOptions as option}
						{@const Icon = option.icon}
						<button
							class="btn btn-ghost justify-start gap-3 h-14"
							onclick={() => handleAddOption(option.action)}
						>
							<div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
								<Icon class="w-5 h-5 text-primary" />
							</div>
							{option.label}
						</button>
					{/each}
				</div>
			</div>
		</div>
	{/if}

	<!-- Main Content -->
	<main class="flex-1 pb-20 lg:pb-0">
		{@render children()}
	</main>

	<!-- Mobile Bottom Nav -->
	<nav
		class="fixed bottom-0 left-0 right-0 lg:hidden border-t border-base-200 bg-base-100/90 backdrop-blur-md z-50"
	>
		<div class="flex items-center justify-around h-16">
			<!-- Home -->
			<a
				href="/"
				class="flex flex-col items-center gap-1 text-base-content/70 hover:text-primary transition-colors min-w-[48px]"
			>
				<Home class="w-5 h-5" />
				<span class="text-xs">Home</span>
			</a>

			<!-- History -->
			<a
				href="/history"
				class="flex flex-col items-center gap-1 text-base-content/70 hover:text-primary transition-colors min-w-[48px]"
			>
				<Clock class="w-5 h-5" />
				<span class="text-xs">History</span>
			</a>

			<!-- Center Add Button (elevated) -->
			<button
				class="w-14 h-14 -mt-6 rounded-full bg-primary text-primary-content shadow-lg flex items-center justify-center active:scale-95 transition-transform"
				onclick={() => (addMenuOpen = !addMenuOpen)}
				aria-label="Add entry"
			>
				<Plus class="w-7 h-7" />
			</button>

			<!-- Chat -->
			<a
				href="/chat"
				class="flex flex-col items-center gap-1 text-base-content/70 hover:text-primary transition-colors min-w-[48px]"
			>
				<MessageCircle class="w-5 h-5" />
				<span class="text-xs">Chat</span>
			</a>

			<!-- Profile -->
			<a
				href="/profile"
				class="flex flex-col items-center gap-1 text-base-content/70 hover:text-primary transition-colors min-w-[48px]"
			>
				<User class="w-5 h-5" />
				<span class="text-xs">Profile</span>
			</a>
		</div>
	</nav>
</div>
