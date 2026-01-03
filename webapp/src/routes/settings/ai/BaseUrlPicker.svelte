<script lang="ts">
	import { Card, Modal } from '$lib/components';
	import { Search } from '@lucide/svelte';

	interface Props {
		onSelect: (baseUrl: string) => void;
	}

	interface Provider {
		name: string;
		url: string;
	}

	const PROVIDERS: Provider[] = [
		{ name: 'OpenAI', url: 'https://api.openai.com/v1' },
		{ name: 'Open Router', url: 'https://openrouter.ai/api/v1' },
		{ name: 'Anthropic', url: 'https://api.anthropic.com/v1' },
		{ name: 'X AI', url: 'https://api.x.ai/v1' },
		{ name: 'Groq', url: 'https://api.groq.com/openai/v1' },
		{ name: 'Cerebras', url: 'https://api.cerebras.ai/v1' },
		{ name: 'Mistral', url: 'https://api.mistral.ai/v1' },
		{ name: 'DeepSeek', url: 'https://api.deepseek.com' },
		{ name: 'Ollama (local)', url: 'http://localhost:11434/v1' },
		{ name: 'LM Studio (local)', url: 'http://localhost:1234/v1' }
	];

	let { onSelect }: Props = $props();

	let modal: Modal | undefined = $state();
	let searchQuery = $state('');

	let filteredProviders = $derived.by(() => {
		if (!searchQuery.trim()) return PROVIDERS;
		const terms = searchQuery.toLowerCase().split(/\s+/).filter(Boolean);
		return PROVIDERS.filter((provider) => {
			const searchText = `${provider.name} ${provider.url}`.toLowerCase();
			return terms.every((term) => searchText.includes(term));
		});
	});

	export function open() {
		searchQuery = '';
		modal?.open();
	}

	function selectProvider(url: string) {
		onSelect(url);
		modal?.close();
	}
</script>

<Modal bind:this={modal} title="Select Provider" size="md">
	{#snippet children()}
		<div class="flex h-125 flex-col gap-4">
			<label class="input w-full">
				<Search class="size-4" />
				<input
					type="text"
					class="grow"
					placeholder="Search providers..."
					bind:value={searchQuery}
				/>
			</label>

			<div class="flex min-h-0 flex-1 flex-col">
				{#if filteredProviders.length === 0}
					<p class="flex flex-1 items-center justify-center text-sm opacity-70">
						No providers match your search
					</p>
				{/if}

				{#if filteredProviders.length > 0}
					<div class="flex flex-col gap-2 overflow-y-auto">
						{#each filteredProviders as provider (provider.url)}
							<Card interactive darker onclick={() => selectProvider(provider.url)}>
								<div class="px-3 py-2">
									<p class="font-medium">{provider.name}</p>
									<p class="text-xs opacity-70">{provider.url}</p>
								</div>
							</Card>
						{/each}
					</div>
				{/if}
			</div>
		</div>
	{/snippet}
	{#snippet actions()}
		<button class="btn" onclick={() => modal?.close()}>Close</button>
	{/snippet}
</Modal>
