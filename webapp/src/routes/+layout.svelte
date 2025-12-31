<script lang="ts">
	import '../app.css';
	import { store } from '$lib/store.svelte';
	import { onMount } from 'svelte';
	import { Mic, Settings, History, Palette, Loader, Sun, Moon, Eclipse } from '@lucide/svelte';
	import { themeChange } from 'theme-change';
	import { page } from '$app/state';
	import { AnimatedLogo, type AppStatus } from '$lib/components';

	let { children } = $props();

	const themes = [
		{ value: 'light', label: 'Light', icon: Sun },
		{ value: 'dim', label: 'Dark', icon: Moon },
		{ value: 'dracula', label: 'Dracula', icon: Eclipse }
	];

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
				<AnimatedLogo status={store.status as AppStatus} size={24} />
				<span class="text-lg font-semibold">Tribar Voice</span>
			</div>

			<div class="dropdown dropdown-end">
				<div tabindex="0" role="button" class="btn btn-ghost btn-sm">
					<Palette class="mr-1 size-4" />
					<span>Theme</span>
				</div>
				<ul
					tabindex="-1"
					class="dropdown-content menu z-10 w-36 rounded-box bg-base-200 p-2 shadow"
				>
					{#each themes as theme}
						<li>
							<button data-set-theme={theme.value}>
								<theme.icon class="size-4" />
								{theme.label}
							</button>
						</li>
					{/each}
				</ul>
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
				<Loader class="size-8 animate-spin text-primary" />
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
			{store.systemInfo.os} / {store.systemInfo.arch}
			{#if store.systemInfo.os === 'Linux'}
				/ {store.systemInfo.displayServer}
			{/if}
		</div>
	</footer>
</div>
