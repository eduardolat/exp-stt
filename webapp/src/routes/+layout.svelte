<script lang="ts">
	import '../app.css';
	import { store } from '$lib/store.svelte';
	import { onMount } from 'svelte';
	import { Mic, Settings, History, Wifi, WifiOff, Sun, Moon, Loader2 } from '@lucide/svelte';
	import { themeChange } from 'theme-change';
	import { page } from '$app/state';

	let { children } = $props();

	function isActive(path: string): boolean {
		const currentPath = page.url.hash.slice(1) || '/';
		if (path === '/') {
			return currentPath === '/' || currentPath === '';
		}
		return currentPath.startsWith(path);
	}

	onMount(() => {
		themeChange(false);
		store.initialize();
		return () => store.destroy();
	});
</script>

<div class="flex h-screen flex-col bg-base-100">
	<header class="navbar shrink-0 border-b border-base-300">
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

	<nav class="shrink-0 border-b border-base-300">
		<div class="mx-auto flex max-w-4xl gap-1 px-4">
			<a
				href="#/"
				class="flex items-center gap-2 border-b-2 px-4 py-3 text-sm font-medium transition-colors {isActive(
					'/'
				)
					? 'border-primary text-primary'
					: 'border-transparent opacity-70 hover:opacity-100'}"
			>
				<Mic class="size-4" />
				Record
			</a>
			<a
				href="#/history"
				class="flex items-center gap-2 border-b-2 px-4 py-3 text-sm font-medium transition-colors {isActive(
					'/history'
				)
					? 'border-primary text-primary'
					: 'border-transparent opacity-70 hover:opacity-100'}"
			>
				<History class="size-4" />
				History
			</a>
			<a
				href="#/settings"
				class="flex items-center gap-2 border-b-2 px-4 py-3 text-sm font-medium transition-colors {isActive(
					'/settings'
				)
					? 'border-primary text-primary'
					: 'border-transparent opacity-70 hover:opacity-100'}"
			>
				<Settings class="size-4" />
				Settings
			</a>
		</div>
	</nav>

	<main class="mx-auto w-full max-w-4xl flex-1 overflow-y-auto px-4 py-6">
		{#if store.isLoading}
			<div class="flex flex-col items-center justify-center py-20">
				<Loader2 class="size-8 animate-spin text-primary" />
				<p class="text-muted-foreground mt-4">Connecting to server...</p>
			</div>
		{:else if store.error}
			<div class="card flex flex-col items-center justify-center py-12 text-center">
				<p class="text-destructive">{store.error}</p>
				<button class="btn mt-4 btn-sm" onclick={() => store.initialize()}>Retry</button>
			</div>
		{:else}
			{@render children()}
		{/if}
	</main>

	<footer class="shrink-0 border-t border-base-300 py-3 text-center text-xs opacity-70">
		<div class="mx-auto max-w-4xl px-4">
			{store.systemInfo.os} / {store.systemInfo.arch} / {store.systemInfo.displayServer}
		</div>
	</footer>
</div>
