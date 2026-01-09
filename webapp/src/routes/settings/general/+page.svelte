<script lang="ts">
	import { Settings } from '@lucide/svelte';
	import { PageHeader } from '$lib/components';
	import KeyboardShortcut from './KeyboardShortcut.svelte';
	import ModelUnloading from './ModelUnloading.svelte';
	import SoundFeedback from './SoundFeedback.svelte';
	import InputDevice from './InputDevice.svelte';
	import OutputMode from './OutputMode.svelte';
	import Notifications from './Notifications.svelte';
	import AdvancedMode from './AdvancedMode.svelte';
</script>

<svelte:head>
	<title>General Settings - Tribar Voice</title>
</svelte:head>

<PageHeader
	icon={Settings}
	title="General Settings"
	description="Configure general Tribar settings"
/>

<div class="columns-2 gap-4 space-y-4">
	<InputDevice />
	<OutputMode />
	<KeyboardShortcut />
	<SoundFeedback />
	<Notifications />
	<ModelUnloading />
	<AdvancedMode />
</div>
