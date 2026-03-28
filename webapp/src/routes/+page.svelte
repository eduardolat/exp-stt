<script lang="ts">
	import { store } from '$lib/store.svelte';
	import {
		Mic,
		MicOff,
		Loader,
		Download,
		Cpu,
		Sparkles,
		Check,
		CircleAlert,
		Clock,
		ArrowRight,
		LayoutDashboard
	} from '@lucide/svelte';
	import { HistoryItem, Card, PageHeader, Modal, CopyButton } from '$lib/components';
	import WaylandShortcutAlert from '$lib/components/WaylandShortcutAlert.svelte';
	import type { HistoryEntry } from '$lib/client.gen';
	import { formatRelativeTime } from '$lib/utils/formatRelativeTime';
	import { Trash2 } from '@lucide/svelte';

	const statusConfig: Record<string, { icon: typeof Mic; pulse: boolean; spin: boolean }> = {
		unknown: { icon: CircleAlert, pulse: false, spin: false },
		unloaded: { icon: MicOff, pulse: false, spin: false },
		downloading: { icon: Download, pulse: true, spin: false },
		loading: { icon: Loader, pulse: false, spin: true },
		loaded: { icon: Check, pulse: false, spin: false },
		listening: { icon: Mic, pulse: true, spin: false },
		transcribing: { icon: Cpu, pulse: true, spin: false },
		post_processing: { icon: Sparkles, pulse: true, spin: false }
	};

	let config = $derived(statusConfig[store.status] ?? statusConfig.unknown);
	let IconComponent = $derived(config.icon);
	let recentHistory = $derived(store.history.slice(0, 3));

	let detailsModal: Modal | undefined = $state();
	let selectedEntry: HistoryEntry | undefined = $state();

	async function handleDelete() {
		if (!selectedEntry) return;
		// Close modal first
		detailsModal?.close();

		if (!confirm('Delete this transcription?')) return;

		await store.deleteHistoryEntry(selectedEntry.id);
		selectedEntry = undefined;
	}
</script>

<svelte:head>
	<title>Dashboard - Tribar Voice</title>
</svelte:head>

<PageHeader
	icon={LayoutDashboard}
	title="Dashboard"
	description="Overview of your transcription activity and system status"
/>

<div class="columns-2 gap-4 space-y-4">
	<!-- Recording Control Panel -->
	<Card class="card-body flex break-inside-avoid flex-col items-center justify-between gap-4">
		<!-- Title -->
		<h3 class="card-title justify-center text-base">Recording</h3>

		<!-- Status Icon + Label -->
		<div
			class="flex size-20 items-center justify-center rounded-full bg-base-300"
			class:animate-pulse={config.pulse}
			class:animate-spin={config.spin}
		>
			<IconComponent class="size-10" />
		</div>
		<p class="font-medium">{store.statusLabel}</p>

		{#if store.status === 'downloading'}
			<div class="w-full max-w-xs">
				<div class="mb-1 flex justify-between text-xs opacity-70">
					<span class="truncate">{store.downloadProgress.fileName}</span>
					<span>{store.downloadProgress.percent.toFixed(1)}%</span>
				</div>
				<progress class="progress progress-primary" value={store.downloadProgress.percent} max="100"
				></progress>
			</div>
		{/if}

		<!-- Action Button + Shortcut -->
		<button
			class="btn btn-lg"
			class:btn-error={store.isRecording}
			class:btn-primary={!store.isRecording}
			onclick={() => store.toggleRecording()}
		>
			<Mic class="size-5" />
			{store.isRecording ? 'Stop Recording' : 'Start Recording'}
		</button>

		<WaylandShortcutAlert>
			{#if store.shortcut.key}
				<p class="text-xs opacity-70">
					Shortcut:
					{#if store.shortcut.modifiers.length > 0}
						<kbd class="kbd kbd-sm">{store.shortcut.modifiers.join(' + ')}</kbd>
						+
					{/if}
					<kbd class="kbd kbd-sm">{store.shortcut.key}</kbd>
				</p>
			{/if}
		</WaylandShortcutAlert>
	</Card>

	<!-- Recent History Panel -->
	<Card class="card-body break-inside-avoid">
		<h3 class="mb-2 card-title text-base">Recent Transcriptions</h3>

		{#if recentHistory.length === 0}
			<div class="flex flex-col items-center justify-center py-26 text-center">
				<Clock class="mb-2 size-10 opacity-50" />
				<p class="mb-2 text-sm opacity-70">No transcriptions yet</p>
				<p class="text-xs opacity-50">Start recording to see your transcriptions here</p>
			</div>
		{:else}
			<div class="space-y-4">
				{#each recentHistory as entry (entry.id)}
					<HistoryItem
						{entry}
						onSelect={(e) => {
							selectedEntry = e;
							detailsModal?.open();
						}}
						darker
					/>
				{/each}
			</div>

			<div class="flex justify-end">
				<a href="#/history" class="btn btn-ghost btn-sm">
					View all history
					<ArrowRight class="size-4" />
				</a>
			</div>
		{/if}
	</Card>
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
