<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Trash2, Clock, Settings, History, Sparkles } from '@lucide/svelte';
	import { HistoryItem, Card, Modal, PageHeader, CopyButton } from '$lib/components';
	import type { HistoryEntry } from '$lib/client.gen';
	import { formatRelativeTime } from '$lib/utils/formatRelativeTime';

	let settingsModal: Modal | undefined = $state();
	let detailsModal: Modal | undefined = $state();
	let selectedEntry: HistoryEntry | undefined = $state();

	async function handleDelete() {
		if (!selectedEntry) return;
		// Close modal first
		detailsModal?.close();

		// Use confirm dialog as before, or rely on user clicking "Delete" in modal which is explicit enough?
		// The previous implementation used confirm inside handleDelete
		if (!confirm('Delete this transcription?')) return;

		await store.deleteHistoryEntry(selectedEntry.id);
		selectedEntry = undefined;
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
			<Clock class="mb-4 size-10 opacity-50" />
			<h3 class="text-lg font-medium">No transcriptions yet</h3>
			<p class="text-xs opacity-50">Your transcription history will appear here</p>
		</Card>
	{:else}
		<div class="space-y-4">
			{#each store.history as entry (entry.id)}
				<HistoryItem
					{entry}
					onSelect={(e) => {
						selectedEntry = e;
						detailsModal?.open();
					}}
				/>
			{/each}
		</div>
	{/if}
</div>

{#snippet textWithFallback(text: string)}
	{#if text.trim()}
		{text}
	{:else}
		<span class="opacity-70">(No text available)</span>
	{/if}
{/snippet}

<Modal bind:this={detailsModal} title="Transcription Details">
	{#if selectedEntry}
		{@const entry = selectedEntry}
		<div class="space-y-4">
			<div class="flex items-center gap-2 text-sm opacity-70">
				<span>{formatRelativeTime(entry.timestamp)}</span>
				<span>•</span>
				<span>Took {Math.floor(entry.durationMs / 1000)}s</span>
				{#if entry.postProcessed}
					<span class="ml-auto badge gap-1 badge-sm badge-secondary">
						<Sparkles class="size-3" />
						AI Enhanced
					</span>
				{/if}
			</div>

			<Card class="p-4">
				<p class="mb-1 text-xs font-medium opacity-70">Audio:</p>
				<audio controls class="w-full" preload="none" src={`/api/v1/audio/${entry.id}`}>
					Your browser does not support audio playback.
				</audio>
			</Card>

			<Card class="p-4">
				<div class="mb-1 flex items-center justify-between">
					<p class="text-xs font-medium opacity-70">Original:</p>
					<CopyButton text={entry.textRaw} showLabel />
				</div>
				<p class="text-sm">{@render textWithFallback(entry.textRaw)}</p>
			</Card>

			{#if entry.postProcessed && entry.textRaw !== entry.textFinal}
				<Card class="p-4">
					<div class="mb-1 flex items-center justify-between">
						<p class="text-xs font-medium opacity-70">Enhanced:</p>
						<CopyButton text={entry.textFinal} showLabel />
					</div>
					<p class="text-sm">{@render textWithFallback(entry.textFinal)}</p>
				</Card>
			{/if}
		</div>
	{/if}

	{#snippet actions()}
		<div class="flex justify-end gap-2 pt-2">
			<button class="btn btn-outline btn-sm btn-error" onclick={handleDelete}>
				<Trash2 class="size-3.5" />
				Delete
			</button>
		</div>
	{/snippet}
</Modal>
