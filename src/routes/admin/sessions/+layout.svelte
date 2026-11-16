<script lang="ts">
	import type { LayoutProps } from './$types';
	import { cn } from '#lib/utils.js';

	let { children }: LayoutProps = $props();
</script>

<svelte:head>
	<title>Sesi Aktif | Admin Optik Pelangi</title>
</svelte:head>

<div class={cn('flex flex-col gap-6', 'animate-in fade-in slide-in-from-top-3 duration-500 ease-out')}>
	<div class="flex flex-col gap-1">
		<h1 class="text-2xl font-semibold tracking-tight text-foreground">Sesi Perangkat</h1>
		<p class="text-sm text-muted-foreground">
			Pantau semua perangkat yang sedang aktif memakai akun Anda. Cabut sesi apa pun yang tidak
			Anda kenali.
		</p>
	</div>

	{@render children()}
</div>
