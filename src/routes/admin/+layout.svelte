<script lang="ts">
	import type { LayoutProps } from './$types';
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import * as Breadcrumb from '#lib/components/ui/breadcrumb/index.js';
	import AdminSidebar from '#lib/components/layout/admin-sidebar.svelte';
	import ModeToggle from '#lib/components/common/mode-toggle.svelte';
	import { page } from '$app/state';
	import { cn } from '#lib/utils.js';

	let { children }: LayoutProps = $props();

	const breadcrumbMap: Record<string, string> = {
		'/admin': 'Dashboard',
		'/admin/analytics': 'Analitik',
		'/admin/orders': 'Pesanan & Lab',
		'/admin/products': 'Katalog Produk',
		'/admin/inventory': 'Inventaris Stok',
		'/admin/users': 'Pengguna',
		'/admin/settings': 'Pengaturan',
		'/admin/sessions': 'Sesi Aktif'
	};

	const currentPageTitle = $derived(breadcrumbMap[page.url.pathname] ?? 'Dashboard');
</script>

<Sidebar.Provider>
	<AdminSidebar />
	<Sidebar.Inset>
		<header
			class={cn(
				'flex h-14 shrink-0 items-center justify-between gap-2 border-b px-4',
				'animate-in duration-500 ease-out fade-in slide-in-from-top-4'
			)}
		>
			<div class="flex items-center gap-2">
				<Sidebar.Trigger class="-ml-1" />

				<Breadcrumb.Root>
					<Breadcrumb.List>
						<Breadcrumb.Item class="hidden md:block">
							<Breadcrumb.Link href="/admin">Admin</Breadcrumb.Link>
						</Breadcrumb.Item>
						<Breadcrumb.Separator class="hidden md:block" />
						<Breadcrumb.Item>
							<Breadcrumb.Page>{currentPageTitle}</Breadcrumb.Page>
						</Breadcrumb.Item>
					</Breadcrumb.List>
				</Breadcrumb.Root>
			</div>
			<ModeToggle />
		</header>
		<div class="flex-1 p-4 md:p-6">
			{@render children()}
		</div>
	</Sidebar.Inset>
</Sidebar.Provider>
