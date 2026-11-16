<script lang="ts">
	import { page } from '$app/state';
	import { refreshAll } from '$app/navigation';
	import { authClient } from '#lib/auth-client.js';
	import { toast } from '#lib/components/ui/sonner/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import * as Card from '#lib/components/ui/card/index.js';
	import { cn } from '#lib/utils.js';

	function banUser(userId: string) {
		toast.promise(
			async () => {
				const { error } = await authClient.admin.banUser({ userId });
				if (error) throw new Error(error.message || 'Gagal memblokir akun.');
				await refreshAll();
			},
			{
				loading: 'Memblokir akun...',
				success: 'Akun berhasil diblokir.',
				error: (err) => (err instanceof Error ? err.message : 'Gagal memblokir akun.')
			}
		);
	}

	function unbanUser(userId: string) {
		toast.promise(
			async () => {
				const { error } = await authClient.admin.unbanUser({ userId });
				if (error) throw new Error(error.message || 'Gagal membuka blokir akun.');
				await refreshAll();
			},
			{
				loading: 'Membuka blokir akun...',
				success: 'Blokir akun berhasil dibuka.',
				error: (err) => (err instanceof Error ? err.message : 'Gagal membuka blokir akun.')
			}
		);
	}

	function toggleRole(userId: string, currentRole?: string | null) {
		toast.promise(
			async () => {
				const { error } = await authClient.admin.setRole({
					userId,
					role: currentRole === 'admin' ? 'user' : 'admin'
				});
				if (error) throw new Error(error.message || 'Gagal mengubah peran pengguna.');
				await refreshAll();
			},
			{
				loading: 'Mengubah peran pengguna...',
				success: `Peran berhasil diubah menjadi ${currentRole === 'admin' ? 'Pelanggan' : 'Admin'}.`,
				error: (err) => (err instanceof Error ? err.message : 'Gagal mengubah peran pengguna.')
			}
		);
	}
	const user = $derived(page.data.user);
</script>

<Card.Root class={cn('animate-in fade-in-50 slide-in-from-bottom-6 duration-300 ease-out')}>
		<Card.Header>
			<Card.Title>Akses Akun</Card.Title>
			<Card.Description>Kelola peran dan status pemblokiran pengguna.</Card.Description>
		</Card.Header>
		<Card.Content class="flex flex-col gap-4">
			<div class="flex flex-col justify-between gap-3 sm:flex-row sm:items-center">
				<div class="flex flex-col gap-0.5">
					<p class="text-sm font-medium">Peran Pengguna</p>
					<p class="text-xs text-muted-foreground">
						{user.role === 'admin'
							? 'Pengguna memiliki hak akses admin.'
							: 'Pengguna memiliki hak akses pelanggan.'}
					</p>
				</div>
				<Button
					variant="secondary"
					disabled={page.data.currentUser?.id === user.id && user.role === 'admin'}
					onclick={() => toggleRole(user.id, user.role)}
				>
					{user.role === 'admin' ? 'Ubah ke Pelanggan' : 'Jadikan Admin'}
				</Button>
			</div>

			<div class="flex flex-col justify-between gap-3 sm:flex-row sm:items-center">
				<div class="flex flex-col gap-0.5">
					<p class="text-sm font-medium">Status Akun</p>
					<p class="text-xs text-muted-foreground">
						{#if user.banned}
							Pengguna ini diblokir dan tidak dapat digunakan untuk masuk.
						{:else}
							Pengguna yang diblokir tidak akan dapat masuk kembali ke akun.
						{/if}
					</p>
				</div>
				{#if user.banned}
					<Button variant="outline" onclick={() => unbanUser(user.id)}>Buka Blokir Akun</Button>
				{:else}
					<Button
						variant="destructive"
						disabled={page.data.currentUser?.id === user.id}
						onclick={() => banUser(user.id)}
					>
						Blokir Akun
					</Button>
				{/if}
			</div>
		</Card.Content>
	</Card.Root>
