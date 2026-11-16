<script lang="ts">
	import { getContext, type ComponentProps, type Component } from 'svelte';
	import { refreshAll } from '$app/navigation';
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import * as Avatar from '#lib/components/ui/avatar/index.js';
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import { authClient } from '#lib/auth-client.js';
	import { page } from '$app/state';
	import LogOut from '@lucide/svelte/icons/log-out';
	import ChevronsUpDown from '@lucide/svelte/icons/chevrons-up-down';
	import UserKey from '@lucide/svelte/icons/user-key';
	import { cn } from '#lib/utils.js';

	const openLogin = getContext<() => void>('openLogin');
	const sidebar = Sidebar.useSidebar();
	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();

	const user = $derived(page.data.currentUser);
	const navItems =
		getContext<Array<{ title: string; href: string; icon?: Component }>>('navItems') ?? [];

	async function handleSignOut() {
		sidebar.setOpenMobile(false);
		await authClient.signOut();
		await refreshAll();
	}
</script>

<Sidebar.Root
	bind:ref
	class={cn('animate-in duration-400 ease-out fade-in', restProps.class)}
	{...restProps}
>
	<Sidebar.Header>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton size="lg">
					{#snippet child({ props })}
						<a href="/" onclick={() => sidebar.setOpenMobile(false)} {...props}>
							<span class="font-serif text-lg font-semibold">Optik Pelangi</span>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Header>

	<Sidebar.Content>
		<Sidebar.Group>
			<Sidebar.GroupContent>
				<Sidebar.Menu>
					{#each navItems as item, i (item.title)}
						<Sidebar.MenuItem
							class="animate-in fade-in slide-in-from-left-2 duration-300 fill-mode-both"
							style="animation-delay: {i * 60}ms;"
						>
							<Sidebar.MenuButton
								isActive={page.url.pathname === item.href}
								onclick={() => sidebar.setOpenMobile(false)}
								tooltipContent={item.title}
							>
								{#snippet child({ props })}
									<a href={item.href} {...props}>
										{#if item.icon}
											<item.icon />
										{/if}
										<span>{item.title}</span>
									</a>
								{/snippet}
							</Sidebar.MenuButton>
						</Sidebar.MenuItem>
					{/each}
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group>
	</Sidebar.Content>

	<Sidebar.Footer>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				{#if user}
					<DropdownMenu.Root>
						<DropdownMenu.Trigger>
							{#snippet child({ props })}
								<Sidebar.MenuButton
									{...props}
									size="lg"
									class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
								>
									<div class="flex items-center gap-3 overflow-hidden">
										<Avatar.Root class="size-8 rounded-lg">
											<Avatar.Image src={user.image ?? undefined} alt={user.name} />
											<Avatar.Fallback class="rounded-lg" />
										</Avatar.Root>
										<div class="grid flex-1 text-sm leading-tight">
											<span class="truncate font-medium">{user.name}</span>
											<span class="truncate text-xs text-muted-foreground">{user.email}</span>
										</div>
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
										<a href="/sessions" {...props}>
											<UserKey class="mr-2 size-4" />
											<span>Sesi Aktif</span>
										</a>
									{/snippet}
								</DropdownMenu.Item>
								<DropdownMenu.Separator />
								<DropdownMenu.Item variant="destructive" onclick={handleSignOut}>
									<LogOut class="size-4" />Keluar
								</DropdownMenu.Item>
							</DropdownMenu.Group>
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				{:else}
					<Sidebar.MenuButton>
						{#snippet child({ props })}
							<Button
								class="w-full"
								{...props}
								onclick={() => {
									sidebar.setOpenMobile(false);
									openLogin();
								}}
							>
								Masuk
							</Button>
						{/snippet}
					</Sidebar.MenuButton>
				{/if}
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Footer>
	<Sidebar.Rail />
</Sidebar.Root>
