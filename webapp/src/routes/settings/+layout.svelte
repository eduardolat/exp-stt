<script lang="ts">
	import { Settings2, Sparkles } from '@lucide/svelte';
	import { page } from '$app/state';

	let { children } = $props();

	function isActive(path: string): boolean {
		const currentPath = page.url.hash.slice(1) || '/';
		return currentPath === path;
	}
</script>

<div class="space-y-4">
	<div role="tablist" class="tabs-box tabs">
		<a
			href="#/settings/general"
			role="tab"
			class="tab gap-2"
			class:tab-active={isActive('/settings/general') || isActive('/settings')}
		>
			<Settings2 class="size-4" />
			General
		</a>
		<a
			href="#/settings/ai"
			role="tab"
			class="tab gap-2"
			class:tab-active={isActive('/settings/ai')}
		>
			<Sparkles class="size-4" />
			AI Enhancement
		</a>
	</div>

	{@render children()}
</div>
