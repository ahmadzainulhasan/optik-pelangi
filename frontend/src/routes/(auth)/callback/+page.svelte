<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	onMount(() => {
		if (window.opener) {
			window.close();
		} else {
			goto('/');
		}
	});
</script>

<svelte:head>
	<title>Berhasil Masuk — Optik Pelangi</title>
</svelte:head>
