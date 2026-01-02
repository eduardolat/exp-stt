<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Trash2, Clock } from '@lucide/svelte';
	import { HistoryItem, Card } from '$lib/components';

	async function handleDelete(entry: (typeof store.history)[0]) {
		if (!confirm('Delete this transcription?')) return;
		await store.deleteHistoryEntry(entry.id);
	}

	async function handleClearAll() {
		if (!confirm('Clear all history? This cannot be undone.')) return;
		await store.clearHistory();
	}
</script>

<svelte:head>
	<title>History - Tribar Voice</title>
</svelte:head>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h2 class="text-lg font-semibold">History</h2>
		{#if store.history.length > 0}
			<button class="btn btn-outline btn-sm btn-error" onclick={handleClearAll}>
				<Trash2 class="size-3.5" />
				Clear All
			</button>
		{/if}
	</div>

	{#if store.history.length === 0}
		<Card class="card-body items-center py-12 text-center">
			<Clock class="mb-3 size-10 opacity-50" />
			<p class="opacity-70">No transcriptions yet</p>
			<p class="text-xs opacity-50">Your transcription history will appear here</p>
		</Card>
	{:else}
		<div class="space-y-4">
			{#each store.history as entry (entry.id)}
				<HistoryItem {entry} onDelete={handleDelete} />
			{/each}
		</div>
	{/if}
</div>
