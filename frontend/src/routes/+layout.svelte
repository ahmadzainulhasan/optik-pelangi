<script lang="ts">
	import './layout.css';
	import favicon from '#lib/assets/favicon.svg';
	import type { LayoutProps } from './$types';
	import AppHeader from '#lib/components/custom/AppHeader.svelte';
	import AppSidebar from '#lib/components/custom/AppSidebar.svelte';
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';

	let { children }: LayoutProps = $props();

	$effect(() => {
		if (typeof window !== 'undefined' && window.name === 'GoogleAuthPopup' && window.opener) {
			window.opener.postMessage({ type: 'AUTH_SUCCESS' }, '*');
			window.close();
		}
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Sidebar.Provider>
	<div class="md:hidden">
		<AppSidebar />
	</div>
	<div class="flex min-h-screen w-full flex-col">
		<AppHeader />
		{@render children()}
	</div>
</Sidebar.Provider>
