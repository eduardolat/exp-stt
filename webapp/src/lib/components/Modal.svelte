<script lang="ts">
	import type { Snippet } from 'svelte';
	import { X } from '@lucide/svelte';

	interface Props {
		title: string;
		size?: 'sm' | 'md' | 'lg';
		children?: Snippet;
		actions?: Snippet;
	}

	let { title, size = 'md', children, actions }: Props = $props();

	let dialog: HTMLDialogElement | undefined = $state();

	const sizeClasses = {
		sm: 'modal-box max-w-md',
		md: 'modal-box',
		lg: 'modal-box max-w-3xl'
	};

	export function open() {
		dialog?.showModal();
	}

	export function close() {
		dialog?.close();
	}
</script>

<dialog bind:this={dialog} class="modal">
	<div class={sizeClasses[size]}>
		<div class="mb-4 flex items-center justify-between">
			<h3 class="text-lg font-bold">{title}</h3>
			<button class="btn btn-circle btn-ghost btn-sm" onclick={close} title="Close">
				<X class="size-4" />
			</button>
		</div>

		<div class="space-y-4">
			{#if children}
				{@render children()}
			{/if}
		</div>

		{#if actions}
			<div class="modal-action">
				{@render actions()}
			</div>
		{/if}
	</div>

	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>
