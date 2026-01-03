<script lang="ts">
	import { Modal, Card, CopyButton } from '$lib/components';
	import { Sparkles, Trash2 } from '@lucide/svelte';
	import { formatTimeAgo } from '$lib/formatters';
	import type { HistoryEntry } from '$lib/client.gen';
	import { store } from '$lib/store.svelte';

	interface Props {
		entry: HistoryEntry | null;
		onClose: () => void;
	}

	let { entry, onClose }: Props = $props();

	let modal: Modal | undefined = $state();

	$effect(() => {
		if (entry) {
			modal?.open();
		} else {
			modal?.close();
		}
	});

	function handleInternalClose() {
		// When the user clicks the X or backdrop, we need to notify the parent
		// to clear the entry, otherwise the modal will be "open" in props but closed in DOM.
		// However, since `modal.close()` is called programmatically when `entry` becomes null,
		// we need to distinguish.
		// Actually, simpler approach: The parent controls `entry`. If `entry` is set, modal opens.
		// If user closes modal, we tell parent to set `entry` to null.
		onClose();
	}

	async function handleDelete() {
		if (!entry) return;
		if (!confirm('Delete this transcription?')) return;

		const id = entry.id;
		handleInternalClose(); // Close first to feel responsive
		await store.deleteHistoryEntry(id);
	}

	function formatTimestamp(iso: string): string {
		const text = formatTimeAgo(iso);
		const capitalized = text.charAt(0).toUpperCase() + text.slice(1);
		return capitalized;
	}

	function getAudioUrl(id: string): string {
		return `/api/v1/audio/${id}`;
	}
</script>

{#snippet textWithFallback(text: string)}
	{#if text.trim()}
		{text}
	{:else}
		<span class="opacity-70">(No text available)</span>
	{/if}
{/snippet}

<!-- We need to bind `this` to access open/close methods, but we also want to intercept close events from the UI -->
<!-- The Modal component exposes open() and close(). When the modal is closed via UI (X button or backdrop), we need to notify parent. -->
<!-- However, the current Modal implementation doesn't seem to emit a close event. -->
<!-- Let's check Modal.svelte again. -->

<Modal bind:this={modal} title="Transcription Details" onClose={handleInternalClose}>
	{#snippet children()}
		{#if entry}
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
					<audio controls class="w-full" preload="none" src={getAudioUrl(entry.id)}>
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
	{/snippet}

	{#snippet actions()}
		{#if entry}
			<div class="flex justify-end gap-2 pt-2">
				<button
					class="btn btn-outline btn-sm btn-error"
					onclick={handleDelete}
				>
					<Trash2 class="size-3.5" />
					Delete
				</button>
			</div>
		{/if}
	{/snippet}
</Modal>
