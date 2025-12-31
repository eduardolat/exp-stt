<script lang="ts">
	import '../app.css';
	import { store } from '$lib/store.svelte';
	import { onMount } from 'svelte';
	import { Mic, Settings, History, Wifi, WifiOff, Sun, Moon } from '@lucide/svelte';
	import { themeChange } from 'theme-change';

	let { children } = $props();

	onMount(() => {
		themeChange(false);
		store.initialize();
		return () => store.destroy();
	});
</script>

<div class="flex min-h-screen flex-col bg-base-100">
	<header class="navbar border-b border-base-300">
		<div class="mx-auto flex w-full max-w-4xl items-center justify-between px-4">
			<div class="flex items-center gap-3">
				<Mic class="size-6 text-primary" />
				<span class="text-lg font-semibold">Tribar Voice</span>
			</div>

			<div class="flex items-center gap-2">
				{#if store.isConnected}
					<div class="badge gap-1 badge-sm badge-success">
						<Wifi class="size-3" />
						Connected
					</div>
				{:else}
					<div class="badge gap-1 badge-sm badge-error">
						<WifiOff class="size-3" />
						Disconnected
					</div>
				{/if}

				<label class="btn swap btn-circle swap-rotate btn-ghost btn-sm">
					<input type="checkbox" data-toggle-theme="dark,light" data-act-class="swap-active" />
					<Sun class="swap-off size-4" />
					<Moon class="swap-on size-4" />
				</label>
			</div>
		</div>
	</header>

	<nav class="border-b border-base-300">
		<div class="mx-auto flex max-w-4xl gap-1 px-4">
			<button
				class="flex items-center gap-2 border-b-2 px-4 py-3 text-sm font-medium transition-colors {store.currentTab ===
				'home'
					? 'border-primary text-primary'
					: 'border-transparent opacity-70 hover:opacity-100'}"
				onclick={() => (store.currentTab = 'home')}
			>
				<Mic class="size-4" />
				Record
			</button>
			<button
				class="flex items-center gap-2 border-b-2 px-4 py-3 text-sm font-medium transition-colors {store.currentTab ===
				'history'
					? 'border-primary text-primary'
					: 'border-transparent opacity-70 hover:opacity-100'}"
				onclick={() => (store.currentTab = 'history')}
			>
				<History class="size-4" />
				History
			</button>
			<button
				class="flex items-center gap-2 border-b-2 px-4 py-3 text-sm font-medium transition-colors {store.currentTab ===
				'settings'
					? 'border-primary text-primary'
					: 'border-transparent opacity-70 hover:opacity-100'}"
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

	<footer class="border-t border-base-300 py-3 text-center text-xs opacity-70">
		<div class="mx-auto max-w-4xl px-4">
			{store.systemInfo.os} / {store.systemInfo.arch} / {store.systemInfo.displayServer}
		</div>
	</footer>
</div>
