<script lang="ts">
	import { getContext } from 'svelte';
	import { refreshAll } from '$app/navigation';
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu/index.js';
	import * as Avatar from '#lib/components/ui/avatar/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import ShopNavbar from '#lib/components/layout/shop-navbar.svelte';
	import ModeToggle from '#lib/components/common/mode-toggle.svelte';
	import { authClient } from '#lib/auth-client.js';
	import { page } from '$app/state';
	import LogOut from '@lucide/svelte/icons/log-out';
	import UserKey from '@lucide/svelte/icons/user-key';
	import { cn } from '#lib/utils.js';

	let scrollY = $state(0);
	const user = $derived(page.data.currentUser);
	const openLogin = getContext<() => void>('openLogin');

	async function handleSignOut() {
		await authClient.signOut();
		await refreshAll();
	}
</script>

<svelte:window bind:scrollY />

<header
	class={cn(
		'sticky top-0 z-40 flex items-center justify-between px-2 py-1.5 md:px-4',
		'animate-in duration-500 ease-out fade-in slide-in-from-top-4',
		'transition-[background-color,box-shadow] duration-300 ease-out',
		scrollY > 0 && 'bg-background shadow-sm'
	)}
>
	<a href="/" class="font-serif text-lg font-semibold">Optik Pelangi</a>
	<ShopNavbar class="hidden md:flex" />
	<div class="flex items-center gap-1 md:gap-3">
		<ModeToggle />
		<Sidebar.Trigger size="icon" class="md:hidden" aria-label="Buka navigasi menu" />
		{#if user}
			<div class="hidden md:flex">
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" {...props}>
								<Avatar.Root>
									<Avatar.Image src={user.image ?? undefined} alt={user.name} />
									<Avatar.Fallback />
								</Avatar.Root>
							</Button>
						{/snippet}
					</DropdownMenu.Trigger>
					<DropdownMenu.Content align="end" class="min-w-44">
						<DropdownMenu.Group>
							<DropdownMenu.Item>
								{#snippet child({ props })}
									<a href="/sessions" {...props}>
										<UserKey class="mr-2 size-4" />
										<span>Sesi Aktif</span>
									</a>
								{/snippet}
							</DropdownMenu.Item>
							<DropdownMenu.Separator />
							<DropdownMenu.Item variant="destructive" onclick={handleSignOut}>
								<LogOut class="mr-2 size-4" />
								<span>Keluar</span>
							</DropdownMenu.Item>
						</DropdownMenu.Group>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</div>
		{:else}
			<Button class="hidden md:flex" variant="ghost" onclick={openLogin}>Masuk</Button>
		{/if}
	</div>
</header>
