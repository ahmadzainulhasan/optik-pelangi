<script lang="ts">
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import * as Avatar from '#lib/components/ui/avatar/index.js';
	import AuthModal from '#lib/components/custom/AuthModal.svelte';
	import type { User } from '#lib/types/user.js';
	import { getInitials } from '#lib/helpers/avatar.js';

	interface Props {
		user?: User | null;
	}

	let { user = null }: Props = $props();

	const sidebar = Sidebar.useSidebar();
</script>

<Sidebar.Root side="right">
	<Sidebar.Header>
		<a href="/" class="px-2 py-4 font-serif text-xl" onclick={() => sidebar.setOpenMobile(false)}>
			Optik Pelangi
		</a>
	</Sidebar.Header>
	<Sidebar.Content />
	<Sidebar.Footer>
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
		{:else}
			<AuthModal />
		{/if}
	</Sidebar.Footer>
</Sidebar.Root>
