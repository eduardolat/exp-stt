<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Trash2, Clock, Settings, History, Sparkles } from '@lucide/svelte';
	import { HistoryItem, Card, Modal, PageHeader, CopyButton } from '$lib/components';
	import type { HistoryEntry } from '$lib/client.gen';
	import { formatRelativeTime } from '$lib/date';

	let settingsModal: Modal | undefined = $state();
	let detailsModal: Modal | undefined = $state();
	let selectedEntry: HistoryEntry | undefined = $state();

	async function handleDelete(entry: HistoryEntry) {
		if (!confirm('Delete this transcription?')) return;
		await store.deleteHistoryEntry(entry.id);
		if (selectedEntry?.id === entry.id) {
			detailsModal?.close();
			selectedEntry = undefined;
		}
	}

	async function handleClearAll() {
		if (!confirm('Clear all history? This cannot be undone.')) return;
		await store.clearHistory();
	}

	function openDetails(entry: HistoryEntry) {
		selectedEntry = entry;
		detailsModal?.open();
	}

	function getAudioUrl(id: string): string {
		return `/api/v1/audio/${id}`;
	}

	function formatTimestamp(iso: string): string {
		const text = formatRelativeTime(iso);
		return text.charAt(0).toUpperCase() + text.slice(1);
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

{#snippet textWithFallback(text: string)}
	{#if text.trim()}
		{text}
	{:else}
		<span class="opacity-70">(No text available)</span>
	{/if}
{/snippet}

<Modal bind:this={detailsModal} title="Transcription Details">
	{#snippet children()}
		{#if selectedEntry}
			<div class="space-y-4">
				<div class="flex items-center gap-2 text-sm opacity-70">
					<span>{formatTimestamp(selectedEntry.timestamp)}</span>
					<span>•</span>
					<span>Took {Math.floor(selectedEntry.durationMs / 1000)}s</span>
					{#if selectedEntry.postProcessed}
						<span class="ml-auto badge gap-1 badge-sm badge-secondary">
							<Sparkles class="size-3" />
							AI Enhanced
						</span>
					{/if}
				</div>

				<Card class="p-4">
					<p class="mb-1 text-xs font-medium opacity-70">Audio:</p>
					<audio controls class="w-full" preload="none">
						<source src={getAudioUrl(selectedEntry.id)} type="audio/wav" />
						Your browser does not support audio playback.
					</audio>
				</Card>

				<Card class="p-4">
					<div class="mb-1 flex items-center justify-between">
						<p class="text-xs font-medium opacity-70">Original:</p>
						<CopyButton text={selectedEntry.textRaw} showLabel />
					</div>
					<p class="text-sm">{@render textWithFallback(selectedEntry.textRaw)}</p>
				</Card>

				{#if selectedEntry.postProcessed && selectedEntry.textRaw !== selectedEntry.textFinal}
					<Card class="p-4">
						<div class="mb-1 flex items-center justify-between">
							<p class="text-xs font-medium opacity-70">Enhanced:</p>
							<CopyButton text={selectedEntry.textFinal} showLabel />
						</div>
						<p class="text-sm">{@render textWithFallback(selectedEntry.textFinal)}</p>
					</Card>
				{/if}
			</div>
		{/if}
	{/snippet}

	{#snippet actions()}
		{#if selectedEntry}
			<div class="flex justify-end gap-2 pt-2">
				<button class="btn btn-outline btn-sm btn-error" onclick={() => handleDelete(selectedEntry!)}>
					<Trash2 class="size-3.5" />
					Delete
				</button>
			</div>
		{/if}
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
				<HistoryItem {entry} onView={openDetails} />
			{/each}
		</div>
	{/if}
</div>
