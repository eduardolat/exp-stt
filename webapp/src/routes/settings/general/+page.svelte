<script lang="ts">
	import { store } from '$lib/store.svelte';
	import {
		Bell,
		BellOff,
		Volume2,
		VolumeX,
		Mic,
		Clipboard,
		ClipboardPaste,
		Ghost,
		Timer,
		Keyboard
	} from '@lucide/svelte';
	import { ShortcutInput } from '$lib/components';
</script>

<svelte:head>
	<title>General Settings - Tribar Voice</title>
</svelte:head>

<div class="space-y-6">
	<div class="card bg-base-200">
		<div class="card-body">
			<h3 class="card-title text-base">
				<Mic class="size-4" />
				Input Device
			</h3>
			<fieldset class="fieldset">
				<select
					class="select w-full"
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
		</div>
	</div>

	<div class="card bg-base-200">
		<div class="card-body">
			<h3 class="card-title text-base">
				<Clipboard class="size-4" />
				Output Mode
			</h3>
			<div class="space-y-2">
				<label
					class="flex cursor-pointer items-start gap-3 rounded-lg border border-base-300 p-3 hover:bg-base-300"
				>
					<input
						type="radio"
						name="outputMode"
						class="radio mt-0.5 radio-sm"
						value="copy_only"
						checked={store.settings.outputMode === 'copy_only'}
						onchange={() => store.updateSettings({ outputMode: 'copy_only' })}
					/>
					<div class="flex-1">
						<div class="flex items-center gap-2 font-medium">
							<Clipboard class="size-4" />
							Copy Only
						</div>
						<p class="text-xs opacity-70">Copies transcription to clipboard without pasting</p>
					</div>
				</label>

				<label
					class="flex cursor-pointer items-start gap-3 rounded-lg border border-base-300 p-3 hover:bg-base-300"
				>
					<input
						type="radio"
						name="outputMode"
						class="radio mt-0.5 radio-sm"
						value="copy_paste"
						checked={store.settings.outputMode === 'copy_paste'}
						onchange={() => store.updateSettings({ outputMode: 'copy_paste' })}
					/>
					<div class="flex-1">
						<div class="flex items-center gap-2 font-medium">
							<ClipboardPaste class="size-4" />
							Copy & Paste
						</div>
						<p class="text-xs opacity-70">Copies to clipboard and pastes at cursor position</p>
					</div>
				</label>

				<label
					class="flex cursor-pointer items-start gap-3 rounded-lg border border-base-300 p-3 hover:bg-base-300"
				>
					<input
						type="radio"
						name="outputMode"
						class="radio mt-0.5 radio-sm"
						value="ghost_paste"
						checked={store.settings.outputMode === 'ghost_paste'}
						onchange={() => store.updateSettings({ outputMode: 'ghost_paste' })}
					/>
					<div class="flex-1">
						<div class="flex items-center gap-2 font-medium">
							<Ghost class="size-4" />
							Ghost Paste
						</div>
						<p class="text-xs opacity-70">Pastes without modifying your clipboard contents</p>
					</div>
				</label>

				<label class="label cursor-pointer justify-start gap-3 pt-2">
					<input
						type="checkbox"
						class="toggle toggle-sm"
						checked={store.settings.outputTrailingSpace}
						onchange={(e) => store.updateSettings({ outputTrailingSpace: e.currentTarget.checked })}
					/>
					<span class="text-sm">Add trailing space after paste</span>
				</label>
			</div>
		</div>
	</div>

	<div class="card bg-base-200">
		<div class="card-body">
			<h3 class="card-title text-base">
				<Keyboard class="size-4" />
				Keyboard Shortcut
			</h3>
			<ShortcutInput />
		</div>
	</div>

	<div class="card bg-base-200">
		<div class="card-body">
			<h3 class="card-title text-base">
				<Bell class="size-4" />
				Notifications
			</h3>
			<div class="space-y-2">
				<label class="label cursor-pointer justify-between">
					<div class="flex items-center gap-2">
						<BellOff class="size-4 opacity-70" />
						<span class="text-sm">Notify on errors</span>
					</div>
					<input
						type="checkbox"
						class="toggle toggle-sm"
						checked={store.settings.notifyOnError}
						onchange={(e) => store.updateSettings({ notifyOnError: e.currentTarget.checked })}
					/>
				</label>
				<label class="label cursor-pointer justify-between">
					<div class="flex items-center gap-2">
						<Bell class="size-4 opacity-70" />
						<span class="text-sm">Notify on recording start</span>
					</div>
					<input
						type="checkbox"
						class="toggle toggle-sm"
						checked={store.settings.notifyOnStart}
						onchange={(e) => store.updateSettings({ notifyOnStart: e.currentTarget.checked })}
					/>
				</label>
				<label class="label cursor-pointer justify-between">
					<div class="flex items-center gap-2">
						<Bell class="size-4 opacity-70" />
						<span class="text-sm">Notify on transcription complete</span>
					</div>
					<input
						type="checkbox"
						class="toggle toggle-sm"
						checked={store.settings.notifyOnFinish}
						onchange={(e) => store.updateSettings({ notifyOnFinish: e.currentTarget.checked })}
					/>
				</label>
			</div>
		</div>
	</div>

	<div class="card bg-base-200">
		<div class="card-body">
			<h3 class="card-title text-base">
				<Volume2 class="size-4" />
				Sound Feedback
			</h3>
			<div class="space-y-4">
				<label class="label cursor-pointer justify-between">
					<div class="flex items-center gap-2">
						{#if store.settings.soundFeedbackEnable}
							<Volume2 class="size-4 opacity-70" />
						{:else}
							<VolumeX class="size-4 opacity-70" />
						{/if}
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
					<fieldset class="fieldset pl-6">
						<label class="label" for="volume">
							<span class="label-text">Volume: {store.settings.soundFeedbackVolume}%</span>
						</label>
						<input
							type="range"
							id="volume"
							class="range range-sm"
							min="0"
							max="100"
							value={store.settings.soundFeedbackVolume}
							onchange={(e) =>
								store.updateSettings({ soundFeedbackVolume: parseInt(e.currentTarget.value) })}
						/>
					</fieldset>
				{/if}
			</div>
		</div>
	</div>

	<div class="card bg-base-200">
		<div class="card-body">
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
					<fieldset class="fieldset pl-6">
						<label class="label" for="unloadSeconds">
							<span class="label-text">Unload after (seconds)</span>
						</label>
						<input
							type="number"
							id="unloadSeconds"
							class="input input-sm w-32"
							min="60"
							max="3600"
							step="60"
							value={store.settings.modelUnloadSeconds}
							onblur={(e) =>
								store.updateSettings({ modelUnloadSeconds: parseInt(e.currentTarget.value) })}
						/>
					</fieldset>
				{/if}
			</div>
		</div>
	</div>

	<div class="card bg-base-200">
		<div class="card-body">
			<h3 class="card-title text-base">History</h3>
			<fieldset class="fieldset">
				<label class="label" for="historyLimit">
					<span class="label-text">Maximum entries to keep</span>
				</label>
				<input
					type="number"
					id="historyLimit"
					class="input input-sm w-32"
					min="10"
					max="1000"
					step="10"
					value={store.settings.historyLimit}
					onblur={(e) => store.updateSettings({ historyLimit: parseInt(e.currentTarget.value) })}
				/>
			</fieldset>
		</div>
	</div>
</div>
