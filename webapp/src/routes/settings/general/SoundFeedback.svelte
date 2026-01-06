<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Volume2 } from '@lucide/svelte';
	import { Card } from '$lib/components';
	import SoundFeedbackSoundPicker from './SoundFeedbackSoundPicker.svelte';
</script>

<Card class="card-body break-inside-avoid">
	<h3 class="card-title text-base">
		<Volume2 class="size-4" />
		Sound Feedback
	</h3>
	<div class="space-y-4">
		<label class="label cursor-pointer justify-between">
			<div class="flex items-center gap-2">
				<span class="text-sm">Enable sound feedback</span>
			</div>
			<input
				type="checkbox"
				class="toggle toggle-sm"
				checked={store.settings.soundFeedbackEnable}
				onchange={(e) => store.updateSettings({ soundFeedbackEnable: e.currentTarget.checked })}
			/>
		</label>

		{#if store.settings.soundFeedbackEnable}
			<fieldset class="fieldset">
				<label class="label" for="volume">
					<span class="label-text">Volume: {store.settings.soundFeedbackVolume}%</span>
				</label>
				<input
					type="range"
					id="volume"
					class="range w-full range-sm"
					min="0"
					max="100"
					value={store.settings.soundFeedbackVolume}
					onchange={(e) =>
						store.updateSettings({ soundFeedbackVolume: parseInt(e.currentTarget.value) })}
				/>
			</fieldset>

			<div>
				<SoundFeedbackSoundPicker
					label="Recording Sound"
					soundType="record"
					options={['1', '2', '3', '4', '5', '6', '7', '8', '9']}
					value={store.settings.soundFeedbackRecordId}
					onchange={(value) => store.updateSettings({ soundFeedbackRecordId: value })}
				/>

				<SoundFeedbackSoundPicker
					label="Success Sound"
					soundType="success"
					options={['1', '2', '3', '4']}
					value={store.settings.soundFeedbackSuccessId}
					onchange={(value) => store.updateSettings({ soundFeedbackSuccessId: value })}
				/>

				<SoundFeedbackSoundPicker
					label="Error Sound"
					soundType="error"
					options={['1', '2', '3', '4', '5', '6', '7', '8']}
					value={store.settings.soundFeedbackErrorId}
					onchange={(value) => store.updateSettings({ soundFeedbackErrorId: value })}
				/>
			</div>
		{/if}
	</div>
</Card>
