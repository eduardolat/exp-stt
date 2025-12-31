<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { StatusDisplay, HistoryPanel, SettingsPanel } from '$lib/components';
	import { Loader2 } from '@lucide/svelte';
</script>

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
{:else if store.currentTab === 'home'}
	<StatusDisplay />
{:else if store.currentTab === 'history'}
	<HistoryPanel />
{:else if store.currentTab === 'settings'}
	<SettingsPanel />
{/if}
