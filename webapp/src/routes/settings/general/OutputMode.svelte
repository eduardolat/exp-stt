<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Clipboard, ClipboardPaste, Ghost } from '@lucide/svelte';
	import { Card } from '$lib/components';
</script>

<Card class="card-body break-inside-avoid">
	<h3 class="card-title text-base">
		<Clipboard class="size-4" />
		Output Mode
	</h3>
	<div class="space-y-2">
		<Card
			interactive
			darker
			active={store.settings.outputMode === 'copy_only'}
			onclick={() => store.updateSettings({ outputMode: 'copy_only' })}
			class="p-3"
		>
			<div class="flex items-start gap-3">
				<input
					type="radio"
					name="outputMode"
					class="pointer-events-none radio mt-0.5 radio-sm"
					value="copy_only"
					checked={store.settings.outputMode === 'copy_only'}
					tabindex="-1"
				/>
				<div class="flex-1">
					<div class="flex items-center gap-2 font-medium">
						<Clipboard class="size-4" />
						Copy Only
					</div>
					<p class="text-xs opacity-70">Copies transcription to clipboard without pasting</p>
				</div>
			</div>
		</Card>

		<Card
			interactive
			darker
			active={store.settings.outputMode === 'copy_paste'}
			onclick={() => store.updateSettings({ outputMode: 'copy_paste' })}
			class="p-3"
		>
			<div class="flex items-start gap-3">
				<input
					type="radio"
					name="outputMode"
					class="pointer-events-none radio mt-0.5 radio-sm"
					value="copy_paste"
					checked={store.settings.outputMode === 'copy_paste'}
					tabindex="-1"
				/>
				<div class="flex-1">
					<div class="flex items-center gap-2 font-medium">
						<ClipboardPaste class="size-4" />
						Copy & Paste
					</div>
					<p class="text-xs opacity-70">Copies to clipboard and pastes at cursor position</p>
				</div>
			</div>
		</Card>

		<Card
			interactive
			darker
			active={store.settings.outputMode === 'ghost_paste'}
			onclick={() => store.updateSettings({ outputMode: 'ghost_paste' })}
			class="p-3"
		>
			<div class="flex items-start gap-3">
				<input
					type="radio"
					name="outputMode"
					class="pointer-events-none radio mt-0.5 radio-sm"
					value="ghost_paste"
					checked={store.settings.outputMode === 'ghost_paste'}
					tabindex="-1"
				/>
				<div class="flex-1">
					<div class="flex items-center gap-2 font-medium">
						<Ghost class="size-4" />
						Ghost Paste
					</div>
					<p class="text-xs opacity-70">Pastes without modifying your clipboard contents</p>
				</div>
			</div>
		</Card>

		{#if store.settings.outputMode !== 'copy_only'}
			<fieldset class="fieldset pt-2">
				<label class="label" for="pasteShortcut">
					<span class="label-text">Paste Shortcut Sequence</span>
				</label>
				<select
					id="pasteShortcut"
					class="select w-full select-sm"
					value={store.settings.pasteShortcut}
					onchange={(e) => store.updateSettings({ pasteShortcut: e.currentTarget.value })}
				>
					<option value="ctrl+v">Ctrl + V</option>
					<option value="ctrl+shift+v">Ctrl + Shift + V</option>
					<option value="shift+insert">Shift + Insert</option>
				</select>
			</fieldset>

			<label class="label cursor-pointer justify-between pt-2">
				<span class="text-sm">Add trailing space after paste</span>
				<input
					type="checkbox"
					class="toggle toggle-sm"
					checked={store.settings.outputTrailingSpace}
					onchange={(e) => store.updateSettings({ outputTrailingSpace: e.currentTarget.checked })}
				/>
			</label>
		{/if}
	</div>
</Card>
