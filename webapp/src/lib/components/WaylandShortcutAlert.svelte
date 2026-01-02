<script lang="ts">
	import type { Snippet } from 'svelte';
	import { store } from '$lib/store.svelte';
	import { Info } from '@lucide/svelte';
	import CopyButton from './CopyButton.svelte';
	import Modal from './Modal.svelte';

	interface Props {
		children?: Snippet;
	}

	let { children }: Props = $props();

	let isWayland = $derived(store.systemInfo.displayServer.toLowerCase() === 'wayland');
	let waylandModal: { open: () => void; close: () => void } | undefined = $state();

	function waylandModalOpen(evt: Event) {
		evt.preventDefault();
		waylandModal?.open();
	}
</script>

{#if !isWayland && children}
	{@render children()}
{/if}

{#if isWayland}
	<div role="alert" class="alert alert-soft">
		<Info class="size-5" />
		<span>
			When you use Wayland, tribar can't set keyboard shortcuts for you, you should do it manually.
			<a class="link" href="#/" onclick={waylandModalOpen}>learn how here.</a>
		</span>
	</div>

	<Modal bind:this={waylandModal} title="Wayland Keyboard Shortcut Limitation" size="md">
		<div class="space-y-4">
			<p class="text-sm">
				<strong>Why global shortcuts don't work on Wayland:</strong>
			</p>
			<p class="text-sm opacity-80">
				Wayland's security model prevents applications from registering global keyboard shortcuts.
				This is a design decision to improve security and prevent keylogging.
			</p>

			<div class="divider my-2"></div>

			<p class="text-sm">
				<strong>Alternative solution:</strong>
			</p>
			<p class="text-sm opacity-80">
				You can configure a keyboard shortcut through your system settings instead. Add the
				following command to your desktop environment's keyboard shortcuts:
			</p>

			<div class="rounded-lg bg-base-300 p-3">
				<div class="mb-2 flex items-center justify-between">
					<code class="font-mono text-sm">tribar --toggle</code>
					<CopyButton text="tribar --toggle" />
				</div>
				<p class="text-xs opacity-70">
					This command toggles the recording on/off when triggered by your system shortcut.
				</p>
			</div>

			<p class="text-sm opacity-80">
				For example, on GNOME, go to
				<strong>Settings → Keyboard → Keyboard Shortcuts → Custom Shortcuts</strong>
				and add a custom shortcut with this command.
			</p>
		</div>
	</Modal>
{/if}
