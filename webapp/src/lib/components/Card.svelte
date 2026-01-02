<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
		lighter?: boolean;
		darker?: boolean;
		interactive?: boolean;
		active?: boolean;
		onclick?: () => void;
		class?: string;
	}

	let {
		children,
		lighter = false,
		darker = false,
		interactive = false,
		active = false,
		onclick,
		class: className = ''
	}: Props = $props();

	let bgClass = $derived(lighter ? 'bg-base-100' : darker ? 'bg-base-300' : 'bg-base-200');

	let interactiveClasses = $derived(
		interactive
			? active
				? 'border-primary bg-primary/10 shadow-md cursor-pointer'
				: 'cursor-pointer hover:border-primary hover:bg-primary/10 hover:shadow-md'
			: ''
	);

	let combinedClasses = $derived(
		`card w-full text-left border border-base-300 shadow-sm transition-all ${bgClass} ${interactiveClasses} ${className}`.trim()
	);
</script>

{#if interactive}
	<button class={combinedClasses} {onclick} type="button">
		{@render children()}
	</button>
{:else}
	<div class={combinedClasses}>
		{@render children()}
	</div>
{/if}
