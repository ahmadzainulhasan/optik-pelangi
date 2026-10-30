<script lang="ts">
	import { getContext } from 'svelte';
	import { API_ORIGIN } from '$app/env/public';
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import * as Avatar from '#lib/components/ui/avatar/index.js';
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import ChevronsUpDownIcon from '@lucide/svelte/icons/chevrons-up-down';
	import LogOutIcon from '@lucide/svelte/icons/log-out';
	import type { User } from '#lib/types/user.js';

	const sidebar = Sidebar.useSidebar();
	const user = getContext<User | null>('user');

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

<Sidebar.Root side="left" collapsible="offcanvas">
	<Sidebar.Header class="h-14 justify-center px-4">
		<a href="/" class="font-serif text-xl" onclick={() => sidebar.setOpenMobile(false)}>
			Optik Pelangi
		</a>
	</Sidebar.Header>

	<Sidebar.Content />

	<Sidebar.Footer>
		{#if user}
			<Sidebar.Menu>
				<Sidebar.MenuItem>
					<DropdownMenu.Root>
						<DropdownMenu.Trigger>
							{#snippet child({ props })}
								<Sidebar.MenuButton {...props}>
									<Avatar.Root>
										<Avatar.Image src={user.picture} alt={user.name} />
										<Avatar.Fallback />
									</Avatar.Root>
									<div class="flex w-full flex-col">
										<span class="font-medium">{user.name}</span>
										<span class="text-[10px] text-muted-foreground">{user.email}</span>
									</div>
									<ChevronsUpDownIcon />
								</Sidebar.MenuButton>
							{/snippet}
						</DropdownMenu.Trigger>
						<DropdownMenu.Content>
							<DropdownMenu.Item variant="destructive" onSelect={handleLogout}>
								<LogOutIcon />
								<span>Keluar Akun</span>
							</DropdownMenu.Item>
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				</Sidebar.MenuItem>
			</Sidebar.Menu>
		{:else}
			<Button href="/auth/login" class="w-full" onclick={() => sidebar.setOpenMobile(false)}>
				Masuk
			</Button>
		{/if}
	</Sidebar.Footer>
</Sidebar.Root>
