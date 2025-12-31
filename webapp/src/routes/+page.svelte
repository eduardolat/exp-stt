<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Mic, MicOff, Loader, Download, Cpu, Sparkles, Check, CircleAlert } from '@lucide/svelte';

	const statusConfig: Record<string, { icon: typeof Mic; pulse: boolean }> = {
		unknown: { icon: CircleAlert, pulse: false },
		unloaded: { icon: MicOff, pulse: false },
		downloading: { icon: Download, pulse: true },
		loading: { icon: Loader, pulse: true },
		loaded: { icon: Check, pulse: false },
		listening: { icon: Mic, pulse: true },
		transcribing: { icon: Cpu, pulse: true },
		post_processing: { icon: Sparkles, pulse: true }
	};

	let config = $derived(statusConfig[store.status] ?? statusConfig.unknown);
	let IconComponent = $derived(config.icon);
</script>

<div class="card bg-base-200">
	<div class="card-body items-center text-center">
		<div
			class="flex size-24 items-center justify-center rounded-full bg-base-300"
			class:animate-pulse={config.pulse}
		>
			<IconComponent class="size-12" />
		</div>

		<h2 class="card-title">{store.statusLabel}</h2>

		{#if store.status === 'downloading'}
			<div class="mt-2 w-64">
				<div class="mb-1 flex justify-between text-xs opacity-70">
					<span>{store.downloadProgress.fileName}</span>
					<span>{store.downloadProgress.percent.toFixed(1)}%</span>
				</div>
				<progress class="progress progress-primary" value={store.downloadProgress.percent} max="100"
				></progress>
			</div>
		{/if}

		<div class="mt-4 card-actions">
			<button
				class="btn btn-lg"
				class:btn-error={store.isRecording}
				class:btn-primary={!store.isRecording}
				onclick={() => store.toggleRecording()}
			>
				<Mic class="size-5" />
				Toggle Recording
			</button>
		</div>

		{#if store.shortcut.key}
			<p class="mt-2 text-xs opacity-70">
				Shortcut:
				{#if store.shortcut.modifiers.length > 0}
					<kbd class="kbd kbd-sm">{store.shortcut.modifiers.join(' + ')}</kbd>
					+
				{/if}
				<kbd class="kbd kbd-sm">{store.shortcut.key}</kbd>
			</p>
		{/if}
	</div>
</div>
