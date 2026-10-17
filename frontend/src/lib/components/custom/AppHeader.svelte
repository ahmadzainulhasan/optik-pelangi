<script lang="ts">
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import AuthModal from '#lib/components/custom/AuthModal.svelte';
</script>

<header
	class="sticky top-0 z-50 flex h-16 w-full items-center justify-between border-b bg-background px-4 md:px-8 lg:px-12"
>
	<a href="/" class="shrink-0 font-serif text-lg font-medium">Optik Pelangi</a>

	<div class="flex items-center gap-1 md:hidden">
		<Sidebar.Trigger />
	</div>

	<div class="hidden shrink-0 items-center gap-2 md:flex">
		<AuthModal />
	</div>
</header>
