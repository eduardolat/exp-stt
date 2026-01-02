<script lang="ts">
	import { Trash2, Sparkles } from '@lucide/svelte';
	import { formatDistanceToNow } from 'date-fns';
	import type { HistoryEntry } from '$lib/client.gen';
	import { Modal, CopyButton } from '$lib/components';

	interface Props {
		entry: HistoryEntry;
		onDelete: (entry: HistoryEntry) => void;
	}

	let { entry, onDelete }: Props = $props();

	let modal: { open: () => void; close: () => void } | undefined = $state();

	function getAudioUrl(id: string): string {
		return `/api/v1/audio/${id}`;
	}

	function formatTimestamp(iso: string): string {
		const text = formatDistanceToNow(new Date(iso), { addSuffix: true });
		const capitalized = text.charAt(0).toUpperCase() + text.slice(1);
		return capitalized;
	}

	function handleItemClick() {
		modal?.open();
	}

	function handleDelete() {
		modal?.close();
		onDelete(entry);
	}
</script>

{#snippet textWithFallback(text: string)}
	{#if text.trim()}
		{text}
	{:else}
		<span class="opacity-70">(No text available)</span>
	{/if}
{/snippet}

<button
	class="
		w-full cursor-pointer rounded-lg border border-base-300 bg-base-300 p-3 text-left transition-all
		hover:border-primary hover:bg-base-100 hover:shadow-md
	"
	onclick={handleItemClick}
>
	<div class="mb-1.5 flex items-center justify-between gap-2">
		<div class="flex items-center gap-2 text-xs opacity-70">
			<span>{formatTimestamp(entry.timestamp)}</span>
			{#if entry.postProcessed}
				<span class="badge gap-1 badge-xs badge-secondary">
					<Sparkles class="size-2.5" />
					AI
				</span>
			{/if}
		</div>
		<CopyButton text={entry.textFinal} showLabel />
	</div>
	<p class="line-clamp-2 text-sm">
		{@render textWithFallback(entry.textFinal)}
	</p>
</button>

<Modal bind:this={modal} title="Transcription Details">
	<div class="space-y-4">
		<div class="flex items-center gap-2 text-sm opacity-70">
			<span>{formatTimestamp(entry.timestamp)}</span>
			<span>•</span>
			<span>Took {Math.floor(entry.durationMs / 1000)}s</span>
			{#if entry.postProcessed}
				<span class="ml-auto badge gap-1 badge-sm badge-secondary">
					<Sparkles class="size-3" />
					AI Enhanced
				</span>
			{/if}
		</div>

		<div>
			<p class="mb-1 text-xs font-medium opacity-70">Audio:</p>
			<audio controls class="w-full" preload="none">
				<source src={getAudioUrl(entry.id)} type="audio/wav" />
				Your browser does not support audio playback.
			</audio>
		</div>

		<div>
			<div class="mb-1 flex items-center justify-between">
				<p class="text-xs font-medium opacity-70">Original:</p>
				<CopyButton text={entry.textRaw} showLabel />
			</div>
			<p class="rounded bg-base-300 p-2 text-sm">{@render textWithFallback(entry.textRaw)}</p>
		</div>

		{#if entry.postProcessed && entry.textRaw !== entry.textFinal}
			<div>
				<div class="mb-1 flex items-center justify-between">
					<p class="text-xs font-medium opacity-70">Enhanced:</p>
					<CopyButton text={entry.textFinal} showLabel />
				</div>
				<p class="rounded bg-base-300 p-2 text-sm">{@render textWithFallback(entry.textFinal)}</p>
			</div>
		{/if}
	</div>

	{#snippet actions()}
		<div class="flex justify-end gap-2 pt-2">
			<button class="btn btn-outline btn-sm btn-error" onclick={handleDelete}>
				<Trash2 class="size-3.5" />
				Delete
			</button>
		</div>
	{/snippet}
</Modal>
