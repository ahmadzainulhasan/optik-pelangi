<script lang="ts">
	import type { LayoutProps } from './$types';
	import AppHeader from '#lib/components/custom/AppHeader.svelte';
	import AppSidebar from '#lib/components/custom/AppSidebar.svelte';
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';

	let { data, children }: LayoutProps = $props();
</script>

<Sidebar.Provider>
	<div class="md:hidden">
		<AppSidebar user={data.user} />
	</div>
	<div class="flex min-h-screen w-full flex-col">
		<AppHeader user={data.user} />
		{@render children()}
	</div>
</Sidebar.Provider>
