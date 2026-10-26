<script lang="ts">
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import { Button } from '#lib/components/ui/button/index.js';

	const sidebar = Sidebar.useSidebar();
</script>

<Sidebar.Root side="right">
	<Sidebar.Header>
		<a href="/" class="px-2 py-4 font-serif text-xl" onclick={() => sidebar.setOpenMobile(false)}>
			Optik Pelangi
		</a>
	</Sidebar.Header>
	<Sidebar.Content />
	<Sidebar.Footer>
		<Button href="/auth/login" class="w-full" onclick={() => sidebar.setOpenMobile(false)}>
			Masuk
		</Button>
	</Sidebar.Footer>
</Sidebar.Root>
