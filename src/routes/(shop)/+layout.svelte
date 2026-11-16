<script lang="ts">
	import { setContext } from 'svelte';
	import type { LayoutProps } from './$types';
	import ShopHeader from '#lib/components/layout/shop-header.svelte';
	import ShopSidebar from '#lib/components/layout/shop-sidebar.svelte';
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import LoginDialog from './components/login-dialog.svelte';
	import House from '@lucide/svelte/icons/house';

	let { children }: LayoutProps = $props();

	const navItems = [
		{
			title: 'Beranda',
			href: '/',
			icon: House
		}
	];

	setContext('navItems', navItems);

	let loginDialog = $state(false);
	setContext('openLogin', () => (loginDialog = true));
</script>

<Sidebar.Provider open={false}>
	<ShopSidebar side="right" class="md:hidden" />
	<Sidebar.Inset>
		<ShopHeader />
		<div class="flex-1 p-6">
			{@render children()}
		</div>
	</Sidebar.Inset>
</Sidebar.Provider>

<LoginDialog bind:open={loginDialog} />
