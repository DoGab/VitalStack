<script lang="ts">
  import { Camera, Sparkles, TrendingUp, Zap } from "lucide-svelte";
  import { Badge } from "$lib/components/ui/badge";
  import * as Card from "$lib/components/ui/card";

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
  <title>VitalStack - AI Food Nutrition Scanner</title>
  <meta
    name="description"
    content="Scan your food with AI to instantly get nutritional macros. Track calories, protein, carbs, and fat with just a photo."
  />
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-4xl">
  <!-- Hero Section -->
  <section class="text-center mb-12">
    <Badge variant="outline" class="mb-4 gap-2">
      <Zap class="h-3 w-3" />
      AI-Powered Nutrition Tracking
    </Badge>
    <h1
      class="text-4xl md:text-5xl font-bold mb-4 leading-tight"
      style="font-family: var(--font-heading);"
    >
      <span
        class="bg-gradient-to-r from-primary via-secondary to-accent bg-clip-text text-transparent"
      >
        Know Your Macros
      </span>
      <br />
      <span class="text-foreground">In Seconds</span>
    </h1>
    <p class="text-muted-foreground text-lg max-w-2xl mx-auto">
      Take a photo of any meal and let AI analyze the nutritional content. Track calories, protein,
      carbs, and fat effortlessly.
    </p>
  </section>

  <!-- Features Grid -->
  <section>
    <h2 class="text-2xl font-bold text-center mb-8" style="font-family: var(--font-heading);">
      How It Works
    </h2>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      {#each features as feature, i (feature.title)}
        {@const Icon = feature.icon}
        <Card.Root class="hover:shadow-lg transition-shadow">
          <Card.Content class="flex flex-col items-center text-center pt-6">
            <div
              class="w-14 h-14 rounded-full bg-gradient-to-br from-primary/20 to-accent/20 flex items-center justify-center mb-4"
            >
              <Icon class="h-7 w-7 text-primary" />
            </div>
            <Badge variant="secondary" class="mb-2 text-xs">
              Step {i + 1}
            </Badge>
            <h3 class="text-lg font-semibold mb-1">{feature.title}</h3>
            <p class="text-muted-foreground text-sm">
              {feature.description}
            </p>
          </Card.Content>
        </Card.Root>
      {/each}
    </div>
  </section>
</div>
