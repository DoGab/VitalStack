<script lang="ts">
	import { Camera, Sparkles, TrendingUp, Zap } from "lucide-svelte";

	const features = [
		{
			icon: Camera,
			title: "Instant Scan",
			description: "Snap a photo of your meal and get nutritional data in seconds"
		},
		{
			icon: Sparkles,
			title: "AI-Powered",
			description: "Advanced AI identifies foods and estimates macros accurately"
		},
		{
			icon: TrendingUp,
			title: "Track Progress",
			description: "Monitor your nutrition journey with detailed insights"
		}
	];
</script>

<svelte:head>
	<title>MacroGuard - AI Food Nutrition Scanner</title>
	<meta
		name="description"
		content="Scan your food with AI to instantly get nutritional macros. Track calories, protein, carbs, and fat with just a photo."
	/>
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-4xl">
	<!-- Hero Section -->
	<section class="text-center mb-12">
		<div class="badge badge-primary badge-outline mb-4 gap-2">
			<Zap class="w-3 h-3" />
			AI-Powered Nutrition Tracking
		</div>
		<h1 class="text-4xl md:text-5xl font-bold mb-4 leading-tight">
			<span
				class="bg-gradient-to-r from-primary via-secondary to-accent bg-clip-text text-transparent"
			>
				Know Your Macros
			</span>
			<br />
			<span class="text-base-content">In Seconds</span>
		</h1>
		<p class="text-base-content/70 text-lg max-w-2xl mx-auto">
			Take a photo of any meal and let AI analyze the nutritional content. Track calories, protein,
			carbs, and fat effortlessly.
		</p>
	</section>

	<!-- Features Grid -->
	<section>
		<h2 class="text-2xl font-bold text-center mb-8">How It Works</h2>
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
			{#each features as feature, i}
				{@const Icon = feature.icon}
				<div
					class="card bg-base-100 shadow-lg hover:shadow-xl transition-shadow border border-base-200"
				>
					<div class="card-body items-center text-center">
						<div
							class="w-14 h-14 rounded-full bg-gradient-to-br from-primary/20 to-accent/20 flex items-center justify-center mb-4"
						>
							<Icon class="w-7 h-7 text-primary" />
						</div>
						<div class="badge badge-ghost badge-sm mb-2">
							Step {i + 1}
						</div>
						<h3 class="card-title text-lg">{feature.title}</h3>
						<p class="text-base-content/60 text-sm">
							{feature.description}
						</p>
					</div>
				</div>
			{/each}
		</div>
	</section>
</div>
