<script lang="ts">
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
</script>

<header
	class="sticky top-0 z-50 flex h-16 w-full items-center justify-between border-b bg-background px-4 md:px-8 lg:px-12"
>
	<a href="/" class="font-serif text-lg font-medium">Optik Pelangi</a>

	<div class="md:hidden">
		<Sidebar.Trigger />
	</div>

	<div class="hidden md:block">
		<Button href="/auth/login">Masuk</Button>
	</div>
</header>
