<script lang="ts">
	import { getContext } from 'svelte';
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import * as Avatar from '#lib/components/ui/avatar/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import type { User } from '#lib/types/user.js';
	import { getInitials } from '#lib/utils.js';

	const user = getContext<User | null>('user');
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
			<Avatar.Root>
				<Avatar.Image src={user.picture} alt={user.name} />
				<Avatar.Fallback>{getInitials(user.name)}</Avatar.Fallback>
			</Avatar.Root>
		{:else}
			<Button href="/auth/login">Masuk</Button>
		{/if}
	</div>
</header>
