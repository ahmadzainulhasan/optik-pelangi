<script lang="ts">
	import type { LayoutProps } from '../../../$types';
	import { cn } from '#lib/utils.js';

	let { children }: LayoutProps = $props();
</script>

<svelte:head>
	<title>Detail Pengguna | Admin Optik Pelangi</title>
</svelte:head>

<div class={cn('flex flex-col gap-6', 'animate-in fade-in duration-500 ease-out')}>
	{@render children()}
</div>
