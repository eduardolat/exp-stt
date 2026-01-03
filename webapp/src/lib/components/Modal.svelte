<script lang="ts">
	import type { Snippet } from 'svelte';
	import { X } from '@lucide/svelte';

	interface Props {
		title: string;
		size?: 'sm' | 'md' | 'lg';
		disableBackdropClose?: boolean;
		children?: Snippet;
		actions?: Snippet;
	}

	let { title, size = 'md', disableBackdropClose = false, children, actions }: Props = $props();

	let dialog: HTMLDialogElement | undefined = $state();

	const sizeClasses = {
		sm: 'max-w-md',
		md: 'max-w-xl',
		lg: 'max-w-3xl'
	};

	export function open() {
		dialog?.showModal();
	}

	export function close() {
		dialog?.close();
	}
</script>

<dialog bind:this={dialog} class="modal">
	<div class="modal-box flex max-h-[90dvh] flex-col p-0 {sizeClasses[size]}">
		<!-- Header -->
		<div class="flex items-center justify-between border-b border-base-300 p-4">
			<h3 class="text-lg font-bold">{title}</h3>
			<button class="btn btn-circle btn-ghost btn-sm" onclick={close} title="Close">
				<X class="size-4" />
			</button>
		</div>

		<!-- Scrollable Content -->
		<div class="grow overflow-y-auto p-4">
			{#if children}
				{@render children()}
			{/if}
		</div>

		<!-- Footer -->
		{#if actions}
			<div class="border-t border-base-300 p-4">
				<div class="m-0 modal-action">
					{@render actions()}
				</div>
			</div>
		{/if}
	</div>

	{#if !disableBackdropClose}
		<form method="dialog" class="modal-backdrop">
			<button>close</button>
		</form>
	{:else}
		<div class="modal-backdrop"></div>
	{/if}
</dialog>
