<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	onMount(() => {
		goto('#/settings/general', { replaceState: true });
	});
</script>

<svelte:head>
	<title>Settings - Tribar Voice</title>
</svelte:head>
