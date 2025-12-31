<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Copy, Trash2, Clock, Sparkles, ChevronDown, ChevronUp } from '@lucide/svelte';
	import type { HistoryEntry } from '$lib/client.gen';

	function getAudioUrl(id: string): string {
		return `/api/v1/audio/${id}`;
	}

	let expandedId = $state<string | null>(null);

	function formatDuration(ms: number): string {
		const seconds = Math.floor(ms / 1000);
		const minutes = Math.floor(seconds / 60);
		const remainingSeconds = seconds % 60;

		if (minutes > 0) {
			return `${minutes}m ${remainingSeconds}s`;
		}
		return `${remainingSeconds}s`;
	}

	function formatTimestamp(iso: string): string {
		const date = new Date(iso);
		return date.toLocaleString(undefined, {
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	async function copyToClipboard(text: string) {
		await navigator.clipboard.writeText(text);
	}

	function toggleExpand(id: string) {
		expandedId = expandedId === id ? null : id;
	}

	async function handleDelete(entry: HistoryEntry) {
		if (!confirm('Delete this transcription?')) return;
		await store.deleteHistoryEntry(entry.id);
	}

	async function handleClearAll() {
		if (!confirm('Clear all history? This cannot be undone.')) return;
		await store.clearHistory();
	}
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h2 class="text-lg font-semibold">Transcription History</h2>
		{#if store.history.length > 0}
			<button class="btn btn-outline btn-sm btn-error" onclick={handleClearAll}>
				<Trash2 class="size-3.5" />
				Clear All
			</button>
		{/if}
	</div>

	{#if store.history.length === 0}
		<div class="card bg-base-200">
			<div class="card-body items-center py-12 text-center">
				<Clock class="mb-3 size-10 opacity-50" />
				<p class="opacity-70">No transcriptions yet</p>
				<p class="text-xs opacity-50">Your transcription history will appear here</p>
			</div>
		</div>
	{/if}

	{#if store.history.length > 0}
		<div class="space-y-2">
			{#each store.history as entry (entry.id)}
				<div class="card overflow-hidden bg-base-200">
					<button
						class="flex w-full items-start gap-3 p-4 text-left hover:bg-base-300"
						onclick={() => toggleExpand(entry.id)}
					>
						<div class="min-w-0 flex-1">
							<div class="mb-1 flex items-center gap-2 text-xs opacity-70">
								<span>{formatTimestamp(entry.timestamp)}</span>
								<span>•</span>
								<span>{formatDuration(entry.durationMs)}</span>
								{#if entry.postProcessed}
									<span class="badge gap-1 badge-sm badge-secondary">
										<Sparkles class="size-3" />
										AI Enhanced
									</span>
								{/if}
							</div>
							<p class="line-clamp-2 text-sm">{entry.textFinal}</p>
						</div>
						{#if expandedId === entry.id}
							<ChevronUp class="size-4 shrink-0 opacity-50" />
						{:else}
							<ChevronDown class="size-4 shrink-0 opacity-50" />
						{/if}
					</button>

					{#if expandedId === entry.id}
						<div class="border-t border-base-300 bg-base-300/50 p-4">
							<div class="mb-3">
								<p class="mb-1 text-xs font-medium opacity-70">Audio:</p>
								<audio controls class="w-full" preload="none">
									<source src={getAudioUrl(entry.id)} type="audio/wav" />
									Your browser does not support audio playback.
								</audio>
							</div>

							<div class="mb-3">
								<div class="mb-1 flex items-center justify-between">
									<p class="text-xs font-medium opacity-70">Original:</p>
									<button
										class="btn btn-ghost btn-xs"
										onclick={() => copyToClipboard(entry.textRaw)}
									>
										<Copy class="size-3" />
										Copy
									</button>
								</div>
								<p class="rounded bg-base-100 p-2 text-sm">{entry.textRaw}</p>
							</div>

							{#if entry.postProcessed && entry.textRaw !== entry.textFinal}
								<div class="mb-3">
									<div class="mb-1 flex items-center justify-between">
										<p class="text-xs font-medium opacity-70">Enhanced:</p>
										<button
											class="btn btn-ghost btn-xs"
											onclick={() => copyToClipboard(entry.textFinal)}
										>
											<Copy class="size-3" />
											Copy
										</button>
									</div>
									<p class="rounded bg-base-100 p-2 text-sm">{entry.textFinal}</p>
								</div>
							{/if}

							<div class="flex justify-end">
								<button
									class="btn btn-outline btn-sm btn-error"
									onclick={() => handleDelete(entry)}
								>
									<Trash2 class="size-3.5" />
									Delete
								</button>
							</div>
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>
