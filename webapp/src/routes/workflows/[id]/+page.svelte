<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { SvelteFlow, Controls, Background, MiniMap, type Node, type Edge } from '@xyflow/svelte';
	import '@xyflow/svelte/dist/style.css';
	import { ArrowLeft, Save, Workflow, Trash2 } from '@lucide/svelte';
	import { Card } from '$lib/components';
	import type { Workflow as WorkflowType, WorkflowNode, WorkflowEdge } from '$lib/client.gen';
	import { toast } from 'svelte-sonner';
	import NodeConfigPanel from '../NodeConfigPanel.svelte';
	import WorkflowNodeComponent from '../WorkflowNode.svelte';

	// Custom node types for SvelteFlow
	const nodeTypes = {
		workflow: WorkflowNodeComponent
	};

	let workflow: WorkflowType | null = $state(null);
	let isLoading = $state(true);
	let isSaving = $state(false);
	let workflowName = $state('');
	let workflowDescription = $state('');

	// Svelte Flow state
	let nodes: Node[] = $state([]);
	let edges: Edge[] = $state([]);

	// Selected node for config panel
	let selectedNode: Node | null = $state(null);

	const workflowId = $derived(page.params.id);

	onMount(async () => {
		await loadWorkflow();
		await store.fetchAvailableNodes();
	});

	async function loadWorkflow() {
		isLoading = true;
		try {
			if (!workflowId) return;
			workflow = await store.getWorkflow(workflowId);
			if (workflow) {
				workflowName = workflow.name;
				workflowDescription = workflow.description ?? '';
				nodes = convertToFlowNodes(workflow.nodes);
				edges = convertToFlowEdges(workflow.edges);
			}
		} finally {
			isLoading = false;
		}
	}

	function convertToFlowNodes(workflowNodes: WorkflowNode[]): Node[] {
		return workflowNodes.map((node) => ({
			id: node.id,
			type: 'workflow',
			position: { x: node.position.x, y: node.position.y },
			data: {
				label: getNodeLabel(node.nodeType),
				nodeType: node.nodeType,
				config: JSON.parse(node.config || '{}')
			}
		}));
	}

	function convertToFlowEdges(workflowEdges: WorkflowEdge[]): Edge[] {
		return workflowEdges.map((edge) => ({
			id: edge.id,
			source: edge.source,
			target: edge.target,
			sourceHandle: edge.sourceHandle ?? undefined,
			targetHandle: edge.targetHandle ?? undefined
		}));
	}

	function getNodeLabel(nodeType: string): string {
		const labels: Record<string, string> = {
			transcribe: '🎤 Transcribe',
			ai_process: '✨ AI Process',
			clipboard_copy: '📋 Copy',
			clipboard_paste: '📋 Paste',
			notify: '🔔 Notify',
			sound: '🔊 Sound',
			condition: '🔀 Condition',
			delay: '⏱️ Delay',
			javascript: '📜 JavaScript',
			terminal: '💻 Terminal',
			http: '🌐 HTTP'
		};
		return labels[nodeType] || nodeType;
	}

	function convertFromFlowNodes(flowNodes: Node[]): WorkflowNode[] {
		return flowNodes.map((node) => ({
			id: node.id,
			nodeType: node.data.nodeType as string,
			position: { x: node.position.x, y: node.position.y },
			config: JSON.stringify(node.data.config || {}),
			data: '{}'
		}));
	}

	function convertFromFlowEdges(flowEdges: Edge[]): WorkflowEdge[] {
		return flowEdges.map((edge) => ({
			id: edge.id,
			source: edge.source,
			target: edge.target,
			sourceHandle: edge.sourceHandle ?? undefined,
			targetHandle: edge.targetHandle ?? undefined
		}));
	}

	async function saveWorkflow() {
		if (!workflow) return;
		isSaving = true;
		try {
			const updated = await store.updateWorkflow(workflow.id, {
				name: workflowName,
				description: workflowDescription || undefined,
				enabled: workflow.enabled,
				trigger: workflow.trigger,
				nodes: convertFromFlowNodes(nodes),
				edges: convertFromFlowEdges(edges)
			});
			if (updated) {
				toast.success('Workflow saved');
			}
		} catch {
			toast.error('Failed to save workflow');
		} finally {
			isSaving = false;
		}
	}

	function addNode(nodeType: string) {
		const newId = crypto.randomUUID();
		const newNode: Node = {
			id: newId,
			type: 'workflow',
			position: { x: 250, y: nodes.length * 100 + 50 },
			data: {
				label: getNodeLabel(nodeType),
				nodeType: nodeType,
				config: {}
			}
		};
		nodes = [...nodes, newNode];
		selectedNode = newNode;
	}

	function isDefault(): boolean {
		return workflow?.id === 'default';
	}

	// Check if selected node is the Transcribe entry node (not deletable)
	function isTranscribeNode(node: Node | null): boolean {
		return node?.data.nodeType === 'transcribe';
	}

	// Filter available nodes to exclude transcribe (it's always present)
	const filteredAvailableNodes = $derived(
		store.availableNodes.filter((n) => n.nodeType !== 'transcribe')
	);

	function handleNodeClick(event: { node: Node }) {
		const clickedNode = nodes.find((n) => n.id === event.node.id);
		selectedNode = clickedNode ?? null;
	}

	function handleNodeConfigUpdate(nodeId: string, config: Record<string, unknown>) {
		nodes = nodes.map((node) => {
			if (node.id === nodeId) {
				return {
					...node,
					data: {
						...node.data,
						config
					}
				};
			}
			return node;
		});
		// Update selected node reference
		if (selectedNode?.id === nodeId) {
			selectedNode = nodes.find((n) => n.id === nodeId) ?? null;
		}
	}

	function deleteSelectedNode() {
		if (!selectedNode) return;
		const nodeId = selectedNode.id;
		nodes = nodes.filter((n) => n.id !== nodeId);
		edges = edges.filter((e) => e.source !== nodeId && e.target !== nodeId);
		selectedNode = null;
	}

	function handlePaneClick() {
		selectedNode = null;
	}

	// Handle edge deletion with Backspace/Delete key when edges are selected
	function handleEdgeClick(event: { edge: Edge }) {
		if (isDefault()) return;
		// Remove the clicked edge
		edges = edges.filter((e) => e.id !== event.edge.id);
	}
