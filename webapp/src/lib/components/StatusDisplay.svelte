<script lang="ts">
	import { store } from '$lib/store.svelte';
	import {
		Mic,
		MicOff,
		Loader2,
		Download,
		Cpu,
		Sparkles,
		Check,
		AlertCircle
	} from '@lucide/svelte';

	const statusConfig: Record<
		string,
		{ icon: typeof Mic; color: string; bgColor: string; pulse: boolean }
	> = {
		unknown: {
			icon: AlertCircle,
			color: 'text-muted-foreground',
			bgColor: 'bg-muted',
			pulse: false
		},
		unloaded: { icon: MicOff, color: 'text-muted-foreground', bgColor: 'bg-muted', pulse: false },
		downloading: { icon: Download, color: 'text-blue-600', bgColor: 'bg-blue-500/10', pulse: true },
		loading: { icon: Loader2, color: 'text-yellow-600', bgColor: 'bg-yellow-500/10', pulse: true },
		loaded: { icon: Check, color: 'text-green-600', bgColor: 'bg-green-500/10', pulse: false },
		listening: { icon: Mic, color: 'text-red-600', bgColor: 'bg-red-500/10', pulse: true },
		transcribing: { icon: Cpu, color: 'text-purple-600', bgColor: 'bg-purple-500/10', pulse: true },
		post_processing: {
			icon: Sparkles,
			color: 'text-orange-600',
			bgColor: 'bg-orange-500/10',
			pulse: true
		}
	};

	let config = $derived(statusConfig[store.status] ?? statusConfig.unknown);
	let IconComponent = $derived(config.icon);
</script>

<div class="card p-6">
	<div class="flex flex-col items-center gap-4">
		<div
			class="relative flex size-24 items-center justify-center rounded-full {config.bgColor}"
			class:animate-pulse={config.pulse}
		>
			<IconComponent class="size-12 {config.color}" />
		</div>

		<div class="text-center">
			<h2 class="text-xl font-semibold">{store.statusLabel}</h2>
			{#if store.status === 'downloading'}
				<div class="mt-3 w-64">
					<div class="mb-1 flex justify-between text-xs text-muted-foreground">
						<span>{store.downloadProgress.fileName}</span>
						<span>{store.downloadProgress.percent.toFixed(1)}%</span>
					</div>
					<div class="progress h-2">
						<div class="progress-bar" style="width: {store.downloadProgress.percent}%"></div>
					</div>
				</div>
			{/if}
		</div>

		<button
			class="btn btn-lg mt-2"
			class:btn-destructive={store.isRecording}
			disabled={store.isProcessing || store.status === 'unknown' || store.status === 'unloaded'}
			onclick={() => store.toggleRecording()}
		>
			{#if store.isRecording}
				<MicOff class="mr-2 size-5" />
				Stop Recording
			{:else}
				<Mic class="mr-2 size-5" />
				Start Recording
			{/if}
		</button>

		{#if store.shortcut.key}
			<p class="text-xs text-muted-foreground">
				Shortcut:
				{#if store.shortcut.modifiers.length > 0}
					<kbd class="kbd">{store.shortcut.modifiers.join(' + ')}</kbd>
					+
				{/if}
				<kbd class="kbd">{store.shortcut.key}</kbd>
			</p>
		{/if}
	</div>
</div>
