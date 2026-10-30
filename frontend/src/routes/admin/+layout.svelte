<script lang="ts">
	import type { LayoutProps } from './$types';
	import AdminSidebar from '#lib/components/custom/AdminSidebar.svelte';
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import * as Breadcrumb from '#lib/components/ui/breadcrumb/index.js';

	let { data, children }: LayoutProps = $props();
</script>

<svelte:head>
	<title>Panel Administrator — Optik Pelangi</title>
</svelte:head>

<Sidebar.Provider>
	<AdminSidebar user={data.user} />
	<Sidebar.Inset>
		<header class="flex h-14 items-center gap-2 border-b px-4">
			<Sidebar.Trigger />
			<Breadcrumb.Root>
				<Breadcrumb.List>
					<Breadcrumb.Item>
						<Breadcrumb.Page>Administrator</Breadcrumb.Page>
					</Breadcrumb.Item>
				</Breadcrumb.List>
			</Breadcrumb.Root>
		</header>

		<main class="p-4">
			{@render children()}
		</main>
	</Sidebar.Inset>
</Sidebar.Provider>
