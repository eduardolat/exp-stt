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
	<table class="w-full">
		<tbody>
			<tr>
				<td class="py-2">
					<span class="text-sm">Notify on errors</span>
				</td>
				<td class="py-2 text-right">
					<input
						type="checkbox"
						class="toggle toggle-sm"
						checked={store.settings.notifyOnError}
						onchange={(e) => store.updateSettings({ notifyOnError: e.currentTarget.checked })}
					/>
				</td>
			</tr>
			<tr>
				<td class="py-2">
					<span class="text-sm">Notify on recording start</span>
				</td>
				<td class="py-2 text-right">
					<input
						type="checkbox"
						class="toggle toggle-sm"
						checked={store.settings.notifyOnStart}
						onchange={(e) => store.updateSettings({ notifyOnStart: e.currentTarget.checked })}
					/>
				</td>
			</tr>
			<tr>
				<td class="py-2">
					<span class="text-sm">Notify on transcription complete</span>
				</td>
				<td class="py-2 text-right">
					<input
						type="checkbox"
						class="toggle toggle-sm"
						checked={store.settings.notifyOnFinish}
						onchange={(e) => store.updateSettings({ notifyOnFinish: e.currentTarget.checked })}
					/>
				</td>
			</tr>
		</tbody>
	</table>
</Card>
