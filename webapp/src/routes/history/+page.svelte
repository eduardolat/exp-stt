<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Trash2, Clock, Settings, History } from '@lucide/svelte';
	import { HistoryItem, Card, Modal, PageHeader } from '$lib/components';

	let settingsModal: Modal | undefined = $state();
	let confirmModal: Modal | undefined = $state();
	let confirmationState = $state<{ type: 'delete' | 'clear'; entry?: (typeof store.history)[0] }>({
		type: 'delete'
	});

	function handleDelete(entry: (typeof store.history)[0]) {
		confirmationState = { type: 'delete', entry };
		confirmModal?.open();
	}

	function handleClearAll() {
		confirmationState = { type: 'clear' };
		confirmModal?.open();
	}

	async function onConfirm() {
		if (confirmationState.type === 'delete' && confirmationState.entry) {
			await store.deleteHistoryEntry(confirmationState.entry.id);
		} else if (confirmationState.type === 'clear') {
			await store.clearHistory();
		}
		confirmModal?.close();
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
				<HistoryItem {entry} onDelete={handleDelete} />
			{/each}
		</div>
	{/if}
</div>

<Modal
	bind:this={confirmModal}
	title={confirmationState.type === 'delete' ? 'Delete Transcription' : 'Clear All History'}
>
	<p class="py-4">
		{confirmationState.type === 'delete'
			? 'Are you sure you want to delete this transcription?'
			: 'Are you sure you want to clear all history? This cannot be undone.'}
	</p>
	{#snippet actions()}
		<button class="btn" onclick={() => confirmModal?.close()}>Cancel</button>
		<button class="btn btn-error" onclick={onConfirm}>
			{confirmationState.type === 'delete' ? 'Delete' : 'Clear All'}
		</button>
	{/snippet}
</Modal>
