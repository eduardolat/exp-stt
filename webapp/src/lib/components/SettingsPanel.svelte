<script lang="ts">
	import { Settings2, Sparkles } from '@lucide/svelte';
	import GeneralSettings from './GeneralSettings.svelte';
	import AISettings from './AISettings.svelte';

	let activeSection = $state<'general' | 'ai'>('general');
</script>

<div class="space-y-4">
	<div class="flex gap-2 border-b pb-2">
		<button
			class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium transition-colors"
			class:bg-primary={activeSection === 'general'}
			class:text-primary-foreground={activeSection === 'general'}
			class:text-muted-foreground={activeSection !== 'general'}
			class:hover:bg-muted={activeSection !== 'general'}
			onclick={() => (activeSection = 'general')}
		>
			<Settings2 class="size-4" />
			General
		</button>
		<button
			class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium transition-colors"
			class:bg-primary={activeSection === 'ai'}
			class:text-primary-foreground={activeSection === 'ai'}
			class:text-muted-foreground={activeSection !== 'ai'}
			class:hover:bg-muted={activeSection !== 'ai'}
			onclick={() => (activeSection = 'ai')}
		>
			<Sparkles class="size-4" />
			AI Enhancement
		</button>
	</div>

	{#if activeSection === 'general'}
		<GeneralSettings />
	{:else}
		<AISettings />
	{/if}
</div>
