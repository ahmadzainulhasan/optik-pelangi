<script lang="ts">
	import { page } from '$app/state';
	import * as Empty from '#lib/components/ui/empty/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import { refreshAll } from '$app/navigation';
	import Users from '@lucide/svelte/icons/users';
	import { cn } from '#lib/utils.js';
</script>

<Empty.Root class={cn('border border-dashed', 'animate-in fade-in zoom-in-90 duration-500 ease-out')}>
	<Empty.Header>
		<Empty.Media variant="icon">
			<Users />
		</Empty.Media>
		<Empty.Title>Gagal memuat data pengguna</Empty.Title>
		<Empty.Description>
			{page.status} - {page.error?.message ||
				'Terjadi kesalahan saat memuat daftar pengguna.'}
		</Empty.Description>
	</Empty.Header>
	<Empty.Content class="flex-row justify-center gap-2">
		<Button size="sm" onclick={() => refreshAll()}>Coba Lagi</Button>
		<Button variant="outline" size="sm" href="/admin">Dashboard</Button>
	</Empty.Content>
</Empty.Root>
