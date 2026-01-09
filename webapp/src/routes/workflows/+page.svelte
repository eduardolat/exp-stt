<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { onMount } from 'svelte';
	import { Plus, Play, Copy, Trash2, ChevronRight, Workflow, Check } from '@lucide/svelte';
	import { Card, PageHeader } from '$lib/components';
	import type { Workflow as WorkflowType } from '$lib/client.gen';
	import { toast } from 'svelte-sonner';

	let isCreating = $state(false);

	onMount(() => {
		store.fetchWorkflows();
	});

	async function createNewWorkflow() {
		isCreating = true;
		try {
			const workflow = await store.createWorkflow({
				name: 'New Workflow',
				enabled: true,
				trigger: { triggerType: 'voice', config: '{}' },
				nodes: [],
				edges: []
			});
			if (workflow) {
				toast.success('Workflow created');
			}
		} finally {
			isCreating = false;
		}
	}

	async function duplicateWorkflow(id: string) {
		const workflow = await store.duplicateWorkflow(id);
		if (workflow) {
			toast.success('Workflow duplicated');
		}
	}

	async function deleteWorkflow(id: string, name: string) {
		if (!confirm(`Delete workflow "${name}"?`)) return;
		try {
			await store.deleteWorkflow(id);
			toast.success('Workflow deleted');
		} catch {
			toast.error('Failed to delete workflow');
		}
	}

	async function setActive(id: string) {
		try {
			await store.setActiveWorkflow(id);
			toast.success('Workflow activated');
		} catch {
			toast.error('Failed to activate workflow');
		}
	}

	function isDefault(workflow: WorkflowType): boolean {
		return workflow.id === 'default';
	}
</script>

<div class="flex flex-col gap-4">
	<PageHeader
		icon={Workflow}
		title="Workflows"
		description="Create and manage custom voice automation workflows"
	/>

	<div class="flex justify-end">
		<button class="btn gap-2 btn-sm btn-primary" onclick={createNewWorkflow} disabled={isCreating}>
			<Plus class="size-4" />
			New Workflow
		</button>
	</div>

	<div class="grid gap-3">
		{#each store.workflows as workflow (workflow.id)}
			<Card class="flex items-center gap-4 p-4">
				<div class="flex-1">
					<div class="flex items-center gap-2">
						<span class="font-medium">{workflow.name}</span>
						{#if isDefault(workflow)}
							<span class="badge badge-sm badge-neutral">Default</span>
						{/if}
						{#if store.activeWorkflowId === workflow.id}
							<span class="badge gap-1 badge-sm badge-primary">
								<Check class="size-3" />
								Active
							</span>
						{/if}
					</div>
					{#if workflow.description}
						<p class="mt-1 text-sm opacity-70">{workflow.description}</p>
					{/if}
					<div class="mt-2 flex gap-4 text-xs opacity-60">
						<span>{workflow.nodes.length} nodes</span>
						<span>{workflow.edges.length} connections</span>
					</div>
				</div>

				<div class="flex items-center gap-2">
					{#if store.activeWorkflowId !== workflow.id}
						<button
							class="btn btn-ghost btn-sm"
							title="Set as active"
							onclick={() => setActive(workflow.id)}
						>
							<Play class="size-4" />
						</button>
					{/if}
					<button
						class="btn btn-ghost btn-sm"
						title="Duplicate"
						onclick={() => duplicateWorkflow(workflow.id)}
					>
						<Copy class="size-4" />
					</button>
					{#if !isDefault(workflow)}
						<button
							class="btn text-error btn-ghost btn-sm"
							title="Delete"
							onclick={() => deleteWorkflow(workflow.id, workflow.name)}
						>
							<Trash2 class="size-4" />
						</button>
					{/if}
					<a href="#/workflows/{workflow.id}" class="btn btn-ghost btn-sm">
						<ChevronRight class="size-4" />
					</a>
				</div>
			</Card>
		{:else}
			<Card class="flex flex-col items-center justify-center py-12 text-center">
				<Workflow class="size-12 opacity-30" />
				<p class="mt-4 text-lg font-medium">No workflows yet</p>
				<p class="mt-1 text-sm opacity-70">Create your first workflow to get started</p>
				<button class="btn btn-primary btn-sm mt-4" onclick={createNewWorkflow}>
					<Plus class="size-4" />
					Create Workflow
				</button>
			</Card>
		{/each}
	</div>
</div>
