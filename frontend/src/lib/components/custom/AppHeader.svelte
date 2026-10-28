<script lang="ts">
	import { getContext } from 'svelte';
	import { API_ORIGIN } from '$app/env/public';
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import * as Avatar from '#lib/components/ui/avatar/index.js';
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import LogOutIcon from '@lucide/svelte/icons/log-out';
	import type { User } from '#lib/types/user.js';

	const user = getContext<User | null>('user');

	function getInitials(name?: string): string {
		if (!name) return '';
		return name
			.trim()
			.split(/\s+/)
			.map((part) => part[0])
			.join('')
			.toUpperCase()
			.slice(0, 2);
	}

	async function handleLogout() {
		try {
			await fetch(`${API_ORIGIN}/auth/logout`, {
				method: 'DELETE',
				credentials: 'include'
			});
		} finally {
			window.location.href = '/';
		}
	}
</script>

<header
	class="sticky top-0 z-50 flex h-16 w-full items-center justify-between border-b bg-background px-4 md:px-8 lg:px-12"
>
	<a href="/" class="font-serif text-lg font-medium">Optik Pelangi</a>

	<div class="md:hidden">
		<Sidebar.Trigger />
	</div>

	<div class="hidden md:block">
		{#if user}
			<DropdownMenu.Root>
				<DropdownMenu.Trigger>
					{#snippet child({ props })}
						<Button {...props} variant="ghost" size="icon" class="rounded-full">
							<Avatar.Root>
								<Avatar.Image src={user.picture} alt={user.name} />
								<Avatar.Fallback>{getInitials(user.name)}</Avatar.Fallback>
							</Avatar.Root>
						</Button>
					{/snippet}
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end">
					<DropdownMenu.Item variant="destructive" onSelect={handleLogout}>
						<LogOutIcon />
						<span>Keluar</span>
					</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		{:else}
			<Button href="/auth/login">Masuk</Button>
		{/if}
	</div>
</header>
