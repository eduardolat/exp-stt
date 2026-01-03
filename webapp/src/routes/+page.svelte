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
	import { HistoryItem, Card, PageHeader } from '$lib/components';
	import WaylandShortcutAlert from '$lib/components/WaylandShortcutAlert.svelte';

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

	async function handleDelete(entry: (typeof store.history)[0]) {
		if (!confirm('Delete this transcription?')) return;
		await store.deleteHistoryEntry(entry.id);
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
			class="flex size-30 items-center justify-center rounded-full bg-base-300"
			class:animate-pulse={config.pulse}
			class:animate-spin={config.spin}
		>
			<IconComponent class="size-20" />
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
			Toggle Recording
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
				<Clock class="mb-2 size-8 opacity-50" />
				<p class="text-sm opacity-70">No transcriptions yet</p>
				<p class="text-xs opacity-50">Start recording to see your transcriptions here</p>
			</div>
		{:else}
			<div class="space-y-4">
				{#each recentHistory as entry (entry.id)}
					<HistoryItem {entry} onDelete={handleDelete} darker />
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
