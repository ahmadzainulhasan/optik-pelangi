<script lang="ts">
	import { page } from '$app/state';
	import UserList from './components/user-list.svelte';
	import Pagination from '#lib/components/common/pagination.svelte';
</script>

<UserList />
<Pagination count={page.data.users.total} />
