<script lang="ts">
	import { getContext } from 'svelte';
	import { API_ORIGIN } from '$app/env/public';
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import * as Avatar from '#lib/components/ui/avatar/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import LogOutIcon from '@lucide/svelte/icons/log-out';
	import type { User } from '#lib/types/user.js';

	const sidebar = Sidebar.useSidebar();
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
		sidebar.setOpenMobile(false);
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

<Sidebar.Root side="right">
	<Sidebar.Header>
		<a href="/" class="px-2 py-4 font-serif text-xl" onclick={() => sidebar.setOpenMobile(false)}>
			Optik Pelangi
		</a>
	</Sidebar.Header>
	<Sidebar.Content />
	<Sidebar.Footer class="space-y-2">
		{#if user}
			<div class="flex items-center gap-3 p-2">
				<Avatar.Root>
					<Avatar.Image src={user.picture} alt={user.name} />
					<Avatar.Fallback>
						{getInitials(user.name)}
					</Avatar.Fallback>
				</Avatar.Root>
				<div class="flex flex-col">
					<span class="text-sm font-medium">{user.name}</span>
					<span class="text-xs text-muted-foreground">{user.email}</span>
				</div>
			</div>
			<Button
				variant="outline"
				class="w-full text-destructive hover:bg-destructive/10"
				onclick={handleLogout}
			>
				<LogOutIcon class="size-4" />
				<span>Keluar</span>
			</Button>
		{:else}
			<Button href="/auth/login" class="w-full" onclick={() => sidebar.setOpenMobile(false)}>
				Masuk
			</Button>
		{/if}
	</Sidebar.Footer>
</Sidebar.Root>
