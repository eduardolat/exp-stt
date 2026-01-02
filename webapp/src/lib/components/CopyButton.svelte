<script lang="ts">
	import { Copy } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';

	interface Props {
		text: string;
		showLabel?: boolean;
		label?: string;
		class?: string;
	}

	let { text, showLabel = false, label = 'Copy', class: className = '' }: Props = $props();

	async function handleCopy(event?: Event) {
		event?.stopPropagation();
		await navigator.clipboard.writeText(text);
		toast.success('Copied to clipboard');
	}
</script>

<button
	class="btn btn-ghost btn-xs {showLabel ? 'gap-1' : 'btn-square'} {className}"
	onclick={handleCopy}
	title="Copy to clipboard"
>
	<Copy class="size-3" />
	{#if showLabel && label}
		{label}
	{/if}
</button>
