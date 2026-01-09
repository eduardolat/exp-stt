<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Card } from '$lib/components';
	import { Workflow, AlertTriangle } from '@lucide/svelte';
</script>

<Card>
	<h3 class="card-title text-sm">
		<Workflow class="size-4" />
		Advanced Mode
	</h3>
	<p class="mt-1 text-xs opacity-70">
		Enable custom workflow automation with a visual node editor.
	</p>

	<div class="mt-4">
		<label class="flex cursor-pointer items-center justify-between">
			<span class="text-sm">Enable Advanced Mode</span>
			<input
				type="checkbox"
				class="toggle toggle-primary"
				checked={store.settings.advancedMode}
				onchange={(e) => {
					const target = e.target as HTMLInputElement;
					store.updateSettings({ advancedMode: target.checked });
				}}
			/>
		</label>
	</div>

	{#if store.settings.advancedMode}
		<div class="mt-4 flex items-start gap-2 rounded-lg bg-warning/10 p-3 text-sm text-warning">
			<AlertTriangle class="size-4 shrink-0" />
			<p>
				In Advanced Mode, transcription is processed through your active workflow instead of the
				default pipeline. Make sure to configure your workflow correctly.
			</p>
		</div>
	{/if}
</Card>
