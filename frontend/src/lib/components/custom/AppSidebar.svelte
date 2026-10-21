<script lang="ts">
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import AuthModal from '#lib/components/custom/AuthModal.svelte';

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
		<AuthModal />
	</Sidebar.Footer>
</Sidebar.Root>
