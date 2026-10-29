<script lang="ts">
	import * as NavigationMenu from '#lib/components/ui/navigation-menu/index.js';
	import { mainNav } from '#lib/config/navigation.js';
	import { cn } from '#lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';

	type ListItemProps = HTMLAttributes<HTMLAnchorElement> & {
		title: string;
		href: string;
		content: string;
	};
</script>

{#snippet ListItem({ title, content, href, class: className, ...restProps }: ListItemProps)}
	<li>
		<NavigationMenu.Link>
			{#snippet child()}
				<a
					{href}
					class={cn(
						'block space-y-1 rounded-md p-3 leading-none no-underline transition-colors outline-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground',
						className
					)}
					{...restProps}
				>
					<div class="text-sm leading-none font-medium">{title}</div>
					<p class="line-clamp-2 text-sm leading-snug text-muted-foreground">
						{content}
					</p>
				</a>
			{/snippet}
		</NavigationMenu.Link>
	</li>
{/snippet}

<NavigationMenu.Root viewport={false}>
	<NavigationMenu.List class="flex-wrap gap-1">
		{#each mainNav as navItem (navItem.title)}
			<NavigationMenu.Item>
				<NavigationMenu.Trigger>{navItem.title}</NavigationMenu.Trigger>
				<NavigationMenu.Content>
					<ul class="grid w-75 gap-2 p-4 md:w-125 md:grid-cols-2 lg:w-150">
						{#each navItem.items as subItem}
							{@render ListItem({
								href: subItem.href,
								title: subItem.title,
								content: subItem.description
							})}
						{/each}

						{#if navItem.title === 'Produk'}
							<li class="col-span-1 md:col-span-2">
								<NavigationMenu.Link>
									{#snippet child()}
										<a
											href={navItem.href}
											class="mt-2 flex w-full items-center justify-center rounded-md bg-muted p-3 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground"
										>
											Lihat Semua {navItem.title} &rarr;
										</a>
									{/snippet}
								</NavigationMenu.Link>
							</li>
						{/if}
					</ul>
				</NavigationMenu.Content>
			</NavigationMenu.Item>
		{/each}
	</NavigationMenu.List>
</NavigationMenu.Root>
