<script lang="ts">
	import { onMount } from 'svelte';
	import { goto, refreshAll } from '$app/navigation';
	import { page } from '$app/state';
	import { API_ORIGIN } from '$app/env/public';
	import { Spinner } from '#lib/components/ui/spinner/index.js';
	import type { WebResponse } from '#lib/types/api';
	import type { LoginResponse } from '#lib/types/auth';
	import type { User } from '#lib/types/user';

	onMount(async () => {
		try {
			const response: Response = await fetch(`${API_ORIGIN}/auth/callback${page.url.search}`, {
				method: 'GET',
				credentials: 'include',
				headers: {
					Accept: 'application/json'
				}
			});

			if (response.ok) {
				const result: WebResponse<LoginResponse> = await response.json();
				const userResponse: Response = await fetch(`${API_ORIGIN}/users/current`, {
					method: 'GET',
					headers: {
						Authorization: `Bearer ${result.data.accessToken}`,
						Accept: 'application/json'
					}
				});
				const userResult: WebResponse<User> = await userResponse.json();

				await refreshAll();

				if (userResult.data.role === 'admin') {
					goto('/admin');
				} else {
					goto('/');
				}
				return;
			}

			if (response.status >= 500) {
				const result: WebResponse<unknown> = await response.json();
				goto(`/auth/login?error=${encodeURIComponent(result.errors ?? '')}`);
				return;
			}

			if (response.status >= 400) {
				goto('/auth/login');
			}
		} catch (err) {
			if (err instanceof Error) {
				goto(
					`/auth/login?error=${encodeURIComponent('Tidak dapat terhubung ke internet. Periksa koneksi Anda.')}`
				);
			}
		}
	});
</script>

<svelte:head>
	<title>Memproses Masuk — Optik Pelangi</title>
</svelte:head>

<main class="grid min-h-screen place-items-center">
	<Spinner class="size-6" />
</main>
