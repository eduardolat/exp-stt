<script lang="ts">
	import { Info } from '@lucide/svelte';

	interface Props {
		nodeType: string;
	}

	let { nodeType }: Props = $props();

	const descriptions: Record<string, string> = {
		transcribe: 'Converts audio input to text using the Parakeet speech recognition model.',
		clipboard_copy: 'Copies the input text to the system clipboard.',
		clipboard_paste: 'Pastes the clipboard content using a keyboard shortcut.'
	};
</script>

<div class="flex flex-col gap-3">
	<div class="rounded-lg bg-base-200 p-4 text-center text-sm opacity-70">
		<Info class="mx-auto mb-2 size-8 opacity-50" />
		<p>This node has no configurable options.</p>
	</div>

	{#if descriptions[nodeType]}
		<div class="rounded-lg bg-info/10 p-3 text-xs text-info">
			{descriptions[nodeType]}
		</div>
	{/if}
</div>
