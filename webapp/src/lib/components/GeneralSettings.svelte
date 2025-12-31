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
	import ShortcutInput from './ShortcutInput.svelte';
</script>

<div class="space-y-6">
	<!-- Input Device -->
	<section class="card p-4">
		<h3 class="mb-3 flex items-center gap-2 font-medium">
			<Mic class="size-4" />
			Input Device
		</h3>
		<div class="field">
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
		</div>
	</section>

	<!-- Output Mode -->
	<section class="card p-4">
		<h3 class="mb-3 flex items-center gap-2 font-medium">
			<Clipboard class="size-4" />
			Output Mode
		</h3>
		<div class="space-y-3">
			<label
				class="flex cursor-pointer items-start gap-3 rounded-lg border p-3 transition-colors hover:bg-muted/50"
			>
				<input
					type="radio"
					name="outputMode"
					value="copy_only"
					checked={store.settings.outputMode === 'copy_only'}
					onchange={() => store.updateSettings({ outputMode: 'copy_only' })}
					class="mt-0.5"
				/>
				<div class="flex-1">
					<div class="flex items-center gap-2">
						<Clipboard class="size-4" />
						<span class="font-medium">Copy Only</span>
					</div>
					<p class="text-xs text-muted-foreground">
						Copies transcription to clipboard without pasting
					</p>
				</div>
			</label>

			<label
				class="flex cursor-pointer items-start gap-3 rounded-lg border p-3 transition-colors hover:bg-muted/50"
			>
				<input
					type="radio"
					name="outputMode"
					value="copy_paste"
					checked={store.settings.outputMode === 'copy_paste'}
					onchange={() => store.updateSettings({ outputMode: 'copy_paste' })}
					class="mt-0.5"
				/>
				<div class="flex-1">
					<div class="flex items-center gap-2">
						<ClipboardPaste class="size-4" />
						<span class="font-medium">Copy & Paste</span>
					</div>
					<p class="text-xs text-muted-foreground">
						Copies to clipboard and pastes at cursor position
					</p>
				</div>
			</label>

			<label
				class="flex cursor-pointer items-start gap-3 rounded-lg border p-3 transition-colors hover:bg-muted/50"
			>
				<input
					type="radio"
					name="outputMode"
					value="ghost_paste"
					checked={store.settings.outputMode === 'ghost_paste'}
					onchange={() => store.updateSettings({ outputMode: 'ghost_paste' })}
					class="mt-0.5"
				/>
				<div class="flex-1">
					<div class="flex items-center gap-2">
						<Ghost class="size-4" />
						<span class="font-medium">Ghost Paste</span>
					</div>
					<p class="text-xs text-muted-foreground">
						Pastes without modifying your clipboard contents
					</p>
				</div>
			</label>

			<label class="flex cursor-pointer items-center gap-3 pt-2">
				<input
					type="checkbox"
					role="switch"
					checked={store.settings.outputTrailingSpace}
					onchange={(e) => store.updateSettings({ outputTrailingSpace: e.currentTarget.checked })}
				/>
				<span class="text-sm">Add trailing space after paste</span>
			</label>
		</div>
	</section>

	<!-- Keyboard Shortcut -->
	<section class="card p-4">
		<h3 class="mb-3 flex items-center gap-2 font-medium">
			<Keyboard class="size-4" />
			Keyboard Shortcut
		</h3>
		<ShortcutInput />
	</section>

	<!-- Notifications -->
	<section class="card p-4">
		<h3 class="mb-3 flex items-center gap-2 font-medium">
			<Bell class="size-4" />
			Notifications
		</h3>
		<div class="space-y-3">
			<label class="flex cursor-pointer items-center justify-between gap-3">
				<div class="flex items-center gap-2">
					<BellOff class="size-4 text-muted-foreground" />
					<span class="text-sm">Notify on errors</span>
				</div>
				<input
					type="checkbox"
					role="switch"
					checked={store.settings.notifyOnError}
					onchange={(e) => store.updateSettings({ notifyOnError: e.currentTarget.checked })}
				/>
			</label>
			<label class="flex cursor-pointer items-center justify-between gap-3">
				<div class="flex items-center gap-2">
					<Bell class="size-4 text-muted-foreground" />
					<span class="text-sm">Notify on recording start</span>
				</div>
				<input
					type="checkbox"
					role="switch"
					checked={store.settings.notifyOnStart}
					onchange={(e) => store.updateSettings({ notifyOnStart: e.currentTarget.checked })}
				/>
			</label>
			<label class="flex cursor-pointer items-center justify-between gap-3">
				<div class="flex items-center gap-2">
					<Bell class="size-4 text-muted-foreground" />
					<span class="text-sm">Notify on transcription complete</span>
				</div>
				<input
					type="checkbox"
					role="switch"
					checked={store.settings.notifyOnFinish}
					onchange={(e) => store.updateSettings({ notifyOnFinish: e.currentTarget.checked })}
				/>
			</label>
		</div>
	</section>

	<!-- Sound Feedback -->
	<section class="card p-4">
		<h3 class="mb-3 flex items-center gap-2 font-medium">
			<Volume2 class="size-4" />
			Sound Feedback
		</h3>
		<div class="space-y-4">
			<label class="flex cursor-pointer items-center justify-between gap-3">
				<div class="flex items-center gap-2">
					{#if store.settings.soundFeedbackEnable}
						<Volume2 class="size-4 text-muted-foreground" />
					{:else}
						<VolumeX class="size-4 text-muted-foreground" />
					{/if}
					<span class="text-sm">Enable sound feedback</span>
				</div>
				<input
					type="checkbox"
					role="switch"
					checked={store.settings.soundFeedbackEnable}
					onchange={(e) => store.updateSettings({ soundFeedbackEnable: e.currentTarget.checked })}
				/>
			</label>

			{#if store.settings.soundFeedbackEnable}
				<div class="space-y-3 pl-6">
					<div class="field">
						<label for="volume" class="text-sm">Volume: {store.settings.soundFeedbackVolume}%</label
						>
						<input
							type="range"
							id="volume"
							min="0"
							max="100"
							value={store.settings.soundFeedbackVolume}
							onchange={(e) =>
								store.updateSettings({ soundFeedbackVolume: parseInt(e.currentTarget.value) })}
							class="w-full"
						/>
					</div>
				</div>
			{/if}
		</div>
	</section>

	<!-- Model Management -->
	<section class="card p-4">
		<h3 class="mb-3 flex items-center gap-2 font-medium">
			<Timer class="size-4" />
			Model Unloading
		</h3>
		<div class="space-y-4">
			<label class="flex cursor-pointer items-center justify-between gap-3">
				<div>
					<span class="text-sm">Auto-unload model when idle</span>
					<p class="text-xs text-muted-foreground">Reduces memory usage when not in use</p>
				</div>
				<input
					type="checkbox"
					role="switch"
					checked={store.settings.modelUnloadEnable}
					onchange={(e) => store.updateSettings({ modelUnloadEnable: e.currentTarget.checked })}
				/>
			</label>

			{#if store.settings.modelUnloadEnable}
				<div class="field pl-6">
					<label for="unloadSeconds" class="text-sm">Unload after (seconds)</label>
					<input
						type="number"
						id="unloadSeconds"
						min="60"
						max="3600"
						step="60"
						value={store.settings.modelUnloadSeconds}
						onchange={(e) =>
							store.updateSettings({ modelUnloadSeconds: parseInt(e.currentTarget.value) })}
						class="input w-32"
					/>
				</div>
			{/if}
		</div>
	</section>

	<!-- History Settings -->
	<section class="card p-4">
		<h3 class="mb-3 font-medium">History</h3>
		<div class="field">
			<label for="historyLimit" class="text-sm">Maximum entries to keep</label>
			<input
				type="number"
				id="historyLimit"
				min="10"
				max="1000"
				step="10"
				value={store.settings.historyLimit}
				onchange={(e) => store.updateSettings({ historyLimit: parseInt(e.currentTarget.value) })}
				class="input w-32"
			/>
		</div>
	</section>
</div>
