<script lang="ts">
	import { Trash2, Sparkles } from '@lucide/svelte';
	import type { HistoryEntry } from '$lib/client.gen';
	import { Card, Modal, CopyButton } from '$lib/components';
	import { formatRelativeTime } from '$lib/utils/formatRelativeTime';

	interface Props {
		entry: HistoryEntry;
		onDelete: (entry: HistoryEntry) => void;
		darker?: boolean;
	}

	let { entry, onDelete, darker }: Props = $props();

	let modal: { open: () => void; close: () => void } | undefined = $state();

	function getAudioUrl(id: string): string {
		return `/api/v1/audio/${id}`;
	}

	function formatTimestamp(iso: string): string {
		const text = formatRelativeTime(iso);
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

<Card {darker} hoverable class="relative p-3">
	<button
		type="button"
		class="absolute inset-0 z-0 h-full w-full cursor-pointer opacity-0"
		aria-label="View transcription details"
		onclick={handleItemClick}
	></button>
	<div class="pointer-events-none relative z-10 mb-1.5 flex items-center justify-between gap-2">
		<div class="flex items-center gap-2 text-xs opacity-70">
			<span>{formatTimestamp(entry.timestamp)}</span>
			{#if entry.postProcessed}
				<span class="badge gap-1 badge-xs badge-secondary">
					<Sparkles class="size-2.5" />
					AI
				</span>
			{/if}
		</div>
		<div class="pointer-events-auto">
			<CopyButton text={entry.textFinal} showLabel />
		</div>
	</div>
	<p class="pointer-events-none relative z-0 line-clamp-2 text-sm">
		{@render textWithFallback(entry.textFinal)}
	</p>
</Card>

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

		<Card class="p-4">
			<p class="mb-1 text-xs font-medium opacity-70">Audio:</p>
			<audio controls class="w-full" preload="none" aria-label="Transcription audio">
				<source src={getAudioUrl(entry.id)} type="audio/wav" />
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

	{#snippet actions()}
		<div class="flex justify-end gap-2 pt-2">
			<button class="btn btn-outline btn-sm btn-error" onclick={handleDelete}>
				<Trash2 class="size-3.5" />
				Delete
			</button>
		</div>
	{/snippet}
</Modal>
