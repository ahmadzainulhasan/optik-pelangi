<script lang="ts">
	import { page } from '$app/state';
	import * as Table from '#lib/components/ui/table/index.js';
	import * as Avatar from '#lib/components/ui/avatar/index.js';
	import { Badge } from '#lib/components/ui/badge/index.js';
	import ChevronRight from '@lucide/svelte/icons/chevron-right';
	import { goto } from '$app/navigation';
	import { cn } from '#lib/utils.js';

	const users = $derived(page.data.users.users);
</script>

<div
	class={cn(
		'flex flex-col divide-y md:hidden',
		'rounded-lg border',
		'animate-in fade-in-50 slide-in-from-bottom-6 duration-300 ease-out'
	)}
>
	{#each users as user (user.id)}
		<a
			class="flex items-center justify-between p-3.5 transition-colors hover:bg-muted/50"
			href="/admin/users/{user.id}"
		>
			<div class="flex items-center gap-3">
				<Avatar.Root class="size-12">
					<Avatar.Image src={user.image} alt={user.name} />
					<Avatar.Fallback />
				</Avatar.Root>
				<div>
					<p class="text-sm font-medium">{user.name}</p>
					<p class="text-xs text-muted-foreground">{user.email}</p>
					<div class="flex items-center gap-2 pt-1">
						<Badge variant="secondary" class="text-[10px]">
							{user.role === 'admin' ? 'Admin' : 'Pelanggan'}
						</Badge>
						<span
							class={cn('text-[10px]', user.banned ? 'text-destructive font-medium' : 'text-muted-foreground')}
						>
							{user.banned ? 'Diblokir' : 'Aktif'}
						</span>
						<span class="text-xs text-muted-foreground">
							{user.createdAt.toLocaleDateString('id-ID', {
								day: 'numeric',
								month: 'short',
								year: 'numeric'
							})}
						</span>
					</div>
				</div>
			</div>
			<ChevronRight class="size-4 text-muted-foreground" />
		</a>
	{/each}
</div>

<div
	class={cn(
		'hidden md:block',
		'rounded-lg border',
		'animate-in fade-in-50 slide-in-from-bottom-6 duration-300 ease-out'
	)}
>
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Pengguna</Table.Head>
				<Table.Head>Role</Table.Head>
				<Table.Head>Status</Table.Head>
				<Table.Head>Bergabung</Table.Head>
				<Table.Head class="w-10"></Table.Head>
			</Table.Row>
		</Table.Header>

		<Table.Body>
			{#each users as user (user.id)}
				<Table.Row class="cursor-pointer" onclick={() => goto(`/admin/users/${user.id}`)}>
					<Table.Cell>
						<div class="flex items-center gap-3">
							<Avatar.Root>
								<Avatar.Image src={user.image} alt={user.name} />
								<Avatar.Fallback />
							</Avatar.Root>
							<div class="min-w-0 leading-tight">
								<p class="truncate font-medium">{user.name}</p>
								<p class="truncate text-xs text-muted-foreground">{user.email}</p>
							</div>
						</div>
					</Table.Cell>

					<Table.Cell>
						<Badge variant="secondary">
							{user.role === 'admin' ? 'Admin' : 'Pelanggan'}
						</Badge>
					</Table.Cell>

					<Table.Cell class={user.banned && 'font-medium text-destructive'}>
						{user.banned ? 'Diblokir' : 'Aktif'}
					</Table.Cell>

					<Table.Cell>
						{user.createdAt.toLocaleDateString('id-ID', {
							day: 'numeric',
							month: 'short',
							year: 'numeric'
						})}
					</Table.Cell>

					<Table.Cell class="text-end">
						<ChevronRight class="ms-auto size-4" />
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
