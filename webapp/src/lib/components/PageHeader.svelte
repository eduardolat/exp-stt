<script lang="ts">
	import type { Component } from 'svelte';

	interface Props {
		icon: Component;
		title: string;
		description?: string;
		children?: import('svelte').Snippet;
		actions?: import('svelte').Snippet;
	}

	let { icon: Icon, title, description, children, actions }: Props = $props();
</script>

<div class="mb-4 flex items-end justify-between gap-4">
	<div class="flex items-center gap-3">
		<div class="flex size-10 items-center justify-center rounded-lg bg-primary/10">
			<Icon class="size-5 text-primary" />
		</div>
		<div class="flex h-12 flex-col justify-center">
			<h2 class="text-lg font-semibold">{title}</h2>
			{#if description}
				<p class="text-sm opacity-70">{description}</p>
			{/if}
		</div>
	</div>
	{#if actions}
		<div class="flex gap-2">
			{@render actions()}
		</div>
	{/if}
</div>

{#if children}
	{@render children()}
{/if}
