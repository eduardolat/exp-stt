<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Mic } from '@lucide/svelte';
	import { Card } from '$lib/components';
</script>

<Card class="card-body break-inside-avoid">
	<h3 class="card-title text-base">
		<Mic class="size-4" />
		Input Device
	</h3>
	<fieldset class="fieldset">
		<select
			class="select w-full select-sm"
			value={store.settings.inputDevice}
			onchange={(e) => store.updateSettings({ inputDevice: e.currentTarget.value })}
		>
			<option value="default">System Default</option>
			{#each store.inputDevices as device (device.id)}
				<option value={device.id}>
					{device.name}
					{device.isDefault ? '(Default)' : ''}
				</option>
			{/each}
		</select>
	</fieldset>
</Card>
