<script lang="ts">
	import type { PageProps } from './$types';
	import { Separator } from '#lib/components/ui/separator/index.js';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>{data.title} — Optik Pelangi</title>
</svelte:head>

<main class="min-h-screen bg-background pb-20">
	<header class="bg-muted py-14 md:py-24">
		<div class="container mx-auto px-4 text-center md:px-6">
			<h1 class="font-serif text-4xl md:text-5xl lg:text-6xl">{data.title}</h1>
			<p class="mt-4 text-sm text-muted-foreground md:text-base">
				Terakhir Diperbarui: {data.lastUpdated}
			</p>
		</div>
	</header>

	<div class="container mx-auto max-w-3xl px-4 py-12 md:px-6">
		<article class="prose prose-gray dark:prose-invert max-w-none">
			{#each data.sections as section, i}
				<div class="mb-8">
					<h2 class="font-serif text-2xl font-medium">{section.title}</h2>
					<div class="mt-4 flex flex-col gap-3 text-muted-foreground">
						{#each section.content as paragraph}
							<p class="leading-relaxed">{paragraph}</p>
						{/each}
					</div>
					{#if i < data.sections.length - 1}
						<Separator class="mt-8" />
					{/if}
				</div>
			{/each}
		</article>
	</div>
</main>

