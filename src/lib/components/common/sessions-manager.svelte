<script lang="ts">
	import { page } from '$app/state';
	import { refreshAll } from '$app/navigation';
	import { authClient } from '#lib/auth-client.js';
	import { toast } from '#lib/components/ui/sonner/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import Laptop from '@lucide/svelte/icons/laptop';
	import Smartphone from '@lucide/svelte/icons/smartphone';
	import { cn } from '#lib/utils.js';
	import Bowser from 'bowser';

	const sessions = $derived(page.data.sessions);
	const currentSession = $derived(page.data.currentSession);

	function parseUserAgent(uaString?: string | null) {
		if (!uaString) return { browser: 'Tidak dikenal', os: 'Tidak dikenal', isMobile: false };

		const parser = Bowser.getParser(uaString);
		return {
			browser: parser.getBrowserName() || 'Tidak dikenal',
			os: parser.getOSName() || 'Tidak dikenal',
			isMobile: ['mobile', 'tablet'].includes(parser.getPlatformType())
		};
	}

	async function handleRevokeSession(token: string) {
		const { error } = await authClient.revokeSession({ token });
		if (error) {
			toast.error(error.message || 'Gagal terhubung ke server. Periksa koneksi internet Anda.');
			return;
		}
		await refreshAll();
	}

	async function handleRevokeOtherSessions() {
		const { error } = await authClient.revokeOtherSessions();
		if (error) {
			toast.error(error.message || 'Gagal terhubung ke server. Periksa koneksi internet Anda.');
			return;
		}
		await refreshAll();
	}
</script>

<div class={cn('flex flex-col gap-4', 'animate-in fade-in-50 slide-in-from-bottom-6 duration-300 ease-out')}>
	{#if sessions.length > 1}
		<div class="flex justify-end">
			<Button
				variant="outline"
				size="sm"
				class="h-8 shrink-0 text-xs font-medium text-destructive hover:border-destructive/30 hover:bg-destructive/10 hover:text-destructive"
				onclick={handleRevokeOtherSessions}
			>
				Cabut semua sesi lainnya
			</Button>
		</div>
	{/if}

	<div class="overflow-hidden rounded-lg border text-card-foreground shadow-xs">
		<ul class="divide-y divide-card">
			{#each sessions as session (session.id || session.token)}
				{@const ua = parseUserAgent(session.userAgent)}
				<li
					class={cn(
						'flex flex-col gap-4 p-4 transition-colors sm:flex-row sm:items-center sm:justify-between',
						session.token === currentSession?.token ? 'bg-card' : 'hover:bg-muted/10'
					)}
				>
					<div class="flex min-w-0 items-start gap-3.5 sm:gap-4">
						<div class="mt-0.5 shrink-0 rounded-md bg-muted/60 p-2 text-muted-foreground">
							{#if ua.isMobile}
								<Smartphone class="size-5" />
							{:else}
								<Laptop class="size-5" />
							{/if}
						</div>

						<div class="min-w-0 flex-1 flex flex-col gap-0.5">
							<div class="text-sm font-medium text-foreground">
								{ua.browser}
							</div>
							<div class="flex flex-wrap items-center gap-x-3 text-xs text-muted-foreground">
								<span>{ua.os}</span>
								<span>
									{session.token === currentSession?.token
										? 'Aktif sekarang'
										: `Terakhir aktif ${session.updatedAt.toLocaleString('id-ID')}`}
								</span>
							</div>
						</div>
					</div>

					<div class="shrink-0 self-end sm:self-center">
						{#if session.token !== currentSession?.token}
							<Button
								variant="outline"
								size="sm"
								class="h-8 px-3 text-xs text-destructive hover:border-destructive/30 hover:bg-destructive/10 hover:text-destructive"
								onclick={() => handleRevokeSession(session.token)}
							>
								Cabut
							</Button>
						{/if}
					</div>
				</li>
			{/each}
		</ul>
	</div>
</div>
