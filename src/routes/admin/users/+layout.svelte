<script lang="ts">
	import type { LayoutProps } from '../../$types';
	import UserToolbar from './components/user-toolbar.svelte';
	import { cn } from '#lib/utils.js';

	let { children }: LayoutProps = $props();
</script>

<svelte:head>
	<title>Pengguna | Admin Optik Pelangi</title>
</svelte:head>

<div class={cn('flex flex-col gap-4', 'animate-in fade-in slide-in-from-top-3 duration-500 ease-out')}>
	<div>
		<h1 class="text-xl font-semibold">Pengguna</h1>
		<p class="text-xs text-muted-foreground">
			Pengguna terdaftar, status akun, dan moderasi akses toko Optik Pelangi.
		</p>
	</div>

	<UserToolbar />

	{@render children()}
</div>
