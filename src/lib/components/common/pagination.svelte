<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import * as Pagination from '#lib/components/ui/pagination/index.js';
	import { SvelteURLSearchParams } from 'svelte/reactivity';
	import type { ComponentProps } from 'svelte';
	import { cn } from '#lib/utils.js';

	let {
		count = 0,
		perPage = 10,
		class: className,
		...restProps
	}: ComponentProps<typeof Pagination.Root> = $props();

	function handlePageChange(newPage: number) {
		const params = new SvelteURLSearchParams(Object.fromEntries(page.url.searchParams));
		params.set('page', String(newPage));
		goto(`${page.url.pathname}?${params}`, { replaceState: true });
	}
</script>

{#if count > 0}
	<Pagination.Root
		{count}
		{perPage}
		page={Number(page.url.searchParams.get('page') || 1)}
		onPageChange={handlePageChange}
		class={cn('animate-in fade-in slide-in-from-bottom-2 duration-500 ease-out', className)}
		{...restProps}
	>
		{#snippet children({ pages, currentPage })}
			<Pagination.Content>
				<Pagination.Item>
					<Pagination.Previous />
				</Pagination.Item>
				{#each pages as p (p.key)}
					{#if p.type === 'ellipsis'}
						<Pagination.Item class="hidden sm:block">
							<Pagination.Ellipsis />
						</Pagination.Item>
					{:else}
						<Pagination.Item>
							<Pagination.Link page={p} isActive={currentPage === p.value}>
								{p.value}
							</Pagination.Link>
						</Pagination.Item>
					{/if}
				{/each}
				<Pagination.Item>
					<Pagination.Next />
				</Pagination.Item>
			</Pagination.Content>
		{/snippet}
	</Pagination.Root>
{/if}
