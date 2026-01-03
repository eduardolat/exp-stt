<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Trash2, Clock, Settings } from '@lucide/svelte';
	import { HistoryItem, Card, Modal, HistoryModal } from '$lib/components';
	import type { HistoryEntry } from '$lib/client.gen';

	let settingsModal: Modal | undefined = $state();
	let selectedEntry: HistoryEntry | null = $state(null);

	async function handleClearAll() {
		if (!confirm('Clear all history? This cannot be undone.')) return;
		await store.clearHistory();
	}

	function handleItemClick(entry: HistoryEntry) {
		selectedEntry = entry;
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

<HistoryModal entry={selectedEntry} onClose={() => (selectedEntry = null)} />

<div class="space-y-4">
	<div class="flex items-end justify-between">
		<h2 class="text-lg font-semibold">History</h2>
		<div class="flex gap-2">
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
		</div>
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
				<HistoryItem {entry} onClick={handleItemClick} />
			{/each}
		</div>
	{/if}
</div>
