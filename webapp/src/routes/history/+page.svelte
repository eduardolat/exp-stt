<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Trash2, Clock, Settings, History } from '@lucide/svelte';
	import { HistoryItem, Card, Modal, PageHeader } from '$lib/components';

	let settingsModal: Modal | undefined = $state();

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

<Modal bind:this={settingsModal} title="History Settings" size="sm">
	{#snippet children()}
		<fieldset class="fieldset">
			<label class="label" for="historyLimit">
				<span class="label-text">Maximum entries to keep</span>
			</label>
			<input
				type="number"
				id="historyLimit"
				class="input input-sm w-full"
				min="10"
				max="1000"
				step="10"
				value={store.settings.historyLimit}
				onblur={(e) => store.updateSettings({ historyLimit: parseInt(e.currentTarget.value) })}
			/>
			<p class="label">
				<span class="label-text-alt opacity-70">Older entries will be automatically removed</span>
			</p>
		</fieldset>
	{/snippet}
</Modal>

<PageHeader icon={History} title="History" description="View and manage your recent transcriptions">
	{#snippet actions()}
		<button class="btn btn-sm" onclick={() => settingsModal?.open()} title="History Settings">
			<Settings class="size-3.5" />
			Settings
		</button>
		{#if store.history.length > 0}
			<button class="btn btn-sm btn-error" onclick={handleClearAll}>
				<Trash2 class="size-3.5" />
				Clear All
			</button>
		{/if}
	{/snippet}
</PageHeader>

<div class="space-y-4">
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
