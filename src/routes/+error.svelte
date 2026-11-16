<script lang="ts">
	import { page } from '$app/state';
	import * as Empty from '#lib/components/ui/empty/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import AlertTriangle from '@lucide/svelte/icons/alert-triangle';
	import { cn } from '#lib/utils.js';
</script>

<div class="flex min-h-[calc(100vh-4rem)] items-center justify-center p-4">
	<Empty.Root class={cn('w-full max-w-md border border-dashed', 'animate-in fade-in zoom-in-90 duration-500 ease-out')}>
		<Empty.Header>
			<Empty.Media variant="icon">
				<AlertTriangle />
			</Empty.Media>
			<Empty.Title>
				{page.status} - {page.error?.message}
			</Empty.Title>
			<Empty.Description>
				Halaman yang Anda tuju tidak ditemukan silahkan tekan tombol dibawah untuk kembali ke
				beranda.
			</Empty.Description>
		</Empty.Header>
		<Empty.Content class="flex-row justify-center gap-2">
			<Button variant="outline" size="sm" href="/">Beranda</Button>
		</Empty.Content>
	</Empty.Root>
</div>
