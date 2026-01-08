<script lang="ts">
	import { Card, Modal } from '$lib/components';
	import { LoaderCircle, TriangleAlert, RefreshCw, Search } from '@lucide/svelte';

	interface Props {
		baseUrl: string;
		apiKey: string;
		onSelect: (modelId: string) => void;
	}

	interface ModelInfo {
		id: string;
		name?: string;
		owned_by?: string;
	}

	interface ModelsResponse {
		data: ModelInfo[];
	}

	let { baseUrl, apiKey, onSelect }: Props = $props();

	let modal: Modal | undefined = $state();
	let models: ModelInfo[] = $state([]);
	let loading = $state(false);
	let error = $state('');
	let searchQuery = $state('');

	let filteredModels = $derived.by(() => {
		if (!searchQuery.trim()) return models;
		const terms = searchQuery.toLowerCase().split(/\s+/).filter(Boolean);
		return models.filter((model) => {
			const searchText = `${getDisplayName(model)} ${model.id}`.toLowerCase();
			return terms.every((term) => searchText.includes(term));
		});
	});

	async function fetchModels() {
		if (!baseUrl.trim() || !apiKey.trim()) {
			error = 'Base URL and API Key are required';
			return;
		}

		loading = true;
		error = '';
		models = [];
		searchQuery = '';

		try {
			const url = baseUrl.replace(/\/$/, '') + '/models';
			const response = await fetch(url, {
				headers: {
					Authorization: `Bearer ${apiKey}`
				}
			});

			if (!response.ok) {
				throw new Error(`HTTP ${response.status}: ${response.statusText}`);
			}

			const data: ModelsResponse = await response.json();
			models = data.data ?? [];
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to fetch models';
		} finally {
			loading = false;
		}
	}

	export function open() {
		searchQuery = '';
		modal?.open();
		fetchModels();
	}

	function selectModel(id: string) {
		onSelect(id);
		modal?.close();
	}

	function getDisplayName(model: ModelInfo): string {
		if (model.name && model.owned_by) {
			return `${model.name} (${model.owned_by})`;
		}
		if (model.name) {
			return model.name;
		}
		if (model.owned_by) {
			return `${model.id} (${model.owned_by})`;
		}
		return model.id;
	}
</script>

<Modal bind:this={modal} title="Select Model" size="md">
	{#snippet children()}
		<div class="flex h-125 flex-col gap-4">
			<label class="input w-full">
				<Search class="size-4" />
				<input
					type="text"
					class="grow"
					placeholder="Search models..."
					aria-label="Search models"
					bind:value={searchQuery}
				/>
			</label>

			<div class="flex min-h-0 flex-1 flex-col">
				{#if loading}
					<div class="flex flex-1 items-center justify-center">
						<LoaderCircle class="size-6 animate-spin opacity-70" />
					</div>
				{/if}

				{#if error}
					<div class="flex flex-1 flex-col items-center justify-center gap-3 text-center">
						<TriangleAlert class="size-8 text-error" />
						<p class="text-sm text-error">{error}</p>
						<button class="btn btn-outline btn-sm" onclick={fetchModels}>
							<RefreshCw class="size-3.5" />
							Retry
						</button>
					</div>
				{/if}

				{#if !loading && !error && models.length === 0}
					<p class="flex flex-1 items-center justify-center text-sm opacity-70">No models found</p>
				{/if}

				{#if !loading && !error && models.length > 0 && filteredModels.length === 0}
					<p class="flex flex-1 items-center justify-center text-sm opacity-70">
						No models match your search
					</p>
				{/if}

				{#if !loading && !error && filteredModels.length > 0}
					<div class="flex flex-col gap-2 overflow-y-auto">
						{#each filteredModels as model (model.id)}
							<Card interactive darker onclick={() => selectModel(model.id)}>
								<div class="px-3 py-2">
									<p class="font-medium">{getDisplayName(model)}</p>
									<p class="text-xs opacity-70">ID: {model.id}</p>
								</div>
							</Card>
						{/each}
					</div>
				{/if}
			</div>
		</div>
	{/snippet}
	{#snippet actions()}
		<button class="btn" onclick={fetchModels} disabled={loading}>
			<RefreshCw class="size-4 {loading ? 'animate-spin' : ''}" />
			Refresh
		</button>
		<button class="btn" onclick={() => modal?.close()}>Close</button>
	{/snippet}
</Modal>
