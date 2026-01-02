<script lang="ts">
	import { store } from '$lib/store.svelte';
	import type { Shortcut } from '$lib/client.gen';

	const AVAILABLE_MODIFIERS = ['ctrl', 'alt', 'shift', 'meta'] as const;
	const AVAILABLE_KEYS = [
		'space',
		'a',
		'b',
		'c',
		'd',
		'e',
		'f',
		'g',
		'h',
		'i',
		'j',
		'k',
		'l',
		'm',
		'n',
		'o',
		'p',
		'q',
		'r',
		's',
		't',
		'u',
		'v',
		'w',
		'x',
		'y',
		'z',
		'0',
		'1',
		'2',
		'3',
		'4',
		'5',
		'6',
		'7',
		'8',
		'9',
		'f1',
		'f2',
		'f3',
		'f4',
		'f5',
		'f6',
		'f7',
		'f8',
		'f9',
		'f10',
		'f11',
		'f12'
	] as const;

	let modifiers = $state<string[]>([...store.shortcut.modifiers]);
	let key = $state(store.shortcut.key);
	let isUpdating = $state(false);
	let updateError = $state<string | null>(null);

	$effect(() => {
		modifiers = [...store.shortcut.modifiers];
		key = store.shortcut.key;
	});

	function toggleModifier(modifier: string) {
		if (modifiers.includes(modifier)) {
			modifiers = modifiers.filter((m) => m !== modifier);
		} else {
			modifiers = [...modifiers, modifier];
		}
	}

	async function saveShortcut() {
		if (!key) {
			updateError = 'Please select a key';
			return;
		}

		isUpdating = true;
		updateError = null;

		const shortcut: Shortcut = {
			modifiers: [...modifiers],
			key
		};

		try {
			await store.updateShortcut(shortcut);
		} catch (err) {
			updateError = err instanceof Error ? err.message : 'Failed to update shortcut';
		} finally {
			isUpdating = false;
		}
	}

	function formatModifierLabel(modifier: string): string {
		const labels: Record<string, string> = {
			ctrl: store.systemInfo.os === 'darwin' ? '⌃ Control' : 'Ctrl',
			alt: store.systemInfo.os === 'darwin' ? '⌥ Option' : 'Alt',
			shift: store.systemInfo.os === 'darwin' ? '⇧ Shift' : 'Shift',
			meta: store.systemInfo.os === 'darwin' ? '⌘ Command' : 'Win'
		};
		return labels[modifier] ?? modifier;
	}
</script>

<div class="space-y-4">
	<div>
		<span class="mb-2 block text-sm font-medium">Modifiers</span>
		<div class="flex flex-wrap gap-2">
			{#each AVAILABLE_MODIFIERS as modifier}
				<button
					type="button"
					class="btn btn-sm"
					class:btn-primary={modifiers.includes(modifier)}
					class:btn-outline={!modifiers.includes(modifier)}
					onclick={() => toggleModifier(modifier)}
				>
					{formatModifierLabel(modifier)}
				</button>
			{/each}
		</div>
	</div>

	<fieldset class="fieldset">
		<label class="label" for="shortcut-key">
			<span class="label-text">Key</span>
		</label>
		<select id="shortcut-key" class="select w-40 select-sm" bind:value={key}>
			<option value="">Select key...</option>
			{#each AVAILABLE_KEYS as k}
				<option value={k}>{k.toUpperCase()}</option>
			{/each}
		</select>
	</fieldset>

	{#if key}
		<div class="flex items-center gap-2 text-sm">
			<span class="opacity-70">Current shortcut:</span>
			{#if modifiers.length > 0}
				{#each modifiers as mod}
					<kbd class="kbd kbd-sm">{mod}</kbd>
					<span>+</span>
				{/each}
			{/if}
			<kbd class="kbd kbd-sm">{key}</kbd>
		</div>
	{/if}

	{#if updateError}
		<p class="text-sm text-error">{updateError}</p>
	{/if}

	<button class="btn btn-sm btn-primary" disabled={isUpdating || !key} onclick={saveShortcut}>
		{isUpdating ? 'Saving...' : 'Save Shortcut'}
	</button>
</div>
