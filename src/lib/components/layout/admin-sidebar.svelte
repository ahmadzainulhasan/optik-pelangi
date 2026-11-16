<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import * as Avatar from '#lib/components/ui/avatar/index.js';
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu/index.js';
	import { Skeleton } from '#lib/components/ui/skeleton/index.js';
	import { page } from '$app/state';
	import { authClient } from '#lib/auth-client.js';
	import { goto } from '$app/navigation';
	import LayoutDashboard from '@lucide/svelte/icons/layout-dashboard';
	import Users from '@lucide/svelte/icons/users';
	import LogOut from '@lucide/svelte/icons/log-out';
	import UserKey from '@lucide/svelte/icons/user-key';
	import ChevronsUpDown from '@lucide/svelte/icons/chevrons-up-down';
	import { cn } from '#lib/utils.js';

	const user = $derived(page.data.currentUser);
	const sidebar = Sidebar.useSidebar();
	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();

	const navGroups = [
		{
			title: 'Ringkasan',
			items: [
				{
					title: 'Dashboard',
					href: '/admin',
					icon: LayoutDashboard,
					get isActive() {
						return page.url.pathname === '/admin';
					}
				}
			]
		},
		{
			title: 'Operasional Toko',
			items: [
				{
					title: 'Pengguna',
					href: '/admin/users',
					icon: Users,
					get isActive() {
						return page.url.pathname.startsWith('/admin/users');
					}
				}
			]
		}
	];

	async function handleSignOut() {
		await authClient.signOut();
		await goto('/');
	}

	function getItemDelay(groupIndex: number, itemIndex: number) {
		let count = 0;
		for (let i = 0; i < groupIndex; i++) {
			count += navGroups[i].items.length;
		}
		return (count + itemIndex) * 60;
	}
</script>

<Sidebar.Root
	bind:ref
	class={cn('animate-in duration-500 ease-out fade-in slide-in-from-left-4', restProps.class)}
	{...restProps}
>
	<Sidebar.Header>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton size="lg">
					{#snippet child({ props })}
						<a href="/admin" onclick={() => sidebar.setOpenMobile(false)} {...props}>
							<span class="font-serif text-lg font-semibold">Optik Pelangi</span>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Header>

	<Sidebar.Content>
		{#each navGroups as group, groupIndex (group.title)}
			<Sidebar.Group>
				<Sidebar.GroupLabel>{group.title}</Sidebar.GroupLabel>
				<Sidebar.GroupContent>
					<Sidebar.Menu>
						{#each group.items as item, itemIndex (item.title)}
							<Sidebar.MenuItem
								class="animate-in fade-in slide-in-from-left-2 duration-300 fill-mode-both"
								style="animation-delay: {getItemDelay(groupIndex, itemIndex)}ms;"
							>
								<Sidebar.MenuButton
									isActive={item.isActive}
									tooltipContent={item.title}
								>
									{#snippet child({ props })}
										<a href={item.href} onclick={() => sidebar.setOpenMobile(false)} {...props}>
											<item.icon />
											<span>{item.title}</span>
										</a>
									{/snippet}
								</Sidebar.MenuButton>
							</Sidebar.MenuItem>
						{/each}
					</Sidebar.Menu>
				</Sidebar.GroupContent>
			</Sidebar.Group>
		{/each}
	</Sidebar.Content>

	<Sidebar.Footer>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						{#snippet child({ props })}
							<Sidebar.MenuButton
								{...props}
								size="lg"
								class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
							>
								<div class="flex items-center gap-3 overflow-hidden">
									{#await user}
										<Skeleton class="size-8 rounded-lg" />
										<div class="grid flex-1 gap-1">
											<Skeleton class="h-3.5 w-24" />
											<Skeleton class="h-3 w-32" />
										</div>
									{:then currentUser}
										<Avatar.Root>
											<Avatar.Image src={currentUser?.image} alt={currentUser?.name} />
											<Avatar.Fallback />
										</Avatar.Root>
										<div class="grid flex-1 text-sm leading-tight">
											<span class="truncate font-medium">{currentUser?.name}</span>
											<span class="truncate text-xs text-muted-foreground"
												>{currentUser?.email}</span
											>
										</div>
									{/await}
								</div>
								<ChevronsUpDown class="ms-auto size-4" />
							</Sidebar.MenuButton>
						{/snippet}
					</DropdownMenu.Trigger>
					<DropdownMenu.Content
						class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
						side="top"
						align="end"
						sideOffset={4}
					>
						<DropdownMenu.Group>
							<DropdownMenu.Item onSelect={() => sidebar.setOpenMobile(false)}>
								{#snippet child({ props })}
									<a href="/admin/sessions" {...props}>
										<UserKey class="size-4" />
										<span>Sesi Aktif</span>
									</a>
								{/snippet}
							</DropdownMenu.Item>
						</DropdownMenu.Group>
						<DropdownMenu.Separator />
						<DropdownMenu.Group>
							<DropdownMenu.Item variant="destructive" onclick={handleSignOut}>
								<LogOut class="size-4" />Keluar
							</DropdownMenu.Item>
						</DropdownMenu.Group>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Footer>

	<Sidebar.Rail />
</Sidebar.Root>
