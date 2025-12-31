<script lang="ts">
	import '../app.css';
	import { store } from '$lib/store.svelte';
	import { onMount } from 'svelte';
	import { Mic, Settings, History, Wifi, WifiOff } from '@lucide/svelte';

	let { children } = $props();

	onMount(() => {
		store.initialize();
		return () => store.destroy();
	});
</script>

<div class="bg-background flex min-h-screen flex-col">
	<header class="border-b">
		<div class="mx-auto flex h-14 max-w-4xl items-center justify-between px-4">
			<div class="flex items-center gap-3">
				<Mic class="size-6 text-primary" />
				<h1 class="text-lg font-semibold">Tribar Voice</h1>
			</div>

			<div class="flex items-center gap-2">
				{#if store.isConnected}
					<div
						class="flex items-center gap-1.5 rounded-full bg-green-500/10 px-2 py-1 text-xs text-green-600"
					>
						<Wifi class="size-3" />
						<span>Connected</span>
					</div>
				{:else}
					<div
						class="flex items-center gap-1.5 rounded-full bg-red-500/10 px-2 py-1 text-xs text-red-600"
					>
						<WifiOff class="size-3" />
						<span>Disconnected</span>
					</div>
				{/if}
			</div>
		</div>
	</header>

	<nav class="border-b">
		<div class="mx-auto flex max-w-4xl gap-1 px-4">
			<button
				class="flex items-center gap-2 border-b-2 px-4 py-3 text-sm font-medium transition-colors {store.currentTab ===
				'home'
					? 'border-primary text-primary'
					: 'text-muted-foreground hover:text-foreground border-transparent'}"
				onclick={() => (store.currentTab = 'home')}
			>
				<Mic class="size-4" />
				Record
			</button>
			<button
				class="flex items-center gap-2 border-b-2 px-4 py-3 text-sm font-medium transition-colors {store.currentTab ===
				'history'
					? 'border-primary text-primary'
					: 'text-muted-foreground hover:text-foreground border-transparent'}"
				onclick={() => (store.currentTab = 'history')}
			>
				<History class="size-4" />
				History
			</button>
			<button
				class="flex items-center gap-2 border-b-2 px-4 py-3 text-sm font-medium transition-colors {store.currentTab ===
				'settings'
					? 'border-primary text-primary'
					: 'text-muted-foreground hover:text-foreground border-transparent'}"
				onclick={() => (store.currentTab = 'settings')}
			>
				<Settings class="size-4" />
				Settings
			</button>
		</div>
	</nav>

	<main class="mx-auto w-full max-w-4xl flex-1 px-4 py-6">
		{@render children()}
	</main>

	<footer class="text-muted-foreground border-t py-3 text-center text-xs">
		<div class="mx-auto max-w-4xl px-4">
			{store.systemInfo.os} / {store.systemInfo.arch} / {store.systemInfo.displayServer}
		</div>
	</footer>
</div>
