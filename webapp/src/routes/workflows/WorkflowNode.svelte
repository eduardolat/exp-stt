<script lang="ts">
	import { Handle, Position } from '@xyflow/svelte';
	import type { NodeProps } from '@xyflow/svelte';

	interface Props extends NodeProps {
		data: {
			label: string;
			nodeType: string;
			config?: Record<string, unknown>;
		};
	}

	let { data, selected }: Props = $props();

	const isTrigger = $derived(data.nodeType === 'transcribe');

	const typeStyles: Record<string, string> = {
		transcribe: 'border-primary bg-primary/10',
		ai_process: 'border-secondary bg-secondary/10',
		clipboard_copy: 'border-info bg-info/10',
		clipboard_paste: 'border-info bg-info/10',
		notify: 'border-warning bg-warning/10',
		sound: 'border-accent bg-accent/10',
		condition: 'border-error bg-error/10',
		delay: 'border-neutral bg-neutral/10',
		javascript: 'border-success bg-success/10',
		terminal: 'border-neutral bg-neutral/10',
		http: 'border-info bg-info/10'
	};

	const nodeStyle = $derived(typeStyles[data.nodeType] || 'border-base-300 bg-base-200');
</script>

<div
	class="workflow-node rounded-lg border-2 px-4 py-3 shadow-md transition-all {nodeStyle}"
	class:ring-2={selected}
	class:ring-primary={selected}
	class:ring-offset-2={selected}
	class:ring-offset-base-100={selected}
>
	<!-- Input Handle (left) - Not for trigger nodes -->
	{#if !isTrigger}
		<Handle type="target" position={Position.Left} class="!size-3 !bg-base-content" />
	{/if}

	<!-- Node Content -->
	<div class="flex flex-col items-center gap-1">
		{#if isTrigger}
			<span class="text-[10px] font-semibold tracking-wide text-primary uppercase opacity-70">
				Trigger
			</span>
		{/if}
		<span class="text-sm font-medium whitespace-nowrap">{data.label}</span>
	</div>

	<!-- Output Handle (right) -->
	<Handle type="source" position={Position.Right} class="!size-3 !bg-base-content" />
</div>

<style>
	.workflow-node {
		min-width: 120px;
	}
</style>
