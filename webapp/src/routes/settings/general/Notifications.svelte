<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Bell } from '@lucide/svelte';
	import { Card } from '$lib/components';
</script>

<Card class="card-body break-inside-avoid">
	<h3 class="card-title text-base">
		<Bell class="size-4" />
		Notifications
	</h3>
	<div class="space-y-4">
		<label class="label cursor-pointer justify-between py-2">
			<span class="text-sm">Notify on errors</span>
			<input
				type="checkbox"
				class="toggle toggle-sm"
				checked={store.settings.notifyOnError}
				onchange={(e) => store.updateSettings({ notifyOnError: e.currentTarget.checked })}
			/>
		</label>
		<label class="label cursor-pointer justify-between py-2">
			<span class="text-sm">Notify on recording start</span>
			<input
				type="checkbox"
				class="toggle toggle-sm"
				checked={store.settings.notifyOnStart}
				onchange={(e) => store.updateSettings({ notifyOnStart: e.currentTarget.checked })}
			/>
		</label>
		<label class="label cursor-pointer justify-between py-2">
			<span class="text-sm">Notify on transcription complete</span>
			<input
				type="checkbox"
				class="toggle toggle-sm"
				checked={store.settings.notifyOnFinish}
				onchange={(e) => store.updateSettings({ notifyOnFinish: e.currentTarget.checked })}
			/>
		</label>
	</div>
</Card>
