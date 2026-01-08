<script lang="ts">
	import { Sparkles } from '@lucide/svelte';
	import type { HistoryEntry } from '$lib/client.gen';
	import { Card, CopyButton } from '$lib/components';
	import { formatRelativeTime } from '$lib/utils/formatRelativeTime';

	interface Props {
		entry: HistoryEntry;
		onSelect: (entry: HistoryEntry) => void;
		darker?: boolean;
	}

	let { entry, onSelect, darker }: Props = $props();

	function formatTimestamp(iso: string): string {
		const text = formatRelativeTime(iso);
		const capitalized = text.charAt(0).toUpperCase() + text.slice(1);
		return capitalized;
	}
</script>

{#snippet textWithFallback(text: string)}
	{#if text.trim()}
		{text}
	{:else}
		<span class="opacity-70">(No text available)</span>
	{/if}
{/snippet}

<Card {darker} interactive onclick={() => onSelect(entry)} class="p-3">
	<div class="mb-1.5 flex items-center justify-between gap-2">
		<div class="flex items-center gap-2 text-xs opacity-70">
			<span>{formatTimestamp(entry.timestamp)}</span>
			{#if entry.postProcessed}
				<span class="badge gap-1 badge-xs badge-secondary">
					<Sparkles class="size-2.5" />
					AI
				</span>
			{/if}
		</div>
		<CopyButton text={entry.textFinal} showLabel />
	</div>
	<p class="line-clamp-2 text-sm">
		{@render textWithFallback(entry.textFinal)}
	</p>
</Card>
