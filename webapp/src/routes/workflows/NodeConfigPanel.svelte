<script lang="ts">
	import { X } from '@lucide/svelte';
	import type { Node } from '@xyflow/svelte';

	// Individual node config components
	import JavaScriptConfig from './config/JavaScriptConfig.svelte';
	import HttpConfig from './config/HttpConfig.svelte';
	import TerminalConfig from './config/TerminalConfig.svelte';
	import DelayConfig from './config/DelayConfig.svelte';
	import NotifyConfig from './config/NotifyConfig.svelte';
	import SoundConfig from './config/SoundConfig.svelte';
	import ConditionConfig from './config/ConditionConfig.svelte';
	import AiProcessConfig from './config/AiProcessConfig.svelte';
	import EmptyConfig from './config/EmptyConfig.svelte';

	interface Props {
		node: Node;
		onUpdate: (nodeId: string, config: Record<string, unknown>) => void;
		onClose: () => void;
		disabled?: boolean;
	}

	let { node, onUpdate, onClose, disabled = false }: Props = $props();

	let config = $state<Record<string, unknown>>({
		...(node.data.config as Record<string, unknown>)
	});

	function handleChange(key: string, value: unknown) {
		config[key] = value;
		onUpdate(node.id, config);
	}

	const nodeType = $derived(node.data.nodeType as string);

	// Mapping of node types to their config components
	const configComponents: Record<string, typeof JavaScriptConfig> = {
		javascript: JavaScriptConfig,
		http: HttpConfig,
		terminal: TerminalConfig,
		delay: DelayConfig,
		notify: NotifyConfig,
		sound: SoundConfig,
		condition: ConditionConfig,
		ai_process: AiProcessConfig
	};

	const emptyConfigTypes = ['transcribe', 'clipboard_copy', 'clipboard_paste'];
</script>

<div class="flex h-full flex-col">
	<div class="flex items-center justify-between border-b border-base-300 p-3">
		<h3 class="text-sm font-medium">Configure Node</h3>
		<button class="btn btn-ghost btn-xs" onclick={onClose}>
			<X class="size-4" />
		</button>
	</div>

	<div class="flex-1 overflow-y-auto p-3">
		<div class="mb-3">
			<span class="badge badge-neutral">{node.data.label}</span>
		</div>

		{#if configComponents[nodeType]}
			<svelte:component
				this={configComponents[nodeType]}
				{config}
				onChange={handleChange}
				{disabled}
			/>
		{:else if emptyConfigTypes.includes(nodeType)}
			<EmptyConfig {nodeType} />
		{:else}
			<div class="rounded-lg bg-base-200 p-4 text-center text-sm opacity-70">
				Unknown node type: {nodeType}
			</div>
		{/if}
	</div>
</div>