</script>

<div class="flex h-full flex-col">
	<!-- Header -->
	<div class="flex shrink-0 items-center justify-between p-4">
		<div class="flex items-center gap-4">
			<a href="#/workflows" class="btn btn-ghost btn-sm">
				<ArrowLeft class="size-4" />
			</a>
			<div>
				<input
					type="text"
					bind:value={workflowName}
					class="input input-ghost text-lg font-semibold"
					placeholder="Workflow name"
					disabled={isDefault()}
				/>
			</div>
		</div>
		<div class="flex gap-2">
			<button
				class="btn gap-2 btn-sm btn-primary"
				onclick={saveWorkflow}
				disabled={isSaving || isDefault()}
			>
				<Save class="size-4" />
				{isSaving ? 'Saving...' : 'Save'}
			</button>
		</div>
	</div>

	{#if isLoading}
		<div class="flex flex-1 items-center justify-center">
			<span class="loading loading-lg loading-spinner"></span>
		</div>
	{:else if !workflow}
		<Card class="m-4 flex flex-1 flex-col items-center justify-center">
			<Workflow class="size-16 opacity-30" />
			<p class="mt-4 text-lg">Workflow not found</p>
			<a href="#/workflows" class="btn mt-4 btn-sm btn-primary">Back to Workflows</a>
		</Card>
	{:else}
		<!-- Main Editor Area - Takes remaining height -->
		<div class="flex min-h-0 flex-1 gap-4 px-4 pb-4">
			<!-- Node Palette -->
			<Card class="w-48 shrink-0 overflow-y-auto p-3">
				<h3 class="mb-3 text-sm font-medium">Add Node</h3>
				<div class="flex flex-col gap-1">
					{#each filteredAvailableNodes as nodeDef (nodeDef.nodeType)}
						<button
							class="btn justify-start text-left btn-ghost btn-sm"
							onclick={() => addNode(nodeDef.nodeType)}
							disabled={isDefault()}
						>
							{getNodeLabel(nodeDef.nodeType)}
						</button>
					{/each}
				</div>
			</Card>

			<!-- Flow Editor Canvas - MUST have explicit height -->
			<div class="svelte-flow-wrapper flex-1 overflow-hidden rounded-lg border border-base-300">
				<SvelteFlow
					bind:nodes
					bind:edges
					{nodeTypes}
					fitView
					onnodeclick={handleNodeClick}
					onpaneclick={handlePaneClick}
					onedgeclick={handleEdgeClick}
				>
					<Controls />
					<Background />
					<MiniMap />
				</SvelteFlow>
			</div>

			<!-- Config Panel -->
			{#if selectedNode}
				<Card class="w-72 shrink-0 overflow-hidden p-0">
					<NodeConfigPanel
						node={selectedNode}
						onUpdate={handleNodeConfigUpdate}
						onClose={() => (selectedNode = null)}
						disabled={isDefault()}
					/>
					{#if !isDefault() && !isTranscribeNode(selectedNode)}
						<div class="border-t border-base-300 p-3">
							<button class="btn w-full gap-2 btn-sm btn-error" onclick={deleteSelectedNode}>
								<Trash2 class="size-4" />
								Delete Node
							</button>
						</div>
					{/if}
				</Card>
			{/if}
		</div>
	{/if}
</div>

<style>
	/* Svelte Flow requires the wrapper to have explicit dimensions */
	.svelte-flow-wrapper {
		height: 100%;
		width: 100%;
	}

	:global(.svelte-flow) {
		height: 100% !important;
		width: 100% !important;
		--xy-background-color: oklch(var(--b1));
		--xy-minimap-background-color: oklch(var(--b2));
		--xy-node-background-color: oklch(var(--b2));
		--xy-node-border-color: oklch(var(--b3));
		--xy-edge-stroke: oklch(var(--p));
		--xy-handle-background-color: oklch(var(--p));
	}

	/* Ensure nodes have proper styling */
	:global(.svelte-flow__node) {
		border-radius: 8px;
		font-size: 12px;
		padding: 10px 14px;
		min-width: 120px;
	}

	:global(.svelte-flow__node.selected) {
		box-shadow: 0 0 0 2px oklch(var(--p));
	}

	:global(.svelte-flow__edge-path) {
		stroke-width: 2;
	}

	:global(.svelte-flow__controls) {
		background: oklch(var(--b2));
		border: 1px solid oklch(var(--b3));
		border-radius: 8px;
	}

	:global(.svelte-flow__controls-button) {
		background: oklch(var(--b2));
		border: none;
		color: oklch(var(--bc));
	}

	:global(.svelte-flow__controls-button:hover) {
		background: oklch(var(--b3));
	}

	:global(.svelte-flow__minimap) {
		border-radius: 8px;
		border: 1px solid oklch(var(--b3));
	}
</style>
