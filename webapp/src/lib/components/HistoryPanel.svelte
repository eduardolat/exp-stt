<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Copy, Trash2, Clock, Sparkles, ChevronDown, ChevronUp } from '@lucide/svelte';
	import type { HistoryEntry } from '$lib/client.gen';

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
		if (confirm('Delete this transcription?')) {
			await store.deleteHistoryEntry(entry.id);
		}
	}

	async function handleClearAll() {
		if (confirm('Clear all history? This cannot be undone.')) {
			await store.clearHistory();
		}
	}
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h2 class="text-lg font-semibold">Transcription History</h2>
		{#if store.history.length > 0}
			<button class="btn btn-sm btn-outline btn-destructive" onclick={handleClearAll}>
				<Trash2 class="mr-1.5 size-3.5" />
				Clear All
			</button>
		{/if}
	</div>

	{#if store.history.length === 0}
		<div class="card flex flex-col items-center justify-center py-12 text-center">
			<Clock class="mb-3 size-10 text-muted-foreground" />
			<p class="text-muted-foreground">No transcriptions yet</p>
			<p class="text-xs text-muted-foreground">Your transcription history will appear here</p>
		</div>
	{:else}
		<div class="space-y-2">
			{#each store.history as entry (entry.id)}
				<div class="card overflow-hidden">
					<button
						class="flex w-full items-start gap-3 p-4 text-left hover:bg-muted/50"
						onclick={() => toggleExpand(entry.id)}
					>
						<div class="min-w-0 flex-1">
							<div class="mb-1 flex items-center gap-2 text-xs text-muted-foreground">
								<span>{formatTimestamp(entry.timestamp)}</span>
								<span>•</span>
								<span>{formatDuration(entry.durationMs)}</span>
								{#if entry.postProcessed}
									<span
										class="inline-flex items-center gap-1 rounded-full bg-purple-500/10 px-2 py-0.5 text-purple-600"
									>
										<Sparkles class="size-3" />
										AI Enhanced
									</span>
								{/if}
							</div>
							<p class="line-clamp-2 text-sm">{entry.textFinal}</p>
						</div>
						{#if expandedId === entry.id}
							<ChevronUp class="size-4 shrink-0 text-muted-foreground" />
						{:else}
							<ChevronDown class="size-4 shrink-0 text-muted-foreground" />
						{/if}
					</button>

					{#if expandedId === entry.id}
						<div class="border-t bg-muted/30 p-4">
							{#if entry.postProcessed && entry.textRaw !== entry.textFinal}
								<div class="mb-3">
									<p class="mb-1 text-xs font-medium text-muted-foreground">Original:</p>
									<p class="rounded bg-background p-2 text-sm">{entry.textRaw}</p>
								</div>
								<div class="mb-3">
									<p class="mb-1 text-xs font-medium text-muted-foreground">Enhanced:</p>
									<p class="rounded bg-background p-2 text-sm">{entry.textFinal}</p>
								</div>
							{:else}
								<p class="mb-3 text-sm">{entry.textFinal}</p>
							{/if}

							<div class="flex gap-2">
								<button
									class="btn btn-sm btn-outline"
									onclick={() => copyToClipboard(entry.textFinal)}
								>
									<Copy class="mr-1.5 size-3.5" />
									Copy
								</button>
								<button
									class="btn btn-sm btn-outline btn-destructive"
									onclick={() => handleDelete(entry)}
								>
									<Trash2 class="mr-1.5 size-3.5" />
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
