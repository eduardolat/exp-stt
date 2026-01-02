<script lang="ts">
	import type { Snippet } from 'svelte';
	import { store } from '$lib/store.svelte';
	import { Info } from '@lucide/svelte';
	import CopyButton from './CopyButton.svelte';
	import Modal from './Modal.svelte';
	import Card from './Card.svelte';

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
			Tribar cannot register global keyboard shortcuts automatically on Wayland.
			<a class="link" href="#/" onclick={waylandModalOpen}>Learn how to configure them manually.</a>
		</span>
	</div>

	<Modal bind:this={waylandModal} title="Configuring Keyboard Shortcuts on Wayland" size="md">
		<div class="space-y-4">
			<p class="text-sm">
				<strong>Why can't Tribar register shortcuts automatically?</strong>
			</p>
			<p class="text-sm opacity-80">
				Unlike X11, Wayland implements a stricter security model where applications cannot register
				global keyboard shortcuts directly. This is by design to prevent malicious applications from
				capturing keystrokes and potentially logging sensitive information like passwords.
			</p>
			<p class="text-sm opacity-80">
				While this enhances security, it means that applications like Tribar must rely on your
				desktop environment to handle global shortcuts.
			</p>

			<p class="text-sm">
				<strong>How to configure shortcuts manually:</strong>
			</p>
			<p class="text-sm opacity-80">
				You need to create a custom keyboard shortcut in your system settings that executes the
				following command:
			</p>

			<Card class="p-3">
				<div class="mb-2 flex items-center justify-between">
					<code class="font-mono text-sm">tribar --toggle</code>
					<CopyButton text="tribar --toggle" />
				</div>
				<p class="text-xs opacity-70">
					This command starts or stops recording when triggered by your chosen keyboard shortcut.
				</p>
			</Card>

			<p class="text-sm">
				<strong>Instructions by desktop environment:</strong>
			</p>

			<div class="space-y-3">
				<Card class="p-3">
					<p class="mb-1 text-sm font-semibold">GNOME / Ubuntu</p>
					<p class="text-xs opacity-80">
						Open <strong>Settings → Keyboard → Keyboard Shortcuts → Custom Shortcuts</strong>. Click
						the <strong>+</strong> button, name it "Toggle Tribar", paste the command above, and
						assign your preferred key combination (e.g., <kbd>Super+Space</kbd>).
					</p>
				</Card>

				<Card class="p-3">
					<p class="mb-1 text-sm font-semibold">KDE Plasma</p>
					<p class="text-xs opacity-80">
						Go to <strong>System Settings → Shortcuts → Custom Shortcuts</strong>. Create a new
						shortcut, select <strong>Command/URL</strong>, paste the command, and set your desired
						key combination.
					</p>
				</Card>

				<Card class="p-3">
					<p class="mb-1 text-sm font-semibold">Hyprland / Sway (Tiling WMs)</p>
					<p class="text-xs opacity-80">
						Add a bind to your config file:
						<code class="ml-1 text-xs">bind = SUPER, SPACE, exec, tribar --toggle</code>
					</p>
				</Card>
			</div>
		</div>
	</Modal>
{/if}
