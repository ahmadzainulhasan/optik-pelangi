<script lang="ts">
	import { onMount, type ComponentProps } from 'svelte';
	import * as Dialog from '#lib/components/ui/dialog/index.js';
	import * as Alert from '#lib/components/ui/alert/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import { Spinner } from '#lib/components/ui/spinner/index.js';
	import { authClient } from '#lib/auth-client.js';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { cn } from '#lib/utils.js';

	let { open = $bindable(false), ...restProps }: ComponentProps<typeof Dialog.Root> = $props();
	let errorMessage = $state<string | null>(null);
	let loading = $state(false);

	onMount(() => {
		if (page.url.searchParams.get('error')) {
			errorMessage = page.url.searchParams.get('error_description');
			open = true;

			const url = new URL(page.url.href);
			url.searchParams.delete('error');
			url.searchParams.delete('error_description');
			goto(url, { shallow: true, replace: true });
		}

		if (page.data.currentUser) {
			open = false;
		}
	});

	async function handleGoogleLogin() {
		errorMessage = null;
		loading = true;

		const { error } = await authClient.signIn.social({
			provider: 'google',
			callbackURL: page.url.href,
			errorCallbackURL: page.url.href
		});

		if (error) {
			errorMessage = error.message || 'Gagal terhubung ke server. Periksa koneksi internet Anda.';
			loading = false;
		}
	}
</script>

<Dialog.Root
	bind:open
	onOpenChange={() => {
		errorMessage = null;
		loading = false;
	}}
	{...restProps}
>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header class="items-center text-center text-balance">
			<Dialog.Title class="text-xl font-semibold">Masuk</Dialog.Title>
			<Dialog.Description>
				Masuk melalui google untuk melanjutkan. Belum punya akun pun tidak masalah, akun akan dibuat
				otomatis.
			</Dialog.Description>
		</Dialog.Header>

		{#if errorMessage}
			<Alert.Root
				variant="destructive"
				class={cn('animate-in fade-in slide-in-from-top-4 duration-400 ease-out')}
			>
				<Alert.Title>Login Gagal</Alert.Title>
				<Alert.Description>
					{errorMessage}
				</Alert.Description>
			</Alert.Root>
		{/if}

		<Button variant="outline" size="lg" disabled={loading} onclick={handleGoogleLogin}>
			{#if loading}
				<Spinner class="size-4" />
				<span class="font-medium">Menghubungkan ke Google...</span>
			{:else}
				<svg class="size-4 shrink-0" viewBox="0 0 24 24" aria-hidden="true">
					<path
						fill="#4285F4"
						d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
					/>
					<path
						fill="#34A853"
						d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
					/>
					<path
						fill="#FBBC05"
						d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.06H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.94l2.85-2.22.81-.63z"
					/>
					<path
						fill="#EA4335"
						d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.06l3.66 2.84c.87-2.6 3.3-4.52 6.16-4.52z"
					/>
				</svg>
				<span class="font-medium">Lanjutkan dengan Google</span>
			{/if}
		</Button>

		<Dialog.Footer>
			<p class="text-center text-muted-foreground">
				Dengan melanjutkan, Anda menyetujui <a
					href="/terms"
					class="underline hover:text-foreground"
					onclick={() => (open = false)}>Ketentuan Layanan</a
				>
				dan
				<a href="/privacy" class="underline hover:text-foreground" onclick={() => (open = false)}
					>Kebijakan Privasi</a
				>
				Optik Pelangi.
			</p>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
