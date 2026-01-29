<script lang="ts">
    import {
        Camera,
        Sparkles,
        TrendingUp,
        Zap,
        Upload,
        ChevronRight,
    } from "lucide-svelte";

    // File upload state
    let fileInput = $state<HTMLInputElement | null>(null);
    let isDragging = $state(false);

    function handleDragOver(e: DragEvent) {
        e.preventDefault();
        isDragging = true;
    }

    function handleDragLeave() {
        isDragging = false;
    }

    function handleDrop(e: DragEvent) {
        e.preventDefault();
        isDragging = false;
        // Handle file drop - will be implemented with API integration
        const files = e.dataTransfer?.files;
        if (files?.length) {
            console.log("Dropped file:", files[0].name);
        }
    }

    function openFileDialog() {
        fileInput?.click();
    }

    function handleFileSelect(e: Event) {
        const input = e.target as HTMLInputElement;
        const file = input.files?.[0];
        if (file) {
            console.log("Selected file:", file.name);
        }
    }

    const features = [
        {
            icon: Camera,
            title: "Instant Scan",
            description:
                "Snap a photo of your meal and get nutritional data in seconds",
        },
        {
            icon: Sparkles,
            title: "AI-Powered",
            description:
                "Advanced AI identifies foods and estimates macros accurately",
        },
        {
            icon: TrendingUp,
            title: "Track Progress",
            description:
                "Monitor your nutrition journey with detailed insights",
        },
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
        <p class="text-base-content/70 text-lg max-w-2xl mx-auto mb-8">
            Take a photo of any meal and let AI analyze the nutritional content.
            Track calories, protein, carbs, and fat effortlessly.
        </p>
    </section>

    <!-- Upload Zone -->
    <section class="mb-16">
        <div
            class="upload-zone p-8 md:p-12 text-center cursor-pointer transition-all duration-300 {isDragging
                ? 'border-primary bg-primary/10 scale-[1.02]'
                : ''}"
            role="button"
            tabindex="0"
            ondragover={handleDragOver}
            ondragleave={handleDragLeave}
            ondrop={handleDrop}
            onclick={openFileDialog}
            onkeydown={(e) => e.key === "Enter" && openFileDialog()}
        >
            <input
                bind:this={fileInput}
                type="file"
                accept="image/*"
                capture="environment"
                class="hidden"
                onchange={handleFileSelect}
            />

            <div class="flex flex-col items-center gap-4">
                <div
                    class="w-20 h-20 rounded-full bg-primary/10 flex items-center justify-center"
                >
                    <Upload class="w-10 h-10 text-primary" />
                </div>
                <div>
                    <h3 class="text-xl font-semibold mb-2">
                        Upload Your Meal Photo
                    </h3>
                    <p class="text-base-content/60">
                        Drag and drop an image here, or click to select
                    </p>
                </div>
                <button class="btn btn-primary btn-lg gap-2 mt-4">
                    <Camera class="w-5 h-5" />
                    Take Photo or Upload
                </button>
            </div>
        </div>
    </section>

    <!-- Features Grid -->
    <section class="mb-16">
        <h2 class="text-2xl font-bold text-center mb-8">How It Works</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            {#each features as feature, i}
                <div
                    class="card bg-base-100 shadow-lg hover:shadow-xl transition-shadow border border-base-200"
                >
                    <div class="card-body items-center text-center">
                        <div
                            class="w-14 h-14 rounded-full bg-gradient-to-br from-primary/20 to-accent/20 flex items-center justify-center mb-4"
                        >
                            <svelte:component
                                this={feature.icon}
                                class="w-7 h-7 text-primary"
                            />
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

    <!-- CTA Section -->
    <section class="text-center">
        <div
            class="card bg-gradient-to-r from-primary to-accent text-primary-content"
        >
            <div class="card-body py-12">
                <h2 class="card-title text-2xl justify-center mb-2">
                    Ready to Start Tracking?
                </h2>
                <p class="mb-6 opacity-90">
                    Join thousands who've simplified their nutrition journey
                </p>
                <div class="card-actions justify-center">
                    <a
                        href="/scan"
                        class="btn btn-outline btn-lg border-white/30 hover:bg-white/20 gap-2"
                    >
                        Get Started
                        <ChevronRight class="w-5 h-5" />
                    </a>
                </div>
            </div>
        </div>
    </section>
</div>
