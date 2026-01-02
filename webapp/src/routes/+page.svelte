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
		Info
	} from '@lucide/svelte';
	import { HistoryItem, Modal, CopyButton } from '$lib/components';
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

<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
	<!-- Recording Control Panel -->
	<div>
		<div class="card bg-base-200">
			<div class="card-body flex flex-col justify-between gap-6">
				<!-- Title -->
				<div class="text-center">
					<h3 class="card-title justify-center text-base">Recording</h3>
				</div>

				<!-- Status Icon + Label -->
				<div class="flex flex-col items-center gap-3">
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
							<progress
								class="progress progress-primary"
								value={store.downloadProgress.percent}
								max="100"
							></progress>
						</div>
					{/if}
				</div>

				<!-- Action Button + Shortcut -->
				<div class="flex flex-col items-center gap-4">
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
				</div>
			</div>
		</div>
	</div>

	<!-- Recent History Panel -->
	<div>
		<div class="card bg-base-200">
			<div class="card-body">
				<h3 class="mb-2 card-title text-base">Recent Transcriptions</h3>

				{#if recentHistory.length === 0}
					<div class="flex flex-col items-center justify-center py-26 text-center">
						<Clock class="mb-2 size-8 opacity-50" />
						<p class="text-sm opacity-70">No transcriptions yet</p>
						<p class="text-xs opacity-50">Start recording to see your transcriptions here</p>
					</div>
				{:else}
					<div class="space-y-2">
						{#each recentHistory as entry (entry.id)}
							<HistoryItem {entry} onDelete={handleDelete} />
						{/each}
					</div>

					<div class="mt-4 flex justify-end">
						<a href="#/history" class="btn gap-1 btn-ghost btn-sm">
							View All History
							<ArrowRight class="size-4" />
						</a>
					</div>
				{/if}
			</div>
		</div>
	</div>
</div>
