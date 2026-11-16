<script lang="ts">
	import { page } from '$app/state';
	import * as Avatar from '#lib/components/ui/avatar/index.js';
	import { Badge } from '#lib/components/ui/badge/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import Mail from '@lucide/svelte/icons/mail';
	import Phone from '@lucide/svelte/icons/phone';
	import Calendar from '@lucide/svelte/icons/calendar';
	import MessageSquare from '@lucide/svelte/icons/message-square';
	import { cn } from '#lib/utils.js';

	const user = $derived(page.data.user);
</script>

<div
	class={cn(
		'flex flex-col justify-between gap-3.5 py-1 sm:flex-row sm:items-center',
		'animate-in fade-in slide-in-from-top-3 duration-500 ease-out'
	)}
>
	<div class="flex min-w-0 items-center gap-3.5">
		<Avatar.Root class="size-13 shrink-0 border sm:size-15">
			<Avatar.Image src={user.image} alt={user.name} />
			<Avatar.Fallback />
		</Avatar.Root>
		<div class="flex min-w-0 flex-col gap-1">
			<div class="flex flex-wrap items-center gap-2">
				<h1 class="truncate text-lg font-bold tracking-tight sm:text-2xl">
					{user.name}
				</h1>
				{#if user.banned}
					<Badge variant="destructive">Diblokir</Badge>
				{/if}
			</div>
			<div class="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-muted-foreground">
				<span class="inline-flex items-center gap-1.5">
					<Mail class="size-3.5 shrink-0" />
					<span>{user.email}</span>
				</span>
				{#if user.contact?.whatsapp}
					<a
						href={`https://wa.me/${user.contact.whatsapp}`}
						target="_blank"
						rel="noopener noreferrer"
						class="inline-flex items-center gap-1.5 hover:text-foreground hover:underline"
					>
						<Phone class="size-3.5 shrink-0" />
						<span>{user.contact.whatsapp}</span>
					</a>
				{/if}
				<span class="inline-flex items-center gap-1.5">
					<Calendar class="size-3.5 shrink-0" />
					<span>
						Bergabung {user.createdAt
							? new Date(user.createdAt).toLocaleDateString('id-ID', {
									day: 'numeric',
									month: 'short',
									year: 'numeric'
								})
							: ''}
					</span>
				</span>
			</div>
		</div>
	</div>

	{#if user.contact?.whatsapp}
		<Button
			variant="outline"
			size="sm"
			href={`https://wa.me/${user.contact.whatsapp}`}
			target="_blank"
			rel="noopener noreferrer"
			class="w-full shrink-0 gap-1.5 sm:w-auto"
		>
			<MessageSquare class="size-3.5" />
			<span>Hubungi WhatsApp</span>
		</Button>
	{/if}
</div>
