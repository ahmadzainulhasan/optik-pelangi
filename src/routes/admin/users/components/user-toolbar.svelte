<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { Input } from '#lib/components/ui/input/index.js';
	import { Button, buttonVariants } from '#lib/components/ui/button/index.js';
	import * as Sheet from '#lib/components/ui/sheet/index.js';
	import Search from '@lucide/svelte/icons/search';
	import SlidersHorizontal from '@lucide/svelte/icons/sliders-horizontal';
	import UserKey from '@lucide/svelte/icons/user-key';
	import AddAdminDialog from './add-admin-dialog.svelte';
	import { SvelteURLSearchParams } from 'svelte/reactivity';
	import { cn } from '#lib/utils.js';

	let isFilterOpen = $state(false);

	const roleOptions = [
		{ value: 'admin', label: 'Admin' },
		{ value: 'user', label: 'Pelanggan' }
	];

	const sortButtons = [
		{ value: 'createdAt-desc', label: 'Terbaru' },
		{ value: 'createdAt-asc', label: 'Terlama' },
		{ value: 'name-asc', label: 'Nama (A - Z)' },
		{ value: 'name-desc', label: 'Nama (Z - A)' }
	];

	function handleSortChange(value: string) {
		const params = new SvelteURLSearchParams(Object.fromEntries(page.url.searchParams));
		params.set('sort', value);
		params.delete('page');
		goto(`/admin/users?${params}`, { replaceState: true });
	}

	function handleChange(e: Event & { currentTarget: HTMLInputElement }) {
		const params = new SvelteURLSearchParams(Object.fromEntries(page.url.searchParams));
		if (e.currentTarget.value) {
			params.set('search', e.currentTarget.value);
		} else {
			params.delete('search');
		}
		params.delete('page');
		goto(`/admin/users?${params}`, { replaceState: true });
	}

	function handleRoleChange(role: string) {
		const params = new SvelteURLSearchParams(Object.fromEntries(page.url.searchParams));
		if (role === 'admin' || role === 'user') {
			params.set('role', role);
		} else {
			params.delete('role');
		}
		params.delete('page');
		goto(`/admin/users?${params}`, { replaceState: true });
		isFilterOpen = false;
	}
</script>

<div class={cn('flex flex-col gap-3', 'animate-in fade-in slide-in-from-top-3 duration-500 ease-out')}>
	<div class="flex items-center gap-2">
		<div class="relative flex-1">
			<Search class="absolute inset-s-2.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
			<Input
				type="search"
				placeholder="Cari nama atau email..."
				value={page.url.searchParams.get('search') ?? ''}
				onchange={handleChange}
				class="ps-8"
			/>
		</div>

		<Sheet.Root bind:open={isFilterOpen}>
			<Sheet.Trigger
				class={buttonVariants({
					variant: 'outline'
				})}
			>
				<SlidersHorizontal class="size-4" />
			</Sheet.Trigger>

			<Sheet.Content>
				<Sheet.Header>
					<Sheet.Title>Filter Pengguna</Sheet.Title>
					<Sheet.Description>Saring daftar pengguna berdasarkan kriteria</Sheet.Description>
				</Sheet.Header>
				<div class="flex flex-col gap-4 px-6">
					<div class="flex items-center gap-2 font-medium">
						<UserKey class="size-4" />
						<span>Akses Pengguna</span>
					</div>
					<div class="flex flex-wrap justify-center gap-2">
						{#each roleOptions as item (item.value)}
							<Button
								variant={(page.url.searchParams.get('role') ?? 'all') === item.value
									? 'default'
									: 'outline'}
								size="sm"
								onclick={() => handleRoleChange(item.value)}
							>
								{item.label}
							</Button>
						{/each}
					</div>
				</div>

				<Sheet.Footer>
					<Button class="w-full" onclick={() => handleRoleChange('')}>Reset Filter</Button>
				</Sheet.Footer>
			</Sheet.Content>
		</Sheet.Root>
	</div>

	<div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
		<div class="no-scrollbar flex items-center gap-2 overflow-x-auto pb-1">
			{#each sortButtons as item (item.value)}
				<Button
					variant={(page.url.searchParams.get('sort') ?? 'createdAt-desc') === item.value
						? 'default'
						: 'outline'}
					class="shrink-0"
					onclick={() => handleSortChange(item.value)}
				>
					{item.label}
				</Button>
			{/each}
		</div>

		{#if page.url.searchParams.get('role') === 'admin'}
			<AddAdminDialog />
		{/if}
	</div>
</div>
