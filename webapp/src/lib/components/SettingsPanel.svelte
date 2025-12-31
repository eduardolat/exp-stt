<script lang="ts">
	import { Settings2, Sparkles } from '@lucide/svelte';
	import GeneralSettings from './GeneralSettings.svelte';
	import AISettings from './AISettings.svelte';

	let activeSection = $state<'general' | 'ai'>('general');
</script>

<div class="space-y-4">
	<div role="tablist" class="tabs-box tabs">
		<button
			role="tab"
			class="tab gap-2"
			class:tab-active={activeSection === 'general'}
			onclick={() => (activeSection = 'general')}
		>
			<Settings2 class="size-4" />
			General
		</button>
		<button
			role="tab"
			class="tab gap-2"
			class:tab-active={activeSection === 'ai'}
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
