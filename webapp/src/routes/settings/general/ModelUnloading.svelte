<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Timer } from '@lucide/svelte';
	import { Card } from '$lib/components';

	const presetValues = [60, 120, 180, 240, 300, 600, 900, 1200, 1800, 3600, 7200];
	let customMode = $state(false);

	let isPreset = $derived(presetValues.includes(store.settings.modelUnloadSeconds));
	let showInput = $derived(customMode || !isPreset);
	let selectValue = $derived(showInput ? 'custom' : String(store.settings.modelUnloadSeconds));
</script>

<Card class="card-body break-inside-avoid">
	<h3 class="card-title text-base">
		<Timer class="size-4" />
		Model Unloading
	</h3>
	<div class="space-y-4">
		<label class="label cursor-pointer justify-between">
			<div>
				<span class="text-sm">Auto-unload model when idle</span>
				<p class="text-xs opacity-70">Reduces memory usage when not in use</p>
			</div>
			<input
				type="checkbox"
				class="toggle toggle-sm"
				checked={store.settings.modelUnloadEnable}
				onchange={(e) => store.updateSettings({ modelUnloadEnable: e.currentTarget.checked })}
			/>
		</label>

		{#if store.settings.modelUnloadEnable}
			<fieldset class="fieldset">
				<label class="label" for="unloadSeconds">
					<span class="label-text">Unload after</span>
				</label>
				<select
					id="unloadSeconds"
					class="select w-full select-sm"
					value={selectValue}
					onchange={(e) => {
						const value = e.currentTarget.value;
						if (value === 'custom') {
							customMode = true;
						} else {
							customMode = false;
							store.updateSettings({ modelUnloadSeconds: parseInt(value) });
						}
					}}
				>
					<option value="60">1 minute</option>
					<option value="120">2 minutes</option>
					<option value="180">3 minutes</option>
					<option value="240">4 minutes</option>
					<option value="300">5 minutes</option>
					<option value="600">10 minutes</option>
					<option value="900">15 minutes</option>
					<option value="1200">20 minutes</option>
					<option value="1800">30 minutes</option>
					<option value="3600">1 hour</option>
					<option value="7200">2 hours</option>
					<option value="custom">Custom duration</option>
				</select>

				{#if showInput}
					<div class="mt-2 flex items-center gap-2">
						<label class="input input-sm flex grow items-center gap-2">
							<input
								type="number"
								class="grow"
								placeholder="Custom seconds"
								value={store.settings.modelUnloadSeconds}
								onchange={(e) =>
									store.updateSettings({ modelUnloadSeconds: parseInt(e.currentTarget.value) })}
							/>
							<span class="text-xs opacity-50">seconds</span>
						</label>
					</div>
				{/if}
			</fieldset>
		{/if}
	</div>
</Card>
