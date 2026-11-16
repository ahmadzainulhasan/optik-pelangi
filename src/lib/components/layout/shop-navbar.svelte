<script lang="ts">
	import { getContext, type ComponentProps } from 'svelte';
	import * as NavigationMenu from '#lib/components/ui/navigation-menu/index.js';
	import { cn } from '#lib/utils.js';

	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof NavigationMenu.Root> =
		$props();

	const navItems = getContext<Array<{ title: string; href: string }>>('navItems') ?? [];
</script>

<NavigationMenu.Root
	bind:ref
	class={cn('animate-in duration-500 ease-out fade-in', restProps.class)}
	{...restProps}
>
	<NavigationMenu.List>
		{#each navItems as item (item.title)}
			<NavigationMenu.Item>
				<NavigationMenu.Link href={item.href}>
					{item.title}
				</NavigationMenu.Link>
			</NavigationMenu.Item>
		{/each}
	</NavigationMenu.List>
</NavigationMenu.Root>
